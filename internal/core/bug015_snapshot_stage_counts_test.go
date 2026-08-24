package core

import (
	"context"
	"fmt"
	"strings"
	"testing"
	"time"
)

func TestBug015SnapshotStageCounts(t *testing.T) {
	engine := NewEngine()
	createBugRecord(t, engine, "snapshot-stage", "snapshot")
	if _, err := engine.Advance(context.Background(), "snapshot-stage", "review", "reviewer"); err != nil {
		t.Fatal(err)
	}
	snapshot := engine.Snapshot()
	if snapshot["stage_review"] != 1 || snapshot["stage_draft"] != 0 {
		t.Fatalf("snapshot counts are stale: %#v", snapshot)
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
