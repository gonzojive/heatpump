// Package db implements a database for storing the history of the state of the
// CX34
package db

import (
	"fmt"
	"sort"
	"strings"
	"time"

	"github.com/dgraph-io/badger/v2"
	"github.com/dgraph-io/badger/v2/options"
	"github.com/golang/glog"
	"github.com/golang/protobuf/proto"
	"github.com/gonzojive/heatpump/proto/chiltrix"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

const (
	schemaVersionKey = "schema-version"
	schemaVersion    = "v2"
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
	db := &Database{bdb}
	if err := db.updateOldVersion(); err != nil {
		return nil, fmt.Errorf("error performing migration: %w", err)
	}
	return db, nil
}

func (db *Database) getSchemaVersion(txn *badger.Txn) (string, error) {
	item, err := txn.Get([]byte(schemaVersion))
	if err == badger.ErrKeyNotFound {
		return "v1", nil
	}
	if err != nil {
		return "", fmt.Errorf("error getting schema version: %v", err)
	}
	val, err := item.ValueCopy(nil)
	if err != nil {
		return "", err
	}
	return string(val), nil
}

func (db *Database) updateOldVersion() error {
	if err := db.badgerDB.Update(func(txn *badger.Txn) error {
		gotVersion, err := db.getSchemaVersion(txn)
		if err != nil {
			return err
		}
		if gotVersion == schemaVersion {
			glog.Infof("schema of database is already up to date (%q)", gotVersion)
			return nil // already up to date
		}
		if err := txn.Set([]byte(schemaVersionKey), []byte(schemaVersion)); err != nil {
			return fmt.Errorf("error setting schema version to %q: %w", schemaVersion, err)
		}
		opts := badger.DefaultIteratorOptions
		it := txn.NewIterator(opts)
		defer it.Close()

		begin := func() {
			it.Rewind()
		}
		for begin(); it.Valid(); it.Next() {
			item := it.Item()
			if !strings.HasPrefix(string(item.Key()), v1KeyPrefix) {
				continue
			}
			err := item.Value(func(v []byte) error {
				state := &chiltrix.State{}
				if err := proto.Unmarshal(v, state); err != nil {
					return err
				}
				if err := db.writeSnapshotInTxn(state, txn); err != nil {
					return err
				}
				return txn.Delete(item.Key())
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

// WriteSnapshot writes a snaptshot of the heatpump state to the database.
func (db *Database) WriteSnapshot(state *chiltrix.State) error {
	if err := db.badgerDB.Update(func(txn *badger.Txn) error {
		return db.writeSnapshotInTxn(state, txn)
	}); err != nil {
		return fmt.Errorf("error writing value to database: %w", err)
	}
	return nil
}

// writeSnapshotInTxn writes a snaptshot of the heatpump state to the database.
func (db *Database) writeSnapshotInTxn(state *chiltrix.State, txn *badger.Txn) error {
	key, err := keyForState(state)
	if err != nil {
		return err
	}
	value, err := proto.Marshal(state)
	if err != nil {
		return err
	}
	if err := txn.Set(key, value); err != nil {
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

const (
	keyTimeLayout = "20060102T150405.999999999"

	v1KeyPrefix = "time/"
)

func keyForTime(t time.Time) []byte {
	//return []byte(fmt.Sprintf("time/%d", t.UnixNano()))
	return []byte(fmt.Sprintf("t/%s", t.In(time.UTC).Format(keyTimeLayout)))
	/*
		key := make([]byte, 8+len("t/"))
		key[0] = 't'
		key[1] = '/'
		binary.LittleEndian.PutUint64(key[2:], uint64(t.UnixNano()))
		return key
	*/
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
	start := time.Now()
	defer func() {
		glog.Infof("QueryStream finished in %s", time.Now().Sub(start))
	}()
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
