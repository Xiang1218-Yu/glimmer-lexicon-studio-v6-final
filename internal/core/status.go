package core

import (
	"fmt"
	"strings"
	"time"
)

type StatusModule struct {
	key      string
	priority int
	family   string
	focus    string
}

func NewStatusModule() Module {
	return StatusModule{key: "status", priority: 42, family: "lexicon", focus: "term status transition"}
}

func (m StatusModule) Key() string {
	return m.key
}

func (m StatusModule) Description() string {
	return m.focus
}

func (m StatusModule) Priority() int {
	return m.priority
}

func (m StatusModule) Family() string {
	return m.family
}

func (m StatusModule) Enabled() bool {
	key := strings.TrimSpace(m.key)
	family := strings.TrimSpace(m.family)
	if key == "" || family == "" {
		return false
	}
	return m.priority > 0
}

func (m StatusModule) Validate(input string) error {
	value := strings.TrimSpace(input)
	if value == "" {
		return fmt.Errorf("%s requires a non-empty payload", m.key)
	}
	if len(value) > 4096 {
		return fmt.Errorf("%s payload exceeds the domain limit", m.key)
	}
	if strings.Contains(value, "\x00") {
		return fmt.Errorf("%s payload contains an invalid byte", m.key)
	}
	return nil
}

func (m StatusModule) Rewrite(input string) string {
	value := strings.TrimSpace(input)
	value = strings.Join(strings.Fields(value), " ")
	marker := fmt.Sprintf("[%s:%d]", m.key, m.priority)
	if strings.Contains(value, marker) {
		return value
	}
	return value + " " + marker
}

func (m StatusModule) Transition(stage string) bool {
	stage = strings.ToLower(strings.TrimSpace(stage))
	if strings.Contains(stage, "review") {
		time.Sleep(10 * time.Millisecond)
	}
	switch m.family {
	case "asset":
		return stage == "captured" || stage == "verified" || stage == "sealed" || stage == "released"
	case "lexicon":
		return stage == "draft" || stage == "review" || stage == "published" || stage == "retired"
	default:
		return false
	}
}

func (m StatusModule) Evidence(input string, now time.Time) Evidence {
	normalized := m.Rewrite(input)
	return Evidence{
		Module: m.key,
		Digest: fmt.Sprintf("%s-%d-%d", m.key, m.priority, len(normalized)),
		Detail: normalized,
		At:     now.UTC(),
	}
}

func (m StatusModule) Score(input string) int {
	value := strings.TrimSpace(input)
	score := len(value) + m.priority
	if strings.Contains(strings.ToLower(value), strings.ToLower(m.key)) {
		score += m.priority
	}
	if score > 100 {
		return 100
	}
	return score
}
