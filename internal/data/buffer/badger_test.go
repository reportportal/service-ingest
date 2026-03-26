package buffer

import (
	"context"
	"encoding/json"
	"testing"
	"time"

	"github.com/dgraph-io/badger/v4"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func newTestBuffer(t *testing.T) *BadgerBuffer {
	t.Helper()
	buf, err := NewBadgerBuffer("")
	require.NoError(t, err)
	t.Cleanup(func() {
		assert.NoError(t, buf.Close())
	})
	return buf
}

func newEnvelope(entityType EntityType) EventEnvelope {
	return EventEnvelope{
		ID:         uuid.New().String(),
		ProjectKey: "project1",
		LaunchUUID: "bc2d0f53-5041-46e8-a14c-267875a49f0c",
		EntityUUID: uuid.New().String(),
		EntityType: entityType,
		Operation:  OperationTypeCreate,
		Timestamp:  time.Now(),
		Data:       json.RawMessage(`{}`),
	}
}

func mustPut(t *testing.T, buf *BadgerBuffer, ctx context.Context, e EventEnvelope) {
	t.Helper()
	require.NoError(t, buf.Put(ctx, e))
}

func TestPut_IncreasesSize(t *testing.T) {
	buf := newTestBuffer(t)
	ctx := context.Background()

	mustPut(t, buf, ctx, newEnvelope(EntityTypeLaunch))

	counter, err := buf.Size(ctx)
	require.NoError(t, err)
	assert.Equal(t, 1, counter)
}

func TestRead_ReturnsEnvelopes(t *testing.T) {
	buf := newTestBuffer(t)
	ctx := context.Background()

	mustPut(t, buf, ctx, newEnvelope(EntityTypeLaunch))
	mustPut(t, buf, ctx, newEnvelope(EntityTypeItem))

	envelopes, err := buf.Read(ctx, 10)
	require.NoError(t, err)
	assert.Len(t, envelopes, 2)
}

func TestRead_RespectsLimit(t *testing.T) {
	buf := newTestBuffer(t)
	ctx := context.Background()

	for range 5 {
		mustPut(t, buf, ctx, newEnvelope(EntityTypeLaunch))
	}

	envelopes, err := buf.Read(ctx, 3)
	require.NoError(t, err)
	assert.Len(t, envelopes, 3)
}

func TestRead_DoesNotReturnLeasedItems(t *testing.T) {
	buf := newTestBuffer(t)
	ctx := context.Background()

	mustPut(t, buf, ctx, newEnvelope(EntityTypeLaunch))

	first, err := buf.Read(ctx, 10)
	require.NoError(t, err)
	require.Len(t, first, 1)

	second, err := buf.Read(ctx, 10)
	require.NoError(t, err)
	assert.Empty(t, second)
}

func TestAck_DeletesItem(t *testing.T) {
	buf := newTestBuffer(t)
	ctx := context.Background()

	mustPut(t, buf, ctx, newEnvelope(EntityTypeLaunch))

	envelopes, err := buf.Read(ctx, 10)
	require.NoError(t, err)

	require.NoError(t, buf.Ack(ctx, envelopes))

	counter, err := buf.Size(ctx)
	require.NoError(t, err)
	assert.Equal(t, 0, counter)
}

func TestRelease_MakesItemAvailableAgain(t *testing.T) {
	buf := newTestBuffer(t)
	ctx := context.Background()

	mustPut(t, buf, ctx, newEnvelope(EntityTypeLaunch))

	envelopes, err := buf.Read(ctx, 10)
	require.NoError(t, err)

	require.NoError(t, buf.Release(ctx, envelopes))

	envelopes2, err := buf.Read(ctx, 10)
	require.NoError(t, err)
	assert.Len(t, envelopes2, 1)
}

func TestSize_ReflectsOnlyUnacked(t *testing.T) {
	buf := newTestBuffer(t)
	ctx := context.Background()

	mustPut(t, buf, ctx, newEnvelope(EntityTypeLaunch))
	mustPut(t, buf, ctx, newEnvelope(EntityTypeLaunch))

	envelopes, err := buf.Read(ctx, 10)
	require.NoError(t, err)

	require.NoError(t, buf.Ack(ctx, envelopes[:1]))

	counter, err := buf.Size(ctx)
	require.NoError(t, err)
	assert.Equal(t, 1, counter)
}

func TestRead_EmptyBuffer(t *testing.T) {
	buf := newTestBuffer(t)
	ctx := context.Background()

	envelopes, err := buf.Read(ctx, 10)
	require.NoError(t, err)
	assert.Empty(t, envelopes)
}

func TestAck_EmptySlice(t *testing.T) {
	buf := newTestBuffer(t)
	ctx := context.Background()

	assert.NoError(t, buf.Ack(ctx, nil))
}

func TestPut_ReadAck(t *testing.T) {
	buf := newTestBuffer(t)
	ctx := context.Background()

	const totalPuts = 40000
	const readLimit = 3000

	// Sequential puts - no concurrency
	for range totalPuts {
		e := newEnvelope(EntityTypeItem)
		require.NoError(t, buf.Put(ctx, e))
	}

	size, err := buf.Size(ctx)
	require.NoError(t, err)
	t.Logf("after puts: size=%d", size)

	// Read and Ack in batches
	totalAcked := 0
	totalRead := 0
	iteration := 0
	for {
		events, err := buf.Read(ctx, readLimit)
		totalRead += len(events)
		require.NoError(t, err)

		if len(events) == 0 {
			break
		}

		require.NoError(t, buf.Ack(ctx, events))
		totalAcked += len(events)
		iterSize, err := buf.Size(ctx)
		require.NoError(t, err)
		iteration++
		t.Logf("iteration %d: read=%d, totalRead=%d, totalAcked=%d, buffeSize=%d",
			iteration,
			len(events),
			totalRead,
			totalAcked,
			iterSize,
		)
	}

	size, err = buf.Size(ctx)
	require.NoError(t, err)

	if size > 0 {
		// Check what's stuck: leased or not?
		leased := 0
		unleased := 0
		_ = buf.db.View(func(txn *badger.Txn) error {
			opts := badger.DefaultIteratorOptions
			opts.Prefix = []byte("event:")
			it := txn.NewIterator(opts)
			defer it.Close()
			for it.Rewind(); it.Valid(); it.Next() {
				var env EventEnvelope
				_ = it.Item().Value(func(val []byte) error {
					return json.Unmarshal(val, &env)
				})
				_, err := txn.Get(getLeaseKey(env))
				if err == nil {
					leased++
				} else {
					unleased++
				}
			}
			return nil
		})
		t.Errorf("stuck events: total=%d, leased=%d, unleased=%d", size, leased, unleased)
	}

	assert.Equal(t, totalPuts, totalAcked, "should ack all events")
	assert.Equal(t, 0, size, "buffer should be empty")
}

func TestConcurrentPut_WhileReadAck(t *testing.T) {
	buf := newTestBuffer(t)
	ctx := context.Background()

	const totalPuts = 40000
	const readLimit = 3000

	// Concurrent puts while reading
	putDone := make(chan struct{})
	go func() {
		defer close(putDone)
		for range totalPuts {
			e := newEnvelope(EntityTypeItem)
			_ = buf.Put(ctx, e)
		}
	}()

	// Read and Ack concurrently with puts
	totalAcked := 0

	for {
		events, err := buf.Read(ctx, readLimit)
		require.NoError(t, err)

		if len(events) == 0 {
			select {
			case <-putDone:
				// Puts finished, do one final read
				events, err = buf.Read(ctx, readLimit)
				require.NoError(t, err)
				if len(events) == 0 {
					goto done
				}
			default:
				// Puts still running, keep trying
				continue
			}
		}

		require.NoError(t, buf.Ack(ctx, events))
		totalAcked += len(events)
	}

done:
	// Drain any remaining events
	for {
		events, err := buf.Read(ctx, readLimit)
		require.NoError(t, err)
		if len(events) == 0 {
			break
		}
		require.NoError(t, buf.Ack(ctx, events))
		totalAcked += len(events)
	}

	size, err := buf.Size(ctx)
	require.NoError(t, err)
	assert.Equal(t, 0, size, "buffer should be empty, got %d stuck events", size)
	t.Logf("total acked: %d out of %d puts", totalAcked, totalPuts)
}
