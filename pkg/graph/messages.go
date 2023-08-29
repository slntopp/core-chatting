package graph

import (
	"context"
	"sync"

	"github.com/slntopp/core-chatting/cc"

	"github.com/arangodb/go-driver"
	"go.uber.org/zap"
)

type MessagesController struct {
	log *zap.Logger
	sync.Mutex

	db  driver.Database
	col driver.Collection
}

func NewMessagesController(logger *zap.Logger, db driver.Database) *MessagesController {
	log := logger.Named("MessagesController")
	log.Debug("Creating Messages Controller")

	ctx := context.TODO()

	col := GetEnsureCollection(log, ctx, db, MESSAGES_COLLECTION)

	return &MessagesController{
		log: log,
		db:  db,
		col: col,
	}
}

func (c *MessagesController) Send(ctx context.Context, msg *cc.Message) (*cc.Message, error) {
	log := c.log.Named("Create")
	log.Debug("Req received")

	document, err := c.col.CreateDocument(ctx, msg)
	if err != nil {
		return nil, err
	}

	msg.Uuid = document.Key

	_, err = c.col.UpdateDocument(ctx, msg.GetUuid(), msg)

	if err != nil {
		return nil, err
	}

	return msg, nil
}

const updateMessageQuery = `
UPDATE @key WITH {
    kind: @kind,
    content: @content,
    edited: @edited,
    under_review: @under_review,
	readers : @readers,
	meta: @meta
} IN @@messages
`

func (c *MessagesController) Update(ctx context.Context, msg *cc.Message) (*cc.Message, error) {
	log := c.log.Named("Update")
	log.Debug("Req received")

	params := map[string]interface{}{
		"kind":         msg.GetKind(),
		"content":      msg.GetContent(),
		"edited":       msg.GetEdited(),
		"under_review": msg.GetUnderReview(),
		"key":          msg.GetUuid(),
		"readers":      msg.GetReaders(),
		"meta":         msg.GetMeta(),
		"@messages":    MESSAGES_COLLECTION,
	}

	_, err := c.db.Query(ctx, updateMessageQuery, params)

	if err != nil {
		return nil, err
	}

	return msg, nil
}

const readMessageQuery = `
LET message = Document(@message)
LET new = UNIQUE(PUSH(message.readers, @reader))
UPDATE message with {readers: new} in @@messages RETURN NEW
`

func (c *MessagesController) Read(ctx context.Context, msg *cc.Message, reader string) (*cc.Message, error) {
	log := c.log.Named("Read")
	log.Debug("Req received")

	c.Lock()
	defer c.Unlock()
	cur, err := c.db.Query(ctx, readMessageQuery, map[string]interface{}{
		"message":   driver.NewDocumentID(MESSAGES_COLLECTION, msg.GetUuid()),
		"reader":    reader,
		"@messages": MESSAGES_COLLECTION,
	})

	if err != nil {
		return nil, err
	}

	var newMessage cc.Message

	_, err = cur.ReadDocument(ctx, &newMessage)
	if err != nil {
		return nil, err
	}

	return &newMessage, nil
}

func (c *MessagesController) Delete(ctx context.Context, msg *cc.Message) (*cc.Message, error) {
	log := c.log.Named("Delete")
	log.Debug("Req received")

	_, err := c.col.RemoveDocument(ctx, msg.GetUuid())

	if err != nil {
		return nil, err
	}

	return msg, nil
}

func (c *MessagesController) Get(ctx context.Context, uuid string) (*cc.Message, error) {
	log := c.log.Named("Get")
	log.Debug("Req received")

	var msg cc.Message

	_, err := c.col.ReadDocument(ctx, uuid, &msg)

	if err != nil {
		return nil, err
	}

	return &msg, nil
}
