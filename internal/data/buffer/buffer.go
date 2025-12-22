package buffer

import "context"

type Buffer interface {
	Put(ctx context.Context, envelope EventEnvelope) error
	Read(ctx context.Context, limit int) ([]EventEnvelope, error)
	Ack(ctx context.Context, events []EventEnvelope) error
	Release(ctx context.Context, events []EventEnvelope) error
	Size(ctx context.Context) (Counter, error)
	Close() error
}

type Counter struct {
	Items int64
}
