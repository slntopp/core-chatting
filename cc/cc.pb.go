// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.28.1
// 	protoc        (unknown)
// source: cc/cc.proto

package cc

import (
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
	_ "google.golang.org/protobuf/types/known/structpb"
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

type EventType int32

const (
	EventType_PING            EventType = 0
	EventType_CHAT_CREATED    EventType = 1
	EventType_CHAT_UPDATED    EventType = 2
	EventType_CHAT_DELETED    EventType = 3
	EventType_MESSAGE_SEND    EventType = 4
	EventType_MESSAGE_UPDATED EventType = 5
	EventType_MESSAGE_DELETED EventType = 6
)

// Enum value maps for EventType.
var (
	EventType_name = map[int32]string{
		0: "PING",
		1: "CHAT_CREATED",
		2: "CHAT_UPDATED",
		3: "CHAT_DELETED",
		4: "MESSAGE_SEND",
		5: "MESSAGE_UPDATED",
		6: "MESSAGE_DELETED",
	}
	EventType_value = map[string]int32{
		"PING":            0,
		"CHAT_CREATED":    1,
		"CHAT_UPDATED":    2,
		"CHAT_DELETED":    3,
		"MESSAGE_SEND":    4,
		"MESSAGE_UPDATED": 5,
		"MESSAGE_DELETED": 6,
	}
)

func (x EventType) Enum() *EventType {
	p := new(EventType)
	*p = x
	return p
}

func (x EventType) String() string {
	return protoimpl.X.EnumStringOf(x.Descriptor(), protoreflect.EnumNumber(x))
}

func (EventType) Descriptor() protoreflect.EnumDescriptor {
	return file_cc_cc_proto_enumTypes[2].Descriptor()
}

func (EventType) Type() protoreflect.EnumType {
	return &file_cc_cc_proto_enumTypes[2]
}

func (x EventType) Number() protoreflect.EnumNumber {
	return protoreflect.EnumNumber(x)
}

// Deprecated: Use EventType.Descriptor instead.
func (EventType) EnumDescriptor() ([]byte, []int) {
	return file_cc_cc_proto_rawDescGZIP(), []int{2}
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

type ChatMeta struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Unread      int32    `protobuf:"varint,1,opt,name=unread,proto3" json:"unread,omitempty"`
	LastMessage *Message `protobuf:"bytes,2,opt,name=last_message,json=lastMessage,proto3" json:"last_message,omitempty"`
}

func (x *ChatMeta) Reset() {
	*x = ChatMeta{}
	if protoimpl.UnsafeEnabled {
		mi := &file_cc_cc_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ChatMeta) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ChatMeta) ProtoMessage() {}

func (x *ChatMeta) ProtoReflect() protoreflect.Message {
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

// Deprecated: Use ChatMeta.ProtoReflect.Descriptor instead.
func (*ChatMeta) Descriptor() ([]byte, []int) {
	return file_cc_cc_proto_rawDescGZIP(), []int{1}
}

func (x *ChatMeta) GetUnread() int32 {
	if x != nil {
		return x.Unread
	}
	return 0
}

func (x *ChatMeta) GetLastMessage() *Message {
	if x != nil {
		return x.LastMessage
	}
	return nil
}

type Chat struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Uuid     string    `protobuf:"bytes,1,opt,name=uuid,proto3" json:"uuid,omitempty"`
	Users    []string  `protobuf:"bytes,2,rep,name=users,proto3" json:"users,omitempty"`
	Admins   []string  `protobuf:"bytes,3,rep,name=admins,proto3" json:"admins,omitempty"`
	Owner    string    `protobuf:"bytes,4,opt,name=owner,proto3" json:"owner,omitempty"`
	Gateways []string  `protobuf:"bytes,5,rep,name=gateways,proto3" json:"gateways,omitempty"`
	Role     Role      `protobuf:"varint,6,opt,name=role,proto3,enum=cc.Role" json:"role,omitempty"`
	Topic    *string   `protobuf:"bytes,7,opt,name=topic,proto3,oneof" json:"topic,omitempty"`
	Meta     *ChatMeta `protobuf:"bytes,8,opt,name=meta,proto3,oneof" json:"meta,omitempty"`
	Created  int64     `protobuf:"varint,9,opt,name=created,proto3" json:"created,omitempty"`
}

func (x *Chat) Reset() {
	*x = Chat{}
	if protoimpl.UnsafeEnabled {
		mi := &file_cc_cc_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Chat) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Chat) ProtoMessage() {}

func (x *Chat) ProtoReflect() protoreflect.Message {
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

// Deprecated: Use Chat.ProtoReflect.Descriptor instead.
func (*Chat) Descriptor() ([]byte, []int) {
	return file_cc_cc_proto_rawDescGZIP(), []int{2}
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

func (x *Chat) GetMeta() *ChatMeta {
	if x != nil {
		return x.Meta
	}
	return nil
}

func (x *Chat) GetCreated() int64 {
	if x != nil {
		return x.Created
	}
	return 0
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
		mi := &file_cc_cc_proto_msgTypes[3]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Chats) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Chats) ProtoMessage() {}

func (x *Chats) ProtoReflect() protoreflect.Message {
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

// Deprecated: Use Chats.ProtoReflect.Descriptor instead.
func (*Chats) Descriptor() ([]byte, []int) {
	return file_cc_cc_proto_rawDescGZIP(), []int{3}
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
		mi := &file_cc_cc_proto_msgTypes[4]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Attachment) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Attachment) ProtoMessage() {}

func (x *Attachment) ProtoReflect() protoreflect.Message {
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

// Deprecated: Use Attachment.ProtoReflect.Descriptor instead.
func (*Attachment) Descriptor() ([]byte, []int) {
	return file_cc_cc_proto_rawDescGZIP(), []int{4}
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
	UnderReview bool          `protobuf:"varint,10,opt,name=under_review,json=underReview,proto3" json:"under_review,omitempty"`
}

func (x *Message) Reset() {
	*x = Message{}
	if protoimpl.UnsafeEnabled {
		mi := &file_cc_cc_proto_msgTypes[5]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Message) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Message) ProtoMessage() {}

func (x *Message) ProtoReflect() protoreflect.Message {
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

// Deprecated: Use Message.ProtoReflect.Descriptor instead.
func (*Message) Descriptor() ([]byte, []int) {
	return file_cc_cc_proto_rawDescGZIP(), []int{5}
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

func (x *Message) GetUnderReview() bool {
	if x != nil {
		return x.UnderReview
	}
	return false
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
		mi := &file_cc_cc_proto_msgTypes[6]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Messages) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Messages) ProtoMessage() {}

func (x *Messages) ProtoReflect() protoreflect.Message {
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

// Deprecated: Use Messages.ProtoReflect.Descriptor instead.
func (*Messages) Descriptor() ([]byte, []int) {
	return file_cc_cc_proto_rawDescGZIP(), []int{6}
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
		mi := &file_cc_cc_proto_msgTypes[7]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *User) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*User) ProtoMessage() {}

func (x *User) ProtoReflect() protoreflect.Message {
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

// Deprecated: Use User.ProtoReflect.Descriptor instead.
func (*User) Descriptor() ([]byte, []int) {
	return file_cc_cc_proto_rawDescGZIP(), []int{7}
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
		mi := &file_cc_cc_proto_msgTypes[8]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Defaults) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Defaults) ProtoMessage() {}

func (x *Defaults) ProtoReflect() protoreflect.Message {
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

// Deprecated: Use Defaults.ProtoReflect.Descriptor instead.
func (*Defaults) Descriptor() ([]byte, []int) {
	return file_cc_cc_proto_rawDescGZIP(), []int{8}
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
		mi := &file_cc_cc_proto_msgTypes[9]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Users) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Users) ProtoMessage() {}

func (x *Users) ProtoReflect() protoreflect.Message {
	mi := &file_cc_cc_proto_msgTypes[9]
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
	return file_cc_cc_proto_rawDescGZIP(), []int{9}
}

func (x *Users) GetUsers() []*User {
	if x != nil {
		return x.Users
	}
	return nil
}

type Event struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Type EventType `protobuf:"varint,1,opt,name=type,proto3,enum=cc.EventType" json:"type,omitempty"`
	// Types that are assignable to Item:
	//
	//	*Event_Chat
	//	*Event_Msg
	Item isEvent_Item `protobuf_oneof:"item"`
}

func (x *Event) Reset() {
	*x = Event{}
	if protoimpl.UnsafeEnabled {
		mi := &file_cc_cc_proto_msgTypes[10]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Event) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Event) ProtoMessage() {}

func (x *Event) ProtoReflect() protoreflect.Message {
	mi := &file_cc_cc_proto_msgTypes[10]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Event.ProtoReflect.Descriptor instead.
func (*Event) Descriptor() ([]byte, []int) {
	return file_cc_cc_proto_rawDescGZIP(), []int{10}
}

func (x *Event) GetType() EventType {
	if x != nil {
		return x.Type
	}
	return EventType_PING
}

func (m *Event) GetItem() isEvent_Item {
	if m != nil {
		return m.Item
	}
	return nil
}

func (x *Event) GetChat() *Chat {
	if x, ok := x.GetItem().(*Event_Chat); ok {
		return x.Chat
	}
	return nil
}

func (x *Event) GetMsg() *Message {
	if x, ok := x.GetItem().(*Event_Msg); ok {
		return x.Msg
	}
	return nil
}

type isEvent_Item interface {
	isEvent_Item()
}

type Event_Chat struct {
	Chat *Chat `protobuf:"bytes,2,opt,name=chat,proto3,oneof"`
}

type Event_Msg struct {
	Msg *Message `protobuf:"bytes,3,opt,name=msg,proto3,oneof"`
}

func (*Event_Chat) isEvent_Item() {}

func (*Event_Msg) isEvent_Item() {}

var File_cc_cc_proto protoreflect.FileDescriptor

var file_cc_cc_proto_rawDesc = []byte{
	0x0a, 0x0b, 0x63, 0x63, 0x2f, 0x63, 0x63, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x02, 0x63,
	0x63, 0x1a, 0x1c, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62,
	0x75, 0x66, 0x2f, 0x73, 0x74, 0x72, 0x75, 0x63, 0x74, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x22,
	0x07, 0x0a, 0x05, 0x45, 0x6d, 0x70, 0x74, 0x79, 0x22, 0x52, 0x0a, 0x08, 0x43, 0x68, 0x61, 0x74,
	0x4d, 0x65, 0x74, 0x61, 0x12, 0x16, 0x0a, 0x06, 0x75, 0x6e, 0x72, 0x65, 0x61, 0x64, 0x18, 0x01,
	0x20, 0x01, 0x28, 0x05, 0x52, 0x06, 0x75, 0x6e, 0x72, 0x65, 0x61, 0x64, 0x12, 0x2e, 0x0a, 0x0c,
	0x6c, 0x61, 0x73, 0x74, 0x5f, 0x6d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x18, 0x02, 0x20, 0x01,
	0x28, 0x0b, 0x32, 0x0b, 0x2e, 0x63, 0x63, 0x2e, 0x4d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x52,
	0x0b, 0x6c, 0x61, 0x73, 0x74, 0x4d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x22, 0x87, 0x02, 0x0a,
	0x04, 0x43, 0x68, 0x61, 0x74, 0x12, 0x12, 0x0a, 0x04, 0x75, 0x75, 0x69, 0x64, 0x18, 0x01, 0x20,
	0x01, 0x28, 0x09, 0x52, 0x04, 0x75, 0x75, 0x69, 0x64, 0x12, 0x14, 0x0a, 0x05, 0x75, 0x73, 0x65,
	0x72, 0x73, 0x18, 0x02, 0x20, 0x03, 0x28, 0x09, 0x52, 0x05, 0x75, 0x73, 0x65, 0x72, 0x73, 0x12,
	0x16, 0x0a, 0x06, 0x61, 0x64, 0x6d, 0x69, 0x6e, 0x73, 0x18, 0x03, 0x20, 0x03, 0x28, 0x09, 0x52,
	0x06, 0x61, 0x64, 0x6d, 0x69, 0x6e, 0x73, 0x12, 0x14, 0x0a, 0x05, 0x6f, 0x77, 0x6e, 0x65, 0x72,
	0x18, 0x04, 0x20, 0x01, 0x28, 0x09, 0x52, 0x05, 0x6f, 0x77, 0x6e, 0x65, 0x72, 0x12, 0x1a, 0x0a,
	0x08, 0x67, 0x61, 0x74, 0x65, 0x77, 0x61, 0x79, 0x73, 0x18, 0x05, 0x20, 0x03, 0x28, 0x09, 0x52,
	0x08, 0x67, 0x61, 0x74, 0x65, 0x77, 0x61, 0x79, 0x73, 0x12, 0x1c, 0x0a, 0x04, 0x72, 0x6f, 0x6c,
	0x65, 0x18, 0x06, 0x20, 0x01, 0x28, 0x0e, 0x32, 0x08, 0x2e, 0x63, 0x63, 0x2e, 0x52, 0x6f, 0x6c,
	0x65, 0x52, 0x04, 0x72, 0x6f, 0x6c, 0x65, 0x12, 0x19, 0x0a, 0x05, 0x74, 0x6f, 0x70, 0x69, 0x63,
	0x18, 0x07, 0x20, 0x01, 0x28, 0x09, 0x48, 0x00, 0x52, 0x05, 0x74, 0x6f, 0x70, 0x69, 0x63, 0x88,
	0x01, 0x01, 0x12, 0x25, 0x0a, 0x04, 0x6d, 0x65, 0x74, 0x61, 0x18, 0x08, 0x20, 0x01, 0x28, 0x0b,
	0x32, 0x0c, 0x2e, 0x63, 0x63, 0x2e, 0x43, 0x68, 0x61, 0x74, 0x4d, 0x65, 0x74, 0x61, 0x48, 0x01,
	0x52, 0x04, 0x6d, 0x65, 0x74, 0x61, 0x88, 0x01, 0x01, 0x12, 0x18, 0x0a, 0x07, 0x63, 0x72, 0x65,
	0x61, 0x74, 0x65, 0x64, 0x18, 0x09, 0x20, 0x01, 0x28, 0x03, 0x52, 0x07, 0x63, 0x72, 0x65, 0x61,
	0x74, 0x65, 0x64, 0x42, 0x08, 0x0a, 0x06, 0x5f, 0x74, 0x6f, 0x70, 0x69, 0x63, 0x42, 0x07, 0x0a,
	0x05, 0x5f, 0x6d, 0x65, 0x74, 0x61, 0x22, 0x27, 0x0a, 0x05, 0x43, 0x68, 0x61, 0x74, 0x73, 0x12,
	0x1e, 0x0a, 0x05, 0x63, 0x68, 0x61, 0x74, 0x73, 0x18, 0x01, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x08,
	0x2e, 0x63, 0x63, 0x2e, 0x43, 0x68, 0x61, 0x74, 0x52, 0x05, 0x63, 0x68, 0x61, 0x74, 0x73, 0x22,
	0x3a, 0x0a, 0x0a, 0x41, 0x74, 0x74, 0x61, 0x63, 0x68, 0x6d, 0x65, 0x6e, 0x74, 0x12, 0x12, 0x0a,
	0x04, 0x74, 0x79, 0x70, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x04, 0x74, 0x79, 0x70,
	0x65, 0x12, 0x18, 0x0a, 0x07, 0x63, 0x6f, 0x6e, 0x74, 0x65, 0x6e, 0x74, 0x18, 0x02, 0x20, 0x01,
	0x28, 0x0c, 0x52, 0x07, 0x63, 0x6f, 0x6e, 0x74, 0x65, 0x6e, 0x74, 0x22, 0xac, 0x02, 0x0a, 0x07,
	0x4d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x12, 0x12, 0x0a, 0x04, 0x75, 0x75, 0x69, 0x64, 0x18,
	0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x04, 0x75, 0x75, 0x69, 0x64, 0x12, 0x1c, 0x0a, 0x04, 0x6b,
	0x69, 0x6e, 0x64, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0e, 0x32, 0x08, 0x2e, 0x63, 0x63, 0x2e, 0x4b,
	0x69, 0x6e, 0x64, 0x52, 0x04, 0x6b, 0x69, 0x6e, 0x64, 0x12, 0x16, 0x0a, 0x06, 0x73, 0x65, 0x6e,
	0x64, 0x65, 0x72, 0x18, 0x03, 0x20, 0x01, 0x28, 0x09, 0x52, 0x06, 0x73, 0x65, 0x6e, 0x64, 0x65,
	0x72, 0x12, 0x18, 0x0a, 0x07, 0x63, 0x6f, 0x6e, 0x74, 0x65, 0x6e, 0x74, 0x18, 0x04, 0x20, 0x01,
	0x28, 0x09, 0x52, 0x07, 0x63, 0x6f, 0x6e, 0x74, 0x65, 0x6e, 0x74, 0x12, 0x30, 0x0a, 0x0b, 0x61,
	0x74, 0x74, 0x61, 0x63, 0x68, 0x6d, 0x65, 0x6e, 0x74, 0x73, 0x18, 0x05, 0x20, 0x03, 0x28, 0x0b,
	0x32, 0x0e, 0x2e, 0x63, 0x63, 0x2e, 0x41, 0x74, 0x74, 0x61, 0x63, 0x68, 0x6d, 0x65, 0x6e, 0x74,
	0x52, 0x0b, 0x61, 0x74, 0x74, 0x61, 0x63, 0x68, 0x6d, 0x65, 0x6e, 0x74, 0x73, 0x12, 0x1a, 0x0a,
	0x08, 0x67, 0x61, 0x74, 0x65, 0x77, 0x61, 0x79, 0x73, 0x18, 0x06, 0x20, 0x03, 0x28, 0x09, 0x52,
	0x08, 0x67, 0x61, 0x74, 0x65, 0x77, 0x61, 0x79, 0x73, 0x12, 0x17, 0x0a, 0x04, 0x63, 0x68, 0x61,
	0x74, 0x18, 0x07, 0x20, 0x01, 0x28, 0x09, 0x48, 0x00, 0x52, 0x04, 0x63, 0x68, 0x61, 0x74, 0x88,
	0x01, 0x01, 0x12, 0x12, 0x0a, 0x04, 0x73, 0x65, 0x6e, 0x74, 0x18, 0x08, 0x20, 0x01, 0x28, 0x03,
	0x52, 0x04, 0x73, 0x65, 0x6e, 0x74, 0x12, 0x16, 0x0a, 0x06, 0x65, 0x64, 0x69, 0x74, 0x65, 0x64,
	0x18, 0x09, 0x20, 0x01, 0x28, 0x03, 0x52, 0x06, 0x65, 0x64, 0x69, 0x74, 0x65, 0x64, 0x12, 0x21,
	0x0a, 0x0c, 0x75, 0x6e, 0x64, 0x65, 0x72, 0x5f, 0x72, 0x65, 0x76, 0x69, 0x65, 0x77, 0x18, 0x0a,
	0x20, 0x01, 0x28, 0x08, 0x52, 0x0b, 0x75, 0x6e, 0x64, 0x65, 0x72, 0x52, 0x65, 0x76, 0x69, 0x65,
	0x77, 0x42, 0x07, 0x0a, 0x05, 0x5f, 0x63, 0x68, 0x61, 0x74, 0x22, 0x33, 0x0a, 0x08, 0x4d, 0x65,
	0x73, 0x73, 0x61, 0x67, 0x65, 0x73, 0x12, 0x27, 0x0a, 0x08, 0x6d, 0x65, 0x73, 0x73, 0x61, 0x67,
	0x65, 0x73, 0x18, 0x01, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x0b, 0x2e, 0x63, 0x63, 0x2e, 0x4d, 0x65,
	0x73, 0x73, 0x61, 0x67, 0x65, 0x52, 0x08, 0x6d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x73, 0x22,
	0x30, 0x0a, 0x04, 0x55, 0x73, 0x65, 0x72, 0x12, 0x12, 0x0a, 0x04, 0x75, 0x75, 0x69, 0x64, 0x18,
	0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x04, 0x75, 0x75, 0x69, 0x64, 0x12, 0x14, 0x0a, 0x05, 0x74,
	0x69, 0x74, 0x6c, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x05, 0x74, 0x69, 0x74, 0x6c,
	0x65, 0x22, 0x3e, 0x0a, 0x08, 0x44, 0x65, 0x66, 0x61, 0x75, 0x6c, 0x74, 0x73, 0x12, 0x1a, 0x0a,
	0x08, 0x67, 0x61, 0x74, 0x65, 0x77, 0x61, 0x79, 0x73, 0x18, 0x01, 0x20, 0x03, 0x28, 0x09, 0x52,
	0x08, 0x67, 0x61, 0x74, 0x65, 0x77, 0x61, 0x79, 0x73, 0x12, 0x16, 0x0a, 0x06, 0x61, 0x64, 0x6d,
	0x69, 0x6e, 0x73, 0x18, 0x02, 0x20, 0x03, 0x28, 0x09, 0x52, 0x06, 0x61, 0x64, 0x6d, 0x69, 0x6e,
	0x73, 0x22, 0x27, 0x0a, 0x05, 0x55, 0x73, 0x65, 0x72, 0x73, 0x12, 0x1e, 0x0a, 0x05, 0x75, 0x73,
	0x65, 0x72, 0x73, 0x18, 0x01, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x08, 0x2e, 0x63, 0x63, 0x2e, 0x55,
	0x73, 0x65, 0x72, 0x52, 0x05, 0x75, 0x73, 0x65, 0x72, 0x73, 0x22, 0x73, 0x0a, 0x05, 0x45, 0x76,
	0x65, 0x6e, 0x74, 0x12, 0x21, 0x0a, 0x04, 0x74, 0x79, 0x70, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28,
	0x0e, 0x32, 0x0d, 0x2e, 0x63, 0x63, 0x2e, 0x45, 0x76, 0x65, 0x6e, 0x74, 0x54, 0x79, 0x70, 0x65,
	0x52, 0x04, 0x74, 0x79, 0x70, 0x65, 0x12, 0x1e, 0x0a, 0x04, 0x63, 0x68, 0x61, 0x74, 0x18, 0x02,
	0x20, 0x01, 0x28, 0x0b, 0x32, 0x08, 0x2e, 0x63, 0x63, 0x2e, 0x43, 0x68, 0x61, 0x74, 0x48, 0x00,
	0x52, 0x04, 0x63, 0x68, 0x61, 0x74, 0x12, 0x1f, 0x0a, 0x03, 0x6d, 0x73, 0x67, 0x18, 0x03, 0x20,
	0x01, 0x28, 0x0b, 0x32, 0x0b, 0x2e, 0x63, 0x63, 0x2e, 0x4d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65,
	0x48, 0x00, 0x52, 0x03, 0x6d, 0x73, 0x67, 0x42, 0x06, 0x0a, 0x04, 0x69, 0x74, 0x65, 0x6d, 0x2a,
	0x34, 0x0a, 0x04, 0x52, 0x6f, 0x6c, 0x65, 0x12, 0x0c, 0x0a, 0x08, 0x4e, 0x4f, 0x41, 0x43, 0x43,
	0x45, 0x53, 0x53, 0x10, 0x00, 0x12, 0x08, 0x0a, 0x04, 0x55, 0x53, 0x45, 0x52, 0x10, 0x01, 0x12,
	0x09, 0x0a, 0x05, 0x4f, 0x57, 0x4e, 0x45, 0x52, 0x10, 0x02, 0x12, 0x09, 0x0a, 0x05, 0x41, 0x44,
	0x4d, 0x49, 0x4e, 0x10, 0x03, 0x2a, 0x23, 0x0a, 0x04, 0x4b, 0x69, 0x6e, 0x64, 0x12, 0x0b, 0x0a,
	0x07, 0x44, 0x45, 0x46, 0x41, 0x55, 0x4c, 0x54, 0x10, 0x00, 0x12, 0x0e, 0x0a, 0x0a, 0x41, 0x44,
	0x4d, 0x49, 0x4e, 0x5f, 0x4f, 0x4e, 0x4c, 0x59, 0x10, 0x01, 0x2a, 0x87, 0x01, 0x0a, 0x09, 0x45,
	0x76, 0x65, 0x6e, 0x74, 0x54, 0x79, 0x70, 0x65, 0x12, 0x08, 0x0a, 0x04, 0x50, 0x49, 0x4e, 0x47,
	0x10, 0x00, 0x12, 0x10, 0x0a, 0x0c, 0x43, 0x48, 0x41, 0x54, 0x5f, 0x43, 0x52, 0x45, 0x41, 0x54,
	0x45, 0x44, 0x10, 0x01, 0x12, 0x10, 0x0a, 0x0c, 0x43, 0x48, 0x41, 0x54, 0x5f, 0x55, 0x50, 0x44,
	0x41, 0x54, 0x45, 0x44, 0x10, 0x02, 0x12, 0x10, 0x0a, 0x0c, 0x43, 0x48, 0x41, 0x54, 0x5f, 0x44,
	0x45, 0x4c, 0x45, 0x54, 0x45, 0x44, 0x10, 0x03, 0x12, 0x10, 0x0a, 0x0c, 0x4d, 0x45, 0x53, 0x53,
	0x41, 0x47, 0x45, 0x5f, 0x53, 0x45, 0x4e, 0x44, 0x10, 0x04, 0x12, 0x13, 0x0a, 0x0f, 0x4d, 0x45,
	0x53, 0x53, 0x41, 0x47, 0x45, 0x5f, 0x55, 0x50, 0x44, 0x41, 0x54, 0x45, 0x44, 0x10, 0x05, 0x12,
	0x13, 0x0a, 0x0f, 0x4d, 0x45, 0x53, 0x53, 0x41, 0x47, 0x45, 0x5f, 0x44, 0x45, 0x4c, 0x45, 0x54,
	0x45, 0x44, 0x10, 0x06, 0x32, 0x82, 0x01, 0x0a, 0x08, 0x43, 0x68, 0x61, 0x74, 0x73, 0x41, 0x50,
	0x49, 0x12, 0x1c, 0x0a, 0x06, 0x43, 0x72, 0x65, 0x61, 0x74, 0x65, 0x12, 0x08, 0x2e, 0x63, 0x63,
	0x2e, 0x43, 0x68, 0x61, 0x74, 0x1a, 0x08, 0x2e, 0x63, 0x63, 0x2e, 0x43, 0x68, 0x61, 0x74, 0x12,
	0x1c, 0x0a, 0x06, 0x55, 0x70, 0x64, 0x61, 0x74, 0x65, 0x12, 0x08, 0x2e, 0x63, 0x63, 0x2e, 0x43,
	0x68, 0x61, 0x74, 0x1a, 0x08, 0x2e, 0x63, 0x63, 0x2e, 0x43, 0x68, 0x61, 0x74, 0x12, 0x1c, 0x0a,
	0x04, 0x4c, 0x69, 0x73, 0x74, 0x12, 0x09, 0x2e, 0x63, 0x63, 0x2e, 0x45, 0x6d, 0x70, 0x74, 0x79,
	0x1a, 0x09, 0x2e, 0x63, 0x63, 0x2e, 0x43, 0x68, 0x61, 0x74, 0x73, 0x12, 0x1c, 0x0a, 0x06, 0x44,
	0x65, 0x6c, 0x65, 0x74, 0x65, 0x12, 0x08, 0x2e, 0x63, 0x63, 0x2e, 0x43, 0x68, 0x61, 0x74, 0x1a,
	0x08, 0x2e, 0x63, 0x63, 0x2e, 0x43, 0x68, 0x61, 0x74, 0x32, 0x96, 0x01, 0x0a, 0x0b, 0x4d, 0x65,
	0x73, 0x73, 0x61, 0x67, 0x65, 0x73, 0x41, 0x50, 0x49, 0x12, 0x1d, 0x0a, 0x03, 0x47, 0x65, 0x74,
	0x12, 0x08, 0x2e, 0x63, 0x63, 0x2e, 0x43, 0x68, 0x61, 0x74, 0x1a, 0x0c, 0x2e, 0x63, 0x63, 0x2e,
	0x4d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x73, 0x12, 0x20, 0x0a, 0x04, 0x53, 0x65, 0x6e, 0x64,
	0x12, 0x0b, 0x2e, 0x63, 0x63, 0x2e, 0x4d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x1a, 0x0b, 0x2e,
	0x63, 0x63, 0x2e, 0x4d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x12, 0x22, 0x0a, 0x06, 0x55, 0x70,
	0x64, 0x61, 0x74, 0x65, 0x12, 0x0b, 0x2e, 0x63, 0x63, 0x2e, 0x4d, 0x65, 0x73, 0x73, 0x61, 0x67,
	0x65, 0x1a, 0x0b, 0x2e, 0x63, 0x63, 0x2e, 0x4d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x12, 0x22,
	0x0a, 0x06, 0x44, 0x65, 0x6c, 0x65, 0x74, 0x65, 0x12, 0x0b, 0x2e, 0x63, 0x63, 0x2e, 0x4d, 0x65,
	0x73, 0x73, 0x61, 0x67, 0x65, 0x1a, 0x0b, 0x2e, 0x63, 0x63, 0x2e, 0x4d, 0x65, 0x73, 0x73, 0x61,
	0x67, 0x65, 0x32, 0x70, 0x0a, 0x08, 0x55, 0x73, 0x65, 0x72, 0x73, 0x41, 0x50, 0x49, 0x12, 0x19,
	0x0a, 0x02, 0x4d, 0x65, 0x12, 0x09, 0x2e, 0x63, 0x63, 0x2e, 0x45, 0x6d, 0x70, 0x74, 0x79, 0x1a,
	0x08, 0x2e, 0x63, 0x63, 0x2e, 0x55, 0x73, 0x65, 0x72, 0x12, 0x28, 0x0a, 0x0d, 0x46, 0x65, 0x74,
	0x63, 0x68, 0x44, 0x65, 0x66, 0x61, 0x75, 0x6c, 0x74, 0x73, 0x12, 0x09, 0x2e, 0x63, 0x63, 0x2e,
	0x45, 0x6d, 0x70, 0x74, 0x79, 0x1a, 0x0c, 0x2e, 0x63, 0x63, 0x2e, 0x44, 0x65, 0x66, 0x61, 0x75,
	0x6c, 0x74, 0x73, 0x12, 0x1f, 0x0a, 0x07, 0x52, 0x65, 0x73, 0x6f, 0x6c, 0x76, 0x65, 0x12, 0x09,
	0x2e, 0x63, 0x63, 0x2e, 0x55, 0x73, 0x65, 0x72, 0x73, 0x1a, 0x09, 0x2e, 0x63, 0x63, 0x2e, 0x55,
	0x73, 0x65, 0x72, 0x73, 0x32, 0x31, 0x0a, 0x0d, 0x53, 0x74, 0x72, 0x65, 0x61, 0x6d, 0x53, 0x65,
	0x72, 0x76, 0x69, 0x63, 0x65, 0x12, 0x20, 0x0a, 0x06, 0x53, 0x74, 0x72, 0x65, 0x61, 0x6d, 0x12,
	0x09, 0x2e, 0x63, 0x63, 0x2e, 0x45, 0x6d, 0x70, 0x74, 0x79, 0x1a, 0x09, 0x2e, 0x63, 0x63, 0x2e,
	0x45, 0x76, 0x65, 0x6e, 0x74, 0x30, 0x01, 0x42, 0x25, 0x5a, 0x23, 0x67, 0x69, 0x74, 0x68, 0x75,
	0x62, 0x2e, 0x63, 0x6f, 0x6d, 0x2f, 0x73, 0x6c, 0x6e, 0x74, 0x6f, 0x70, 0x70, 0x2f, 0x63, 0x6f,
	0x72, 0x65, 0x2d, 0x63, 0x68, 0x61, 0x74, 0x74, 0x69, 0x6e, 0x67, 0x2f, 0x63, 0x63, 0x62, 0x06,
	0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
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

var file_cc_cc_proto_enumTypes = make([]protoimpl.EnumInfo, 3)
var file_cc_cc_proto_msgTypes = make([]protoimpl.MessageInfo, 11)
var file_cc_cc_proto_goTypes = []interface{}{
	(Role)(0),          // 0: cc.Role
	(Kind)(0),          // 1: cc.Kind
	(EventType)(0),     // 2: cc.EventType
	(*Empty)(nil),      // 3: cc.Empty
	(*ChatMeta)(nil),   // 4: cc.ChatMeta
	(*Chat)(nil),       // 5: cc.Chat
	(*Chats)(nil),      // 6: cc.Chats
	(*Attachment)(nil), // 7: cc.Attachment
	(*Message)(nil),    // 8: cc.Message
	(*Messages)(nil),   // 9: cc.Messages
	(*User)(nil),       // 10: cc.User
	(*Defaults)(nil),   // 11: cc.Defaults
	(*Users)(nil),      // 12: cc.Users
	(*Event)(nil),      // 13: cc.Event
}
var file_cc_cc_proto_depIdxs = []int32{
	8,  // 0: cc.ChatMeta.last_message:type_name -> cc.Message
	0,  // 1: cc.Chat.role:type_name -> cc.Role
	4,  // 2: cc.Chat.meta:type_name -> cc.ChatMeta
	5,  // 3: cc.Chats.chats:type_name -> cc.Chat
	1,  // 4: cc.Message.kind:type_name -> cc.Kind
	7,  // 5: cc.Message.attachments:type_name -> cc.Attachment
	8,  // 6: cc.Messages.messages:type_name -> cc.Message
	10, // 7: cc.Users.users:type_name -> cc.User
	2,  // 8: cc.Event.type:type_name -> cc.EventType
	5,  // 9: cc.Event.chat:type_name -> cc.Chat
	8,  // 10: cc.Event.msg:type_name -> cc.Message
	5,  // 11: cc.ChatsAPI.Create:input_type -> cc.Chat
	5,  // 12: cc.ChatsAPI.Update:input_type -> cc.Chat
	3,  // 13: cc.ChatsAPI.List:input_type -> cc.Empty
	5,  // 14: cc.ChatsAPI.Delete:input_type -> cc.Chat
	5,  // 15: cc.MessagesAPI.Get:input_type -> cc.Chat
	8,  // 16: cc.MessagesAPI.Send:input_type -> cc.Message
	8,  // 17: cc.MessagesAPI.Update:input_type -> cc.Message
	8,  // 18: cc.MessagesAPI.Delete:input_type -> cc.Message
	3,  // 19: cc.UsersAPI.Me:input_type -> cc.Empty
	3,  // 20: cc.UsersAPI.FetchDefaults:input_type -> cc.Empty
	12, // 21: cc.UsersAPI.Resolve:input_type -> cc.Users
	3,  // 22: cc.StreamService.Stream:input_type -> cc.Empty
	5,  // 23: cc.ChatsAPI.Create:output_type -> cc.Chat
	5,  // 24: cc.ChatsAPI.Update:output_type -> cc.Chat
	6,  // 25: cc.ChatsAPI.List:output_type -> cc.Chats
	5,  // 26: cc.ChatsAPI.Delete:output_type -> cc.Chat
	9,  // 27: cc.MessagesAPI.Get:output_type -> cc.Messages
	8,  // 28: cc.MessagesAPI.Send:output_type -> cc.Message
	8,  // 29: cc.MessagesAPI.Update:output_type -> cc.Message
	8,  // 30: cc.MessagesAPI.Delete:output_type -> cc.Message
	10, // 31: cc.UsersAPI.Me:output_type -> cc.User
	11, // 32: cc.UsersAPI.FetchDefaults:output_type -> cc.Defaults
	12, // 33: cc.UsersAPI.Resolve:output_type -> cc.Users
	13, // 34: cc.StreamService.Stream:output_type -> cc.Event
	23, // [23:35] is the sub-list for method output_type
	11, // [11:23] is the sub-list for method input_type
	11, // [11:11] is the sub-list for extension type_name
	11, // [11:11] is the sub-list for extension extendee
	0,  // [0:11] is the sub-list for field type_name
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
			switch v := v.(*ChatMeta); i {
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
		file_cc_cc_proto_msgTypes[3].Exporter = func(v interface{}, i int) interface{} {
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
		file_cc_cc_proto_msgTypes[4].Exporter = func(v interface{}, i int) interface{} {
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
		file_cc_cc_proto_msgTypes[5].Exporter = func(v interface{}, i int) interface{} {
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
		file_cc_cc_proto_msgTypes[6].Exporter = func(v interface{}, i int) interface{} {
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
		file_cc_cc_proto_msgTypes[7].Exporter = func(v interface{}, i int) interface{} {
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
		file_cc_cc_proto_msgTypes[8].Exporter = func(v interface{}, i int) interface{} {
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
		file_cc_cc_proto_msgTypes[9].Exporter = func(v interface{}, i int) interface{} {
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
		file_cc_cc_proto_msgTypes[10].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Event); i {
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
	file_cc_cc_proto_msgTypes[2].OneofWrappers = []interface{}{}
	file_cc_cc_proto_msgTypes[5].OneofWrappers = []interface{}{}
	file_cc_cc_proto_msgTypes[10].OneofWrappers = []interface{}{
		(*Event_Chat)(nil),
		(*Event_Msg)(nil),
	}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_cc_cc_proto_rawDesc,
			NumEnums:      3,
			NumMessages:   11,
			NumExtensions: 0,
			NumServices:   4,
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
