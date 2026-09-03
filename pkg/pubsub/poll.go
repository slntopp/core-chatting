package pubsub

import (
	"fmt"
	"strings"

	"github.com/slntopp/core-chatting/cc"
	"google.golang.org/protobuf/proto"
)

// withPollText is what keeps a poll from breaking the channels that know
// nothing about polls.
//
// Inside the panel a poll is a card with buttons. A person reading the same
// ticket over telegram, email or WHMCS gets the message text and nothing else,
// so their copy carries the options as a numbered list and they answer in
// words — an operator reads "3" perfectly well. The structured answer only
// exists for whoever clicked a button.
//
// The event is cloned: the message it points at is the one every other
// subscriber is looking at, and appending to it in place would put the list
// into the panel too, under the card that already shows it.
func withPollText(event *cc.Event) *cc.Event {
	msg := event.GetMsg()
	if msg == nil || len(msg.GetPoll().GetOptions()) == 0 {
		return event
	}

	out := proto.Clone(event).(*cc.Event)
	copied := out.GetMsg()

	var text strings.Builder
	text.WriteString(copied.GetContent())
	if question := copied.GetPoll().GetQuestion(); question != "" {
		text.WriteString("\n\n")
		text.WriteString(question)
	}
	for i, option := range copied.GetPoll().GetOptions() {
		text.WriteString(fmt.Sprintf("\n%d. %s", i+1, option.GetLabel()))
	}

	copied.Content = text.String()
	return out
}
