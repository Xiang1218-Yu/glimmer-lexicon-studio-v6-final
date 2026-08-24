package core

import (
	"context"
	"fmt"
	"strings"
	"testing"
	"time"
)

func TestBug016ValidationWarningsIsolation(t *testing.T) {
	engine := NewEngine()
	_, warnings, err := engine.ValidatePayload(context.Background(), "x")
	if err != nil || len(warnings) == 0 {
		t.Fatalf("expected first validation warnings: %v %#v", err, warnings)
	}
	_, warnings, err = engine.ValidatePayload(context.Background(), strings.Repeat("complete ", 80))
	if err != nil {
		t.Fatal(err)
	}
	if len(warnings) != 0 {
		t.Fatalf("warnings leaked across requests: %#v", warnings)
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
