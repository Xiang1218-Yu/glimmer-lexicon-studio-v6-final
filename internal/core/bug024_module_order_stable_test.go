package core

import (
	"context"
	"fmt"
	"strings"
	"testing"
	"time"
)

func TestBug024ModuleOrderStable(t *testing.T) {
	engine := NewEngine()
	for i := 0; i < 4; i++ {
		record := createBugRecord(t, engine, fmt.Sprintf("order-%d", i), "stable order")
		if len(record.Evidence) == 0 {
			t.Fatal("expected evidence")
		}
		if record.Evidence[0].Module != "concept" {
			t.Fatalf("unexpected first module on run %d: %s", i, record.Evidence[0].Module)
		}
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
