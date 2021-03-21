package fancoil

import (
	"context"

	"github.com/golang/glog"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"

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
	glog.Infof("SetState request received: %s", req)
	if req.GetPreferenceFanSetting() != pb.FanSetting_FAN_SETTING_UNSPECIFIED {
		registerValue, err := fanSettingToModbusValue(req.GetPreferenceFanSetting())
		if err != nil {
			return nil, err
		}
		if err := s.c.setRegisterValue(Register(pb.RegisterName_REGISTER_NAME_FANSPEED), registerValue); err != nil {
			return nil, grpc.Errorf(codes.Unavailable, "error setting modbus value: %v", err)
		}
	}
	if req.GetPowerStatus() != pb.PowerStatus_POWER_STATUS_UNSPECIFIED {
		registerValue, err := powerStatusToModbusValue(req.GetPowerStatus())
		if err != nil {
			return nil, err
		}
		if err := s.c.setRegisterValue(Register(pb.RegisterName_REGISTER_NAME_ON_OFF), registerValue); err != nil {
			return nil, grpc.Errorf(codes.Unavailable, "error setting modbus value: %v", err)
		}
	}
	return &pb.SetStateResponse{}, nil
}
