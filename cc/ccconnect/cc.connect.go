// Code generated by protoc-gen-connect-go. DO NOT EDIT.
//
// Source: cc/cc.proto

package ccconnect

import (
	context "context"
	errors "errors"
	connect_go "github.com/bufbuild/connect-go"
	cc "github.com/slntopp/core-chatting/cc"
	http "net/http"
	strings "strings"
)

// This is a compile-time assertion to ensure that this generated file and the connect package are
// compatible. If you get a compiler error that this constant is not defined, this code was
// generated with a version of connect newer than the one compiled into your binary. You can fix the
// problem by either regenerating this code with an older version of connect or updating the connect
// version compiled into your binary.
const _ = connect_go.IsAtLeastVersion0_1_0

const (
	// ChatsAPIName is the fully-qualified name of the ChatsAPI service.
	ChatsAPIName = "cc.ChatsAPI"
	// MessagesAPIName is the fully-qualified name of the MessagesAPI service.
	MessagesAPIName = "cc.MessagesAPI"
)

// These constants are the fully-qualified names of the RPCs defined in this package. They're
// exposed at runtime as Spec.Procedure and as the final two segments of the HTTP route.
//
// Note that these are different from the fully-qualified method names used by
// google.golang.org/protobuf/reflect/protoreflect. To convert from these constants to
// reflection-formatted method names, remove the leading slash and convert the remaining slash to a
// period.
const (
	// ChatsAPICreateProcedure is the fully-qualified name of the ChatsAPI's Create RPC.
	ChatsAPICreateProcedure = "/cc.ChatsAPI/Create"
	// ChatsAPIUpdateProcedure is the fully-qualified name of the ChatsAPI's Update RPC.
	ChatsAPIUpdateProcedure = "/cc.ChatsAPI/Update"
	// ChatsAPIListProcedure is the fully-qualified name of the ChatsAPI's List RPC.
	ChatsAPIListProcedure = "/cc.ChatsAPI/List"
	// ChatsAPIDeleteProcedure is the fully-qualified name of the ChatsAPI's Delete RPC.
	ChatsAPIDeleteProcedure = "/cc.ChatsAPI/Delete"
	// MessagesAPIGetProcedure is the fully-qualified name of the MessagesAPI's Get RPC.
	MessagesAPIGetProcedure = "/cc.MessagesAPI/Get"
	// MessagesAPISendProcedure is the fully-qualified name of the MessagesAPI's Send RPC.
	MessagesAPISendProcedure = "/cc.MessagesAPI/Send"
	// MessagesAPIUpdateProcedure is the fully-qualified name of the MessagesAPI's Update RPC.
	MessagesAPIUpdateProcedure = "/cc.MessagesAPI/Update"
	// MessagesAPIDeleteProcedure is the fully-qualified name of the MessagesAPI's Delete RPC.
	MessagesAPIDeleteProcedure = "/cc.MessagesAPI/Delete"
)

// ChatsAPIClient is a client for the cc.ChatsAPI service.
type ChatsAPIClient interface {
	Create(context.Context, *connect_go.Request[cc.Chat]) (*connect_go.Response[cc.Chat], error)
	Update(context.Context, *connect_go.Request[cc.Chat]) (*connect_go.Response[cc.Chat], error)
	List(context.Context, *connect_go.Request[cc.Empty]) (*connect_go.Response[cc.Chats], error)
	Delete(context.Context, *connect_go.Request[cc.Chat]) (*connect_go.Response[cc.Chat], error)
}

// NewChatsAPIClient constructs a client for the cc.ChatsAPI service. By default, it uses the
// Connect protocol with the binary Protobuf Codec, asks for gzipped responses, and sends
// uncompressed requests. To use the gRPC or gRPC-Web protocols, supply the connect.WithGRPC() or
// connect.WithGRPCWeb() options.
//
// The URL supplied here should be the base URL for the Connect or gRPC server (for example,
// http://api.acme.com or https://acme.com/grpc).
func NewChatsAPIClient(httpClient connect_go.HTTPClient, baseURL string, opts ...connect_go.ClientOption) ChatsAPIClient {
	baseURL = strings.TrimRight(baseURL, "/")
	return &chatsAPIClient{
		create: connect_go.NewClient[cc.Chat, cc.Chat](
			httpClient,
			baseURL+ChatsAPICreateProcedure,
			opts...,
		),
		update: connect_go.NewClient[cc.Chat, cc.Chat](
			httpClient,
			baseURL+ChatsAPIUpdateProcedure,
			opts...,
		),
		list: connect_go.NewClient[cc.Empty, cc.Chats](
			httpClient,
			baseURL+ChatsAPIListProcedure,
			opts...,
		),
		delete: connect_go.NewClient[cc.Chat, cc.Chat](
			httpClient,
			baseURL+ChatsAPIDeleteProcedure,
			opts...,
		),
	}
}

// chatsAPIClient implements ChatsAPIClient.
type chatsAPIClient struct {
	create *connect_go.Client[cc.Chat, cc.Chat]
	update *connect_go.Client[cc.Chat, cc.Chat]
	list   *connect_go.Client[cc.Empty, cc.Chats]
	delete *connect_go.Client[cc.Chat, cc.Chat]
}

// Create calls cc.ChatsAPI.Create.
func (c *chatsAPIClient) Create(ctx context.Context, req *connect_go.Request[cc.Chat]) (*connect_go.Response[cc.Chat], error) {
	return c.create.CallUnary(ctx, req)
}

// Update calls cc.ChatsAPI.Update.
func (c *chatsAPIClient) Update(ctx context.Context, req *connect_go.Request[cc.Chat]) (*connect_go.Response[cc.Chat], error) {
	return c.update.CallUnary(ctx, req)
}

// List calls cc.ChatsAPI.List.
func (c *chatsAPIClient) List(ctx context.Context, req *connect_go.Request[cc.Empty]) (*connect_go.Response[cc.Chats], error) {
	return c.list.CallUnary(ctx, req)
}

// Delete calls cc.ChatsAPI.Delete.
func (c *chatsAPIClient) Delete(ctx context.Context, req *connect_go.Request[cc.Chat]) (*connect_go.Response[cc.Chat], error) {
	return c.delete.CallUnary(ctx, req)
}

// ChatsAPIHandler is an implementation of the cc.ChatsAPI service.
type ChatsAPIHandler interface {
	Create(context.Context, *connect_go.Request[cc.Chat]) (*connect_go.Response[cc.Chat], error)
	Update(context.Context, *connect_go.Request[cc.Chat]) (*connect_go.Response[cc.Chat], error)
	List(context.Context, *connect_go.Request[cc.Empty]) (*connect_go.Response[cc.Chats], error)
	Delete(context.Context, *connect_go.Request[cc.Chat]) (*connect_go.Response[cc.Chat], error)
}

// NewChatsAPIHandler builds an HTTP handler from the service implementation. It returns the path on
// which to mount the handler and the handler itself.
//
// By default, handlers support the Connect, gRPC, and gRPC-Web protocols with the binary Protobuf
// and JSON codecs. They also support gzip compression.
func NewChatsAPIHandler(svc ChatsAPIHandler, opts ...connect_go.HandlerOption) (string, http.Handler) {
	mux := http.NewServeMux()
	mux.Handle(ChatsAPICreateProcedure, connect_go.NewUnaryHandler(
		ChatsAPICreateProcedure,
		svc.Create,
		opts...,
	))
	mux.Handle(ChatsAPIUpdateProcedure, connect_go.NewUnaryHandler(
		ChatsAPIUpdateProcedure,
		svc.Update,
		opts...,
	))
	mux.Handle(ChatsAPIListProcedure, connect_go.NewUnaryHandler(
		ChatsAPIListProcedure,
		svc.List,
		opts...,
	))
	mux.Handle(ChatsAPIDeleteProcedure, connect_go.NewUnaryHandler(
		ChatsAPIDeleteProcedure,
		svc.Delete,
		opts...,
	))
	return "/cc.ChatsAPI/", mux
}

// UnimplementedChatsAPIHandler returns CodeUnimplemented from all methods.
type UnimplementedChatsAPIHandler struct{}

func (UnimplementedChatsAPIHandler) Create(context.Context, *connect_go.Request[cc.Chat]) (*connect_go.Response[cc.Chat], error) {
	return nil, connect_go.NewError(connect_go.CodeUnimplemented, errors.New("cc.ChatsAPI.Create is not implemented"))
}

func (UnimplementedChatsAPIHandler) Update(context.Context, *connect_go.Request[cc.Chat]) (*connect_go.Response[cc.Chat], error) {
	return nil, connect_go.NewError(connect_go.CodeUnimplemented, errors.New("cc.ChatsAPI.Update is not implemented"))
}

func (UnimplementedChatsAPIHandler) List(context.Context, *connect_go.Request[cc.Empty]) (*connect_go.Response[cc.Chats], error) {
	return nil, connect_go.NewError(connect_go.CodeUnimplemented, errors.New("cc.ChatsAPI.List is not implemented"))
}

func (UnimplementedChatsAPIHandler) Delete(context.Context, *connect_go.Request[cc.Chat]) (*connect_go.Response[cc.Chat], error) {
	return nil, connect_go.NewError(connect_go.CodeUnimplemented, errors.New("cc.ChatsAPI.Delete is not implemented"))
}

// MessagesAPIClient is a client for the cc.MessagesAPI service.
type MessagesAPIClient interface {
	Get(context.Context, *connect_go.Request[cc.Chat]) (*connect_go.Response[cc.Messages], error)
	Send(context.Context, *connect_go.Request[cc.Message]) (*connect_go.Response[cc.Message], error)
	Update(context.Context, *connect_go.Request[cc.Message]) (*connect_go.Response[cc.Message], error)
	Delete(context.Context, *connect_go.Request[cc.Message]) (*connect_go.Response[cc.Message], error)
}

// NewMessagesAPIClient constructs a client for the cc.MessagesAPI service. By default, it uses the
// Connect protocol with the binary Protobuf Codec, asks for gzipped responses, and sends
// uncompressed requests. To use the gRPC or gRPC-Web protocols, supply the connect.WithGRPC() or
// connect.WithGRPCWeb() options.
//
// The URL supplied here should be the base URL for the Connect or gRPC server (for example,
// http://api.acme.com or https://acme.com/grpc).
func NewMessagesAPIClient(httpClient connect_go.HTTPClient, baseURL string, opts ...connect_go.ClientOption) MessagesAPIClient {
	baseURL = strings.TrimRight(baseURL, "/")
	return &messagesAPIClient{
		get: connect_go.NewClient[cc.Chat, cc.Messages](
			httpClient,
			baseURL+MessagesAPIGetProcedure,
			opts...,
		),
		send: connect_go.NewClient[cc.Message, cc.Message](
			httpClient,
			baseURL+MessagesAPISendProcedure,
			opts...,
		),
		update: connect_go.NewClient[cc.Message, cc.Message](
			httpClient,
			baseURL+MessagesAPIUpdateProcedure,
			opts...,
		),
		delete: connect_go.NewClient[cc.Message, cc.Message](
			httpClient,
			baseURL+MessagesAPIDeleteProcedure,
			opts...,
		),
	}
}

// messagesAPIClient implements MessagesAPIClient.
type messagesAPIClient struct {
	get    *connect_go.Client[cc.Chat, cc.Messages]
	send   *connect_go.Client[cc.Message, cc.Message]
	update *connect_go.Client[cc.Message, cc.Message]
	delete *connect_go.Client[cc.Message, cc.Message]
}

// Get calls cc.MessagesAPI.Get.
func (c *messagesAPIClient) Get(ctx context.Context, req *connect_go.Request[cc.Chat]) (*connect_go.Response[cc.Messages], error) {
	return c.get.CallUnary(ctx, req)
}

// Send calls cc.MessagesAPI.Send.
func (c *messagesAPIClient) Send(ctx context.Context, req *connect_go.Request[cc.Message]) (*connect_go.Response[cc.Message], error) {
	return c.send.CallUnary(ctx, req)
}

// Update calls cc.MessagesAPI.Update.
func (c *messagesAPIClient) Update(ctx context.Context, req *connect_go.Request[cc.Message]) (*connect_go.Response[cc.Message], error) {
	return c.update.CallUnary(ctx, req)
}

// Delete calls cc.MessagesAPI.Delete.
func (c *messagesAPIClient) Delete(ctx context.Context, req *connect_go.Request[cc.Message]) (*connect_go.Response[cc.Message], error) {
	return c.delete.CallUnary(ctx, req)
}

// MessagesAPIHandler is an implementation of the cc.MessagesAPI service.
type MessagesAPIHandler interface {
	Get(context.Context, *connect_go.Request[cc.Chat]) (*connect_go.Response[cc.Messages], error)
	Send(context.Context, *connect_go.Request[cc.Message]) (*connect_go.Response[cc.Message], error)
	Update(context.Context, *connect_go.Request[cc.Message]) (*connect_go.Response[cc.Message], error)
	Delete(context.Context, *connect_go.Request[cc.Message]) (*connect_go.Response[cc.Message], error)
}

// NewMessagesAPIHandler builds an HTTP handler from the service implementation. It returns the path
// on which to mount the handler and the handler itself.
//
// By default, handlers support the Connect, gRPC, and gRPC-Web protocols with the binary Protobuf
// and JSON codecs. They also support gzip compression.
func NewMessagesAPIHandler(svc MessagesAPIHandler, opts ...connect_go.HandlerOption) (string, http.Handler) {
	mux := http.NewServeMux()
	mux.Handle(MessagesAPIGetProcedure, connect_go.NewUnaryHandler(
		MessagesAPIGetProcedure,
		svc.Get,
		opts...,
	))
	mux.Handle(MessagesAPISendProcedure, connect_go.NewUnaryHandler(
		MessagesAPISendProcedure,
		svc.Send,
		opts...,
	))
	mux.Handle(MessagesAPIUpdateProcedure, connect_go.NewUnaryHandler(
		MessagesAPIUpdateProcedure,
		svc.Update,
		opts...,
	))
	mux.Handle(MessagesAPIDeleteProcedure, connect_go.NewUnaryHandler(
		MessagesAPIDeleteProcedure,
		svc.Delete,
		opts...,
	))
	return "/cc.MessagesAPI/", mux
}

// UnimplementedMessagesAPIHandler returns CodeUnimplemented from all methods.
type UnimplementedMessagesAPIHandler struct{}

func (UnimplementedMessagesAPIHandler) Get(context.Context, *connect_go.Request[cc.Chat]) (*connect_go.Response[cc.Messages], error) {
	return nil, connect_go.NewError(connect_go.CodeUnimplemented, errors.New("cc.MessagesAPI.Get is not implemented"))
}

func (UnimplementedMessagesAPIHandler) Send(context.Context, *connect_go.Request[cc.Message]) (*connect_go.Response[cc.Message], error) {
	return nil, connect_go.NewError(connect_go.CodeUnimplemented, errors.New("cc.MessagesAPI.Send is not implemented"))
}

func (UnimplementedMessagesAPIHandler) Update(context.Context, *connect_go.Request[cc.Message]) (*connect_go.Response[cc.Message], error) {
	return nil, connect_go.NewError(connect_go.CodeUnimplemented, errors.New("cc.MessagesAPI.Update is not implemented"))
}

func (UnimplementedMessagesAPIHandler) Delete(context.Context, *connect_go.Request[cc.Message]) (*connect_go.Response[cc.Message], error) {
	return nil, connect_go.NewError(connect_go.CodeUnimplemented, errors.New("cc.MessagesAPI.Delete is not implemented"))
}
