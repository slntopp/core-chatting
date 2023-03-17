package server

import (
	"context"

	c "github.com/bufbuild/connect-go"
	"github.com/slntopp/core-chatting/cc"
	"github.com/slntopp/core-chatting/cc/ccconnect"
	"github.com/slntopp/core-chatting/storage"
)

type MessagesAPIServer struct {
	accKey any
	s      storage.StorageClient

	ccconnect.UnimplementedMessagesAPIHandler
}

func NewMessagesAPIServer(key any, st storage.StorageClient) *MessagesAPIServer {
	return &MessagesAPIServer{
		accKey: key, s: st,
	}
}

func (srv *MessagesAPIServer) Get(ctx context.Context, req *c.Request[cc.Chat]) (*c.Response[cc.Messages], error) {

	requestor, err := Requestor(srv.accKey, ctx)
	if err != nil {
		return nil, err
	}

	chat, err := srv.s.GetChat(requestor, req.Msg.GetUuid())
	if err != nil {
		return nil, err
	}

	messages, err := srv.s.GetMessages(chat)
	return c.NewResponse(&cc.Messages{
		Messages: messages,
	}), err
}

func (srv *MessagesAPIServer) Send(ctx context.Context, req *c.Request[cc.Message]) (*c.Response[cc.Message], error) {
	requestor, err := Requestor(srv.accKey, ctx)
	if err != nil {
		return nil, err
	}

	if req.Msg.Chat == nil {
		return nil, NewError(c.CodeInvalidArgument, "chat uuid is not provided")
	}

	chat, err := srv.s.GetChat(requestor, req.Msg.GetChat())
	if err != nil {
		return nil, err
	}

	msg := req.Msg

	if chat.Role != cc.Role_ADMIN || msg.Sender == "" {
		msg.Sender = requestor
	}

	msg.Uuid = ""

	if chat.Role != cc.Role_ADMIN {
		if msg.Kind == cc.Kind_ADMIN_ONLY {
			return nil, NewError(c.CodePermissionDenied, "")
		}
		if len(msg.Gateways) > 0 {
			return nil, NewError(c.CodePermissionDenied, "")
		}
	}

	err = srv.s.SaveMessage(chat, msg)
	return c.NewResponse(msg), err
}
