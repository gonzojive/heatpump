// Package stateservice implements a centralized IoT device state storage
// service. See the controller.StateService proto.
package stateservice

import (
	"context"
	"fmt"

	"github.com/golang/glog"
	"github.com/gonzojive/heatpump/cloud/acls"
	"github.com/gonzojive/heatpump/cloud/acls/server2serverauth"
	"github.com/gonzojive/heatpump/cloud/google/cloudconfig"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/encoding/prototext"

	cpb "github.com/gonzojive/heatpump/proto/controller"
)

const (
	statServiceAudience = "stateservice"
	permittedEmail      = "google-actions-http-job@heatpump-dev.iam.gserviceaccount.com"
)

// Service implements controller.StateService.
type Service struct {
	cpb.UnimplementedStateServiceServer

	aclsService      *acls.Service
	db               Store
	s2sAuthValidator *server2serverauth.Validator
}

// New returns a new StateService implementation.
func New(ctx context.Context, projectParams *cloudconfig.Params) (*Service, error) {
	store, err := newFirestoreBackedStore(ctx, projectParams)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize db: %w", err)
	}

	validator, err := server2serverauth.NewValidator(ctx, statServiceAudience)
	if err != nil {
		return nil, fmt.Errorf("failed to create validator: %w", err)
	}

	return &Service{
		aclsService:      acls.NewService(projectParams),
		db:               store,
		s2sAuthValidator: validator,
	}, nil
}

func (s *Service) GetDeviceState(ctx context.Context, req *cpb.GetDeviceStateRequest) (*cpb.DeviceState, error) {
	identity, err := s.aclsService.IdentityFromContext(ctx)
	if err != nil {
		clientInfo, err2 := s.s2sAuthValidator.ValidateFromIncomingContext(ctx)
		if err2 != nil {
			glog.Errorf("Authentication failed for both device and server type clients; device err = %v, server err = %v", err, err2)

			return nil, status.Errorf(codes.Unauthenticated, "Client must pass a valid %q gRPC metadata header to identify the device using a DeviceAccessToken or an \"Authoriation: Bearer xxx\" type header with a token provided by the idtoken package", acls.DeviceAccessTokenMetadataKey)
		}
		glog.Infof("authenticated sender is %q", clientInfo.Email())
		if clientInfo.Email() != permittedEmail {
			glog.Errorf("unexpected authenticated sender %q", clientInfo.Email())
			return nil, fmt.Errorf("unexpected authenticated sender %q", clientInfo.Email())
		}
		identity = acls.FixmeMainHardcodedIdentity()
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

	glog.Infof("set device state for %q to %s", state.GetName(), prototext.Format(state))

	if err := s.db.StoreDeviceState(ctx, identity, state); err != nil {
		return nil, status.Errorf(codes.Internal, "error storing device state: %v", err)
	}

	return &cpb.SetDeviceStateResponse{}, nil
}
