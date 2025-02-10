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
LET members = ( FOR a in @@accounts
  %s
  %s
  RETURN MERGE(a, {
  	uuid: a._key
  })
)
LET total = LENGTH(members)
RETURN { pool: @count>0 ? SLICE(members, @offset, @count) : members, total: total }
`

func (c *UsersController) GetMembers(ctx context.Context, req *cc.GetMembersRequest) ([]*cc.User, int64, error) {
	log := c.log.Named("GetMembers")
	log.Debug("Request received")

	vars := map[string]interface{}{
		"@accounts": c.colname,
	}
	filters := ""
	if req.GetFilters() != nil {
		for key, value := range req.GetFilters() {
			if key == "uuid" || key == "uuids" {
				values := value.GetListValue().AsSlice()
				if len(values) == 0 {
					continue
				}
				filters += fmt.Sprintf(` FILTER a._key in @%s`, key)
				vars[key] = values
			} else if key == "search_param" {
				filters += fmt.Sprintf(` FILTER a._key LIKE "%s"`,
					"%"+value.GetStringValue()+"%")
			} else if key == "exclude_uuids" {
				values := value.GetListValue().AsSlice()
				if len(values) == 0 {
					continue
				}
				filters += fmt.Sprintf(` FILTER a._key NOT IN @%s`, key)
				vars[key] = values
			} else {
				values := value.GetListValue().AsSlice()
				if len(values) == 0 {
					continue
				}
				filters += fmt.Sprintf(` FILTER a["%s"] in @%s`, key, key)
				vars[key] = values
			}
		}
	}
	if len(req.GetUuids()) > 0 {
		key := "includedUuids"
		filters += fmt.Sprintf(` FILTER a._key in @%s`, key)
		vars[key] = req.GetUuids()
	}

	sorts := ""
	if req.Field != nil && req.Sort != nil {
		subQuery := ` SORT c.%s %s`
		field, sort := req.GetField(), req.GetSort()
		sorts += fmt.Sprintf(subQuery, field, sort)
	}

	if req.Page != nil && req.Limit != nil && req.GetLimit() != 0 {
		limit, page := req.GetLimit(), req.GetPage()
		offset := (page - 1) * limit
		vars["offset"] = offset
		vars["count"] = limit
	} else {
		vars["offset"] = 0
		vars["count"] = 0
	}

	query := fmt.Sprintf(getMembers, filters, sorts)
	cur, err := c.db.Query(ctx, query, vars)
	if err != nil {
		return nil, 0, err
	}

	resp := struct {
		Pool  []*cc.User `json:"pool"`
		Total int64      `json:"total"`
	}{}

	if cur.HasMore() {
		_, err := cur.ReadDocument(ctx, &resp)
		if err != nil {
			return nil, 0, err
		}
	} else {
		return nil, 0, fmt.Errorf("not found or internal")
	}

	return resp.Pool, resp.Total, nil
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
