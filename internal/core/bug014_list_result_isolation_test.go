package core

import (
	"context"
	"fmt"
	"strings"
	"testing"
	"time"
)

func TestBug014ListResultIsolation(t *testing.T) {
	engine := NewEngine()
	createBugRecord(t, engine, "list-history", "list history")
	records, err := engine.List(context.Background(), "", 100)
	if err != nil {
		t.Fatal(err)
	}
	records[0].History[0].By = "page-mutated"
	again, err := engine.Get(context.Background(), "list-history")
	if err != nil {
		t.Fatal(err)
	}
	if again.History[0].By == "page-mutated" {
		t.Fatal("list result mutation polluted stored history")
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
