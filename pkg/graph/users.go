package graph

import (
	"context"

	"github.com/arangodb/go-driver"
	"github.com/slntopp/core-chatting/cc"
	"go.uber.org/zap"
)

type UsersController struct {
	log *zap.Logger

	usersColname   string
	devicesColname string

	db         driver.Database
	usersCol   driver.Collection
	devicesCol driver.Collection
}

func NewUsersController(logger *zap.Logger, db driver.Database, usersColname, devicesColname string) *UsersController {
	log := logger.Named("UsersController")
	log.Debug("Creating Users Controller", zap.String("users collection", usersColname), zap.String("devices collection", devicesColname))

	ctx := context.TODO()
	usersCol := GetEnsureCollection(log, ctx, db, usersColname)
	devicesCol := GetEnsureCollection(log, ctx, db, usersColname)

	return &UsersController{
		log:            log,
		usersColname:   usersColname,
		devicesColname: devicesColname,
		db:             db,
		usersCol:       usersCol,
		devicesCol:     devicesCol,
	}
}

func (c *UsersController) Resolve(ctx context.Context, uuids []string) ([]*cc.User, error) {
	log := c.log.Named("Resolve")
	log.Debug("Request received", zap.Strings("uuids", uuids))

	var users []*cc.User
	// TODO: Check Access i.e. search on Graph
	for _, uuid := range uuids {
		user := cc.User{}
		_, err := c.usersCol.ReadDocument(ctx, uuid, &user)
		if err != nil {
			return nil, err
		}

		user.Uuid = uuid
		users = append(users, &user)
	}

	log.Debug("Found users", zap.Any("users", users))

	return users, nil
}

const getCollection = `
FOR c in @@col
  RETURN MERGE(c, {
  	uuid: c._key
  })
`

func (c *UsersController) GetMembers(ctx context.Context) ([]*cc.User, error) {
	log := c.log.Named("GetMembers")
	log.Debug("Request received")

	cur, err := c.db.Query(ctx, getCollection, map[string]interface{}{
		"@col": c.usersColname,
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

	log.Debug("Len", zap.Int("len", len(members)))

	return members, nil
}

func (c *UsersController) GetDevices(ctx context.Context) ([]*cc.Device, error) {
	log := c.log.Named("GetMembers")
	log.Debug("Request received")

	cur, err := c.db.Query(ctx, getCollection, map[string]interface{}{
		"@col": c.devicesCol,
	})
	if err != nil {
		return nil, err
	}

	var members []*cc.Device

	for cur.HasMore() {
		var member cc.Device

		_, err := cur.ReadDocument(ctx, &member)
		if err != nil {
			return nil, err
		}

		members = append(members, &member)
	}

	log.Debug("Len", zap.Int("len", len(members)))

	return members, nil
}
