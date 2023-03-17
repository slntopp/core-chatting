package im

import (
	"errors"
	"fmt"

	"github.com/bufbuild/connect-go"
	"github.com/slntopp/core-chatting/cc"
	"github.com/slntopp/core-chatting/storage"
)

type Chat struct {
	Chat     *cc.Chat
	Messages []*cc.Message
}

type InMemoryStorage struct {
	Chats map[string]Chat

	storage.UnimplementedStorageClient
}

func New() *InMemoryStorage {
	return &InMemoryStorage{
		Chats: make(map[string]Chat),
	}
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

func (s *InMemoryStorage) GetChat(acc string, chat_id string) (res *cc.Chat, err error) {
	chat, ok := s.Chats[chat_id]
	if !ok {
		return nil, connect.NewError(connect.CodeNotFound, nil)
	}

	if acc == "admin" || In(acc, chat.Chat.Admins) {
		chat.Chat.Role = cc.Role_ADMIN
	} else if In(acc, chat.Chat.Users) {
		chat.Chat.Role = cc.Role_USER
	} else {
		return nil, connect.NewError(connect.CodePermissionDenied, nil)
	}

	return chat.Chat, nil
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
			return connect.NewError(connect.CodeNotFound, errors.New("not found"))
		} else if !In(acc, chat.Chat.Admins) {
			return connect.NewError(connect.CodePermissionDenied, errors.New("not enough access rights"))
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

func (s *InMemoryStorage) DeleteChat(chat *cc.Chat) error {
	delete(s.Chats, chat.GetUuid())
	return nil
}

func (s *InMemoryStorage) GetMessages(chat *cc.Chat) ([]*cc.Message, error) {
	c, ok := s.Chats[chat.GetUuid()]
	if !ok {
		return nil, connect.NewError(connect.CodeNotFound, errors.New("not found"))
	}

	return c.Messages, nil
}

func (s *InMemoryStorage) SaveMessage(chat *cc.Chat, msg *cc.Message) error {
	c, ok := s.Chats[chat.GetUuid()]
	if !ok {
		return connect.NewError(connect.CodeNotFound, errors.New("not found"))
	}

	c.Messages = append(c.Messages, msg)
	s.Chats[chat.GetUuid()] = c
	return nil
}
