package buffer

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/dgraph-io/badger/v4"
	"github.com/dgraph-io/ristretto/v2/z"
	"github.com/google/uuid"
)

const (
	eventPrefix = "event:"
	leasePrefix = "_lease:"
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

// Read
// TODO: Consider option to read by EventEnvelope.EntityType
func (b *BadgerBuffer) Read(ctx context.Context, limit int) (envelopes []EventEnvelope, err error) {
	readID := uuid.New().String()

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

			leaseKey := getLeaseKey(envelope)
			if _, err := txn.Get(leaseKey); err == nil {
				continue
			} else if !errors.Is(err, badger.ErrKeyNotFound) {
				return err
			}

			entry := badger.NewEntry(leaseKey, []byte(readID)).WithTTL(30 * time.Second)
			if err := txn.SetEntry(entry); err != nil {
				return fmt.Errorf("failed set lease for envelop %v: %w", envelope.ID, err)
			}

			envelopes = append(envelopes, envelope)
		}

		return nil
	})

	return envelopes, err
}

// Stream
// TODO: Consider to use a channel as a return value instead slice. Relevant for 100k+ events.
func (b *BadgerBuffer) Stream(ctx context.Context) (envelops []EventEnvelope, err error) {
	stream := b.db.NewStream()
	stream.Prefix = []byte(eventPrefix)

	stream.Send = func(buf *z.Buffer) error {
		list, err := badger.BufferToKVList(buf)
		if err != nil {
			return err
		}
		for _, kv := range list.Kv {
			var envelop EventEnvelope
			if err := json.Unmarshal(kv.Value, &envelop); err != nil {
				continue
			}
			envelops = append(envelops, envelop)
		}
		return nil
	}

	if err := stream.Orchestrate(ctx); err != nil {
		return nil, err
	}

	return envelops, nil
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
		// TODO: Remove if use only streams for reading
		if err := wb.Delete(getLeaseKey(envelope)); err != nil {
			return fmt.Errorf("failed to delete lease %s: %w", envelope.ID, err)
		}
	}

	return wb.Flush()
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
			if err := txn.Delete(getLeaseKey(envelope)); err != nil {
				return fmt.Errorf("failed to release envelop %s: %w", envelope.ID, err)
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

func getLeaseKey(envelope EventEnvelope) []byte {
	return append([]byte(leasePrefix), envelope.BufferKey...)
}
