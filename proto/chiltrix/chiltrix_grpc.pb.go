// Code generated by protoc-gen-go-grpc. DO NOT EDIT.

package chiltrix

import (
	context "context"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
)

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
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
	stream, err := c.cc.NewStream(ctx, &_Historian_serviceDesc.Streams[0], "/chiltrix.Historian/QueryStream", opts...)
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
	s.RegisterService(&_Historian_serviceDesc, srv)
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

var _Historian_serviceDesc = grpc.ServiceDesc{
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