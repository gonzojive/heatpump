package db

import (
	"io/ioutil"
	"testing"
	"time"

	"github.com/gonzojive/heatpump/proto/chiltrix"
	"github.com/google/go-cmp/cmp"
	"google.golang.org/protobuf/testing/protocmp"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func TestDatabase_ReadSnapshots(t *testing.T) {
	t0 := time.Date(2020, time.January, 10, 0, 0, 0, 0, time.UTC)
	t1 := time.Date(2020, time.January, 11, 0, 0, 0, 0, time.UTC)
	t2 := time.Date(2020, time.January, 12, 0, 0, 0, 0, time.UTC)
	t3 := time.Date(2020, time.January, 13, 0, 0, 0, 0, time.UTC)
	t4 := time.Date(2020, time.January, 14, 0, 0, 0, 0, time.UTC)
	tests := []struct {
		name        string
		inputStates []*chiltrix.State
		start       time.Time
		end         time.Time
		want        []*chiltrix.State
		wantErr     bool
	}{
		{
			"one",
			[]*chiltrix.State{
				simpleState(t0, 0),
				simpleState(t1, 1),
				simpleState(t2, 2),
				simpleState(t3, 3),
				simpleState(t4, 4),
			},
			t1,
			t3,
			[]*chiltrix.State{
				simpleState(t1, 1),
				simpleState(t2, 2),
			},
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			dir, err := ioutil.TempDir("", "readsnapshots-*")
			if err != nil {
				t.Fatalf("error setting up test: %v", err)
			}
			db, err := Open(dir)
			if err != nil {
				t.Fatalf("error setting up test: %v", err)
			}
			for _, s := range tt.inputStates {
				if err := db.WriteSnapshot(s); err != nil {
					t.Fatalf("error writing record to database: %v", err)
				}
			}
			got, err := db.ReadSnapshots(tt.start, tt.end)
			if (err != nil) != tt.wantErr {
				t.Errorf("Database.ReadSnapshots() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			diff := cmp.Diff(tt.want, got, protocmp.Transform())
			if diff != "" {
				t.Errorf("Database.ReadSnapshots(%s, %s) got diff (-want, +got):\n%s", tt.start, tt.end, diff)
			}
		})
	}
}

func simpleState(t time.Time, value uint32) *chiltrix.State {
	return &chiltrix.State{
		CollectionTime: timestamppb.New(t),
		RegisterValues: &chiltrix.RegisterSnapshot{
			HoldingRegisterValues: map[uint32]uint32{value: value + 1000},
		},
	}
}
