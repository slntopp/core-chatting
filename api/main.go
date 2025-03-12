package main

import (
	"context"
	"fmt"
	"github.com/slntopp/core-chatting/pkg/notifications"
	http_server "github.com/slntopp/nocloud/pkg/nocloud/http"
	"github.com/slntopp/nocloud/pkg/nocloud/rabbitmq"
	"github.com/slntopp/nocloud/pkg/nocloud/schema"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/metadata"
	"gopkg.in/yaml.v2"
	"net/http"
	"os"
	"sync"

	"github.com/gorilla/mux"
	"github.com/slntopp/core-chatting/pkg/attachments"
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

	settingspb "github.com/slntopp/nocloud-proto/settings"
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

	whmcsTickets bool

	s3host, s3bucket, s3port, s3AccKey, s3SecKey string

	SIGNING_KEY []byte

	settingsHost string

	gatewaysConfigYaml string
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

	viper.SetDefault("WHMCS_TICKETS", "true")
	whmcsTickets = viper.GetBool("WHMCS_TICKETS")

	viper.SetDefault("SIGNING_KEY", "secret")
	SIGNING_KEY = []byte(viper.GetString("SIGNING_KEY"))

	viper.SetDefault("S3_HOST", "host")
	viper.SetDefault("S3_BUCKET", "bucket")
	viper.SetDefault("S3_PORT", "port")
	viper.SetDefault("ACC_KEY", "key")
	viper.SetDefault("SEC_KEY", "key")
	s3host = viper.GetString("S3_HOST")
	s3bucket = viper.GetString("S3_BUCKET")
	s3port = viper.GetString("S3_PORT")
	s3AccKey = viper.GetString("ACC_KEY")
	s3SecKey = viper.GetString("SEC_KEY")

	viper.SetDefault("SETTINGS_HOST", "settings:8000")
	settingsHost = viper.GetString("SETTINGS_HOST")

	viper.SetDefault("NOTIFY_GATEWAYS_CONFIG", "./gateways_config.yaml")
	gatewaysConfigYaml = viper.GetString("NOTIFY_GATEWAYS_CONFIG")

}

func main() {
	workers := &sync.WaitGroup{}

	conn, err := amqp091.Dial(rbmq)
	if err != nil {
		log.Fatal("Failed to connect to RabbitMQ", zap.Error(err))
	}
	defer conn.Close()
	rabbitmq.FatalOnConnectionClose(log, conn)

	ps := pubsub.NewPubSub(log, conn)

	db := graph.ConnectDb(log, arangodbHost, arangodbCred, dbName)

	chatCtrl := graph.NewChatsController(log, db)
	msgCtrl := graph.NewMessagesController(log, db)
	usersCtrl := graph.NewUsersController(log, db, usersCol)
	attachmentsCtrl := graph.NewAttachmentsController(log, db)

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

	token, err := authInterceptor.MakeToken(schema.ROOT_ACCOUNT_KEY)
	if err != nil {
		log.Fatal("Can't generate token", zap.Error(err))
	}
	ctx := metadata.AppendToOutgoingContext(context.Background(), "authorization", "bearer "+token)
	ctx, cancel := context.WithCancel(ctx)

	settingsConn, err := grpc.Dial(settingsHost, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		panic(err)
	}
	settingsClient := settingspb.NewSettingsServiceClient(settingsConn)
	log.Info("Check settings server")
	if _, err = settingsClient.Get(ctx, &settingspb.GetRequest{}); err != nil {
		log.Fatal("Can't check settings connection", zap.Error(err))
	}
	log.Info("Settings server connection established")

	log.Info("Obtaining notify gateways configurations")
	notifyConfigYaml, err := os.ReadFile(gatewaysConfigYaml)
	if err != nil {
		log.Fatal("Failed to read gateways config file", zap.Error(err))
	}
	var configs []notifications.GatewayConfig
	err = yaml.Unmarshal(notifyConfigYaml, &configs)
	if err != nil {
		log.Fatal("Failed to parse yaml config", zap.Error(err))
	}
	notify := notifications.NewNotificationDispatcher(log, configs)
	log.Info("Notification dispatcher initialized")

	chatServer := chats.NewChatsServer(log, chatCtrl, usersCtrl, msgCtrl, ps, whmcsTickets, settingsClient, conn, notify)
	path, handler := cc.NewChatsAPIHandler(chatServer, interceptors)
	router.PathPrefix(path).Handler(handler)
	go chatServer.CloseInactiveChatsRoutine(ctx, worker(workers))
	go chatServer.SLAViolationRoutine(ctx, worker(workers))

	messagesServer := messages.NewMessagesServer(log, chatCtrl, msgCtrl, attachmentsCtrl, ps, whmcsTickets)
	path, handler = cc.NewMessagesAPIHandler(messagesServer, interceptors)
	router.PathPrefix(path).Handler(handler)

	usersServer := users.NewUsersServer(log, usersCtrl)
	path, handler = cc.NewUsersAPIHandler(usersServer, interceptors)
	router.PathPrefix(path).Handler(handler)

	streamServer := stream.NewStreamServer(log, usersCtrl, msgCtrl, chatCtrl, ps, SIGNING_KEY)
	path, handler = cc.NewStreamServiceHandler(streamServer, interceptors)
	router.PathPrefix(path).Handler(handler)

	attServer := attachments.NewAttacmentsServer(log, attachmentsCtrl, s3host, s3port, s3bucket, s3AccKey, s3SecKey)
	attServer.Hander(router)

	host := fmt.Sprintf("0.0.0.0:%s", port)

	handler = cors.New(cors.Options{
		AllowedOrigins:      []string{"*"},
		AllowedMethods:      []string{"GET", "POST", "OPTIONS", "PUT", "DELETE"},
		AllowedHeaders:      []string{"*", "Connect-Protocol-Version"},
		AllowCredentials:    true,
		AllowPrivateNetwork: true,
	}).Handler(h2c.NewHandler(router, &http2.Server{}))

	http_server.Serve(log, host, handler)
	log.Info("Stopping workers.")
	cancel()
	workers.Wait()
	log.Info("All workers were stopped.")
}

func worker(wg *sync.WaitGroup) *sync.WaitGroup {
	wg.Add(1)
	return wg
}
