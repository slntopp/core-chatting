package storage

import (
	"errors"

	c "github.com/bufbuild/connect-go"
	"github.com/slntopp/core-chatting/cc"
)

type Chat = *cc.Chat
type Message = *cc.Message

type StorageClient interface {
	// Must retrieve Chats available to the given Account and set the ROLE
	// Error has to comply with connect-go error with status code
	GetChats(account string) ([]Chat, error)
	// Must retrieve requested Chat and set the ROLE or return error if NoAccess
	// Error has to comply with connect-go error with status code
	GetChat(account, chat string) (Chat, error)
	// Must store or update the Chat if account has enough access rights
	// Error has to comply with connect-go error with status code
	SaveChat(account string, chat Chat) error
	// Must delete or mark the Chat and all it's Messages as deleted
	// Error has to comply with connect-go error with status code
	DeleteChat(chat Chat) error

	// Must retrieve the Messages sent to the given Chat
	// Error has to comply with connect-go error with status code
	GetMessages(chat Chat) ([]Message, error)
	// Must retrieve the Message
	// Error has to comply with connect-go error with status code
	GetMessage(chat Chat, msg string) (Message, error)
	// Must save and bind the message to the given Chat or Update it if GetUuid is not empty
	// Error has to comply with connect-go error with status code
	SaveMessage(chat Chat, msg Message) error
	// Must delete or mark the Message as deleted
	// Error has to comply with connect-go error with status code
	DeleteMessage(chat Chat, msg Message) error

	mustEmbedUnimplementedStorageClient()
}

type UnimplementedStorageClient struct {
}

func (UnimplementedStorageClient) GetChats(account string) ([]Chat, error) {
	return nil, c.NewError(c.CodeUnimplemented, errors.New("method GetChats unimplemented"))
}
func (UnimplementedStorageClient) GetChat(account, chat string) (Chat, error) {
	return nil, c.NewError(c.CodeUnimplemented, errors.New("method GetChat unimplemented"))
}
func (UnimplementedStorageClient) SaveChat(chat Chat) error {
	return c.NewError(c.CodeUnimplemented, errors.New("method SaveChat unimplemented"))
}
func (UnimplementedStorageClient) DeleteChat(chat Chat) error {
	return c.NewError(c.CodeUnimplemented, errors.New("method DeleteChat unimplemented"))
}
func (UnimplementedStorageClient) GetMessages(chat Chat) ([]Message, error) {
	return nil, c.NewError(c.CodeUnimplemented, errors.New("method GetMessages unimplemented"))
}
func (UnimplementedStorageClient) GetMessage(chat Chat, msg string) (Message, error) {
	return nil, c.NewError(c.CodeUnimplemented, errors.New("method GetMessage unimplemented"))
}
func (UnimplementedStorageClient) SaveMessage(chat Chat, msg Message) error {
	return c.NewError(c.CodeUnimplemented, errors.New("method SaveMessage unimplemented"))
}
func (UnimplementedStorageClient) DeleteMessage(chat Chat, msg Message) error {
	return c.NewError(c.CodeUnimplemented, errors.New("method DeleteMessage unimplemented"))
}
func (UnimplementedStorageClient) mustEmbedUnimplementedStorageClient() {}
