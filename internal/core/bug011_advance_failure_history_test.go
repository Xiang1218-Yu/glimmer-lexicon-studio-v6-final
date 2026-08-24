package core

import (
	"context"
	"fmt"
	"strings"
	"testing"
	"time"
)

func TestBug011AdvanceFailureHistory(t *testing.T) {
	engine := NewEngine()
	createBugRecord(t, engine, "history-term", "history")
	before, _ := engine.Get(context.Background(), "history-term")
	if _, err := engine.Advance(context.Background(), "history-term", "not-a-stage", "reviewer"); err == nil {
		t.Fatal("expected invalid stage error")
	}
	after, err := engine.Get(context.Background(), "history-term")
	if err != nil {
		t.Fatal(err)
	}
	if len(after.History) != len(before.History) || after.Version != before.Version {
		t.Fatalf("failed advance changed history: before=%d/%d after=%d/%d", len(before.History), before.Version, len(after.History), after.Version)
	}
}

func createBugRecord(t *testing.T, engine *Engine, id, payload string) Record {
	t.Helper()
	record, err := engine.Create(context.Background(), id, payload, "tester")
	if err != nil {
		t.Fatalf("Engine.Create failed: %v", err)
	}
	return record
}

var (
	_ = fmt.Sprintf
	_ = strings.Contains
	_ = time.Second
)
