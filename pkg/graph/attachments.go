package graph

import (
	"context"
	"path/filepath"

	"github.com/arangodb/go-driver"
	"go.uber.org/zap"
)

type Attachment struct {
	Uuid     string `json:"uuid"`
	Title    string `json:"title"`
	Chat     string `json:"chat"`
	ObjectId string `json:"object_id"`
	Ext      string `json:"ext"`
}

type AttachmentsController struct {
	log *zap.Logger

	db  driver.Database
	col driver.Collection
}

func NewAttachmentsController(logger *zap.Logger, db driver.Database) *AttachmentsController {
	log := logger.Named("Attachments")
	log.Debug("Creating attachments controller")

	ctx := context.TODO()
	col := GetEnsureCollection(log, ctx, db, ATTACHMENTS_COLLECTION)

	return &AttachmentsController{
		log: log,
		db:  db,
		col: col,
	}
}

func (c *AttachmentsController) Upload(ctx context.Context, a *Attachment) (*Attachment, error) {
	log := c.log.Named("Upload")
	a.ObjectId = ""
	a.Ext = filepath.Ext(a.Title)

	document, err := c.col.CreateDocument(ctx, a)
	if err != nil {
		log.Error("Failed to create document", zap.Error(err))
		return nil, err
	}

	a.Uuid = document.Key

	_, err = c.col.UpdateDocument(ctx, a.Uuid, a)
	if err != nil {
		c.col.RemoveDocument(ctx, a.Uuid)
		log.Error("Failed to update document", zap.Error(err))
		return nil, err
	}
	return a, nil
}

func (c *AttachmentsController) Get(ctx context.Context, uuid string) (*Attachment, error) {
	log := c.log.Named("Get")

	var attachment Attachment
	_, err := c.col.ReadDocument(ctx, uuid, &attachment)
	if err != nil {
		log.Error("Failed to read document", zap.Error(err))
		return nil, err
	}
	return &attachment, nil
}

func (c *AttachmentsController) Delete(ctx context.Context, uuid string) error {
	log := c.log.Named("Get")

	_, err := c.col.RemoveDocument(ctx, uuid)
	if err != nil {
		log.Error("Failed to delete document", zap.Error(err))
		return err
	}
	return nil
}
