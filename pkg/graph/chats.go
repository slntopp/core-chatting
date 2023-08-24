package graph

import (
	"context"
	"errors"
	"fmt"
	"time"

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

	chat.Created = time.Now().UnixMilli()
	chat.Status = cc.Status_NEW

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

const getChatQuery = `
LET chat = Document(@chat)

LET role = (
 @requestor in chat.admins ? 3 : (
  chat.owner == @requestor ? 2 : (
   @requestor in chat.users ? 1 : 0
  )
 )
)

RETURN MERGE(chat, {
  role: role
})
`

func (c *ChatsController) Get(ctx context.Context, uuid, requestor string) (*cc.Chat, error) {
	log := c.log.Named("Get")
	log.Debug("Req received")

	cur, err := c.db.Query(ctx, getChatQuery, map[string]interface{}{
		"chat":      driver.NewDocumentID(CHATS_COLLECTION, uuid),
		"requestor": requestor,
	})

	if err != nil {
		return nil, err
	}
	defer cur.Close()

	var chat cc.Chat

	_, err = cur.ReadDocument(ctx, &chat)
	if err != nil {
		return nil, err
	}

	return &chat, nil
}

const listChatsQuery = `
FOR c in @@chats
FILTER @requestor in c.admins || @requestor in c.users
	LET role = (
     @requestor in c.admins ? 3 : (
      c.owner == @requestor ? 2 : (
       @requestor in c.users ? 1 : 0
      )
     )
    )
	LET message = LAST(FOR m in @@messages FILTER m.chat == c._key SORT m.sent ASC RETURN m)
	LET unread = LENGTH(
		FOR m in @@messages 
			FILTER m.chat == c._key
			FILTER @requestor not in m.readers
			FILTER m.sender != @requestor
			RETURN m
		)
	RETURN MERGE(c, {
	  role: role,
      meta: {
		unread: unread,
		last_message: message,
		data: c.meta.data
	  }
	})
`

func (c *ChatsController) List(ctx context.Context, requestor string) ([]*cc.Chat, error) {
	log := c.log.Named("List")
	log.Debug("Req received")

	cur, err := c.db.Query(ctx, listChatsQuery, map[string]interface{}{
		"@chats":    CHATS_COLLECTION,
		"@messages": MESSAGES_COLLECTION,
		"requestor": requestor,
	})

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
FOR m in @@messages
	FILTER m.chat == @chat
	REMOVE m IN @@messages
`

func (c *ChatsController) Delete(ctx context.Context, chat *cc.Chat) (*cc.Chat, error) {
	log := c.log.Named("Delete")
	log.Debug("Req received")

	var deletedChat cc.Chat

	_, err := c.col.ReadDocument(ctx, chat.GetUuid(), &deletedChat)

	if err != nil {
		return nil, err
	}

	_, err = c.col.RemoveDocument(ctx, chat.GetUuid())

	if err != nil {
		return nil, err
	}

	_, err = c.db.Query(ctx, deleteChatMessages, map[string]interface{}{
		"chat":      chat.GetUuid(),
		"@messages": MESSAGES_COLLECTION,
	})

	if err != nil {
		return nil, err
	}

	return &deletedChat, nil
}

const getChatMessages = `
FOR m in @@messages
	FILTER m.chat == @chat
	%s
	RETURN m
`

func (c *ChatsController) GetMessages(ctx context.Context, chat *cc.Chat, is_admin bool) ([]*cc.Message, error) {
	log := c.log.Named("List")
	log.Debug("Request received", zap.String("chat", chat.GetUuid()), zap.Bool("is_admin", is_admin))

	bind_vars := map[string]interface{}{
		"chat":      chat.GetUuid(),
		"@messages": MESSAGES_COLLECTION,
	}

	extra_filter := ""
	if !is_admin {
		extra_filter = "FILTER m.kind != @admin_only\nFILTER !m.under_review"
		bind_vars["admin_only"] = cc.Kind_ADMIN_ONLY
	}

	query := fmt.Sprintf(getChatMessages, extra_filter)

	cur, err := c.db.Query(ctx, query, bind_vars)

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

const getChatByGateway = `
FOR c in @@chats
    FILTER c.meta.data[@gateway] == @id
    RETURN c
`

func (c *ChatsController) GetByGateway(ctx context.Context, req *cc.GetawayRequest) (*cc.Chat, error) {
	log := c.log.Named("GetByGateway")
	log.Debug("Req received")

	queryContext := driver.WithQueryCount(ctx)

	cur, err := c.db.Query(queryContext, getChatByGateway, map[string]interface{}{
		"@chats":  CHATS_COLLECTION,
		"id":      req.GatewayId,
		"gateway": req.Gateway,
	})

	if err != nil {
		return nil, err
	}
	defer cur.Close()

	if cur.Count() == 0 {
		return nil, errors.New("no chat")
	}

	var chat cc.Chat

	_, err = cur.ReadDocument(ctx, &chat)
	if err != nil {
		return nil, err
	}

	return &chat, nil
}
