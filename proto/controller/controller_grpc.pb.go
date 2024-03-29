// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.2.0
// - protoc             v3.12.4
// source: controller.proto

package controller

import (
	context "context"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
)

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
// Requires gRPC-Go v1.32.0 or later.
const _ = grpc.SupportPackageIsVersion7

// StateServiceClient is the client API for StateService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type StateServiceClient interface {
	// Get a snapshot of the state of a single fan coil unit.
	GetDeviceState(ctx context.Context, in *GetDeviceStateRequest, opts ...grpc.CallOption) (*DeviceState, error)
	// Set some parameters of a device.
	SetDeviceState(ctx context.Context, in *SetDeviceStateRequest, opts ...grpc.CallOption) (*SetDeviceStateResponse, error)
}

type stateServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewStateServiceClient(cc grpc.ClientConnInterface) StateServiceClient {
	return &stateServiceClient{cc}
}

func (c *stateServiceClient) GetDeviceState(ctx context.Context, in *GetDeviceStateRequest, opts ...grpc.CallOption) (*DeviceState, error) {
	out := new(DeviceState)
	err := c.cc.Invoke(ctx, "/controller.StateService/GetDeviceState", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *stateServiceClient) SetDeviceState(ctx context.Context, in *SetDeviceStateRequest, opts ...grpc.CallOption) (*SetDeviceStateResponse, error) {
	out := new(SetDeviceStateResponse)
	err := c.cc.Invoke(ctx, "/controller.StateService/SetDeviceState", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// StateServiceServer is the server API for StateService service.
// All implementations must embed UnimplementedStateServiceServer
// for forward compatibility
type StateServiceServer interface {
	// Get a snapshot of the state of a single fan coil unit.
	GetDeviceState(context.Context, *GetDeviceStateRequest) (*DeviceState, error)
	// Set some parameters of a device.
	SetDeviceState(context.Context, *SetDeviceStateRequest) (*SetDeviceStateResponse, error)
	mustEmbedUnimplementedStateServiceServer()
}

// UnimplementedStateServiceServer must be embedded to have forward compatible implementations.
type UnimplementedStateServiceServer struct {
}

func (UnimplementedStateServiceServer) GetDeviceState(context.Context, *GetDeviceStateRequest) (*DeviceState, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetDeviceState not implemented")
}
func (UnimplementedStateServiceServer) SetDeviceState(context.Context, *SetDeviceStateRequest) (*SetDeviceStateResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method SetDeviceState not implemented")
}
func (UnimplementedStateServiceServer) mustEmbedUnimplementedStateServiceServer() {}

// UnsafeStateServiceServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to StateServiceServer will
// result in compilation errors.
type UnsafeStateServiceServer interface {
	mustEmbedUnimplementedStateServiceServer()
}

func RegisterStateServiceServer(s grpc.ServiceRegistrar, srv StateServiceServer) {
	s.RegisterService(&StateService_ServiceDesc, srv)
}

func _StateService_GetDeviceState_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetDeviceStateRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(StateServiceServer).GetDeviceState(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/controller.StateService/GetDeviceState",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(StateServiceServer).GetDeviceState(ctx, req.(*GetDeviceStateRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _StateService_SetDeviceState_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(SetDeviceStateRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(StateServiceServer).SetDeviceState(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/controller.StateService/SetDeviceState",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(StateServiceServer).SetDeviceState(ctx, req.(*SetDeviceStateRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// StateService_ServiceDesc is the grpc.ServiceDesc for StateService service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var StateService_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "controller.StateService",
	HandlerType: (*StateServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "GetDeviceState",
			Handler:    _StateService_GetDeviceState_Handler,
		},
		{
			MethodName: "SetDeviceState",
			Handler:    _StateService_SetDeviceState_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "controller.proto",
}

// AuthServiceClient is the client API for AuthService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type AuthServiceClient interface {
	// Get a snapshot of the state of a single fan coil unit.
	ExtendToken(ctx context.Context, in *ExtendTokenRequest, opts ...grpc.CallOption) (*ExtendTokenResponse, error)
}

type authServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewAuthServiceClient(cc grpc.ClientConnInterface) AuthServiceClient {
	return &authServiceClient{cc}
}

func (c *authServiceClient) ExtendToken(ctx context.Context, in *ExtendTokenRequest, opts ...grpc.CallOption) (*ExtendTokenResponse, error) {
	out := new(ExtendTokenResponse)
	err := c.cc.Invoke(ctx, "/controller.AuthService/ExtendToken", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// AuthServiceServer is the server API for AuthService service.
// All implementations must embed UnimplementedAuthServiceServer
// for forward compatibility
type AuthServiceServer interface {
	// Get a snapshot of the state of a single fan coil unit.
	ExtendToken(context.Context, *ExtendTokenRequest) (*ExtendTokenResponse, error)
	mustEmbedUnimplementedAuthServiceServer()
}

// UnimplementedAuthServiceServer must be embedded to have forward compatible implementations.
type UnimplementedAuthServiceServer struct {
}

func (UnimplementedAuthServiceServer) ExtendToken(context.Context, *ExtendTokenRequest) (*ExtendTokenResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ExtendToken not implemented")
}
func (UnimplementedAuthServiceServer) mustEmbedUnimplementedAuthServiceServer() {}

// UnsafeAuthServiceServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to AuthServiceServer will
// result in compilation errors.
type UnsafeAuthServiceServer interface {
	mustEmbedUnimplementedAuthServiceServer()
}

func RegisterAuthServiceServer(s grpc.ServiceRegistrar, srv AuthServiceServer) {
	s.RegisterService(&AuthService_ServiceDesc, srv)
}

func _AuthService_ExtendToken_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ExtendTokenRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(AuthServiceServer).ExtendToken(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/controller.AuthService/ExtendToken",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(AuthServiceServer).ExtendToken(ctx, req.(*ExtendTokenRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// AuthService_ServiceDesc is the grpc.ServiceDesc for AuthService service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var AuthService_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "controller.AuthService",
	HandlerType: (*AuthServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "ExtendToken",
			Handler:    _AuthService_ExtendToken_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "controller.proto",
}
