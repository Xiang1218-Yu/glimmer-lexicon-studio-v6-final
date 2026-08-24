package core

import (
	"context"
	"fmt"
	"strings"
	"sync"
	"testing"
	"time"
)

func TestBug002ConcurrentAdvanceHistory(t *testing.T) {
	engine := NewEngine()
	createBugRecord(t, engine, "advance-term", "race-advance")
	var wg sync.WaitGroup
	wg.Add(2)
	for _, stage := range []string{"review", "published"} {
		stage := stage
		go func() {
			defer wg.Done()
			if _, err := engine.Advance(context.Background(), "advance-term", stage, "reviewer"); err != nil {
				t.Errorf("advance %s: %v", stage, err)
			}
		}()
	}
	wg.Wait()
	got, err := engine.Get(context.Background(), "advance-term")
	if err != nil {
		t.Fatal(err)
	}
	if got.Version != 3 || len(got.History) != 3 {
		t.Fatalf("lost concurrent transition: version=%d history=%d", got.Version, len(got.History))
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
