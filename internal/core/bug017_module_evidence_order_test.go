package core

import (
	"context"
	"fmt"
	"strings"
	"testing"
	"time"
)

func TestBug017ModuleEvidenceOrder(t *testing.T) {
	engine := NewEngine()
	first := createBugRecord(t, engine, "order-one", "order")
	second := createBugRecord(t, engine, "order-two", "order")
	if len(first.Evidence) == 0 || len(second.Evidence) == 0 {
		t.Fatal("expected evidence")
	}
	if first.Evidence[0].Module != second.Evidence[0].Module {
		t.Fatalf("module order changed: first=%s second=%s", first.Evidence[0].Module, second.Evidence[0].Module)
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
