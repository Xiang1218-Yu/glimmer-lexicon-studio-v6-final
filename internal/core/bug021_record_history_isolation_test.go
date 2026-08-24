package core

import (
	"context"
	"fmt"
	"strings"
	"testing"
	"time"
)

func TestBug021RecordHistoryIsolation(t *testing.T) {
	engine := NewEngine()
	createBugRecord(t, engine, "history-isolation", "history")
	if _, err := engine.Advance(context.Background(), "history-isolation", "review", "alice"); err != nil {
		t.Fatal(err)
	}
	record, _ := engine.Get(context.Background(), "history-isolation")
	record.History[0].By = "mutated"
	again, _ := engine.Get(context.Background(), "history-isolation")
	if again.History[0].By == "mutated" {
		t.Fatal("history slice aliases stored state")
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
