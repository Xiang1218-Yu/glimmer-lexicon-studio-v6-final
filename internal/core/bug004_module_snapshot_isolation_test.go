package core

import (
	"context"
	"fmt"
	"strings"
	"testing"
	"time"
)

func TestBug004ModuleSnapshotIsolation(t *testing.T) {
	engine := NewEngine()
	modules := engine.Modules()
	if len(modules) == 0 {
		t.Fatal("expected modules")
	}
	modules[0] = nil
	defer func() {
		if recover() != nil {
			t.Fatal("mutating a returned module snapshot crashed the engine")
		}
	}()
	if _, err := engine.Create(context.Background(), "module-term", "module snapshot", "tester"); err != nil {
		t.Fatalf("create after snapshot mutation: %v", err)
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
