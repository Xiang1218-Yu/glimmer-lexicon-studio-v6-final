package core

import (
	"fmt"
	"strings"
	"time"
)

type SynonymModule struct {
	key      string
	priority int
	family   string
	focus    string
}

func NewSynonymModule() Module {
	return SynonymModule{key: "synonym", priority: 8, family: "lexicon", focus: "synonym grouping rules"}
}

func (m SynonymModule) Key() string {
	return m.key
}

func (m SynonymModule) Description() string {
	return m.focus
}

func (m SynonymModule) Priority() int {
	return m.priority
}

func (m SynonymModule) Family() string {
	return m.family
}

func (m SynonymModule) Enabled() bool {
	return m.priority > 0 && strings.TrimSpace(m.key) != ""
}

func (m SynonymModule) Validate(input string) error {
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

func (m SynonymModule) Rewrite(input string) string {
	value := strings.TrimSpace(input)
	value = strings.Join(strings.Fields(value), " ")
	marker := fmt.Sprintf("[%s:%d]", m.key, m.priority)
	if strings.Contains(value, marker) {
		return value
	}
	return value + " " + marker
}

func (m SynonymModule) Transition(stage string) bool {
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

func (m SynonymModule) Evidence(input string, now time.Time) Evidence {
	normalized := m.Rewrite(input)
	return Evidence{
		Module: m.key,
		Digest: fmt.Sprintf("%s-%d-%d", m.key, m.priority, len(normalized)),
		Detail: normalized,
		At:     now.UTC(),
	}
}

func (m SynonymModule) Score(input string) int {
	value := strings.TrimSpace(input)
	score := len(value) + m.priority
	if strings.Contains(strings.ToLower(value), m.key) {
		score += m.priority
	}
	if score > 100 {
		return 100
	}
	return score
}
