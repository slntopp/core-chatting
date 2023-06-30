package main

import (
	"fmt"
	"github.com/slntopp/core-chatting/pkg/core/auth"
	"net/http"

	"github.com/rabbitmq/amqp091-go"
	"github.com/slntopp/core-chatting/pkg/pubsub"
	"github.com/slntopp/core-chatting/pkg/stream"

	cc "github.com/slntopp/core-chatting/cc/ccconnect"

	"github.com/slntopp/core-chatting/pkg/chats"
	"github.com/slntopp/core-chatting/pkg/core"
	"github.com/slntopp/core-chatting/pkg/graph"
	"github.com/slntopp/core-chatting/pkg/messages"
	"github.com/slntopp/core-chatting/pkg/users"

	"github.com/bufbuild/connect-go"
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
	rbmq         string

	SIGNING_KEY []byte
)

func init() {
	viper.AutomaticEnv()
	log = core.NewLogger()

	viper.SetDefault("PORT", "8080")
	viper.SetDefault("DB_HOST", "localhost:8529")
	viper.SetDefault("DB_CRED", "root:password")
	viper.SetDefault("RABBITMQ_CONN", "amqp://guest:guest@localhost:5672/")
	viper.SetDefault("DB_NAME", "chats")
	viper.SetDefault("USERS_COL", "Accounts")

	viper.SetDefault("SIGNING_KEY", "secret")

	port = viper.GetString("PORT")

	arangodbHost = viper.GetString("DB_HOST")
	arangodbCred = viper.GetString("DB_CRED")
	dbName = viper.GetString("DB_NAME")
	usersCol = viper.GetString("USERS_COL")
	rbmq = viper.GetString("RABBITMQ_CONN")

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
	msgCtrs := graph.NewMessagesController(log, db)
	usersCtrl := graph.NewUsersController(log, db, usersCol)

	mux := http.NewServeMux()

	authInterceptor := auth.NewAuthInterceptor(log, SIGNING_KEY)

	interceptors := connect.WithInterceptors(authInterceptor)

	chatServer := chats.NewChatsServer(log, chatCtrl, ps)
	path, handler := cc.NewChatsAPIHandler(chatServer, interceptors)
	mux.Handle(path, handler)

	messagesServer := messages.NewMessagesServer(log, chatCtrl, msgCtrs, ps)
	path, handler = cc.NewMessagesAPIHandler(messagesServer, interceptors)
	mux.Handle(path, handler)

	usersServer := users.NewUsersServer(log, usersCtrl)
	path, handler = cc.NewUsersAPIHandler(usersServer, interceptors)
	mux.Handle(path, handler)

	streamServer := stream.NewStreamServer(log, usersCtrl, ps, SIGNING_KEY)
	path, handler = cc.NewStreamServiceHandler(streamServer, interceptors)
	mux.Handle(path, handler)

	host := fmt.Sprintf("0.0.0.0:%s", port)

	handler = cors.New(cors.Options{
		AllowedOrigins:      []string{"*"},
		AllowedMethods:      []string{"GET", "POST", "OPTIONS", "PUT", "DELETE"},
		AllowedHeaders:      []string{"*", "Connect-Protocol-Version"},
		AllowCredentials:    true,
		AllowPrivateNetwork: true,
	}).Handler(h2c.NewHandler(mux, &http2.Server{}))

	log.Debug("Start server", zap.String("host", host))
	err = http.ListenAndServe(host, handler)
	if err != nil {
		log.Fatal("Failed to start server", zap.Error(err))
	}
}
