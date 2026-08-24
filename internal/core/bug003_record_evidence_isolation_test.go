package core

import (
	"context"
	"fmt"
	"strings"
	"testing"
	"time"
)

func TestBug003RecordEvidenceIsolation(t *testing.T) {
	engine := NewEngine()
	createBugRecord(t, engine, "evidence-term", "evidence")
	first, err := engine.Get(context.Background(), "evidence-term")
	if err != nil {
		t.Fatal(err)
	}
	if len(first.Evidence) < 2 {
		t.Fatalf("expected module evidence, got %d", len(first.Evidence))
	}
	first.Evidence[0].Detail = "changed-by-caller"
	second, err := engine.Get(context.Background(), "evidence-term")
	if err != nil {
		t.Fatal(err)
	}
	if second.Evidence[0].Detail == "changed-by-caller" {
		t.Fatal("caller mutation changed the stored evidence slice")
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
