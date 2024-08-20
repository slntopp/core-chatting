# Generated by the gRPC Python protocol compiler plugin. DO NOT EDIT!
"""Client and server classes corresponding to protobuf-defined services."""
import grpc

from cc import cc_pb2 as cc_dot_cc__pb2


class ChatsAPIStub(object):
    """Missing associated documentation comment in .proto file."""

    def __init__(self, channel):
        """Constructor.

        Args:
            channel: A grpc.Channel.
        """
        self.Create = channel.unary_unary(
                '/cc.ChatsAPI/Create',
                request_serializer=cc_dot_cc__pb2.Chat.SerializeToString,
                response_deserializer=cc_dot_cc__pb2.Chat.FromString,
                )
        self.Update = channel.unary_unary(
                '/cc.ChatsAPI/Update',
                request_serializer=cc_dot_cc__pb2.Chat.SerializeToString,
                response_deserializer=cc_dot_cc__pb2.Chat.FromString,
                )
        self.Get = channel.unary_unary(
                '/cc.ChatsAPI/Get',
                request_serializer=cc_dot_cc__pb2.Chat.SerializeToString,
                response_deserializer=cc_dot_cc__pb2.Chat.FromString,
                )
        self.List = channel.unary_unary(
                '/cc.ChatsAPI/List',
                request_serializer=cc_dot_cc__pb2.Empty.SerializeToString,
                response_deserializer=cc_dot_cc__pb2.Chats.FromString,
                )
        self.Delete = channel.unary_unary(
                '/cc.ChatsAPI/Delete',
                request_serializer=cc_dot_cc__pb2.Chat.SerializeToString,
                response_deserializer=cc_dot_cc__pb2.Chat.FromString,
                )
        self.SetBotState = channel.unary_unary(
                '/cc.ChatsAPI/SetBotState',
                request_serializer=cc_dot_cc__pb2.Chat.SerializeToString,
                response_deserializer=cc_dot_cc__pb2.Chat.FromString,
                )
        self.GetBotState = channel.unary_unary(
                '/cc.ChatsAPI/GetBotState',
                request_serializer=cc_dot_cc__pb2.Chat.SerializeToString,
                response_deserializer=cc_dot_cc__pb2.Chat.FromString,
                )
        self.ChangeDepartment = channel.unary_unary(
                '/cc.ChatsAPI/ChangeDepartment',
                request_serializer=cc_dot_cc__pb2.Chat.SerializeToString,
                response_deserializer=cc_dot_cc__pb2.Chat.FromString,
                )
        self.ChangeGateway = channel.unary_unary(
                '/cc.ChatsAPI/ChangeGateway',
                request_serializer=cc_dot_cc__pb2.Chat.SerializeToString,
                response_deserializer=cc_dot_cc__pb2.Chat.FromString,
                )
        self.ChangeStatus = channel.unary_unary(
                '/cc.ChatsAPI/ChangeStatus',
                request_serializer=cc_dot_cc__pb2.Chat.SerializeToString,
                response_deserializer=cc_dot_cc__pb2.Chat.FromString,
                )
        self.MergeChats = channel.unary_unary(
                '/cc.ChatsAPI/MergeChats',
                request_serializer=cc_dot_cc__pb2.Merge.SerializeToString,
                response_deserializer=cc_dot_cc__pb2.Chat.FromString,
                )
        self.SyncChats = channel.unary_unary(
                '/cc.ChatsAPI/SyncChats',
                request_serializer=cc_dot_cc__pb2.Empty.SerializeToString,
                response_deserializer=cc_dot_cc__pb2.Empty.FromString,
                )


class ChatsAPIServicer(object):
    """Missing associated documentation comment in .proto file."""

    def Create(self, request, context):
        """Missing associated documentation comment in .proto file."""
        context.set_code(grpc.StatusCode.UNIMPLEMENTED)
        context.set_details('Method not implemented!')
        raise NotImplementedError('Method not implemented!')

    def Update(self, request, context):
        """Missing associated documentation comment in .proto file."""
        context.set_code(grpc.StatusCode.UNIMPLEMENTED)
        context.set_details('Method not implemented!')
        raise NotImplementedError('Method not implemented!')

    def Get(self, request, context):
        """Missing associated documentation comment in .proto file."""
        context.set_code(grpc.StatusCode.UNIMPLEMENTED)
        context.set_details('Method not implemented!')
        raise NotImplementedError('Method not implemented!')

    def List(self, request, context):
        """Missing associated documentation comment in .proto file."""
        context.set_code(grpc.StatusCode.UNIMPLEMENTED)
        context.set_details('Method not implemented!')
        raise NotImplementedError('Method not implemented!')

    def Delete(self, request, context):
        """Missing associated documentation comment in .proto file."""
        context.set_code(grpc.StatusCode.UNIMPLEMENTED)
        context.set_details('Method not implemented!')
        raise NotImplementedError('Method not implemented!')

    def SetBotState(self, request, context):
        """Missing associated documentation comment in .proto file."""
        context.set_code(grpc.StatusCode.UNIMPLEMENTED)
        context.set_details('Method not implemented!')
        raise NotImplementedError('Method not implemented!')

    def GetBotState(self, request, context):
        """Missing associated documentation comment in .proto file."""
        context.set_code(grpc.StatusCode.UNIMPLEMENTED)
        context.set_details('Method not implemented!')
        raise NotImplementedError('Method not implemented!')

    def ChangeDepartment(self, request, context):
        """Missing associated documentation comment in .proto file."""
        context.set_code(grpc.StatusCode.UNIMPLEMENTED)
        context.set_details('Method not implemented!')
        raise NotImplementedError('Method not implemented!')

    def ChangeGateway(self, request, context):
        """Missing associated documentation comment in .proto file."""
        context.set_code(grpc.StatusCode.UNIMPLEMENTED)
        context.set_details('Method not implemented!')
        raise NotImplementedError('Method not implemented!')

    def ChangeStatus(self, request, context):
        """Missing associated documentation comment in .proto file."""
        context.set_code(grpc.StatusCode.UNIMPLEMENTED)
        context.set_details('Method not implemented!')
        raise NotImplementedError('Method not implemented!')

    def MergeChats(self, request, context):
        """Missing associated documentation comment in .proto file."""
        context.set_code(grpc.StatusCode.UNIMPLEMENTED)
        context.set_details('Method not implemented!')
        raise NotImplementedError('Method not implemented!')

    def SyncChats(self, request, context):
        """Missing associated documentation comment in .proto file."""
        context.set_code(grpc.StatusCode.UNIMPLEMENTED)
        context.set_details('Method not implemented!')
        raise NotImplementedError('Method not implemented!')


def add_ChatsAPIServicer_to_server(servicer, server):
    rpc_method_handlers = {
            'Create': grpc.unary_unary_rpc_method_handler(
                    servicer.Create,
                    request_deserializer=cc_dot_cc__pb2.Chat.FromString,
                    response_serializer=cc_dot_cc__pb2.Chat.SerializeToString,
            ),
            'Update': grpc.unary_unary_rpc_method_handler(
                    servicer.Update,
                    request_deserializer=cc_dot_cc__pb2.Chat.FromString,
                    response_serializer=cc_dot_cc__pb2.Chat.SerializeToString,
            ),
            'Get': grpc.unary_unary_rpc_method_handler(
                    servicer.Get,
                    request_deserializer=cc_dot_cc__pb2.Chat.FromString,
                    response_serializer=cc_dot_cc__pb2.Chat.SerializeToString,
            ),
            'List': grpc.unary_unary_rpc_method_handler(
                    servicer.List,
                    request_deserializer=cc_dot_cc__pb2.Empty.FromString,
                    response_serializer=cc_dot_cc__pb2.Chats.SerializeToString,
            ),
            'Delete': grpc.unary_unary_rpc_method_handler(
                    servicer.Delete,
                    request_deserializer=cc_dot_cc__pb2.Chat.FromString,
                    response_serializer=cc_dot_cc__pb2.Chat.SerializeToString,
            ),
            'SetBotState': grpc.unary_unary_rpc_method_handler(
                    servicer.SetBotState,
                    request_deserializer=cc_dot_cc__pb2.Chat.FromString,
                    response_serializer=cc_dot_cc__pb2.Chat.SerializeToString,
            ),
            'GetBotState': grpc.unary_unary_rpc_method_handler(
                    servicer.GetBotState,
                    request_deserializer=cc_dot_cc__pb2.Chat.FromString,
                    response_serializer=cc_dot_cc__pb2.Chat.SerializeToString,
            ),
            'ChangeDepartment': grpc.unary_unary_rpc_method_handler(
                    servicer.ChangeDepartment,
                    request_deserializer=cc_dot_cc__pb2.Chat.FromString,
                    response_serializer=cc_dot_cc__pb2.Chat.SerializeToString,
            ),
            'ChangeGateway': grpc.unary_unary_rpc_method_handler(
                    servicer.ChangeGateway,
                    request_deserializer=cc_dot_cc__pb2.Chat.FromString,
                    response_serializer=cc_dot_cc__pb2.Chat.SerializeToString,
            ),
            'ChangeStatus': grpc.unary_unary_rpc_method_handler(
                    servicer.ChangeStatus,
                    request_deserializer=cc_dot_cc__pb2.Chat.FromString,
                    response_serializer=cc_dot_cc__pb2.Chat.SerializeToString,
            ),
            'MergeChats': grpc.unary_unary_rpc_method_handler(
                    servicer.MergeChats,
                    request_deserializer=cc_dot_cc__pb2.Merge.FromString,
                    response_serializer=cc_dot_cc__pb2.Chat.SerializeToString,
            ),
            'SyncChats': grpc.unary_unary_rpc_method_handler(
                    servicer.SyncChats,
                    request_deserializer=cc_dot_cc__pb2.Empty.FromString,
                    response_serializer=cc_dot_cc__pb2.Empty.SerializeToString,
            ),
    }
    generic_handler = grpc.method_handlers_generic_handler(
            'cc.ChatsAPI', rpc_method_handlers)
    server.add_generic_rpc_handlers((generic_handler,))


 # This class is part of an EXPERIMENTAL API.
class ChatsAPI(object):
    """Missing associated documentation comment in .proto file."""

    @staticmethod
    def Create(request,
            target,
            options=(),
            channel_credentials=None,
            call_credentials=None,
            insecure=False,
            compression=None,
            wait_for_ready=None,
            timeout=None,
            metadata=None):
        return grpc.experimental.unary_unary(request, target, '/cc.ChatsAPI/Create',
            cc_dot_cc__pb2.Chat.SerializeToString,
            cc_dot_cc__pb2.Chat.FromString,
            options, channel_credentials,
            insecure, call_credentials, compression, wait_for_ready, timeout, metadata)

    @staticmethod
    def Update(request,
            target,
            options=(),
            channel_credentials=None,
            call_credentials=None,
            insecure=False,
            compression=None,
            wait_for_ready=None,
            timeout=None,
            metadata=None):
        return grpc.experimental.unary_unary(request, target, '/cc.ChatsAPI/Update',
            cc_dot_cc__pb2.Chat.SerializeToString,
            cc_dot_cc__pb2.Chat.FromString,
            options, channel_credentials,
            insecure, call_credentials, compression, wait_for_ready, timeout, metadata)

    @staticmethod
    def Get(request,
            target,
            options=(),
            channel_credentials=None,
            call_credentials=None,
            insecure=False,
            compression=None,
            wait_for_ready=None,
            timeout=None,
            metadata=None):
        return grpc.experimental.unary_unary(request, target, '/cc.ChatsAPI/Get',
            cc_dot_cc__pb2.Chat.SerializeToString,
            cc_dot_cc__pb2.Chat.FromString,
            options, channel_credentials,
            insecure, call_credentials, compression, wait_for_ready, timeout, metadata)

    @staticmethod
    def List(request,
            target,
            options=(),
            channel_credentials=None,
            call_credentials=None,
            insecure=False,
            compression=None,
            wait_for_ready=None,
            timeout=None,
            metadata=None):
        return grpc.experimental.unary_unary(request, target, '/cc.ChatsAPI/List',
            cc_dot_cc__pb2.Empty.SerializeToString,
            cc_dot_cc__pb2.Chats.FromString,
            options, channel_credentials,
            insecure, call_credentials, compression, wait_for_ready, timeout, metadata)

    @staticmethod
    def Delete(request,
            target,
            options=(),
            channel_credentials=None,
            call_credentials=None,
            insecure=False,
            compression=None,
            wait_for_ready=None,
            timeout=None,
            metadata=None):
        return grpc.experimental.unary_unary(request, target, '/cc.ChatsAPI/Delete',
            cc_dot_cc__pb2.Chat.SerializeToString,
            cc_dot_cc__pb2.Chat.FromString,
            options, channel_credentials,
            insecure, call_credentials, compression, wait_for_ready, timeout, metadata)

    @staticmethod
    def SetBotState(request,
            target,
            options=(),
            channel_credentials=None,
            call_credentials=None,
            insecure=False,
            compression=None,
            wait_for_ready=None,
            timeout=None,
            metadata=None):
        return grpc.experimental.unary_unary(request, target, '/cc.ChatsAPI/SetBotState',
            cc_dot_cc__pb2.Chat.SerializeToString,
            cc_dot_cc__pb2.Chat.FromString,
            options, channel_credentials,
            insecure, call_credentials, compression, wait_for_ready, timeout, metadata)

    @staticmethod
    def GetBotState(request,
            target,
            options=(),
            channel_credentials=None,
            call_credentials=None,
            insecure=False,
            compression=None,
            wait_for_ready=None,
            timeout=None,
            metadata=None):
        return grpc.experimental.unary_unary(request, target, '/cc.ChatsAPI/GetBotState',
            cc_dot_cc__pb2.Chat.SerializeToString,
            cc_dot_cc__pb2.Chat.FromString,
            options, channel_credentials,
            insecure, call_credentials, compression, wait_for_ready, timeout, metadata)

    @staticmethod
    def ChangeDepartment(request,
            target,
            options=(),
            channel_credentials=None,
            call_credentials=None,
            insecure=False,
            compression=None,
            wait_for_ready=None,
            timeout=None,
            metadata=None):
        return grpc.experimental.unary_unary(request, target, '/cc.ChatsAPI/ChangeDepartment',
            cc_dot_cc__pb2.Chat.SerializeToString,
            cc_dot_cc__pb2.Chat.FromString,
            options, channel_credentials,
            insecure, call_credentials, compression, wait_for_ready, timeout, metadata)

    @staticmethod
    def ChangeGateway(request,
            target,
            options=(),
            channel_credentials=None,
            call_credentials=None,
            insecure=False,
            compression=None,
            wait_for_ready=None,
            timeout=None,
            metadata=None):
        return grpc.experimental.unary_unary(request, target, '/cc.ChatsAPI/ChangeGateway',
            cc_dot_cc__pb2.Chat.SerializeToString,
            cc_dot_cc__pb2.Chat.FromString,
            options, channel_credentials,
            insecure, call_credentials, compression, wait_for_ready, timeout, metadata)

    @staticmethod
    def ChangeStatus(request,
            target,
            options=(),
            channel_credentials=None,
            call_credentials=None,
            insecure=False,
            compression=None,
            wait_for_ready=None,
            timeout=None,
            metadata=None):
        return grpc.experimental.unary_unary(request, target, '/cc.ChatsAPI/ChangeStatus',
            cc_dot_cc__pb2.Chat.SerializeToString,
            cc_dot_cc__pb2.Chat.FromString,
            options, channel_credentials,
            insecure, call_credentials, compression, wait_for_ready, timeout, metadata)

    @staticmethod
    def MergeChats(request,
            target,
            options=(),
            channel_credentials=None,
            call_credentials=None,
            insecure=False,
            compression=None,
            wait_for_ready=None,
            timeout=None,
            metadata=None):
        return grpc.experimental.unary_unary(request, target, '/cc.ChatsAPI/MergeChats',
            cc_dot_cc__pb2.Merge.SerializeToString,
            cc_dot_cc__pb2.Chat.FromString,
            options, channel_credentials,
            insecure, call_credentials, compression, wait_for_ready, timeout, metadata)

    @staticmethod
    def SyncChats(request,
            target,
            options=(),
            channel_credentials=None,
            call_credentials=None,
            insecure=False,
            compression=None,
            wait_for_ready=None,
            timeout=None,
            metadata=None):
        return grpc.experimental.unary_unary(request, target, '/cc.ChatsAPI/SyncChats',
            cc_dot_cc__pb2.Empty.SerializeToString,
            cc_dot_cc__pb2.Empty.FromString,
            options, channel_credentials,
            insecure, call_credentials, compression, wait_for_ready, timeout, metadata)


class MessagesAPIStub(object):
    """Missing associated documentation comment in .proto file."""

    def __init__(self, channel):
        """Constructor.

        Args:
            channel: A grpc.Channel.
        """
        self.Get = channel.unary_unary(
                '/cc.MessagesAPI/Get',
                request_serializer=cc_dot_cc__pb2.Chat.SerializeToString,
                response_deserializer=cc_dot_cc__pb2.Messages.FromString,
                )
        self.Send = channel.unary_unary(
                '/cc.MessagesAPI/Send',
                request_serializer=cc_dot_cc__pb2.Message.SerializeToString,
                response_deserializer=cc_dot_cc__pb2.Message.FromString,
                )
        self.Update = channel.unary_unary(
                '/cc.MessagesAPI/Update',
                request_serializer=cc_dot_cc__pb2.Message.SerializeToString,
                response_deserializer=cc_dot_cc__pb2.Message.FromString,
                )
        self.Delete = channel.unary_unary(
                '/cc.MessagesAPI/Delete',
                request_serializer=cc_dot_cc__pb2.Message.SerializeToString,
                response_deserializer=cc_dot_cc__pb2.Message.FromString,
                )


class MessagesAPIServicer(object):
    """Missing associated documentation comment in .proto file."""

    def Get(self, request, context):
        """Missing associated documentation comment in .proto file."""
        context.set_code(grpc.StatusCode.UNIMPLEMENTED)
        context.set_details('Method not implemented!')
        raise NotImplementedError('Method not implemented!')

    def Send(self, request, context):
        """Missing associated documentation comment in .proto file."""
        context.set_code(grpc.StatusCode.UNIMPLEMENTED)
        context.set_details('Method not implemented!')
        raise NotImplementedError('Method not implemented!')

    def Update(self, request, context):
        """Missing associated documentation comment in .proto file."""
        context.set_code(grpc.StatusCode.UNIMPLEMENTED)
        context.set_details('Method not implemented!')
        raise NotImplementedError('Method not implemented!')

    def Delete(self, request, context):
        """Missing associated documentation comment in .proto file."""
        context.set_code(grpc.StatusCode.UNIMPLEMENTED)
        context.set_details('Method not implemented!')
        raise NotImplementedError('Method not implemented!')


def add_MessagesAPIServicer_to_server(servicer, server):
    rpc_method_handlers = {
            'Get': grpc.unary_unary_rpc_method_handler(
                    servicer.Get,
                    request_deserializer=cc_dot_cc__pb2.Chat.FromString,
                    response_serializer=cc_dot_cc__pb2.Messages.SerializeToString,
            ),
            'Send': grpc.unary_unary_rpc_method_handler(
                    servicer.Send,
                    request_deserializer=cc_dot_cc__pb2.Message.FromString,
                    response_serializer=cc_dot_cc__pb2.Message.SerializeToString,
            ),
            'Update': grpc.unary_unary_rpc_method_handler(
                    servicer.Update,
                    request_deserializer=cc_dot_cc__pb2.Message.FromString,
                    response_serializer=cc_dot_cc__pb2.Message.SerializeToString,
            ),
            'Delete': grpc.unary_unary_rpc_method_handler(
                    servicer.Delete,
                    request_deserializer=cc_dot_cc__pb2.Message.FromString,
                    response_serializer=cc_dot_cc__pb2.Message.SerializeToString,
            ),
    }
    generic_handler = grpc.method_handlers_generic_handler(
            'cc.MessagesAPI', rpc_method_handlers)
    server.add_generic_rpc_handlers((generic_handler,))


 # This class is part of an EXPERIMENTAL API.
class MessagesAPI(object):
    """Missing associated documentation comment in .proto file."""

    @staticmethod
    def Get(request,
            target,
            options=(),
            channel_credentials=None,
            call_credentials=None,
            insecure=False,
            compression=None,
            wait_for_ready=None,
            timeout=None,
            metadata=None):
        return grpc.experimental.unary_unary(request, target, '/cc.MessagesAPI/Get',
            cc_dot_cc__pb2.Chat.SerializeToString,
            cc_dot_cc__pb2.Messages.FromString,
            options, channel_credentials,
            insecure, call_credentials, compression, wait_for_ready, timeout, metadata)

    @staticmethod
    def Send(request,
            target,
            options=(),
            channel_credentials=None,
            call_credentials=None,
            insecure=False,
            compression=None,
            wait_for_ready=None,
            timeout=None,
            metadata=None):
        return grpc.experimental.unary_unary(request, target, '/cc.MessagesAPI/Send',
            cc_dot_cc__pb2.Message.SerializeToString,
            cc_dot_cc__pb2.Message.FromString,
            options, channel_credentials,
            insecure, call_credentials, compression, wait_for_ready, timeout, metadata)

    @staticmethod
    def Update(request,
            target,
            options=(),
            channel_credentials=None,
            call_credentials=None,
            insecure=False,
            compression=None,
            wait_for_ready=None,
            timeout=None,
            metadata=None):
        return grpc.experimental.unary_unary(request, target, '/cc.MessagesAPI/Update',
            cc_dot_cc__pb2.Message.SerializeToString,
            cc_dot_cc__pb2.Message.FromString,
            options, channel_credentials,
            insecure, call_credentials, compression, wait_for_ready, timeout, metadata)

    @staticmethod
    def Delete(request,
            target,
            options=(),
            channel_credentials=None,
            call_credentials=None,
            insecure=False,
            compression=None,
            wait_for_ready=None,
            timeout=None,
            metadata=None):
        return grpc.experimental.unary_unary(request, target, '/cc.MessagesAPI/Delete',
            cc_dot_cc__pb2.Message.SerializeToString,
            cc_dot_cc__pb2.Message.FromString,
            options, channel_credentials,
            insecure, call_credentials, compression, wait_for_ready, timeout, metadata)


class UsersAPIStub(object):
    """Missing associated documentation comment in .proto file."""

    def __init__(self, channel):
        """Constructor.

        Args:
            channel: A grpc.Channel.
        """
        self.Me = channel.unary_unary(
                '/cc.UsersAPI/Me',
                request_serializer=cc_dot_cc__pb2.Empty.SerializeToString,
                response_deserializer=cc_dot_cc__pb2.User.FromString,
                )
        self.FetchDefaults = channel.unary_unary(
                '/cc.UsersAPI/FetchDefaults',
                request_serializer=cc_dot_cc__pb2.FetchDefaultsRequest.SerializeToString,
                response_deserializer=cc_dot_cc__pb2.Defaults.FromString,
                )
        self.GetConfig = channel.unary_unary(
                '/cc.UsersAPI/GetConfig',
                request_serializer=cc_dot_cc__pb2.Empty.SerializeToString,
                response_deserializer=cc_dot_cc__pb2.Defaults.FromString,
                )
        self.SetConfig = channel.unary_unary(
                '/cc.UsersAPI/SetConfig',
                request_serializer=cc_dot_cc__pb2.Defaults.SerializeToString,
                response_deserializer=cc_dot_cc__pb2.Defaults.FromString,
                )
        self.Resolve = channel.unary_unary(
                '/cc.UsersAPI/Resolve',
                request_serializer=cc_dot_cc__pb2.Users.SerializeToString,
                response_deserializer=cc_dot_cc__pb2.Users.FromString,
                )
        self.GetMembers = channel.unary_unary(
                '/cc.UsersAPI/GetMembers',
                request_serializer=cc_dot_cc__pb2.Empty.SerializeToString,
                response_deserializer=cc_dot_cc__pb2.Users.FromString,
                )


class UsersAPIServicer(object):
    """Missing associated documentation comment in .proto file."""

    def Me(self, request, context):
        """Missing associated documentation comment in .proto file."""
        context.set_code(grpc.StatusCode.UNIMPLEMENTED)
        context.set_details('Method not implemented!')
        raise NotImplementedError('Method not implemented!')

    def FetchDefaults(self, request, context):
        """Missing associated documentation comment in .proto file."""
        context.set_code(grpc.StatusCode.UNIMPLEMENTED)
        context.set_details('Method not implemented!')
        raise NotImplementedError('Method not implemented!')

    def GetConfig(self, request, context):
        """Missing associated documentation comment in .proto file."""
        context.set_code(grpc.StatusCode.UNIMPLEMENTED)
        context.set_details('Method not implemented!')
        raise NotImplementedError('Method not implemented!')

    def SetConfig(self, request, context):
        """Missing associated documentation comment in .proto file."""
        context.set_code(grpc.StatusCode.UNIMPLEMENTED)
        context.set_details('Method not implemented!')
        raise NotImplementedError('Method not implemented!')

    def Resolve(self, request, context):
        """Resolves given Users data by their UUIDs
        And returns all accessible Users for Requestor
        """
        context.set_code(grpc.StatusCode.UNIMPLEMENTED)
        context.set_details('Method not implemented!')
        raise NotImplementedError('Method not implemented!')

    def GetMembers(self, request, context):
        """Missing associated documentation comment in .proto file."""
        context.set_code(grpc.StatusCode.UNIMPLEMENTED)
        context.set_details('Method not implemented!')
        raise NotImplementedError('Method not implemented!')


def add_UsersAPIServicer_to_server(servicer, server):
    rpc_method_handlers = {
            'Me': grpc.unary_unary_rpc_method_handler(
                    servicer.Me,
                    request_deserializer=cc_dot_cc__pb2.Empty.FromString,
                    response_serializer=cc_dot_cc__pb2.User.SerializeToString,
            ),
            'FetchDefaults': grpc.unary_unary_rpc_method_handler(
                    servicer.FetchDefaults,
                    request_deserializer=cc_dot_cc__pb2.FetchDefaultsRequest.FromString,
                    response_serializer=cc_dot_cc__pb2.Defaults.SerializeToString,
            ),
            'GetConfig': grpc.unary_unary_rpc_method_handler(
                    servicer.GetConfig,
                    request_deserializer=cc_dot_cc__pb2.Empty.FromString,
                    response_serializer=cc_dot_cc__pb2.Defaults.SerializeToString,
            ),
            'SetConfig': grpc.unary_unary_rpc_method_handler(
                    servicer.SetConfig,
                    request_deserializer=cc_dot_cc__pb2.Defaults.FromString,
                    response_serializer=cc_dot_cc__pb2.Defaults.SerializeToString,
            ),
            'Resolve': grpc.unary_unary_rpc_method_handler(
                    servicer.Resolve,
                    request_deserializer=cc_dot_cc__pb2.Users.FromString,
                    response_serializer=cc_dot_cc__pb2.Users.SerializeToString,
            ),
            'GetMembers': grpc.unary_unary_rpc_method_handler(
                    servicer.GetMembers,
                    request_deserializer=cc_dot_cc__pb2.Empty.FromString,
                    response_serializer=cc_dot_cc__pb2.Users.SerializeToString,
            ),
    }
    generic_handler = grpc.method_handlers_generic_handler(
            'cc.UsersAPI', rpc_method_handlers)
    server.add_generic_rpc_handlers((generic_handler,))


 # This class is part of an EXPERIMENTAL API.
class UsersAPI(object):
    """Missing associated documentation comment in .proto file."""

    @staticmethod
    def Me(request,
            target,
            options=(),
            channel_credentials=None,
            call_credentials=None,
            insecure=False,
            compression=None,
            wait_for_ready=None,
            timeout=None,
            metadata=None):
        return grpc.experimental.unary_unary(request, target, '/cc.UsersAPI/Me',
            cc_dot_cc__pb2.Empty.SerializeToString,
            cc_dot_cc__pb2.User.FromString,
            options, channel_credentials,
            insecure, call_credentials, compression, wait_for_ready, timeout, metadata)

    @staticmethod
    def FetchDefaults(request,
            target,
            options=(),
            channel_credentials=None,
            call_credentials=None,
            insecure=False,
            compression=None,
            wait_for_ready=None,
            timeout=None,
            metadata=None):
        return grpc.experimental.unary_unary(request, target, '/cc.UsersAPI/FetchDefaults',
            cc_dot_cc__pb2.FetchDefaultsRequest.SerializeToString,
            cc_dot_cc__pb2.Defaults.FromString,
            options, channel_credentials,
            insecure, call_credentials, compression, wait_for_ready, timeout, metadata)

    @staticmethod
    def GetConfig(request,
            target,
            options=(),
            channel_credentials=None,
            call_credentials=None,
            insecure=False,
            compression=None,
            wait_for_ready=None,
            timeout=None,
            metadata=None):
        return grpc.experimental.unary_unary(request, target, '/cc.UsersAPI/GetConfig',
            cc_dot_cc__pb2.Empty.SerializeToString,
            cc_dot_cc__pb2.Defaults.FromString,
            options, channel_credentials,
            insecure, call_credentials, compression, wait_for_ready, timeout, metadata)

    @staticmethod
    def SetConfig(request,
            target,
            options=(),
            channel_credentials=None,
            call_credentials=None,
            insecure=False,
            compression=None,
            wait_for_ready=None,
            timeout=None,
            metadata=None):
        return grpc.experimental.unary_unary(request, target, '/cc.UsersAPI/SetConfig',
            cc_dot_cc__pb2.Defaults.SerializeToString,
            cc_dot_cc__pb2.Defaults.FromString,
            options, channel_credentials,
            insecure, call_credentials, compression, wait_for_ready, timeout, metadata)

    @staticmethod
    def Resolve(request,
            target,
            options=(),
            channel_credentials=None,
            call_credentials=None,
            insecure=False,
            compression=None,
            wait_for_ready=None,
            timeout=None,
            metadata=None):
        return grpc.experimental.unary_unary(request, target, '/cc.UsersAPI/Resolve',
            cc_dot_cc__pb2.Users.SerializeToString,
            cc_dot_cc__pb2.Users.FromString,
            options, channel_credentials,
            insecure, call_credentials, compression, wait_for_ready, timeout, metadata)

    @staticmethod
    def GetMembers(request,
            target,
            options=(),
            channel_credentials=None,
            call_credentials=None,
            insecure=False,
            compression=None,
            wait_for_ready=None,
            timeout=None,
            metadata=None):
        return grpc.experimental.unary_unary(request, target, '/cc.UsersAPI/GetMembers',
            cc_dot_cc__pb2.Empty.SerializeToString,
            cc_dot_cc__pb2.Users.FromString,
            options, channel_credentials,
            insecure, call_credentials, compression, wait_for_ready, timeout, metadata)


class StreamServiceStub(object):
    """Missing associated documentation comment in .proto file."""

    def __init__(self, channel):
        """Constructor.

        Args:
            channel: A grpc.Channel.
        """
        self.Stream = channel.unary_stream(
                '/cc.StreamService/Stream',
                request_serializer=cc_dot_cc__pb2.StreamRequest.SerializeToString,
                response_deserializer=cc_dot_cc__pb2.Event.FromString,
                )


class StreamServiceServicer(object):
    """Missing associated documentation comment in .proto file."""

    def Stream(self, request, context):
        """Missing associated documentation comment in .proto file."""
        context.set_code(grpc.StatusCode.UNIMPLEMENTED)
        context.set_details('Method not implemented!')
        raise NotImplementedError('Method not implemented!')


def add_StreamServiceServicer_to_server(servicer, server):
    rpc_method_handlers = {
            'Stream': grpc.unary_stream_rpc_method_handler(
                    servicer.Stream,
                    request_deserializer=cc_dot_cc__pb2.StreamRequest.FromString,
                    response_serializer=cc_dot_cc__pb2.Event.SerializeToString,
            ),
    }
    generic_handler = grpc.method_handlers_generic_handler(
            'cc.StreamService', rpc_method_handlers)
    server.add_generic_rpc_handlers((generic_handler,))


 # This class is part of an EXPERIMENTAL API.
class StreamService(object):
    """Missing associated documentation comment in .proto file."""

    @staticmethod
    def Stream(request,
            target,
            options=(),
            channel_credentials=None,
            call_credentials=None,
            insecure=False,
            compression=None,
            wait_for_ready=None,
            timeout=None,
            metadata=None):
        return grpc.experimental.unary_stream(request, target, '/cc.StreamService/Stream',
            cc_dot_cc__pb2.StreamRequest.SerializeToString,
            cc_dot_cc__pb2.Event.FromString,
            options, channel_credentials,
            insecure, call_credentials, compression, wait_for_ready, timeout, metadata)
