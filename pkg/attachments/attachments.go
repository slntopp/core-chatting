package attachments

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/slntopp/core-chatting/pkg/graph"
	"go.uber.org/zap"
)

type AttachmentsServer struct {
	log *zap.Logger

	host, bucket string

	ctrl *graph.AttachmentsController
}

func NewAttacmentsServer(logger *zap.Logger, ctrl *graph.AttachmentsController, host, bucket string) *AttachmentsServer {
	return &AttachmentsServer{log: logger.Named("AttachmentsServer"), ctrl: ctrl, host: host, bucket: bucket}
}

func (s *AttachmentsServer) Hander(router *mux.Router) {
	s.log.Info("Registring handlers")
	router.Path("/attachments").HandlerFunc(s.upload).Methods(http.MethodPut)
	router.Path("/attachments/{uuid}").HandlerFunc(s.download).Methods(http.MethodGet)
	router.Path("/attachments/{uuid}").HandlerFunc(s.delete).Methods(http.MethodDelete)
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
	log := s.log.Named("Upload")
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

	http.Redirect(w, r, fmt.Sprintf("%s/%s/%s", s.host, s.bucket, result.Uuid), http.StatusPermanentRedirect)
}

func (s *AttachmentsServer) delete(w http.ResponseWriter, r *http.Request) {
	log := s.log.Named("Upload")
	ctx := r.Context()

	vars := mux.Vars(r)
	uuid := vars["uuid"]

	err := s.ctrl.Delete(ctx, uuid)
	if err != nil {
		log.Error("Failed to save attachment", zap.Error(err))
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(fmt.Sprintf(`{"error": "%s"}`, err.Error())))
		return
	}

	w.WriteHeader(http.StatusOK)
}
