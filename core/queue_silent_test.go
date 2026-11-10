package core

import (
	"fmt"
	"testing"
)

func TestQueueSilent_PreservesFIFOAndErrorReplies(t *testing.T) {
	for _, silent := range []bool{false, true} {
		t.Run(fmt.Sprintf("silent=%t", silent), func(t *testing.T) {
			p := &stubPlatformEngine{n: "test"}
			sess := newQueuingSession("queue-silent")
			e := NewEngine("test", &controllableAgent{nextSession: sess}, []Platform{p}, "", LangEnglish)
			if silent {
				e.SetQueueSilent(true)
			}
			e.SetMaxQueuedMessages(2)
			key := "test:queue-silent"
			state := &interactiveState{agentSession: sess, platform: p, replyCtx: "ctx"}
			e.interactiveStates[key] = state
			for _, content := range []string{"first", "second"} {
				if !e.queueMessageForBusySession(p, &Message{SessionKey: key, Content: content, ReplyCtx: "ctx"}, key) {
					t.Fatal("message was not accepted")
				}
			}
			wantReplies := 2
			if silent {
				wantReplies = 0
			}
			if got := len(p.getSent()); got != wantReplies {
				t.Fatalf("acceptance replies = %d, want %d", got, wantReplies)
			}
			e.queueMessageForBusySession(p, &Message{SessionKey: key, Content: "overflow", ReplyCtx: "ctx"}, key)
			if got := p.getSent(); len(got) != wantReplies+1 || got[len(got)-1] != fmt.Sprintf(e.i18n.T(MsgQueueFull), 2) {
				t.Fatalf("queue-full reply missing: %v", got)
			}
			if len(state.pendingMessages) != 2 || state.pendingMessages[0].content != "first" || state.pendingMessages[1].content != "second" {
				t.Fatal("queue contents or FIFO order changed")
			}
			e.notifyDroppedQueuedMessages(state, fmt.Errorf("agent exited"))
			if got := len(p.getSent()); got != wantReplies+3 {
				t.Fatalf("failure replies = %d, want %d", got, wantReplies+3)
			}
		})
	}
}
