// Package db implements a database for storing the history of the state of the
// CX34
package db

import (
	"fmt"
	"sort"
	"time"

	"github.com/dgraph-io/badger/v2"
	"github.com/dgraph-io/badger/v2/options"
	"github.com/golang/protobuf/proto"
	"github.com/gonzojive/heatpump/proto/chiltrix"
)

// Database stores historical values of the heat pump's state.
type Database struct {
	badgerDB *badger.DB
}

// Open returns a databased stored in the given directory.
func Open(dir string) (*Database, error) {
	opts := badger.DefaultOptions(dir)
	opts.ValueLogLoadingMode = options.FileIO
	bdb, err := badger.Open(opts)
	if err != nil {
		return nil, fmt.Errorf("badger database open error: %w", err)
	}
	return &Database{bdb}, nil
}

// WriteSnapshot writes a snaptshot of the heatpump state to the database.
func (db *Database) WriteSnapshot(state *chiltrix.State) error {
	key, err := keyForState(state)
	if err != nil {
		return err
	}
	value, err := proto.Marshal(state)
	if err != nil {
		return err
	}
	if err := db.badgerDB.Update(func(txn *badger.Txn) error {
		return txn.Set(key, value)
	}); err != nil {
		return fmt.Errorf("error writing value to database: %w", err)
	}
	return nil
}

// ReadSnapshots reads a set of snapshots from the database based on a time range.
func (db *Database) ReadSnapshots(start, end time.Time) ([]*chiltrix.State, error) {
	var states []*chiltrix.State

	timeMatches := func(t time.Time) bool {
		return !t.Before(start) && end.After(t)
	}

	prefix := sharedPrefix(keyForTime(start), keyForTime(end))

	if err := db.badgerDB.View(func(txn *badger.Txn) error {
		states = nil

		opts := badger.DefaultIteratorOptions
		opts.Prefix = prefix
		it := txn.NewIterator(opts)
		defer it.Close()
		iterIsAfterEndTime := false
		for it.Rewind(); it.Valid() && !iterIsAfterEndTime; it.Next() {
			item := it.Item()
			err := item.Value(func(v []byte) error {
				state := &chiltrix.State{}
				if err := proto.Unmarshal(v, state); err != nil {
					return err
				}
				iterIsAfterEndTime = !state.GetCollectionTime().AsTime().Before(end)
				if timeMatches(state.GetCollectionTime().AsTime()) {
					states = append(states, state)
				}
				return nil
			})
			if err != nil {
				return err
			}
		}
		return nil
	}); err != nil {
		return nil, fmt.Errorf("error reading values from database: %w", err)
	}
	sort.Slice(states, func(i, j int) bool {
		return states[i].CollectionTime.AsTime().Before(states[j].CollectionTime.AsTime())
	})
	return states, nil
}

func keyForState(state *chiltrix.State) ([]byte, error) {
	if state.GetCollectionTime() == nil {
		return nil, fmt.Errorf("cannot store state with no timestamp")
	}
	return keyForTime(state.GetCollectionTime().AsTime()), nil
}

func keyForTime(t time.Time) []byte {
	return []byte(fmt.Sprintf("time/%d", t.UnixNano()))
}

func sharedPrefix(a, b []byte) []byte {
	var out []byte
	for i := 0; i < len(a) && i < len(b); i++ {
		if a[i] != b[i] {
			break
		}
		out = a[0 : i+1]
	}
	return out
}
