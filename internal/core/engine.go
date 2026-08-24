package core

import (
	"context"
	"errors"
	"sort"
	"strings"
	"sync"
	"time"
)

type Module interface {
	Key() string
	Description() string
	Priority() int
	Family() string
	Enabled() bool
	Validate(string) error
	Rewrite(string) string
	Transition(string) bool
	Evidence(string, time.Time) Evidence
	Score(string) int
}

type Evidence struct {
	Module string    `json:"module"`
	Digest string    `json:"digest"`
	Detail string    `json:"detail"`
	At     time.Time `json:"at"`
}

type Record struct {
	ID        string       `json:"id"`
	Stage     string       `json:"stage"`
	Payload   string       `json:"payload"`
	Score     int          `json:"score"`
	Version   int64        `json:"version"`
	Evidence  []Evidence   `json:"evidence"`
	History   []Transition `json:"history"`
	CreatedAt time.Time    `json:"created_at"`
	UpdatedAt time.Time    `json:"updated_at"`
}

type Transition struct {
	From string    `json:"from"`
	To   string    `json:"to"`
	By   string    `json:"by"`
	At   time.Time `json:"at"`
}

type Engine struct {
	mu       sync.RWMutex
	modules  []Module
	records  map[string]Record
	sequence uint64
	ready    bool
}

// NewEngine returns an Engine whose module registry has not yet been
// assembled. The returned engine is not ready: callers must run Init (in the
// main goroutine or a background one) before record operations are meaningful.
// Keeping construction and initialization separate lets the HTTP layer serve a
// stable "not ready" response during startup instead of returning an empty
// module list that monitoring cannot distinguish from a crash.
func NewEngine() *Engine {
	return &Engine{records: make(map[string]Record)}
}

// Init assembles and sorts the module registry and marks the engine ready. It
// is safe to call from a goroutine started before the HTTP server begins
// serving; once it returns, Ready reports true. Calling Init more than once is
// a no-op after the first successful assembly.
func (e *Engine) Init() {
	modules := []Module{
		NewConceptModule(),
		NewTermModule(),
		NewVariantModule(),
		NewContextModule(),
		NewProposalModule(),
		NewReviewModule(),
		NewRelationModule(),
		NewSynonymModule(),
		NewAntonymModule(),
		NewForbiddenModule(),
		NewPreferredModule(),
		NewDeprecatedModule(),
		NewReplacementModule(),
		NewLocaleModule(),
		NewScriptModule(),
		NewGrammarModule(),
		NewGenderModule(),
		NewNumberModule(),
		NewToneModule(),
		NewBrandModule(),
		NewProductModule(),
		NewFeatureModule(),
		NewInterfaceModule(),
		NewErrorModule(),
		NewCommandModule(),
		NewMeasurementModule(),
		NewUnitModule(),
		NewDateModule(),
		NewNumberformatModule(),
		NewPunctuationModule(),
		NewCapitalizationModule(),
		NewAbbreviationModule(),
		NewAcronymModule(),
		NewTokenModule(),
		NewPhraseModule(),
		NewExampleModule(),
		NewSourceModule(),
		NewTranslatorModule(),
		NewReviewerModule(),
		NewCommentModule(),
		NewDecisionModule(),
		NewStatusModule(),
		NewVersionModule(),
		NewReleaseModule(),
		NewBundleModule(),
		NewManifestModule(),
		NewChecksumModule(),
		NewImporterModule(),
		NewExporterModule(),
		NewMemoryModule(),
		NewMatchModule(),
		NewAmbiguityModule(),
		NewExceptionModule(),
		NewExpiryModule(),
		NewNotificationModule(),
		NewHistoryModule(),
		NewSearchModule(),
		NewAccessModule(),
		NewWorkspaceModule(),
		NewLanguageModule(),
		NewScriptmapModule(),
		NewQualityModule(),
	}
	sort.Slice(modules, func(i, j int) bool { return modules[i].Priority() < modules[j].Priority() })
	e.mu.Lock()
	if !e.ready {
		e.modules = modules
		e.ready = true
	}
	e.mu.Unlock()
}

// Ready reports whether Init has completed and the module registry is available.
func (e *Engine) Ready() bool {
	e.mu.RLock()
	defer e.mu.RUnlock()
	return e.ready
}

func (e *Engine) Modules() []Module {
	e.mu.RLock()
	defer e.mu.RUnlock()
	return append([]Module(nil), e.modules...)
}

func (e *Engine) Create(ctx context.Context, id, payload, actor string) (Record, error) {
	if err := contextError(ctx); err != nil {
		return Record{}, err
	}
	id = strings.TrimSpace(id)
	if id == "" || strings.TrimSpace(payload) == "" {
		return Record{}, errors.New("id and payload are required")
	}
	e.mu.Lock()
	defer e.mu.Unlock()
	if _, exists := e.records[id]; exists {
		return Record{}, errors.New("record already exists")
	}
	now := time.Now().UTC()
	result := Record{ID: id, Stage: "draft", Payload: strings.TrimSpace(payload), Version: 1, CreatedAt: now, UpdatedAt: now}
	for _, module := range e.modules {
		if !module.Enabled() {
			continue
		}
		if err := module.Validate(result.Payload); err != nil {
			return Record{}, err
		}
		result.Payload = module.Rewrite(result.Payload)
		result.Score += module.Score(result.Payload)
		result.Evidence = append(result.Evidence, module.Evidence(result.Payload, now))
	}
	result.History = append(result.History, Transition{From: "", To: result.Stage, By: actor, At: now})
	e.records[id] = cloneRecord(result)
	return cloneRecord(result), nil
}

func (e *Engine) Get(ctx context.Context, id string) (Record, error) {
	if err := contextError(ctx); err != nil {
		return Record{}, err
	}
	e.mu.RLock()
	defer e.mu.RUnlock()
	record, ok := e.records[strings.TrimSpace(id)]
	if !ok {
		return Record{}, errors.New("record not found")
	}
	return cloneRecord(record), nil
}

func (e *Engine) Advance(ctx context.Context, id, stage, actor string) (Record, error) {
	if err := contextError(ctx); err != nil {
		return Record{}, err
	}
	e.mu.Lock()
	defer e.mu.Unlock()
	record, ok := e.records[strings.TrimSpace(id)]
	if !ok {
		return Record{}, errors.New("record not found")
	}
	stage = strings.ToLower(strings.TrimSpace(stage))
	if stage == "" || stage == record.Stage {
		return Record{}, errors.New("a different target stage is required")
	}
	for _, module := range e.modules {
		if !module.Transition(stage) {
			return Record{}, errors.New("stage is not accepted by module " + module.Key())
		}
	}
	now := time.Now().UTC()
	record.History = append(record.History, Transition{From: record.Stage, To: stage, By: actor, At: now})
	record.Stage = stage
	record.Version++
	record.UpdatedAt = now
	record.Payload = strings.TrimSpace(record.Payload)
	e.records[record.ID] = cloneRecord(record)
	return cloneRecord(record), nil
}

func (e *Engine) List(ctx context.Context, stage string, limit int) ([]Record, error) {
	if err := contextError(ctx); err != nil {
		return nil, err
	}
	e.mu.RLock()
	defer e.mu.RUnlock()
	result := make([]Record, 0, len(e.records))
	for _, record := range e.records {
		if stage != "" && record.Stage != stage {
			continue
		}
		result = append(result, cloneRecord(record))
	}
	sort.Slice(result, func(i, j int) bool { return result[i].UpdatedAt.After(result[j].UpdatedAt) })
	if limit > 0 && len(result) > limit {
		result = result[:limit]
	}
	return result, nil
}

func (e *Engine) ValidatePayload(ctx context.Context, payload string) (int, []string, error) {
	if err := contextError(ctx); err != nil {
		return 0, nil, err
	}
	score := 0
	warnings := make([]string, 0)
	for _, module := range e.modules {
		if err := module.Validate(payload); err != nil {
			return 0, warnings, err
		}
		score += module.Score(payload)
		if len(payload) < module.Priority() {
			warnings = append(warnings, module.Key()+" expects richer context")
		}
	}
	return score, warnings, nil
}

func (e *Engine) Snapshot() map[string]int {
	e.mu.RLock()
	defer e.mu.RUnlock()
	snapshot := map[string]int{"modules": len(e.modules), "records": len(e.records)}
	for _, record := range e.records {
		snapshot["stage_"+record.Stage]++
	}
	return snapshot
}

func contextError(ctx context.Context) error {
	if ctx == nil {
		return errors.New("context is required")
	}
	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
		return nil
	}
}

func cloneRecord(record Record) Record {
	record.Evidence = append([]Evidence(nil), record.Evidence...)
	record.History = append([]Transition(nil), record.History...)
	return record
}
