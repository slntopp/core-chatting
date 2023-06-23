// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.28.1
// 	protoc        (unknown)
// source: cc/cc.proto

package cc

import (
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

type Role int32

const (
	Role_NOACCESS Role = 0
	Role_USER     Role = 1
	Role_OWNER    Role = 2
	Role_ADMIN    Role = 3
)

// Enum value maps for Role.
var (
	Role_name = map[int32]string{
		0: "NOACCESS",
		1: "USER",
		2: "OWNER",
		3: "ADMIN",
	}
	Role_value = map[string]int32{
		"NOACCESS": 0,
		"USER":     1,
		"OWNER":    2,
		"ADMIN":    3,
	}
)

func (x Role) Enum() *Role {
	p := new(Role)
	*p = x
	return p
}

func (x Role) String() string {
	return protoimpl.X.EnumStringOf(x.Descriptor(), protoreflect.EnumNumber(x))
}

func (Role) Descriptor() protoreflect.EnumDescriptor {
	return file_cc_cc_proto_enumTypes[0].Descriptor()
}

func (Role) Type() protoreflect.EnumType {
	return &file_cc_cc_proto_enumTypes[0]
}

func (x Role) Number() protoreflect.EnumNumber {
	return protoreflect.EnumNumber(x)
}

// Deprecated: Use Role.Descriptor instead.
func (Role) EnumDescriptor() ([]byte, []int) {
	return file_cc_cc_proto_rawDescGZIP(), []int{0}
}

type Kind int32

const (
	Kind_DEFAULT    Kind = 0
	Kind_ADMIN_ONLY Kind = 1
)

// Enum value maps for Kind.
var (
	Kind_name = map[int32]string{
		0: "DEFAULT",
		1: "ADMIN_ONLY",
	}
	Kind_value = map[string]int32{
		"DEFAULT":    0,
		"ADMIN_ONLY": 1,
	}
)

func (x Kind) Enum() *Kind {
	p := new(Kind)
	*p = x
	return p
}

func (x Kind) String() string {
	return protoimpl.X.EnumStringOf(x.Descriptor(), protoreflect.EnumNumber(x))
}

func (Kind) Descriptor() protoreflect.EnumDescriptor {
	return file_cc_cc_proto_enumTypes[1].Descriptor()
}

func (Kind) Type() protoreflect.EnumType {
	return &file_cc_cc_proto_enumTypes[1]
}

func (x Kind) Number() protoreflect.EnumNumber {
	return protoreflect.EnumNumber(x)
}

// Deprecated: Use Kind.Descriptor instead.
func (Kind) EnumDescriptor() ([]byte, []int) {
	return file_cc_cc_proto_rawDescGZIP(), []int{1}
}

type Empty struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields
}

func (x *Empty) Reset() {
	*x = Empty{}
	if protoimpl.UnsafeEnabled {
		mi := &file_cc_cc_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Empty) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Empty) ProtoMessage() {}

func (x *Empty) ProtoReflect() protoreflect.Message {
	mi := &file_cc_cc_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Empty.ProtoReflect.Descriptor instead.
func (*Empty) Descriptor() ([]byte, []int) {
	return file_cc_cc_proto_rawDescGZIP(), []int{0}
}

type Chat struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Uuid     string   `protobuf:"bytes,1,opt,name=uuid,proto3" json:"uuid,omitempty"`
	Users    []string `protobuf:"bytes,2,rep,name=users,proto3" json:"users,omitempty"`
	Admins   []string `protobuf:"bytes,3,rep,name=admins,proto3" json:"admins,omitempty"`
	Owner    string   `protobuf:"bytes,4,opt,name=owner,proto3" json:"owner,omitempty"`
	Gateways []string `protobuf:"bytes,5,rep,name=gateways,proto3" json:"gateways,omitempty"`
	Role     Role     `protobuf:"varint,6,opt,name=role,proto3,enum=cc.Role" json:"role,omitempty"`
	Topic    *string  `protobuf:"bytes,7,opt,name=topic,proto3,oneof" json:"topic,omitempty"`
}

func (x *Chat) Reset() {
	*x = Chat{}
	if protoimpl.UnsafeEnabled {
		mi := &file_cc_cc_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Chat) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Chat) ProtoMessage() {}

func (x *Chat) ProtoReflect() protoreflect.Message {
	mi := &file_cc_cc_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Chat.ProtoReflect.Descriptor instead.
func (*Chat) Descriptor() ([]byte, []int) {
	return file_cc_cc_proto_rawDescGZIP(), []int{1}
}

func (x *Chat) GetUuid() string {
	if x != nil {
		return x.Uuid
	}
	return ""
}

func (x *Chat) GetUsers() []string {
	if x != nil {
		return x.Users
	}
	return nil
}

func (x *Chat) GetAdmins() []string {
	if x != nil {
		return x.Admins
	}
	return nil
}

func (x *Chat) GetOwner() string {
	if x != nil {
		return x.Owner
	}
	return ""
}

func (x *Chat) GetGateways() []string {
	if x != nil {
		return x.Gateways
	}
	return nil
}

func (x *Chat) GetRole() Role {
	if x != nil {
		return x.Role
	}
	return Role_NOACCESS
}

func (x *Chat) GetTopic() string {
	if x != nil && x.Topic != nil {
		return *x.Topic
	}
	return ""
}

type Chats struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Chats []*Chat `protobuf:"bytes,1,rep,name=chats,proto3" json:"chats,omitempty"`
}

func (x *Chats) Reset() {
	*x = Chats{}
	if protoimpl.UnsafeEnabled {
		mi := &file_cc_cc_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Chats) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Chats) ProtoMessage() {}

func (x *Chats) ProtoReflect() protoreflect.Message {
	mi := &file_cc_cc_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Chats.ProtoReflect.Descriptor instead.
func (*Chats) Descriptor() ([]byte, []int) {
	return file_cc_cc_proto_rawDescGZIP(), []int{2}
}

func (x *Chats) GetChats() []*Chat {
	if x != nil {
		return x.Chats
	}
	return nil
}

type Attachment struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Type    string `protobuf:"bytes,1,opt,name=type,proto3" json:"type,omitempty"`
	Content []byte `protobuf:"bytes,2,opt,name=content,proto3" json:"content,omitempty"`
}

func (x *Attachment) Reset() {
	*x = Attachment{}
	if protoimpl.UnsafeEnabled {
		mi := &file_cc_cc_proto_msgTypes[3]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Attachment) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Attachment) ProtoMessage() {}

func (x *Attachment) ProtoReflect() protoreflect.Message {
	mi := &file_cc_cc_proto_msgTypes[3]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Attachment.ProtoReflect.Descriptor instead.
func (*Attachment) Descriptor() ([]byte, []int) {
	return file_cc_cc_proto_rawDescGZIP(), []int{3}
}

func (x *Attachment) GetType() string {
	if x != nil {
		return x.Type
	}
	return ""
}

func (x *Attachment) GetContent() []byte {
	if x != nil {
		return x.Content
	}
	return nil
}

type Message struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Uuid        string        `protobuf:"bytes,1,opt,name=uuid,proto3" json:"uuid,omitempty"`
	Kind        Kind          `protobuf:"varint,2,opt,name=kind,proto3,enum=cc.Kind" json:"kind,omitempty"`
	Sender      string        `protobuf:"bytes,3,opt,name=sender,proto3" json:"sender,omitempty"`
	Content     string        `protobuf:"bytes,4,opt,name=content,proto3" json:"content,omitempty"`
	Attachments []*Attachment `protobuf:"bytes,5,rep,name=attachments,proto3" json:"attachments,omitempty"`
	Gateways    []string      `protobuf:"bytes,6,rep,name=gateways,proto3" json:"gateways,omitempty"`
	Chat        *string       `protobuf:"bytes,7,opt,name=chat,proto3,oneof" json:"chat,omitempty"` // Required for Send, Update and Delete requests
	Sent        int64         `protobuf:"varint,8,opt,name=sent,proto3" json:"sent,omitempty"`
	Edited      int64         `protobuf:"varint,9,opt,name=edited,proto3" json:"edited,omitempty"`
}

func (x *Message) Reset() {
	*x = Message{}
	if protoimpl.UnsafeEnabled {
		mi := &file_cc_cc_proto_msgTypes[4]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Message) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Message) ProtoMessage() {}

func (x *Message) ProtoReflect() protoreflect.Message {
	mi := &file_cc_cc_proto_msgTypes[4]
	if protoimpl.UnsafeEnabled && x != nil {
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
	return file_cc_cc_proto_rawDescGZIP(), []int{4}
}

func (x *Message) GetUuid() string {
	if x != nil {
		return x.Uuid
	}
	return ""
}

func (x *Message) GetKind() Kind {
	if x != nil {
		return x.Kind
	}
	return Kind_DEFAULT
}

func (x *Message) GetSender() string {
	if x != nil {
		return x.Sender
	}
	return ""
}

func (x *Message) GetContent() string {
	if x != nil {
		return x.Content
	}
	return ""
}

func (x *Message) GetAttachments() []*Attachment {
	if x != nil {
		return x.Attachments
	}
	return nil
}

func (x *Message) GetGateways() []string {
	if x != nil {
		return x.Gateways
	}
	return nil
}

func (x *Message) GetChat() string {
	if x != nil && x.Chat != nil {
		return *x.Chat
	}
	return ""
}

func (x *Message) GetSent() int64 {
	if x != nil {
		return x.Sent
	}
	return 0
}

func (x *Message) GetEdited() int64 {
	if x != nil {
		return x.Edited
	}
	return 0
}

type Messages struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Messages []*Message `protobuf:"bytes,1,rep,name=messages,proto3" json:"messages,omitempty"`
}

func (x *Messages) Reset() {
	*x = Messages{}
	if protoimpl.UnsafeEnabled {
		mi := &file_cc_cc_proto_msgTypes[5]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Messages) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Messages) ProtoMessage() {}

func (x *Messages) ProtoReflect() protoreflect.Message {
	mi := &file_cc_cc_proto_msgTypes[5]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Messages.ProtoReflect.Descriptor instead.
func (*Messages) Descriptor() ([]byte, []int) {
	return file_cc_cc_proto_rawDescGZIP(), []int{5}
}

func (x *Messages) GetMessages() []*Message {
	if x != nil {
		return x.Messages
	}
	return nil
}

type User struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Uuid  string `protobuf:"bytes,1,opt,name=uuid,proto3" json:"uuid,omitempty"`
	Title string `protobuf:"bytes,2,opt,name=title,proto3" json:"title,omitempty"`
}

func (x *User) Reset() {
	*x = User{}
	if protoimpl.UnsafeEnabled {
		mi := &file_cc_cc_proto_msgTypes[6]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *User) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*User) ProtoMessage() {}

func (x *User) ProtoReflect() protoreflect.Message {
	mi := &file_cc_cc_proto_msgTypes[6]
	if protoimpl.UnsafeEnabled && x != nil {
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
	return file_cc_cc_proto_rawDescGZIP(), []int{6}
}

func (x *User) GetUuid() string {
	if x != nil {
		return x.Uuid
	}
	return ""
}

func (x *User) GetTitle() string {
	if x != nil {
		return x.Title
	}
	return ""
}

type Defaults struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Gateways []string `protobuf:"bytes,1,rep,name=gateways,proto3" json:"gateways,omitempty"`
	Admins   []string `protobuf:"bytes,2,rep,name=admins,proto3" json:"admins,omitempty"`
}

func (x *Defaults) Reset() {
	*x = Defaults{}
	if protoimpl.UnsafeEnabled {
		mi := &file_cc_cc_proto_msgTypes[7]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Defaults) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Defaults) ProtoMessage() {}

func (x *Defaults) ProtoReflect() protoreflect.Message {
	mi := &file_cc_cc_proto_msgTypes[7]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Defaults.ProtoReflect.Descriptor instead.
func (*Defaults) Descriptor() ([]byte, []int) {
	return file_cc_cc_proto_rawDescGZIP(), []int{7}
}

func (x *Defaults) GetGateways() []string {
	if x != nil {
		return x.Gateways
	}
	return nil
}

func (x *Defaults) GetAdmins() []string {
	if x != nil {
		return x.Admins
	}
	return nil
}

type Users struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Users []*User `protobuf:"bytes,1,rep,name=users,proto3" json:"users,omitempty"`
}

func (x *Users) Reset() {
	*x = Users{}
	if protoimpl.UnsafeEnabled {
		mi := &file_cc_cc_proto_msgTypes[8]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Users) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Users) ProtoMessage() {}

func (x *Users) ProtoReflect() protoreflect.Message {
	mi := &file_cc_cc_proto_msgTypes[8]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Users.ProtoReflect.Descriptor instead.
func (*Users) Descriptor() ([]byte, []int) {
	return file_cc_cc_proto_rawDescGZIP(), []int{8}
}

func (x *Users) GetUsers() []*User {
	if x != nil {
		return x.Users
	}
	return nil
}

var File_cc_cc_proto protoreflect.FileDescriptor

var file_cc_cc_proto_rawDesc = []byte{
	0x0a, 0x0b, 0x63, 0x63, 0x2f, 0x63, 0x63, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x02, 0x63,
	0x63, 0x22, 0x07, 0x0a, 0x05, 0x45, 0x6d, 0x70, 0x74, 0x79, 0x22, 0xbd, 0x01, 0x0a, 0x04, 0x43,
	0x68, 0x61, 0x74, 0x12, 0x12, 0x0a, 0x04, 0x75, 0x75, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28,
	0x09, 0x52, 0x04, 0x75, 0x75, 0x69, 0x64, 0x12, 0x14, 0x0a, 0x05, 0x75, 0x73, 0x65, 0x72, 0x73,
	0x18, 0x02, 0x20, 0x03, 0x28, 0x09, 0x52, 0x05, 0x75, 0x73, 0x65, 0x72, 0x73, 0x12, 0x16, 0x0a,
	0x06, 0x61, 0x64, 0x6d, 0x69, 0x6e, 0x73, 0x18, 0x03, 0x20, 0x03, 0x28, 0x09, 0x52, 0x06, 0x61,
	0x64, 0x6d, 0x69, 0x6e, 0x73, 0x12, 0x14, 0x0a, 0x05, 0x6f, 0x77, 0x6e, 0x65, 0x72, 0x18, 0x04,
	0x20, 0x01, 0x28, 0x09, 0x52, 0x05, 0x6f, 0x77, 0x6e, 0x65, 0x72, 0x12, 0x1a, 0x0a, 0x08, 0x67,
	0x61, 0x74, 0x65, 0x77, 0x61, 0x79, 0x73, 0x18, 0x05, 0x20, 0x03, 0x28, 0x09, 0x52, 0x08, 0x67,
	0x61, 0x74, 0x65, 0x77, 0x61, 0x79, 0x73, 0x12, 0x1c, 0x0a, 0x04, 0x72, 0x6f, 0x6c, 0x65, 0x18,
	0x06, 0x20, 0x01, 0x28, 0x0e, 0x32, 0x08, 0x2e, 0x63, 0x63, 0x2e, 0x52, 0x6f, 0x6c, 0x65, 0x52,
	0x04, 0x72, 0x6f, 0x6c, 0x65, 0x12, 0x19, 0x0a, 0x05, 0x74, 0x6f, 0x70, 0x69, 0x63, 0x18, 0x07,
	0x20, 0x01, 0x28, 0x09, 0x48, 0x00, 0x52, 0x05, 0x74, 0x6f, 0x70, 0x69, 0x63, 0x88, 0x01, 0x01,
	0x42, 0x08, 0x0a, 0x06, 0x5f, 0x74, 0x6f, 0x70, 0x69, 0x63, 0x22, 0x27, 0x0a, 0x05, 0x43, 0x68,
	0x61, 0x74, 0x73, 0x12, 0x1e, 0x0a, 0x05, 0x63, 0x68, 0x61, 0x74, 0x73, 0x18, 0x01, 0x20, 0x03,
	0x28, 0x0b, 0x32, 0x08, 0x2e, 0x63, 0x63, 0x2e, 0x43, 0x68, 0x61, 0x74, 0x52, 0x05, 0x63, 0x68,
	0x61, 0x74, 0x73, 0x22, 0x3a, 0x0a, 0x0a, 0x41, 0x74, 0x74, 0x61, 0x63, 0x68, 0x6d, 0x65, 0x6e,
	0x74, 0x12, 0x12, 0x0a, 0x04, 0x74, 0x79, 0x70, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52,
	0x04, 0x74, 0x79, 0x70, 0x65, 0x12, 0x18, 0x0a, 0x07, 0x63, 0x6f, 0x6e, 0x74, 0x65, 0x6e, 0x74,
	0x18, 0x02, 0x20, 0x01, 0x28, 0x0c, 0x52, 0x07, 0x63, 0x6f, 0x6e, 0x74, 0x65, 0x6e, 0x74, 0x22,
	0x89, 0x02, 0x0a, 0x07, 0x4d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x12, 0x12, 0x0a, 0x04, 0x75,
	0x75, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x04, 0x75, 0x75, 0x69, 0x64, 0x12,
	0x1c, 0x0a, 0x04, 0x6b, 0x69, 0x6e, 0x64, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0e, 0x32, 0x08, 0x2e,
	0x63, 0x63, 0x2e, 0x4b, 0x69, 0x6e, 0x64, 0x52, 0x04, 0x6b, 0x69, 0x6e, 0x64, 0x12, 0x16, 0x0a,
	0x06, 0x73, 0x65, 0x6e, 0x64, 0x65, 0x72, 0x18, 0x03, 0x20, 0x01, 0x28, 0x09, 0x52, 0x06, 0x73,
	0x65, 0x6e, 0x64, 0x65, 0x72, 0x12, 0x18, 0x0a, 0x07, 0x63, 0x6f, 0x6e, 0x74, 0x65, 0x6e, 0x74,
	0x18, 0x04, 0x20, 0x01, 0x28, 0x09, 0x52, 0x07, 0x63, 0x6f, 0x6e, 0x74, 0x65, 0x6e, 0x74, 0x12,
	0x30, 0x0a, 0x0b, 0x61, 0x74, 0x74, 0x61, 0x63, 0x68, 0x6d, 0x65, 0x6e, 0x74, 0x73, 0x18, 0x05,
	0x20, 0x03, 0x28, 0x0b, 0x32, 0x0e, 0x2e, 0x63, 0x63, 0x2e, 0x41, 0x74, 0x74, 0x61, 0x63, 0x68,
	0x6d, 0x65, 0x6e, 0x74, 0x52, 0x0b, 0x61, 0x74, 0x74, 0x61, 0x63, 0x68, 0x6d, 0x65, 0x6e, 0x74,
	0x73, 0x12, 0x1a, 0x0a, 0x08, 0x67, 0x61, 0x74, 0x65, 0x77, 0x61, 0x79, 0x73, 0x18, 0x06, 0x20,
	0x03, 0x28, 0x09, 0x52, 0x08, 0x67, 0x61, 0x74, 0x65, 0x77, 0x61, 0x79, 0x73, 0x12, 0x17, 0x0a,
	0x04, 0x63, 0x68, 0x61, 0x74, 0x18, 0x07, 0x20, 0x01, 0x28, 0x09, 0x48, 0x00, 0x52, 0x04, 0x63,
	0x68, 0x61, 0x74, 0x88, 0x01, 0x01, 0x12, 0x12, 0x0a, 0x04, 0x73, 0x65, 0x6e, 0x74, 0x18, 0x08,
	0x20, 0x01, 0x28, 0x03, 0x52, 0x04, 0x73, 0x65, 0x6e, 0x74, 0x12, 0x16, 0x0a, 0x06, 0x65, 0x64,
	0x69, 0x74, 0x65, 0x64, 0x18, 0x09, 0x20, 0x01, 0x28, 0x03, 0x52, 0x06, 0x65, 0x64, 0x69, 0x74,
	0x65, 0x64, 0x42, 0x07, 0x0a, 0x05, 0x5f, 0x63, 0x68, 0x61, 0x74, 0x22, 0x33, 0x0a, 0x08, 0x4d,
	0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x73, 0x12, 0x27, 0x0a, 0x08, 0x6d, 0x65, 0x73, 0x73, 0x61,
	0x67, 0x65, 0x73, 0x18, 0x01, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x0b, 0x2e, 0x63, 0x63, 0x2e, 0x4d,
	0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x52, 0x08, 0x6d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x73,
	0x22, 0x30, 0x0a, 0x04, 0x55, 0x73, 0x65, 0x72, 0x12, 0x12, 0x0a, 0x04, 0x75, 0x75, 0x69, 0x64,
	0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x04, 0x75, 0x75, 0x69, 0x64, 0x12, 0x14, 0x0a, 0x05,
	0x74, 0x69, 0x74, 0x6c, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x05, 0x74, 0x69, 0x74,
	0x6c, 0x65, 0x22, 0x3e, 0x0a, 0x08, 0x44, 0x65, 0x66, 0x61, 0x75, 0x6c, 0x74, 0x73, 0x12, 0x1a,
	0x0a, 0x08, 0x67, 0x61, 0x74, 0x65, 0x77, 0x61, 0x79, 0x73, 0x18, 0x01, 0x20, 0x03, 0x28, 0x09,
	0x52, 0x08, 0x67, 0x61, 0x74, 0x65, 0x77, 0x61, 0x79, 0x73, 0x12, 0x16, 0x0a, 0x06, 0x61, 0x64,
	0x6d, 0x69, 0x6e, 0x73, 0x18, 0x02, 0x20, 0x03, 0x28, 0x09, 0x52, 0x06, 0x61, 0x64, 0x6d, 0x69,
	0x6e, 0x73, 0x22, 0x27, 0x0a, 0x05, 0x55, 0x73, 0x65, 0x72, 0x73, 0x12, 0x1e, 0x0a, 0x05, 0x75,
	0x73, 0x65, 0x72, 0x73, 0x18, 0x01, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x08, 0x2e, 0x63, 0x63, 0x2e,
	0x55, 0x73, 0x65, 0x72, 0x52, 0x05, 0x75, 0x73, 0x65, 0x72, 0x73, 0x2a, 0x34, 0x0a, 0x04, 0x52,
	0x6f, 0x6c, 0x65, 0x12, 0x0c, 0x0a, 0x08, 0x4e, 0x4f, 0x41, 0x43, 0x43, 0x45, 0x53, 0x53, 0x10,
	0x00, 0x12, 0x08, 0x0a, 0x04, 0x55, 0x53, 0x45, 0x52, 0x10, 0x01, 0x12, 0x09, 0x0a, 0x05, 0x4f,
	0x57, 0x4e, 0x45, 0x52, 0x10, 0x02, 0x12, 0x09, 0x0a, 0x05, 0x41, 0x44, 0x4d, 0x49, 0x4e, 0x10,
	0x03, 0x2a, 0x23, 0x0a, 0x04, 0x4b, 0x69, 0x6e, 0x64, 0x12, 0x0b, 0x0a, 0x07, 0x44, 0x45, 0x46,
	0x41, 0x55, 0x4c, 0x54, 0x10, 0x00, 0x12, 0x0e, 0x0a, 0x0a, 0x41, 0x44, 0x4d, 0x49, 0x4e, 0x5f,
	0x4f, 0x4e, 0x4c, 0x59, 0x10, 0x01, 0x32, 0x82, 0x01, 0x0a, 0x08, 0x43, 0x68, 0x61, 0x74, 0x73,
	0x41, 0x50, 0x49, 0x12, 0x1c, 0x0a, 0x06, 0x43, 0x72, 0x65, 0x61, 0x74, 0x65, 0x12, 0x08, 0x2e,
	0x63, 0x63, 0x2e, 0x43, 0x68, 0x61, 0x74, 0x1a, 0x08, 0x2e, 0x63, 0x63, 0x2e, 0x43, 0x68, 0x61,
	0x74, 0x12, 0x1c, 0x0a, 0x06, 0x55, 0x70, 0x64, 0x61, 0x74, 0x65, 0x12, 0x08, 0x2e, 0x63, 0x63,
	0x2e, 0x43, 0x68, 0x61, 0x74, 0x1a, 0x08, 0x2e, 0x63, 0x63, 0x2e, 0x43, 0x68, 0x61, 0x74, 0x12,
	0x1c, 0x0a, 0x04, 0x4c, 0x69, 0x73, 0x74, 0x12, 0x09, 0x2e, 0x63, 0x63, 0x2e, 0x45, 0x6d, 0x70,
	0x74, 0x79, 0x1a, 0x09, 0x2e, 0x63, 0x63, 0x2e, 0x43, 0x68, 0x61, 0x74, 0x73, 0x12, 0x1c, 0x0a,
	0x06, 0x44, 0x65, 0x6c, 0x65, 0x74, 0x65, 0x12, 0x08, 0x2e, 0x63, 0x63, 0x2e, 0x43, 0x68, 0x61,
	0x74, 0x1a, 0x08, 0x2e, 0x63, 0x63, 0x2e, 0x43, 0x68, 0x61, 0x74, 0x32, 0x96, 0x01, 0x0a, 0x0b,
	0x4d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x73, 0x41, 0x50, 0x49, 0x12, 0x1d, 0x0a, 0x03, 0x47,
	0x65, 0x74, 0x12, 0x08, 0x2e, 0x63, 0x63, 0x2e, 0x43, 0x68, 0x61, 0x74, 0x1a, 0x0c, 0x2e, 0x63,
	0x63, 0x2e, 0x4d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x73, 0x12, 0x20, 0x0a, 0x04, 0x53, 0x65,
	0x6e, 0x64, 0x12, 0x0b, 0x2e, 0x63, 0x63, 0x2e, 0x4d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x1a,
	0x0b, 0x2e, 0x63, 0x63, 0x2e, 0x4d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x12, 0x22, 0x0a, 0x06,
	0x55, 0x70, 0x64, 0x61, 0x74, 0x65, 0x12, 0x0b, 0x2e, 0x63, 0x63, 0x2e, 0x4d, 0x65, 0x73, 0x73,
	0x61, 0x67, 0x65, 0x1a, 0x0b, 0x2e, 0x63, 0x63, 0x2e, 0x4d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65,
	0x12, 0x22, 0x0a, 0x06, 0x44, 0x65, 0x6c, 0x65, 0x74, 0x65, 0x12, 0x0b, 0x2e, 0x63, 0x63, 0x2e,
	0x4d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x1a, 0x0b, 0x2e, 0x63, 0x63, 0x2e, 0x4d, 0x65, 0x73,
	0x73, 0x61, 0x67, 0x65, 0x32, 0x70, 0x0a, 0x08, 0x55, 0x73, 0x65, 0x72, 0x73, 0x41, 0x50, 0x49,
	0x12, 0x19, 0x0a, 0x02, 0x4d, 0x65, 0x12, 0x09, 0x2e, 0x63, 0x63, 0x2e, 0x45, 0x6d, 0x70, 0x74,
	0x79, 0x1a, 0x08, 0x2e, 0x63, 0x63, 0x2e, 0x55, 0x73, 0x65, 0x72, 0x12, 0x28, 0x0a, 0x0d, 0x46,
	0x65, 0x74, 0x63, 0x68, 0x44, 0x65, 0x66, 0x61, 0x75, 0x6c, 0x74, 0x73, 0x12, 0x09, 0x2e, 0x63,
	0x63, 0x2e, 0x45, 0x6d, 0x70, 0x74, 0x79, 0x1a, 0x0c, 0x2e, 0x63, 0x63, 0x2e, 0x44, 0x65, 0x66,
	0x61, 0x75, 0x6c, 0x74, 0x73, 0x12, 0x1f, 0x0a, 0x07, 0x52, 0x65, 0x73, 0x6f, 0x6c, 0x76, 0x65,
	0x12, 0x09, 0x2e, 0x63, 0x63, 0x2e, 0x55, 0x73, 0x65, 0x72, 0x73, 0x1a, 0x09, 0x2e, 0x63, 0x63,
	0x2e, 0x55, 0x73, 0x65, 0x72, 0x73, 0x42, 0x25, 0x5a, 0x23, 0x67, 0x69, 0x74, 0x68, 0x75, 0x62,
	0x2e, 0x63, 0x6f, 0x6d, 0x2f, 0x73, 0x6c, 0x6e, 0x74, 0x6f, 0x70, 0x70, 0x2f, 0x63, 0x6f, 0x72,
	0x65, 0x2d, 0x63, 0x68, 0x61, 0x74, 0x74, 0x69, 0x6e, 0x67, 0x2f, 0x63, 0x63, 0x62, 0x06, 0x70,
	0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_cc_cc_proto_rawDescOnce sync.Once
	file_cc_cc_proto_rawDescData = file_cc_cc_proto_rawDesc
)

func file_cc_cc_proto_rawDescGZIP() []byte {
	file_cc_cc_proto_rawDescOnce.Do(func() {
		file_cc_cc_proto_rawDescData = protoimpl.X.CompressGZIP(file_cc_cc_proto_rawDescData)
	})
	return file_cc_cc_proto_rawDescData
}

var file_cc_cc_proto_enumTypes = make([]protoimpl.EnumInfo, 2)
var file_cc_cc_proto_msgTypes = make([]protoimpl.MessageInfo, 9)
var file_cc_cc_proto_goTypes = []interface{}{
	(Role)(0),          // 0: cc.Role
	(Kind)(0),          // 1: cc.Kind
	(*Empty)(nil),      // 2: cc.Empty
	(*Chat)(nil),       // 3: cc.Chat
	(*Chats)(nil),      // 4: cc.Chats
	(*Attachment)(nil), // 5: cc.Attachment
	(*Message)(nil),    // 6: cc.Message
	(*Messages)(nil),   // 7: cc.Messages
	(*User)(nil),       // 8: cc.User
	(*Defaults)(nil),   // 9: cc.Defaults
	(*Users)(nil),      // 10: cc.Users
}
var file_cc_cc_proto_depIdxs = []int32{
	0,  // 0: cc.Chat.role:type_name -> cc.Role
	3,  // 1: cc.Chats.chats:type_name -> cc.Chat
	1,  // 2: cc.Message.kind:type_name -> cc.Kind
	5,  // 3: cc.Message.attachments:type_name -> cc.Attachment
	6,  // 4: cc.Messages.messages:type_name -> cc.Message
	8,  // 5: cc.Users.users:type_name -> cc.User
	3,  // 6: cc.ChatsAPI.Create:input_type -> cc.Chat
	3,  // 7: cc.ChatsAPI.Update:input_type -> cc.Chat
	2,  // 8: cc.ChatsAPI.List:input_type -> cc.Empty
	3,  // 9: cc.ChatsAPI.Delete:input_type -> cc.Chat
	3,  // 10: cc.MessagesAPI.Get:input_type -> cc.Chat
	6,  // 11: cc.MessagesAPI.Send:input_type -> cc.Message
	6,  // 12: cc.MessagesAPI.Update:input_type -> cc.Message
	6,  // 13: cc.MessagesAPI.Delete:input_type -> cc.Message
	2,  // 14: cc.UsersAPI.Me:input_type -> cc.Empty
	2,  // 15: cc.UsersAPI.FetchDefaults:input_type -> cc.Empty
	10, // 16: cc.UsersAPI.Resolve:input_type -> cc.Users
	3,  // 17: cc.ChatsAPI.Create:output_type -> cc.Chat
	3,  // 18: cc.ChatsAPI.Update:output_type -> cc.Chat
	4,  // 19: cc.ChatsAPI.List:output_type -> cc.Chats
	3,  // 20: cc.ChatsAPI.Delete:output_type -> cc.Chat
	7,  // 21: cc.MessagesAPI.Get:output_type -> cc.Messages
	6,  // 22: cc.MessagesAPI.Send:output_type -> cc.Message
	6,  // 23: cc.MessagesAPI.Update:output_type -> cc.Message
	6,  // 24: cc.MessagesAPI.Delete:output_type -> cc.Message
	8,  // 25: cc.UsersAPI.Me:output_type -> cc.User
	9,  // 26: cc.UsersAPI.FetchDefaults:output_type -> cc.Defaults
	10, // 27: cc.UsersAPI.Resolve:output_type -> cc.Users
	17, // [17:28] is the sub-list for method output_type
	6,  // [6:17] is the sub-list for method input_type
	6,  // [6:6] is the sub-list for extension type_name
	6,  // [6:6] is the sub-list for extension extendee
	0,  // [0:6] is the sub-list for field type_name
}

func init() { file_cc_cc_proto_init() }
func file_cc_cc_proto_init() {
	if File_cc_cc_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_cc_cc_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Empty); i {
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
		file_cc_cc_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Chat); i {
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
		file_cc_cc_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Chats); i {
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
		file_cc_cc_proto_msgTypes[3].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Attachment); i {
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
		file_cc_cc_proto_msgTypes[4].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Message); i {
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
		file_cc_cc_proto_msgTypes[5].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Messages); i {
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
		file_cc_cc_proto_msgTypes[6].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*User); i {
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
		file_cc_cc_proto_msgTypes[7].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Defaults); i {
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
		file_cc_cc_proto_msgTypes[8].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Users); i {
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
	file_cc_cc_proto_msgTypes[1].OneofWrappers = []interface{}{}
	file_cc_cc_proto_msgTypes[4].OneofWrappers = []interface{}{}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_cc_cc_proto_rawDesc,
			NumEnums:      2,
			NumMessages:   9,
			NumExtensions: 0,
			NumServices:   3,
		},
		GoTypes:           file_cc_cc_proto_goTypes,
		DependencyIndexes: file_cc_cc_proto_depIdxs,
		EnumInfos:         file_cc_cc_proto_enumTypes,
		MessageInfos:      file_cc_cc_proto_msgTypes,
	}.Build()
	File_cc_cc_proto = out.File
	file_cc_cc_proto_rawDesc = nil
	file_cc_cc_proto_goTypes = nil
	file_cc_cc_proto_depIdxs = nil
}
