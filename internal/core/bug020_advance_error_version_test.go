package core

import (
	"context"
	"fmt"
	"strings"
	"testing"
	"time"
)

func TestBug020AdvanceErrorVersion(t *testing.T) {
	engine := NewEngine()
	createBugRecord(t, engine, "version-term", "version")
	before, _ := engine.Get(context.Background(), "version-term")
	if _, err := engine.Advance(context.Background(), "version-term", "invalid-stage", "alice"); err == nil {
		t.Fatal("expected stage error")
	}
	after, _ := engine.Get(context.Background(), "version-term")
	if after.Version != before.Version {
		t.Fatalf("error path changed version: before=%d after=%d", before.Version, after.Version)
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
