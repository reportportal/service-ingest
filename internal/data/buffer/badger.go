package buffer

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"

	"github.com/dgraph-io/badger/v4"
	"github.com/dgraph-io/ristretto/v2/z"
)

const eventPrefix = "event:"

type BadgerBuffer struct {
	db     *badger.DB
	logger *slog.Logger
}

// NewBadgerBuffer creates a new BadgerBuffer
// If path is empty, uses in-memory mode
func NewBadgerBuffer(opts badger.Options) (*BadgerBuffer, error) {
	db, err := badger.Open(opts)
	if err != nil {
		return nil, fmt.Errorf("failed to open badger: %w", err)
	}

	return &BadgerBuffer{
		db:     db,
		logger: slog.Default(),
	}, nil
}

func (b *BadgerBuffer) Put(ctx context.Context, envelope EventEnvelope) error {
	envelope.BufferKey = buildKey(envelope)

	data, err := json.Marshal(envelope)
	if err != nil {
		return fmt.Errorf("failed to marshal envelope: %w", err)
	}

	return b.db.Update(func(txn *badger.Txn) error {
		if err := ctx.Err(); err != nil {
			return err
		}

		if err := txn.Set(envelope.BufferKey, data); err != nil {
			return err
		}

		return nil
	})
}

// Read
// TODO: Consider option to read by EventEnvelope.EntityType
func (b *BadgerBuffer) Read(ctx context.Context, limit int) (envelopes []EventEnvelope, err error) {
	err = b.db.Update(func(txn *badger.Txn) error {
		if err := ctx.Err(); err != nil {
			return err
		}

		opts := badger.DefaultIteratorOptions
		opts.PrefetchSize = limit
		opts.Prefix = []byte(eventPrefix)
		it := txn.NewIterator(opts)
		defer it.Close()

		for it.Rewind(); it.Valid() && len(envelopes) < limit; it.Next() {
			item := it.Item()

			var envelope EventEnvelope
			if err := item.Value(func(val []byte) error {
				return json.Unmarshal(val, &envelope)
			}); err != nil {
				continue
			}

			envelopes = append(envelopes, envelope)
		}

		return nil
	})

	return envelopes, err
}

func (b *BadgerBuffer) Stream(ctx context.Context) (<-chan EventEnvelope, <-chan error) {
	ch := make(chan EventEnvelope, 1024)
	errCh := make(chan error, 1)

	stream := b.db.NewStream()
	stream.Prefix = []byte(eventPrefix)

	stream.Send = func(buf *z.Buffer) error {
		list, err := badger.BufferToKVList(buf)
		if err != nil {
			return err
		}
		for _, kv := range list.Kv {
			var envelope EventEnvelope
			if err := json.Unmarshal(kv.Value, &envelope); err != nil {
				b.logger.Warn("can't convert entry to envelope", "err", err)
				continue
			}
			select {
			case ch <- envelope:
			case <-ctx.Done():
				return ctx.Err()
			}
		}
		return nil
	}

	go func() {
		defer close(ch)
		if err := stream.Orchestrate(ctx); err != nil {
			errCh <- err
		}
		close(errCh)
	}()

	return ch, errCh
}

func (b *BadgerBuffer) Ack(ctx context.Context, events []EventEnvelope) error {
	if len(events) == 0 {
		return nil
	}

	wb := b.db.NewWriteBatch()
	defer wb.Cancel()

	for _, envelope := range events {
		if err := ctx.Err(); err != nil {
			return err
		}
		if err := wb.Delete(envelope.BufferKey); err != nil {
			return fmt.Errorf("failed to delete envelope %s: %w", envelope.ID, err)
		}
	}

	return wb.Flush()
}

func (b *BadgerBuffer) Size(ctx context.Context) (counter int, err error) {
	err = b.db.View(func(txn *badger.Txn) error {
		if ctxErr := ctx.Err(); ctxErr != nil {
			return ctxErr
		}

		opts := badger.DefaultIteratorOptions
		opts.Prefix = []byte(eventPrefix)
		opts.PrefetchValues = false
		it := txn.NewIterator(opts)
		defer it.Close()
		for it.Rewind(); it.Valid(); it.Next() {
			counter++
		}
		return nil
	})

	return counter, err
}

func (b *BadgerBuffer) Close() error {
	return b.db.Close()
}

func buildKey(envelope EventEnvelope) []byte {
	return fmt.Appendf(nil, eventPrefix+"%s:%020d:%s",
		envelope.EntityType,
		envelope.Timestamp.UnixNano(),
		envelope.ID,
	)
}
