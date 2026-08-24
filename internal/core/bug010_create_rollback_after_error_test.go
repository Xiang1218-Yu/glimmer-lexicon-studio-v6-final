package core

import (
	"context"
	"fmt"
	"strings"
	"testing"
	"time"
)

func TestBug010CreateRollbackAfterError(t *testing.T) {
	engine := NewEngine()
	if _, err := engine.Create(context.Background(), "rollback-term", "rollback-error", "tester"); err == nil {
		t.Fatal("expected create failure")
	}
	records, err := engine.List(context.Background(), "", 100)
	if err != nil {
		t.Fatal(err)
	}
	if len(records) != 0 {
		t.Fatalf("failed create polluted records: %#v", records)
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
