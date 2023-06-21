package graph

import (
	"context"

	"github.com/slntopp/core-chatting/cc"

	"github.com/arangodb/go-driver"
	"go.uber.org/zap"
)

type ChatsController struct {
	log *zap.Logger

	db  driver.Database
	col driver.Collection
}

func NewChatsController(logger *zap.Logger, db driver.Database) *ChatsController {
	log := logger.Named("ChatsController")
	log.Debug("Creating Chats Controller")

	ctx := context.TODO()

	col := GetEnsureCollection(log, ctx, db, CHATS_COLLECTION)

	return &ChatsController{
		log: log,
		db:  db,
		col: col,
	}
}

func (c *ChatsController) Create(ctx context.Context, chat *cc.Chat) (*cc.Chat, error) {
	log := c.log.Named("Create")
	log.Debug("Req received")

	document, err := c.col.CreateDocument(ctx, chat)
	if err != nil {
		return nil, err
	}

	chat.Uuid = document.Key

	_, err = c.col.UpdateDocument(ctx, chat.GetUuid(), chat)

	if err != nil {
		return nil, err
	}

	return chat, nil
}

func (c *ChatsController) Update(ctx context.Context, chat *cc.Chat) (*cc.Chat, error) {
	log := c.log.Named("Update")
	log.Debug("Req received")

	_, err := c.col.UpdateDocument(ctx, chat.GetUuid(), chat)

	if err != nil {
		return nil, err
	}

	return chat, nil
}

func (c *ChatsController) Get(ctx context.Context, uuid string) (*cc.Chat, error) {
	log := c.log.Named("Get")
	log.Debug("Req received")

	var chat cc.Chat

	_, err := c.col.ReadDocument(ctx, uuid, &chat)

	if err != nil {
		return nil, err
	}

	return &chat, nil
}

const getChatsQuery = `
FOR c in Chats
	RETURN c
`

func (c *ChatsController) List(ctx context.Context) ([]*cc.Chat, error) {
	log := c.log.Named("List")
	log.Debug("Req received")

	cur, err := c.db.Query(ctx, getChatsQuery, map[string]interface{}{})

	if err != nil {
		return nil, err
	}

	var chats []*cc.Chat

	for cur.HasMore() {
		var chat cc.Chat

		_, err := cur.ReadDocument(ctx, &chat)
		if err != nil {
			return nil, err
		}

		chats = append(chats, &chat)
	}

	return chats, nil
}

const deleteChatMessages = `
FOR m in Messages
	FILTER m.chat == @chat
	REMOVE m IN Messages
`

func (c *ChatsController) Delete(ctx context.Context, chat *cc.Chat) (*cc.Chat, error) {
	log := c.log.Named("Delete")
	log.Debug("Req received")

	_, err := c.col.RemoveDocument(ctx, chat.GetUuid())

	if err != nil {
		return nil, err
	}

	_, err = c.db.Query(ctx, deleteChatMessages, map[string]interface{}{
		"chat": chat.GetUuid(),
	})

	if err != nil {
		return nil, err
	}

	return chat, nil
}

const getChatMessages = `
FOR m in Messages
	FILTER m.chat == @chat
	RETURN m
`

func (c *ChatsController) GetMessages(ctx context.Context, chat *cc.Chat) ([]*cc.Message, error) {
	log := c.log.Named("List")
	log.Debug("Req received")

	cur, err := c.db.Query(ctx, getChatMessages, map[string]interface{}{
		"chat": chat.GetUuid(),
	})

	if err != nil {
		return nil, err
	}

	var messages []*cc.Message

	for cur.HasMore() {
		var message cc.Message

		_, err := cur.ReadDocument(ctx, &message)
		if err != nil {
			return nil, err
		}

		messages = append(messages, &message)
	}

	return messages, nil
}
