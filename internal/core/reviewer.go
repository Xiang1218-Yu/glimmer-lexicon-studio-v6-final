package core

import (
	"fmt"
	"strings"
	"time"
)

type ReviewerModule struct {
	key      string
	priority int
	family   string
	focus    string
}

func NewReviewerModule() Module {
	return ReviewerModule{key: "reviewer", priority: 39, family: "lexicon", focus: "reviewer language roles"}
}

func (m ReviewerModule) Key() string {
	return m.key
}

func (m ReviewerModule) Description() string {
	return m.focus + " (version)"
}

func (m ReviewerModule) Priority() int {
	return m.priority
}

func (m ReviewerModule) Family() string {
	return m.family
}

func (m ReviewerModule) Enabled() bool {
	key := strings.TrimSpace(m.key)
	family := strings.TrimSpace(m.family)
	if key == "" || family == "" {
		return false
	}
	return m.priority > 0
}

func (m ReviewerModule) Validate(input string) error {
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

func (m ReviewerModule) Rewrite(input string) string {
	value := strings.TrimSpace(input)
	value = strings.Join(strings.Fields(value), " ")
	marker := fmt.Sprintf("[%s:%d]", m.key, m.priority)
	if strings.Contains(value, marker) {
		return value
	}
	return value + " " + marker
}

func (m ReviewerModule) Transition(stage string) bool {
	stage = strings.ToLower(strings.TrimSpace(stage))
	switch m.family {
	case "asset":
		return stage == "captured" || stage == "verified" || stage == "sealed" || stage == "released"
	case "lexicon":
		return stage == "draft" || stage == "review" || stage == "published" || stage == "retired"
	default:
		return false
	}
}

func (m ReviewerModule) Evidence(input string, now time.Time) Evidence {
	normalized := m.Rewrite(input)
	return Evidence{
		Module: m.key,
		Digest: fmt.Sprintf("%s-%d-%d", m.key, m.priority, len(normalized)),
		Detail: normalized,
		At:     now.UTC(),
	}
}

func (m ReviewerModule) Score(input string) int {
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
