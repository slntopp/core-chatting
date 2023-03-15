package client

import (
	"fmt"

	"github.com/slntopp/core-chatting/cc"
	"github.com/slntopp/core-chatting/schema"
)

type Chat = *cc.Chat
type Message = *cc.Message

type DefaultClient struct {
	Storage schema.StorageClient
}

func (p *DefaultClient) GetChats(account string) ([]Chat, error) {
	return p.Storage.GetChats(account)
}

func (p *DefaultClient) GetChat(account, chat_id string) (Chat, error) {
	return p.Storage.GetChat(account, chat_id)
}

func (p *DefaultClient) CreateChat(account string, chat Chat) (Chat, error) {
	if chat.GetUuid() != "" {
		return chat, fmt.Errorf("id must be unset")
	}

	err := p.Storage.SaveChat(&chat)
	return chat, err
}

func (p *DefaultClient) UpdateChat(account string, chat Chat) (Chat, error) {
	if chat.GetUuid() == "" {
		return chat, fmt.Errorf("id is unset")
	}

	curr, err := p.Storage.GetChat(account, chat.GetUuid())
	if err != nil {
		return curr, err
	}

	if curr.GetRole() != cc.Role_ADMIN {
		return curr, fmt.Errorf("must be admin")
	}

	err = p.Storage.SaveChat(&chat)
	return chat, err
}

func (p *DefaultClient) SendMessage(account, chat_id string, msg Message) error {
	chat, err := p.Storage.GetChat(account, chat_id)
	if err != nil {
		return err
	}

	if msg.GetKind() == cc.Kind_ADMIN_ONLY && chat.GetRole() != cc.Role_ADMIN {
		return fmt.Errorf("admin action is not available for chat %s", chat_id)
	}

	return p.Storage.SaveMessage(chat, msg)
}

func (p *DefaultClient) GetMessages(account, chat_id string) (res []schema.Message, err error) {
	chat, err := p.Storage.GetChat(account, chat_id)
	if err != nil {
		return
	}

	pool, err := p.Storage.GetMessages(chat)
	if err != nil {
		return
	}

	for _, msg := range pool {
		if msg.GetKind() == cc.Kind_ADMIN_ONLY && chat.GetRole() != cc.Role_ADMIN {
			continue
		}
		res = append(res, msg)
	}

	return
}

func (p *DefaultClient) DeleteMessage(account, chat_id, msg_id string) (err error) {
	chat, err := p.Storage.GetChat(account, chat_id)
	if err != nil {
		return
	}

	msg, err := p.Storage.GetMessage(chat, msg_id)
	if err != nil {
		return
	}

	if msg.GetKind() == cc.Kind_ADMIN_ONLY && chat.GetRole() != cc.Role_ADMIN {
		return fmt.Errorf("admin action is not available for chat %s", chat_id)
	}

	return p.Storage.DeleteMessage(chat, msg)
}
