// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.2.0
// - protoc             v3.21.12
// source: chiltrix.proto

package chiltrix

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

// HistorianClient is the client API for Historian service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type HistorianClient interface {
	// Provides a value for each key request
	QueryStream(ctx context.Context, in *QueryStreamRequest, opts ...grpc.CallOption) (Historian_QueryStreamClient, error)
}

type historianClient struct {
	cc grpc.ClientConnInterface
}

func NewHistorianClient(cc grpc.ClientConnInterface) HistorianClient {
	return &historianClient{cc}
}

func (c *historianClient) QueryStream(ctx context.Context, in *QueryStreamRequest, opts ...grpc.CallOption) (Historian_QueryStreamClient, error) {
	stream, err := c.cc.NewStream(ctx, &Historian_ServiceDesc.Streams[0], "/chiltrix.Historian/QueryStream", opts...)
	if err != nil {
		return nil, err
	}
	x := &historianQueryStreamClient{stream}
	if err := x.ClientStream.SendMsg(in); err != nil {
		return nil, err
	}
	if err := x.ClientStream.CloseSend(); err != nil {
		return nil, err
	}
	return x, nil
}

type Historian_QueryStreamClient interface {
	Recv() (*QueryStreamResponse, error)
	grpc.ClientStream
}

type historianQueryStreamClient struct {
	grpc.ClientStream
}

func (x *historianQueryStreamClient) Recv() (*QueryStreamResponse, error) {
	m := new(QueryStreamResponse)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

// HistorianServer is the server API for Historian service.
// All implementations must embed UnimplementedHistorianServer
// for forward compatibility
type HistorianServer interface {
	// Provides a value for each key request
	QueryStream(*QueryStreamRequest, Historian_QueryStreamServer) error
	mustEmbedUnimplementedHistorianServer()
}

// UnimplementedHistorianServer must be embedded to have forward compatible implementations.
type UnimplementedHistorianServer struct {
}

func (UnimplementedHistorianServer) QueryStream(*QueryStreamRequest, Historian_QueryStreamServer) error {
	return status.Errorf(codes.Unimplemented, "method QueryStream not implemented")
}
func (UnimplementedHistorianServer) mustEmbedUnimplementedHistorianServer() {}

// UnsafeHistorianServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to HistorianServer will
// result in compilation errors.
type UnsafeHistorianServer interface {
	mustEmbedUnimplementedHistorianServer()
}

func RegisterHistorianServer(s grpc.ServiceRegistrar, srv HistorianServer) {
	s.RegisterService(&Historian_ServiceDesc, srv)
}

func _Historian_QueryStream_Handler(srv interface{}, stream grpc.ServerStream) error {
	m := new(QueryStreamRequest)
	if err := stream.RecvMsg(m); err != nil {
		return err
	}
	return srv.(HistorianServer).QueryStream(m, &historianQueryStreamServer{stream})
}

type Historian_QueryStreamServer interface {
	Send(*QueryStreamResponse) error
	grpc.ServerStream
}

type historianQueryStreamServer struct {
	grpc.ServerStream
}

func (x *historianQueryStreamServer) Send(m *QueryStreamResponse) error {
	return x.ServerStream.SendMsg(m)
}

// Historian_ServiceDesc is the grpc.ServiceDesc for Historian service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var Historian_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "chiltrix.Historian",
	HandlerType: (*HistorianServer)(nil),
	Methods:     []grpc.MethodDesc{},
	Streams: []grpc.StreamDesc{
		{
			StreamName:    "QueryStream",
			Handler:       _Historian_QueryStream_Handler,
			ServerStreams: true,
		},
	},
	Metadata: "chiltrix.proto",
}

// ReadWriteServiceClient is the client API for ReadWriteService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type ReadWriteServiceClient interface {
	// Change a parameter.
	SetParameter(ctx context.Context, in *SetParameterRequest, opts ...grpc.CallOption) (*SetParameterResponse, error)
	// Retrieve the state of the heatpump (all modbus parameters).
	GetState(ctx context.Context, in *GetStateRequest, opts ...grpc.CallOption) (*State, error)
}

type readWriteServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewReadWriteServiceClient(cc grpc.ClientConnInterface) ReadWriteServiceClient {
	return &readWriteServiceClient{cc}
}

func (c *readWriteServiceClient) SetParameter(ctx context.Context, in *SetParameterRequest, opts ...grpc.CallOption) (*SetParameterResponse, error) {
	out := new(SetParameterResponse)
	err := c.cc.Invoke(ctx, "/chiltrix.ReadWriteService/SetParameter", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *readWriteServiceClient) GetState(ctx context.Context, in *GetStateRequest, opts ...grpc.CallOption) (*State, error) {
	out := new(State)
	err := c.cc.Invoke(ctx, "/chiltrix.ReadWriteService/GetState", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// ReadWriteServiceServer is the server API for ReadWriteService service.
// All implementations must embed UnimplementedReadWriteServiceServer
// for forward compatibility
type ReadWriteServiceServer interface {
	// Change a parameter.
	SetParameter(context.Context, *SetParameterRequest) (*SetParameterResponse, error)
	// Retrieve the state of the heatpump (all modbus parameters).
	GetState(context.Context, *GetStateRequest) (*State, error)
	mustEmbedUnimplementedReadWriteServiceServer()
}

// UnimplementedReadWriteServiceServer must be embedded to have forward compatible implementations.
type UnimplementedReadWriteServiceServer struct {
}

func (UnimplementedReadWriteServiceServer) SetParameter(context.Context, *SetParameterRequest) (*SetParameterResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method SetParameter not implemented")
}
func (UnimplementedReadWriteServiceServer) GetState(context.Context, *GetStateRequest) (*State, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetState not implemented")
}
func (UnimplementedReadWriteServiceServer) mustEmbedUnimplementedReadWriteServiceServer() {}

// UnsafeReadWriteServiceServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to ReadWriteServiceServer will
// result in compilation errors.
type UnsafeReadWriteServiceServer interface {
	mustEmbedUnimplementedReadWriteServiceServer()
}

func RegisterReadWriteServiceServer(s grpc.ServiceRegistrar, srv ReadWriteServiceServer) {
	s.RegisterService(&ReadWriteService_ServiceDesc, srv)
}

func _ReadWriteService_SetParameter_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(SetParameterRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ReadWriteServiceServer).SetParameter(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/chiltrix.ReadWriteService/SetParameter",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ReadWriteServiceServer).SetParameter(ctx, req.(*SetParameterRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _ReadWriteService_GetState_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetStateRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ReadWriteServiceServer).GetState(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/chiltrix.ReadWriteService/GetState",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ReadWriteServiceServer).GetState(ctx, req.(*GetStateRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// ReadWriteService_ServiceDesc is the grpc.ServiceDesc for ReadWriteService service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var ReadWriteService_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "chiltrix.ReadWriteService",
	HandlerType: (*ReadWriteServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "SetParameter",
			Handler:    _ReadWriteService_SetParameter_Handler,
		},
		{
			MethodName: "GetState",
			Handler:    _ReadWriteService_GetState_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "chiltrix.proto",
}
