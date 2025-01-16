package graph

import (
	"context"
	"fmt"
	"github.com/slntopp/nocloud/pkg/nocloud/schema"
	"slices"
	"strconv"
	"time"

	"google.golang.org/protobuf/types/known/structpb"

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

	if chat.Responsible != nil {
		if !slices.Contains(chat.GetAdmins(), chat.GetResponsible()) {
			chat.Admins = append(chat.GetAdmins(), chat.GetResponsible())
		}
	}

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

	if chat.Responsible != nil {
		if !slices.Contains(chat.GetAdmins(), chat.GetResponsible()) {
			chat.Admins = append(chat.GetAdmins(), chat.GetResponsible())
		}
	}

	_, err := c.col.UpdateDocument(ctx, chat.GetUuid(), chat)

	if err != nil {
		log.Error("Failed to update chat", zap.Error(err))
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
LET chats = ( FOR c in @@chats
FILTER @requestor in c.admins || @requestor in c.users || @requestor == @root_account
	LET role = (
     @requestor in c.admins or @requestor == @root_account ? 3 : (
      c.owner == @requestor ? 2 : (
       @requestor in c.users ? 1 : 0
      )
     )
    )
    %s
	LET messages = (
         FOR m in @@messages 
         FILTER m.chat == c._key 
		 FILTER m.kind != @admin_only || role == 3
         FILTER !m.under_review || role == 3
         SORT m.sent ASC RETURN m
    )
	LET first_message = FIRST(messages)
	LET last_message = LAST(messages)
	LET unread = c.status == @closed_status ? 0 : LENGTH(
		FOR m in @@messages 
			FILTER m.chat == c._key
			FILTER @requestor not in m.readers
			FILTER m.sender != @requestor
			RETURN m
		)
    %s
	RETURN MERGE(c, {
	  role: role,
      meta: {
		unread: unread,
		last_message: last_message,
		first_message: first_message,
		data: c.meta.data
	  }
	}))
LET total = COUNT(FOR c in @@chats %s RETURN c)
RETURN { pool: chats, total: total }
`

func (c *ChatsController) List(ctx context.Context, requester string, req *cc.ListChatsRequest) ([]*cc.Chat, int64, error) {
	log := c.log.Named("List")
	log.Debug("Req received")

	vars := map[string]any{
		"@chats":        CHATS_COLLECTION,
		"@messages":     MESSAGES_COLLECTION,
		"requestor":     requester,
		"root_account":  schema.ROOT_ACCOUNT_KEY,
		"closed_status": cc.Status_CLOSE,
		"admin_only":    cc.Kind_ADMIN_ONLY,
	}
	filters := ""
	if req.GetFilters() != nil {
		for key, value := range req.GetFilters() {
			if key == "created" {
				values := value.GetStructValue().AsMap()
				if val, ok := values["from"]; ok {
					from := val.(float64)
					filters += fmt.Sprintf(` FILTER c["%s"] >= %f`, key, from)
				}

				if val, ok := values["to"]; ok {
					to := val.(float64)
					filters += fmt.Sprintf(` FILTER c["%s"] <= %f`, key, to)
				}
			} else if key == "admins" {
				values := value.GetListValue().AsSlice()
				if len(values) == 0 {
					continue
				}
				filters += fmt.Sprintf(` FILTER CONTAINS_ARRAY(c.admins, @%s)`, key)
				vars[key] = values
			} else if key == "status" {
				values := value.GetListValue().AsSlice()
				if len(values) == 0 {
					continue
				}
				filters += fmt.Sprintf(` FILTER c.status in @%s`, key)
				vars[key] = values
			} else if key == "department" {
				values := value.GetListValue().AsSlice()
				if len(values) == 0 {
					continue
				}
				filters += fmt.Sprintf(` FILTER c.department in @%s`, key)
				vars[key] = values
			} else if key == "responsible" {
				values := value.GetListValue().AsSlice()
				if len(values) == 0 {
					continue
				}
				filters += fmt.Sprintf(` FILTER c.responsible in @%s`, key)
				vars[key] = values
			} else if key == "account" {
				values := value.GetStringValue()
				if len(values) == 0 {
					continue
				}
				filters += fmt.Sprintf(` FILTER @%s in c.admins || @%s in c.users || @%s == c.owner || @%s == c.responsible`, key, key, key, key)
				vars[key] = values
			} else if key == "search_param" {
				filters += fmt.Sprintf(` FILTER c._key LIKE "%s" || LOWER(c.topic) LIKE LOWER("%s")`,
					"%"+value.GetStringValue()+"%", "%"+value.GetStringValue()+"%")
			} else {
				values := value.GetListValue().AsSlice()
				if len(values) == 0 {
					continue
				}
				filters += fmt.Sprintf(` FILTER c["%s"] in @%s`, key, key)
				vars[key] = values
			}
		}
	}

	sorts := ""
	if req.Field != nil && req.Sort != nil {
		subQuery := ` SORT c.%s %s`
		field, sort := req.GetField(), req.GetSort()

		if field == "sent" || field == "updated" {
			sorts += fmt.Sprintf(` SORT %s %s, c.created %s`, "last_message.sent", sort, sort)
		} else if field == "unread" {
			sorts += fmt.Sprintf(` SORT %s %s, c.created DESC`, "unread", sort)
		} else {
			sorts += fmt.Sprintf(subQuery, field, sort)
		}
	}

	if req.Page != nil && req.Limit != nil {
		if req.GetLimit() != 0 {
			limit, page := req.GetLimit(), req.GetPage()
			offset := (page - 1) * limit

			filters += ` LIMIT @offset, @count`
			vars["offset"] = offset
			vars["count"] = limit
		}
	}

	query := fmt.Sprintf(listChatsQuery, filters, sorts, filters)
	cur, err := c.db.Query(ctx, query, vars)
	if err != nil {
		return nil, 0, err
	}

	resp := struct {
		Pool  []*cc.Chat `json:"pool"`
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

const countChats = `
FOR c in @@chats
FILTER @requestor in c.admins || @requestor in c.users || @requestor == @root_account
COLLECT status = c.status WITH COUNT INTO times
RETURN { [status]: times }
`

func (c *ChatsController) Count(ctx context.Context, requester string) (map[int32]int64, error) {
	log := c.log.Named("Count")
	log.Debug("Req received")

	vars := map[string]any{
		"@chats":       CHATS_COLLECTION,
		"requestor":    requester,
		"root_account": schema.ROOT_ACCOUNT_KEY,
	}
	cur, err := c.db.Query(ctx, countChats, vars)
	if err != nil {
		return nil, err
	}

	resp := make(map[string]any)
	if cur.HasMore() {
		_, err := cur.ReadDocument(ctx, &resp)
		if err != nil {
			return nil, err
		}
	} else {
		return nil, fmt.Errorf("not found or internal")
	}

	res := make(map[int32]int64)
	for key, val := range resp {
		status, _ := strconv.Atoi(key)
		res[int32(status)] = int64(val.(float64))
	}
	return res, nil
}

const deleteChatMessages = `
FOR m in @@messages
	FILTER m.chat == @chat
	LET attachments = m.attachments == null ? [] : m.attachments 
	FOR a in attachments	
		REMOVE a in @@attachments
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
		"chat":         chat.GetUuid(),
		"@messages":    MESSAGES_COLLECTION,
		"@attachments": ATTACHMENTS_COLLECTION,
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
	SORT m.sent ASC
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

const getByGatewayQuery = `
FOR c in @@chats
	FILTER c.meta.data["@gate"] == @gate_id
	RETURN c
`

func (c *ChatsController) GetByGateway(ctx context.Context, gate string, gateId *structpb.Value) (*cc.Chat, error) {
	value := gateId.GetNumberValue()
	stringValue := gateId.GetStringValue()

	var queryId interface{}
	if stringValue == "" {
		queryId = value
	} else {
		queryId = stringValue
	}

	cur, err := c.db.Query(ctx, getByGatewayQuery, map[string]interface{}{
		"@chats":  CHATS_COLLECTION,
		"gate":    gate,
		"gate_id": queryId,
	})

	if err != nil {
		return nil, err
	}

	var gateChat cc.Chat

	_, err = cur.ReadDocument(ctx, &gateChat)
	if err != nil {
		return nil, err
	}

	return &gateChat, nil
}

const deleteGateways = `
FOR c in @@chats
	FILTER c.meta.data[@gate] == @gate_id
	UPDATE c with {meta: {data: { @gate : null }}} in @@chats
    OPTIONS { keepNull: false }
`

const resetState = `
UPDATE DOCUMENT(@key) WITH { bot_state: null } IN @@chats
OPTIONS { keepNull: false }
`

const setState = `
UPDATE DOCUMENT(@key) WITH { bot_state: @bot_state } IN @@chats 
`

func (c *ChatsController) DeleteGateways(ctx context.Context, fields map[string]*structpb.Value) error {
	for key, val := range fields {
		numValue := val.GetNumberValue()
		stringValue := val.GetStringValue()

		var queryValue interface{}

		if stringValue == "" {
			queryValue = numValue
		} else {
			queryValue = stringValue
		}

		_, err := c.db.Query(ctx, deleteGateways, map[string]interface{}{
			"@chats":  CHATS_COLLECTION,
			"gate":    key,
			"gate_id": queryValue,
		})
		if err != nil {
			return err
		}
	}
	return nil
}

func (c *ChatsController) SetBotState(ctx context.Context, chat *cc.Chat) (*cc.Chat, error) {
	log := c.log.Named("Set state")
	log.Debug("Req received")

	_, err := c.col.Database().Query(ctx, resetState, map[string]interface{}{
		"key":    driver.NewDocumentID(CHATS_COLLECTION, chat.GetUuid()),
		"@chats": CHATS_COLLECTION,
	})

	if err != nil {
		log.Error("Failed to reset state", zap.Error(err))
		return nil, err
	}

	_, err = c.col.Database().Query(ctx, setState, map[string]interface{}{
		"key":       driver.NewDocumentID(CHATS_COLLECTION, chat.GetUuid()),
		"@chats":    CHATS_COLLECTION,
		"bot_state": chat.GetBotState(),
	})

	if err != nil {
		log.Error("Failed to set state", zap.Error(err))
		return nil, err
	}

	return chat, nil
}

const getEarliest = `
LET earliest = FIRST(
	FOR c in @@chats
	FILTER c._key in @chats
	SORT c.created ASC
	RETURN c
)

RETURN earliest
`

const megreMessages = `
FOR m in @@messages
	FILTER m.chat in @chats && m.chat != @earliest
	UPDATE m with {chat: @earliest} in @@messages
`

const deleteOldChats = `
FOR c in @chats
	FILTER c != @earliest
	REMOVE c in @@chats
`

func (c *ChatsController) Merge(ctx context.Context, chats []string) (*cc.Chat, error) {
	log := c.log.Named("Merge")

	cur, err := c.db.Query(ctx, getEarliest, map[string]interface{}{
		"@chats": CHATS_COLLECTION,
		"chats":  chats,
	})

	if err != nil {
		log.Error("Failed to get earliest chats", zap.Error(err))
		return nil, err
	}

	var chat cc.Chat
	_, err = cur.ReadDocument(ctx, &chat)
	if err != nil {
		log.Error("Failed to read chat", zap.Error(err))
		return nil, err
	}

	_, err = c.db.Query(ctx, megreMessages, map[string]interface{}{
		"@messages": MESSAGES_COLLECTION,
		"chats":     chats,
		"earliest":  chat.GetUuid(),
	})

	if err != nil {
		log.Error("Failed to get merge messages", zap.Error(err))
		return nil, err
	}

	_, err = c.db.Query(ctx, deleteOldChats, map[string]interface{}{
		"@chats":   CHATS_COLLECTION,
		"chats":    chats,
		"earliest": chat.GetUuid(),
	})

	if err != nil {
		log.Error("Failed to delete chats", zap.Error(err))
		return nil, err
	}

	return &chat, nil
}

const syncQuery = `
LET deps = @deps

FOR c in @@chats
	FILTER c.department != null
	LET admins = deps[c.department]
	FILTER admins != null
	UPDATE c with {admins: admins} in @@chats
`

func (c *ChatsController) Sync(ctx context.Context, cfg *cc.Defaults) error {
	log := c.log.Named("Merge")

	var deps = make(map[string][]string)

	for _, val := range cfg.Departments {
		deps[val.Key] = val.Admins
	}

	_, err := c.db.Query(ctx, syncQuery, map[string]interface{}{
		"@chats": CHATS_COLLECTION,
		"deps":   deps,
	})

	if err != nil {
		log.Error("Failed to sync chats", zap.Error(err))
		return err
	}

	return nil
}
