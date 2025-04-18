// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.36.6
// 	protoc        v3.20.3
// source: proto/rpc.proto

package rpc

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

type User struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	Account       string                 `protobuf:"bytes,1,opt,name=account,proto3" json:"account,omitempty"`
	Alias         string                 `protobuf:"bytes,2,opt,name=alias,proto3" json:"alias,omitempty"`
	Avatar        string                 `protobuf:"bytes,3,opt,name=avatar,proto3" json:"avatar,omitempty"`
	CreatedAt     int64                  `protobuf:"varint,4,opt,name=created_at,json=createdAt,proto3" json:"created_at,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *User) Reset() {
	*x = User{}
	mi := &file_proto_rpc_proto_msgTypes[0]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *User) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*User) ProtoMessage() {}

func (x *User) ProtoReflect() protoreflect.Message {
	mi := &file_proto_rpc_proto_msgTypes[0]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use User.ProtoReflect.Descriptor instead.
func (*User) Descriptor() ([]byte, []int) {
	return file_proto_rpc_proto_rawDescGZIP(), []int{0}
}

func (x *User) GetAccount() string {
	if x != nil {
		return x.Account
	}
	return ""
}

func (x *User) GetAlias() string {
	if x != nil {
		return x.Alias
	}
	return ""
}

func (x *User) GetAvatar() string {
	if x != nil {
		return x.Avatar
	}
	return ""
}

func (x *User) GetCreatedAt() int64 {
	if x != nil {
		return x.CreatedAt
	}
	return 0
}

type Message struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	Id            int64                  `protobuf:"varint,1,opt,name=id,proto3" json:"id,omitempty"`
	Type          int32                  `protobuf:"varint,2,opt,name=type,proto3" json:"type,omitempty"`
	Body          string                 `protobuf:"bytes,3,opt,name=body,proto3" json:"body,omitempty"`
	Extra         string                 `protobuf:"bytes,4,opt,name=extra,proto3" json:"extra,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *Message) Reset() {
	*x = Message{}
	mi := &file_proto_rpc_proto_msgTypes[1]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *Message) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Message) ProtoMessage() {}

func (x *Message) ProtoReflect() protoreflect.Message {
	mi := &file_proto_rpc_proto_msgTypes[1]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Message.ProtoReflect.Descriptor instead.
func (*Message) Descriptor() ([]byte, []int) {
	return file_proto_rpc_proto_rawDescGZIP(), []int{1}
}

func (x *Message) GetId() int64 {
	if x != nil {
		return x.Id
	}
	return 0
}

func (x *Message) GetType() int32 {
	if x != nil {
		return x.Type
	}
	return 0
}

func (x *Message) GetBody() string {
	if x != nil {
		return x.Body
	}
	return ""
}

func (x *Message) GetExtra() string {
	if x != nil {
		return x.Extra
	}
	return ""
}

type Member struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	Account       string                 `protobuf:"bytes,1,opt,name=account,proto3" json:"account,omitempty"`
	Alias         string                 `protobuf:"bytes,2,opt,name=alias,proto3" json:"alias,omitempty"`
	Avatar        string                 `protobuf:"bytes,3,opt,name=avatar,proto3" json:"avatar,omitempty"`
	JoinTime      int64                  `protobuf:"varint,4,opt,name=join_time,json=joinTime,proto3" json:"join_time,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *Member) Reset() {
	*x = Member{}
	mi := &file_proto_rpc_proto_msgTypes[2]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *Member) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Member) ProtoMessage() {}

func (x *Member) ProtoReflect() protoreflect.Message {
	mi := &file_proto_rpc_proto_msgTypes[2]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Member.ProtoReflect.Descriptor instead.
func (*Member) Descriptor() ([]byte, []int) {
	return file_proto_rpc_proto_rawDescGZIP(), []int{2}
}

func (x *Member) GetAccount() string {
	if x != nil {
		return x.Account
	}
	return ""
}

func (x *Member) GetAlias() string {
	if x != nil {
		return x.Alias
	}
	return ""
}

func (x *Member) GetAvatar() string {
	if x != nil {
		return x.Avatar
	}
	return ""
}

func (x *Member) GetJoinTime() int64 {
	if x != nil {
		return x.JoinTime
	}
	return 0
}

type InsertMessageReq struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	Sender        string                 `protobuf:"bytes,1,opt,name=sender,proto3" json:"sender,omitempty"`
	Dest          string                 `protobuf:"bytes,2,opt,name=dest,proto3" json:"dest,omitempty"`
	SendTime      int64                  `protobuf:"varint,3,opt,name=send_time,json=sendTime,proto3" json:"send_time,omitempty"`
	Message       *Message               `protobuf:"bytes,4,opt,name=message,proto3" json:"message,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *InsertMessageReq) Reset() {
	*x = InsertMessageReq{}
	mi := &file_proto_rpc_proto_msgTypes[3]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *InsertMessageReq) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*InsertMessageReq) ProtoMessage() {}

func (x *InsertMessageReq) ProtoReflect() protoreflect.Message {
	mi := &file_proto_rpc_proto_msgTypes[3]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use InsertMessageReq.ProtoReflect.Descriptor instead.
func (*InsertMessageReq) Descriptor() ([]byte, []int) {
	return file_proto_rpc_proto_rawDescGZIP(), []int{3}
}

func (x *InsertMessageReq) GetSender() string {
	if x != nil {
		return x.Sender
	}
	return ""
}

func (x *InsertMessageReq) GetDest() string {
	if x != nil {
		return x.Dest
	}
	return ""
}

func (x *InsertMessageReq) GetSendTime() int64 {
	if x != nil {
		return x.SendTime
	}
	return 0
}

func (x *InsertMessageReq) GetMessage() *Message {
	if x != nil {
		return x.Message
	}
	return nil
}

type InsertMessageResp struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	MessageId     int64                  `protobuf:"varint,1,opt,name=message_id,json=messageId,proto3" json:"message_id,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *InsertMessageResp) Reset() {
	*x = InsertMessageResp{}
	mi := &file_proto_rpc_proto_msgTypes[4]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *InsertMessageResp) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*InsertMessageResp) ProtoMessage() {}

func (x *InsertMessageResp) ProtoReflect() protoreflect.Message {
	mi := &file_proto_rpc_proto_msgTypes[4]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use InsertMessageResp.ProtoReflect.Descriptor instead.
func (*InsertMessageResp) Descriptor() ([]byte, []int) {
	return file_proto_rpc_proto_rawDescGZIP(), []int{4}
}

func (x *InsertMessageResp) GetMessageId() int64 {
	if x != nil {
		return x.MessageId
	}
	return 0
}

type AckMessageReq struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	Account       string                 `protobuf:"bytes,1,opt,name=account,proto3" json:"account,omitempty"`
	MessageId     int64                  `protobuf:"varint,2,opt,name=message_id,json=messageId,proto3" json:"message_id,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *AckMessageReq) Reset() {
	*x = AckMessageReq{}
	mi := &file_proto_rpc_proto_msgTypes[5]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *AckMessageReq) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*AckMessageReq) ProtoMessage() {}

func (x *AckMessageReq) ProtoReflect() protoreflect.Message {
	mi := &file_proto_rpc_proto_msgTypes[5]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use AckMessageReq.ProtoReflect.Descriptor instead.
func (*AckMessageReq) Descriptor() ([]byte, []int) {
	return file_proto_rpc_proto_rawDescGZIP(), []int{5}
}

func (x *AckMessageReq) GetAccount() string {
	if x != nil {
		return x.Account
	}
	return ""
}

func (x *AckMessageReq) GetMessageId() int64 {
	if x != nil {
		return x.MessageId
	}
	return 0
}

type CreateGroupReq struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	App           string                 `protobuf:"bytes,1,opt,name=app,proto3" json:"app,omitempty"`
	Name          string                 `protobuf:"bytes,2,opt,name=name,proto3" json:"name,omitempty"`
	Avatar        string                 `protobuf:"bytes,3,opt,name=avatar,proto3" json:"avatar,omitempty"`
	Introduction  string                 `protobuf:"bytes,4,opt,name=introduction,proto3" json:"introduction,omitempty"`
	Owner         string                 `protobuf:"bytes,5,opt,name=owner,proto3" json:"owner,omitempty"`
	Members       []string               `protobuf:"bytes,6,rep,name=members,proto3" json:"members,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *CreateGroupReq) Reset() {
	*x = CreateGroupReq{}
	mi := &file_proto_rpc_proto_msgTypes[6]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *CreateGroupReq) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*CreateGroupReq) ProtoMessage() {}

func (x *CreateGroupReq) ProtoReflect() protoreflect.Message {
	mi := &file_proto_rpc_proto_msgTypes[6]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use CreateGroupReq.ProtoReflect.Descriptor instead.
func (*CreateGroupReq) Descriptor() ([]byte, []int) {
	return file_proto_rpc_proto_rawDescGZIP(), []int{6}
}

func (x *CreateGroupReq) GetApp() string {
	if x != nil {
		return x.App
	}
	return ""
}

func (x *CreateGroupReq) GetName() string {
	if x != nil {
		return x.Name
	}
	return ""
}

func (x *CreateGroupReq) GetAvatar() string {
	if x != nil {
		return x.Avatar
	}
	return ""
}

func (x *CreateGroupReq) GetIntroduction() string {
	if x != nil {
		return x.Introduction
	}
	return ""
}

func (x *CreateGroupReq) GetOwner() string {
	if x != nil {
		return x.Owner
	}
	return ""
}

func (x *CreateGroupReq) GetMembers() []string {
	if x != nil {
		return x.Members
	}
	return nil
}

type CreateGroupResp struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	GroupId       string                 `protobuf:"bytes,1,opt,name=group_id,json=groupId,proto3" json:"group_id,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *CreateGroupResp) Reset() {
	*x = CreateGroupResp{}
	mi := &file_proto_rpc_proto_msgTypes[7]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *CreateGroupResp) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*CreateGroupResp) ProtoMessage() {}

func (x *CreateGroupResp) ProtoReflect() protoreflect.Message {
	mi := &file_proto_rpc_proto_msgTypes[7]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use CreateGroupResp.ProtoReflect.Descriptor instead.
func (*CreateGroupResp) Descriptor() ([]byte, []int) {
	return file_proto_rpc_proto_rawDescGZIP(), []int{7}
}

func (x *CreateGroupResp) GetGroupId() string {
	if x != nil {
		return x.GroupId
	}
	return ""
}

type JoinGroupReq struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	Account       string                 `protobuf:"bytes,1,opt,name=account,proto3" json:"account,omitempty"`
	GroupId       string                 `protobuf:"bytes,2,opt,name=group_id,json=groupId,proto3" json:"group_id,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *JoinGroupReq) Reset() {
	*x = JoinGroupReq{}
	mi := &file_proto_rpc_proto_msgTypes[8]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *JoinGroupReq) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*JoinGroupReq) ProtoMessage() {}

func (x *JoinGroupReq) ProtoReflect() protoreflect.Message {
	mi := &file_proto_rpc_proto_msgTypes[8]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use JoinGroupReq.ProtoReflect.Descriptor instead.
func (*JoinGroupReq) Descriptor() ([]byte, []int) {
	return file_proto_rpc_proto_rawDescGZIP(), []int{8}
}

func (x *JoinGroupReq) GetAccount() string {
	if x != nil {
		return x.Account
	}
	return ""
}

func (x *JoinGroupReq) GetGroupId() string {
	if x != nil {
		return x.GroupId
	}
	return ""
}

type QuitGroupReq struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	Account       string                 `protobuf:"bytes,1,opt,name=account,proto3" json:"account,omitempty"`
	GroupId       string                 `protobuf:"bytes,2,opt,name=group_id,json=groupId,proto3" json:"group_id,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *QuitGroupReq) Reset() {
	*x = QuitGroupReq{}
	mi := &file_proto_rpc_proto_msgTypes[9]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *QuitGroupReq) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*QuitGroupReq) ProtoMessage() {}

func (x *QuitGroupReq) ProtoReflect() protoreflect.Message {
	mi := &file_proto_rpc_proto_msgTypes[9]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use QuitGroupReq.ProtoReflect.Descriptor instead.
func (*QuitGroupReq) Descriptor() ([]byte, []int) {
	return file_proto_rpc_proto_rawDescGZIP(), []int{9}
}

func (x *QuitGroupReq) GetAccount() string {
	if x != nil {
		return x.Account
	}
	return ""
}

func (x *QuitGroupReq) GetGroupId() string {
	if x != nil {
		return x.GroupId
	}
	return ""
}

type GetGroupReq struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	GroupId       string                 `protobuf:"bytes,1,opt,name=group_id,json=groupId,proto3" json:"group_id,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *GetGroupReq) Reset() {
	*x = GetGroupReq{}
	mi := &file_proto_rpc_proto_msgTypes[10]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *GetGroupReq) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GetGroupReq) ProtoMessage() {}

func (x *GetGroupReq) ProtoReflect() protoreflect.Message {
	mi := &file_proto_rpc_proto_msgTypes[10]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use GetGroupReq.ProtoReflect.Descriptor instead.
func (*GetGroupReq) Descriptor() ([]byte, []int) {
	return file_proto_rpc_proto_rawDescGZIP(), []int{10}
}

func (x *GetGroupReq) GetGroupId() string {
	if x != nil {
		return x.GroupId
	}
	return ""
}

type GetGroupResp struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	Id            string                 `protobuf:"bytes,1,opt,name=id,proto3" json:"id,omitempty"`
	Name          string                 `protobuf:"bytes,2,opt,name=name,proto3" json:"name,omitempty"`
	Avatar        string                 `protobuf:"bytes,3,opt,name=avatar,proto3" json:"avatar,omitempty"`
	Introduction  string                 `protobuf:"bytes,4,opt,name=introduction,proto3" json:"introduction,omitempty"`
	Owner         string                 `protobuf:"bytes,5,opt,name=owner,proto3" json:"owner,omitempty"`
	CreatedAt     int64                  `protobuf:"varint,6,opt,name=created_at,json=createdAt,proto3" json:"created_at,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *GetGroupResp) Reset() {
	*x = GetGroupResp{}
	mi := &file_proto_rpc_proto_msgTypes[11]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *GetGroupResp) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GetGroupResp) ProtoMessage() {}

func (x *GetGroupResp) ProtoReflect() protoreflect.Message {
	mi := &file_proto_rpc_proto_msgTypes[11]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use GetGroupResp.ProtoReflect.Descriptor instead.
func (*GetGroupResp) Descriptor() ([]byte, []int) {
	return file_proto_rpc_proto_rawDescGZIP(), []int{11}
}

func (x *GetGroupResp) GetId() string {
	if x != nil {
		return x.Id
	}
	return ""
}

func (x *GetGroupResp) GetName() string {
	if x != nil {
		return x.Name
	}
	return ""
}

func (x *GetGroupResp) GetAvatar() string {
	if x != nil {
		return x.Avatar
	}
	return ""
}

func (x *GetGroupResp) GetIntroduction() string {
	if x != nil {
		return x.Introduction
	}
	return ""
}

func (x *GetGroupResp) GetOwner() string {
	if x != nil {
		return x.Owner
	}
	return ""
}

func (x *GetGroupResp) GetCreatedAt() int64 {
	if x != nil {
		return x.CreatedAt
	}
	return 0
}

type GroupMembersReq struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	GroupId       string                 `protobuf:"bytes,1,opt,name=group_id,json=groupId,proto3" json:"group_id,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *GroupMembersReq) Reset() {
	*x = GroupMembersReq{}
	mi := &file_proto_rpc_proto_msgTypes[12]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *GroupMembersReq) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GroupMembersReq) ProtoMessage() {}

func (x *GroupMembersReq) ProtoReflect() protoreflect.Message {
	mi := &file_proto_rpc_proto_msgTypes[12]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use GroupMembersReq.ProtoReflect.Descriptor instead.
func (*GroupMembersReq) Descriptor() ([]byte, []int) {
	return file_proto_rpc_proto_rawDescGZIP(), []int{12}
}

func (x *GroupMembersReq) GetGroupId() string {
	if x != nil {
		return x.GroupId
	}
	return ""
}

type GroupMembersResp struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	Users         []*Member              `protobuf:"bytes,1,rep,name=users,proto3" json:"users,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *GroupMembersResp) Reset() {
	*x = GroupMembersResp{}
	mi := &file_proto_rpc_proto_msgTypes[13]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *GroupMembersResp) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GroupMembersResp) ProtoMessage() {}

func (x *GroupMembersResp) ProtoReflect() protoreflect.Message {
	mi := &file_proto_rpc_proto_msgTypes[13]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use GroupMembersResp.ProtoReflect.Descriptor instead.
func (*GroupMembersResp) Descriptor() ([]byte, []int) {
	return file_proto_rpc_proto_rawDescGZIP(), []int{13}
}

func (x *GroupMembersResp) GetUsers() []*Member {
	if x != nil {
		return x.Users
	}
	return nil
}

type GetOfflineMessageIndexReq struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	Account       string                 `protobuf:"bytes,1,opt,name=account,proto3" json:"account,omitempty"`
	MessageId     int64                  `protobuf:"varint,2,opt,name=message_id,json=messageId,proto3" json:"message_id,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *GetOfflineMessageIndexReq) Reset() {
	*x = GetOfflineMessageIndexReq{}
	mi := &file_proto_rpc_proto_msgTypes[14]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *GetOfflineMessageIndexReq) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GetOfflineMessageIndexReq) ProtoMessage() {}

func (x *GetOfflineMessageIndexReq) ProtoReflect() protoreflect.Message {
	mi := &file_proto_rpc_proto_msgTypes[14]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use GetOfflineMessageIndexReq.ProtoReflect.Descriptor instead.
func (*GetOfflineMessageIndexReq) Descriptor() ([]byte, []int) {
	return file_proto_rpc_proto_rawDescGZIP(), []int{14}
}

func (x *GetOfflineMessageIndexReq) GetAccount() string {
	if x != nil {
		return x.Account
	}
	return ""
}

func (x *GetOfflineMessageIndexReq) GetMessageId() int64 {
	if x != nil {
		return x.MessageId
	}
	return 0
}

type GetOfflineMessageIndexResp struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	List          []*MessageIndex        `protobuf:"bytes,1,rep,name=list,proto3" json:"list,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *GetOfflineMessageIndexResp) Reset() {
	*x = GetOfflineMessageIndexResp{}
	mi := &file_proto_rpc_proto_msgTypes[15]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *GetOfflineMessageIndexResp) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GetOfflineMessageIndexResp) ProtoMessage() {}

func (x *GetOfflineMessageIndexResp) ProtoReflect() protoreflect.Message {
	mi := &file_proto_rpc_proto_msgTypes[15]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use GetOfflineMessageIndexResp.ProtoReflect.Descriptor instead.
func (*GetOfflineMessageIndexResp) Descriptor() ([]byte, []int) {
	return file_proto_rpc_proto_rawDescGZIP(), []int{15}
}

func (x *GetOfflineMessageIndexResp) GetList() []*MessageIndex {
	if x != nil {
		return x.List
	}
	return nil
}

type MessageIndex struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	MessageId     int64                  `protobuf:"varint,1,opt,name=message_id,json=messageId,proto3" json:"message_id,omitempty"`
	Direction     int32                  `protobuf:"varint,2,opt,name=direction,proto3" json:"direction,omitempty"`
	SendTime      int64                  `protobuf:"varint,3,opt,name=send_time,json=sendTime,proto3" json:"send_time,omitempty"`
	AccountB      string                 `protobuf:"bytes,4,opt,name=accountB,proto3" json:"accountB,omitempty"`
	Group         string                 `protobuf:"bytes,5,opt,name=group,proto3" json:"group,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *MessageIndex) Reset() {
	*x = MessageIndex{}
	mi := &file_proto_rpc_proto_msgTypes[16]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *MessageIndex) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*MessageIndex) ProtoMessage() {}

func (x *MessageIndex) ProtoReflect() protoreflect.Message {
	mi := &file_proto_rpc_proto_msgTypes[16]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use MessageIndex.ProtoReflect.Descriptor instead.
func (*MessageIndex) Descriptor() ([]byte, []int) {
	return file_proto_rpc_proto_rawDescGZIP(), []int{16}
}

func (x *MessageIndex) GetMessageId() int64 {
	if x != nil {
		return x.MessageId
	}
	return 0
}

func (x *MessageIndex) GetDirection() int32 {
	if x != nil {
		return x.Direction
	}
	return 0
}

func (x *MessageIndex) GetSendTime() int64 {
	if x != nil {
		return x.SendTime
	}
	return 0
}

func (x *MessageIndex) GetAccountB() string {
	if x != nil {
		return x.AccountB
	}
	return ""
}

func (x *MessageIndex) GetGroup() string {
	if x != nil {
		return x.Group
	}
	return ""
}

type GetOfflineMessageContentReq struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	MessageIds    []int64                `protobuf:"varint,1,rep,packed,name=message_ids,json=messageIds,proto3" json:"message_ids,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *GetOfflineMessageContentReq) Reset() {
	*x = GetOfflineMessageContentReq{}
	mi := &file_proto_rpc_proto_msgTypes[17]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *GetOfflineMessageContentReq) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GetOfflineMessageContentReq) ProtoMessage() {}

func (x *GetOfflineMessageContentReq) ProtoReflect() protoreflect.Message {
	mi := &file_proto_rpc_proto_msgTypes[17]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use GetOfflineMessageContentReq.ProtoReflect.Descriptor instead.
func (*GetOfflineMessageContentReq) Descriptor() ([]byte, []int) {
	return file_proto_rpc_proto_rawDescGZIP(), []int{17}
}

func (x *GetOfflineMessageContentReq) GetMessageIds() []int64 {
	if x != nil {
		return x.MessageIds
	}
	return nil
}

type GetOfflineMessageContentResp struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	List          []*Message             `protobuf:"bytes,1,rep,name=list,proto3" json:"list,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *GetOfflineMessageContentResp) Reset() {
	*x = GetOfflineMessageContentResp{}
	mi := &file_proto_rpc_proto_msgTypes[18]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *GetOfflineMessageContentResp) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GetOfflineMessageContentResp) ProtoMessage() {}

func (x *GetOfflineMessageContentResp) ProtoReflect() protoreflect.Message {
	mi := &file_proto_rpc_proto_msgTypes[18]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use GetOfflineMessageContentResp.ProtoReflect.Descriptor instead.
func (*GetOfflineMessageContentResp) Descriptor() ([]byte, []int) {
	return file_proto_rpc_proto_rawDescGZIP(), []int{18}
}

func (x *GetOfflineMessageContentResp) GetList() []*Message {
	if x != nil {
		return x.List
	}
	return nil
}

var File_proto_rpc_proto protoreflect.FileDescriptor

const file_proto_rpc_proto_rawDesc = "" +
	"\n" +
	"\x0fproto/rpc.proto\x12\x03rpc\"m\n" +
	"\x04User\x12\x18\n" +
	"\aaccount\x18\x01 \x01(\tR\aaccount\x12\x14\n" +
	"\x05alias\x18\x02 \x01(\tR\x05alias\x12\x16\n" +
	"\x06avatar\x18\x03 \x01(\tR\x06avatar\x12\x1d\n" +
	"\n" +
	"created_at\x18\x04 \x01(\x03R\tcreatedAt\"W\n" +
	"\aMessage\x12\x0e\n" +
	"\x02id\x18\x01 \x01(\x03R\x02id\x12\x12\n" +
	"\x04type\x18\x02 \x01(\x05R\x04type\x12\x12\n" +
	"\x04body\x18\x03 \x01(\tR\x04body\x12\x14\n" +
	"\x05extra\x18\x04 \x01(\tR\x05extra\"m\n" +
	"\x06Member\x12\x18\n" +
	"\aaccount\x18\x01 \x01(\tR\aaccount\x12\x14\n" +
	"\x05alias\x18\x02 \x01(\tR\x05alias\x12\x16\n" +
	"\x06avatar\x18\x03 \x01(\tR\x06avatar\x12\x1b\n" +
	"\tjoin_time\x18\x04 \x01(\x03R\bjoinTime\"\x83\x01\n" +
	"\x10InsertMessageReq\x12\x16\n" +
	"\x06sender\x18\x01 \x01(\tR\x06sender\x12\x12\n" +
	"\x04dest\x18\x02 \x01(\tR\x04dest\x12\x1b\n" +
	"\tsend_time\x18\x03 \x01(\x03R\bsendTime\x12&\n" +
	"\amessage\x18\x04 \x01(\v2\f.rpc.MessageR\amessage\"2\n" +
	"\x11InsertMessageResp\x12\x1d\n" +
	"\n" +
	"message_id\x18\x01 \x01(\x03R\tmessageId\"H\n" +
	"\rAckMessageReq\x12\x18\n" +
	"\aaccount\x18\x01 \x01(\tR\aaccount\x12\x1d\n" +
	"\n" +
	"message_id\x18\x02 \x01(\x03R\tmessageId\"\xa2\x01\n" +
	"\x0eCreateGroupReq\x12\x10\n" +
	"\x03app\x18\x01 \x01(\tR\x03app\x12\x12\n" +
	"\x04name\x18\x02 \x01(\tR\x04name\x12\x16\n" +
	"\x06avatar\x18\x03 \x01(\tR\x06avatar\x12\"\n" +
	"\fintroduction\x18\x04 \x01(\tR\fintroduction\x12\x14\n" +
	"\x05owner\x18\x05 \x01(\tR\x05owner\x12\x18\n" +
	"\amembers\x18\x06 \x03(\tR\amembers\",\n" +
	"\x0fCreateGroupResp\x12\x19\n" +
	"\bgroup_id\x18\x01 \x01(\tR\agroupId\"C\n" +
	"\fJoinGroupReq\x12\x18\n" +
	"\aaccount\x18\x01 \x01(\tR\aaccount\x12\x19\n" +
	"\bgroup_id\x18\x02 \x01(\tR\agroupId\"C\n" +
	"\fQuitGroupReq\x12\x18\n" +
	"\aaccount\x18\x01 \x01(\tR\aaccount\x12\x19\n" +
	"\bgroup_id\x18\x02 \x01(\tR\agroupId\"(\n" +
	"\vGetGroupReq\x12\x19\n" +
	"\bgroup_id\x18\x01 \x01(\tR\agroupId\"\xa3\x01\n" +
	"\fGetGroupResp\x12\x0e\n" +
	"\x02id\x18\x01 \x01(\tR\x02id\x12\x12\n" +
	"\x04name\x18\x02 \x01(\tR\x04name\x12\x16\n" +
	"\x06avatar\x18\x03 \x01(\tR\x06avatar\x12\"\n" +
	"\fintroduction\x18\x04 \x01(\tR\fintroduction\x12\x14\n" +
	"\x05owner\x18\x05 \x01(\tR\x05owner\x12\x1d\n" +
	"\n" +
	"created_at\x18\x06 \x01(\x03R\tcreatedAt\",\n" +
	"\x0fGroupMembersReq\x12\x19\n" +
	"\bgroup_id\x18\x01 \x01(\tR\agroupId\"5\n" +
	"\x10GroupMembersResp\x12!\n" +
	"\x05users\x18\x01 \x03(\v2\v.rpc.MemberR\x05users\"T\n" +
	"\x19GetOfflineMessageIndexReq\x12\x18\n" +
	"\aaccount\x18\x01 \x01(\tR\aaccount\x12\x1d\n" +
	"\n" +
	"message_id\x18\x02 \x01(\x03R\tmessageId\"C\n" +
	"\x1aGetOfflineMessageIndexResp\x12%\n" +
	"\x04list\x18\x01 \x03(\v2\x11.rpc.MessageIndexR\x04list\"\x9a\x01\n" +
	"\fMessageIndex\x12\x1d\n" +
	"\n" +
	"message_id\x18\x01 \x01(\x03R\tmessageId\x12\x1c\n" +
	"\tdirection\x18\x02 \x01(\x05R\tdirection\x12\x1b\n" +
	"\tsend_time\x18\x03 \x01(\x03R\bsendTime\x12\x1a\n" +
	"\baccountB\x18\x04 \x01(\tR\baccountB\x12\x14\n" +
	"\x05group\x18\x05 \x01(\tR\x05group\">\n" +
	"\x1bGetOfflineMessageContentReq\x12\x1f\n" +
	"\vmessage_ids\x18\x01 \x03(\x03R\n" +
	"messageIds\"@\n" +
	"\x1cGetOfflineMessageContentResp\x12 \n" +
	"\x04list\x18\x01 \x03(\v2\f.rpc.MessageR\x04listB\aZ\x05./rpcb\x06proto3"

var (
	file_proto_rpc_proto_rawDescOnce sync.Once
	file_proto_rpc_proto_rawDescData []byte
)

func file_proto_rpc_proto_rawDescGZIP() []byte {
	file_proto_rpc_proto_rawDescOnce.Do(func() {
		file_proto_rpc_proto_rawDescData = protoimpl.X.CompressGZIP(unsafe.Slice(unsafe.StringData(file_proto_rpc_proto_rawDesc), len(file_proto_rpc_proto_rawDesc)))
	})
	return file_proto_rpc_proto_rawDescData
}

var file_proto_rpc_proto_msgTypes = make([]protoimpl.MessageInfo, 19)
var file_proto_rpc_proto_goTypes = []any{
	(*User)(nil),                         // 0: rpc.User
	(*Message)(nil),                      // 1: rpc.Message
	(*Member)(nil),                       // 2: rpc.Member
	(*InsertMessageReq)(nil),             // 3: rpc.InsertMessageReq
	(*InsertMessageResp)(nil),            // 4: rpc.InsertMessageResp
	(*AckMessageReq)(nil),                // 5: rpc.AckMessageReq
	(*CreateGroupReq)(nil),               // 6: rpc.CreateGroupReq
	(*CreateGroupResp)(nil),              // 7: rpc.CreateGroupResp
	(*JoinGroupReq)(nil),                 // 8: rpc.JoinGroupReq
	(*QuitGroupReq)(nil),                 // 9: rpc.QuitGroupReq
	(*GetGroupReq)(nil),                  // 10: rpc.GetGroupReq
	(*GetGroupResp)(nil),                 // 11: rpc.GetGroupResp
	(*GroupMembersReq)(nil),              // 12: rpc.GroupMembersReq
	(*GroupMembersResp)(nil),             // 13: rpc.GroupMembersResp
	(*GetOfflineMessageIndexReq)(nil),    // 14: rpc.GetOfflineMessageIndexReq
	(*GetOfflineMessageIndexResp)(nil),   // 15: rpc.GetOfflineMessageIndexResp
	(*MessageIndex)(nil),                 // 16: rpc.MessageIndex
	(*GetOfflineMessageContentReq)(nil),  // 17: rpc.GetOfflineMessageContentReq
	(*GetOfflineMessageContentResp)(nil), // 18: rpc.GetOfflineMessageContentResp
}
var file_proto_rpc_proto_depIdxs = []int32{
	1,  // 0: rpc.InsertMessageReq.message:type_name -> rpc.Message
	2,  // 1: rpc.GroupMembersResp.users:type_name -> rpc.Member
	16, // 2: rpc.GetOfflineMessageIndexResp.list:type_name -> rpc.MessageIndex
	1,  // 3: rpc.GetOfflineMessageContentResp.list:type_name -> rpc.Message
	4,  // [4:4] is the sub-list for method output_type
	4,  // [4:4] is the sub-list for method input_type
	4,  // [4:4] is the sub-list for extension type_name
	4,  // [4:4] is the sub-list for extension extendee
	0,  // [0:4] is the sub-list for field type_name
}

func init() { file_proto_rpc_proto_init() }
func file_proto_rpc_proto_init() {
	if File_proto_rpc_proto != nil {
		return
	}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: unsafe.Slice(unsafe.StringData(file_proto_rpc_proto_rawDesc), len(file_proto_rpc_proto_rawDesc)),
			NumEnums:      0,
			NumMessages:   19,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_proto_rpc_proto_goTypes,
		DependencyIndexes: file_proto_rpc_proto_depIdxs,
		MessageInfos:      file_proto_rpc_proto_msgTypes,
	}.Build()
	File_proto_rpc_proto = out.File
	file_proto_rpc_proto_goTypes = nil
	file_proto_rpc_proto_depIdxs = nil
}
