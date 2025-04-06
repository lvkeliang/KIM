// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.36.6
// 	protoc        v3.20.3
// source: proto/login.proto

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

type LoginReq struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	Token         string                 `protobuf:"bytes,1,opt,name=token,proto3" json:"token,omitempty"`
	Isp           string                 `protobuf:"bytes,2,opt,name=isp,proto3" json:"isp,omitempty"`
	Zone          string                 `protobuf:"bytes,3,opt,name=zone,proto3" json:"zone,omitempty"` // location code
	Tags          []string               `protobuf:"bytes,4,rep,name=tags,proto3" json:"tags,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *LoginReq) Reset() {
	*x = LoginReq{}
	mi := &file_proto_login_proto_msgTypes[0]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *LoginReq) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*LoginReq) ProtoMessage() {}

func (x *LoginReq) ProtoReflect() protoreflect.Message {
	mi := &file_proto_login_proto_msgTypes[0]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use LoginReq.ProtoReflect.Descriptor instead.
func (*LoginReq) Descriptor() ([]byte, []int) {
	return file_proto_login_proto_rawDescGZIP(), []int{0}
}

func (x *LoginReq) GetToken() string {
	if x != nil {
		return x.Token
	}
	return ""
}

func (x *LoginReq) GetIsp() string {
	if x != nil {
		return x.Isp
	}
	return ""
}

func (x *LoginReq) GetZone() string {
	if x != nil {
		return x.Zone
	}
	return ""
}

func (x *LoginReq) GetTags() []string {
	if x != nil {
		return x.Tags
	}
	return nil
}

type LoginResp struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	ChannelId     string                 `protobuf:"bytes,1,opt,name=channelId,proto3" json:"channelId,omitempty"`
	Account       string                 `protobuf:"bytes,2,opt,name=account,proto3" json:"account,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *LoginResp) Reset() {
	*x = LoginResp{}
	mi := &file_proto_login_proto_msgTypes[1]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *LoginResp) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*LoginResp) ProtoMessage() {}

func (x *LoginResp) ProtoReflect() protoreflect.Message {
	mi := &file_proto_login_proto_msgTypes[1]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use LoginResp.ProtoReflect.Descriptor instead.
func (*LoginResp) Descriptor() ([]byte, []int) {
	return file_proto_login_proto_rawDescGZIP(), []int{1}
}

func (x *LoginResp) GetChannelId() string {
	if x != nil {
		return x.ChannelId
	}
	return ""
}

func (x *LoginResp) GetAccount() string {
	if x != nil {
		return x.Account
	}
	return ""
}

type KickOutNotify struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	ChannelId     string                 `protobuf:"bytes,1,opt,name=channelId,proto3" json:"channelId,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *KickOutNotify) Reset() {
	*x = KickOutNotify{}
	mi := &file_proto_login_proto_msgTypes[2]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *KickOutNotify) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*KickOutNotify) ProtoMessage() {}

func (x *KickOutNotify) ProtoReflect() protoreflect.Message {
	mi := &file_proto_login_proto_msgTypes[2]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use KickOutNotify.ProtoReflect.Descriptor instead.
func (*KickOutNotify) Descriptor() ([]byte, []int) {
	return file_proto_login_proto_rawDescGZIP(), []int{2}
}

func (x *KickOutNotify) GetChannelId() string {
	if x != nil {
		return x.ChannelId
	}
	return ""
}

var File_proto_login_proto protoreflect.FileDescriptor

const file_proto_login_proto_rawDesc = "" +
	"\n" +
	"\x11proto/login.proto\x12\bprotocol\"Z\n" +
	"\bLoginReq\x12\x14\n" +
	"\x05token\x18\x01 \x01(\tR\x05token\x12\x10\n" +
	"\x03isp\x18\x02 \x01(\tR\x03isp\x12\x12\n" +
	"\x04zone\x18\x03 \x01(\tR\x04zone\x12\x12\n" +
	"\x04tags\x18\x04 \x03(\tR\x04tags\"C\n" +
	"\tLoginResp\x12\x1c\n" +
	"\tchannelId\x18\x01 \x01(\tR\tchannelId\x12\x18\n" +
	"\aaccount\x18\x02 \x01(\tR\aaccount\"-\n" +
	"\rKickOutNotify\x12\x1c\n" +
	"\tchannelId\x18\x01 \x01(\tR\tchannelIdB\rZ\v./protoImplb\x06proto3"

var (
	file_proto_login_proto_rawDescOnce sync.Once
	file_proto_login_proto_rawDescData []byte
)

func file_proto_login_proto_rawDescGZIP() []byte {
	file_proto_login_proto_rawDescOnce.Do(func() {
		file_proto_login_proto_rawDescData = protoimpl.X.CompressGZIP(unsafe.Slice(unsafe.StringData(file_proto_login_proto_rawDesc), len(file_proto_login_proto_rawDesc)))
	})
	return file_proto_login_proto_rawDescData
}

var file_proto_login_proto_msgTypes = make([]protoimpl.MessageInfo, 3)
var file_proto_login_proto_goTypes = []any{
	(*LoginReq)(nil),      // 0: protocol.LoginReq
	(*LoginResp)(nil),     // 1: protocol.LoginResp
	(*KickOutNotify)(nil), // 2: protocol.KickOutNotify
}
var file_proto_login_proto_depIdxs = []int32{
	0, // [0:0] is the sub-list for method output_type
	0, // [0:0] is the sub-list for method input_type
	0, // [0:0] is the sub-list for extension type_name
	0, // [0:0] is the sub-list for extension extendee
	0, // [0:0] is the sub-list for field type_name
}

func init() { file_proto_login_proto_init() }
func file_proto_login_proto_init() {
	if File_proto_login_proto != nil {
		return
	}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: unsafe.Slice(unsafe.StringData(file_proto_login_proto_rawDesc), len(file_proto_login_proto_rawDesc)),
			NumEnums:      0,
			NumMessages:   3,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_proto_login_proto_goTypes,
		DependencyIndexes: file_proto_login_proto_depIdxs,
		MessageInfos:      file_proto_login_proto_msgTypes,
	}.Build()
	File_proto_login_proto = out.File
	file_proto_login_proto_goTypes = nil
	file_proto_login_proto_depIdxs = nil
}
