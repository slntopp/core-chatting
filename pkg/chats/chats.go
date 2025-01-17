package chats

import (
	"context"
	"errors"
	"github.com/rabbitmq/amqp091-go"
	settingspb "github.com/slntopp/nocloud-proto/settings"
	"slices"

	"github.com/slntopp/core-chatting/pkg/core"
	"github.com/slntopp/core-chatting/pkg/pubsub"
	"google.golang.org/protobuf/types/known/structpb"

	"github.com/slntopp/core-chatting/cc"
	"github.com/slntopp/core-chatting/pkg/graph"

	"connectrpc.com/connect"
	"go.uber.org/zap"
)

type ChatsServer struct {
	log  *zap.Logger
	conn *amqp091.Connection

	ctrl       *graph.ChatsController
	users_ctrl *graph.UsersController
	msgsCtrl   *graph.MessagesController
	ps         *pubsub.PubSub

	settingsClient settingspb.SettingsServiceClient

	whmcsTickets bool
}

func NewChatsServer(logger *zap.Logger, ctrl *graph.ChatsController, users_ctrl *graph.UsersController, msgsCtrl *graph.MessagesController,
	ps *pubsub.PubSub, whmcsTickets bool, settingsClient settingspb.SettingsServiceClient, conn *amqp091.Connection) *ChatsServer {
	return &ChatsServer{log: logger.Named("ChatsServer"), ctrl: ctrl, users_ctrl: users_ctrl,
		ps: ps, whmcsTickets: whmcsTickets, settingsClient: settingsClient, conn: conn, msgsCtrl: msgsCtrl}
}

func (s *ChatsServer) Create(ctx context.Context, req *connect.Request[cc.Chat]) (*connect.Response[cc.Chat], error) {
	log := s.log.Named("Create")
	log.Debug("Request received", zap.Any("req", req.Msg))

	requestor := ctx.Value(core.ChatAccount).(string)

	req.Msg.Owner = requestor

	msg := req.Msg
	resolve, err := s.users_ctrl.Resolve(ctx, []string{requestor})
	if err != nil {
		return nil, err
	}
	user := resolve[0]

	if user.GetData() != nil {
		fields := user.GetData().GetFields()

		log.Debug("Field", zap.Any("f", fields))

		if msg.GetMeta() == nil {
			msg.Meta = &cc.ChatMeta{
				Data: map[string]*structpb.Value{},
			}
		}

		if msg.GetMeta().GetData() == nil {
			msg.Meta.Data = map[string]*structpb.Value{}
		}

		if fields != nil {
			for _, gate := range msg.GetGateways() {
				log.Debug("gate", zap.String("g", gate))
				if val, ok := fields[gate]; ok {
					log.Debug("ok")
					msg.Meta.Data[gate] = val
				}
			}
		}
		err := s.ctrl.DeleteGateways(ctx, fields)
		if err != nil {
			log.Error("Failed to delete old gates", zap.Error(err))
			return nil, err
		}
	}

	chat, err := s.ctrl.Create(ctx, req.Msg)
	if err != nil {
		return nil, err
	}

	go pubsub.HandleNotifyChat(ctx, log, s.ps, chat, cc.EventType_CHAT_CREATED)

	if s.whmcsTickets && chat.GetDepartment() != "openai" {
		go s.ps.PubWhmcs(ctx, &cc.Event{
			Type: cc.EventType_CHAT_CREATED,
			Item: &cc.Event_Chat{Chat: chat},
		})
	}

	resp := connect.NewResponse[cc.Chat](chat)

	return resp, nil
}

func (s *ChatsServer) Update(ctx context.Context, req *connect.Request[cc.Chat]) (*connect.Response[cc.Chat], error) {
	log := s.log.Named("Create")
	log.Debug("Request received", zap.Any("req", req.Msg))

	requestor := ctx.Value(core.ChatAccount).(string)

	chat, err := s.ctrl.Get(ctx, req.Msg.Uuid, requestor)
	if err != nil {
		return nil, err
	}

	if chat.Role < cc.Role_OWNER {
		return nil, connect.NewError(connect.CodePermissionDenied, errors.New("no access to chat"))
	}

	chat, err = s.ctrl.Update(ctx, req.Msg)
	if err != nil {
		log.Error("Failed to update chat", zap.Error(err))
		return nil, err
	}

	go pubsub.HandleNotifyChat(ctx, log, s.ps, chat, cc.EventType_CHAT_UPDATED)

	if s.whmcsTickets && chat.GetDepartment() != "openai" {
		go s.ps.PubWhmcs(ctx, &cc.Event{
			Type: cc.EventType_CHAT_UPDATED,
			Item: &cc.Event_Chat{Chat: chat},
		})
	}

	resp := connect.NewResponse(chat)

	return resp, nil
}

func (s *ChatsServer) Get(ctx context.Context, req *connect.Request[cc.Chat]) (*connect.Response[cc.Chat], error) {
	log := s.log.Named("Create")
	log.Debug("Request received", zap.Any("req", req.Msg))

	requestor := ctx.Value(core.ChatAccount).(string)

	chat, err := s.ctrl.Get(ctx, req.Msg.Uuid, requestor)
	if err != nil {
		return nil, err
	}

	if chat.Role == cc.Role_NOACCESS {
		return nil, connect.NewError(connect.CodePermissionDenied, errors.New("no access to chat"))
	}

	resp := connect.NewResponse[cc.Chat](chat)

	return resp, nil
}

func (s *ChatsServer) List(ctx context.Context, req *connect.Request[cc.ListChatsRequest]) (*connect.Response[cc.ListChatsResponse], error) {
	log := s.log.Named("List")
	log.Debug("Request received", zap.Any("req", req.Msg))

	requester := ctx.Value(core.ChatAccount).(string)

	chats, total, err := s.ctrl.List(ctx, requester, req.Msg)
	if err != nil {
		return nil, err
	}

	resp := connect.NewResponse[cc.ListChatsResponse](&cc.ListChatsResponse{
		Pool:  chats,
		Total: total,
	})

	return resp, nil
}

func (s *ChatsServer) Count(ctx context.Context, req *connect.Request[cc.Empty]) (*connect.Response[cc.CountChatsResponse], error) {
	log := s.log.Named("Count")
	log.Debug("Request received", zap.Any("req", req.Msg))

	requester := ctx.Value(core.ChatAccount).(string)

	count, err := s.ctrl.Count(ctx, requester)
	if err != nil {
		return nil, err
	}

	resp := connect.NewResponse[cc.CountChatsResponse](&cc.CountChatsResponse{
		Statuses: count,
	})

	return resp, nil
}

func (s *ChatsServer) Delete(ctx context.Context, req *connect.Request[cc.Chat]) (*connect.Response[cc.Chat], error) {
	log := s.log.Named("Create")
	log.Debug("Request received", zap.Any("req", req.Msg))

	requestor := ctx.Value(core.ChatAccount).(string)

	chat, err := s.ctrl.Get(ctx, req.Msg.Uuid, requestor)
	if err != nil {
		return nil, err
	}

	if chat.Role < cc.Role_OWNER {
		return nil, connect.NewError(connect.CodePermissionDenied, errors.New("no access to chat"))
	}

	chat, err = s.ctrl.Delete(ctx, req.Msg)
	if err != nil {
		return nil, err
	}

	go pubsub.HandleNotifyChat(ctx, log, s.ps, chat, cc.EventType_CHAT_DELETED)

	if s.whmcsTickets && chat.GetDepartment() != "openai" {
		go s.ps.PubWhmcs(ctx, &cc.Event{
			Type: cc.EventType_CHAT_DELETED,
			Item: &cc.Event_Chat{Chat: chat},
		})
	}

	resp := connect.NewResponse[cc.Chat](chat)

	return resp, nil
}

func (s *ChatsServer) SetBotState(ctx context.Context, req *connect.Request[cc.Chat]) (*connect.Response[cc.Chat], error) {
	log := s.log.Named("SetBotState")
	log.Debug("Request received", zap.Any("req", req.Msg))

	requestor := ctx.Value(core.ChatAccount).(string)

	chat, err := s.ctrl.Get(ctx, req.Msg.Uuid, requestor)
	if err != nil {
		return nil, err
	}

	if chat.Role < cc.Role_OWNER {
		return nil, connect.NewError(connect.CodePermissionDenied, errors.New("no access to chat"))
	}

	chat, err = s.ctrl.SetBotState(ctx, req.Msg)
	if err != nil {
		return nil, err
	}

	go pubsub.HandleNotifyChat(ctx, log, s.ps, chat, cc.EventType_CHAT_UPDATED)

	resp := connect.NewResponse[cc.Chat](chat)

	return resp, nil
}

func (s *ChatsServer) GetBotState(ctx context.Context, req *connect.Request[cc.Chat]) (*connect.Response[cc.Chat], error) {
	log := s.log.Named("GetBotState")
	log.Debug("Request received", zap.Any("req", req.Msg))

	requestor := ctx.Value(core.ChatAccount).(string)

	chat, err := s.ctrl.Get(ctx, req.Msg.Uuid, requestor)
	if err != nil {
		return nil, err
	}

	if chat.Role == cc.Role_NOACCESS {
		return nil, connect.NewError(connect.CodePermissionDenied, errors.New("no access to chat"))
	}

	resp := connect.NewResponse[cc.Chat](chat)

	return resp, nil
}

func (s *ChatsServer) ChangeDepartment(ctx context.Context, req *connect.Request[cc.Chat]) (*connect.Response[cc.Chat], error) {
	log := s.log.Named("GetBotState")
	log.Debug("Request received", zap.Any("req", req.Msg))

	requestor := ctx.Value(core.ChatAccount).(string)

	chat, err := s.ctrl.Get(ctx, req.Msg.Uuid, requestor)
	if err != nil {
		return nil, err
	}

	if chat.Role == cc.Role_NOACCESS {
		return nil, connect.NewError(connect.CodePermissionDenied, errors.New("no access to chat"))
	}

	if chat.GetRole() != cc.Role_ADMIN {
		return nil, connect.NewError(connect.CodePermissionDenied, errors.New("not enough rights"))
	}

	newDepartment := req.Msg.GetDepartment()

	chat.Department = newDepartment

	config, err := core.Config()
	if err != nil {
		return nil, err
	}

	for _, department := range config.GetDepartments() {
		if department.GetKey() == newDepartment {
			chat.Admins = department.GetAdmins()
			break
		}
	}

	chat, err = s.ctrl.Update(ctx, req.Msg)
	if err != nil {
		return nil, err
	}

	go pubsub.HandleNotifyChat(ctx, log, s.ps, chat, cc.EventType_CHAT_DEPARTMENT_CHANGED)

	if s.whmcsTickets && chat.GetDepartment() != "openai" {
		go s.ps.PubWhmcs(ctx, &cc.Event{
			Type: cc.EventType_CHAT_DEPARTMENT_CHANGED,
			Item: &cc.Event_Chat{Chat: chat},
		})
	}

	resp := connect.NewResponse[cc.Chat](chat)

	return resp, nil
}

func (s *ChatsServer) ChangeGateway(ctx context.Context, req *connect.Request[cc.Chat]) (*connect.Response[cc.Chat], error) {
	log := s.log.Named("ChangeGateway")
	log.Debug("Request received", zap.Any("req", req.Msg))

	requestor := ctx.Value(core.ChatAccount).(string)

	chat, err := s.ctrl.Get(ctx, req.Msg.Uuid, requestor)
	if err != nil {
		return nil, err
	}

	if chat.Role == cc.Role_NOACCESS {
		return nil, connect.NewError(connect.CodePermissionDenied, errors.New("no access to chat"))
	}

	msg := req.Msg
	resolve, err := s.users_ctrl.Resolve(ctx, []string{requestor})
	if err != nil {
		return nil, err
	}
	user := resolve[0]

	if user.GetData() != nil {
		fields := user.GetData().GetFields()

		log.Debug("Field", zap.Any("f", fields))

		if msg.GetMeta() == nil {
			msg.Meta = &cc.ChatMeta{
				Data: map[string]*structpb.Value{},
			}
		}

		if msg.GetMeta().GetData() == nil {
			msg.Meta.Data = map[string]*structpb.Value{}
		}

		if fields != nil {
			for _, gate := range msg.GetGateways() {
				log.Debug("gate", zap.String("g", gate))
				if val, ok := fields[gate]; ok {
					log.Debug("ok")
					msg.Meta.Data[gate] = val
				}
			}
		}
		err := s.ctrl.DeleteGateways(ctx, fields)
		if err != nil {
			log.Error("Failed to delete old gates", zap.Error(err))
			return nil, err
		}
	}

	update, err := s.ctrl.Update(ctx, msg)
	if err != nil {
		return nil, err
	}

	go pubsub.HandleNotifyChat(ctx, log, s.ps, update, cc.EventType_CHAT_UPDATED)

	resp := connect.NewResponse[cc.Chat](update)

	return resp, nil
}

func (s *ChatsServer) ChangeStatus(ctx context.Context, req *connect.Request[cc.Chat]) (*connect.Response[cc.Chat], error) {
	log := s.log.Named("ChangeGateway")
	log.Debug("Request received", zap.Any("req", req.Msg))

	requestor := ctx.Value(core.ChatAccount).(string)

	chat, err := s.ctrl.Get(ctx, req.Msg.Uuid, requestor)
	if err != nil {
		return nil, err
	}

	if chat.Role == cc.Role_NOACCESS {
		return nil, connect.NewError(connect.CodePermissionDenied, errors.New("no access to chat"))
	}

	msg := req.Msg
	chat.Status = msg.GetStatus()

	update, err := s.ctrl.Update(ctx, msg)
	if err != nil {
		return nil, err
	}

	go pubsub.HandleNotifyChat(ctx, log, s.ps, update, cc.EventType_CHAT_STATUS_CHANGED)

	if s.whmcsTickets && chat.GetDepartment() != "openai" {
		go s.ps.PubWhmcs(ctx, &cc.Event{
			Type: cc.EventType_CHAT_STATUS_CHANGED,
			Item: &cc.Event_Chat{Chat: update},
		})
	}

	resp := connect.NewResponse[cc.Chat](update)

	return resp, nil
}

func (s *ChatsServer) MergeChats(ctx context.Context, req *connect.Request[cc.Merge]) (*connect.Response[cc.Chat], error) {
	log := s.log.Named("ChangeGateway")
	log.Debug("Request received", zap.Any("req", req.Msg))

	requestor := ctx.Value(core.ChatAccount).(string)

	chats := req.Msg.GetChats()

	for _, uuid := range chats {
		chat, err := s.ctrl.Get(ctx, uuid, requestor)
		if err != nil {
			return nil, err
		}

		if chat.Role == cc.Role_NOACCESS {
			return nil, connect.NewError(connect.CodePermissionDenied, errors.New("no access to chat"))
		}
	}

	chat, err := s.ctrl.Merge(ctx, chats)
	if err != nil {
		return nil, err
	}

	go pubsub.HandleNotifyChat(ctx, log, s.ps, chat, cc.EventType_CHATS_MERGED)

	resp := connect.NewResponse[cc.Chat](chat)

	return resp, nil
}

func (s *ChatsServer) SyncChats(ctx context.Context, req *connect.Request[cc.Empty]) (*connect.Response[cc.Empty], error) {
	log := s.log.Named("SyncChats")
	log.Debug("Request received", zap.Any("req", req.Msg))

	requestor := ctx.Value(core.ChatAccount).(string)

	config, err := core.Config()
	if err != nil {
		return nil, err
	}

	if !slices.Contains(config.GetAdmins(), requestor) {
		return nil, connect.NewError(connect.CodePermissionDenied, errors.New("not enough rights"))
	}

	return &connect.Response[cc.Empty]{}, nil
}
