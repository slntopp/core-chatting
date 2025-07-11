package attachments

import (
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/gorilla/mux"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	"github.com/slntopp/core-chatting/pkg/graph"
	"go.uber.org/zap"
)

type AttachmentsServer struct {
	log *zap.Logger

	host, bucket, port string

	ctrl     *graph.AttachmentsController
	s3Client *minio.Client
}

func NewAttacmentsServer(logger *zap.Logger, ctrl *graph.AttachmentsController, host, port, bucket, accessKey, secretKey string) *AttachmentsServer {
	client, err := minio.New(fmt.Sprintf("%s:%s", host, port), &minio.Options{
		Creds:  credentials.NewStaticV4(accessKey, secretKey, ""),
		Secure: true,
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		},
	})
	if err != nil {
		logger.Error("S3 Client not defined", zap.Error(err))
	}
	return &AttachmentsServer{log: logger.Named("AttachmentsServer"), ctrl: ctrl, host: host, port: port, bucket: bucket, s3Client: client}
}

func (s *AttachmentsServer) Hander(router *mux.Router) {
	s.log.Info("Registring handlers")
	router.Path("/attachments").HandlerFunc(s.upload).Methods(http.MethodPut)
	router.Path("/attachments/{uuid}").HandlerFunc(s.download).Methods(http.MethodGet)
	router.Path("/attachments/{uuid}").HandlerFunc(s.delete).Methods(http.MethodDelete)
	router.Path("/attachments").HandlerFunc(s.downloadBatch).Methods(http.MethodGet)
}

func (s *AttachmentsServer) upload(w http.ResponseWriter, r *http.Request) {
	log := s.log.Named("Upload")
	ctx := r.Context()

	body, err := io.ReadAll(r.Body)
	if err != nil {
		log.Error("Failed to read body", zap.Error(err))
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(fmt.Sprintf(`{"error": "%s"}`, err.Error())))
		return
	}

	var attachment graph.Attachment
	err = json.Unmarshal(body, &attachment)
	if err != nil {
		log.Error("Failed to parse body", zap.Error(err))
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(fmt.Sprintf(`{"error": "%s"}`, err.Error())))
		return
	}

	result, err := s.ctrl.Upload(ctx, &attachment)
	if err != nil {
		log.Error("Failed to save attachment", zap.Error(err))
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(fmt.Sprintf(`{"error": "%s"}`, err.Error())))
		return
	}

	presignedUrl, err := s.s3Client.PresignedPutObject(ctx, s.bucket, fmt.Sprintf("%s%s", result.Uuid, result.Ext), time.Second*60)
	if err != nil {
		log.Error("Failed to create url", zap.Error(err))
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(fmt.Sprintf(`{"error": "%s"}`, err.Error())))
		return
	}

	result.ObjectId = presignedUrl.String()

	response, err := json.Marshal(result)
	if err != nil {
		log.Error("Failed to parse response", zap.Error(err))
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(fmt.Sprintf(`{"error": "%s"}`, err.Error())))
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write(response)
}

func (s *AttachmentsServer) download(w http.ResponseWriter, r *http.Request) {
	log := s.log.Named("Download")
	ctx := r.Context()

	vars := mux.Vars(r)
	uuid := vars["uuid"]

	result, err := s.ctrl.Get(ctx, uuid)
	if err != nil {
		log.Error("Failed to save attachment", zap.Error(err))
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(fmt.Sprintf(`{"error": "%s"}`, err.Error())))
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte(fmt.Sprintf(`{"url": "%s:%s/%s/%s%s","title":"%s"}`, s.host, s.port, s.bucket, result.Uuid, result.Ext, result.Title)))

	//http.Redirect(w, r, fmt.Sprintf("%s/%s/%s", s.host, s.bucket, result.Uuid), http.StatusPermanentRedirect)
}

func (s *AttachmentsServer) downloadBatch(w http.ResponseWriter, r *http.Request) {
	log := s.log.Named("DownloadBatch")
	ctx := r.Context()

	const queryField = "ids"
	if !r.URL.Query().Has(queryField) {
		w.WriteHeader(http.StatusBadRequest)
		_, _ = w.Write([]byte(fmt.Sprintf(`{"error": "%s"}`, "no required query with name: "+queryField)))
		return
	}
	uuids := strings.Split(r.URL.Query().Get(queryField), ",")

	responseValues := make([]string, 0)
	for _, uuid := range uuids {
		if uuid == "" {
			continue
		}
		result, err := s.ctrl.Get(ctx, uuid)
		if err != nil {
			log.Error("Failed to get attachment", zap.Error(err), zap.String("uuid", uuid))
			continue
		}
		responseValues = append(responseValues, fmt.Sprintf(`"%s":{"url":"%s:%s/%s/%s%s","title":"%s"}`, uuid, s.host, s.port, s.bucket, result.Uuid, result.Ext, result.Title))
	}

	w.WriteHeader(http.StatusOK)
	result := fmt.Sprintf(`{%s}`, strings.Join(responseValues, ","))
	_, _ = w.Write([]byte(result))
}

func (s *AttachmentsServer) delete(w http.ResponseWriter, r *http.Request) {
	log := s.log.Named("Delete")
	ctx := r.Context()

	uuid := mux.Vars(r)["uuid"]

	att, err := s.ctrl.Get(ctx, uuid)
	if err != nil {
		log.Error("Failed to get attachment", zap.Error(err))
		http.Error(w, fmt.Sprintf(`{"error": "%s"}`, err.Error()), http.StatusNotFound)
		return
	}

	objectName := fmt.Sprintf("%s%s", att.Uuid, att.Ext)
	err = s.s3Client.RemoveObject(ctx, s.bucket, objectName, minio.RemoveObjectOptions{})
	if err != nil {
		log.Warn("Failed to remove object from S3", zap.String("object", objectName), zap.Error(err))
	}

	if err = s.ctrl.Delete(ctx, uuid); err != nil {
		log.Error("Failed to delete attachment record", zap.Error(err))
		http.Error(w, fmt.Sprintf(`{"error": "%s"}`, err.Error()), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	_, _ = w.Write([]byte(`{"status":"ok"}`))
}
