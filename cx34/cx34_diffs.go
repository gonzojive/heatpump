package cx34

import (
	"fmt"

	"github.com/golang/protobuf/proto"
	"github.com/inhies/go-bytesize"

	"github.com/gonzojive/heatpump/proto/chiltrix"
)

func DebugSequenceInfo(states []*State) string {
	var protos []*chiltrix.State
	for _, x := range states {
		protos = append(protos, x.Proto())
	}
	totalSizeA := 0
	for _, x := range protos {
		totalSizeA += proto.Size(x)
	}
	asSeq, err := encodeSeq(protos)
	if err != nil {
		return fmt.Sprintf("error computing differential encoding: %v", err)
	}
	return fmt.Sprintf("normal size: %s; diff size: %s",
		bytesize.New(float64(totalSizeA)),
		bytesize.New(float64(proto.Size(asSeq))))
}

// Code for differential encoding of State objects.
func decodeSeq(seq *chiltrix.StateSequence) ([]*chiltrix.State, error) {
	var states []*chiltrix.State
	var prev *chiltrix.State
	for i, d := range seq.GetDiffs() {
		ith, err := decodeDiff(prev, d)
		if err != nil {
			return nil, fmt.Errorf("error decoding %dth diff: %w", i, err)
		}
		states = append(states, ith)
	}
	return states, nil
}

func encodeSeq(states []*chiltrix.State) (*chiltrix.StateSequence, error) {
	var diffs []*chiltrix.StateDiff
	for i, s := range states {
		var prev *chiltrix.State
		if i > 0 {
			prev = states[i-1]
		}
		diffs = append(diffs, diffStates(prev, s))
	}
	return &chiltrix.StateSequence{
		Diffs: diffs,
	}, nil
}

func diffStates(a, b *chiltrix.State) *chiltrix.StateDiff {
	diff := &chiltrix.StateDiff{
		CollectionTime: b.GetCollectionTime(),
		UpdatedValues:  &chiltrix.RegisterSnapshot{HoldingRegisterValues: map[uint32]uint32{}},
	}
	for k, v := range b.GetRegisterValues().GetHoldingRegisterValues() {
		prev, hasPrev := a.GetRegisterValues().GetHoldingRegisterValues()[k]
		if !hasPrev || prev != v {
			diff.GetUpdatedValues().GetHoldingRegisterValues()[k] = v
		}
	}
	for k := range a.GetRegisterValues().GetHoldingRegisterValues() {
		if _, alreadySet := diff.GetUpdatedValues().GetHoldingRegisterValues()[k]; alreadySet {
			continue
		}
		diff.DeletedRegisters = append(diff.DeletedRegisters, k)
	}
	return diff
}

func decodeDiff(a *chiltrix.State, diff *chiltrix.StateDiff) (*chiltrix.State, error) {
	if diff.GetCollectionTime() == nil {
		return nil, fmt.Errorf("diff must specify collection time")
	}
	out := &chiltrix.State{
		CollectionTime: diff.GetCollectionTime(),
		RegisterValues: &chiltrix.RegisterSnapshot{HoldingRegisterValues: make(map[uint32]uint32)},
	}
	for k, v := range a.GetRegisterValues().GetHoldingRegisterValues() {
		out.GetRegisterValues().GetHoldingRegisterValues()[k] = v
	}
	for _, deleted := range diff.GetDeletedRegisters() {
		if _, ok := out.GetRegisterValues().GetHoldingRegisterValues()[deleted]; !ok {
			return nil, fmt.Errorf("diff says to delete key %d, but that key isn't present in the basis", deleted)
		}
	}
	for k, v := range diff.GetUpdatedValues().GetHoldingRegisterValues() {
		out.GetRegisterValues().GetHoldingRegisterValues()[k] = v
	}
	return out, nil
}
