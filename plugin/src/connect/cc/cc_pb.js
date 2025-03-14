// @generated by protoc-gen-es v1.10.0 with parameter "target=js+dts"
// @generated from file cc/cc.proto (package cc, syntax proto3)
/* eslint-disable */
// @ts-nocheck

import { proto3, Struct, Value } from "@bufbuild/protobuf";

/**
 * @generated from enum cc.Role
 */
export const Role = /*@__PURE__*/ proto3.makeEnum(
  "cc.Role",
  [
    {no: 0, name: "NOACCESS"},
    {no: 1, name: "USER"},
    {no: 2, name: "OWNER"},
    {no: 3, name: "ADMIN"},
  ],
);

/**
 * @generated from enum cc.Status
 */
export const Status = /*@__PURE__*/ proto3.makeEnum(
  "cc.Status",
  [
    {no: 0, name: "NEW"},
    {no: 1, name: "OPEN"},
    {no: 2, name: "RESOLVE"},
    {no: 3, name: "CLOSE"},
    {no: 4, name: "ANSWERED"},
    {no: 5, name: "CUSTOMER_REPLY"},
    {no: 6, name: "WAITING_FOR_REPLY"},
    {no: 7, name: "ON_HOLD"},
    {no: 8, name: "IN_PROGRESS"},
  ],
);

/**
 * @generated from enum cc.Kind
 */
export const Kind = /*@__PURE__*/ proto3.makeEnum(
  "cc.Kind",
  [
    {no: 0, name: "DEFAULT"},
    {no: 1, name: "ADMIN_ONLY"},
    {no: 2, name: "FOR_BOT"},
  ],
);

/**
 * @generated from enum cc.EventType
 */
export const EventType = /*@__PURE__*/ proto3.makeEnum(
  "cc.EventType",
  [
    {no: 0, name: "PING"},
    {no: 1, name: "CHAT_CREATED"},
    {no: 2, name: "CHAT_UPDATED"},
    {no: 3, name: "CHAT_DELETED"},
    {no: 4, name: "CHAT_READ"},
    {no: 5, name: "MESSAGE_SENT"},
    {no: 6, name: "MESSAGE_UPDATED"},
    {no: 7, name: "MESSAGE_DELETED"},
    {no: 8, name: "CHAT_DEPARTMENT_CHANGED"},
    {no: 9, name: "CHAT_STATUS_CHANGED"},
    {no: 10, name: "CHATS_MERGED"},
  ],
);

/**
 * @generated from message cc.Empty
 */
export const Empty = /*@__PURE__*/ proto3.makeMessageType(
  "cc.Empty",
  [],
);

/**
 * @generated from message cc.ChatMeta
 */
export const ChatMeta = /*@__PURE__*/ proto3.makeMessageType(
  "cc.ChatMeta",
  () => [
    { no: 1, name: "unread", kind: "scalar", T: 5 /* ScalarType.INT32 */ },
    { no: 2, name: "first_message", kind: "message", T: Message },
    { no: 3, name: "last_message", kind: "message", T: Message },
    { no: 4, name: "data", kind: "map", K: 9 /* ScalarType.STRING */, V: {kind: "message", T: Value} },
  ],
);

/**
 * @generated from message cc.Chat
 */
export const Chat = /*@__PURE__*/ proto3.makeMessageType(
  "cc.Chat",
  () => [
    { no: 1, name: "uuid", kind: "scalar", T: 9 /* ScalarType.STRING */ },
    { no: 2, name: "users", kind: "scalar", T: 9 /* ScalarType.STRING */, repeated: true },
    { no: 3, name: "admins", kind: "scalar", T: 9 /* ScalarType.STRING */, repeated: true },
    { no: 4, name: "owner", kind: "scalar", T: 9 /* ScalarType.STRING */ },
    { no: 5, name: "gateways", kind: "scalar", T: 9 /* ScalarType.STRING */, repeated: true },
    { no: 6, name: "role", kind: "enum", T: proto3.getEnumType(Role) },
    { no: 7, name: "topic", kind: "scalar", T: 9 /* ScalarType.STRING */, opt: true },
    { no: 8, name: "meta", kind: "message", T: ChatMeta, opt: true },
    { no: 9, name: "created", kind: "scalar", T: 3 /* ScalarType.INT64 */ },
    { no: 10, name: "status", kind: "enum", T: proto3.getEnumType(Status) },
    { no: 11, name: "department", kind: "scalar", T: 9 /* ScalarType.STRING */ },
    { no: 12, name: "responsible", kind: "scalar", T: 9 /* ScalarType.STRING */, opt: true },
    { no: 13, name: "bot_state", kind: "map", K: 9 /* ScalarType.STRING */, V: {kind: "message", T: Value} },
  ],
);

/**
 * @generated from message cc.Chats
 */
export const Chats = /*@__PURE__*/ proto3.makeMessageType(
  "cc.Chats",
  () => [
    { no: 1, name: "chats", kind: "message", T: Chat, repeated: true },
  ],
);

/**
 * @generated from message cc.Merge
 */
export const Merge = /*@__PURE__*/ proto3.makeMessageType(
  "cc.Merge",
  () => [
    { no: 1, name: "chats", kind: "scalar", T: 9 /* ScalarType.STRING */, repeated: true },
  ],
);

/**
 * @generated from message cc.ListChatsRequest
 */
export const ListChatsRequest = /*@__PURE__*/ proto3.makeMessageType(
  "cc.ListChatsRequest",
  () => [
    { no: 1, name: "page", kind: "scalar", T: 4 /* ScalarType.UINT64 */, opt: true },
    { no: 2, name: "limit", kind: "scalar", T: 4 /* ScalarType.UINT64 */, opt: true },
    { no: 3, name: "field", kind: "scalar", T: 9 /* ScalarType.STRING */, opt: true },
    { no: 4, name: "sort", kind: "scalar", T: 9 /* ScalarType.STRING */, opt: true },
    { no: 5, name: "filters", kind: "map", K: 9 /* ScalarType.STRING */, V: {kind: "message", T: Value} },
  ],
);

/**
 * @generated from message cc.SetBotStateRequest
 */
export const SetBotStateRequest = /*@__PURE__*/ proto3.makeMessageType(
  "cc.SetBotStateRequest",
  () => [
    { no: 1, name: "chat", kind: "scalar", T: 9 /* ScalarType.STRING */ },
    { no: 2, name: "state", kind: "map", K: 9 /* ScalarType.STRING */, V: {kind: "message", T: Value} },
    { no: 3, name: "disabled", kind: "scalar", T: 8 /* ScalarType.BOOL */, opt: true },
    { no: 4, name: "escalated", kind: "scalar", T: 8 /* ScalarType.BOOL */, opt: true },
    { no: 5, name: "skip_review", kind: "scalar", T: 8 /* ScalarType.BOOL */, opt: true },
  ],
);

/**
 * @generated from message cc.ListChatsResponse
 */
export const ListChatsResponse = /*@__PURE__*/ proto3.makeMessageType(
  "cc.ListChatsResponse",
  () => [
    { no: 1, name: "pool", kind: "message", T: Chat, repeated: true },
    { no: 2, name: "total", kind: "scalar", T: 3 /* ScalarType.INT64 */ },
  ],
);

/**
 * @generated from message cc.CountChatsRequest
 */
export const CountChatsRequest = /*@__PURE__*/ proto3.makeMessageType(
  "cc.CountChatsRequest",
  () => [
    { no: 1, name: "filters", kind: "map", K: 9 /* ScalarType.STRING */, V: {kind: "message", T: Value} },
  ],
);

/**
 * @generated from message cc.CountChatsResponse
 */
export const CountChatsResponse = /*@__PURE__*/ proto3.makeMessageType(
  "cc.CountChatsResponse",
  () => [
    { no: 1, name: "statuses", kind: "map", K: 5 /* ScalarType.INT32 */, V: {kind: "scalar", T: 3 /* ScalarType.INT64 */} },
  ],
);

/**
 * @generated from message cc.Attachment
 */
export const Attachment = /*@__PURE__*/ proto3.makeMessageType(
  "cc.Attachment",
  () => [
    { no: 1, name: "type", kind: "scalar", T: 9 /* ScalarType.STRING */ },
    { no: 2, name: "content", kind: "scalar", T: 12 /* ScalarType.BYTES */ },
  ],
);

/**
 * @generated from message cc.Message
 */
export const Message = /*@__PURE__*/ proto3.makeMessageType(
  "cc.Message",
  () => [
    { no: 1, name: "uuid", kind: "scalar", T: 9 /* ScalarType.STRING */ },
    { no: 2, name: "kind", kind: "enum", T: proto3.getEnumType(Kind) },
    { no: 3, name: "sender", kind: "scalar", T: 9 /* ScalarType.STRING */ },
    { no: 4, name: "content", kind: "scalar", T: 9 /* ScalarType.STRING */ },
    { no: 5, name: "attachments", kind: "scalar", T: 9 /* ScalarType.STRING */, repeated: true },
    { no: 6, name: "gateways", kind: "scalar", T: 9 /* ScalarType.STRING */, repeated: true },
    { no: 7, name: "chat", kind: "scalar", T: 9 /* ScalarType.STRING */, opt: true },
    { no: 8, name: "sent", kind: "scalar", T: 3 /* ScalarType.INT64 */ },
    { no: 9, name: "edited", kind: "scalar", T: 3 /* ScalarType.INT64 */ },
    { no: 10, name: "under_review", kind: "scalar", T: 8 /* ScalarType.BOOL */ },
    { no: 11, name: "readers", kind: "scalar", T: 9 /* ScalarType.STRING */, repeated: true },
    { no: 12, name: "meta", kind: "map", K: 9 /* ScalarType.STRING */, V: {kind: "message", T: Value} },
    { no: 13, name: "mentioned", kind: "scalar", T: 9 /* ScalarType.STRING */, repeated: true },
  ],
);

/**
 * @generated from message cc.Messages
 */
export const Messages = /*@__PURE__*/ proto3.makeMessageType(
  "cc.Messages",
  () => [
    { no: 1, name: "messages", kind: "message", T: Message, repeated: true },
  ],
);

/**
 * @generated from message cc.MessagesListRequest
 */
export const MessagesListRequest = /*@__PURE__*/ proto3.makeMessageType(
  "cc.MessagesListRequest",
  [],
);

/**
 * @generated from message cc.User
 */
export const User = /*@__PURE__*/ proto3.makeMessageType(
  "cc.User",
  () => [
    { no: 1, name: "uuid", kind: "scalar", T: 9 /* ScalarType.STRING */ },
    { no: 2, name: "title", kind: "scalar", T: 9 /* ScalarType.STRING */ },
    { no: 3, name: "data", kind: "message", T: Struct },
    { no: 4, name: "cc_is_bot", kind: "scalar", T: 8 /* ScalarType.BOOL */ },
    { no: 5, name: "cc_username", kind: "scalar", T: 9 /* ScalarType.STRING */ },
    { no: 6, name: "cc_commands", kind: "map", K: 9 /* ScalarType.STRING */, V: {kind: "scalar", T: 9 /* ScalarType.STRING */} },
  ],
);

/**
 * @generated from message cc.Department
 */
export const Department = /*@__PURE__*/ proto3.makeMessageType(
  "cc.Department",
  () => [
    { no: 1, name: "key", kind: "scalar", T: 9 /* ScalarType.STRING */ },
    { no: 2, name: "title", kind: "scalar", T: 9 /* ScalarType.STRING */ },
    { no: 3, name: "description", kind: "scalar", T: 9 /* ScalarType.STRING */ },
    { no: 4, name: "admins", kind: "scalar", T: 9 /* ScalarType.STRING */, repeated: true },
    { no: 5, name: "public", kind: "scalar", T: 8 /* ScalarType.BOOL */ },
    { no: 6, name: "whmcs_id", kind: "scalar", T: 9 /* ScalarType.STRING */, opt: true },
  ],
);

/**
 * @generated from message cc.Option
 */
export const Option = /*@__PURE__*/ proto3.makeMessageType(
  "cc.Option",
  () => [
    { no: 1, name: "key", kind: "scalar", T: 9 /* ScalarType.STRING */ },
    { no: 2, name: "value", kind: "scalar", T: 2 /* ScalarType.FLOAT */ },
  ],
);

/**
 * @generated from message cc.Metric
 */
export const Metric = /*@__PURE__*/ proto3.makeMessageType(
  "cc.Metric",
  () => [
    { no: 1, name: "title", kind: "scalar", T: 9 /* ScalarType.STRING */ },
    { no: 2, name: "options", kind: "message", T: Option, repeated: true },
  ],
);

/**
 * @generated from message cc.Defaults
 */
export const Defaults = /*@__PURE__*/ proto3.makeMessageType(
  "cc.Defaults",
  () => [
    { no: 1, name: "gateways", kind: "scalar", T: 9 /* ScalarType.STRING */, repeated: true },
    { no: 2, name: "admins", kind: "scalar", T: 9 /* ScalarType.STRING */, repeated: true },
    { no: 3, name: "departments", kind: "message", T: Department, repeated: true },
    { no: 4, name: "metrics", kind: "map", K: 9 /* ScalarType.STRING */, V: {kind: "message", T: Metric} },
    { no: 5, name: "templates", kind: "map", K: 9 /* ScalarType.STRING */, V: {kind: "scalar", T: 9 /* ScalarType.STRING */} },
    { no: 6, name: "bot", kind: "message", T: Bot },
  ],
);

/**
 * @generated from message cc.Bot
 */
export const Bot = /*@__PURE__*/ proto3.makeMessageType(
  "cc.Bot",
  () => [
    { no: 1, name: "prompt", kind: "scalar", T: 9 /* ScalarType.STRING */ },
    { no: 2, name: "review", kind: "scalar", T: 8 /* ScalarType.BOOL */ },
    { no: 3, name: "enable", kind: "scalar", T: 8 /* ScalarType.BOOL */ },
    { no: 4, name: "values", kind: "map", K: 9 /* ScalarType.STRING */, V: {kind: "scalar", T: 9 /* ScalarType.STRING */} },
    { no: 5, name: "initiator", kind: "scalar", T: 8 /* ScalarType.BOOL */ },
    { no: 6, name: "emergency", kind: "scalar", T: 8 /* ScalarType.BOOL */ },
  ],
);

/**
 * @generated from message cc.Users
 */
export const Users = /*@__PURE__*/ proto3.makeMessageType(
  "cc.Users",
  () => [
    { no: 1, name: "users", kind: "message", T: User, repeated: true },
    { no: 2, name: "total", kind: "scalar", T: 3 /* ScalarType.INT64 */ },
  ],
);

/**
 * @generated from message cc.FetchDefaultsRequest
 */
export const FetchDefaultsRequest = /*@__PURE__*/ proto3.makeMessageType(
  "cc.FetchDefaultsRequest",
  () => [
    { no: 1, name: "fetch_templates", kind: "scalar", T: 8 /* ScalarType.BOOL */ },
  ],
);

/**
 * @generated from message cc.GetMembersRequest
 */
export const GetMembersRequest = /*@__PURE__*/ proto3.makeMessageType(
  "cc.GetMembersRequest",
  () => [
    { no: 1, name: "uuids", kind: "scalar", T: 9 /* ScalarType.STRING */, repeated: true },
    { no: 2, name: "page", kind: "scalar", T: 4 /* ScalarType.UINT64 */, opt: true },
    { no: 3, name: "limit", kind: "scalar", T: 4 /* ScalarType.UINT64 */, opt: true },
    { no: 4, name: "field", kind: "scalar", T: 9 /* ScalarType.STRING */, opt: true },
    { no: 5, name: "sort", kind: "scalar", T: 9 /* ScalarType.STRING */, opt: true },
    { no: 6, name: "filters", kind: "map", K: 9 /* ScalarType.STRING */, V: {kind: "message", T: Value} },
  ],
);

/**
 * @generated from message cc.Event
 */
export const Event = /*@__PURE__*/ proto3.makeMessageType(
  "cc.Event",
  () => [
    { no: 1, name: "type", kind: "enum", T: proto3.getEnumType(EventType) },
    { no: 2, name: "chat", kind: "message", T: Chat, oneof: "item" },
    { no: 3, name: "msg", kind: "message", T: Message, oneof: "item" },
  ],
);

/**
 * @generated from message cc.StreamRequest
 */
export const StreamRequest = /*@__PURE__*/ proto3.makeMessageType(
  "cc.StreamRequest",
  () => [
    { no: 1, name: "commands", kind: "map", K: 9 /* ScalarType.STRING */, V: {kind: "scalar", T: 9 /* ScalarType.STRING */} },
  ],
);

