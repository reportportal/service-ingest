package buffer

import (
	"context"
	"encoding/json"
	"strconv"
	"testing"
	"time"
)

func newTestBuffer(t *testing.T) *BadgerBuffer {
	t.Helper()
	buf, err := NewBadgerBuffer("")
	if err != nil {
		t.Fatalf("failed to create buffer: %v", err)
	}
	t.Cleanup(func() {
		if err := buf.Close(); err != nil {
			t.Errorf("failed to close buffer: %v", err)
		}
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
	if err := buf.Put(ctx, e); err != nil {
		t.Fatalf("Put(%s): %v", e.ID, err)
	}
}

func TestPut_IncreasesSize(t *testing.T) {
	buf := newTestBuffer(t)
	ctx := context.Background()

	mustPut(t, buf, ctx, newEnvelope("1", EntityTypeLaunch))

	counter, err := buf.Size(ctx)
	if err != nil {
		t.Fatalf("Size: %v", err)
	}
	if counter.Items != 1 {
		t.Errorf("expected 1 item, got %d", counter.Items)
	}
}

func TestRead_ReturnsEnvelopes(t *testing.T) {
	buf := newTestBuffer(t)
	ctx := context.Background()

	mustPut(t, buf, ctx, newEnvelope("1", EntityTypeLaunch))
	mustPut(t, buf, ctx, newEnvelope("2", EntityTypeItem))

	envelopes, err := buf.Read(ctx, 10)
	if err != nil {
		t.Fatalf("Read: %v", err)
	}
	if len(envelopes) != 2 {
		t.Errorf("expected 2 envelopes, got %d", len(envelopes))
	}
}

func TestRead_RespectsLimit(t *testing.T) {
	buf := newTestBuffer(t)
	ctx := context.Background()

	for i := range 5 {
		mustPut(t, buf, ctx, newEnvelope(strconv.Itoa(i), EntityTypeLaunch))
	}

	envelopes, err := buf.Read(ctx, 3)
	if err != nil {
		t.Fatalf("Read: %v", err)
	}
	if len(envelopes) != 3 {
		t.Errorf("expected 3 envelopes, got %d", len(envelopes))
	}
}

func TestRead_SetsLease(t *testing.T) {
	buf := newTestBuffer(t)
	ctx := context.Background()

	mustPut(t, buf, ctx, newEnvelope("1", EntityTypeLaunch))

	envelopes, err := buf.Read(ctx, 10)
	if err != nil {
		t.Fatalf("Read: %v", err)
	}
	if envelopes[0].LeaseID == "" {
		t.Error("expected LeaseID to be set after Read")
	}
}

func TestRead_DoesNotReturnLeasedItems(t *testing.T) {
	buf := newTestBuffer(t)
	ctx := context.Background()

	mustPut(t, buf, ctx, newEnvelope("1", EntityTypeLaunch))

	first, err := buf.Read(ctx, 10)
	if err != nil {
		t.Fatalf("first Read: %v", err)
	}
	if len(first) != 1 {
		t.Fatalf("expected 1 envelope, got %d", len(first))
	}

	second, err := buf.Read(ctx, 10)
	if err != nil {
		t.Fatalf("second Read: %v", err)
	}
	if len(second) != 0 {
		t.Errorf("expected 0 envelopes, got %d", len(second))
	}
}

func TestAck_DeletesItem(t *testing.T) {
	buf := newTestBuffer(t)
	ctx := context.Background()

	mustPut(t, buf, ctx, newEnvelope("1", EntityTypeLaunch))

	envelopes, err := buf.Read(ctx, 10)
	if err != nil {
		t.Fatalf("Read: %v", err)
	}

	if err := buf.Ack(ctx, envelopes); err != nil {
		t.Fatalf("Ack: %v", err)
	}

	counter, err := buf.Size(ctx)
	if err != nil {
		t.Fatalf("Size: %v", err)
	}
	if counter.Items != 0 {
		t.Errorf("expected 0 items after Ack, got %d", counter.Items)
	}
}

func TestRelease_MakesItemAvailableAgain(t *testing.T) {
	buf := newTestBuffer(t)
	ctx := context.Background()

	mustPut(t, buf, ctx, newEnvelope("1", EntityTypeLaunch))

	envelopes, err := buf.Read(ctx, 10)
	if err != nil {
		t.Fatalf("Read: %v", err)
	}

	if err := buf.Release(ctx, envelopes); err != nil {
		t.Fatalf("Release: %v", err)
	}

	envelopes2, err := buf.Read(ctx, 10)
	if err != nil {
		t.Fatalf("Read after Release: %v", err)
	}
	if len(envelopes2) != 1 {
		t.Errorf("expected 1 envelope after Release, got %d", len(envelopes2))
	}
}

func TestSize_ReflectsOnlyUnacked(t *testing.T) {
	buf := newTestBuffer(t)
	ctx := context.Background()

	mustPut(t, buf, ctx, newEnvelope("1", EntityTypeLaunch))
	mustPut(t, buf, ctx, newEnvelope("2", EntityTypeLaunch))

	envelopes, err := buf.Read(ctx, 10)
	if err != nil {
		t.Fatalf("Read: %v", err)
	}

	if err := buf.Ack(ctx, envelopes[:1]); err != nil {
		t.Fatalf("Ack: %v", err)
	}

	counter, err := buf.Size(ctx)
	if err != nil {
		t.Fatalf("Size: %v", err)
	}
	if counter.Items != 1 {
		t.Errorf("expected 1 item after partial Ack, got %d", counter.Items)
	}
}

func TestRead_EmptyBuffer(t *testing.T) {
	buf := newTestBuffer(t)
	ctx := context.Background()

	envelopes, err := buf.Read(ctx, 10)
	if err != nil {
		t.Fatalf("Read on empty buffer: %v", err)
	}
	if len(envelopes) != 0 {
		t.Errorf("expected 0 envelopes, got %d", len(envelopes))
	}
}

func TestAck_EmptySlice(t *testing.T) {
	buf := newTestBuffer(t)
	ctx := context.Background()

	if err := buf.Ack(ctx, nil); err != nil {
		t.Errorf("Ack with nil should be no-op, got: %v", err)
	}
}
