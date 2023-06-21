package api

import (
	"fmt"
	"net/http"

	"github.com/bufbuild/connect-go"
	"github.com/rs/cors"
	cc "github.com/slntopp/core-chatting/cc/ccconnect"
	"golang.org/x/net/http2"
	"golang.org/x/net/http2/h2c"

	"github.com/slntopp/core-chatting/pkg/chats"
	"github.com/slntopp/core-chatting/pkg/core"
	"github.com/slntopp/core-chatting/pkg/core/auth"
	"github.com/slntopp/core-chatting/pkg/graph"
	"github.com/slntopp/core-chatting/pkg/messages"
	"github.com/spf13/viper"
	"go.uber.org/zap"
)

var (
	log *zap.Logger

	port         string
	arangodbHost string
	arangodbCred string
	dbName       string
)

func init() {
	viper.AutomaticEnv()
	log = core.NewLogger()

	viper.SetDefault("PORT", "8080")
	viper.SetDefault("DB_HOST", "localhost:8529")
	viper.SetDefault("DB_CRED", "root:openSesame")
	viper.SetDefault("DB_NAME", "name")

	port = viper.GetString("PORT")

	arangodbHost = viper.GetString("DB_HOST")
	arangodbCred = viper.GetString("DB_CRED")
	dbName = viper.GetString("DB_NAME")
}

func main() {
	db := graph.ConnectDb(log, arangodbHost, arangodbCred, dbName)

	chatCtrl := graph.NewChatsController(log, db)
	msgCtrs := graph.NewMessagesController(log, db)

	mux := http.NewServeMux()

	interceptors := connect.WithInterceptors(auth.NewUserInterceptor())

	chatServer := chats.NewChatsServer(log, chatCtrl)
	path, handler := cc.NewChatsAPIHandler(chatServer)
	mux.Handle(path, handler)

	messagesServer := messages.NewMessagesServer(log, chatCtrl, msgCtrs)
	path, handler = cc.NewMessagesAPIHandler(messagesServer, interceptors)
	mux.Handle(path, handler)

	host := fmt.Sprintf("0.0.0.0:%s", port)

	handler = cors.New(cors.Options{
		AllowedOrigins:      []string{"*"},
		AllowedMethods:      []string{"GET", "POST", "OPTIONS", "PUT", "DELETE"},
		AllowedHeaders:      []string{"*", "Connect-Protocol-Version"},
		AllowCredentials:    true,
		AllowPrivateNetwork: true,
		Debug:               true,
	}).Handler(h2c.NewHandler(mux, &http2.Server{}))

	log.Debug("Start server", zap.String("host", host))
	err := http.ListenAndServe(host, handler)
	if err != nil {
		log.Fatal("Failed to start server", zap.Error(err))
	}
}
