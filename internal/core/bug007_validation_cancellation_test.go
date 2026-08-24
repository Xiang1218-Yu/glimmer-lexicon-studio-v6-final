package core

import (
	"context"
	"fmt"
	"strings"
	"testing"
	"time"
)

func TestBug007ValidationCancellation(t *testing.T) {
	engine := NewEngine()
	ctx, cancel := context.WithCancel(context.Background())
	done := make(chan error, 1)
	go func() {
		_, _, err := engine.ValidatePayload(ctx, "validate-window")
		done <- err
	}()
	time.Sleep(3 * time.Millisecond)
	cancel()
	if err := <-done; err == nil {
		t.Fatal("expected validation cancellation")
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
