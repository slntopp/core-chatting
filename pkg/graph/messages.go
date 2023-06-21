package graph

import (
	"context"

	"github.com/arangodb/go-driver"
	"github.com/slntopp/core-chatting/cc"
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
	return nil, nil
}

func (c *MessagesController) Update(ctx context.Context, msg *cc.Message) (*cc.Message, error) {
	return nil, nil
}

func (c *MessagesController) Delete(ctx context.Context, msg *cc.Message) (*cc.Message, error) {
	return nil, nil
}
