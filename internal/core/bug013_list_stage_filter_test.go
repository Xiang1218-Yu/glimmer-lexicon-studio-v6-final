package core

import (
	"context"
	"fmt"
	"strings"
	"testing"
	"time"
)

func TestBug013ListStageFilter(t *testing.T) {
	engine := NewEngine()
	createBugRecord(t, engine, "stage-draft", "stage filter")
	createBugRecord(t, engine, "stage-review", "stage filter")
	if _, err := engine.Advance(context.Background(), "stage-review", "review", "reviewer"); err != nil {
		t.Fatal(err)
	}
	records, err := engine.List(context.Background(), "review", 100)
	if err != nil {
		t.Fatal(err)
	}
	if len(records) != 1 || records[0].Stage != "review" {
		t.Fatalf("wrong stage filter result: %#v", records)
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
