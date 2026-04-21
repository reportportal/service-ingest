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
	opts := badger.DefaultOptions("").WithInMemory(true)
	buf, err := NewBadgerBuffer(opts)
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

func collectStream(t *testing.T, ch <-chan EventEnvelope, errCh <-chan error) []EventEnvelope {
	t.Helper()
	var envelopes []EventEnvelope
	for e := range ch {
		envelopes = append(envelopes, e)
	}
	require.NoError(t, <-errCh)
	return envelopes
}

func TestStream_ReturnsAllEnvelopes(t *testing.T) {
	buf := newTestBuffer(t)
	ctx := context.Background()

	mustPut(t, buf, ctx, newEnvelope(EntityTypeLaunch))
	mustPut(t, buf, ctx, newEnvelope(EntityTypeItem))
	mustPut(t, buf, ctx, newEnvelope(EntityTypeLog))

	ch, errCh := buf.Stream(ctx)
	envelopes := collectStream(t, ch, errCh)
	assert.Len(t, envelopes, 3)
}

func TestStream_EmptyBuffer(t *testing.T) {
	buf := newTestBuffer(t)
	ctx := context.Background()

	ch, errCh := buf.Stream(ctx)
	envelopes := collectStream(t, ch, errCh)
	assert.Empty(t, envelopes)
}

func TestStream_PreservesData(t *testing.T) {
	buf := newTestBuffer(t)
	ctx := context.Background()

	original := newEnvelope(EntityTypeLaunch)
	original.Data = json.RawMessage(`{"name":"test-launch"}`)
	mustPut(t, buf, ctx, original)

	ch, errCh := buf.Stream(ctx)
	envelopes := collectStream(t, ch, errCh)
	require.Len(t, envelopes, 1)

	assert.Equal(t, original.ID, envelopes[0].ID)
	assert.Equal(t, original.ProjectKey, envelopes[0].ProjectKey)
	assert.Equal(t, original.LaunchUUID, envelopes[0].LaunchUUID)
	assert.Equal(t, original.EntityType, envelopes[0].EntityType)
	assert.JSONEq(t, `{"name":"test-launch"}`, string(envelopes[0].Data))
}

func TestLargeVolume_ReadWithAck(t *testing.T) {
	buf := newTestBuffer(t)
	ctx := context.Background()

	const totalPuts = 100000
	const readLimit = 10000

	for range totalPuts {
		e := newEnvelope(EntityTypeItem)
		require.NoError(t, buf.Put(ctx, e))
	}

	size, err := buf.Size(ctx)
	require.NoError(t, err)
	t.Logf("after puts: size=%d", size)

	totalRead := 0

	for {
		events, err := buf.Read(ctx, readLimit)
		require.NoError(t, err)
		if len(events) == 0 {
			break
		}
		require.NoError(t, buf.Ack(ctx, events))
		totalRead += len(events)
	}

	size, err = buf.Size(ctx)
	require.NoError(t, err)

	assert.Equal(t, totalPuts, totalRead, "should read all events")
	assert.Equal(t, 0, size, "buffer should be empty")
}

func TestLargeVolume_StreamWithAck(t *testing.T) {
	buf := newTestBuffer(t)
	ctx := context.Background()

	const total = 100000
	for range total {
		mustPut(t, buf, ctx, newEnvelope(EntityTypeLog))
	}

	ch, errCh := buf.Stream(ctx)
	envelopes := collectStream(t, ch, errCh)
	assert.Len(t, envelopes, total)

	require.NoError(t, buf.Ack(ctx, envelopes))
	size, err := buf.Size(ctx)
	require.NoError(t, err)
	assert.Equal(t, 0, size)
}
