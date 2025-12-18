package buffer

import "context"

type Buffer interface {
	Put(ctx context.Context, envelope EventEnvelope) error
	Read(ctx context.Context, limit int) ([]EventEnvelope, error)
	Ack(ctx context.Context, ids []string) error
	Release(ctx context.Context, ids []string) error
	Size(ctx context.Context) (count int, bytes int64, err error)
	Close() error
}
