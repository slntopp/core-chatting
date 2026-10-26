package pubsub

import (
	"testing"

	"github.com/slntopp/core-chatting/cc"
)

func TestChatForUsersDropsOperatorOnlyLastMessage(t *testing.T) {
	note := &cc.Message{Kind: cc.Kind_ADMIN_ONLY, Content: "internal"}
	first := &cc.Message{Content: "hello"}
	chat := &cc.Chat{Meta: &cc.ChatMeta{LastMessage: note, FirstMessage: first}}

	got := chatForUsers(chat)
	if got.GetMeta().GetLastMessage() != nil {
		t.Fatal("operator-only last message reached the customer")
	}
	if got.GetMeta().GetFirstMessage().GetContent() != "hello" {
		t.Fatal("visible first message was dropped")
	}
	if chat.GetMeta().GetLastMessage() != note {
		t.Fatal("admins' copy of the chat was modified")
	}
}
