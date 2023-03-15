package im

import (
	"github.com/slntopp/core-chatting/cc"
	"github.com/slntopp/core-chatting/schema"
)

type Chat struct {
	Chat     *cc.Chat
	Messages []*cc.Messages
}

type InMemoryStorage struct {
	Chats map[string]Chat

	schema.UnimplementedStorageClient
}
