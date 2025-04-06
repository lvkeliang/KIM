// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.36.6
// 	protoc        v3.20.3
// source: proto/innerHandshake.proto

package protoImpl

import (
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
	reflect "reflect"
	sync "sync"
	unsafe "unsafe"
)

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

type InnerHandshakeReq struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	ServiceId     string                 `protobuf:"bytes,1,opt,name=serviceId,proto3" json:"serviceId,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *InnerHandshakeReq) Reset() {
	*x = InnerHandshakeReq{}
	mi := &file_proto_innerHandshake_proto_msgTypes[0]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *InnerHandshakeReq) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*InnerHandshakeReq) ProtoMessage() {}

func (x *InnerHandshakeReq) ProtoReflect() protoreflect.Message {
	mi := &file_proto_innerHandshake_proto_msgTypes[0]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use InnerHandshakeReq.ProtoReflect.Descriptor instead.
func (*InnerHandshakeReq) Descriptor() ([]byte, []int) {
	return file_proto_innerHandshake_proto_rawDescGZIP(), []int{0}
}

func (x *InnerHandshakeReq) GetServiceId() string {
	if x != nil {
		return x.ServiceId
	}
	return ""
}

type InnerHandshakeResponse struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	Code          uint32                 `protobuf:"varint,1,opt,name=Code,proto3" json:"Code,omitempty"`
	Error         string                 `protobuf:"bytes,2,opt,name=Error,proto3" json:"Error,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *InnerHandshakeResponse) Reset() {
	*x = InnerHandshakeResponse{}
	mi := &file_proto_innerHandshake_proto_msgTypes[1]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *InnerHandshakeResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*InnerHandshakeResponse) ProtoMessage() {}

func (x *InnerHandshakeResponse) ProtoReflect() protoreflect.Message {
	mi := &file_proto_innerHandshake_proto_msgTypes[1]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use InnerHandshakeResponse.ProtoReflect.Descriptor instead.
func (*InnerHandshakeResponse) Descriptor() ([]byte, []int) {
	return file_proto_innerHandshake_proto_rawDescGZIP(), []int{1}
}

func (x *InnerHandshakeResponse) GetCode() uint32 {
	if x != nil {
		return x.Code
	}
	return 0
}

func (x *InnerHandshakeResponse) GetError() string {
	if x != nil {
		return x.Error
	}
	return ""
}

var File_proto_innerHandshake_proto protoreflect.FileDescriptor

const file_proto_innerHandshake_proto_rawDesc = "" +
	"\n" +
	"\x1aproto/innerHandshake.proto\x12\bprotocol\"1\n" +
	"\x11InnerHandshakeReq\x12\x1c\n" +
	"\tserviceId\x18\x01 \x01(\tR\tserviceId\"B\n" +
	"\x16InnerHandshakeResponse\x12\x12\n" +
	"\x04Code\x18\x01 \x01(\rR\x04Code\x12\x14\n" +
	"\x05Error\x18\x02 \x01(\tR\x05ErrorB\rZ\v./protoImplb\x06proto3"

var (
	file_proto_innerHandshake_proto_rawDescOnce sync.Once
	file_proto_innerHandshake_proto_rawDescData []byte
)

func file_proto_innerHandshake_proto_rawDescGZIP() []byte {
	file_proto_innerHandshake_proto_rawDescOnce.Do(func() {
		file_proto_innerHandshake_proto_rawDescData = protoimpl.X.CompressGZIP(unsafe.Slice(unsafe.StringData(file_proto_innerHandshake_proto_rawDesc), len(file_proto_innerHandshake_proto_rawDesc)))
	})
	return file_proto_innerHandshake_proto_rawDescData
}

var file_proto_innerHandshake_proto_msgTypes = make([]protoimpl.MessageInfo, 2)
var file_proto_innerHandshake_proto_goTypes = []any{
	(*InnerHandshakeReq)(nil),      // 0: protocol.InnerHandshakeReq
	(*InnerHandshakeResponse)(nil), // 1: protocol.InnerHandshakeResponse
}
var file_proto_innerHandshake_proto_depIdxs = []int32{
	0, // [0:0] is the sub-list for method output_type
	0, // [0:0] is the sub-list for method input_type
	0, // [0:0] is the sub-list for extension type_name
	0, // [0:0] is the sub-list for extension extendee
	0, // [0:0] is the sub-list for field type_name
}

func init() { file_proto_innerHandshake_proto_init() }
func file_proto_innerHandshake_proto_init() {
	if File_proto_innerHandshake_proto != nil {
		return
	}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: unsafe.Slice(unsafe.StringData(file_proto_innerHandshake_proto_rawDesc), len(file_proto_innerHandshake_proto_rawDesc)),
			NumEnums:      0,
			NumMessages:   2,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_proto_innerHandshake_proto_goTypes,
		DependencyIndexes: file_proto_innerHandshake_proto_depIdxs,
		MessageInfos:      file_proto_innerHandshake_proto_msgTypes,
	}.Build()
	File_proto_innerHandshake_proto = out.File
	file_proto_innerHandshake_proto_goTypes = nil
	file_proto_innerHandshake_proto_depIdxs = nil
}
