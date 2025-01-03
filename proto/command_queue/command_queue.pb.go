// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.36.1
// 	protoc        v5.29.0
// source: proto/command_queue/command_queue.proto

package command_queue

import (
	context "context"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
	reflect "reflect"
	sync "sync"
)

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

type ListenRequest struct {
	state     protoimpl.MessageState `protogen:"open.v1"`
	RequestId string                 `protobuf:"bytes,1,opt,name=request_id,json=requestId,proto3" json:"request_id,omitempty"`
	// Types that are valid to be assigned to Request:
	//
	//	*ListenRequest_SubscribeRequest_
	//	*ListenRequest_AckRequest_
	Request       isListenRequest_Request `protobuf_oneof:"request"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *ListenRequest) Reset() {
	*x = ListenRequest{}
	mi := &file_proto_command_queue_command_queue_proto_msgTypes[0]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *ListenRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ListenRequest) ProtoMessage() {}

func (x *ListenRequest) ProtoReflect() protoreflect.Message {
	mi := &file_proto_command_queue_command_queue_proto_msgTypes[0]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ListenRequest.ProtoReflect.Descriptor instead.
func (*ListenRequest) Descriptor() ([]byte, []int) {
	return file_proto_command_queue_command_queue_proto_rawDescGZIP(), []int{0}
}

func (x *ListenRequest) GetRequestId() string {
	if x != nil {
		return x.RequestId
	}
	return ""
}

func (x *ListenRequest) GetRequest() isListenRequest_Request {
	if x != nil {
		return x.Request
	}
	return nil
}

func (x *ListenRequest) GetSubscribeRequest() *ListenRequest_SubscribeRequest {
	if x != nil {
		if x, ok := x.Request.(*ListenRequest_SubscribeRequest_); ok {
			return x.SubscribeRequest
		}
	}
	return nil
}

func (x *ListenRequest) GetAckRequest() *ListenRequest_AckRequest {
	if x != nil {
		if x, ok := x.Request.(*ListenRequest_AckRequest_); ok {
			return x.AckRequest
		}
	}
	return nil
}

type isListenRequest_Request interface {
	isListenRequest_Request()
}

type ListenRequest_SubscribeRequest_ struct {
	SubscribeRequest *ListenRequest_SubscribeRequest `protobuf:"bytes,2,opt,name=subscribe_request,json=subscribeRequest,proto3,oneof"`
}

type ListenRequest_AckRequest_ struct {
	AckRequest *ListenRequest_AckRequest `protobuf:"bytes,3,opt,name=ack_request,json=ackRequest,proto3,oneof"`
}

func (*ListenRequest_SubscribeRequest_) isListenRequest_Request() {}

func (*ListenRequest_AckRequest_) isListenRequest_Request() {}

type ListenResponse struct {
	state     protoimpl.MessageState `protogen:"open.v1"`
	RequestId string                 `protobuf:"bytes,1,opt,name=request_id,json=requestId,proto3" json:"request_id,omitempty"`
	// Types that are valid to be assigned to Response:
	//
	//	*ListenResponse_MessageResponse
	//	*ListenResponse_AckResponse_
	Response      isListenResponse_Response `protobuf_oneof:"response"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *ListenResponse) Reset() {
	*x = ListenResponse{}
	mi := &file_proto_command_queue_command_queue_proto_msgTypes[1]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *ListenResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ListenResponse) ProtoMessage() {}

func (x *ListenResponse) ProtoReflect() protoreflect.Message {
	mi := &file_proto_command_queue_command_queue_proto_msgTypes[1]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ListenResponse.ProtoReflect.Descriptor instead.
func (*ListenResponse) Descriptor() ([]byte, []int) {
	return file_proto_command_queue_command_queue_proto_rawDescGZIP(), []int{1}
}

func (x *ListenResponse) GetRequestId() string {
	if x != nil {
		return x.RequestId
	}
	return ""
}

func (x *ListenResponse) GetResponse() isListenResponse_Response {
	if x != nil {
		return x.Response
	}
	return nil
}

func (x *ListenResponse) GetMessageResponse() *MessageResponse {
	if x != nil {
		if x, ok := x.Response.(*ListenResponse_MessageResponse); ok {
			return x.MessageResponse
		}
	}
	return nil
}

func (x *ListenResponse) GetAckResponse() *ListenResponse_AckResponse {
	if x != nil {
		if x, ok := x.Response.(*ListenResponse_AckResponse_); ok {
			return x.AckResponse
		}
	}
	return nil
}

type isListenResponse_Response interface {
	isListenResponse_Response()
}

type ListenResponse_MessageResponse struct {
	MessageResponse *MessageResponse `protobuf:"bytes,2,opt,name=message_response,json=messageResponse,proto3,oneof"`
}

type ListenResponse_AckResponse_ struct {
	AckResponse *ListenResponse_AckResponse `protobuf:"bytes,3,opt,name=ack_response,json=ackResponse,proto3,oneof"`
}

func (*ListenResponse_MessageResponse) isListenResponse_Response() {}

func (*ListenResponse_AckResponse_) isListenResponse_Response() {}

type MessageResponse struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	Id            string                 `protobuf:"bytes,1,opt,name=id,proto3" json:"id,omitempty"`
	Payload       []byte                 `protobuf:"bytes,2,opt,name=payload,proto3" json:"payload,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *MessageResponse) Reset() {
	*x = MessageResponse{}
	mi := &file_proto_command_queue_command_queue_proto_msgTypes[2]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *MessageResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*MessageResponse) ProtoMessage() {}

func (x *MessageResponse) ProtoReflect() protoreflect.Message {
	mi := &file_proto_command_queue_command_queue_proto_msgTypes[2]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use MessageResponse.ProtoReflect.Descriptor instead.
func (*MessageResponse) Descriptor() ([]byte, []int) {
	return file_proto_command_queue_command_queue_proto_rawDescGZIP(), []int{2}
}

func (x *MessageResponse) GetId() string {
	if x != nil {
		return x.Id
	}
	return ""
}

func (x *MessageResponse) GetPayload() []byte {
	if x != nil {
		return x.Payload
	}
	return nil
}

type ListTopicsRequest struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *ListTopicsRequest) Reset() {
	*x = ListTopicsRequest{}
	mi := &file_proto_command_queue_command_queue_proto_msgTypes[3]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *ListTopicsRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ListTopicsRequest) ProtoMessage() {}

func (x *ListTopicsRequest) ProtoReflect() protoreflect.Message {
	mi := &file_proto_command_queue_command_queue_proto_msgTypes[3]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ListTopicsRequest.ProtoReflect.Descriptor instead.
func (*ListTopicsRequest) Descriptor() ([]byte, []int) {
	return file_proto_command_queue_command_queue_proto_rawDescGZIP(), []int{3}
}

type ListTopicsResponse struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	Topics        []string               `protobuf:"bytes,1,rep,name=topics,proto3" json:"topics,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *ListTopicsResponse) Reset() {
	*x = ListTopicsResponse{}
	mi := &file_proto_command_queue_command_queue_proto_msgTypes[4]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *ListTopicsResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ListTopicsResponse) ProtoMessage() {}

func (x *ListTopicsResponse) ProtoReflect() protoreflect.Message {
	mi := &file_proto_command_queue_command_queue_proto_msgTypes[4]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ListTopicsResponse.ProtoReflect.Descriptor instead.
func (*ListTopicsResponse) Descriptor() ([]byte, []int) {
	return file_proto_command_queue_command_queue_proto_rawDescGZIP(), []int{4}
}

func (x *ListTopicsResponse) GetTopics() []string {
	if x != nil {
		return x.Topics
	}
	return nil
}

type ListenRequest_SubscribeRequest struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	Topics        []string               `protobuf:"bytes,1,rep,name=topics,proto3" json:"topics,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *ListenRequest_SubscribeRequest) Reset() {
	*x = ListenRequest_SubscribeRequest{}
	mi := &file_proto_command_queue_command_queue_proto_msgTypes[5]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *ListenRequest_SubscribeRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ListenRequest_SubscribeRequest) ProtoMessage() {}

func (x *ListenRequest_SubscribeRequest) ProtoReflect() protoreflect.Message {
	mi := &file_proto_command_queue_command_queue_proto_msgTypes[5]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ListenRequest_SubscribeRequest.ProtoReflect.Descriptor instead.
func (*ListenRequest_SubscribeRequest) Descriptor() ([]byte, []int) {
	return file_proto_command_queue_command_queue_proto_rawDescGZIP(), []int{0, 0}
}

func (x *ListenRequest_SubscribeRequest) GetTopics() []string {
	if x != nil {
		return x.Topics
	}
	return nil
}

type ListenRequest_AckRequest struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	MessageId     string                 `protobuf:"bytes,1,opt,name=message_id,json=messageId,proto3" json:"message_id,omitempty"`
	Nack          bool                   `protobuf:"varint,2,opt,name=nack,proto3" json:"nack,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *ListenRequest_AckRequest) Reset() {
	*x = ListenRequest_AckRequest{}
	mi := &file_proto_command_queue_command_queue_proto_msgTypes[6]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *ListenRequest_AckRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ListenRequest_AckRequest) ProtoMessage() {}

func (x *ListenRequest_AckRequest) ProtoReflect() protoreflect.Message {
	mi := &file_proto_command_queue_command_queue_proto_msgTypes[6]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ListenRequest_AckRequest.ProtoReflect.Descriptor instead.
func (*ListenRequest_AckRequest) Descriptor() ([]byte, []int) {
	return file_proto_command_queue_command_queue_proto_rawDescGZIP(), []int{0, 1}
}

func (x *ListenRequest_AckRequest) GetMessageId() string {
	if x != nil {
		return x.MessageId
	}
	return ""
}

func (x *ListenRequest_AckRequest) GetNack() bool {
	if x != nil {
		return x.Nack
	}
	return false
}

type ListenResponse_AckResponse struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	Id            string                 `protobuf:"bytes,1,opt,name=id,proto3" json:"id,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *ListenResponse_AckResponse) Reset() {
	*x = ListenResponse_AckResponse{}
	mi := &file_proto_command_queue_command_queue_proto_msgTypes[7]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *ListenResponse_AckResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ListenResponse_AckResponse) ProtoMessage() {}

func (x *ListenResponse_AckResponse) ProtoReflect() protoreflect.Message {
	mi := &file_proto_command_queue_command_queue_proto_msgTypes[7]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ListenResponse_AckResponse.ProtoReflect.Descriptor instead.
func (*ListenResponse_AckResponse) Descriptor() ([]byte, []int) {
	return file_proto_command_queue_command_queue_proto_rawDescGZIP(), []int{1, 0}
}

func (x *ListenResponse_AckResponse) GetId() string {
	if x != nil {
		return x.Id
	}
	return ""
}

var File_proto_command_queue_command_queue_proto protoreflect.FileDescriptor

var file_proto_command_queue_command_queue_proto_rawDesc = []byte{
	0x0a, 0x27, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2f, 0x63, 0x6f, 0x6d, 0x6d, 0x61, 0x6e, 0x64, 0x5f,
	0x71, 0x75, 0x65, 0x75, 0x65, 0x2f, 0x63, 0x6f, 0x6d, 0x6d, 0x61, 0x6e, 0x64, 0x5f, 0x71, 0x75,
	0x65, 0x75, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x16, 0x68, 0x65, 0x61, 0x74, 0x70,
	0x75, 0x6d, 0x70, 0x2e, 0x63, 0x6f, 0x6d, 0x6d, 0x61, 0x6e, 0x64, 0x5f, 0x71, 0x75, 0x65, 0x75,
	0x65, 0x22, 0xe2, 0x02, 0x0a, 0x0d, 0x4c, 0x69, 0x73, 0x74, 0x65, 0x6e, 0x52, 0x65, 0x71, 0x75,
	0x65, 0x73, 0x74, 0x12, 0x1d, 0x0a, 0x0a, 0x72, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x5f, 0x69,
	0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x09, 0x72, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74,
	0x49, 0x64, 0x12, 0x65, 0x0a, 0x11, 0x73, 0x75, 0x62, 0x73, 0x63, 0x72, 0x69, 0x62, 0x65, 0x5f,
	0x72, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x36, 0x2e,
	0x68, 0x65, 0x61, 0x74, 0x70, 0x75, 0x6d, 0x70, 0x2e, 0x63, 0x6f, 0x6d, 0x6d, 0x61, 0x6e, 0x64,
	0x5f, 0x71, 0x75, 0x65, 0x75, 0x65, 0x2e, 0x4c, 0x69, 0x73, 0x74, 0x65, 0x6e, 0x52, 0x65, 0x71,
	0x75, 0x65, 0x73, 0x74, 0x2e, 0x53, 0x75, 0x62, 0x73, 0x63, 0x72, 0x69, 0x62, 0x65, 0x52, 0x65,
	0x71, 0x75, 0x65, 0x73, 0x74, 0x48, 0x00, 0x52, 0x10, 0x73, 0x75, 0x62, 0x73, 0x63, 0x72, 0x69,
	0x62, 0x65, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x53, 0x0a, 0x0b, 0x61, 0x63, 0x6b,
	0x5f, 0x72, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x18, 0x03, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x30,
	0x2e, 0x68, 0x65, 0x61, 0x74, 0x70, 0x75, 0x6d, 0x70, 0x2e, 0x63, 0x6f, 0x6d, 0x6d, 0x61, 0x6e,
	0x64, 0x5f, 0x71, 0x75, 0x65, 0x75, 0x65, 0x2e, 0x4c, 0x69, 0x73, 0x74, 0x65, 0x6e, 0x52, 0x65,
	0x71, 0x75, 0x65, 0x73, 0x74, 0x2e, 0x41, 0x63, 0x6b, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74,
	0x48, 0x00, 0x52, 0x0a, 0x61, 0x63, 0x6b, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x2a,
	0x0a, 0x10, 0x53, 0x75, 0x62, 0x73, 0x63, 0x72, 0x69, 0x62, 0x65, 0x52, 0x65, 0x71, 0x75, 0x65,
	0x73, 0x74, 0x12, 0x16, 0x0a, 0x06, 0x74, 0x6f, 0x70, 0x69, 0x63, 0x73, 0x18, 0x01, 0x20, 0x03,
	0x28, 0x09, 0x52, 0x06, 0x74, 0x6f, 0x70, 0x69, 0x63, 0x73, 0x1a, 0x3f, 0x0a, 0x0a, 0x41, 0x63,
	0x6b, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x1d, 0x0a, 0x0a, 0x6d, 0x65, 0x73, 0x73,
	0x61, 0x67, 0x65, 0x5f, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x09, 0x6d, 0x65,
	0x73, 0x73, 0x61, 0x67, 0x65, 0x49, 0x64, 0x12, 0x12, 0x0a, 0x04, 0x6e, 0x61, 0x63, 0x6b, 0x18,
	0x02, 0x20, 0x01, 0x28, 0x08, 0x52, 0x04, 0x6e, 0x61, 0x63, 0x6b, 0x42, 0x09, 0x0a, 0x07, 0x72,
	0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x22, 0x89, 0x02, 0x0a, 0x0e, 0x4c, 0x69, 0x73, 0x74, 0x65,
	0x6e, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x1d, 0x0a, 0x0a, 0x72, 0x65, 0x71,
	0x75, 0x65, 0x73, 0x74, 0x5f, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x09, 0x72,
	0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x49, 0x64, 0x12, 0x54, 0x0a, 0x10, 0x6d, 0x65, 0x73, 0x73,
	0x61, 0x67, 0x65, 0x5f, 0x72, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x18, 0x02, 0x20, 0x01,
	0x28, 0x0b, 0x32, 0x27, 0x2e, 0x68, 0x65, 0x61, 0x74, 0x70, 0x75, 0x6d, 0x70, 0x2e, 0x63, 0x6f,
	0x6d, 0x6d, 0x61, 0x6e, 0x64, 0x5f, 0x71, 0x75, 0x65, 0x75, 0x65, 0x2e, 0x4d, 0x65, 0x73, 0x73,
	0x61, 0x67, 0x65, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x48, 0x00, 0x52, 0x0f, 0x6d,
	0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x57,
	0x0a, 0x0c, 0x61, 0x63, 0x6b, 0x5f, 0x72, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x18, 0x03,
	0x20, 0x01, 0x28, 0x0b, 0x32, 0x32, 0x2e, 0x68, 0x65, 0x61, 0x74, 0x70, 0x75, 0x6d, 0x70, 0x2e,
	0x63, 0x6f, 0x6d, 0x6d, 0x61, 0x6e, 0x64, 0x5f, 0x71, 0x75, 0x65, 0x75, 0x65, 0x2e, 0x4c, 0x69,
	0x73, 0x74, 0x65, 0x6e, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x2e, 0x41, 0x63, 0x6b,
	0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x48, 0x00, 0x52, 0x0b, 0x61, 0x63, 0x6b, 0x52,
	0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x1a, 0x1d, 0x0a, 0x0b, 0x41, 0x63, 0x6b, 0x52, 0x65,
	0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x0e, 0x0a, 0x02, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01,
	0x28, 0x09, 0x52, 0x02, 0x69, 0x64, 0x42, 0x0a, 0x0a, 0x08, 0x72, 0x65, 0x73, 0x70, 0x6f, 0x6e,
	0x73, 0x65, 0x22, 0x3b, 0x0a, 0x0f, 0x4d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x52, 0x65, 0x73,
	0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x0e, 0x0a, 0x02, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28,
	0x09, 0x52, 0x02, 0x69, 0x64, 0x12, 0x18, 0x0a, 0x07, 0x70, 0x61, 0x79, 0x6c, 0x6f, 0x61, 0x64,
	0x18, 0x02, 0x20, 0x01, 0x28, 0x0c, 0x52, 0x07, 0x70, 0x61, 0x79, 0x6c, 0x6f, 0x61, 0x64, 0x22,
	0x13, 0x0a, 0x11, 0x4c, 0x69, 0x73, 0x74, 0x54, 0x6f, 0x70, 0x69, 0x63, 0x73, 0x52, 0x65, 0x71,
	0x75, 0x65, 0x73, 0x74, 0x22, 0x2c, 0x0a, 0x12, 0x4c, 0x69, 0x73, 0x74, 0x54, 0x6f, 0x70, 0x69,
	0x63, 0x73, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x16, 0x0a, 0x06, 0x74, 0x6f,
	0x70, 0x69, 0x63, 0x73, 0x18, 0x01, 0x20, 0x03, 0x28, 0x09, 0x52, 0x06, 0x74, 0x6f, 0x70, 0x69,
	0x63, 0x73, 0x32, 0xdb, 0x01, 0x0a, 0x13, 0x43, 0x6f, 0x6d, 0x6d, 0x61, 0x6e, 0x64, 0x51, 0x75,
	0x65, 0x75, 0x65, 0x53, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x12, 0x5d, 0x0a, 0x06, 0x4c, 0x69,
	0x73, 0x74, 0x65, 0x6e, 0x12, 0x25, 0x2e, 0x68, 0x65, 0x61, 0x74, 0x70, 0x75, 0x6d, 0x70, 0x2e,
	0x63, 0x6f, 0x6d, 0x6d, 0x61, 0x6e, 0x64, 0x5f, 0x71, 0x75, 0x65, 0x75, 0x65, 0x2e, 0x4c, 0x69,
	0x73, 0x74, 0x65, 0x6e, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x26, 0x2e, 0x68, 0x65,
	0x61, 0x74, 0x70, 0x75, 0x6d, 0x70, 0x2e, 0x63, 0x6f, 0x6d, 0x6d, 0x61, 0x6e, 0x64, 0x5f, 0x71,
	0x75, 0x65, 0x75, 0x65, 0x2e, 0x4c, 0x69, 0x73, 0x74, 0x65, 0x6e, 0x52, 0x65, 0x73, 0x70, 0x6f,
	0x6e, 0x73, 0x65, 0x22, 0x00, 0x28, 0x01, 0x30, 0x01, 0x12, 0x65, 0x0a, 0x0a, 0x4c, 0x69, 0x73,
	0x74, 0x54, 0x6f, 0x70, 0x69, 0x63, 0x73, 0x12, 0x29, 0x2e, 0x68, 0x65, 0x61, 0x74, 0x70, 0x75,
	0x6d, 0x70, 0x2e, 0x63, 0x6f, 0x6d, 0x6d, 0x61, 0x6e, 0x64, 0x5f, 0x71, 0x75, 0x65, 0x75, 0x65,
	0x2e, 0x4c, 0x69, 0x73, 0x74, 0x54, 0x6f, 0x70, 0x69, 0x63, 0x73, 0x52, 0x65, 0x71, 0x75, 0x65,
	0x73, 0x74, 0x1a, 0x2a, 0x2e, 0x68, 0x65, 0x61, 0x74, 0x70, 0x75, 0x6d, 0x70, 0x2e, 0x63, 0x6f,
	0x6d, 0x6d, 0x61, 0x6e, 0x64, 0x5f, 0x71, 0x75, 0x65, 0x75, 0x65, 0x2e, 0x4c, 0x69, 0x73, 0x74,
	0x54, 0x6f, 0x70, 0x69, 0x63, 0x73, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x22, 0x00,
	0x42, 0x33, 0x5a, 0x31, 0x67, 0x69, 0x74, 0x68, 0x75, 0x62, 0x2e, 0x63, 0x6f, 0x6d, 0x2f, 0x67,
	0x6f, 0x6e, 0x7a, 0x6f, 0x6a, 0x69, 0x76, 0x65, 0x2f, 0x68, 0x65, 0x61, 0x74, 0x70, 0x75, 0x6d,
	0x70, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2f, 0x63, 0x6f, 0x6d, 0x6d, 0x61, 0x6e, 0x64, 0x5f,
	0x71, 0x75, 0x65, 0x75, 0x65, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_proto_command_queue_command_queue_proto_rawDescOnce sync.Once
	file_proto_command_queue_command_queue_proto_rawDescData = file_proto_command_queue_command_queue_proto_rawDesc
)

func file_proto_command_queue_command_queue_proto_rawDescGZIP() []byte {
	file_proto_command_queue_command_queue_proto_rawDescOnce.Do(func() {
		file_proto_command_queue_command_queue_proto_rawDescData = protoimpl.X.CompressGZIP(file_proto_command_queue_command_queue_proto_rawDescData)
	})
	return file_proto_command_queue_command_queue_proto_rawDescData
}

var file_proto_command_queue_command_queue_proto_msgTypes = make([]protoimpl.MessageInfo, 8)
var file_proto_command_queue_command_queue_proto_goTypes = []any{
	(*ListenRequest)(nil),                  // 0: heatpump.command_queue.ListenRequest
	(*ListenResponse)(nil),                 // 1: heatpump.command_queue.ListenResponse
	(*MessageResponse)(nil),                // 2: heatpump.command_queue.MessageResponse
	(*ListTopicsRequest)(nil),              // 3: heatpump.command_queue.ListTopicsRequest
	(*ListTopicsResponse)(nil),             // 4: heatpump.command_queue.ListTopicsResponse
	(*ListenRequest_SubscribeRequest)(nil), // 5: heatpump.command_queue.ListenRequest.SubscribeRequest
	(*ListenRequest_AckRequest)(nil),       // 6: heatpump.command_queue.ListenRequest.AckRequest
	(*ListenResponse_AckResponse)(nil),     // 7: heatpump.command_queue.ListenResponse.AckResponse
}
var file_proto_command_queue_command_queue_proto_depIdxs = []int32{
	5, // 0: heatpump.command_queue.ListenRequest.subscribe_request:type_name -> heatpump.command_queue.ListenRequest.SubscribeRequest
	6, // 1: heatpump.command_queue.ListenRequest.ack_request:type_name -> heatpump.command_queue.ListenRequest.AckRequest
	2, // 2: heatpump.command_queue.ListenResponse.message_response:type_name -> heatpump.command_queue.MessageResponse
	7, // 3: heatpump.command_queue.ListenResponse.ack_response:type_name -> heatpump.command_queue.ListenResponse.AckResponse
	0, // 4: heatpump.command_queue.CommandQueueService.Listen:input_type -> heatpump.command_queue.ListenRequest
	3, // 5: heatpump.command_queue.CommandQueueService.ListTopics:input_type -> heatpump.command_queue.ListTopicsRequest
	1, // 6: heatpump.command_queue.CommandQueueService.Listen:output_type -> heatpump.command_queue.ListenResponse
	4, // 7: heatpump.command_queue.CommandQueueService.ListTopics:output_type -> heatpump.command_queue.ListTopicsResponse
	6, // [6:8] is the sub-list for method output_type
	4, // [4:6] is the sub-list for method input_type
	4, // [4:4] is the sub-list for extension type_name
	4, // [4:4] is the sub-list for extension extendee
	0, // [0:4] is the sub-list for field type_name
}

func init() { file_proto_command_queue_command_queue_proto_init() }
func file_proto_command_queue_command_queue_proto_init() {
	if File_proto_command_queue_command_queue_proto != nil {
		return
	}
	file_proto_command_queue_command_queue_proto_msgTypes[0].OneofWrappers = []any{
		(*ListenRequest_SubscribeRequest_)(nil),
		(*ListenRequest_AckRequest_)(nil),
	}
	file_proto_command_queue_command_queue_proto_msgTypes[1].OneofWrappers = []any{
		(*ListenResponse_MessageResponse)(nil),
		(*ListenResponse_AckResponse_)(nil),
	}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_proto_command_queue_command_queue_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   8,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_proto_command_queue_command_queue_proto_goTypes,
		DependencyIndexes: file_proto_command_queue_command_queue_proto_depIdxs,
		MessageInfos:      file_proto_command_queue_command_queue_proto_msgTypes,
	}.Build()
	File_proto_command_queue_command_queue_proto = out.File
	file_proto_command_queue_command_queue_proto_rawDesc = nil
	file_proto_command_queue_command_queue_proto_goTypes = nil
	file_proto_command_queue_command_queue_proto_depIdxs = nil
}

// Reference imports to suppress errors if they are not otherwise used.
var _ context.Context
var _ grpc.ClientConnInterface

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
const _ = grpc.SupportPackageIsVersion6

// CommandQueueServiceClient is the client API for CommandQueueService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://godoc.org/google.golang.org/grpc#ClientConn.NewStream.
type CommandQueueServiceClient interface {
	Listen(ctx context.Context, opts ...grpc.CallOption) (CommandQueueService_ListenClient, error)
	ListTopics(ctx context.Context, in *ListTopicsRequest, opts ...grpc.CallOption) (*ListTopicsResponse, error)
}

type commandQueueServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewCommandQueueServiceClient(cc grpc.ClientConnInterface) CommandQueueServiceClient {
	return &commandQueueServiceClient{cc}
}

func (c *commandQueueServiceClient) Listen(ctx context.Context, opts ...grpc.CallOption) (CommandQueueService_ListenClient, error) {
	stream, err := c.cc.NewStream(ctx, &_CommandQueueService_serviceDesc.Streams[0], "/heatpump.command_queue.CommandQueueService/Listen", opts...)
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
type CommandQueueServiceServer interface {
	Listen(CommandQueueService_ListenServer) error
	ListTopics(context.Context, *ListTopicsRequest) (*ListTopicsResponse, error)
}

// UnimplementedCommandQueueServiceServer can be embedded to have forward compatible implementations.
type UnimplementedCommandQueueServiceServer struct {
}

func (*UnimplementedCommandQueueServiceServer) Listen(CommandQueueService_ListenServer) error {
	return status.Errorf(codes.Unimplemented, "method Listen not implemented")
}
func (*UnimplementedCommandQueueServiceServer) ListTopics(context.Context, *ListTopicsRequest) (*ListTopicsResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ListTopics not implemented")
}

func RegisterCommandQueueServiceServer(s *grpc.Server, srv CommandQueueServiceServer) {
	s.RegisterService(&_CommandQueueService_serviceDesc, srv)
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

var _CommandQueueService_serviceDesc = grpc.ServiceDesc{
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
	Metadata: "proto/command_queue/command_queue.proto",
}
