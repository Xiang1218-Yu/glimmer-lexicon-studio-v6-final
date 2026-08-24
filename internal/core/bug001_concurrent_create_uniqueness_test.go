package core

import (
	"context"
	"fmt"
	"strings"
	"sync"
	"testing"
	"time"
)

func TestBug001ConcurrentCreateUniqueness(t *testing.T) {
	engine := NewEngine()
	const callers = 20
	var wg sync.WaitGroup
	var mu sync.Mutex
	successes := 0
	wg.Add(callers)
	for i := 0; i < callers; i++ {
		go func() {
			defer wg.Done()
			if _, err := engine.Create(context.Background(), "shared-term", "race-create", "tester"); err == nil {
				mu.Lock()
				successes++
				mu.Unlock()
			}
		}()
	}
	wg.Wait()
	if successes != 1 {
		t.Fatalf("expected one successful create, got %d", successes)
	}
	if got, err := engine.Get(context.Background(), "shared-term"); err != nil || got.ID != "shared-term" {
		t.Fatalf("stored record is inconsistent: %v %#v", err, got)
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
