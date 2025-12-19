package buffer

import (
	"encoding/binary"
	"errors"
	"fmt"

	"github.com/dgraph-io/badger/v4"
)

var (
	counterKeyCount = []byte("_meta:count")
	counterKeySize  = []byte("_meta:size")
)

func getCounter(txn *badger.Txn, key []byte) (value int64, err error) {
	item, err := txn.Get(key)
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

func updateCounter(txn *badger.Txn, key []byte, delta int64) error {
	current, err := getCounter(txn, key)
	if err != nil {
		return err
	}

	newValue := current + delta

	if newValue < 0 {
		return fmt.Errorf("counter cannot be negative: %d", newValue)
	}

	data := make([]byte, 8)
	binary.BigEndian.PutUint64(data, uint64(newValue))
	return txn.Set(key, data)
}
