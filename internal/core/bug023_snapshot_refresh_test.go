package core

import (
	"context"
	"fmt"
	"strings"
	"testing"
	"time"
)

func TestBug023SnapshotRefresh(t *testing.T) {
	engine := NewEngine()
	first := engine.Snapshot()
	if first["records"] != 0 {
		t.Fatal("expected empty initial snapshot")
	}
	createBugRecord(t, engine, "refresh-snapshot", "snapshot")
	second := engine.Snapshot()
	if second["records"] != 1 {
		t.Fatalf("snapshot did not refresh: %#v", second)
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
