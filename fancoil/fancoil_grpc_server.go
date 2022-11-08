package fancoil

import (
	"context"

	"github.com/golang/glog"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

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
			return nil, status.Errorf(codes.Unavailable, "error setting modbus value: %v", err)
		}
	}
	if req.GetPowerStatus() != pb.PowerStatus_POWER_STATUS_UNSPECIFIED {
		registerValue, err := powerStatusToModbusValue(req.GetPowerStatus())
		if err != nil {
			return nil, err
		}
		if err := s.c.setRegisterValue(Register(pb.RegisterName_REGISTER_NAME_ON_OFF), registerValue); err != nil {
			return nil, status.Errorf(codes.Unavailable, "error setting modbus value: %v", err)
		}
	}
	if req.GetValveSetting() != pb.ValveSetting_VALVE_SETTING_UNSPECIFIED {
		registerValue, ok := valveSettingMap.getByEnumNumber(req.GetValveSetting().Number())
		if !ok {
			return nil, status.Errorf(codes.Internal, "enum %s has no known encoding as a modbus register", req.GetValveSetting())
		}
		if err := s.c.setRegisterValue(Register(pb.RegisterName_REGISTER_NAME_USE_VALVE), registerValue); err != nil {
			return nil, status.Errorf(codes.Unavailable, "error setting modbus value: %v", err)
		}
	}
	if req.GetHeatingTargetTemperature() != nil {
		registerValue, err := tempToModbusValue(req.GetHeatingTargetTemperature())
		if err != nil {
			return nil, status.Errorf(codes.InvalidArgument, "invalid heating_target_temperature value: %v", err)
		}
		if err := s.c.setRegisterValue(Register(pb.RegisterName_REGISTER_NAME_HEATING_SET_TEMPERATURE), registerValue); err != nil {
			return nil, status.Errorf(codes.Unavailable, "error setting modbus value: %v", err)
		}
	}
	if req.GetCoolingTargetTemperature() != nil {
		registerValue, err := tempToModbusValue(req.GetCoolingTargetTemperature())
		if err != nil {
			return nil, status.Errorf(codes.InvalidArgument, "invalid cooling_target_temperature value: %v", err)
		}
		if err := s.c.setRegisterValue(Register(pb.RegisterName_REGISTER_NAME_COOLING_SET_TEMPERATURE), registerValue); err != nil {
			return nil, status.Errorf(codes.Unavailable, "error setting modbus value: %v", err)
		}
	}
	return &pb.SetStateResponse{}, nil
}
