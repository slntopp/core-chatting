// Package otus proxies the Learn window to the Otus admin HTTP API.
//
// Otus is a separate process guarded by a static shared secret
// (OTUS_ADMIN_TOKEN), not by the NoCloud JWT that ai-bot-manager and the
// NoCloud gateway accept. That secret must not reach the browser, so the Vue
// app calls these routes and this handler adds the bearer.
//
// These are plain mux handlers, so the Connect auth interceptor does not run:
// each one validates the caller's JWT itself and refuses anyone who is not an
// admin of the chat. Learn makes Otus read a whole ticket and Save makes it
// open a pull request; a client must not be able to trigger either.
package otus

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/gorilla/mux"
	"github.com/slntopp/core-chatting/cc"
	"github.com/slntopp/core-chatting/pkg/core"
	"github.com/slntopp/core-chatting/pkg/core/auth"
	"github.com/slntopp/core-chatting/pkg/graph"
	"go.uber.org/zap"
)

// Propose and Save run a model over the whole ticket and can take a minute.
const callTimeout = 180 * time.Second

// A Learn body is a chat uuid and, on Save, the takes the operator kept. The
// takes are forwarded untouched: Otus owns their shape and their limits.
type learnBody struct {
	Chat string `json:"chat"`
}

type LearnServer struct {
	log *zap.Logger

	ctrl *graph.ChatsController

	url        string
	token      string
	signingKey []byte
	client     *http.Client
}

func NewLearnServer(logger *zap.Logger, ctrl *graph.ChatsController, otusUrl, token string, signingKey []byte) *LearnServer {
	return &LearnServer{
		log:        logger.Named("OtusLearnServer"),
		ctrl:       ctrl,
		url:        strings.TrimRight(otusUrl, "/"),
		token:      token,
		signingKey: signingKey,
		client:     &http.Client{Timeout: callTimeout},
	}
}

func (s *LearnServer) Hander(router *mux.Router) {
	s.log.Info("Registring handlers")
	router.Path("/otus/learn").HandlerFunc(s.learn).Methods(http.MethodGet, http.MethodPost)
	router.Path("/otus/learn/save").HandlerFunc(s.save).Methods(http.MethodPost)
}

func fail(w http.ResponseWriter, status int, why string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	body, _ := json.Marshal(map[string]any{"ok": false, "why": why})
	w.Write(body)
}

// requestor returns the account uuid the bearer token was minted for. The
// Connect interceptor does the same thing for the RPC handlers.
func (s *LearnServer) requestor(r *http.Request) (string, error) {
	segments := strings.Split(r.Header.Get("Authorization"), " ")
	if len(segments) != 2 {
		return "", fmt.Errorf("invalid token")
	}
	claims, err := auth.ValidateToken(s.signingKey, segments[1])
	if err != nil {
		return "", err
	}
	account, ok := claims[core.JWT_ACCOUNT_CLAIM].(string)
	if !ok || account == "" {
		return "", fmt.Errorf("no account claim")
	}
	return account, nil
}

// allowed answers whether this caller may drive Learn on this chat. Admins of
// the chat may; its owner and its users may not, however the UI is served.
func (s *LearnServer) allowed(r *http.Request, chat string) (int, string) {
	if s.url == "" {
		return http.StatusNotImplemented, "otus is not configured"
	}
	if chat == "" {
		return http.StatusBadRequest, "chat is required"
	}
	account, err := s.requestor(r)
	if err != nil {
		return http.StatusUnauthorized, "unauthorized"
	}
	found, err := s.ctrl.Get(r.Context(), chat, account)
	if err != nil {
		s.log.Debug("Failed to get chat", zap.String("chat", chat), zap.Error(err))
		return http.StatusForbidden, "no access to this chat"
	}
	if found.GetRole() != cc.Role_ADMIN {
		return http.StatusForbidden, "only an admin of this chat may use Learn"
	}
	return 0, ""
}

// forward passes the call to Otus with the admin bearer and copies its answer
// back verbatim: the UI switches on Otus's own why values and status codes.
func (s *LearnServer) forward(w http.ResponseWriter, method, path string, body []byte) {
	req, err := http.NewRequest(method, s.url+path, bytes.NewReader(body))
	if err != nil {
		fail(w, http.StatusInternalServerError, "failed to build the request")
		return
	}
	req.Header.Set("Authorization", "Bearer "+s.token)
	req.Header.Set("Content-Type", "application/json")

	res, err := s.client.Do(req)
	if err != nil {
		s.log.Error("Failed to reach otus", zap.String("path", path), zap.Error(err))
		fail(w, http.StatusBadGateway, "otus is unreachable")
		return
	}
	defer res.Body.Close()

	answer, err := io.ReadAll(res.Body)
	if err != nil {
		fail(w, http.StatusBadGateway, "otus answered with a broken body")
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(res.StatusCode)
	w.Write(answer)
}

// learn is both the propose (POST) and the draft read (GET).
func (s *LearnServer) learn(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		chat := r.URL.Query().Get("chat")
		if status, why := s.allowed(r, chat); status != 0 {
			fail(w, status, why)
			return
		}
		s.forward(w, http.MethodGet, "/learn?chat="+url.QueryEscape(chat), nil)
		return
	}

	body, err := io.ReadAll(r.Body)
	if err != nil {
		fail(w, http.StatusBadRequest, "failed to read the body")
		return
	}
	var parsed learnBody
	if err = json.Unmarshal(body, &parsed); err != nil {
		fail(w, http.StatusBadRequest, "body must be JSON")
		return
	}
	if status, why := s.allowed(r, parsed.Chat); status != 0 {
		fail(w, status, why)
		return
	}
	s.forward(w, http.MethodPost, "/learn", body)
}

func (s *LearnServer) save(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		fail(w, http.StatusBadRequest, "failed to read the body")
		return
	}
	var parsed learnBody
	if err = json.Unmarshal(body, &parsed); err != nil {
		fail(w, http.StatusBadRequest, "body must be JSON")
		return
	}
	if status, why := s.allowed(r, parsed.Chat); status != 0 {
		fail(w, status, why)
		return
	}
	s.forward(w, http.MethodPost, "/learn/save", body)
}
