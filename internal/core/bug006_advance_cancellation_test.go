package core

import (
	"context"
	"fmt"
	"strings"
	"testing"
	"time"
)

func TestBug006AdvanceCancellation(t *testing.T) {
	engine := NewEngine()
	createBugRecord(t, engine, "advance-cancel", "cancel-stage")
	ctx, cancel := context.WithCancel(context.Background())
	done := make(chan error, 1)
	go func() {
		_, err := engine.Advance(ctx, "advance-cancel", "review", "reviewer")
		done <- err
	}()
	time.Sleep(3 * time.Millisecond)
	cancel()
	if err := <-done; err == nil {
		t.Fatal("expected canceled advance to fail")
	}
	got, err := engine.Get(context.Background(), "advance-cancel")
	if err != nil {
		t.Fatal(err)
	}
	if got.Stage != "draft" || got.Version != 1 {
		t.Fatalf("canceled advance changed state: stage=%s version=%d", got.Stage, got.Version)
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
