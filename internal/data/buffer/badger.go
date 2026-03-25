package buffer

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"

	"github.com/dgraph-io/badger/v4"
	"github.com/google/uuid"
)

type BadgerBuffer struct {
	db *badger.DB
}

// NewBadgerBuffer creates a new BadgerBuffer
// If path is empty, uses in-memory mode
func NewBadgerBuffer(path string) (*BadgerBuffer, error) {
	var opts badger.Options

	if path == "" {
		opts = badger.DefaultOptions("").WithInMemory(true)
	} else {
		opts = badger.DefaultOptions(path)
	}

	db, err := badger.Open(opts)
	if err != nil {
		return nil, fmt.Errorf("failed to open badger: %w", err)
	}

	return &BadgerBuffer{
		db: db,
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

// TODO: Consider option to read by EventEnvelope.EntityType
func (b *BadgerBuffer) Read(ctx context.Context, limit int) (envelopes []EventEnvelope, err error) {
	readID := uuid.New().String()

	err = b.db.Update(func(txn *badger.Txn) error {
		if err := ctx.Err(); err != nil {
			return err
		}

		opts := badger.DefaultIteratorOptions
		opts.PrefetchSize = limit
		opts.Prefix = []byte("event:")
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

			if !envelope.IsAvailable() {
				continue
			}

			envelope.LeaseID = readID

			data, err := json.Marshal(envelope)
			if err != nil {
				return fmt.Errorf("failed to marshal envelope: %w", err)
			}

			if err := txn.Set(item.Key(), data); err != nil {
				return fmt.Errorf("failed to set lease: %w", err)
			}

			envelopes = append(envelopes, envelope)
		}

		return nil
	})

	return envelopes, err
}

func (b *BadgerBuffer) Ack(ctx context.Context, events []EventEnvelope) error {
	if len(events) == 0 {
		return nil
	}

	return b.db.Update(func(txn *badger.Txn) error {
		if err := ctx.Err(); err != nil {
			return err
		}

		for _, envelope := range events {
			key := envelope.BufferKey
			if _, err := txn.Get(key); err != nil {
				slog.Error("ack key not found, skipping delete",
					"envelope_id", envelope.ID,
					"key", string(key),
				)
				continue
			}
			if err := txn.Delete(key); err != nil {
				return fmt.Errorf("failed to delete envelope %s: %w", envelope.ID, err)
			}
		}

		return nil
	})
}

func (b *BadgerBuffer) Release(ctx context.Context, events []EventEnvelope) error {
	if len(events) == 0 {
		return nil
	}

	return b.db.Update(func(txn *badger.Txn) error {
		if err := ctx.Err(); err != nil {
			return err
		}

		for _, envelope := range events {
			envelope.LeaseID = ""
			//envelope.LeaseExpiresAt = nil

			data, err := json.Marshal(envelope)
			if err != nil {
				return fmt.Errorf("failed to marshal envelope: %w", err)
			}

			key := envelope.BufferKey

			if err := txn.Set(key, data); err != nil {
				return fmt.Errorf("failed to release envelope %s: %w", envelope.ID, err)
			}
		}

		return nil
	})
}

func (b *BadgerBuffer) Size(ctx context.Context) (counter int, err error) {
	err = b.db.View(func(txn *badger.Txn) error {
		if ctxErr := ctx.Err(); ctxErr != nil {
			return ctxErr
		}

		opts := badger.DefaultIteratorOptions
		opts.Prefix = []byte("event:")
		opts.PrefetchValues = false
		it := txn.NewIterator(opts)
		defer it.Close()
		for it.Rewind(); it.Valid(); it.Next() {
			counter = counter + 1
		}
		return nil
	})

	return counter, err
}

func (b *BadgerBuffer) Close() error {
	return b.db.Close()
}

func buildKey(envelope EventEnvelope) []byte {
	return fmt.Appendf(nil, "event:%s:%020d:%s",
		envelope.EntityType,
		envelope.Timestamp.UnixNano(),
		envelope.ID,
	)
}
