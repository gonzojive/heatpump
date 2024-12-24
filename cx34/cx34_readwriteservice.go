package cx34

import (
	"context"

	"github.com/gonzojive/heatpump/proto/chiltrix"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// ReadWriteServiceServer is an implementation of the Chiltrix ReadWriteService
// that uses a local cx34 client to communicate with the heat pump.
type ReadWriteServiceServer struct {
	chiltrix.UnimplementedReadWriteServiceServer
	client *Client
}

var (
	_ chiltrix.ReadWriteServiceServer = (*ReadWriteServiceServer)(nil)
)

// SetParameter implements the SetParameter method of the
// [chiltrix.ReadWriteServiceServer] interface.
func (s *ReadWriteServiceServer) SetParameter(ctx context.Context, req *chiltrix.SetParameterRequest) (*chiltrix.SetParameterResponse, error) {
	return s.client.SetParameter(ctx, req)
}

// GetState implements the GetState method of the
// [chiltrix.ReadWriteServiceServer] interface.
func (s *ReadWriteServiceServer) GetState(ctx context.Context, _ *chiltrix.GetStateRequest) (*chiltrix.State, error) {
	state, err := s.client.ReadState()
	if err != nil {
		status.Errorf(codes.Internal, "error reading state of Chiltrix heat pump: %v", err)
	}
	return state.Proto(), nil
}
