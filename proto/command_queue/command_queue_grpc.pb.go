// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.2.0
// - protoc             v3.12.4
// source: command_queue.proto

package command_queue

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

// CommandQueueServiceClient is the client API for CommandQueueService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type CommandQueueServiceClient interface {
	// Listens for events
	Listen(ctx context.Context, opts ...grpc.CallOption) (CommandQueueService_ListenClient, error)
	// Lists the topics the client is eligible to subscribe to.
	ListTopics(ctx context.Context, in *ListTopicsRequest, opts ...grpc.CallOption) (*ListTopicsResponse, error)
}

type commandQueueServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewCommandQueueServiceClient(cc grpc.ClientConnInterface) CommandQueueServiceClient {
	return &commandQueueServiceClient{cc}
}

func (c *commandQueueServiceClient) Listen(ctx context.Context, opts ...grpc.CallOption) (CommandQueueService_ListenClient, error) {
	stream, err := c.cc.NewStream(ctx, &CommandQueueService_ServiceDesc.Streams[0], "/heatpump.command_queue.CommandQueueService/Listen", opts...)
	if err != nil {
		return nil, err
	}
	x := &commandQueueServiceListenClient{stream}
	return x, nil
}

type CommandQueueService_ListenClient interface {
	Send(*ListenRequest) error
	Recv() (*ListenResponse, error)
	grpc.ClientStream
}

type commandQueueServiceListenClient struct {
	grpc.ClientStream
}

func (x *commandQueueServiceListenClient) Send(m *ListenRequest) error {
	return x.ClientStream.SendMsg(m)
}

func (x *commandQueueServiceListenClient) Recv() (*ListenResponse, error) {
	m := new(ListenResponse)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

func (c *commandQueueServiceClient) ListTopics(ctx context.Context, in *ListTopicsRequest, opts ...grpc.CallOption) (*ListTopicsResponse, error) {
	out := new(ListTopicsResponse)
	err := c.cc.Invoke(ctx, "/heatpump.command_queue.CommandQueueService/ListTopics", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// CommandQueueServiceServer is the server API for CommandQueueService service.
// All implementations must embed UnimplementedCommandQueueServiceServer
// for forward compatibility
type CommandQueueServiceServer interface {
	// Listens for events
	Listen(CommandQueueService_ListenServer) error
	// Lists the topics the client is eligible to subscribe to.
	ListTopics(context.Context, *ListTopicsRequest) (*ListTopicsResponse, error)
	mustEmbedUnimplementedCommandQueueServiceServer()
}

// UnimplementedCommandQueueServiceServer must be embedded to have forward compatible implementations.
type UnimplementedCommandQueueServiceServer struct {
}

func (UnimplementedCommandQueueServiceServer) Listen(CommandQueueService_ListenServer) error {
	return status.Errorf(codes.Unimplemented, "method Listen not implemented")
}
func (UnimplementedCommandQueueServiceServer) ListTopics(context.Context, *ListTopicsRequest) (*ListTopicsResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ListTopics not implemented")
}
func (UnimplementedCommandQueueServiceServer) mustEmbedUnimplementedCommandQueueServiceServer() {}

// UnsafeCommandQueueServiceServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to CommandQueueServiceServer will
// result in compilation errors.
type UnsafeCommandQueueServiceServer interface {
	mustEmbedUnimplementedCommandQueueServiceServer()
}

func RegisterCommandQueueServiceServer(s grpc.ServiceRegistrar, srv CommandQueueServiceServer) {
	s.RegisterService(&CommandQueueService_ServiceDesc, srv)
}

func _CommandQueueService_Listen_Handler(srv interface{}, stream grpc.ServerStream) error {
	return srv.(CommandQueueServiceServer).Listen(&commandQueueServiceListenServer{stream})
}

type CommandQueueService_ListenServer interface {
	Send(*ListenResponse) error
	Recv() (*ListenRequest, error)
	grpc.ServerStream
}

type commandQueueServiceListenServer struct {
	grpc.ServerStream
}

func (x *commandQueueServiceListenServer) Send(m *ListenResponse) error {
	return x.ServerStream.SendMsg(m)
}

func (x *commandQueueServiceListenServer) Recv() (*ListenRequest, error) {
	m := new(ListenRequest)
	if err := x.ServerStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

func _CommandQueueService_ListTopics_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ListTopicsRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(CommandQueueServiceServer).ListTopics(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/heatpump.command_queue.CommandQueueService/ListTopics",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(CommandQueueServiceServer).ListTopics(ctx, req.(*ListTopicsRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// CommandQueueService_ServiceDesc is the grpc.ServiceDesc for CommandQueueService service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var CommandQueueService_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "heatpump.command_queue.CommandQueueService",
	HandlerType: (*CommandQueueServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "ListTopics",
			Handler:    _CommandQueueService_ListTopics_Handler,
		},
	},
	Streams: []grpc.StreamDesc{
		{
			StreamName:    "Listen",
			Handler:       _CommandQueueService_Listen_Handler,
			ServerStreams: true,
			ClientStreams: true,
		},
	},
	Metadata: "command_queue.proto",
}
