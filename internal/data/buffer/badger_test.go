package buffer

import (
	"context"
	"encoding/json"
	"strconv"
	"testing"
	"time"

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

func newEnvelope(id string, entityType EntityType) EventEnvelope {
	return EventEnvelope{
		ID:         id,
		ProjectKey: "project1",
		LaunchUUID: "launch1",
		EntityUUID: "entity1",
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

	mustPut(t, buf, ctx, newEnvelope("1", EntityTypeLaunch))

	counter, err := buf.Size(ctx)
	require.NoError(t, err)
	assert.Equal(t, int64(1), counter.Items)
}

func TestRead_ReturnsEnvelopes(t *testing.T) {
	buf := newTestBuffer(t)
	ctx := context.Background()

	mustPut(t, buf, ctx, newEnvelope("1", EntityTypeLaunch))
	mustPut(t, buf, ctx, newEnvelope("2", EntityTypeItem))

	envelopes, err := buf.Read(ctx, 10)
	require.NoError(t, err)
	assert.Len(t, envelopes, 2)
}

func TestRead_RespectsLimit(t *testing.T) {
	buf := newTestBuffer(t)
	ctx := context.Background()

	for i := range 5 {
		mustPut(t, buf, ctx, newEnvelope(strconv.Itoa(i), EntityTypeLaunch))
	}

	envelopes, err := buf.Read(ctx, 3)
	require.NoError(t, err)
	assert.Len(t, envelopes, 3)
}

func TestRead_SetsLease(t *testing.T) {
	buf := newTestBuffer(t)
	ctx := context.Background()

	mustPut(t, buf, ctx, newEnvelope("1", EntityTypeLaunch))

	envelopes, err := buf.Read(ctx, 10)
	require.NoError(t, err)
	assert.NotEmpty(t, envelopes[0].LeaseID)
}

func TestRead_DoesNotReturnLeasedItems(t *testing.T) {
	buf := newTestBuffer(t)
	ctx := context.Background()

	mustPut(t, buf, ctx, newEnvelope("1", EntityTypeLaunch))

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

	mustPut(t, buf, ctx, newEnvelope("1", EntityTypeLaunch))

	envelopes, err := buf.Read(ctx, 10)
	require.NoError(t, err)

	require.NoError(t, buf.Ack(ctx, envelopes))

	counter, err := buf.Size(ctx)
	require.NoError(t, err)
	assert.Equal(t, int64(0), counter.Items)
}

func TestRelease_MakesItemAvailableAgain(t *testing.T) {
	buf := newTestBuffer(t)
	ctx := context.Background()

	mustPut(t, buf, ctx, newEnvelope("1", EntityTypeLaunch))

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

	mustPut(t, buf, ctx, newEnvelope("1", EntityTypeLaunch))
	mustPut(t, buf, ctx, newEnvelope("2", EntityTypeLaunch))

	envelopes, err := buf.Read(ctx, 10)
	require.NoError(t, err)

	require.NoError(t, buf.Ack(ctx, envelopes[:1]))

	counter, err := buf.Size(ctx)
	require.NoError(t, err)
	assert.Equal(t, int64(1), counter.Items)
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
