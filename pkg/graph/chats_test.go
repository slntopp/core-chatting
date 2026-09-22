package graph

import (
	"testing"

	"github.com/slntopp/core-chatting/cc"
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/types/known/structpb"
)

func boolPtr(v bool) *bool { return &v }

// The handoff is the bot's own write (disabled+escalated, no state), and it
// happens while a chat may already carry a mode an operator set from the UI.
// Replacing the map instead of merging it silently un-set that mode.
func TestMergeBotStateKeepsKeysTheCallerDidNotSend(t *testing.T) {
	current := map[string]*structpb.Value{
		"otus.mode":   structpb.NewStringValue("copilot_mode"),
		"skip_review": structpb.NewBoolValue(true),
	}
	got := MergeBotState(current, &cc.SetBotStateRequest{
		Chat:      "chat-1",
		Disabled:  boolPtr(true),
		Escalated: boolPtr(true),
	})
	if got["otus.mode"].GetStringValue() != "copilot_mode" {
		t.Fatalf("the handoff must not clear the chat's mode, got %v", got["otus.mode"])
	}
	if !got["skip_review"].GetBoolValue() {
		t.Fatalf("a key nobody sent must survive, got %v", got["skip_review"])
	}
	if !got["disabled"].GetBoolValue() || !got["escalated"].GetBoolValue() {
		t.Fatalf("the handoff must still land, got %v", got)
	}
}

// The other direction: an operator moving the mode radio sends that one key and
// must not drop the flag the bot wrote when it handed the chat over.
func TestMergeBotStateKeepsTheHandoffWhenOnlyTheModeIsSent(t *testing.T) {
	current := map[string]*structpb.Value{
		"disabled":  structpb.NewBoolValue(true),
		"escalated": structpb.NewBoolValue(true),
	}
	got := MergeBotState(current, &cc.SetBotStateRequest{
		Chat:  "chat-1",
		State: map[string]*structpb.Value{"otus.mode": structpb.NewStringValue("self_mode")},
	})
	if got["otus.mode"].GetStringValue() != "self_mode" {
		t.Fatalf("the mode must land, got %v", got["otus.mode"])
	}
	if !got["escalated"].GetBoolValue() || !got["disabled"].GetBoolValue() {
		t.Fatalf("a chat already handed to a human must stay that way, got %v", got)
	}
}

// Clearing an override writes an empty string, because merging cannot delete a
// key. Every reader treats empty as "no override", so this is what must land.
func TestMergeBotStateClearsTheModeWithAnEmptyString(t *testing.T) {
	current := map[string]*structpb.Value{"otus.mode": structpb.NewStringValue("god_mode")}
	got := MergeBotState(current, &cc.SetBotStateRequest{
		Chat:  "chat-1",
		State: map[string]*structpb.Value{"otus.mode": structpb.NewStringValue("")},
	})
	if got["otus.mode"].GetStringValue() != "" {
		t.Fatalf("the override must be cleared, got %v", got["otus.mode"])
	}
}

// The plugin sends both the switches and (when it moved) the mode in one
// request, on a chat whose state was never written.
func TestMergeBotStateOnAChatWithNoStateYet(t *testing.T) {
	got := MergeBotState(nil, &cc.SetBotStateRequest{
		Chat:       "chat-1",
		Disabled:   boolPtr(false),
		SkipReview: boolPtr(false),
		State:      map[string]*structpb.Value{"otus.mode": structpb.NewStringValue("copilot_mode")},
	})
	if got["disabled"].GetBoolValue() || got["skip_review"].GetBoolValue() {
		t.Fatalf("the switches must land as sent, got %v", got)
	}
	if got["otus.mode"].GetStringValue() != "copilot_mode" {
		t.Fatalf("the mode must land, got %v", got["otus.mode"])
	}
}

// The typed fields are the contract; a caller that sends both must not be able
// to disable a chat through the freeform bag while saying the opposite.
func TestMergeBotStateNamedFieldsWinOverTheSameKeyInState(t *testing.T) {
	got := MergeBotState(nil, &cc.SetBotStateRequest{
		Chat:     "chat-1",
		Disabled: boolPtr(false),
		State:    map[string]*structpb.Value{"disabled": structpb.NewBoolValue(true)},
	})
	if got["disabled"].GetBoolValue() {
		t.Fatalf("the typed field must win, got %v", got["disabled"])
	}
}

// The exact body the plugin builds when an operator moves the radio and saves
// (SetBotStateRequest.fromJson in chat_actions.vue, dumped from protobuf-es).
// Kept literal: this is the wire between the two languages, and a change in the
// dialog that breaks it must fail here rather than in production.
func TestMergeBotStateTakesThePluginsOwnRequest(t *testing.T) {
	const fromTheDialog = `{"chat":"chat-1","state":{"otus.mode":"self_mode"},"disabled":false,"skipReview":true}`

	req := &cc.SetBotStateRequest{}
	if err := protojson.Unmarshal([]byte(fromTheDialog), req); err != nil {
		t.Fatalf("the dialog's own request must parse: %v", err)
	}
	got := MergeBotState(map[string]*structpb.Value{
		"escalated": structpb.NewBoolValue(true),
	}, req)
	if got["otus.mode"].GetStringValue() != "self_mode" {
		t.Fatalf("the operator's choice must land, got %v", got["otus.mode"])
	}
	if got["disabled"].GetBoolValue() || !got["skip_review"].GetBoolValue() {
		t.Fatalf("the two switches must land as the dialog sent them, got %v", got)
	}
	if !got["escalated"].GetBoolValue() {
		t.Fatalf("saving the dialog must not un-escalate a chat, got %v", got)
	}
}
