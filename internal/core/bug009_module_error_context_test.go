package core

import (
	"context"
	"fmt"
	"strings"
	"testing"
	"time"
)

func TestBug009ModuleErrorContext(t *testing.T) {
	engine := NewEngine()
	_, err := engine.Create(context.Background(), "error-term", "module-error", "tester")
	if err == nil {
		t.Fatal("expected module validation error")
	}
	if !strings.Contains(err.Error(), "concept") {
		t.Fatalf("module context was lost: %v", err)
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
