package server

import (
	"context"

	"github.com/bufbuild/connect-go"
	"github.com/slntopp/core-chatting/cc"
	"github.com/slntopp/core-chatting/cc/ccconnect"
	client "github.com/slntopp/core-chatting/default"
)

type ChatsAPIServer struct {
	AccountKey any
	Client     client.DefaultClient

	ccconnect.UnimplementedChatsAPIHandler
}

func (srv *ChatsAPIServer) List(ctx context.Context, _ *connect.Request[cc.Empty]) (*connect.Response[cc.Chats], error) {

	requestor, err := Requestor(srv.AccountKey, ctx)
	if err != nil {
		return nil, err
	}

	chats, err := srv.Client.GetChats(requestor)
	if err != nil {
		return nil, connect.NewError(connect.CodeInternal, err)
	}

	return connect.NewResponse(&cc.Chats{Chats: chats}), nil
}

func (srv *ChatsAPIServer) Create(ctx context.Context, req *connect.Request[cc.Chat]) (*connect.Response[cc.Chat], error) {
	requestor, err := Requestor(srv.AccountKey, ctx)
	if err != nil {
		return nil, err
	}

	chat, err := srv.Client.CreateChat(requestor, req.Msg)
	return connect.NewResponse(chat), err
}
