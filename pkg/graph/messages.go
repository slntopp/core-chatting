package graph

import (
	"context"
	"sync"

	"github.com/slntopp/core-chatting/cc"
	"github.com/slntopp/nocloud/pkg/nocloud/schema"

	"github.com/arangodb/go-driver"
	"go.uber.org/zap"
)

type MessagesController struct {
	log *zap.Logger
	m   *sync.Mutex

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
		m:   &sync.Mutex{},
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

func (c *MessagesController) Update(ctx context.Context, msg *cc.Message, withTransaction ...bool) (*cc.Message, error) {
	_ = c.log.Named("Update")

	var (
		inTransaction = len(withTransaction) > 0 && withTransaction[0]
		trID          driver.TransactionID
		err           error
	)
	if inTransaction {
		trID, err = c.db.BeginTransaction(ctx, driver.TransactionCollections{
			Exclusive: []string{MESSAGES_COLLECTION},
		}, &driver.BeginTransactionOptions{AllowImplicit: true})
		if err != nil {
			return nil, err
		}
		ctx = driver.WithTransactionID(ctx, trID)
	}
	commit := func() error {
		if inTransaction {
			return c.db.CommitTransaction(ctx, trID, nil)
		}
		return nil
	}
	abort := func() {
		if inTransaction {
			_ = c.db.AbortTransaction(ctx, trID, nil)
		}
	}

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
	_, err = c.db.Query(ctx, updateMessageQuery, params)
	if err != nil {
		abort()
		return nil, err
	}

	err = commit()
	return msg, err
}

const readMessageQuery = `
LET message = Document(@message)
LET new = UNIQUE(PUSH(message.readers, @reader))
UPDATE message with {readers: new} in @@messages RETURN NEW
`

func (c *MessagesController) Read(ctx context.Context, msg *cc.Message, reader string) (*cc.Message, error) {
	log := c.log.Named("Read")
	log.Debug("Req received")

	c.m.Lock()
	defer c.m.Unlock()
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

const voteQuery = `
LET msg = DOCUMENT(@@messages, @message)
FILTER msg != null
UPDATE msg WITH { poll: { votes: { [@account]: @vote } } } IN @@messages
	OPTIONS { keepNull: false, mergeObjects: true }
RETURN NEW
`

// Vote records one person's answer to a poll, or removes it when vote is nil.
//
// One statement on one document: an answer is a merge into the votes map, so
// two people answering at the same time cannot overwrite each other and there
// is nothing to lock. keepNull: false is what makes a nil vote a retraction
// rather than an empty answer.
func (c *MessagesController) Vote(ctx context.Context, message, account string, vote *cc.PollVote) (*cc.Message, error) {
	log := c.log.Named("Vote")
	log.Debug("Req received", zap.String("message", message), zap.String("account", account))

	// The vote goes in whole, not field by field. A hand-written map here was
	// one forgotten field away from dropping data silently, and did exactly
	// that: it kept the options and the timestamp and quietly lost which
	// message states the answer, so every change of answer posted another one.
	// The generated json tags carry the proto field names, so this is the same
	// document shape.
	var value interface{}
	if vote != nil {
		value = vote
	}

	cur, err := c.db.Query(ctx, voteQuery, map[string]interface{}{
		"@messages": MESSAGES_COLLECTION,
		"message":   message,
		"account":   account,
		"vote":      value,
	})
	if err != nil {
		return nil, err
	}
	defer cur.Close()

	var updated cc.Message
	if _, err = cur.ReadDocument(ctx, &updated); err != nil {
		return nil, err
	}
	return &updated, nil
}

const pollMessagesQuery = `
FOR m IN @@messages
	FILTER m._key IN @messages
	LET chat = DOCUMENT(@@chats, m.chat)
	FILTER chat != null
	// The same access rule the rest of the service uses, applied per message:
	// this returns only what the caller could have opened by hand.
	FILTER @requestor IN chat.admins OR chat.owner == @requestor
	    OR @requestor IN chat.users OR @requestor == @root_account
	RETURN m
`

// Polls returns the given messages without touching them.
//
// Deliberately not built on Get: that one marks messages read for the caller,
// which is right for someone opening a chat and wrong for a service reading
// answers in the background.
func (c *MessagesController) Polls(ctx context.Context, requestor string, messages []string) ([]*cc.Message, error) {
	log := c.log.Named("Polls")
	log.Debug("Req received", zap.Int("messages", len(messages)))

	cur, err := c.db.Query(ctx, pollMessagesQuery, map[string]interface{}{
		"@messages":    MESSAGES_COLLECTION,
		"@chats":       CHATS_COLLECTION,
		"messages":     messages,
		"requestor":    requestor,
		"root_account": schema.ROOT_ACCOUNT_KEY,
	})
	if err != nil {
		return nil, err
	}
	defer cur.Close()

	var out []*cc.Message
	for cur.HasMore() {
		var msg cc.Message
		if _, err = cur.ReadDocument(ctx, &msg); err != nil {
			return nil, err
		}
		out = append(out, &msg)
	}
	return out, nil
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

const listMessages = `
LET messages = (
	FOR m IN @@messages
	RETURN m
)
RETURN messages
`

func (c *MessagesController) List(ctx context.Context) ([]*cc.Message, error) {
	log := c.log.Named("List")
	log.Debug("Req received")

	var msg []*cc.Message

	cur, err := c.db.Query(ctx, listMessages, map[string]interface{}{
		"@messages": MESSAGES_COLLECTION,
	})
	if err != nil {
		return nil, err
	}

	_, err = cur.ReadDocument(ctx, &msg)
	if err != nil {
		return nil, err
	}

	return msg, nil
}
