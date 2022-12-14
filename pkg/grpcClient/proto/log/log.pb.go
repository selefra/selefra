// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.26.0
// 	protoc        v3.17.3
// source: pkg/grpcClient/proto/log/log.proto

package log

import (
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
	anypb "google.golang.org/protobuf/types/known/anypb"
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

type BaseConnectionInfo struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Token  string `protobuf:"bytes,1,opt,name=Token,proto3" json:"Token,omitempty"`
	TaskId string `protobuf:"bytes,2,opt,name=TaskId,proto3" json:"TaskId,omitempty"`
}

func (x *BaseConnectionInfo) Reset() {
	*x = BaseConnectionInfo{}
	if protoimpl.UnsafeEnabled {
		mi := &file_pkg_grpcClient_proto_log_log_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *BaseConnectionInfo) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*BaseConnectionInfo) ProtoMessage() {}

func (x *BaseConnectionInfo) ProtoReflect() protoreflect.Message {
	mi := &file_pkg_grpcClient_proto_log_log_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use BaseConnectionInfo.ProtoReflect.Descriptor instead.
func (*BaseConnectionInfo) Descriptor() ([]byte, []int) {
	return file_pkg_grpcClient_proto_log_log_proto_rawDescGZIP(), []int{0}
}

func (x *BaseConnectionInfo) GetToken() string {
	if x != nil {
		return x.Token
	}
	return ""
}

func (x *BaseConnectionInfo) GetTaskId() string {
	if x != nil {
		return x.TaskId
	}
	return ""
}

type LogJOSN struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Cmd   string                 `protobuf:"bytes,1,opt,name=Cmd,proto3" json:"Cmd,omitempty"`
	Stag  string                 `protobuf:"bytes,2,opt,name=Stag,proto3" json:"Stag,omitempty"`
	Msg   string                 `protobuf:"bytes,3,opt,name=Msg,proto3" json:"Msg,omitempty"`
	Time  *timestamppb.Timestamp `protobuf:"bytes,4,opt,name=Time,proto3" json:"Time,omitempty"`
	Level string                 `protobuf:"bytes,5,opt,name=Level,proto3" json:"Level,omitempty"`
}

func (x *LogJOSN) Reset() {
	*x = LogJOSN{}
	if protoimpl.UnsafeEnabled {
		mi := &file_pkg_grpcClient_proto_log_log_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *LogJOSN) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*LogJOSN) ProtoMessage() {}

func (x *LogJOSN) ProtoReflect() protoreflect.Message {
	mi := &file_pkg_grpcClient_proto_log_log_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use LogJOSN.ProtoReflect.Descriptor instead.
func (*LogJOSN) Descriptor() ([]byte, []int) {
	return file_pkg_grpcClient_proto_log_log_proto_rawDescGZIP(), []int{1}
}

func (x *LogJOSN) GetCmd() string {
	if x != nil {
		return x.Cmd
	}
	return ""
}

func (x *LogJOSN) GetStag() string {
	if x != nil {
		return x.Stag
	}
	return ""
}

func (x *LogJOSN) GetMsg() string {
	if x != nil {
		return x.Msg
	}
	return ""
}

func (x *LogJOSN) GetTime() *timestamppb.Timestamp {
	if x != nil {
		return x.Time
	}
	return nil
}

func (x *LogJOSN) GetLevel() string {
	if x != nil {
		return x.Level
	}
	return ""
}

type ConnectMsg struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	ActionName string              `protobuf:"bytes,1,opt,name=ActionName,proto3" json:"ActionName,omitempty"`
	Data       *LogJOSN            `protobuf:"bytes,2,opt,name=Data,proto3" json:"Data,omitempty"`
	Index      int32               `protobuf:"varint,3,opt,name=index,proto3" json:"index,omitempty"`
	Msg        string              `protobuf:"bytes,4,opt,name=Msg,proto3" json:"Msg,omitempty"`
	BaseInfo   *BaseConnectionInfo `protobuf:"bytes,5,opt,name=BaseInfo,proto3" json:"BaseInfo,omitempty"`
}

func (x *ConnectMsg) Reset() {
	*x = ConnectMsg{}
	if protoimpl.UnsafeEnabled {
		mi := &file_pkg_grpcClient_proto_log_log_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ConnectMsg) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ConnectMsg) ProtoMessage() {}

func (x *ConnectMsg) ProtoReflect() protoreflect.Message {
	mi := &file_pkg_grpcClient_proto_log_log_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ConnectMsg.ProtoReflect.Descriptor instead.
func (*ConnectMsg) Descriptor() ([]byte, []int) {
	return file_pkg_grpcClient_proto_log_log_proto_rawDescGZIP(), []int{2}
}

func (x *ConnectMsg) GetActionName() string {
	if x != nil {
		return x.ActionName
	}
	return ""
}

func (x *ConnectMsg) GetData() *LogJOSN {
	if x != nil {
		return x.Data
	}
	return nil
}

func (x *ConnectMsg) GetIndex() int32 {
	if x != nil {
		return x.Index
	}
	return 0
}

func (x *ConnectMsg) GetMsg() string {
	if x != nil {
		return x.Msg
	}
	return ""
}

func (x *ConnectMsg) GetBaseInfo() *BaseConnectionInfo {
	if x != nil {
		return x.BaseInfo
	}
	return nil
}

type StatusInfo struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	BaseInfo *BaseConnectionInfo    `protobuf:"bytes,1,opt,name=BaseInfo,proto3" json:"BaseInfo,omitempty"`
	Stag     string                 `protobuf:"bytes,2,opt,name=Stag,proto3" json:"Stag,omitempty"`
	Status   string                 `protobuf:"bytes,3,opt,name=Status,proto3" json:"Status,omitempty"`
	Time     *timestamppb.Timestamp `protobuf:"bytes,4,opt,name=Time,proto3" json:"Time,omitempty"`
}

func (x *StatusInfo) Reset() {
	*x = StatusInfo{}
	if protoimpl.UnsafeEnabled {
		mi := &file_pkg_grpcClient_proto_log_log_proto_msgTypes[3]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *StatusInfo) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*StatusInfo) ProtoMessage() {}

func (x *StatusInfo) ProtoReflect() protoreflect.Message {
	mi := &file_pkg_grpcClient_proto_log_log_proto_msgTypes[3]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use StatusInfo.ProtoReflect.Descriptor instead.
func (*StatusInfo) Descriptor() ([]byte, []int) {
	return file_pkg_grpcClient_proto_log_log_proto_rawDescGZIP(), []int{3}
}

func (x *StatusInfo) GetBaseInfo() *BaseConnectionInfo {
	if x != nil {
		return x.BaseInfo
	}
	return nil
}

func (x *StatusInfo) GetStag() string {
	if x != nil {
		return x.Stag
	}
	return ""
}

func (x *StatusInfo) GetStatus() string {
	if x != nil {
		return x.Status
	}
	return ""
}

func (x *StatusInfo) GetTime() *timestamppb.Timestamp {
	if x != nil {
		return x.Time
	}
	return nil
}

type Res struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Code int32      `protobuf:"varint,1,opt,name=Code,proto3" json:"Code,omitempty"`
	Msg  string     `protobuf:"bytes,2,opt,name=Msg,proto3" json:"Msg,omitempty"`
	Data *anypb.Any `protobuf:"bytes,3,opt,name=Data,proto3" json:"Data,omitempty"`
}

func (x *Res) Reset() {
	*x = Res{}
	if protoimpl.UnsafeEnabled {
		mi := &file_pkg_grpcClient_proto_log_log_proto_msgTypes[4]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Res) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Res) ProtoMessage() {}

func (x *Res) ProtoReflect() protoreflect.Message {
	mi := &file_pkg_grpcClient_proto_log_log_proto_msgTypes[4]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Res.ProtoReflect.Descriptor instead.
func (*Res) Descriptor() ([]byte, []int) {
	return file_pkg_grpcClient_proto_log_log_proto_rawDescGZIP(), []int{4}
}

func (x *Res) GetCode() int32 {
	if x != nil {
		return x.Code
	}
	return 0
}

func (x *Res) GetMsg() string {
	if x != nil {
		return x.Msg
	}
	return ""
}

func (x *Res) GetData() *anypb.Any {
	if x != nil {
		return x.Data
	}
	return nil
}

var File_pkg_grpcClient_proto_log_log_proto protoreflect.FileDescriptor

var file_pkg_grpcClient_proto_log_log_proto_rawDesc = []byte{
	0x0a, 0x22, 0x70, 0x6b, 0x67, 0x2f, 0x67, 0x72, 0x70, 0x63, 0x43, 0x6c, 0x69, 0x65, 0x6e, 0x74,
	0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2f, 0x6c, 0x6f, 0x67, 0x2f, 0x6c, 0x6f, 0x67, 0x2e, 0x70,
	0x72, 0x6f, 0x74, 0x6f, 0x12, 0x03, 0x6c, 0x6f, 0x67, 0x1a, 0x1f, 0x67, 0x6f, 0x6f, 0x67, 0x6c,
	0x65, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2f, 0x74, 0x69, 0x6d, 0x65, 0x73,
	0x74, 0x61, 0x6d, 0x70, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x1a, 0x19, 0x67, 0x6f, 0x6f, 0x67,
	0x6c, 0x65, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2f, 0x61, 0x6e, 0x79, 0x2e,
	0x70, 0x72, 0x6f, 0x74, 0x6f, 0x22, 0x42, 0x0a, 0x12, 0x42, 0x61, 0x73, 0x65, 0x43, 0x6f, 0x6e,
	0x6e, 0x65, 0x63, 0x74, 0x69, 0x6f, 0x6e, 0x49, 0x6e, 0x66, 0x6f, 0x12, 0x14, 0x0a, 0x05, 0x54,
	0x6f, 0x6b, 0x65, 0x6e, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x05, 0x54, 0x6f, 0x6b, 0x65,
	0x6e, 0x12, 0x16, 0x0a, 0x06, 0x54, 0x61, 0x73, 0x6b, 0x49, 0x64, 0x18, 0x02, 0x20, 0x01, 0x28,
	0x09, 0x52, 0x06, 0x54, 0x61, 0x73, 0x6b, 0x49, 0x64, 0x22, 0x87, 0x01, 0x0a, 0x07, 0x4c, 0x6f,
	0x67, 0x4a, 0x4f, 0x53, 0x4e, 0x12, 0x10, 0x0a, 0x03, 0x43, 0x6d, 0x64, 0x18, 0x01, 0x20, 0x01,
	0x28, 0x09, 0x52, 0x03, 0x43, 0x6d, 0x64, 0x12, 0x12, 0x0a, 0x04, 0x53, 0x74, 0x61, 0x67, 0x18,
	0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x04, 0x53, 0x74, 0x61, 0x67, 0x12, 0x10, 0x0a, 0x03, 0x4d,
	0x73, 0x67, 0x18, 0x03, 0x20, 0x01, 0x28, 0x09, 0x52, 0x03, 0x4d, 0x73, 0x67, 0x12, 0x2e, 0x0a,
	0x04, 0x54, 0x69, 0x6d, 0x65, 0x18, 0x04, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x1a, 0x2e, 0x67, 0x6f,
	0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x54, 0x69,
	0x6d, 0x65, 0x73, 0x74, 0x61, 0x6d, 0x70, 0x52, 0x04, 0x54, 0x69, 0x6d, 0x65, 0x12, 0x14, 0x0a,
	0x05, 0x4c, 0x65, 0x76, 0x65, 0x6c, 0x18, 0x05, 0x20, 0x01, 0x28, 0x09, 0x52, 0x05, 0x4c, 0x65,
	0x76, 0x65, 0x6c, 0x22, 0xab, 0x01, 0x0a, 0x0a, 0x43, 0x6f, 0x6e, 0x6e, 0x65, 0x63, 0x74, 0x4d,
	0x73, 0x67, 0x12, 0x1e, 0x0a, 0x0a, 0x41, 0x63, 0x74, 0x69, 0x6f, 0x6e, 0x4e, 0x61, 0x6d, 0x65,
	0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0a, 0x41, 0x63, 0x74, 0x69, 0x6f, 0x6e, 0x4e, 0x61,
	0x6d, 0x65, 0x12, 0x20, 0x0a, 0x04, 0x44, 0x61, 0x74, 0x61, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0b,
	0x32, 0x0c, 0x2e, 0x6c, 0x6f, 0x67, 0x2e, 0x4c, 0x6f, 0x67, 0x4a, 0x4f, 0x53, 0x4e, 0x52, 0x04,
	0x44, 0x61, 0x74, 0x61, 0x12, 0x14, 0x0a, 0x05, 0x69, 0x6e, 0x64, 0x65, 0x78, 0x18, 0x03, 0x20,
	0x01, 0x28, 0x05, 0x52, 0x05, 0x69, 0x6e, 0x64, 0x65, 0x78, 0x12, 0x10, 0x0a, 0x03, 0x4d, 0x73,
	0x67, 0x18, 0x04, 0x20, 0x01, 0x28, 0x09, 0x52, 0x03, 0x4d, 0x73, 0x67, 0x12, 0x33, 0x0a, 0x08,
	0x42, 0x61, 0x73, 0x65, 0x49, 0x6e, 0x66, 0x6f, 0x18, 0x05, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x17,
	0x2e, 0x6c, 0x6f, 0x67, 0x2e, 0x42, 0x61, 0x73, 0x65, 0x43, 0x6f, 0x6e, 0x6e, 0x65, 0x63, 0x74,
	0x69, 0x6f, 0x6e, 0x49, 0x6e, 0x66, 0x6f, 0x52, 0x08, 0x42, 0x61, 0x73, 0x65, 0x49, 0x6e, 0x66,
	0x6f, 0x22, 0x9d, 0x01, 0x0a, 0x0a, 0x53, 0x74, 0x61, 0x74, 0x75, 0x73, 0x49, 0x6e, 0x66, 0x6f,
	0x12, 0x33, 0x0a, 0x08, 0x42, 0x61, 0x73, 0x65, 0x49, 0x6e, 0x66, 0x6f, 0x18, 0x01, 0x20, 0x01,
	0x28, 0x0b, 0x32, 0x17, 0x2e, 0x6c, 0x6f, 0x67, 0x2e, 0x42, 0x61, 0x73, 0x65, 0x43, 0x6f, 0x6e,
	0x6e, 0x65, 0x63, 0x74, 0x69, 0x6f, 0x6e, 0x49, 0x6e, 0x66, 0x6f, 0x52, 0x08, 0x42, 0x61, 0x73,
	0x65, 0x49, 0x6e, 0x66, 0x6f, 0x12, 0x12, 0x0a, 0x04, 0x53, 0x74, 0x61, 0x67, 0x18, 0x02, 0x20,
	0x01, 0x28, 0x09, 0x52, 0x04, 0x53, 0x74, 0x61, 0x67, 0x12, 0x16, 0x0a, 0x06, 0x53, 0x74, 0x61,
	0x74, 0x75, 0x73, 0x18, 0x03, 0x20, 0x01, 0x28, 0x09, 0x52, 0x06, 0x53, 0x74, 0x61, 0x74, 0x75,
	0x73, 0x12, 0x2e, 0x0a, 0x04, 0x54, 0x69, 0x6d, 0x65, 0x18, 0x04, 0x20, 0x01, 0x28, 0x0b, 0x32,
	0x1a, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75,
	0x66, 0x2e, 0x54, 0x69, 0x6d, 0x65, 0x73, 0x74, 0x61, 0x6d, 0x70, 0x52, 0x04, 0x54, 0x69, 0x6d,
	0x65, 0x22, 0x55, 0x0a, 0x03, 0x52, 0x65, 0x73, 0x12, 0x12, 0x0a, 0x04, 0x43, 0x6f, 0x64, 0x65,
	0x18, 0x01, 0x20, 0x01, 0x28, 0x05, 0x52, 0x04, 0x43, 0x6f, 0x64, 0x65, 0x12, 0x10, 0x0a, 0x03,
	0x4d, 0x73, 0x67, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x03, 0x4d, 0x73, 0x67, 0x12, 0x28,
	0x0a, 0x04, 0x44, 0x61, 0x74, 0x61, 0x18, 0x03, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x14, 0x2e, 0x67,
	0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x41,
	0x6e, 0x79, 0x52, 0x04, 0x44, 0x61, 0x74, 0x61, 0x32, 0xa5, 0x01, 0x0a, 0x03, 0x4c, 0x6f, 0x67,
	0x12, 0x30, 0x0a, 0x0f, 0x55, 0x70, 0x6c, 0x6f, 0x61, 0x64, 0x4c, 0x6f, 0x67, 0x53, 0x74, 0x72,
	0x65, 0x61, 0x6d, 0x12, 0x0f, 0x2e, 0x6c, 0x6f, 0x67, 0x2e, 0x43, 0x6f, 0x6e, 0x6e, 0x65, 0x63,
	0x74, 0x4d, 0x73, 0x67, 0x1a, 0x08, 0x2e, 0x6c, 0x6f, 0x67, 0x2e, 0x52, 0x65, 0x73, 0x22, 0x00,
	0x28, 0x01, 0x12, 0x2e, 0x0a, 0x0f, 0x55, 0x70, 0x6c, 0x6f, 0x61, 0x64, 0x4c, 0x6f, 0x67, 0x53,
	0x74, 0x61, 0x74, 0x75, 0x73, 0x12, 0x0f, 0x2e, 0x6c, 0x6f, 0x67, 0x2e, 0x53, 0x74, 0x61, 0x74,
	0x75, 0x73, 0x49, 0x6e, 0x66, 0x6f, 0x1a, 0x08, 0x2e, 0x6c, 0x6f, 0x67, 0x2e, 0x52, 0x65, 0x73,
	0x22, 0x00, 0x12, 0x3c, 0x0a, 0x0c, 0x47, 0x65, 0x74, 0x4c, 0x6f, 0x67, 0x53, 0x74, 0x72, 0x65,
	0x61, 0x6d, 0x12, 0x17, 0x2e, 0x6c, 0x6f, 0x67, 0x2e, 0x42, 0x61, 0x73, 0x65, 0x43, 0x6f, 0x6e,
	0x6e, 0x65, 0x63, 0x74, 0x69, 0x6f, 0x6e, 0x49, 0x6e, 0x66, 0x6f, 0x1a, 0x0f, 0x2e, 0x6c, 0x6f,
	0x67, 0x2e, 0x43, 0x6f, 0x6e, 0x6e, 0x65, 0x63, 0x74, 0x4d, 0x73, 0x67, 0x22, 0x00, 0x30, 0x01,
	0x42, 0x0b, 0x5a, 0x09, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2f, 0x6c, 0x6f, 0x67, 0x62, 0x06, 0x70,
	0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_pkg_grpcClient_proto_log_log_proto_rawDescOnce sync.Once
	file_pkg_grpcClient_proto_log_log_proto_rawDescData = file_pkg_grpcClient_proto_log_log_proto_rawDesc
)

func file_pkg_grpcClient_proto_log_log_proto_rawDescGZIP() []byte {
	file_pkg_grpcClient_proto_log_log_proto_rawDescOnce.Do(func() {
		file_pkg_grpcClient_proto_log_log_proto_rawDescData = protoimpl.X.CompressGZIP(file_pkg_grpcClient_proto_log_log_proto_rawDescData)
	})
	return file_pkg_grpcClient_proto_log_log_proto_rawDescData
}

var file_pkg_grpcClient_proto_log_log_proto_msgTypes = make([]protoimpl.MessageInfo, 5)
var file_pkg_grpcClient_proto_log_log_proto_goTypes = []interface{}{
	(*BaseConnectionInfo)(nil),    // 0: log.BaseConnectionInfo
	(*LogJOSN)(nil),               // 1: log.LogJOSN
	(*ConnectMsg)(nil),            // 2: log.ConnectMsg
	(*StatusInfo)(nil),            // 3: log.StatusInfo
	(*Res)(nil),                   // 4: log.Res
	(*timestamppb.Timestamp)(nil), // 5: google.protobuf.Timestamp
	(*anypb.Any)(nil),             // 6: google.protobuf.Any
}
var file_pkg_grpcClient_proto_log_log_proto_depIdxs = []int32{
	5, // 0: log.LogJOSN.Time:type_name -> google.protobuf.Timestamp
	1, // 1: log.ConnectMsg.Data:type_name -> log.LogJOSN
	0, // 2: log.ConnectMsg.BaseInfo:type_name -> log.BaseConnectionInfo
	0, // 3: log.StatusInfo.BaseInfo:type_name -> log.BaseConnectionInfo
	5, // 4: log.StatusInfo.Time:type_name -> google.protobuf.Timestamp
	6, // 5: log.Res.Data:type_name -> google.protobuf.Any
	2, // 6: log.Log.UploadLogStream:input_type -> log.ConnectMsg
	3, // 7: log.Log.UploadLogStatus:input_type -> log.StatusInfo
	0, // 8: log.Log.GetLogStream:input_type -> log.BaseConnectionInfo
	4, // 9: log.Log.UploadLogStream:output_type -> log.Res
	4, // 10: log.Log.UploadLogStatus:output_type -> log.Res
	2, // 11: log.Log.GetLogStream:output_type -> log.ConnectMsg
	9, // [9:12] is the sub-list for method output_type
	6, // [6:9] is the sub-list for method input_type
	6, // [6:6] is the sub-list for extension type_name
	6, // [6:6] is the sub-list for extension extendee
	0, // [0:6] is the sub-list for field type_name
}

func init() { file_pkg_grpcClient_proto_log_log_proto_init() }
func file_pkg_grpcClient_proto_log_log_proto_init() {
	if File_pkg_grpcClient_proto_log_log_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_pkg_grpcClient_proto_log_log_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*BaseConnectionInfo); i {
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
		file_pkg_grpcClient_proto_log_log_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*LogJOSN); i {
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
		file_pkg_grpcClient_proto_log_log_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*ConnectMsg); i {
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
		file_pkg_grpcClient_proto_log_log_proto_msgTypes[3].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*StatusInfo); i {
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
		file_pkg_grpcClient_proto_log_log_proto_msgTypes[4].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Res); i {
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
			RawDescriptor: file_pkg_grpcClient_proto_log_log_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   5,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_pkg_grpcClient_proto_log_log_proto_goTypes,
		DependencyIndexes: file_pkg_grpcClient_proto_log_log_proto_depIdxs,
		MessageInfos:      file_pkg_grpcClient_proto_log_log_proto_msgTypes,
	}.Build()
	File_pkg_grpcClient_proto_log_log_proto = out.File
	file_pkg_grpcClient_proto_log_log_proto_rawDesc = nil
	file_pkg_grpcClient_proto_log_log_proto_goTypes = nil
	file_pkg_grpcClient_proto_log_log_proto_depIdxs = nil
}
