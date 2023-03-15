package schema

import (
	"errors"

	"github.com/slntopp/core-chatting/cc"
)

type Chat = *cc.Chat
type Message = *cc.Message

type StorageClient interface {
	// Must retrieve Chats available to the given Account and set the ROLE
	GetChats(account string) ([]Chat, error)
	// Must retrieve requested Chat and set the ROLE or return error if NoAccess
	GetChat(account, chat string) (Chat, error)
	// Must store or update the Chat
	SaveChat(chat *Chat) error
	// Must delete or mark the Chat and all it's Messages as deleted
	DeleteChat(chat Chat) error

	// Must retrieve the Messages sent to the given Chat
	GetMessages(chat Chat) ([]Message, error)
	// Must retrieve the Message
	GetMessage(chat Chat, msg string) (Message, error)
	// Must save and bind the message to the given Chat or Update it if GetUuid is not empty
	SaveMessage(chat Chat, msg Message) error
	// Must delete or mark the Message as deleted
	DeleteMessage(chat Chat, msg Message) error

	mustEmbedUnimplementedStorageClient()
}

type UnimplementedStorageClient struct {
}

func (UnimplementedStorageClient) GetChats(account string) ([]Chat, error) {
	return nil, errors.New("method GetChats unimplemented")
}
func (UnimplementedStorageClient) GetChat(account, chat string) (Chat, error) {
	return nil, errors.New("method GetChat unimplemented")
}
func (UnimplementedStorageClient) SaveChat(chat Chat) error {
	return errors.New("method SaveChat unimplemented")
}
func (UnimplementedStorageClient) DeleteChat(chat Chat) error {
	return errors.New("method DeleteChat unimplemented")
}
func (UnimplementedStorageClient) GetMessages(chat Chat) ([]Message, error) {
	return nil, errors.New("method GetMessages unimplemented")
}
func (UnimplementedStorageClient) GetMessage(chat Chat, msg string) (Message, error) {
	return nil, errors.New("method GetMessage unimplemented")
}
func (UnimplementedStorageClient) SaveMessage(chat Chat, msg Message) error {
	return errors.New("method SaveMessage unimplemented")
}
func (UnimplementedStorageClient) DeleteMessage(chat Chat, msg Message) error {
	return errors.New("method DeleteMessage unimplemented")
}
func (UnimplementedStorageClient) mustEmbedUnimplementedStorageClient() {}
