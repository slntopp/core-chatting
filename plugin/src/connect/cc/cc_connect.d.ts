// @generated by protoc-gen-connect-es v1.1.4 with parameter "target=js+dts"
// @generated from file cc/cc.proto (package cc, syntax proto3)
/* eslint-disable */
// @ts-nocheck

import { Chat, Chats, Defaults, Empty, Event, Merge, Message, Messages, StreamRequest, User, Users } from "./cc_pb.js";
import { MethodKind } from "@bufbuild/protobuf";

/**
 * @generated from service cc.ChatsAPI
 */
export declare const ChatsAPI: {
  readonly typeName: "cc.ChatsAPI",
  readonly methods: {
    /**
     * @generated from rpc cc.ChatsAPI.Create
     */
    readonly create: {
      readonly name: "Create",
      readonly I: typeof Chat,
      readonly O: typeof Chat,
      readonly kind: MethodKind.Unary,
    },
    /**
     * @generated from rpc cc.ChatsAPI.Update
     */
    readonly update: {
      readonly name: "Update",
      readonly I: typeof Chat,
      readonly O: typeof Chat,
      readonly kind: MethodKind.Unary,
    },
    /**
     * @generated from rpc cc.ChatsAPI.Get
     */
    readonly get: {
      readonly name: "Get",
      readonly I: typeof Chat,
      readonly O: typeof Chat,
      readonly kind: MethodKind.Unary,
    },
    /**
     * @generated from rpc cc.ChatsAPI.List
     */
    readonly list: {
      readonly name: "List",
      readonly I: typeof Empty,
      readonly O: typeof Chats,
      readonly kind: MethodKind.Unary,
    },
    /**
     * @generated from rpc cc.ChatsAPI.Delete
     */
    readonly delete: {
      readonly name: "Delete",
      readonly I: typeof Chat,
      readonly O: typeof Chat,
      readonly kind: MethodKind.Unary,
    },
    /**
     * @generated from rpc cc.ChatsAPI.SetBotState
     */
    readonly setBotState: {
      readonly name: "SetBotState",
      readonly I: typeof Chat,
      readonly O: typeof Chat,
      readonly kind: MethodKind.Unary,
    },
    /**
     * @generated from rpc cc.ChatsAPI.GetBotState
     */
    readonly getBotState: {
      readonly name: "GetBotState",
      readonly I: typeof Chat,
      readonly O: typeof Chat,
      readonly kind: MethodKind.Unary,
    },
    /**
     * @generated from rpc cc.ChatsAPI.ChangeDepartment
     */
    readonly changeDepartment: {
      readonly name: "ChangeDepartment",
      readonly I: typeof Chat,
      readonly O: typeof Chat,
      readonly kind: MethodKind.Unary,
    },
    /**
     * @generated from rpc cc.ChatsAPI.ChangeGateway
     */
    readonly changeGateway: {
      readonly name: "ChangeGateway",
      readonly I: typeof Chat,
      readonly O: typeof Chat,
      readonly kind: MethodKind.Unary,
    },
    /**
     * @generated from rpc cc.ChatsAPI.ChangeStatus
     */
    readonly changeStatus: {
      readonly name: "ChangeStatus",
      readonly I: typeof Chat,
      readonly O: typeof Chat,
      readonly kind: MethodKind.Unary,
    },
    /**
     * @generated from rpc cc.ChatsAPI.MergeChats
     */
    readonly mergeChats: {
      readonly name: "MergeChats",
      readonly I: typeof Merge,
      readonly O: typeof Chat,
      readonly kind: MethodKind.Unary,
    },
    /**
     * @generated from rpc cc.ChatsAPI.SyncChats
     */
    readonly syncChats: {
      readonly name: "SyncChats",
      readonly I: typeof Empty,
      readonly O: typeof Empty,
      readonly kind: MethodKind.Unary,
    },
  }
};

/**
 * @generated from service cc.MessagesAPI
 */
export declare const MessagesAPI: {
  readonly typeName: "cc.MessagesAPI",
  readonly methods: {
    /**
     * @generated from rpc cc.MessagesAPI.Get
     */
    readonly get: {
      readonly name: "Get",
      readonly I: typeof Chat,
      readonly O: typeof Messages,
      readonly kind: MethodKind.Unary,
    },
    /**
     * @generated from rpc cc.MessagesAPI.Send
     */
    readonly send: {
      readonly name: "Send",
      readonly I: typeof Message,
      readonly O: typeof Message,
      readonly kind: MethodKind.Unary,
    },
    /**
     * @generated from rpc cc.MessagesAPI.Update
     */
    readonly update: {
      readonly name: "Update",
      readonly I: typeof Message,
      readonly O: typeof Message,
      readonly kind: MethodKind.Unary,
    },
    /**
     * @generated from rpc cc.MessagesAPI.Delete
     */
    readonly delete: {
      readonly name: "Delete",
      readonly I: typeof Message,
      readonly O: typeof Message,
      readonly kind: MethodKind.Unary,
    },
  }
};

/**
 * @generated from service cc.UsersAPI
 */
export declare const UsersAPI: {
  readonly typeName: "cc.UsersAPI",
  readonly methods: {
    /**
     * @generated from rpc cc.UsersAPI.Me
     */
    readonly me: {
      readonly name: "Me",
      readonly I: typeof Empty,
      readonly O: typeof User,
      readonly kind: MethodKind.Unary,
    },
    /**
     * @generated from rpc cc.UsersAPI.FetchDefaults
     */
    readonly fetchDefaults: {
      readonly name: "FetchDefaults",
      readonly I: typeof Empty,
      readonly O: typeof Defaults,
      readonly kind: MethodKind.Unary,
    },
    /**
     * @generated from rpc cc.UsersAPI.GetConfig
     */
    readonly getConfig: {
      readonly name: "GetConfig",
      readonly I: typeof Empty,
      readonly O: typeof Defaults,
      readonly kind: MethodKind.Unary,
    },
    /**
     * @generated from rpc cc.UsersAPI.SetConfig
     */
    readonly setConfig: {
      readonly name: "SetConfig",
      readonly I: typeof Defaults,
      readonly O: typeof Defaults,
      readonly kind: MethodKind.Unary,
    },
    /**
     * Resolves given Users data by their UUIDs
     * And returns all accessible Users for Requestor
     *
     * @generated from rpc cc.UsersAPI.Resolve
     */
    readonly resolve: {
      readonly name: "Resolve",
      readonly I: typeof Users,
      readonly O: typeof Users,
      readonly kind: MethodKind.Unary,
    },
    /**
     * @generated from rpc cc.UsersAPI.GetMembers
     */
    readonly getMembers: {
      readonly name: "GetMembers",
      readonly I: typeof Empty,
      readonly O: typeof Users,
      readonly kind: MethodKind.Unary,
    },
  }
};

/**
 * @generated from service cc.StreamService
 */
export declare const StreamService: {
  readonly typeName: "cc.StreamService",
  readonly methods: {
    /**
     * @generated from rpc cc.StreamService.Stream
     */
    readonly stream: {
      readonly name: "Stream",
      readonly I: typeof StreamRequest,
      readonly O: typeof Event,
      readonly kind: MethodKind.ServerStreaming,
    },
  }
};

