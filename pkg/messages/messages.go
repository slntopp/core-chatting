package messages

import (
	"context"
	"errors"
	"fmt"
	"reflect"
	"slices"
	"time"

	"google.golang.org/protobuf/types/known/structpb"

	"github.com/slntopp/core-chatting/pkg/pubsub"

	"github.com/slntopp/core-chatting/cc"
	"github.com/slntopp/core-chatting/pkg/core"
	"github.com/slntopp/core-chatting/pkg/graph"

	"connectrpc.com/connect"
	"go.uber.org/zap"
)

type MessagesServer struct {
	log *zap.Logger

	chatCtrl        *graph.ChatsController
	msgCtrl         *graph.MessagesController
	attachmentsCtrl *graph.AttachmentsController
	ps              *pubsub.PubSub

	whmcsTickets bool
}

func NewMessagesServer(logger *zap.Logger, chatCtrl *graph.ChatsController, msgCtrl *graph.MessagesController, attachmentsCtrl *graph.AttachmentsController, ps *pubsub.PubSub, whmcsTickets bool) *MessagesServer {
	return &MessagesServer{log: logger.Named("MessagesServer"), chatCtrl: chatCtrl, msgCtrl: msgCtrl, attachmentsCtrl: attachmentsCtrl, ps: ps, whmcsTickets: whmcsTickets}
}

func (s *MessagesServer) Get(ctx context.Context, req *connect.Request[cc.Chat]) (*connect.Response[cc.Messages], error) {
	log := s.log.Named("Get")
	log.Debug("Request received", zap.Any("req", req.Msg))

	requestor := ctx.Value(core.ChatAccount).(string)

	chat, err := s.chatCtrl.Get(ctx, req.Msg.GetUuid(), requestor)
	if err != nil {
		return nil, err
	}

	if chat.Role < cc.Role_USER {
		return nil, connect.NewError(connect.CodePermissionDenied, errors.New("no access to chat"))
	}

	messages, err := s.chatCtrl.GetMessages(ctx, req.Msg, chat.Role == cc.Role_ADMIN)
	if err != nil {
		return nil, err
	}

	for index := range messages {
		if !core.In(requestor, messages[index].Readers) {
			newMessage, err := s.msgCtrl.Read(ctx, messages[index], requestor)
			if err != nil {
				log.Error("Failed to update reader", zap.Error(err))
			}
			messages[index] = newMessage
		}
	}

	resp := connect.NewResponse[cc.Messages](&cc.Messages{
		Messages: messages,
	})

	go s.ps.Pub(ctx, requestor, &cc.Event{Type: cc.EventType_CHAT_READ, Item: &cc.Event_Chat{Chat: req.Msg}})

	return resp, nil
}

func (s *MessagesServer) Send(ctx context.Context, req *connect.Request[cc.Message]) (*connect.Response[cc.Message], error) {
	log := s.log.Named("Send")
	log.Debug("Request received", zap.Any("req", req.Msg))

	requestor := ctx.Value(core.ChatAccount).(string)

	msg := req.Msg

	chat, err := s.chatCtrl.Get(ctx, msg.GetChat(), requestor)
	if err != nil {
		return nil, err
	}
	log.Debug("Chat retrieved", zap.String("requestor", requestor), zap.String("role", chat.Role.String()))

	if chat.Role < cc.Role_USER {
		return nil, connect.NewError(connect.CodePermissionDenied, errors.New("no access to chat"))
	}

	msg.Sender = requestor

	if (msg.Kind == cc.Kind_ADMIN_ONLY || msg.UnderReview) && chat.Role != cc.Role_ADMIN {
		return nil, connect.NewError(connect.CodeUnauthenticated, errors.New("can't send admin only message"))
	}

	msg.Sent = time.Now().UnixMilli()

	if chat.GetGateways() != nil {
		if msg.GetMeta() == nil {
			msg.Meta = map[string]*structpb.Value{}
		}
		data := chat.GetMeta().GetData()
		if data != nil {
			for _, gate := range chat.GetGateways() {
				if _, ok := data[gate]; !ok {
					continue
				}
				msg.Meta[fmt.Sprintf("%s_chat_id", gate)] = data[gate]
			}
		}
	}

	message, err := s.msgCtrl.Send(ctx, msg)
	if err != nil {
		return nil, err
	}

	log.Info("Sending result message", zap.Any("message", message))

	if chat.GetResponsible() == "" && slices.Contains(chat.GetAdmins(), requestor) {
		chat.Responsible = &requestor
		chat, err = s.chatCtrl.Update(ctx, chat)
		if err != nil {
			log.Error("Failed to update chat", zap.Error(err))
			go pubsub.HandleNotifyChat(ctx, log, s.ps, chat, cc.EventType_CHAT_UPDATED)
		}
	}

	go pubsub.HandleNotifyMessage(ctx, log, s.ps, message, chat, cc.EventType_MESSAGE_SENT)

	if s.whmcsTickets && chat.GetDepartment() != "openai" {
		go s.ps.PubWhmcs(ctx, &cc.Event{
			Type: cc.EventType_MESSAGE_SENT,
			Item: &cc.Event_Msg{Msg: message},
		})
	}

	if slices.Contains(chat.GetAdmins(), requestor) {
		chat.Status = cc.Status_OPEN
	} else {
		chat.Status = cc.Status_CUSTOMER_REPLY
	}

	update, err := s.chatCtrl.Update(ctx, chat)
	if err != nil {
		return nil, err
	}
	go pubsub.HandleNotifyChat(ctx, log, s.ps, update, cc.EventType_CHAT_STATUS_CHANGED)

	resp := connect.NewResponse[cc.Message](message)

	return resp, nil
}

func (s *MessagesServer) Update(ctx context.Context, req *connect.Request[cc.Message]) (*connect.Response[cc.Message], error) {
	log := s.log.Named("Update")
	log.Debug("Request received", zap.Any("req", req.Msg))

	requestor := ctx.Value(core.ChatAccount).(string)

	chat, err := s.chatCtrl.Get(ctx, req.Msg.GetChat(), requestor)
	if err != nil {
		return nil, err
	}

	if chat.Role < cc.Role_USER {
		return nil, connect.NewError(connect.CodePermissionDenied, errors.New("no access to chat"))
	}

	if req.Msg.Kind == cc.Kind_ADMIN_ONLY && chat.Role != cc.Role_ADMIN {
		return nil, connect.NewError(connect.CodeUnauthenticated, errors.New("can't send admin only message"))
	}

	if requestor != req.Msg.GetSender() && !core.In(requestor, chat.GetAdmins()) {
		return nil, connect.NewError(connect.CodeUnauthenticated, errors.New("you are not message sender or chat admin"))
	}

	req.Msg.Edited = time.Now().UnixMilli()

	oldMessage, err := s.msgCtrl.Get(ctx, req.Msg.GetUuid())

	if !reflect.DeepEqual(oldMessage.GetAttachments(), req.Msg.GetAttachments()) {
		oldAttachments := map[string]struct{}{}
		newAttachments := map[string]struct{}{}

		for _, attachment := range oldMessage.GetAttachments() {
			oldAttachments[attachment] = struct{}{}
		}
		for _, attachment := range req.Msg.GetAttachments() {
			newAttachments[attachment] = struct{}{}
		}
		for key := range oldAttachments {
			if _, ok := newAttachments[key]; !ok {
				s.attachmentsCtrl.Delete(ctx, key)
			}
		}
	}

	if err != nil {
		return nil, err
	}

	message, err := s.msgCtrl.Update(ctx, req.Msg)
	if err != nil {
		return nil, err
	}

	if oldMessage.GetKind() != message.GetKind() || oldMessage.GetUnderReview() != message.GetUnderReview() {
		go pubsub.HandleSpecialNotify(ctx, log, s.ps, message, oldMessage, chat)
	} else {
		go pubsub.HandleNotifyMessage(ctx, log, s.ps, message, chat, cc.EventType_MESSAGE_UPDATED)
	}

	if s.whmcsTickets && chat.GetDepartment() != "openai" {
		go s.ps.PubWhmcs(ctx, &cc.Event{
			Type: cc.EventType_MESSAGE_UPDATED,
			Item: &cc.Event_Msg{Msg: message},
		})
	}

	resp := connect.NewResponse[cc.Message](message)

	return resp, nil
}

func (s *MessagesServer) Delete(ctx context.Context, req *connect.Request[cc.Message]) (*connect.Response[cc.Message], error) {
	log := s.log.Named("Delete")
	log.Debug("Request received", zap.Any("req", req.Msg))

	requestor := ctx.Value(core.ChatAccount).(string)

	chat, err := s.chatCtrl.Get(ctx, req.Msg.GetChat(), requestor)
	if err != nil {
		return nil, err
	}

	if chat.Role < cc.Role_USER {
		return nil, connect.NewError(connect.CodePermissionDenied, errors.New("no access to chat"))
	}

	if requestor != req.Msg.GetSender() && !core.In(requestor, chat.GetAdmins()) {
		return nil, connect.NewError(connect.CodeUnauthenticated, errors.New("you are not message sender or chat admin"))
	}

	msg, err := s.msgCtrl.Get(ctx, req.Msg.GetUuid())
	if err != nil {
		return nil, err
	}

	for _, attachment := range msg.GetAttachments() {
		s.attachmentsCtrl.Delete(ctx, attachment)
	}

	message, err := s.msgCtrl.Delete(ctx, req.Msg)
	if err != nil {
		return nil, err
	}

	go pubsub.HandleNotifyMessage(ctx, log, s.ps, message, chat, cc.EventType_MESSAGE_DELETED)

	if s.whmcsTickets && chat.GetDepartment() != "openai" {
		go s.ps.PubWhmcs(ctx, &cc.Event{
			Type: cc.EventType_MESSAGE_DELETED,
			Item: &cc.Event_Msg{Msg: message},
		})
	}

	resp := connect.NewResponse[cc.Message](message)

	return resp, nil
}
