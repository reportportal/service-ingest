package buffer

import (
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

func (b *BadgerBuffer) Close() error {
	return b.db.Close()
}
