package core

import (
	"context"
	"fmt"
	"strings"
	"testing"
	"time"
)

func TestBug018HistoryActorPreserved(t *testing.T) {
	engine := NewEngine()
	createBugRecord(t, engine, "actor-term", "actor")
	if _, err := engine.Advance(context.Background(), "actor-term", "review", "alice"); err != nil {
		t.Fatal(err)
	}
	got, err := engine.Get(context.Background(), "actor-term")
	if err != nil {
		t.Fatal(err)
	}
	if got.History[len(got.History)-1].By != "alice" {
		t.Fatalf("actor was lost: %#v", got.History)
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
