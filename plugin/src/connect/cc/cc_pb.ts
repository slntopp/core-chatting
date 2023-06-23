// @generated by protoc-gen-es v1.2.1 with parameter "target=ts"
// @generated from file cc/cc.proto (package cc, syntax proto3)
/* eslint-disable */
// @ts-nocheck

import type { BinaryReadOptions, FieldList, JsonReadOptions, JsonValue, PartialMessage, PlainMessage } from "@bufbuild/protobuf";
import { Message as Message$1, proto3 } from "@bufbuild/protobuf";

/**
 * @generated from enum cc.Role
 */
export enum Role {
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
// Retrieve enum metadata with: proto3.getEnumType(Role)
proto3.util.setEnumType(Role, "cc.Role", [
  { no: 0, name: "NOACCESS" },
  { no: 1, name: "USER" },
  { no: 2, name: "OWNER" },
  { no: 3, name: "ADMIN" },
]);

/**
 * @generated from enum cc.Kind
 */
export enum Kind {
  /**
   * @generated from enum value: DEFAULT = 0;
   */
  DEFAULT = 0,

  /**
   * @generated from enum value: ADMIN_ONLY = 1;
   */
  ADMIN_ONLY = 1,
}
// Retrieve enum metadata with: proto3.getEnumType(Kind)
proto3.util.setEnumType(Kind, "cc.Kind", [
  { no: 0, name: "DEFAULT" },
  { no: 1, name: "ADMIN_ONLY" },
]);

/**
 * @generated from message cc.Empty
 */
export class Empty extends Message$1<Empty> {
  constructor(data?: PartialMessage<Empty>) {
    super();
    proto3.util.initPartial(data, this);
  }

  static readonly runtime: typeof proto3 = proto3;
  static readonly typeName = "cc.Empty";
  static readonly fields: FieldList = proto3.util.newFieldList(() => [
  ]);

  static fromBinary(bytes: Uint8Array, options?: Partial<BinaryReadOptions>): Empty {
    return new Empty().fromBinary(bytes, options);
  }

  static fromJson(jsonValue: JsonValue, options?: Partial<JsonReadOptions>): Empty {
    return new Empty().fromJson(jsonValue, options);
  }

  static fromJsonString(jsonString: string, options?: Partial<JsonReadOptions>): Empty {
    return new Empty().fromJsonString(jsonString, options);
  }

  static equals(a: Empty | PlainMessage<Empty> | undefined, b: Empty | PlainMessage<Empty> | undefined): boolean {
    return proto3.util.equals(Empty, a, b);
  }
}

/**
 * @generated from message cc.Chat
 */
export class Chat extends Message$1<Chat> {
  /**
   * @generated from field: string uuid = 1;
   */
  uuid = "";

  /**
   * @generated from field: repeated string users = 2;
   */
  users: string[] = [];

  /**
   * @generated from field: repeated string admins = 3;
   */
  admins: string[] = [];

  /**
   * @generated from field: string owner = 4;
   */
  owner = "";

  /**
   * @generated from field: repeated string gateways = 5;
   */
  gateways: string[] = [];

  /**
   * @generated from field: cc.Role role = 6;
   */
  role = Role.NOACCESS;

  /**
   * @generated from field: optional string topic = 7;
   */
  topic?: string;

  constructor(data?: PartialMessage<Chat>) {
    super();
    proto3.util.initPartial(data, this);
  }

  static readonly runtime: typeof proto3 = proto3;
  static readonly typeName = "cc.Chat";
  static readonly fields: FieldList = proto3.util.newFieldList(() => [
    { no: 1, name: "uuid", kind: "scalar", T: 9 /* ScalarType.STRING */ },
    { no: 2, name: "users", kind: "scalar", T: 9 /* ScalarType.STRING */, repeated: true },
    { no: 3, name: "admins", kind: "scalar", T: 9 /* ScalarType.STRING */, repeated: true },
    { no: 4, name: "owner", kind: "scalar", T: 9 /* ScalarType.STRING */ },
    { no: 5, name: "gateways", kind: "scalar", T: 9 /* ScalarType.STRING */, repeated: true },
    { no: 6, name: "role", kind: "enum", T: proto3.getEnumType(Role) },
    { no: 7, name: "topic", kind: "scalar", T: 9 /* ScalarType.STRING */, opt: true },
  ]);

  static fromBinary(bytes: Uint8Array, options?: Partial<BinaryReadOptions>): Chat {
    return new Chat().fromBinary(bytes, options);
  }

  static fromJson(jsonValue: JsonValue, options?: Partial<JsonReadOptions>): Chat {
    return new Chat().fromJson(jsonValue, options);
  }

  static fromJsonString(jsonString: string, options?: Partial<JsonReadOptions>): Chat {
    return new Chat().fromJsonString(jsonString, options);
  }

  static equals(a: Chat | PlainMessage<Chat> | undefined, b: Chat | PlainMessage<Chat> | undefined): boolean {
    return proto3.util.equals(Chat, a, b);
  }
}

/**
 * @generated from message cc.Chats
 */
export class Chats extends Message$1<Chats> {
  /**
   * @generated from field: repeated cc.Chat chats = 1;
   */
  chats: Chat[] = [];

  constructor(data?: PartialMessage<Chats>) {
    super();
    proto3.util.initPartial(data, this);
  }

  static readonly runtime: typeof proto3 = proto3;
  static readonly typeName = "cc.Chats";
  static readonly fields: FieldList = proto3.util.newFieldList(() => [
    { no: 1, name: "chats", kind: "message", T: Chat, repeated: true },
  ]);

  static fromBinary(bytes: Uint8Array, options?: Partial<BinaryReadOptions>): Chats {
    return new Chats().fromBinary(bytes, options);
  }

  static fromJson(jsonValue: JsonValue, options?: Partial<JsonReadOptions>): Chats {
    return new Chats().fromJson(jsonValue, options);
  }

  static fromJsonString(jsonString: string, options?: Partial<JsonReadOptions>): Chats {
    return new Chats().fromJsonString(jsonString, options);
  }

  static equals(a: Chats | PlainMessage<Chats> | undefined, b: Chats | PlainMessage<Chats> | undefined): boolean {
    return proto3.util.equals(Chats, a, b);
  }
}

/**
 * @generated from message cc.Attachment
 */
export class Attachment extends Message$1<Attachment> {
  /**
   * @generated from field: string type = 1;
   */
  type = "";

  /**
   * @generated from field: bytes content = 2;
   */
  content = new Uint8Array(0);

  constructor(data?: PartialMessage<Attachment>) {
    super();
    proto3.util.initPartial(data, this);
  }

  static readonly runtime: typeof proto3 = proto3;
  static readonly typeName = "cc.Attachment";
  static readonly fields: FieldList = proto3.util.newFieldList(() => [
    { no: 1, name: "type", kind: "scalar", T: 9 /* ScalarType.STRING */ },
    { no: 2, name: "content", kind: "scalar", T: 12 /* ScalarType.BYTES */ },
  ]);

  static fromBinary(bytes: Uint8Array, options?: Partial<BinaryReadOptions>): Attachment {
    return new Attachment().fromBinary(bytes, options);
  }

  static fromJson(jsonValue: JsonValue, options?: Partial<JsonReadOptions>): Attachment {
    return new Attachment().fromJson(jsonValue, options);
  }

  static fromJsonString(jsonString: string, options?: Partial<JsonReadOptions>): Attachment {
    return new Attachment().fromJsonString(jsonString, options);
  }

  static equals(a: Attachment | PlainMessage<Attachment> | undefined, b: Attachment | PlainMessage<Attachment> | undefined): boolean {
    return proto3.util.equals(Attachment, a, b);
  }
}

/**
 * @generated from message cc.Message
 */
export class Message extends Message$1<Message> {
  /**
   * @generated from field: string uuid = 1;
   */
  uuid = "";

  /**
   * @generated from field: cc.Kind kind = 2;
   */
  kind = Kind.DEFAULT;

  /**
   * @generated from field: string sender = 3;
   */
  sender = "";

  /**
   * @generated from field: string content = 4;
   */
  content = "";

  /**
   * @generated from field: repeated cc.Attachment attachments = 5;
   */
  attachments: Attachment[] = [];

  /**
   * @generated from field: repeated string gateways = 6;
   */
  gateways: string[] = [];

  /**
   * Required for Send, Update and Delete requests
   *
   * @generated from field: optional string chat = 7;
   */
  chat?: string;

  constructor(data?: PartialMessage<Message>) {
    super();
    proto3.util.initPartial(data, this);
  }

  static readonly runtime: typeof proto3 = proto3;
  static readonly typeName = "cc.Message";
  static readonly fields: FieldList = proto3.util.newFieldList(() => [
    { no: 1, name: "uuid", kind: "scalar", T: 9 /* ScalarType.STRING */ },
    { no: 2, name: "kind", kind: "enum", T: proto3.getEnumType(Kind) },
    { no: 3, name: "sender", kind: "scalar", T: 9 /* ScalarType.STRING */ },
    { no: 4, name: "content", kind: "scalar", T: 9 /* ScalarType.STRING */ },
    { no: 5, name: "attachments", kind: "message", T: Attachment, repeated: true },
    { no: 6, name: "gateways", kind: "scalar", T: 9 /* ScalarType.STRING */, repeated: true },
    { no: 7, name: "chat", kind: "scalar", T: 9 /* ScalarType.STRING */, opt: true },
  ]);

  static fromBinary(bytes: Uint8Array, options?: Partial<BinaryReadOptions>): Message {
    return new Message().fromBinary(bytes, options);
  }

  static fromJson(jsonValue: JsonValue, options?: Partial<JsonReadOptions>): Message {
    return new Message().fromJson(jsonValue, options);
  }

  static fromJsonString(jsonString: string, options?: Partial<JsonReadOptions>): Message {
    return new Message().fromJsonString(jsonString, options);
  }

  static equals(a: Message | PlainMessage<Message> | undefined, b: Message | PlainMessage<Message> | undefined): boolean {
    return proto3.util.equals(Message, a, b);
  }
}

/**
 * @generated from message cc.Messages
 */
export class Messages extends Message$1<Messages> {
  /**
   * @generated from field: repeated cc.Message messages = 1;
   */
  messages: Message[] = [];

  constructor(data?: PartialMessage<Messages>) {
    super();
    proto3.util.initPartial(data, this);
  }

  static readonly runtime: typeof proto3 = proto3;
  static readonly typeName = "cc.Messages";
  static readonly fields: FieldList = proto3.util.newFieldList(() => [
    { no: 1, name: "messages", kind: "message", T: Message, repeated: true },
  ]);

  static fromBinary(bytes: Uint8Array, options?: Partial<BinaryReadOptions>): Messages {
    return new Messages().fromBinary(bytes, options);
  }

  static fromJson(jsonValue: JsonValue, options?: Partial<JsonReadOptions>): Messages {
    return new Messages().fromJson(jsonValue, options);
  }

  static fromJsonString(jsonString: string, options?: Partial<JsonReadOptions>): Messages {
    return new Messages().fromJsonString(jsonString, options);
  }

  static equals(a: Messages | PlainMessage<Messages> | undefined, b: Messages | PlainMessage<Messages> | undefined): boolean {
    return proto3.util.equals(Messages, a, b);
  }
}

/**
 * @generated from message cc.User
 */
export class User extends Message$1<User> {
  /**
   * @generated from field: string uuid = 1;
   */
  uuid = "";

  /**
   * @generated from field: string title = 2;
   */
  title = "";

  constructor(data?: PartialMessage<User>) {
    super();
    proto3.util.initPartial(data, this);
  }

  static readonly runtime: typeof proto3 = proto3;
  static readonly typeName = "cc.User";
  static readonly fields: FieldList = proto3.util.newFieldList(() => [
    { no: 1, name: "uuid", kind: "scalar", T: 9 /* ScalarType.STRING */ },
    { no: 2, name: "title", kind: "scalar", T: 9 /* ScalarType.STRING */ },
  ]);

  static fromBinary(bytes: Uint8Array, options?: Partial<BinaryReadOptions>): User {
    return new User().fromBinary(bytes, options);
  }

  static fromJson(jsonValue: JsonValue, options?: Partial<JsonReadOptions>): User {
    return new User().fromJson(jsonValue, options);
  }

  static fromJsonString(jsonString: string, options?: Partial<JsonReadOptions>): User {
    return new User().fromJsonString(jsonString, options);
  }

  static equals(a: User | PlainMessage<User> | undefined, b: User | PlainMessage<User> | undefined): boolean {
    return proto3.util.equals(User, a, b);
  }
}

/**
 * @generated from message cc.Defaults
 */
export class Defaults extends Message$1<Defaults> {
  /**
   * @generated from field: repeated string gateways = 1;
   */
  gateways: string[] = [];

  /**
   * @generated from field: repeated cc.User admins = 2;
   */
  admins: User[] = [];

  constructor(data?: PartialMessage<Defaults>) {
    super();
    proto3.util.initPartial(data, this);
  }

  static readonly runtime: typeof proto3 = proto3;
  static readonly typeName = "cc.Defaults";
  static readonly fields: FieldList = proto3.util.newFieldList(() => [
    { no: 1, name: "gateways", kind: "scalar", T: 9 /* ScalarType.STRING */, repeated: true },
    { no: 2, name: "admins", kind: "message", T: User, repeated: true },
  ]);

  static fromBinary(bytes: Uint8Array, options?: Partial<BinaryReadOptions>): Defaults {
    return new Defaults().fromBinary(bytes, options);
  }

  static fromJson(jsonValue: JsonValue, options?: Partial<JsonReadOptions>): Defaults {
    return new Defaults().fromJson(jsonValue, options);
  }

  static fromJsonString(jsonString: string, options?: Partial<JsonReadOptions>): Defaults {
    return new Defaults().fromJsonString(jsonString, options);
  }

  static equals(a: Defaults | PlainMessage<Defaults> | undefined, b: Defaults | PlainMessage<Defaults> | undefined): boolean {
    return proto3.util.equals(Defaults, a, b);
  }
}

/**
 * @generated from message cc.Users
 */
export class Users extends Message$1<Users> {
  /**
   * @generated from field: repeated cc.User users = 1;
   */
  users: User[] = [];

  constructor(data?: PartialMessage<Users>) {
    super();
    proto3.util.initPartial(data, this);
  }

  static readonly runtime: typeof proto3 = proto3;
  static readonly typeName = "cc.Users";
  static readonly fields: FieldList = proto3.util.newFieldList(() => [
    { no: 1, name: "users", kind: "message", T: User, repeated: true },
  ]);

  static fromBinary(bytes: Uint8Array, options?: Partial<BinaryReadOptions>): Users {
    return new Users().fromBinary(bytes, options);
  }

  static fromJson(jsonValue: JsonValue, options?: Partial<JsonReadOptions>): Users {
    return new Users().fromJson(jsonValue, options);
  }

  static fromJsonString(jsonString: string, options?: Partial<JsonReadOptions>): Users {
    return new Users().fromJsonString(jsonString, options);
  }

  static equals(a: Users | PlainMessage<Users> | undefined, b: Users | PlainMessage<Users> | undefined): boolean {
    return proto3.util.equals(Users, a, b);
  }
}

