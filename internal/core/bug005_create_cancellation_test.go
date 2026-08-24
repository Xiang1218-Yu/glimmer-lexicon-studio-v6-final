package core

import (
	"context"
	"fmt"
	"strings"
	"testing"
	"time"
)

func TestBug005CreateCancellation(t *testing.T) {
	engine := NewEngine()
	ctx, cancel := context.WithCancel(context.Background())
	done := make(chan error, 1)
	go func() {
		_, err := engine.Create(ctx, "cancel-term", "context-window", "tester")
		done <- err
	}()
	time.Sleep(3 * time.Millisecond)
	cancel()
	if err := <-done; err == nil {
		t.Fatal("expected cancellation to stop the create call")
	}
	if _, err := engine.Get(context.Background(), "cancel-term"); err == nil {
		t.Fatal("canceled create left a stored record")
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
