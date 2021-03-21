package fancoil

import (
	"context"

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
func (s *server) GetState(ctx context.Context, req *pb.GetStateRequest) (*pb.GetStateResponse, error) {
	return s.c.GetState(ctx, req)
}

// Get a snapshot of the state of a single fan coil unit.
func (s *server) SetState(ctx context.Context, req *pb.SetStateRequest) (*pb.SetStateResponse, error) {
	if req.GetPreferenceFanSetting() != pb.FanSetting_FAN_SETTING_UNSPECIFIED {
		registerValue, err := fanSettingToModbusValue(req.GetPreferenceFanSetting())
		if err != nil {
			return nil, err
		}
		s.c.setRegisterValue(Register(pb.RegisterName_REGISTER_NAME_FANSPEED), registerValue)
	}
	return &pb.SetStateResponse{}, nil
}
