package server

import (
	"context"

	c "github.com/bufbuild/connect-go"

	"github.com/slntopp/core-chatting/cc"
	"github.com/slntopp/core-chatting/cc/ccconnect"
	"github.com/slntopp/core-chatting/storage"
)

type ChatsAPIServer struct {
	accKey any
	s      storage.StorageClient

	ccconnect.UnimplementedChatsAPIHandler
}

func NewChatsAPIServer(key any, st storage.StorageClient) *ChatsAPIServer {
	return &ChatsAPIServer{
		accKey: key, s: st,
	}
}

func (srv *ChatsAPIServer) List(ctx context.Context, _ *c.Request[cc.Empty]) (*c.Response[cc.Chats], error) {

	requestor, err := Requestor(srv.accKey, ctx)
	if err != nil {
		return nil, err
	}

	chats, err := srv.s.GetChats(requestor)
	if err != nil {
		return nil, err
	}

	return c.NewResponse(&cc.Chats{Chats: chats}), nil
}

func (srv *ChatsAPIServer) Create(ctx context.Context, req *c.Request[cc.Chat]) (*c.Response[cc.Chat], error) {
	requestor, err := Requestor(srv.accKey, ctx)
	if err != nil {
		return nil, err
	}

	req.Msg.Uuid = ""

	err = srv.s.SaveChat(requestor, req.Msg)
	return c.NewResponse(req.Msg), err
}

func (srv *ChatsAPIServer) Update(ctx context.Context, req *c.Request[cc.Chat]) (*c.Response[cc.Chat], error) {
	requestor, err := Requestor(srv.accKey, ctx)
	if err != nil {
		return nil, err
	}

	if req.Msg.GetUuid() == "" {
		return nil, NewError(c.CodeInvalidArgument, "must provide uuid to update")
	}

	err = srv.s.SaveChat(requestor, req.Msg)
	return c.NewResponse(req.Msg), err
}

func (srv *ChatsAPIServer) Delete(ctx context.Context, req *c.Request[cc.Chat]) (*c.Response[cc.Chat], error) {
	requestor, err := Requestor(srv.accKey, ctx)
	if err != nil {
		return nil, err
	}

	if req.Msg.GetUuid() == "" {
		return nil, NewError(c.CodeInvalidArgument, "must provide uuid to delete")
	}

	chat, err := srv.s.GetChat(requestor, req.Msg.Uuid)
	if err != nil {
		return nil, err
	}

	if chat.Role != cc.Role_ADMIN && chat.Role != cc.Role_OWNER {
		return nil, NewError(c.CodePermissionDenied, "only owners and admins can perform this operation")
	}

	return c.NewResponse(chat), srv.s.DeleteChat(chat)
}
