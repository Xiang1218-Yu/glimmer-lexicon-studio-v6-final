package core

import (
	"context"
	"fmt"
	"strings"
	"testing"
	"time"
)

func TestBug022ListRecordIsolation(t *testing.T) {
	engine := NewEngine()
	createBugRecord(t, engine, "list-record", "list record")
	records, err := engine.List(context.Background(), "", 100)
	if err != nil {
		t.Fatal(err)
	}
	records[0].Evidence[0].Detail = "changed-by-page"
	got, _ := engine.Get(context.Background(), "list-record")
	if got.Evidence[0].Detail == "changed-by-page" {
		t.Fatal("list result mutation changed stored record")
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
