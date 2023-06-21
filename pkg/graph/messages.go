package graph

import (
	"context"

	"github.com/slntopp/core-chatting/cc"

	"github.com/arangodb/go-driver"
	"go.uber.org/zap"
)

type MessagesController struct {
	log *zap.Logger

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

func (c *MessagesController) Update(ctx context.Context, msg *cc.Message) (*cc.Message, error) {
	log := c.log.Named("Update")
	log.Debug("Req received")

	_, err := c.col.UpdateDocument(ctx, msg.GetUuid(), msg)

	if err != nil {
		return nil, err
	}

	return msg, nil
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
