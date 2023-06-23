// @generated by protoc-gen-connect-es v0.9.1 with parameter "target=ts"
// @generated from file cc/cc.proto (package cc, syntax proto3)
/* eslint-disable */
// @ts-nocheck

import { Chat, Chats, Defaults, Empty, Message, Messages, Users } from "./cc_pb.js";
import { MethodKind } from "@bufbuild/protobuf";

/**
 * @generated from service cc.ChatsAPI
 */
export const ChatsAPI = {
  typeName: "cc.ChatsAPI",
  methods: {
    /**
     * @generated from rpc cc.ChatsAPI.Create
     */
    create: {
      name: "Create",
      I: Chat,
      O: Chat,
      kind: MethodKind.Unary,
    },
    /**
     * @generated from rpc cc.ChatsAPI.Update
     */
    update: {
      name: "Update",
      I: Chat,
      O: Chat,
      kind: MethodKind.Unary,
    },
    /**
     * @generated from rpc cc.ChatsAPI.List
     */
    list: {
      name: "List",
      I: Empty,
      O: Chats,
      kind: MethodKind.Unary,
    },
    /**
     * @generated from rpc cc.ChatsAPI.Delete
     */
    delete: {
      name: "Delete",
      I: Chat,
      O: Chat,
      kind: MethodKind.Unary,
    },
  }
} as const;

/**
 * @generated from service cc.MessagesAPI
 */
export const MessagesAPI = {
  typeName: "cc.MessagesAPI",
  methods: {
    /**
     * @generated from rpc cc.MessagesAPI.Get
     */
    get: {
      name: "Get",
      I: Chat,
      O: Messages,
      kind: MethodKind.Unary,
    },
    /**
     * @generated from rpc cc.MessagesAPI.Send
     */
    send: {
      name: "Send",
      I: Message,
      O: Message,
      kind: MethodKind.Unary,
    },
    /**
     * @generated from rpc cc.MessagesAPI.Update
     */
    update: {
      name: "Update",
      I: Message,
      O: Message,
      kind: MethodKind.Unary,
    },
    /**
     * @generated from rpc cc.MessagesAPI.Delete
     */
    delete: {
      name: "Delete",
      I: Message,
      O: Message,
      kind: MethodKind.Unary,
    },
  }
} as const;

/**
 * @generated from service cc.UsersAPI
 */
export const UsersAPI = {
  typeName: "cc.UsersAPI",
  methods: {
    /**
     * @generated from rpc cc.UsersAPI.FetchDefaults
     */
    fetchDefaults: {
      name: "FetchDefaults",
      I: Empty,
      O: Defaults,
      kind: MethodKind.Unary,
    },
    /**
     * @generated from rpc cc.UsersAPI.Resolve
     */
    resolve: {
      name: "Resolve",
      I: Empty,
      O: Users,
      kind: MethodKind.Unary,
    },
  }
} as const;

