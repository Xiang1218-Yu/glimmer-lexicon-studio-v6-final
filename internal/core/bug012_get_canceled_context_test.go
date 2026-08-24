package core

import (
	"context"
	"fmt"
	"strings"
	"testing"
	"time"
)

func TestBug012GetCanceledContext(t *testing.T) {
	engine := NewEngine()
	createBugRecord(t, engine, "get-cancel", "get context")
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	if _, err := engine.Get(ctx, "get-cancel"); err == nil {
		t.Fatal("canceled get returned a record")
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
