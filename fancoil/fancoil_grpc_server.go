package fancoil

import (
	"context"

	fancoil "github.com/gonzojive/heatpump/proto/fancoil"
	pb "github.com/gonzojive/heatpump/proto/fancoil"
)

// GRPCServiceFromClient returns an implementation of FanCoilServiceServer using
// a *Client object.
func GRPCServiceFromClient(c *Client) pb.FanCoilServiceServer {
	return &server{c: c}
}

type server struct {
	pb.UnimplementedFanCoilServiceServer

	c *Client
}

// Get a snapshot of the state of a single fan coil unit.
func (s *server) GetState(ctx context.Context, req *fancoil.GetStateRequest) (*fancoil.GetStateResponse, error) {
	return s.c.GetState(ctx, req)
}
