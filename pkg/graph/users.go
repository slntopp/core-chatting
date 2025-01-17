package graph

import (
	"context"
	"fmt"

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
  %s
  RETURN MERGE(a, {
  	uuid: a._key
  })
`

func (c *UsersController) GetMembers(ctx context.Context, uuids []string) ([]*cc.User, error) {
	log := c.log.Named("GetMembers")
	log.Debug("Request received")

	vars := map[string]interface{}{
		"@accounts": c.colname,
	}
	filters := ""
	if len(uuids) > 0 {
		filters += " FILTER a._key IN @uuids"
		vars["uuids"] = uuids
	}

	query := fmt.Sprintf(getMembers, filters)
	cur, err := c.db.Query(ctx, query, vars)
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

	log.Debug("Len", zap.Int("len", len(members)))

	return members, nil
}

const removeCommandsQuery = `
UPDATE DOCUMENT(@key) WITH { cc_commands: null } IN @@collection 
`

const updateCommandsQuery = `
UPDATE DOCUMENT(@key) WITH { cc_commands: @commands } IN @@collection
`

func (c *UsersController) UpdateCommands(ctx context.Context, uuid string, commands map[string]string) error {
	log := c.log.Named("GetCommands")
	log.Debug("Request received")

	_, err := c.db.Query(ctx, removeCommandsQuery, map[string]interface{}{
		"@collection": c.colname,
		"key":         driver.NewDocumentID(c.colname, uuid),
	})
	if err != nil {
		log.Error("Failed to clear commands", zap.Error(err))
		return err
	}

	_, err = c.db.Query(ctx, updateCommandsQuery, map[string]interface{}{
		"@collection": c.colname,
		"key":         driver.NewDocumentID(c.colname, uuid),
		"commands":    commands,
	})

	return nil
}
