package processor

import (
	"context"
	"encoding/json"
	"errors"
	"log/slog"
	"path/filepath"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/reportportal/service-ingest/internal/data/buffer"
	"github.com/reportportal/service-ingest/internal/data/catalog"
	"github.com/reportportal/service-ingest/internal/data/parquet"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type mockBuffer struct {
	sizeFunc    func(ctx context.Context) (buffer.Counter, error)
	readFunc    func(ctx context.Context, limit int) ([]buffer.EventEnvelope, error)
	ackFunc     func(ctx context.Context, events []buffer.EventEnvelope) error
	releaseFunc func(ctx context.Context, events []buffer.EventEnvelope) error
}

func (m *mockBuffer) Put(_ context.Context, _ buffer.EventEnvelope) error {
	return nil
}

func (m *mockBuffer) Size(ctx context.Context) (buffer.Counter, error) {
	return m.sizeFunc(ctx)
}

func (m *mockBuffer) Read(ctx context.Context, limit int) ([]buffer.EventEnvelope, error) {
	return m.readFunc(ctx, limit)
}

func (m *mockBuffer) Ack(ctx context.Context, events []buffer.EventEnvelope) error {
	return m.ackFunc(ctx, events)
}

func (m *mockBuffer) Release(ctx context.Context, events []buffer.EventEnvelope) error {
	return m.releaseFunc(ctx, events)
}

func (m *mockBuffer) Close() error {
	return nil
}

func newTestEnvelope(id string, entityType buffer.EntityType, projectKey, launchUUID string) buffer.EventEnvelope {
	return buffer.EventEnvelope{
		ID:         id,
		ProjectKey: projectKey,
		LaunchUUID: launchUUID,
		EntityUUID: uuid.New().String(),
		EntityType: entityType,
		Operation:  buffer.OperationTypeCreate,
		Timestamp:  time.Date(2026, 3, 18, 12, 0, 0, 0, time.UTC),
		Data:       json.RawMessage(`{"key":"value"}`),
	}
}

func newProcessor(t *testing.T, buf buffer.Buffer) *BatchProcessor {
	t.Helper()
	writer := parquet.NewWriter(t.TempDir(), "snappy")
	return NewBatchProcessor(BatchProcessorOptions{
		Buffer:        buf,
		Writer:        writer,
		FlushInterval: time.Second,
		ReadLimit:     100,
		Logger:        slog.Default(),
	})
}

func TestProcessBatch_CancelledContext(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	buf := &mockBuffer{}
	bp := newProcessor(t, buf)

	err := bp.processBatch(ctx)
	assert.ErrorIs(t, err, context.Canceled)
}

func TestProcessBatch_SizeError(t *testing.T) {
	buf := &mockBuffer{
		sizeFunc: func(ctx context.Context) (buffer.Counter, error) {
			return buffer.Counter{}, errors.New("buffer unavailable")
		},
	}
	bp := newProcessor(t, buf)

	err := bp.processBatch(context.Background())
	assert.ErrorContains(t, err, "failed to get buffer size")
}

func TestProcessBatch_EmptyBuffer(t *testing.T) {
	buf := &mockBuffer{
		sizeFunc: func(ctx context.Context) (buffer.Counter, error) {
			return buffer.Counter{Items: 0}, nil
		},
	}
	bp := newProcessor(t, buf)

	err := bp.processBatch(context.Background())
	assert.NoError(t, err)
}

func TestProcessBatch_ReadError(t *testing.T) {
	buf := &mockBuffer{
		sizeFunc: func(ctx context.Context) (buffer.Counter, error) {
			return buffer.Counter{Items: 5}, nil
		},
		readFunc: func(ctx context.Context, limit int) ([]buffer.EventEnvelope, error) {
			return nil, errors.New("read failed")
		},
	}
	bp := newProcessor(t, buf)

	err := bp.processBatch(context.Background())
	assert.ErrorContains(t, err, "failed to read from buffer")
}

func TestProcessBatch_ReadReturnsEmpty(t *testing.T) {
	buf := &mockBuffer{
		sizeFunc: func(ctx context.Context) (buffer.Counter, error) {
			return buffer.Counter{Items: 5}, nil
		},
		readFunc: func(ctx context.Context, limit int) ([]buffer.EventEnvelope, error) {
			return []buffer.EventEnvelope{}, nil
		},
	}
	bp := newProcessor(t, buf)

	err := bp.processBatch(context.Background())
	assert.NoError(t, err)
}

func TestProcessBatch_Ack(t *testing.T) {
	events := []buffer.EventEnvelope{
		newTestEnvelope("1", buffer.EntityTypeLaunch, "proj1", "launch1"),
		newTestEnvelope("2", buffer.EntityTypeLaunch, "proj1", "launch1"),
	}

	var ackedEvents []buffer.EventEnvelope
	buf := &mockBuffer{
		sizeFunc: func(ctx context.Context) (buffer.Counter, error) {
			return buffer.Counter{Items: int64(len(events))}, nil
		},
		readFunc: func(ctx context.Context, limit int) ([]buffer.EventEnvelope, error) {
			return events, nil
		},
		ackFunc: func(ctx context.Context, events []buffer.EventEnvelope) error {
			ackedEvents = events
			return nil
		},
	}
	bp := newProcessor(t, buf)

	err := bp.processBatch(context.Background())
	require.NoError(t, err)
	assert.Equal(t, events, ackedEvents)
}

func TestProcessBatch_AckError(t *testing.T) {
	events := []buffer.EventEnvelope{
		newTestEnvelope("1", buffer.EntityTypeLaunch, "proj1", "launch1"),
	}

	buf := &mockBuffer{
		sizeFunc: func(ctx context.Context) (buffer.Counter, error) {
			return buffer.Counter{Items: 1}, nil
		},
		readFunc: func(ctx context.Context, limit int) ([]buffer.EventEnvelope, error) {
			return events, nil
		},
		ackFunc: func(ctx context.Context, events []buffer.EventEnvelope) error {
			return errors.New("ack failed")
		},
	}
	bp := newProcessor(t, buf)

	err := bp.processBatch(context.Background())
	assert.ErrorContains(t, err, "failed to ack events")
}

func TestProcessBatch_WritesParquetFiles(t *testing.T) {
	events := []buffer.EventEnvelope{
		newTestEnvelope("1", buffer.EntityTypeLaunch, "proj1", "launch1"),
		newTestEnvelope("2", buffer.EntityTypeItem, "proj1", "launch1"),
		newTestEnvelope("3", buffer.EntityTypeLog, "proj1", "launch1"),
	}

	buf := &mockBuffer{
		sizeFunc: func(ctx context.Context) (buffer.Counter, error) {
			return buffer.Counter{Items: int64(len(events))}, nil
		},
		readFunc: func(ctx context.Context, limit int) ([]buffer.EventEnvelope, error) {
			return events, nil
		},
		ackFunc: func(ctx context.Context, events []buffer.EventEnvelope) error {
			return nil
		},
	}
	bp := newProcessor(t, buf)

	err := bp.processBatch(context.Background())
	require.NoError(t, err)

	for _, event := range events {
		path := filepath.Join(
			bp.writer.BasePath,
			catalog.BuildPath(event.ProjectKey, event.LaunchUUID, string(event.EntityType), event.Timestamp.Format("2006-01-02"), "*"),
		)
		parquetPattern := path + "/*.parquet"
		t.Logf("parquet pattern: %s", parquetPattern)
		matches, err := filepath.Glob(parquetPattern)
		require.NoError(t, err)
		assert.NotEmpty(t, matches, "expected parquet file for entity type %s", event.EntityType)
		t.Logf("parquet files: %v", matches)

		successPattern := path + "/_SUCCESS"
		successMatches, err := filepath.Glob(successPattern)
		require.NoError(t, err)
		assert.NotEmpty(t, successMatches, "expected _SUCCESS marker for entity type %s", event.EntityType)
		t.Logf("success files: %v", successMatches)
	}
}
