// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.36.1
// 	protoc        v5.29.0
// source: proto/controller/controller.proto

package controller

import (
	context "context"
	fancoil "github.com/gonzojive/heatpump/proto/fancoil"
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

type GetDeviceStateRequest struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	Name          string                 `protobuf:"bytes,1,opt,name=name,proto3" json:"name,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *GetDeviceStateRequest) Reset() {
	*x = GetDeviceStateRequest{}
	mi := &file_proto_controller_controller_proto_msgTypes[0]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *GetDeviceStateRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GetDeviceStateRequest) ProtoMessage() {}

func (x *GetDeviceStateRequest) ProtoReflect() protoreflect.Message {
	mi := &file_proto_controller_controller_proto_msgTypes[0]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use GetDeviceStateRequest.ProtoReflect.Descriptor instead.
func (*GetDeviceStateRequest) Descriptor() ([]byte, []int) {
	return file_proto_controller_controller_proto_rawDescGZIP(), []int{0}
}

func (x *GetDeviceStateRequest) GetName() string {
	if x != nil {
		return x.Name
	}
	return ""
}

type DeviceState struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	Name          string                 `protobuf:"bytes,1,opt,name=name,proto3" json:"name,omitempty"`
	FancoilState  *fancoil.State         `protobuf:"bytes,2,opt,name=fancoil_state,json=fancoilState,proto3" json:"fancoil_state,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *DeviceState) Reset() {
	*x = DeviceState{}
	mi := &file_proto_controller_controller_proto_msgTypes[1]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *DeviceState) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*DeviceState) ProtoMessage() {}

func (x *DeviceState) ProtoReflect() protoreflect.Message {
	mi := &file_proto_controller_controller_proto_msgTypes[1]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use DeviceState.ProtoReflect.Descriptor instead.
func (*DeviceState) Descriptor() ([]byte, []int) {
	return file_proto_controller_controller_proto_rawDescGZIP(), []int{1}
}

func (x *DeviceState) GetName() string {
	if x != nil {
		return x.Name
	}
	return ""
}

func (x *DeviceState) GetFancoilState() *fancoil.State {
	if x != nil {
		return x.FancoilState
	}
	return nil
}

type SetDeviceStateRequest struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	State         *DeviceState           `protobuf:"bytes,1,opt,name=state,proto3" json:"state,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *SetDeviceStateRequest) Reset() {
	*x = SetDeviceStateRequest{}
	mi := &file_proto_controller_controller_proto_msgTypes[2]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *SetDeviceStateRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*SetDeviceStateRequest) ProtoMessage() {}

func (x *SetDeviceStateRequest) ProtoReflect() protoreflect.Message {
	mi := &file_proto_controller_controller_proto_msgTypes[2]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use SetDeviceStateRequest.ProtoReflect.Descriptor instead.
func (*SetDeviceStateRequest) Descriptor() ([]byte, []int) {
	return file_proto_controller_controller_proto_rawDescGZIP(), []int{2}
}

func (x *SetDeviceStateRequest) GetState() *DeviceState {
	if x != nil {
		return x.State
	}
	return nil
}

type SetDeviceStateResponse struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *SetDeviceStateResponse) Reset() {
	*x = SetDeviceStateResponse{}
	mi := &file_proto_controller_controller_proto_msgTypes[3]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *SetDeviceStateResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*SetDeviceStateResponse) ProtoMessage() {}

func (x *SetDeviceStateResponse) ProtoReflect() protoreflect.Message {
	mi := &file_proto_controller_controller_proto_msgTypes[3]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use SetDeviceStateResponse.ProtoReflect.Descriptor instead.
func (*SetDeviceStateResponse) Descriptor() ([]byte, []int) {
	return file_proto_controller_controller_proto_rawDescGZIP(), []int{3}
}

type Configuration struct {
	state                    protoimpl.MessageState `protogen:"open.v1"`
	CloudIotDeviceName       string                 `protobuf:"bytes,1,opt,name=cloud_iot_device_name,json=cloudIotDeviceName,proto3" json:"cloud_iot_device_name,omitempty"`
	HeatpumpModbusDevicePath string                 `protobuf:"bytes,2,opt,name=heatpump_modbus_device_path,json=heatpumpModbusDevicePath,proto3" json:"heatpump_modbus_device_path,omitempty"`
	FanCoilsModbusDevicePath string                 `protobuf:"bytes,3,opt,name=fan_coils_modbus_device_path,json=fanCoilsModbusDevicePath,proto3" json:"fan_coils_modbus_device_path,omitempty"`
	unknownFields            protoimpl.UnknownFields
	sizeCache                protoimpl.SizeCache
}

func (x *Configuration) Reset() {
	*x = Configuration{}
	mi := &file_proto_controller_controller_proto_msgTypes[4]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *Configuration) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Configuration) ProtoMessage() {}

func (x *Configuration) ProtoReflect() protoreflect.Message {
	mi := &file_proto_controller_controller_proto_msgTypes[4]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Configuration.ProtoReflect.Descriptor instead.
func (*Configuration) Descriptor() ([]byte, []int) {
	return file_proto_controller_controller_proto_rawDescGZIP(), []int{4}
}

func (x *Configuration) GetCloudIotDeviceName() string {
	if x != nil {
		return x.CloudIotDeviceName
	}
	return ""
}

func (x *Configuration) GetHeatpumpModbusDevicePath() string {
	if x != nil {
		return x.HeatpumpModbusDevicePath
	}
	return ""
}

func (x *Configuration) GetFanCoilsModbusDevicePath() string {
	if x != nil {
		return x.FanCoilsModbusDevicePath
	}
	return ""
}

type Command struct {
	state protoimpl.MessageState `protogen:"open.v1"`
	// Types that are valid to be assigned to Command:
	//
	//	*Command_SetStateRequest
	Command       isCommand_Command `protobuf_oneof:"command"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *Command) Reset() {
	*x = Command{}
	mi := &file_proto_controller_controller_proto_msgTypes[5]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *Command) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Command) ProtoMessage() {}

func (x *Command) ProtoReflect() protoreflect.Message {
	mi := &file_proto_controller_controller_proto_msgTypes[5]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Command.ProtoReflect.Descriptor instead.
func (*Command) Descriptor() ([]byte, []int) {
	return file_proto_controller_controller_proto_rawDescGZIP(), []int{5}
}

func (x *Command) GetCommand() isCommand_Command {
	if x != nil {
		return x.Command
	}
	return nil
}

func (x *Command) GetSetStateRequest() *fancoil.SetStateRequest {
	if x != nil {
		if x, ok := x.Command.(*Command_SetStateRequest); ok {
			return x.SetStateRequest
		}
	}
	return nil
}

type isCommand_Command interface {
	isCommand_Command()
}

type Command_SetStateRequest struct {
	SetStateRequest *fancoil.SetStateRequest `protobuf:"bytes,1,opt,name=set_state_request,json=setStateRequest,proto3,oneof"`
}

func (*Command_SetStateRequest) isCommand_Command() {}

type DeviceAccessToken struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	UserId        string                 `protobuf:"bytes,1,opt,name=user_id,json=userId,proto3" json:"user_id,omitempty"`
	Expiration    *timestamppb.Timestamp `protobuf:"bytes,2,opt,name=expiration,proto3" json:"expiration,omitempty"`
	Signature     string                 `protobuf:"bytes,3,opt,name=signature,proto3" json:"signature,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *DeviceAccessToken) Reset() {
	*x = DeviceAccessToken{}
	mi := &file_proto_controller_controller_proto_msgTypes[6]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *DeviceAccessToken) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*DeviceAccessToken) ProtoMessage() {}

func (x *DeviceAccessToken) ProtoReflect() protoreflect.Message {
	mi := &file_proto_controller_controller_proto_msgTypes[6]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use DeviceAccessToken.ProtoReflect.Descriptor instead.
func (*DeviceAccessToken) Descriptor() ([]byte, []int) {
	return file_proto_controller_controller_proto_rawDescGZIP(), []int{6}
}

func (x *DeviceAccessToken) GetUserId() string {
	if x != nil {
		return x.UserId
	}
	return ""
}

func (x *DeviceAccessToken) GetExpiration() *timestamppb.Timestamp {
	if x != nil {
		return x.Expiration
	}
	return nil
}

func (x *DeviceAccessToken) GetSignature() string {
	if x != nil {
		return x.Signature
	}
	return ""
}

type ExtendTokenRequest struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	Token         string                 `protobuf:"bytes,1,opt,name=token,proto3" json:"token,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *ExtendTokenRequest) Reset() {
	*x = ExtendTokenRequest{}
	mi := &file_proto_controller_controller_proto_msgTypes[7]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *ExtendTokenRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ExtendTokenRequest) ProtoMessage() {}

func (x *ExtendTokenRequest) ProtoReflect() protoreflect.Message {
	mi := &file_proto_controller_controller_proto_msgTypes[7]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ExtendTokenRequest.ProtoReflect.Descriptor instead.
func (*ExtendTokenRequest) Descriptor() ([]byte, []int) {
	return file_proto_controller_controller_proto_rawDescGZIP(), []int{7}
}

func (x *ExtendTokenRequest) GetToken() string {
	if x != nil {
		return x.Token
	}
	return ""
}

type ExtendTokenResponse struct {
	state          protoimpl.MessageState `protogen:"open.v1"`
	RefreshedToken string                 `protobuf:"bytes,1,opt,name=refreshed_token,json=refreshedToken,proto3" json:"refreshed_token,omitempty"`
	Expiration     *timestamppb.Timestamp `protobuf:"bytes,2,opt,name=expiration,proto3" json:"expiration,omitempty"`
	unknownFields  protoimpl.UnknownFields
	sizeCache      protoimpl.SizeCache
}

func (x *ExtendTokenResponse) Reset() {
	*x = ExtendTokenResponse{}
	mi := &file_proto_controller_controller_proto_msgTypes[8]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *ExtendTokenResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ExtendTokenResponse) ProtoMessage() {}

func (x *ExtendTokenResponse) ProtoReflect() protoreflect.Message {
	mi := &file_proto_controller_controller_proto_msgTypes[8]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ExtendTokenResponse.ProtoReflect.Descriptor instead.
func (*ExtendTokenResponse) Descriptor() ([]byte, []int) {
	return file_proto_controller_controller_proto_rawDescGZIP(), []int{8}
}

func (x *ExtendTokenResponse) GetRefreshedToken() string {
	if x != nil {
		return x.RefreshedToken
	}
	return ""
}

func (x *ExtendTokenResponse) GetExpiration() *timestamppb.Timestamp {
	if x != nil {
		return x.Expiration
	}
	return nil
}

var File_proto_controller_controller_proto protoreflect.FileDescriptor

var file_proto_controller_controller_proto_rawDesc = []byte{
	0x0a, 0x21, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2f, 0x63, 0x6f, 0x6e, 0x74, 0x72, 0x6f, 0x6c, 0x6c,
	0x65, 0x72, 0x2f, 0x63, 0x6f, 0x6e, 0x74, 0x72, 0x6f, 0x6c, 0x6c, 0x65, 0x72, 0x2e, 0x70, 0x72,
	0x6f, 0x74, 0x6f, 0x12, 0x0a, 0x63, 0x6f, 0x6e, 0x74, 0x72, 0x6f, 0x6c, 0x6c, 0x65, 0x72, 0x1a,
	0x1b, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2f, 0x66, 0x61, 0x6e, 0x63, 0x6f, 0x69, 0x6c, 0x2f, 0x66,
	0x61, 0x6e, 0x63, 0x6f, 0x69, 0x6c, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x1a, 0x1f, 0x67, 0x6f,
	0x6f, 0x67, 0x6c, 0x65, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2f, 0x74, 0x69,
	0x6d, 0x65, 0x73, 0x74, 0x61, 0x6d, 0x70, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x22, 0x2b, 0x0a,
	0x15, 0x47, 0x65, 0x74, 0x44, 0x65, 0x76, 0x69, 0x63, 0x65, 0x53, 0x74, 0x61, 0x74, 0x65, 0x52,
	0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x12, 0x0a, 0x04, 0x6e, 0x61, 0x6d, 0x65, 0x18, 0x01,
	0x20, 0x01, 0x28, 0x09, 0x52, 0x04, 0x6e, 0x61, 0x6d, 0x65, 0x22, 0x56, 0x0a, 0x0b, 0x44, 0x65,
	0x76, 0x69, 0x63, 0x65, 0x53, 0x74, 0x61, 0x74, 0x65, 0x12, 0x12, 0x0a, 0x04, 0x6e, 0x61, 0x6d,
	0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x04, 0x6e, 0x61, 0x6d, 0x65, 0x12, 0x33, 0x0a,
	0x0d, 0x66, 0x61, 0x6e, 0x63, 0x6f, 0x69, 0x6c, 0x5f, 0x73, 0x74, 0x61, 0x74, 0x65, 0x18, 0x02,
	0x20, 0x01, 0x28, 0x0b, 0x32, 0x0e, 0x2e, 0x66, 0x61, 0x6e, 0x63, 0x6f, 0x69, 0x6c, 0x2e, 0x53,
	0x74, 0x61, 0x74, 0x65, 0x52, 0x0c, 0x66, 0x61, 0x6e, 0x63, 0x6f, 0x69, 0x6c, 0x53, 0x74, 0x61,
	0x74, 0x65, 0x22, 0x46, 0x0a, 0x15, 0x53, 0x65, 0x74, 0x44, 0x65, 0x76, 0x69, 0x63, 0x65, 0x53,
	0x74, 0x61, 0x74, 0x65, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x2d, 0x0a, 0x05, 0x73,
	0x74, 0x61, 0x74, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x17, 0x2e, 0x63, 0x6f, 0x6e,
	0x74, 0x72, 0x6f, 0x6c, 0x6c, 0x65, 0x72, 0x2e, 0x44, 0x65, 0x76, 0x69, 0x63, 0x65, 0x53, 0x74,
	0x61, 0x74, 0x65, 0x52, 0x05, 0x73, 0x74, 0x61, 0x74, 0x65, 0x22, 0x18, 0x0a, 0x16, 0x53, 0x65,
	0x74, 0x44, 0x65, 0x76, 0x69, 0x63, 0x65, 0x53, 0x74, 0x61, 0x74, 0x65, 0x52, 0x65, 0x73, 0x70,
	0x6f, 0x6e, 0x73, 0x65, 0x22, 0xc1, 0x01, 0x0a, 0x0d, 0x43, 0x6f, 0x6e, 0x66, 0x69, 0x67, 0x75,
	0x72, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x12, 0x31, 0x0a, 0x15, 0x63, 0x6c, 0x6f, 0x75, 0x64, 0x5f,
	0x69, 0x6f, 0x74, 0x5f, 0x64, 0x65, 0x76, 0x69, 0x63, 0x65, 0x5f, 0x6e, 0x61, 0x6d, 0x65, 0x18,
	0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x12, 0x63, 0x6c, 0x6f, 0x75, 0x64, 0x49, 0x6f, 0x74, 0x44,
	0x65, 0x76, 0x69, 0x63, 0x65, 0x4e, 0x61, 0x6d, 0x65, 0x12, 0x3d, 0x0a, 0x1b, 0x68, 0x65, 0x61,
	0x74, 0x70, 0x75, 0x6d, 0x70, 0x5f, 0x6d, 0x6f, 0x64, 0x62, 0x75, 0x73, 0x5f, 0x64, 0x65, 0x76,
	0x69, 0x63, 0x65, 0x5f, 0x70, 0x61, 0x74, 0x68, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x18,
	0x68, 0x65, 0x61, 0x74, 0x70, 0x75, 0x6d, 0x70, 0x4d, 0x6f, 0x64, 0x62, 0x75, 0x73, 0x44, 0x65,
	0x76, 0x69, 0x63, 0x65, 0x50, 0x61, 0x74, 0x68, 0x12, 0x3e, 0x0a, 0x1c, 0x66, 0x61, 0x6e, 0x5f,
	0x63, 0x6f, 0x69, 0x6c, 0x73, 0x5f, 0x6d, 0x6f, 0x64, 0x62, 0x75, 0x73, 0x5f, 0x64, 0x65, 0x76,
	0x69, 0x63, 0x65, 0x5f, 0x70, 0x61, 0x74, 0x68, 0x18, 0x03, 0x20, 0x01, 0x28, 0x09, 0x52, 0x18,
	0x66, 0x61, 0x6e, 0x43, 0x6f, 0x69, 0x6c, 0x73, 0x4d, 0x6f, 0x64, 0x62, 0x75, 0x73, 0x44, 0x65,
	0x76, 0x69, 0x63, 0x65, 0x50, 0x61, 0x74, 0x68, 0x22, 0x5c, 0x0a, 0x07, 0x43, 0x6f, 0x6d, 0x6d,
	0x61, 0x6e, 0x64, 0x12, 0x46, 0x0a, 0x11, 0x73, 0x65, 0x74, 0x5f, 0x73, 0x74, 0x61, 0x74, 0x65,
	0x5f, 0x72, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x18,
	0x2e, 0x66, 0x61, 0x6e, 0x63, 0x6f, 0x69, 0x6c, 0x2e, 0x53, 0x65, 0x74, 0x53, 0x74, 0x61, 0x74,
	0x65, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x48, 0x00, 0x52, 0x0f, 0x73, 0x65, 0x74, 0x53,
	0x74, 0x61, 0x74, 0x65, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x42, 0x09, 0x0a, 0x07, 0x63,
	0x6f, 0x6d, 0x6d, 0x61, 0x6e, 0x64, 0x22, 0x86, 0x01, 0x0a, 0x11, 0x44, 0x65, 0x76, 0x69, 0x63,
	0x65, 0x41, 0x63, 0x63, 0x65, 0x73, 0x73, 0x54, 0x6f, 0x6b, 0x65, 0x6e, 0x12, 0x17, 0x0a, 0x07,
	0x75, 0x73, 0x65, 0x72, 0x5f, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x06, 0x75,
	0x73, 0x65, 0x72, 0x49, 0x64, 0x12, 0x3a, 0x0a, 0x0a, 0x65, 0x78, 0x70, 0x69, 0x72, 0x61, 0x74,
	0x69, 0x6f, 0x6e, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x1a, 0x2e, 0x67, 0x6f, 0x6f, 0x67,
	0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x54, 0x69, 0x6d, 0x65,
	0x73, 0x74, 0x61, 0x6d, 0x70, 0x52, 0x0a, 0x65, 0x78, 0x70, 0x69, 0x72, 0x61, 0x74, 0x69, 0x6f,
	0x6e, 0x12, 0x1c, 0x0a, 0x09, 0x73, 0x69, 0x67, 0x6e, 0x61, 0x74, 0x75, 0x72, 0x65, 0x18, 0x03,
	0x20, 0x01, 0x28, 0x09, 0x52, 0x09, 0x73, 0x69, 0x67, 0x6e, 0x61, 0x74, 0x75, 0x72, 0x65, 0x22,
	0x2a, 0x0a, 0x12, 0x45, 0x78, 0x74, 0x65, 0x6e, 0x64, 0x54, 0x6f, 0x6b, 0x65, 0x6e, 0x52, 0x65,
	0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x14, 0x0a, 0x05, 0x74, 0x6f, 0x6b, 0x65, 0x6e, 0x18, 0x01,
	0x20, 0x01, 0x28, 0x09, 0x52, 0x05, 0x74, 0x6f, 0x6b, 0x65, 0x6e, 0x22, 0x7a, 0x0a, 0x13, 0x45,
	0x78, 0x74, 0x65, 0x6e, 0x64, 0x54, 0x6f, 0x6b, 0x65, 0x6e, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e,
	0x73, 0x65, 0x12, 0x27, 0x0a, 0x0f, 0x72, 0x65, 0x66, 0x72, 0x65, 0x73, 0x68, 0x65, 0x64, 0x5f,
	0x74, 0x6f, 0x6b, 0x65, 0x6e, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0e, 0x72, 0x65, 0x66,
	0x72, 0x65, 0x73, 0x68, 0x65, 0x64, 0x54, 0x6f, 0x6b, 0x65, 0x6e, 0x12, 0x3a, 0x0a, 0x0a, 0x65,
	0x78, 0x70, 0x69, 0x72, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0b, 0x32,
	0x1a, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75,
	0x66, 0x2e, 0x54, 0x69, 0x6d, 0x65, 0x73, 0x74, 0x61, 0x6d, 0x70, 0x52, 0x0a, 0x65, 0x78, 0x70,
	0x69, 0x72, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x32, 0xb9, 0x01, 0x0a, 0x0c, 0x53, 0x74, 0x61, 0x74,
	0x65, 0x53, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x12, 0x4e, 0x0a, 0x0e, 0x47, 0x65, 0x74, 0x44,
	0x65, 0x76, 0x69, 0x63, 0x65, 0x53, 0x74, 0x61, 0x74, 0x65, 0x12, 0x21, 0x2e, 0x63, 0x6f, 0x6e,
	0x74, 0x72, 0x6f, 0x6c, 0x6c, 0x65, 0x72, 0x2e, 0x47, 0x65, 0x74, 0x44, 0x65, 0x76, 0x69, 0x63,
	0x65, 0x53, 0x74, 0x61, 0x74, 0x65, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x17, 0x2e,
	0x63, 0x6f, 0x6e, 0x74, 0x72, 0x6f, 0x6c, 0x6c, 0x65, 0x72, 0x2e, 0x44, 0x65, 0x76, 0x69, 0x63,
	0x65, 0x53, 0x74, 0x61, 0x74, 0x65, 0x22, 0x00, 0x12, 0x59, 0x0a, 0x0e, 0x53, 0x65, 0x74, 0x44,
	0x65, 0x76, 0x69, 0x63, 0x65, 0x53, 0x74, 0x61, 0x74, 0x65, 0x12, 0x21, 0x2e, 0x63, 0x6f, 0x6e,
	0x74, 0x72, 0x6f, 0x6c, 0x6c, 0x65, 0x72, 0x2e, 0x53, 0x65, 0x74, 0x44, 0x65, 0x76, 0x69, 0x63,
	0x65, 0x53, 0x74, 0x61, 0x74, 0x65, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x22, 0x2e,
	0x63, 0x6f, 0x6e, 0x74, 0x72, 0x6f, 0x6c, 0x6c, 0x65, 0x72, 0x2e, 0x53, 0x65, 0x74, 0x44, 0x65,
	0x76, 0x69, 0x63, 0x65, 0x53, 0x74, 0x61, 0x74, 0x65, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73,
	0x65, 0x22, 0x00, 0x32, 0x5f, 0x0a, 0x0b, 0x41, 0x75, 0x74, 0x68, 0x53, 0x65, 0x72, 0x76, 0x69,
	0x63, 0x65, 0x12, 0x50, 0x0a, 0x0b, 0x45, 0x78, 0x74, 0x65, 0x6e, 0x64, 0x54, 0x6f, 0x6b, 0x65,
	0x6e, 0x12, 0x1e, 0x2e, 0x63, 0x6f, 0x6e, 0x74, 0x72, 0x6f, 0x6c, 0x6c, 0x65, 0x72, 0x2e, 0x45,
	0x78, 0x74, 0x65, 0x6e, 0x64, 0x54, 0x6f, 0x6b, 0x65, 0x6e, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73,
	0x74, 0x1a, 0x1f, 0x2e, 0x63, 0x6f, 0x6e, 0x74, 0x72, 0x6f, 0x6c, 0x6c, 0x65, 0x72, 0x2e, 0x45,
	0x78, 0x74, 0x65, 0x6e, 0x64, 0x54, 0x6f, 0x6b, 0x65, 0x6e, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e,
	0x73, 0x65, 0x22, 0x00, 0x42, 0x30, 0x5a, 0x2e, 0x67, 0x69, 0x74, 0x68, 0x75, 0x62, 0x2e, 0x63,
	0x6f, 0x6d, 0x2f, 0x67, 0x6f, 0x6e, 0x7a, 0x6f, 0x6a, 0x69, 0x76, 0x65, 0x2f, 0x68, 0x65, 0x61,
	0x74, 0x70, 0x75, 0x6d, 0x70, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2f, 0x63, 0x6f, 0x6e, 0x74,
	0x72, 0x6f, 0x6c, 0x6c, 0x65, 0x72, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_proto_controller_controller_proto_rawDescOnce sync.Once
	file_proto_controller_controller_proto_rawDescData = file_proto_controller_controller_proto_rawDesc
)

func file_proto_controller_controller_proto_rawDescGZIP() []byte {
	file_proto_controller_controller_proto_rawDescOnce.Do(func() {
		file_proto_controller_controller_proto_rawDescData = protoimpl.X.CompressGZIP(file_proto_controller_controller_proto_rawDescData)
	})
	return file_proto_controller_controller_proto_rawDescData
}

var file_proto_controller_controller_proto_msgTypes = make([]protoimpl.MessageInfo, 9)
var file_proto_controller_controller_proto_goTypes = []any{
	(*GetDeviceStateRequest)(nil),   // 0: controller.GetDeviceStateRequest
	(*DeviceState)(nil),             // 1: controller.DeviceState
	(*SetDeviceStateRequest)(nil),   // 2: controller.SetDeviceStateRequest
	(*SetDeviceStateResponse)(nil),  // 3: controller.SetDeviceStateResponse
	(*Configuration)(nil),           // 4: controller.Configuration
	(*Command)(nil),                 // 5: controller.Command
	(*DeviceAccessToken)(nil),       // 6: controller.DeviceAccessToken
	(*ExtendTokenRequest)(nil),      // 7: controller.ExtendTokenRequest
	(*ExtendTokenResponse)(nil),     // 8: controller.ExtendTokenResponse
	(*fancoil.State)(nil),           // 9: fancoil.State
	(*fancoil.SetStateRequest)(nil), // 10: fancoil.SetStateRequest
	(*timestamppb.Timestamp)(nil),   // 11: google.protobuf.Timestamp
}
var file_proto_controller_controller_proto_depIdxs = []int32{
	9,  // 0: controller.DeviceState.fancoil_state:type_name -> fancoil.State
	1,  // 1: controller.SetDeviceStateRequest.state:type_name -> controller.DeviceState
	10, // 2: controller.Command.set_state_request:type_name -> fancoil.SetStateRequest
	11, // 3: controller.DeviceAccessToken.expiration:type_name -> google.protobuf.Timestamp
	11, // 4: controller.ExtendTokenResponse.expiration:type_name -> google.protobuf.Timestamp
	0,  // 5: controller.StateService.GetDeviceState:input_type -> controller.GetDeviceStateRequest
	2,  // 6: controller.StateService.SetDeviceState:input_type -> controller.SetDeviceStateRequest
	7,  // 7: controller.AuthService.ExtendToken:input_type -> controller.ExtendTokenRequest
	1,  // 8: controller.StateService.GetDeviceState:output_type -> controller.DeviceState
	3,  // 9: controller.StateService.SetDeviceState:output_type -> controller.SetDeviceStateResponse
	8,  // 10: controller.AuthService.ExtendToken:output_type -> controller.ExtendTokenResponse
	8,  // [8:11] is the sub-list for method output_type
	5,  // [5:8] is the sub-list for method input_type
	5,  // [5:5] is the sub-list for extension type_name
	5,  // [5:5] is the sub-list for extension extendee
	0,  // [0:5] is the sub-list for field type_name
}

func init() { file_proto_controller_controller_proto_init() }
func file_proto_controller_controller_proto_init() {
	if File_proto_controller_controller_proto != nil {
		return
	}
	file_proto_controller_controller_proto_msgTypes[5].OneofWrappers = []any{
		(*Command_SetStateRequest)(nil),
	}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_proto_controller_controller_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   9,
			NumExtensions: 0,
			NumServices:   2,
		},
		GoTypes:           file_proto_controller_controller_proto_goTypes,
		DependencyIndexes: file_proto_controller_controller_proto_depIdxs,
		MessageInfos:      file_proto_controller_controller_proto_msgTypes,
	}.Build()
	File_proto_controller_controller_proto = out.File
	file_proto_controller_controller_proto_rawDesc = nil
	file_proto_controller_controller_proto_goTypes = nil
	file_proto_controller_controller_proto_depIdxs = nil
}

// Reference imports to suppress errors if they are not otherwise used.
var _ context.Context
var _ grpc.ClientConnInterface

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
const _ = grpc.SupportPackageIsVersion6

// StateServiceClient is the client API for StateService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://godoc.org/google.golang.org/grpc#ClientConn.NewStream.
type StateServiceClient interface {
	GetDeviceState(ctx context.Context, in *GetDeviceStateRequest, opts ...grpc.CallOption) (*DeviceState, error)
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
type StateServiceServer interface {
	GetDeviceState(context.Context, *GetDeviceStateRequest) (*DeviceState, error)
	SetDeviceState(context.Context, *SetDeviceStateRequest) (*SetDeviceStateResponse, error)
}

// UnimplementedStateServiceServer can be embedded to have forward compatible implementations.
type UnimplementedStateServiceServer struct {
}

func (*UnimplementedStateServiceServer) GetDeviceState(context.Context, *GetDeviceStateRequest) (*DeviceState, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetDeviceState not implemented")
}
func (*UnimplementedStateServiceServer) SetDeviceState(context.Context, *SetDeviceStateRequest) (*SetDeviceStateResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method SetDeviceState not implemented")
}

func RegisterStateServiceServer(s *grpc.Server, srv StateServiceServer) {
	s.RegisterService(&_StateService_serviceDesc, srv)
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

var _StateService_serviceDesc = grpc.ServiceDesc{
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
	Metadata: "proto/controller/controller.proto",
}

// AuthServiceClient is the client API for AuthService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://godoc.org/google.golang.org/grpc#ClientConn.NewStream.
type AuthServiceClient interface {
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
type AuthServiceServer interface {
	ExtendToken(context.Context, *ExtendTokenRequest) (*ExtendTokenResponse, error)
}

// UnimplementedAuthServiceServer can be embedded to have forward compatible implementations.
type UnimplementedAuthServiceServer struct {
}

func (*UnimplementedAuthServiceServer) ExtendToken(context.Context, *ExtendTokenRequest) (*ExtendTokenResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ExtendToken not implemented")
}

func RegisterAuthServiceServer(s *grpc.Server, srv AuthServiceServer) {
	s.RegisterService(&_AuthService_serviceDesc, srv)
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

var _AuthService_serviceDesc = grpc.ServiceDesc{
	ServiceName: "controller.AuthService",
	HandlerType: (*AuthServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "ExtendToken",
			Handler:    _AuthService_ExtendToken_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "proto/controller/controller.proto",
}
