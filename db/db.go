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
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
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

	if err := db.readSnapshotsStreaming(start, end, func(s *chiltrix.State) error {
		states = append(states, s)
		return nil
	}, func() error {
		states = nil
		return nil
	}); err != nil {
		return nil, err
	}
	sort.Slice(states, func(i, j int) bool {
		return states[i].CollectionTime.AsTime().Before(states[j].CollectionTime.AsTime())
	})
	return states, nil
}

func (db *Database) readSnapshotsStreaming(start, end time.Time, callback func(*chiltrix.State) error, restart func() error) error {
	timeMatches := func(t time.Time) bool {
		return !t.Before(start) && end.After(t)
	}

	prefix := sharedPrefix(keyForTime(start), keyForTime(end))

	if err := db.badgerDB.View(func(txn *badger.Txn) error {
		if err := restart(); err != nil {
			return fmt.Errorf("callback error: %w", err)
		}

		opts := badger.DefaultIteratorOptions
		opts.Prefix = prefix
		opts.PrefetchValues = true
		opts.PrefetchSize = 1500
		it := txn.NewIterator(opts)
		defer it.Close()
		iterIsAfterEndTime := false

		begin := func() {
			it.Rewind()
			it.Seek(keyForTime(start))
		}
		for begin(); it.Valid() && !iterIsAfterEndTime; it.Next() {
			item := it.Item()
			err := item.Value(func(v []byte) error {
				state := &chiltrix.State{}
				if err := proto.Unmarshal(v, state); err != nil {
					return err
				}
				iterIsAfterEndTime = !state.GetCollectionTime().AsTime().Before(end)
				if timeMatches(state.GetCollectionTime().AsTime()) {
					if err := callback(state); err != nil {
						return fmt.Errorf("callback error: %w", err)
					}
				}
				return nil
			})
			if err != nil {
				return err
			}
		}
		return nil
	}); err != nil {
		return fmt.Errorf("error reading values from database: %w", err)
	}
	return nil
}

// HistorianService returns an implementation of chiltrix.Historian.
func (db *Database) HistorianService() *Service {
	return &Service{db: db}
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

// Service implements chiltrix.HistorianServer and is backed by a Database object.
type Service struct {
	chiltrix.UnimplementedHistorianServer

	db *Database
}

var _ chiltrix.HistorianServer = (*Service)(nil)

// QueryStream returns a stream of State values based on a query.
func (s *Service) QueryStream(req *chiltrix.QueryStreamRequest, srv chiltrix.Historian_QueryStreamServer) error {
	deliveredKeys := map[int64]bool{}

	err := s.db.readSnapshotsStreaming(req.StartTime.AsTime(), req.GetEndTime().AsTime(), func(s *chiltrix.State) error {
		key := s.GetCollectionTime().AsTime().UnixNano()
		if deliveredKeys[key] {
			return nil
		}
		deliveredKeys[key] = true

		return srv.Send(&chiltrix.QueryStreamResponse{
			State: s,
		})
	}, func() error {
		return nil
	})

	if err != nil {
		return status.Errorf(codes.Internal, "database retrieval error: %v", err)
	}
	return nil
}
