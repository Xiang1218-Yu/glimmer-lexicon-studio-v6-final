package core

import (
	"context"
	"fmt"
	"strings"
	"testing"
	"time"
)

func TestBug019StoredRewrittenPayload(t *testing.T) {
	engine := NewEngine()
	created := createBugRecord(t, engine, "payload-term", "rewritten")
	if !strings.Contains(created.Payload, "[concept:1]") {
		t.Fatalf("expected rewritten create result: %s", created.Payload)
	}
	stored, err := engine.Get(context.Background(), "payload-term")
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(stored.Payload, "[concept:1]") {
		t.Fatalf("stored payload lost module rewrite: %s", stored.Payload)
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
