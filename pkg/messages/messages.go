package messages

import (
	"context"
	"errors"
	"fmt"
	"github.com/slntopp/nocloud/pkg/nocloud/schema"
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

	log.Debug("Messages retrieved", zap.Any("messages", messages))
	resp := connect.NewResponse[cc.Messages](&cc.Messages{
		Messages: messages,
	})

	go s.ps.Pub(ctx, requestor, &cc.Event{Type: cc.EventType_CHAT_READ, Item: &cc.Event_Chat{Chat: req.Msg}})

	return resp, nil
}

func (s *MessagesServer) List(ctx context.Context, req *connect.Request[cc.MessagesListRequest]) (*connect.Response[cc.Messages], error) {
	log := s.log.Named("List")
	log.Debug("Request received", zap.Any("req", req.Msg))
	requester := ctx.Value(core.ChatAccount).(string)
	if schema.ROOT_ACCOUNT_KEY != requester {
		return nil, connect.NewError(connect.CodePermissionDenied, errors.New("no access rights"))
	}

	messages, err := s.msgCtrl.List(ctx)
	if err != nil {
		return nil, err
	}

	resp := connect.NewResponse[cc.Messages](&cc.Messages{
		Messages: messages,
	})

	return resp, nil
}

// pollLimit is how many answers one question may offer. Past a dozen nobody
// picks, they type.
const pollLimit = 12

// validatePoll checks a poll on its way in. The votes map is never accepted
// from a client: answers are recorded by Vote, and a sender who could preseed
// them could invent them.
func validatePoll(poll *cc.Poll) error {
	if poll == nil {
		return nil
	}
	if len(poll.GetOptions()) < 2 {
		return errors.New("a poll needs at least two options")
	}
	if len(poll.GetOptions()) > pollLimit {
		return fmt.Errorf("a poll takes at most %d options", pollLimit)
	}

	seen := map[string]bool{}
	for _, option := range poll.GetOptions() {
		if option.GetId() == "" || option.GetLabel() == "" {
			return errors.New("every poll option needs an id and a label")
		}
		if seen[option.GetId()] {
			return fmt.Errorf("duplicate poll option id %q", option.GetId())
		}
		seen[option.GetId()] = true
	}

	poll.Votes = nil
	poll.Closed = false
	return nil
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

	if err = validatePoll(msg.GetPoll()); err != nil {
		return nil, connect.NewError(connect.CodeInvalidArgument, err)
	}

	if (cc.IsOperatorOnly(msg.Kind) || msg.UnderReview) && chat.Role != cc.Role_ADMIN {
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

	if chat.GetResponsible() == "" && slices.Contains(chat.GetAdmins(), requestor) && !msg.GetUnderReview() {
		chat.Responsible = &requestor
		chat, err = s.chatCtrl.Update(ctx, chat)
		if err != nil {
			log.Error("Failed to update chat", zap.Error(err))
			go pubsub.HandleNotifyChat(ctx, log, s.ps, chat, cc.EventType_CHAT_UPDATED)
		}
	}

	go pubsub.HandleNotifyMessage(ctx, log, s.ps, message, chat, cc.EventType_MESSAGE_SENT)

	if s.whmcsTickets && chat.GetDepartment() != "openai" && !msg.GetUnderReview() {
		go s.ps.PubWhmcs(ctx, &cc.Event{
			Type: cc.EventType_MESSAGE_SENT,
			Item: &cc.Event_Msg{Msg: message},
		})
	}

	// An onboarding ticket keeps its status through the whole conversation.
	// We started it, not the customer, and a campaign of thousands moving
	// through OPEN and CUSTOMER_REPLY would bury the tickets people actually
	// opened themselves. It leaves this status only when someone closes it.
	if chat.GetStatus() != cc.Status_ONBOARDING {
		if slices.Contains(chat.GetAdmins(), requestor) {
			chat.Status = cc.Status_OPEN
		} else if chat.Status != cc.Status_NEW {
			chat.Status = cc.Status_CUSTOMER_REPLY
		}
	}

	update, err := s.chatCtrl.Update(ctx, chat)
	if err != nil {
		return nil, err
	}
	go pubsub.HandleNotifyChat(ctx, log, s.ps, update, cc.EventType_CHAT_STATUS_CHANGED)

	resp := connect.NewResponse[cc.Message](message)

	return resp, nil
}

// Vote answers the poll on a message.
//
// Answering again replaces the previous answer, and an empty list retracts it:
// a person who clicked the wrong button should be able to fix it, and the
// answer we act on is the last one they gave.
//
// The answer is not published to WHMCS. A vote is not an edit of the message,
// and the client that voted sends the choice as an ordinary message right
// after — that is what reaches the gateways, the WHMCS sync, the unread
// counter and the operator, through the paths that already exist.
func (s *MessagesServer) Vote(ctx context.Context, req *connect.Request[cc.VoteRequest]) (*connect.Response[cc.Message], error) {
	log := s.log.Named("Vote")
	log.Debug("Request received", zap.Any("req", req.Msg))

	requestor := ctx.Value(core.ChatAccount).(string)

	msg, err := s.msgCtrl.Get(ctx, req.Msg.GetMessage())
	if err != nil {
		return nil, connect.NewError(connect.CodeNotFound, errors.New("message not found"))
	}

	chat, err := s.chatCtrl.Get(ctx, msg.GetChat(), requestor)
	if err != nil {
		return nil, err
	}
	if chat.Role < cc.Role_USER {
		return nil, connect.NewError(connect.CodePermissionDenied, errors.New("no access to chat"))
	}

	poll := msg.GetPoll()
	if poll == nil {
		return nil, connect.NewError(connect.CodeInvalidArgument, errors.New("this message has no poll"))
	}
	if poll.GetClosed() {
		return nil, connect.NewError(connect.CodeFailedPrecondition, errors.New("this poll is closed"))
	}

	chosen := req.Msg.GetOptions()
	if len(chosen) > 1 && !poll.GetMultiple() {
		return nil, connect.NewError(connect.CodeInvalidArgument, errors.New("this poll takes one answer"))
	}
	for _, id := range chosen {
		if !slices.ContainsFunc(poll.GetOptions(), func(o *cc.PollOption) bool { return o.GetId() == id }) {
			return nil, connect.NewError(connect.CodeInvalidArgument, fmt.Errorf("no such option %q", id))
		}
	}

	var vote *cc.PollVote
	if len(chosen) > 0 {
		vote = &cc.PollVote{Options: chosen, Ts: time.Now().UnixMilli()}
	}

	updated, err := s.msgCtrl.Vote(ctx, msg.GetUuid(), requestor, vote)
	if err != nil {
		log.Error("Failed to record the vote", zap.Error(err))
		return nil, err
	}

	go pubsub.HandleNotifyMessage(ctx, log, s.ps, updated, chat, cc.EventType_MESSAGE_UPDATED)

	log.Info("Poll answered", zap.String("message", msg.GetUuid()),
		zap.String("account", requestor), zap.Strings("options", chosen))
	return connect.NewResponse[cc.Message](updated), nil
}

// pollsLimit caps one batch. A campaign asking for its answers pages through
// its own journal; an unbounded list would be a way to ask for the whole
// collection at once.
const pollsLimit = 500

// Polls returns the given messages, for reading the answers off their polls,
// and touches nothing: no read receipts, no events. See the RPC comment in
// cc.proto for why this is not Get.
func (s *MessagesServer) Polls(ctx context.Context, req *connect.Request[cc.PollsRequest]) (*connect.Response[cc.Messages], error) {
	log := s.log.Named("Polls")
	log.Debug("Request received", zap.Int("messages", len(req.Msg.GetMessages())))

	requestor := ctx.Value(core.ChatAccount).(string)

	messages := req.Msg.GetMessages()
	if len(messages) == 0 {
		return connect.NewResponse[cc.Messages](&cc.Messages{}), nil
	}
	if len(messages) > pollsLimit {
		return nil, connect.NewError(connect.CodeInvalidArgument,
			fmt.Errorf("at most %d messages at a time", pollsLimit))
	}

	found, err := s.msgCtrl.Polls(ctx, requestor, messages)
	if err != nil {
		log.Error("Failed to read messages", zap.Error(err))
		return nil, err
	}

	return connect.NewResponse[cc.Messages](&cc.Messages{Messages: found}), nil
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

	if cc.IsOperatorOnly(req.Msg.Kind) && chat.Role != cc.Role_ADMIN {
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

	message, err := s.msgCtrl.Update(ctx, req.Msg, true)
	if err != nil {
		return nil, err
	}

	if oldMessage.GetKind() != message.GetKind() || oldMessage.GetUnderReview() != message.GetUnderReview() {
		go pubsub.HandleSpecialNotify(ctx, log, s.ps, message, oldMessage, chat)
		sender := message.GetSender()
		if chat.GetResponsible() == "" && slices.Contains(chat.GetAdmins(), sender) && !message.GetUnderReview() {
			chat.Responsible = &sender
			chat, err = s.chatCtrl.Update(ctx, chat)
			if err != nil {
				log.Error("Failed to update chat", zap.Error(err))
				go pubsub.HandleNotifyChat(ctx, log, s.ps, chat, cc.EventType_CHAT_UPDATED)
			}
		}
	} else {
		go pubsub.HandleNotifyMessage(ctx, log, s.ps, message, chat, cc.EventType_MESSAGE_UPDATED)
	}

	if s.whmcsTickets && chat.GetDepartment() != "openai" {
		var tpe = cc.EventType_MESSAGE_UPDATED
		if oldMessage.GetUnderReview() && !message.GetUnderReview() {
			tpe = cc.EventType_MESSAGE_SENT
		}
		go s.ps.PubWhmcs(ctx, &cc.Event{
			Type: tpe,
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
