package buffer

import (
	"encoding/binary"
	"errors"
	"fmt"

	"github.com/dgraph-io/badger/v4"
)

var (
	counterKeyItems = []byte("_meta:items")
)

func getCounter(txn *badger.Txn) (value int64, err error) {
	item, err := txn.Get(counterKeyItems)
	if err != nil {
		if errors.Is(err, badger.ErrKeyNotFound) {
			return 0, nil
		}
		return 0, err
	}

	err = item.Value(func(val []byte) error {
		if len(val) != 8 {
			return fmt.Errorf("invalid counter value size: %d", len(val))
		}
		value = int64(binary.BigEndian.Uint64(val))
		return nil
	})

	return value, err
}

func updateCounter(txn *badger.Txn, delta int64) error {
	current, err := getCounter(txn)
	if err != nil {
		return err
	}

	newValue := current + delta

	if newValue < 0 {
		return fmt.Errorf("counter cannot be negative: %d", newValue)
	}

	data := make([]byte, 8)
	binary.BigEndian.PutUint64(data, uint64(newValue))
	return txn.Set(counterKeyItems, data)
}
