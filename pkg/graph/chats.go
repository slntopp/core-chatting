package graph

import (
	"context"

	"github.com/arangodb/go-driver"
	"github.com/slntopp/core-chatting/cc"
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
	return nil, nil
}
func (c *ChatsController) Update(ctx context.Context, chat *cc.Chat) (*cc.Chat, error) {
	return nil, nil
}
func (c *ChatsController) List(ctx context.Context) ([]*cc.Chat, error) {
	return nil, nil
}
func (c *ChatsController) Delete(ctx context.Context, chat *cc.Chat) (*cc.Chat, error) {
	return nil, nil
}
