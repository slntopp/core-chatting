package main

import (
	"net/http"

	"github.com/slntopp/core-chatting/cc/ccconnect"
	"github.com/slntopp/core-chatting/server"
	"golang.org/x/net/http2"
	"golang.org/x/net/http2/h2c"
)

func main() {
	srv := &server.ChatsAPIServer{}

	mux := http.NewServeMux()
	path, handler := ccconnect.NewChatsAPIHandler(srv)
	mux.Handle(path, handler)

	http.ListenAndServe(":8000", h2c.NewHandler(mux, &http2.Server{}))
}
