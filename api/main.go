package main

import (
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/slntopp/core-chatting/pkg/core/auth"

	"github.com/rabbitmq/amqp091-go"
	"github.com/slntopp/core-chatting/pkg/pubsub"
	"github.com/slntopp/core-chatting/pkg/stream"

	cc "github.com/slntopp/core-chatting/cc/ccconnect"

	"github.com/slntopp/core-chatting/pkg/chats"
	"github.com/slntopp/core-chatting/pkg/core"
	"github.com/slntopp/core-chatting/pkg/graph"
	"github.com/slntopp/core-chatting/pkg/messages"
	"github.com/slntopp/core-chatting/pkg/users"

	"connectrpc.com/connect"
	"github.com/rs/cors"
	"github.com/spf13/viper"
	"go.uber.org/zap"
	"golang.org/x/net/http2"
	"golang.org/x/net/http2/h2c"
)

var (
	log *zap.Logger

	port string

	arangodbHost string
	arangodbCred string
	dbName       string
	usersCol     string

	rbmq string

	dist string

	SIGNING_KEY []byte
)

func init() {
	viper.AutomaticEnv()
	log = core.NewLogger()

	viper.SetDefault("PORT", "8080")
	port = viper.GetString("PORT")

	viper.SetDefault("DB_HOST", "localhost:8529")
	viper.SetDefault("DB_CRED", "root:password")
	viper.SetDefault("DB_NAME", "chats")
	viper.SetDefault("USERS_COL", "Accounts")
	arangodbHost = viper.GetString("DB_HOST")
	arangodbCred = viper.GetString("DB_CRED")
	dbName = viper.GetString("DB_NAME")
	usersCol = viper.GetString("USERS_COL")

	viper.SetDefault("RABBITMQ_CONN", "amqp://guest:guest@localhost:5672/")
	rbmq = viper.GetString("RABBITMQ_CONN")

	viper.SetDefault("DIST", "dist")
	dist = viper.GetString("DIST")

	viper.SetDefault("SIGNING_KEY", "secret")
	SIGNING_KEY = []byte(viper.GetString("SIGNING_KEY"))
}

func main() {
	rbmq, err := amqp091.Dial(rbmq)
	if err != nil {
		log.Fatal("Failed to connect to RabbitMQ", zap.Error(err))
	}
	defer rbmq.Close()

	ps := pubsub.NewPubSub(log, rbmq)

	db := graph.ConnectDb(log, arangodbHost, arangodbCred, dbName)

	chatCtrl := graph.NewChatsController(log, db)
	msgCtrl := graph.NewMessagesController(log, db)
	usersCtrl := graph.NewUsersController(log, db, usersCol)

	router := mux.NewRouter()
	router.Use(func(h http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			log.Debug("Request", zap.String("method", r.Method), zap.String("path", r.URL.Path))
			h.ServeHTTP(w, r)
		})
	})

	router.PathPrefix("/cc.ui/").Handler(
		http.StripPrefix("/cc.ui/", http.FileServer(http.Dir(dist))),
	)

	authInterceptor := auth.NewAuthInterceptor(log, SIGNING_KEY)

	interceptors := connect.WithInterceptors(authInterceptor)

	chatServer := chats.NewChatsServer(log, chatCtrl, usersCtrl, ps)
	path, handler := cc.NewChatsAPIHandler(chatServer, interceptors)
	router.PathPrefix(path).Handler(handler)

	messagesServer := messages.NewMessagesServer(log, chatCtrl, msgCtrl, ps)
	path, handler = cc.NewMessagesAPIHandler(messagesServer, interceptors)
	router.PathPrefix(path).Handler(handler)

	usersServer := users.NewUsersServer(log, usersCtrl)
	path, handler = cc.NewUsersAPIHandler(usersServer, interceptors)
	router.PathPrefix(path).Handler(handler)

	streamServer := stream.NewStreamServer(log, usersCtrl, msgCtrl, ps, SIGNING_KEY)
	path, handler = cc.NewStreamServiceHandler(streamServer, interceptors)
	router.PathPrefix(path).Handler(handler)

	host := fmt.Sprintf("0.0.0.0:%s", port)

	handler = cors.New(cors.Options{
		AllowedOrigins:      []string{"*"},
		AllowedMethods:      []string{"GET", "POST", "OPTIONS", "PUT", "DELETE"},
		AllowedHeaders:      []string{"*", "Connect-Protocol-Version"},
		AllowCredentials:    true,
		AllowPrivateNetwork: true,
	}).Handler(h2c.NewHandler(router, &http2.Server{}))

	log.Debug("Start server", zap.String("host", host))
	err = http.ListenAndServe(host, handler)
	if err != nil {
		log.Fatal("Failed to start server", zap.Error(err))
	}
}
