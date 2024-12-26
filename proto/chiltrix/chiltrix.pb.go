// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.34.2
// 	protoc        v5.29.0
// source: proto/chiltrix/chiltrix.proto

package chiltrix

import (
	context "context"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
	timestamppb "google.golang.org/protobuf/types/known/timestamppb"
	reflect "reflect"
	sync "sync"
)

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

type State struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	CollectionTime *timestamppb.Timestamp `protobuf:"bytes,1,opt,name=collection_time,json=collectionTime,proto3" json:"collection_time,omitempty"`
	RegisterValues *RegisterSnapshot      `protobuf:"bytes,2,opt,name=register_values,json=registerValues,proto3" json:"register_values,omitempty"`
}

func (x *State) Reset() {
	*x = State{}
	if protoimpl.UnsafeEnabled {
		mi := &file_proto_chiltrix_chiltrix_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *State) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*State) ProtoMessage() {}

func (x *State) ProtoReflect() protoreflect.Message {
	mi := &file_proto_chiltrix_chiltrix_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use State.ProtoReflect.Descriptor instead.
func (*State) Descriptor() ([]byte, []int) {
	return file_proto_chiltrix_chiltrix_proto_rawDescGZIP(), []int{0}
}

func (x *State) GetCollectionTime() *timestamppb.Timestamp {
	if x != nil {
		return x.CollectionTime
	}
	return nil
}

func (x *State) GetRegisterValues() *RegisterSnapshot {
	if x != nil {
		return x.RegisterValues
	}
	return nil
}

type RegisterSnapshot struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	HoldingRegisterValues map[uint32]uint32 `protobuf:"bytes,1,rep,name=holding_register_values,json=holdingRegisterValues,proto3" json:"holding_register_values,omitempty" protobuf_key:"varint,1,opt,name=key,proto3" protobuf_val:"varint,2,opt,name=value,proto3"`
}

func (x *RegisterSnapshot) Reset() {
	*x = RegisterSnapshot{}
	if protoimpl.UnsafeEnabled {
		mi := &file_proto_chiltrix_chiltrix_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *RegisterSnapshot) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*RegisterSnapshot) ProtoMessage() {}

func (x *RegisterSnapshot) ProtoReflect() protoreflect.Message {
	mi := &file_proto_chiltrix_chiltrix_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use RegisterSnapshot.ProtoReflect.Descriptor instead.
func (*RegisterSnapshot) Descriptor() ([]byte, []int) {
	return file_proto_chiltrix_chiltrix_proto_rawDescGZIP(), []int{1}
}

func (x *RegisterSnapshot) GetHoldingRegisterValues() map[uint32]uint32 {
	if x != nil {
		return x.HoldingRegisterValues
	}
	return nil
}

type GetStateRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields
}

func (x *GetStateRequest) Reset() {
	*x = GetStateRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_proto_chiltrix_chiltrix_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *GetStateRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GetStateRequest) ProtoMessage() {}

func (x *GetStateRequest) ProtoReflect() protoreflect.Message {
	mi := &file_proto_chiltrix_chiltrix_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use GetStateRequest.ProtoReflect.Descriptor instead.
func (*GetStateRequest) Descriptor() ([]byte, []int) {
	return file_proto_chiltrix_chiltrix_proto_rawDescGZIP(), []int{2}
}

type QueryStreamRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	StartTime *timestamppb.Timestamp `protobuf:"bytes,1,opt,name=start_time,json=startTime,proto3" json:"start_time,omitempty"`
	EndTime   *timestamppb.Timestamp `protobuf:"bytes,2,opt,name=end_time,json=endTime,proto3" json:"end_time,omitempty"`
}

func (x *QueryStreamRequest) Reset() {
	*x = QueryStreamRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_proto_chiltrix_chiltrix_proto_msgTypes[3]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *QueryStreamRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*QueryStreamRequest) ProtoMessage() {}

func (x *QueryStreamRequest) ProtoReflect() protoreflect.Message {
	mi := &file_proto_chiltrix_chiltrix_proto_msgTypes[3]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use QueryStreamRequest.ProtoReflect.Descriptor instead.
func (*QueryStreamRequest) Descriptor() ([]byte, []int) {
	return file_proto_chiltrix_chiltrix_proto_rawDescGZIP(), []int{3}
}

func (x *QueryStreamRequest) GetStartTime() *timestamppb.Timestamp {
	if x != nil {
		return x.StartTime
	}
	return nil
}

func (x *QueryStreamRequest) GetEndTime() *timestamppb.Timestamp {
	if x != nil {
		return x.EndTime
	}
	return nil
}

type QueryStreamResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	State *State `protobuf:"bytes,1,opt,name=state,proto3" json:"state,omitempty"`
}

func (x *QueryStreamResponse) Reset() {
	*x = QueryStreamResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_proto_chiltrix_chiltrix_proto_msgTypes[4]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *QueryStreamResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*QueryStreamResponse) ProtoMessage() {}

func (x *QueryStreamResponse) ProtoReflect() protoreflect.Message {
	mi := &file_proto_chiltrix_chiltrix_proto_msgTypes[4]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use QueryStreamResponse.ProtoReflect.Descriptor instead.
func (*QueryStreamResponse) Descriptor() ([]byte, []int) {
	return file_proto_chiltrix_chiltrix_proto_rawDescGZIP(), []int{4}
}

func (x *QueryStreamResponse) GetState() *State {
	if x != nil {
		return x.State
	}
	return nil
}

type StateSequence struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Diffs []*StateDiff `protobuf:"bytes,1,rep,name=diffs,proto3" json:"diffs,omitempty"`
}

func (x *StateSequence) Reset() {
	*x = StateSequence{}
	if protoimpl.UnsafeEnabled {
		mi := &file_proto_chiltrix_chiltrix_proto_msgTypes[5]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *StateSequence) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*StateSequence) ProtoMessage() {}

func (x *StateSequence) ProtoReflect() protoreflect.Message {
	mi := &file_proto_chiltrix_chiltrix_proto_msgTypes[5]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use StateSequence.ProtoReflect.Descriptor instead.
func (*StateSequence) Descriptor() ([]byte, []int) {
	return file_proto_chiltrix_chiltrix_proto_rawDescGZIP(), []int{5}
}

func (x *StateSequence) GetDiffs() []*StateDiff {
	if x != nil {
		return x.Diffs
	}
	return nil
}

type StateDiff struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	CollectionTime   *timestamppb.Timestamp `protobuf:"bytes,1,opt,name=collection_time,json=collectionTime,proto3" json:"collection_time,omitempty"`
	UpdatedValues    *RegisterSnapshot      `protobuf:"bytes,2,opt,name=updated_values,json=updatedValues,proto3" json:"updated_values,omitempty"`
	DeletedRegisters []uint32               `protobuf:"varint,3,rep,packed,name=deleted_registers,json=deletedRegisters,proto3" json:"deleted_registers,omitempty"`
}

func (x *StateDiff) Reset() {
	*x = StateDiff{}
	if protoimpl.UnsafeEnabled {
		mi := &file_proto_chiltrix_chiltrix_proto_msgTypes[6]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *StateDiff) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*StateDiff) ProtoMessage() {}

func (x *StateDiff) ProtoReflect() protoreflect.Message {
	mi := &file_proto_chiltrix_chiltrix_proto_msgTypes[6]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use StateDiff.ProtoReflect.Descriptor instead.
func (*StateDiff) Descriptor() ([]byte, []int) {
	return file_proto_chiltrix_chiltrix_proto_rawDescGZIP(), []int{6}
}

func (x *StateDiff) GetCollectionTime() *timestamppb.Timestamp {
	if x != nil {
		return x.CollectionTime
	}
	return nil
}

func (x *StateDiff) GetUpdatedValues() *RegisterSnapshot {
	if x != nil {
		return x.UpdatedValues
	}
	return nil
}

func (x *StateDiff) GetDeletedRegisters() []uint32 {
	if x != nil {
		return x.DeletedRegisters
	}
	return nil
}

type SetParameterRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	TargetHeatingModeTemperature *Temperature   `protobuf:"bytes,1,opt,name=target_heating_mode_temperature,json=targetHeatingModeTemperature,proto3" json:"target_heating_mode_temperature,omitempty"`
	RegisterValue                *RegisterValue `protobuf:"bytes,2,opt,name=register_value,json=registerValue,proto3" json:"register_value,omitempty"`
}

func (x *SetParameterRequest) Reset() {
	*x = SetParameterRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_proto_chiltrix_chiltrix_proto_msgTypes[7]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *SetParameterRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*SetParameterRequest) ProtoMessage() {}

func (x *SetParameterRequest) ProtoReflect() protoreflect.Message {
	mi := &file_proto_chiltrix_chiltrix_proto_msgTypes[7]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use SetParameterRequest.ProtoReflect.Descriptor instead.
func (*SetParameterRequest) Descriptor() ([]byte, []int) {
	return file_proto_chiltrix_chiltrix_proto_rawDescGZIP(), []int{7}
}

func (x *SetParameterRequest) GetTargetHeatingModeTemperature() *Temperature {
	if x != nil {
		return x.TargetHeatingModeTemperature
	}
	return nil
}

func (x *SetParameterRequest) GetRegisterValue() *RegisterValue {
	if x != nil {
		return x.RegisterValue
	}
	return nil
}

type RegisterValue struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Register uint32 `protobuf:"varint,1,opt,name=register,proto3" json:"register,omitempty"`
	Value    uint32 `protobuf:"varint,2,opt,name=value,proto3" json:"value,omitempty"`
}

func (x *RegisterValue) Reset() {
	*x = RegisterValue{}
	if protoimpl.UnsafeEnabled {
		mi := &file_proto_chiltrix_chiltrix_proto_msgTypes[8]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *RegisterValue) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*RegisterValue) ProtoMessage() {}

func (x *RegisterValue) ProtoReflect() protoreflect.Message {
	mi := &file_proto_chiltrix_chiltrix_proto_msgTypes[8]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use RegisterValue.ProtoReflect.Descriptor instead.
func (*RegisterValue) Descriptor() ([]byte, []int) {
	return file_proto_chiltrix_chiltrix_proto_rawDescGZIP(), []int{8}
}

func (x *RegisterValue) GetRegister() uint32 {
	if x != nil {
		return x.Register
	}
	return 0
}

func (x *RegisterValue) GetValue() uint32 {
	if x != nil {
		return x.Value
	}
	return 0
}

type SetParameterResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields
}

func (x *SetParameterResponse) Reset() {
	*x = SetParameterResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_proto_chiltrix_chiltrix_proto_msgTypes[9]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *SetParameterResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*SetParameterResponse) ProtoMessage() {}

func (x *SetParameterResponse) ProtoReflect() protoreflect.Message {
	mi := &file_proto_chiltrix_chiltrix_proto_msgTypes[9]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use SetParameterResponse.ProtoReflect.Descriptor instead.
func (*SetParameterResponse) Descriptor() ([]byte, []int) {
	return file_proto_chiltrix_chiltrix_proto_rawDescGZIP(), []int{9}
}

type Temperature struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	DegreesCelcius float64 `protobuf:"fixed64,1,opt,name=degrees_celcius,json=degreesCelcius,proto3" json:"degrees_celcius,omitempty"`
}

func (x *Temperature) Reset() {
	*x = Temperature{}
	if protoimpl.UnsafeEnabled {
		mi := &file_proto_chiltrix_chiltrix_proto_msgTypes[10]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Temperature) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Temperature) ProtoMessage() {}

func (x *Temperature) ProtoReflect() protoreflect.Message {
	mi := &file_proto_chiltrix_chiltrix_proto_msgTypes[10]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Temperature.ProtoReflect.Descriptor instead.
func (*Temperature) Descriptor() ([]byte, []int) {
	return file_proto_chiltrix_chiltrix_proto_rawDescGZIP(), []int{10}
}

func (x *Temperature) GetDegreesCelcius() float64 {
	if x != nil {
		return x.DegreesCelcius
	}
	return 0
}

var File_proto_chiltrix_chiltrix_proto protoreflect.FileDescriptor

var file_proto_chiltrix_chiltrix_proto_rawDesc = []byte{
	0x0a, 0x1d, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2f, 0x63, 0x68, 0x69, 0x6c, 0x74, 0x72, 0x69, 0x78,
	0x2f, 0x63, 0x68, 0x69, 0x6c, 0x74, 0x72, 0x69, 0x78, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12,
	0x08, 0x63, 0x68, 0x69, 0x6c, 0x74, 0x72, 0x69, 0x78, 0x1a, 0x1f, 0x67, 0x6f, 0x6f, 0x67, 0x6c,
	0x65, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2f, 0x74, 0x69, 0x6d, 0x65, 0x73,
	0x74, 0x61, 0x6d, 0x70, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x22, 0x91, 0x01, 0x0a, 0x05, 0x53,
	0x74, 0x61, 0x74, 0x65, 0x12, 0x43, 0x0a, 0x0f, 0x63, 0x6f, 0x6c, 0x6c, 0x65, 0x63, 0x74, 0x69,
	0x6f, 0x6e, 0x5f, 0x74, 0x69, 0x6d, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x1a, 0x2e,
	0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e,
	0x54, 0x69, 0x6d, 0x65, 0x73, 0x74, 0x61, 0x6d, 0x70, 0x52, 0x0e, 0x63, 0x6f, 0x6c, 0x6c, 0x65,
	0x63, 0x74, 0x69, 0x6f, 0x6e, 0x54, 0x69, 0x6d, 0x65, 0x12, 0x43, 0x0a, 0x0f, 0x72, 0x65, 0x67,
	0x69, 0x73, 0x74, 0x65, 0x72, 0x5f, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x73, 0x18, 0x02, 0x20, 0x01,
	0x28, 0x0b, 0x32, 0x1a, 0x2e, 0x63, 0x68, 0x69, 0x6c, 0x74, 0x72, 0x69, 0x78, 0x2e, 0x52, 0x65,
	0x67, 0x69, 0x73, 0x74, 0x65, 0x72, 0x53, 0x6e, 0x61, 0x70, 0x73, 0x68, 0x6f, 0x74, 0x52, 0x0e,
	0x72, 0x65, 0x67, 0x69, 0x73, 0x74, 0x65, 0x72, 0x56, 0x61, 0x6c, 0x75, 0x65, 0x73, 0x22, 0xcb,
	0x01, 0x0a, 0x10, 0x52, 0x65, 0x67, 0x69, 0x73, 0x74, 0x65, 0x72, 0x53, 0x6e, 0x61, 0x70, 0x73,
	0x68, 0x6f, 0x74, 0x12, 0x6d, 0x0a, 0x17, 0x68, 0x6f, 0x6c, 0x64, 0x69, 0x6e, 0x67, 0x5f, 0x72,
	0x65, 0x67, 0x69, 0x73, 0x74, 0x65, 0x72, 0x5f, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x73, 0x18, 0x01,
	0x20, 0x03, 0x28, 0x0b, 0x32, 0x35, 0x2e, 0x63, 0x68, 0x69, 0x6c, 0x74, 0x72, 0x69, 0x78, 0x2e,
	0x52, 0x65, 0x67, 0x69, 0x73, 0x74, 0x65, 0x72, 0x53, 0x6e, 0x61, 0x70, 0x73, 0x68, 0x6f, 0x74,
	0x2e, 0x48, 0x6f, 0x6c, 0x64, 0x69, 0x6e, 0x67, 0x52, 0x65, 0x67, 0x69, 0x73, 0x74, 0x65, 0x72,
	0x56, 0x61, 0x6c, 0x75, 0x65, 0x73, 0x45, 0x6e, 0x74, 0x72, 0x79, 0x52, 0x15, 0x68, 0x6f, 0x6c,
	0x64, 0x69, 0x6e, 0x67, 0x52, 0x65, 0x67, 0x69, 0x73, 0x74, 0x65, 0x72, 0x56, 0x61, 0x6c, 0x75,
	0x65, 0x73, 0x1a, 0x48, 0x0a, 0x1a, 0x48, 0x6f, 0x6c, 0x64, 0x69, 0x6e, 0x67, 0x52, 0x65, 0x67,
	0x69, 0x73, 0x74, 0x65, 0x72, 0x56, 0x61, 0x6c, 0x75, 0x65, 0x73, 0x45, 0x6e, 0x74, 0x72, 0x79,
	0x12, 0x10, 0x0a, 0x03, 0x6b, 0x65, 0x79, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0d, 0x52, 0x03, 0x6b,
	0x65, 0x79, 0x12, 0x14, 0x0a, 0x05, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28,
	0x0d, 0x52, 0x05, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x3a, 0x02, 0x38, 0x01, 0x22, 0x11, 0x0a, 0x0f,
	0x47, 0x65, 0x74, 0x53, 0x74, 0x61, 0x74, 0x65, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x22,
	0x86, 0x01, 0x0a, 0x12, 0x51, 0x75, 0x65, 0x72, 0x79, 0x53, 0x74, 0x72, 0x65, 0x61, 0x6d, 0x52,
	0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x39, 0x0a, 0x0a, 0x73, 0x74, 0x61, 0x72, 0x74, 0x5f,
	0x74, 0x69, 0x6d, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x1a, 0x2e, 0x67, 0x6f, 0x6f,
	0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x54, 0x69, 0x6d,
	0x65, 0x73, 0x74, 0x61, 0x6d, 0x70, 0x52, 0x09, 0x73, 0x74, 0x61, 0x72, 0x74, 0x54, 0x69, 0x6d,
	0x65, 0x12, 0x35, 0x0a, 0x08, 0x65, 0x6e, 0x64, 0x5f, 0x74, 0x69, 0x6d, 0x65, 0x18, 0x02, 0x20,
	0x01, 0x28, 0x0b, 0x32, 0x1a, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f,
	0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x54, 0x69, 0x6d, 0x65, 0x73, 0x74, 0x61, 0x6d, 0x70, 0x52,
	0x07, 0x65, 0x6e, 0x64, 0x54, 0x69, 0x6d, 0x65, 0x22, 0x3c, 0x0a, 0x13, 0x51, 0x75, 0x65, 0x72,
	0x79, 0x53, 0x74, 0x72, 0x65, 0x61, 0x6d, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12,
	0x25, 0x0a, 0x05, 0x73, 0x74, 0x61, 0x74, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x0f,
	0x2e, 0x63, 0x68, 0x69, 0x6c, 0x74, 0x72, 0x69, 0x78, 0x2e, 0x53, 0x74, 0x61, 0x74, 0x65, 0x52,
	0x05, 0x73, 0x74, 0x61, 0x74, 0x65, 0x22, 0x3a, 0x0a, 0x0d, 0x53, 0x74, 0x61, 0x74, 0x65, 0x53,
	0x65, 0x71, 0x75, 0x65, 0x6e, 0x63, 0x65, 0x12, 0x29, 0x0a, 0x05, 0x64, 0x69, 0x66, 0x66, 0x73,
	0x18, 0x01, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x13, 0x2e, 0x63, 0x68, 0x69, 0x6c, 0x74, 0x72, 0x69,
	0x78, 0x2e, 0x53, 0x74, 0x61, 0x74, 0x65, 0x44, 0x69, 0x66, 0x66, 0x52, 0x05, 0x64, 0x69, 0x66,
	0x66, 0x73, 0x22, 0xc0, 0x01, 0x0a, 0x09, 0x53, 0x74, 0x61, 0x74, 0x65, 0x44, 0x69, 0x66, 0x66,
	0x12, 0x43, 0x0a, 0x0f, 0x63, 0x6f, 0x6c, 0x6c, 0x65, 0x63, 0x74, 0x69, 0x6f, 0x6e, 0x5f, 0x74,
	0x69, 0x6d, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x1a, 0x2e, 0x67, 0x6f, 0x6f, 0x67,
	0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x54, 0x69, 0x6d, 0x65,
	0x73, 0x74, 0x61, 0x6d, 0x70, 0x52, 0x0e, 0x63, 0x6f, 0x6c, 0x6c, 0x65, 0x63, 0x74, 0x69, 0x6f,
	0x6e, 0x54, 0x69, 0x6d, 0x65, 0x12, 0x41, 0x0a, 0x0e, 0x75, 0x70, 0x64, 0x61, 0x74, 0x65, 0x64,
	0x5f, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x73, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x1a, 0x2e,
	0x63, 0x68, 0x69, 0x6c, 0x74, 0x72, 0x69, 0x78, 0x2e, 0x52, 0x65, 0x67, 0x69, 0x73, 0x74, 0x65,
	0x72, 0x53, 0x6e, 0x61, 0x70, 0x73, 0x68, 0x6f, 0x74, 0x52, 0x0d, 0x75, 0x70, 0x64, 0x61, 0x74,
	0x65, 0x64, 0x56, 0x61, 0x6c, 0x75, 0x65, 0x73, 0x12, 0x2b, 0x0a, 0x11, 0x64, 0x65, 0x6c, 0x65,
	0x74, 0x65, 0x64, 0x5f, 0x72, 0x65, 0x67, 0x69, 0x73, 0x74, 0x65, 0x72, 0x73, 0x18, 0x03, 0x20,
	0x03, 0x28, 0x0d, 0x52, 0x10, 0x64, 0x65, 0x6c, 0x65, 0x74, 0x65, 0x64, 0x52, 0x65, 0x67, 0x69,
	0x73, 0x74, 0x65, 0x72, 0x73, 0x22, 0xb3, 0x01, 0x0a, 0x13, 0x53, 0x65, 0x74, 0x50, 0x61, 0x72,
	0x61, 0x6d, 0x65, 0x74, 0x65, 0x72, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x5c, 0x0a,
	0x1f, 0x74, 0x61, 0x72, 0x67, 0x65, 0x74, 0x5f, 0x68, 0x65, 0x61, 0x74, 0x69, 0x6e, 0x67, 0x5f,
	0x6d, 0x6f, 0x64, 0x65, 0x5f, 0x74, 0x65, 0x6d, 0x70, 0x65, 0x72, 0x61, 0x74, 0x75, 0x72, 0x65,
	0x18, 0x01, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x15, 0x2e, 0x63, 0x68, 0x69, 0x6c, 0x74, 0x72, 0x69,
	0x78, 0x2e, 0x54, 0x65, 0x6d, 0x70, 0x65, 0x72, 0x61, 0x74, 0x75, 0x72, 0x65, 0x52, 0x1c, 0x74,
	0x61, 0x72, 0x67, 0x65, 0x74, 0x48, 0x65, 0x61, 0x74, 0x69, 0x6e, 0x67, 0x4d, 0x6f, 0x64, 0x65,
	0x54, 0x65, 0x6d, 0x70, 0x65, 0x72, 0x61, 0x74, 0x75, 0x72, 0x65, 0x12, 0x3e, 0x0a, 0x0e, 0x72,
	0x65, 0x67, 0x69, 0x73, 0x74, 0x65, 0x72, 0x5f, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x18, 0x02, 0x20,
	0x01, 0x28, 0x0b, 0x32, 0x17, 0x2e, 0x63, 0x68, 0x69, 0x6c, 0x74, 0x72, 0x69, 0x78, 0x2e, 0x52,
	0x65, 0x67, 0x69, 0x73, 0x74, 0x65, 0x72, 0x56, 0x61, 0x6c, 0x75, 0x65, 0x52, 0x0d, 0x72, 0x65,
	0x67, 0x69, 0x73, 0x74, 0x65, 0x72, 0x56, 0x61, 0x6c, 0x75, 0x65, 0x22, 0x41, 0x0a, 0x0d, 0x52,
	0x65, 0x67, 0x69, 0x73, 0x74, 0x65, 0x72, 0x56, 0x61, 0x6c, 0x75, 0x65, 0x12, 0x1a, 0x0a, 0x08,
	0x72, 0x65, 0x67, 0x69, 0x73, 0x74, 0x65, 0x72, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0d, 0x52, 0x08,
	0x72, 0x65, 0x67, 0x69, 0x73, 0x74, 0x65, 0x72, 0x12, 0x14, 0x0a, 0x05, 0x76, 0x61, 0x6c, 0x75,
	0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0d, 0x52, 0x05, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x22, 0x16,
	0x0a, 0x14, 0x53, 0x65, 0x74, 0x50, 0x61, 0x72, 0x61, 0x6d, 0x65, 0x74, 0x65, 0x72, 0x52, 0x65,
	0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x22, 0x36, 0x0a, 0x0b, 0x54, 0x65, 0x6d, 0x70, 0x65, 0x72,
	0x61, 0x74, 0x75, 0x72, 0x65, 0x12, 0x27, 0x0a, 0x0f, 0x64, 0x65, 0x67, 0x72, 0x65, 0x65, 0x73,
	0x5f, 0x63, 0x65, 0x6c, 0x63, 0x69, 0x75, 0x73, 0x18, 0x01, 0x20, 0x01, 0x28, 0x01, 0x52, 0x0e,
	0x64, 0x65, 0x67, 0x72, 0x65, 0x65, 0x73, 0x43, 0x65, 0x6c, 0x63, 0x69, 0x75, 0x73, 0x32, 0x5b,
	0x0a, 0x09, 0x48, 0x69, 0x73, 0x74, 0x6f, 0x72, 0x69, 0x61, 0x6e, 0x12, 0x4e, 0x0a, 0x0b, 0x51,
	0x75, 0x65, 0x72, 0x79, 0x53, 0x74, 0x72, 0x65, 0x61, 0x6d, 0x12, 0x1c, 0x2e, 0x63, 0x68, 0x69,
	0x6c, 0x74, 0x72, 0x69, 0x78, 0x2e, 0x51, 0x75, 0x65, 0x72, 0x79, 0x53, 0x74, 0x72, 0x65, 0x61,
	0x6d, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x1d, 0x2e, 0x63, 0x68, 0x69, 0x6c, 0x74,
	0x72, 0x69, 0x78, 0x2e, 0x51, 0x75, 0x65, 0x72, 0x79, 0x53, 0x74, 0x72, 0x65, 0x61, 0x6d, 0x52,
	0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x22, 0x00, 0x30, 0x01, 0x32, 0x9d, 0x01, 0x0a, 0x10,
	0x52, 0x65, 0x61, 0x64, 0x57, 0x72, 0x69, 0x74, 0x65, 0x53, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65,
	0x12, 0x4f, 0x0a, 0x0c, 0x53, 0x65, 0x74, 0x50, 0x61, 0x72, 0x61, 0x6d, 0x65, 0x74, 0x65, 0x72,
	0x12, 0x1d, 0x2e, 0x63, 0x68, 0x69, 0x6c, 0x74, 0x72, 0x69, 0x78, 0x2e, 0x53, 0x65, 0x74, 0x50,
	0x61, 0x72, 0x61, 0x6d, 0x65, 0x74, 0x65, 0x72, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a,
	0x1e, 0x2e, 0x63, 0x68, 0x69, 0x6c, 0x74, 0x72, 0x69, 0x78, 0x2e, 0x53, 0x65, 0x74, 0x50, 0x61,
	0x72, 0x61, 0x6d, 0x65, 0x74, 0x65, 0x72, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x22,
	0x00, 0x12, 0x38, 0x0a, 0x08, 0x47, 0x65, 0x74, 0x53, 0x74, 0x61, 0x74, 0x65, 0x12, 0x19, 0x2e,
	0x63, 0x68, 0x69, 0x6c, 0x74, 0x72, 0x69, 0x78, 0x2e, 0x47, 0x65, 0x74, 0x53, 0x74, 0x61, 0x74,
	0x65, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x0f, 0x2e, 0x63, 0x68, 0x69, 0x6c, 0x74,
	0x72, 0x69, 0x78, 0x2e, 0x53, 0x74, 0x61, 0x74, 0x65, 0x22, 0x00, 0x42, 0x2e, 0x5a, 0x2c, 0x67,
	0x69, 0x74, 0x68, 0x75, 0x62, 0x2e, 0x63, 0x6f, 0x6d, 0x2f, 0x67, 0x6f, 0x6e, 0x7a, 0x6f, 0x6a,
	0x69, 0x76, 0x65, 0x2f, 0x68, 0x65, 0x61, 0x74, 0x70, 0x75, 0x6d, 0x70, 0x2f, 0x70, 0x72, 0x6f,
	0x74, 0x6f, 0x2f, 0x63, 0x68, 0x69, 0x6c, 0x74, 0x72, 0x69, 0x78, 0x62, 0x06, 0x70, 0x72, 0x6f,
	0x74, 0x6f, 0x33,
}

var (
	file_proto_chiltrix_chiltrix_proto_rawDescOnce sync.Once
	file_proto_chiltrix_chiltrix_proto_rawDescData = file_proto_chiltrix_chiltrix_proto_rawDesc
)

func file_proto_chiltrix_chiltrix_proto_rawDescGZIP() []byte {
	file_proto_chiltrix_chiltrix_proto_rawDescOnce.Do(func() {
		file_proto_chiltrix_chiltrix_proto_rawDescData = protoimpl.X.CompressGZIP(file_proto_chiltrix_chiltrix_proto_rawDescData)
	})
	return file_proto_chiltrix_chiltrix_proto_rawDescData
}

var file_proto_chiltrix_chiltrix_proto_msgTypes = make([]protoimpl.MessageInfo, 12)
var file_proto_chiltrix_chiltrix_proto_goTypes = []any{
	(*State)(nil),                 // 0: chiltrix.State
	(*RegisterSnapshot)(nil),      // 1: chiltrix.RegisterSnapshot
	(*GetStateRequest)(nil),       // 2: chiltrix.GetStateRequest
	(*QueryStreamRequest)(nil),    // 3: chiltrix.QueryStreamRequest
	(*QueryStreamResponse)(nil),   // 4: chiltrix.QueryStreamResponse
	(*StateSequence)(nil),         // 5: chiltrix.StateSequence
	(*StateDiff)(nil),             // 6: chiltrix.StateDiff
	(*SetParameterRequest)(nil),   // 7: chiltrix.SetParameterRequest
	(*RegisterValue)(nil),         // 8: chiltrix.RegisterValue
	(*SetParameterResponse)(nil),  // 9: chiltrix.SetParameterResponse
	(*Temperature)(nil),           // 10: chiltrix.Temperature
	nil,                           // 11: chiltrix.RegisterSnapshot.HoldingRegisterValuesEntry
	(*timestamppb.Timestamp)(nil), // 12: google.protobuf.Timestamp
}
var file_proto_chiltrix_chiltrix_proto_depIdxs = []int32{
	12, // 0: chiltrix.State.collection_time:type_name -> google.protobuf.Timestamp
	1,  // 1: chiltrix.State.register_values:type_name -> chiltrix.RegisterSnapshot
	11, // 2: chiltrix.RegisterSnapshot.holding_register_values:type_name -> chiltrix.RegisterSnapshot.HoldingRegisterValuesEntry
	12, // 3: chiltrix.QueryStreamRequest.start_time:type_name -> google.protobuf.Timestamp
	12, // 4: chiltrix.QueryStreamRequest.end_time:type_name -> google.protobuf.Timestamp
	0,  // 5: chiltrix.QueryStreamResponse.state:type_name -> chiltrix.State
	6,  // 6: chiltrix.StateSequence.diffs:type_name -> chiltrix.StateDiff
	12, // 7: chiltrix.StateDiff.collection_time:type_name -> google.protobuf.Timestamp
	1,  // 8: chiltrix.StateDiff.updated_values:type_name -> chiltrix.RegisterSnapshot
	10, // 9: chiltrix.SetParameterRequest.target_heating_mode_temperature:type_name -> chiltrix.Temperature
	8,  // 10: chiltrix.SetParameterRequest.register_value:type_name -> chiltrix.RegisterValue
	3,  // 11: chiltrix.Historian.QueryStream:input_type -> chiltrix.QueryStreamRequest
	7,  // 12: chiltrix.ReadWriteService.SetParameter:input_type -> chiltrix.SetParameterRequest
	2,  // 13: chiltrix.ReadWriteService.GetState:input_type -> chiltrix.GetStateRequest
	4,  // 14: chiltrix.Historian.QueryStream:output_type -> chiltrix.QueryStreamResponse
	9,  // 15: chiltrix.ReadWriteService.SetParameter:output_type -> chiltrix.SetParameterResponse
	0,  // 16: chiltrix.ReadWriteService.GetState:output_type -> chiltrix.State
	14, // [14:17] is the sub-list for method output_type
	11, // [11:14] is the sub-list for method input_type
	11, // [11:11] is the sub-list for extension type_name
	11, // [11:11] is the sub-list for extension extendee
	0,  // [0:11] is the sub-list for field type_name
}

func init() { file_proto_chiltrix_chiltrix_proto_init() }
func file_proto_chiltrix_chiltrix_proto_init() {
	if File_proto_chiltrix_chiltrix_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_proto_chiltrix_chiltrix_proto_msgTypes[0].Exporter = func(v any, i int) any {
			switch v := v.(*State); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_proto_chiltrix_chiltrix_proto_msgTypes[1].Exporter = func(v any, i int) any {
			switch v := v.(*RegisterSnapshot); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_proto_chiltrix_chiltrix_proto_msgTypes[2].Exporter = func(v any, i int) any {
			switch v := v.(*GetStateRequest); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_proto_chiltrix_chiltrix_proto_msgTypes[3].Exporter = func(v any, i int) any {
			switch v := v.(*QueryStreamRequest); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_proto_chiltrix_chiltrix_proto_msgTypes[4].Exporter = func(v any, i int) any {
			switch v := v.(*QueryStreamResponse); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_proto_chiltrix_chiltrix_proto_msgTypes[5].Exporter = func(v any, i int) any {
			switch v := v.(*StateSequence); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_proto_chiltrix_chiltrix_proto_msgTypes[6].Exporter = func(v any, i int) any {
			switch v := v.(*StateDiff); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_proto_chiltrix_chiltrix_proto_msgTypes[7].Exporter = func(v any, i int) any {
			switch v := v.(*SetParameterRequest); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_proto_chiltrix_chiltrix_proto_msgTypes[8].Exporter = func(v any, i int) any {
			switch v := v.(*RegisterValue); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_proto_chiltrix_chiltrix_proto_msgTypes[9].Exporter = func(v any, i int) any {
			switch v := v.(*SetParameterResponse); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_proto_chiltrix_chiltrix_proto_msgTypes[10].Exporter = func(v any, i int) any {
			switch v := v.(*Temperature); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
	}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_proto_chiltrix_chiltrix_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   12,
			NumExtensions: 0,
			NumServices:   2,
		},
		GoTypes:           file_proto_chiltrix_chiltrix_proto_goTypes,
		DependencyIndexes: file_proto_chiltrix_chiltrix_proto_depIdxs,
		MessageInfos:      file_proto_chiltrix_chiltrix_proto_msgTypes,
	}.Build()
	File_proto_chiltrix_chiltrix_proto = out.File
	file_proto_chiltrix_chiltrix_proto_rawDesc = nil
	file_proto_chiltrix_chiltrix_proto_goTypes = nil
	file_proto_chiltrix_chiltrix_proto_depIdxs = nil
}

// Reference imports to suppress errors if they are not otherwise used.
var _ context.Context
var _ grpc.ClientConnInterface

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
const _ = grpc.SupportPackageIsVersion6

// HistorianClient is the client API for Historian service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://godoc.org/google.golang.org/grpc#ClientConn.NewStream.
type HistorianClient interface {
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
type HistorianServer interface {
	QueryStream(*QueryStreamRequest, Historian_QueryStreamServer) error
}

// UnimplementedHistorianServer can be embedded to have forward compatible implementations.
type UnimplementedHistorianServer struct {
}

func (*UnimplementedHistorianServer) QueryStream(*QueryStreamRequest, Historian_QueryStreamServer) error {
	return status.Errorf(codes.Unimplemented, "method QueryStream not implemented")
}

func RegisterHistorianServer(s *grpc.Server, srv HistorianServer) {
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
	Metadata: "proto/chiltrix/chiltrix.proto",
}

// ReadWriteServiceClient is the client API for ReadWriteService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://godoc.org/google.golang.org/grpc#ClientConn.NewStream.
type ReadWriteServiceClient interface {
	SetParameter(ctx context.Context, in *SetParameterRequest, opts ...grpc.CallOption) (*SetParameterResponse, error)
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
type ReadWriteServiceServer interface {
	SetParameter(context.Context, *SetParameterRequest) (*SetParameterResponse, error)
	GetState(context.Context, *GetStateRequest) (*State, error)
}

// UnimplementedReadWriteServiceServer can be embedded to have forward compatible implementations.
type UnimplementedReadWriteServiceServer struct {
}

func (*UnimplementedReadWriteServiceServer) SetParameter(context.Context, *SetParameterRequest) (*SetParameterResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method SetParameter not implemented")
}
func (*UnimplementedReadWriteServiceServer) GetState(context.Context, *GetStateRequest) (*State, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetState not implemented")
}

func RegisterReadWriteServiceServer(s *grpc.Server, srv ReadWriteServiceServer) {
	s.RegisterService(&_ReadWriteService_serviceDesc, srv)
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

var _ReadWriteService_serviceDesc = grpc.ServiceDesc{
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
	Metadata: "proto/chiltrix/chiltrix.proto",
}
