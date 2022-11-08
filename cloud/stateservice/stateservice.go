// Package stateservice implements a centralized IoT device state storage
// service. See the controller.StateService proto.
package stateservice

import (
	"context"
	"fmt"

	"github.com/gonzojive/heatpump/cloud/acls"
	"github.com/gonzojive/heatpump/cloud/google/cloudconfig"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	cpb "github.com/gonzojive/heatpump/proto/controller"
)

// Service implements controller.StateService.
type Service struct {
	cpb.UnimplementedStateServiceServer

	aclsService *acls.Service
	db          Store
}

// New returns a new StateService implementation.
func New(ctx context.Context, projectParams *cloudconfig.Params) (*Service, error) {
	store, err := newFirestoreBackedStore(ctx, projectParams)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize db: %w", err)
	}

	return &Service{
		aclsService: acls.NewService(projectParams),
		db:          store,
	}, nil
}

func (s *Service) GetDeviceState(ctx context.Context, req *cpb.GetDeviceStateRequest) (*cpb.DeviceState, error) {
	identity, err := s.aclsService.IdentityFromContext(ctx)
	if err != nil {
		return nil, status.Errorf(codes.Unauthenticated, "Client must pass a %q gRPC metadata header to identify the device using a DeviceAccessToken: %v", acls.DeviceAccessTokenMetadataKey, err)
	}
	state, err := s.db.GetDeviceState(ctx, identity, req.GetName())
	if err != nil {
		return nil, err
	}
	return state, nil
}

func (s *Service) SetDeviceState(ctx context.Context, req *cpb.SetDeviceStateRequest) (*cpb.SetDeviceStateResponse, error) {
	identity, err := s.aclsService.IdentityFromContext(ctx)
	if err != nil {
		return nil, status.Errorf(codes.Unauthenticated, "Client must pass a %q gRPC metadata header to identify the device using a DeviceAccessToken: %v", acls.DeviceAccessTokenMetadataKey, err)
	}

	state := req.GetState()
	if state == nil {
		return nil, status.Errorf(codes.InvalidArgument, "must specify state field")
	}

	if err := s.db.StoreDeviceState(ctx, identity, state); err != nil {
		return nil, status.Errorf(codes.Internal, "error storing device state: %v", err)
	}

	return &cpb.SetDeviceStateResponse{}, nil
}
