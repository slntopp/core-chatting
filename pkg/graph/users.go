package graph

import (
	"context"

	"github.com/arangodb/go-driver"
	"github.com/slntopp/core-chatting/cc"
	"go.uber.org/zap"
)

type UsersController struct {
	log *zap.Logger

	colname string

	db  driver.Database
	col driver.Collection
}

func NewUsersController(logger *zap.Logger, db driver.Database, colname string) *UsersController {
	log := logger.Named("UsersController")
	log.Debug("Creating Users Controller", zap.String("collection", colname))

	ctx := context.TODO()
	col := GetEnsureCollection(log, ctx, db, colname)

	return &UsersController{
		log:     log,
		colname: colname,
		db:      db,
		col:     col,
	}
}

func (c *UsersController) Resolve(ctx context.Context, uuids []string) ([]*cc.User, error) {
	log := c.log.Named("Resolve")
	log.Debug("Request received", zap.Strings("uuids", uuids))

	var users []*cc.User
	// TODO: Check Access i.e. search on Graph
	for _, uuid := range uuids {
		user := cc.User{}
		_, err := c.col.ReadDocument(ctx, uuid, &user)
		if err != nil {
			return nil, err
		}

		user.Uuid = uuid
		users = append(users, &user)
	}

	log.Debug("Found users", zap.Any("users", users))

	return users, nil
}

const getMembers = `
FOR a in @@accounts
  RETURN MERGE(a, {
  	uuid: a._key
  })
`

func (c *UsersController) GetMembers(ctx context.Context) ([]*cc.User, error) {
	log := c.log.Named("GetMembers")
	log.Debug("Request received")

	cur, err := c.db.Query(ctx, getMembers, map[string]interface{}{
		"@accounts": c.colname,
	})
	if err != nil {
		return nil, err
	}

	var members []*cc.User

	for cur.HasMore() {
		var member cc.User

		_, err := cur.ReadDocument(ctx, &member)
		if err != nil {
			return nil, err
		}

		members = append(members, &member)
	}

	return members, nil
}
