// @generated by protoc-gen-es v1.3.1 with parameter "target=js+dts"
// @generated from file cc/cc.proto (package cc, syntax proto3)
/* eslint-disable */
// @ts-nocheck

import type { BinaryReadOptions, FieldList, JsonReadOptions, JsonValue, PartialMessage, PlainMessage, Struct, Value } from "@bufbuild/protobuf";
import { Message as Message$1, proto3 } from "@bufbuild/protobuf";

/**
 * @generated from enum cc.Role
 */
export declare enum Role {
  /**
   * @generated from enum value: NOACCESS = 0;
   */
  NOACCESS = 0,

  /**
   * @generated from enum value: USER = 1;
   */
  USER = 1,

  /**
   * @generated from enum value: OWNER = 2;
   */
  OWNER = 2,

  /**
   * @generated from enum value: ADMIN = 3;
   */
  ADMIN = 3,
}

/**
 * @generated from enum cc.Status
 */
export declare enum Status {
  /**
   * @generated from enum value: NEW = 0;
   */
  NEW = 0,

  /**
   * @generated from enum value: OPEN = 1;
   */
  OPEN = 1,

  /**
   * @generated from enum value: RESOLVE = 2;
   */
  RESOLVE = 2,

  /**
   * @generated from enum value: CLOSE = 3;
   */
  CLOSE = 3,
}

/**
 * @generated from enum cc.Kind
 */
export declare enum Kind {
  /**
   * @generated from enum value: DEFAULT = 0;
   */
  DEFAULT = 0,

  /**
   * @generated from enum value: FOR_BOT = 1;
   */
  FOR_BOT = 1,

  /**
   * @generated from enum value: ADMIN_ONLY = 2;
   */
  ADMIN_ONLY = 2,
}

/**
 * @generated from enum cc.EventType
 */
export declare enum EventType {
  /**
   * @generated from enum value: PING = 0;
   */
  PING = 0,

  /**
   * @generated from enum value: CHAT_CREATED = 1;
   */
  CHAT_CREATED = 1,

  /**
   * @generated from enum value: CHAT_UPDATED = 2;
   */
  CHAT_UPDATED = 2,

  /**
   * @generated from enum value: CHAT_DELETED = 3;
   */
  CHAT_DELETED = 3,

  /**
   * @generated from enum value: CHAT_READ = 4;
   */
  CHAT_READ = 4,

  /**
   * @generated from enum value: MESSAGE_SENT = 5;
   */
  MESSAGE_SENT = 5,

  /**
   * @generated from enum value: MESSAGE_UPDATED = 6;
   */
  MESSAGE_UPDATED = 6,

  /**
   * @generated from enum value: MESSAGE_DELETED = 7;
   */
  MESSAGE_DELETED = 7,
}

/**
 * @generated from message cc.Empty
 */
export declare class Empty extends Message$1<Empty> {
  constructor(data?: PartialMessage<Empty>);

  static readonly runtime: typeof proto3;
  static readonly typeName = "cc.Empty";
  static readonly fields: FieldList;

  static fromBinary(bytes: Uint8Array, options?: Partial<BinaryReadOptions>): Empty;

  static fromJson(jsonValue: JsonValue, options?: Partial<JsonReadOptions>): Empty;

  static fromJsonString(jsonString: string, options?: Partial<JsonReadOptions>): Empty;

  static equals(a: Empty | PlainMessage<Empty> | undefined, b: Empty | PlainMessage<Empty> | undefined): boolean;
}

/**
 * @generated from message cc.ChatMeta
 */
export declare class ChatMeta extends Message$1<ChatMeta> {
  /**
   * @generated from field: int32 unread = 1;
   */
  unread: number;

  /**
   * @generated from field: cc.Message first_message = 2;
   */
  firstMessage?: Message;

  /**
   * @generated from field: cc.Message last_message = 3;
   */
  lastMessage?: Message;

  /**
   * @generated from field: map<string, google.protobuf.Value> data = 4;
   */
  data: { [key: string]: Value };

  constructor(data?: PartialMessage<ChatMeta>);

  static readonly runtime: typeof proto3;
  static readonly typeName = "cc.ChatMeta";
  static readonly fields: FieldList;

  static fromBinary(bytes: Uint8Array, options?: Partial<BinaryReadOptions>): ChatMeta;

  static fromJson(jsonValue: JsonValue, options?: Partial<JsonReadOptions>): ChatMeta;

  static fromJsonString(jsonString: string, options?: Partial<JsonReadOptions>): ChatMeta;

  static equals(a: ChatMeta | PlainMessage<ChatMeta> | undefined, b: ChatMeta | PlainMessage<ChatMeta> | undefined): boolean;
}

/**
 * @generated from message cc.Chat
 */
export declare class Chat extends Message$1<Chat> {
  /**
   * @generated from field: string uuid = 1;
   */
  uuid: string;

  /**
   * @generated from field: repeated string users = 2;
   */
  users: string[];

  /**
   * @generated from field: repeated string admins = 3;
   */
  admins: string[];

  /**
   * @generated from field: string owner = 4;
   */
  owner: string;

  /**
   * @generated from field: repeated string gateways = 5;
   */
  gateways: string[];

  /**
   * @generated from field: cc.Role role = 6;
   */
  role: Role;

  /**
   * @generated from field: optional string topic = 7;
   */
  topic?: string;

  /**
   * @generated from field: optional cc.ChatMeta meta = 8;
   */
  meta?: ChatMeta;

  /**
   * @generated from field: int64 created = 9;
   */
  created: bigint;

  /**
   * @generated from field: cc.Status status = 10;
   */
  status: Status;

  /**
   * @generated from field: string department = 11;
   */
  department: string;

  /**
   * @generated from field: optional string responsible = 12;
   */
  responsible?: string;

  /**
   * @generated from field: map<string, google.protobuf.Value> bot_state = 13;
   */
  botState: { [key: string]: Value };

  constructor(data?: PartialMessage<Chat>);

  static readonly runtime: typeof proto3;
  static readonly typeName = "cc.Chat";
  static readonly fields: FieldList;

  static fromBinary(bytes: Uint8Array, options?: Partial<BinaryReadOptions>): Chat;

  static fromJson(jsonValue: JsonValue, options?: Partial<JsonReadOptions>): Chat;

  static fromJsonString(jsonString: string, options?: Partial<JsonReadOptions>): Chat;

  static equals(a: Chat | PlainMessage<Chat> | undefined, b: Chat | PlainMessage<Chat> | undefined): boolean;
}

/**
 * @generated from message cc.Chats
 */
export declare class Chats extends Message$1<Chats> {
  /**
   * @generated from field: repeated cc.Chat chats = 1;
   */
  chats: Chat[];

  constructor(data?: PartialMessage<Chats>);

  static readonly runtime: typeof proto3;
  static readonly typeName = "cc.Chats";
  static readonly fields: FieldList;

  static fromBinary(bytes: Uint8Array, options?: Partial<BinaryReadOptions>): Chats;

  static fromJson(jsonValue: JsonValue, options?: Partial<JsonReadOptions>): Chats;

  static fromJsonString(jsonString: string, options?: Partial<JsonReadOptions>): Chats;

  static equals(a: Chats | PlainMessage<Chats> | undefined, b: Chats | PlainMessage<Chats> | undefined): boolean;
}

/**
 * @generated from message cc.Attachment
 */
export declare class Attachment extends Message$1<Attachment> {
  /**
   * @generated from field: string type = 1;
   */
  type: string;

  /**
   * @generated from field: bytes content = 2;
   */
  content: Uint8Array;

  constructor(data?: PartialMessage<Attachment>);

  static readonly runtime: typeof proto3;
  static readonly typeName = "cc.Attachment";
  static readonly fields: FieldList;

  static fromBinary(bytes: Uint8Array, options?: Partial<BinaryReadOptions>): Attachment;

  static fromJson(jsonValue: JsonValue, options?: Partial<JsonReadOptions>): Attachment;

  static fromJsonString(jsonString: string, options?: Partial<JsonReadOptions>): Attachment;

  static equals(a: Attachment | PlainMessage<Attachment> | undefined, b: Attachment | PlainMessage<Attachment> | undefined): boolean;
}

/**
 * @generated from message cc.Message
 */
export declare class Message extends Message$1<Message> {
  /**
   * @generated from field: string uuid = 1;
   */
  uuid: string;

  /**
   * @generated from field: cc.Kind kind = 2;
   */
  kind: Kind;

  /**
   * @generated from field: string sender = 3;
   */
  sender: string;

  /**
   * @generated from field: string content = 4;
   */
  content: string;

  /**
   * @generated from field: repeated cc.Attachment attachments = 5;
   */
  attachments: Attachment[];

  /**
   * @generated from field: repeated string gateways = 6;
   */
  gateways: string[];

  /**
   * Required for Send, Update and Delete requests
   *
   * @generated from field: optional string chat = 7;
   */
  chat?: string;

  /**
   * @generated from field: int64 sent = 8;
   */
  sent: bigint;

  /**
   * @generated from field: int64 edited = 9;
   */
  edited: bigint;

  /**
   * @generated from field: bool under_review = 10;
   */
  underReview: boolean;

  /**
   * @generated from field: repeated string readers = 11;
   */
  readers: string[];

  /**
   * @generated from field: map<string, google.protobuf.Value> meta = 12;
   */
  meta: { [key: string]: Value };

  /**
   * @generated from field: repeated string mentioned = 13;
   */
  mentioned: string[];

  constructor(data?: PartialMessage<Message>);

  static readonly runtime: typeof proto3;
  static readonly typeName = "cc.Message";
  static readonly fields: FieldList;

  static fromBinary(bytes: Uint8Array, options?: Partial<BinaryReadOptions>): Message;

  static fromJson(jsonValue: JsonValue, options?: Partial<JsonReadOptions>): Message;

  static fromJsonString(jsonString: string, options?: Partial<JsonReadOptions>): Message;

  static equals(a: Message | PlainMessage<Message> | undefined, b: Message | PlainMessage<Message> | undefined): boolean;
}

/**
 * @generated from message cc.Messages
 */
export declare class Messages extends Message$1<Messages> {
  /**
   * @generated from field: repeated cc.Message messages = 1;
   */
  messages: Message[];

  constructor(data?: PartialMessage<Messages>);

  static readonly runtime: typeof proto3;
  static readonly typeName = "cc.Messages";
  static readonly fields: FieldList;

  static fromBinary(bytes: Uint8Array, options?: Partial<BinaryReadOptions>): Messages;

  static fromJson(jsonValue: JsonValue, options?: Partial<JsonReadOptions>): Messages;

  static fromJsonString(jsonString: string, options?: Partial<JsonReadOptions>): Messages;

  static equals(a: Messages | PlainMessage<Messages> | undefined, b: Messages | PlainMessage<Messages> | undefined): boolean;
}

/**
 * @generated from message cc.User
 */
export declare class User extends Message$1<User> {
  /**
   * @generated from field: string uuid = 1;
   */
  uuid: string;

  /**
   * @generated from field: string title = 2;
   */
  title: string;

  /**
   * @generated from field: google.protobuf.Struct data = 3;
   */
  data?: Struct;

  /**
   * @generated from field: bool cc_is_bot = 4;
   */
  ccIsBot: boolean;

  /**
   * @generated from field: string cc_username = 5;
   */
  ccUsername: string;

  /**
   * @generated from field: map<string, string> cc_commands = 6;
   */
  ccCommands: { [key: string]: string };

  constructor(data?: PartialMessage<User>);

  static readonly runtime: typeof proto3;
  static readonly typeName = "cc.User";
  static readonly fields: FieldList;

  static fromBinary(bytes: Uint8Array, options?: Partial<BinaryReadOptions>): User;

  static fromJson(jsonValue: JsonValue, options?: Partial<JsonReadOptions>): User;

  static fromJsonString(jsonString: string, options?: Partial<JsonReadOptions>): User;

  static equals(a: User | PlainMessage<User> | undefined, b: User | PlainMessage<User> | undefined): boolean;
}

/**
 * @generated from message cc.Department
 */
export declare class Department extends Message$1<Department> {
  /**
   * @generated from field: string key = 1;
   */
  key: string;

  /**
   * @generated from field: string title = 2;
   */
  title: string;

  /**
   * @generated from field: string description = 3;
   */
  description: string;

  /**
   * @generated from field: repeated string admins = 4;
   */
  admins: string[];

  constructor(data?: PartialMessage<Department>);

  static readonly runtime: typeof proto3;
  static readonly typeName = "cc.Department";
  static readonly fields: FieldList;

  static fromBinary(bytes: Uint8Array, options?: Partial<BinaryReadOptions>): Department;

  static fromJson(jsonValue: JsonValue, options?: Partial<JsonReadOptions>): Department;

  static fromJsonString(jsonString: string, options?: Partial<JsonReadOptions>): Department;

  static equals(a: Department | PlainMessage<Department> | undefined, b: Department | PlainMessage<Department> | undefined): boolean;
}

/**
 * @generated from message cc.Option
 */
export declare class Option extends Message$1<Option> {
  /**
   * @generated from field: string key = 1;
   */
  key: string;

  /**
   * @generated from field: float value = 2;
   */
  value: number;

  constructor(data?: PartialMessage<Option>);

  static readonly runtime: typeof proto3;
  static readonly typeName = "cc.Option";
  static readonly fields: FieldList;

  static fromBinary(bytes: Uint8Array, options?: Partial<BinaryReadOptions>): Option;

  static fromJson(jsonValue: JsonValue, options?: Partial<JsonReadOptions>): Option;

  static fromJsonString(jsonString: string, options?: Partial<JsonReadOptions>): Option;

  static equals(a: Option | PlainMessage<Option> | undefined, b: Option | PlainMessage<Option> | undefined): boolean;
}

/**
 * @generated from message cc.Metric
 */
export declare class Metric extends Message$1<Metric> {
  /**
   * @generated from field: string title = 1;
   */
  title: string;

  /**
   * @generated from field: repeated cc.Option options = 2;
   */
  options: Option[];

  constructor(data?: PartialMessage<Metric>);

  static readonly runtime: typeof proto3;
  static readonly typeName = "cc.Metric";
  static readonly fields: FieldList;

  static fromBinary(bytes: Uint8Array, options?: Partial<BinaryReadOptions>): Metric;

  static fromJson(jsonValue: JsonValue, options?: Partial<JsonReadOptions>): Metric;

  static fromJsonString(jsonString: string, options?: Partial<JsonReadOptions>): Metric;

  static equals(a: Metric | PlainMessage<Metric> | undefined, b: Metric | PlainMessage<Metric> | undefined): boolean;
}

/**
 * @generated from message cc.Defaults
 */
export declare class Defaults extends Message$1<Defaults> {
  /**
   * @generated from field: repeated string gateways = 1;
   */
  gateways: string[];

  /**
   * @generated from field: repeated string admins = 2;
   */
  admins: string[];

  /**
   * @generated from field: repeated cc.Department departments = 3;
   */
  departments: Department[];

  /**
   * @generated from field: map<string, cc.Metric> metrics = 4;
   */
  metrics: { [key: string]: Metric };

  constructor(data?: PartialMessage<Defaults>);

  static readonly runtime: typeof proto3;
  static readonly typeName = "cc.Defaults";
  static readonly fields: FieldList;

  static fromBinary(bytes: Uint8Array, options?: Partial<BinaryReadOptions>): Defaults;

  static fromJson(jsonValue: JsonValue, options?: Partial<JsonReadOptions>): Defaults;

  static fromJsonString(jsonString: string, options?: Partial<JsonReadOptions>): Defaults;

  static equals(a: Defaults | PlainMessage<Defaults> | undefined, b: Defaults | PlainMessage<Defaults> | undefined): boolean;
}

/**
 * @generated from message cc.Users
 */
export declare class Users extends Message$1<Users> {
  /**
   * @generated from field: repeated cc.User users = 1;
   */
  users: User[];

  constructor(data?: PartialMessage<Users>);

  static readonly runtime: typeof proto3;
  static readonly typeName = "cc.Users";
  static readonly fields: FieldList;

  static fromBinary(bytes: Uint8Array, options?: Partial<BinaryReadOptions>): Users;

  static fromJson(jsonValue: JsonValue, options?: Partial<JsonReadOptions>): Users;

  static fromJsonString(jsonString: string, options?: Partial<JsonReadOptions>): Users;

  static equals(a: Users | PlainMessage<Users> | undefined, b: Users | PlainMessage<Users> | undefined): boolean;
}

/**
 * @generated from message cc.Event
 */
export declare class Event extends Message$1<Event> {
  /**
   * @generated from field: cc.EventType type = 1;
   */
  type: EventType;

  /**
   * @generated from oneof cc.Event.item
   */
  item: {
    /**
     * @generated from field: cc.Chat chat = 2;
     */
    value: Chat;
    case: "chat";
  } | {
    /**
     * @generated from field: cc.Message msg = 3;
     */
    value: Message;
    case: "msg";
  } | { case: undefined; value?: undefined };

  constructor(data?: PartialMessage<Event>);

  static readonly runtime: typeof proto3;
  static readonly typeName = "cc.Event";
  static readonly fields: FieldList;

  static fromBinary(bytes: Uint8Array, options?: Partial<BinaryReadOptions>): Event;

  static fromJson(jsonValue: JsonValue, options?: Partial<JsonReadOptions>): Event;

  static fromJsonString(jsonString: string, options?: Partial<JsonReadOptions>): Event;

  static equals(a: Event | PlainMessage<Event> | undefined, b: Event | PlainMessage<Event> | undefined): boolean;
}

/**
 * @generated from message cc.StreamRequest
 */
export declare class StreamRequest extends Message$1<StreamRequest> {
  /**
   * @generated from field: map<string, string> commands = 1;
   */
  commands: { [key: string]: string };

  constructor(data?: PartialMessage<StreamRequest>);

  static readonly runtime: typeof proto3;
  static readonly typeName = "cc.StreamRequest";
  static readonly fields: FieldList;

  static fromBinary(bytes: Uint8Array, options?: Partial<BinaryReadOptions>): StreamRequest;

  static fromJson(jsonValue: JsonValue, options?: Partial<JsonReadOptions>): StreamRequest;

  static fromJsonString(jsonString: string, options?: Partial<JsonReadOptions>): StreamRequest;

  static equals(a: StreamRequest | PlainMessage<StreamRequest> | undefined, b: StreamRequest | PlainMessage<StreamRequest> | undefined): boolean;
}

