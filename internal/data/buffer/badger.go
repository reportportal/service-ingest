package buffer

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/dgraph-io/badger/v4"
	"github.com/google/uuid"
)

type BadgerBuffer struct {
	db            *badger.DB
	processorID   string
	leaseDuration time.Duration
}

// NewBadgerBuffer creates a new BadgerBuffer
// If path is empty, uses in-memory mode
func NewBadgerBuffer(path string, leaseDuration time.Duration) (*BadgerBuffer, error) {
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

	if leaseDuration == 0 {
		leaseDuration = 5 * time.Minute
	}

	return &BadgerBuffer{
		db:            db,
		processorID:   uuid.New().String(),
		leaseDuration: leaseDuration,
	}, nil
}

func (b *BadgerBuffer) Put(ctx context.Context, envelope EventEnvelope) error {
	key := fmt.Sprintf("%s/%d-%s",
		envelope.EntityType,
		envelope.Timestamp.UnixNano(),
		envelope.ID,
	)

	data, err := json.Marshal(envelope)
	if err != nil {
		return fmt.Errorf("failed to marshal envelope: %w", err)
	}

	return b.db.Update(func(txn *badger.Txn) error {
		if err := ctx.Err(); err != nil {
			return err
		}
		return txn.Set([]byte(key), data)
	})
}

func (b *BadgerBuffer) Read(ctx context.Context, limit int) ([]EventEnvelope, error) {
	var envelopes []EventEnvelope

	err := b.db.Update(func(txn *badger.Txn) error {
		opts := badger.DefaultIteratorOptions
		opts.PrefetchSize = limit
		it := txn.NewIterator(opts)
		defer it.Close()

		now := time.Now()
		leaseExpiry := now.Add(b.leaseDuration)

		for it.Rewind(); it.Valid() && len(envelopes) < limit; it.Next() {
			item := it.Item()

			// Read value
			var envelope EventEnvelope
			err := item.Value(func(val []byte) error {
				return json.Unmarshal(val, &envelope)
			})
			if err != nil {
				continue // skip malformed entries
			}

			// Check if available
			if !envelope.IsAvailable() {
				continue
			}

			// Set lease
			envelope.LeaseID = b.processorID
			envelope.LeaseExpiresAt = &leaseExpiry

			// Update in BadgerDB
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

func (b *BadgerBuffer) Ack(ctx context.Context, ids []string) error {
	//TODO implement me
	panic("implement me")
}

func (b *BadgerBuffer) Release(ctx context.Context, ids []string) error {
	//TODO implement me
	panic("implement me")
}

func (b *BadgerBuffer) Size(ctx context.Context) (count int, bytes int64, err error) {
	//TODO implement me
	panic("implement me")
}

func (b *BadgerBuffer) Close() error {
	//TODO implement me
	panic("implement me")
}
