package buffer

import "context"

type Buffer interface {
	Put(ctx context.Context, envelope EventEnvelope) error
	Read(ctx context.Context, limit int) ([]EventEnvelope, error)
	Stream(ctx context.Context) (<-chan EventEnvelope, <-chan error)
	Ack(ctx context.Context, events []EventEnvelope) error
	Release(ctx context.Context, events []EventEnvelope) error
	Size(ctx context.Context) (int, error)
	Close() error
}
