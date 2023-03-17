package main

import (
	"context"
	"errors"
	"fmt"
	"net/http"

	"github.com/bufbuild/connect-go"
	"github.com/slntopp/core-chatting/cc/ccconnect"
	"github.com/slntopp/core-chatting/server"
	"github.com/slntopp/core-chatting/storage/im"
	"golang.org/x/net/http2"
	"golang.org/x/net/http2/h2c"
)

type AccountKey string

const AccountType = AccountKey("account-id")

func NewAuthInterceptor() connect.UnaryInterceptorFunc {
	return func(next connect.UnaryFunc) connect.UnaryFunc {
		return func(ctx context.Context, req connect.AnyRequest) (connect.AnyResponse, error) {
			fmt.Println("Headers", req.Header())
			acc := req.Header().Get("X-Account-ID")
			fmt.Println("Account", acc)
			if acc == "" {
				return nil, connect.NewError(connect.CodeUnauthenticated, errors.New("account type is not provided"))
			}

			ctx = context.WithValue(ctx, AccountType, acc)
			return next(ctx, req)
		}
	}
}

func main() {
	im := im.New()
	cht_srv := server.NewChatsAPIServer(AccountType, im)
	msg_srv := server.NewMessagesAPIServer(AccountType, im)

	mux := http.NewServeMux()
	path, handler := ccconnect.NewChatsAPIHandler(cht_srv, connect.WithInterceptors(NewAuthInterceptor()))
	mux.Handle(path, handler)
	path, handler = ccconnect.NewMessagesAPIHandler(msg_srv, connect.WithInterceptors(NewAuthInterceptor()))
	mux.Handle(path, handler)

	http.ListenAndServe(":8000", h2c.NewHandler(mux, &http2.Server{}))
}
