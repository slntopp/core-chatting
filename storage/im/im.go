package im

import (
	"errors"
	"fmt"

	"github.com/slntopp/core-chatting/cc"
	"github.com/slntopp/core-chatting/schema"
)

type Chat struct {
	Chat     *cc.Chat
	Messages []*cc.Message
}

type InMemoryStorage struct {
	Chats map[string]Chat

	schema.UnimplementedStorageClient
}

var counter = 0

func In[T comparable](o T, a []T) bool {
	for _, el := range a {
		if el == o {
			return true
		}
	}
	return false
}

func (s *InMemoryStorage) GetChats(acc string) (res []*cc.Chat, err error) {
	for _, chat := range s.Chats {
		if acc == "admin" || In(acc, chat.Chat.Admins) {
			chat.Chat.Role = cc.Role_ADMIN
			res = append(res, chat.Chat)
		} else if In(acc, chat.Chat.Users) {
			chat.Chat.Role = cc.Role_USER
			res = append(res, chat.Chat)
		}
	}

	return
}

func (s *InMemoryStorage) SaveChat(acc string, chat *cc.Chat) error {
	var msgs []*cc.Message
	if chat.GetUuid() != "" {
		if chat, ok := s.Chats[chat.Uuid]; !ok {
			return errors.New("not found")
		} else if !In(acc, chat.Chat.Admins) {
			return errors.New("not enough access rights")
		} else {
			msgs = chat.Messages
		}
	} else {
		chat.Uuid = fmt.Sprintf("chat-%d", counter)
		counter++
	}

	s.Chats[chat.Uuid] = Chat{
		Chat:     chat,
		Messages: msgs,
	}

	return nil
}
