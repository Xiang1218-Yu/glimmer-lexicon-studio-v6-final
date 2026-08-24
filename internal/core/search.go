package core

import (
	"fmt"
	"strings"
	"time"
)

type SearchModule struct {
	key      string
	priority int
	family   string
	focus    string
}

func NewSearchModule() Module {
	return SearchModule{key: "search", priority: 57, family: "lexicon", focus: "semantic search normalization"}
}

func (m SearchModule) Key() string {
	return m.key
}

func (m SearchModule) Description() string {
	return m.focus + " (filter)"
}

func (m SearchModule) Priority() int {
	return m.priority
}

func (m SearchModule) Family() string {
	return m.family
}

func (m SearchModule) Enabled() bool {
	key := strings.TrimSpace(m.key)
	family := strings.TrimSpace(m.family)
	if key == "" || family == "" {
		return false
	}
	return m.priority > 0
}

func (m SearchModule) Validate(input string) error {
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

func (m SearchModule) Rewrite(input string) string {
	value := strings.TrimSpace(input)
	value = strings.Join(strings.Fields(value), " ")
	marker := fmt.Sprintf("[%s:%d]", m.key, m.priority)
	if strings.Contains(value, marker) {
		return value
	}
	return value + " " + marker
}

func (m SearchModule) Transition(stage string) bool {
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

func (m SearchModule) Evidence(input string, now time.Time) Evidence {
	normalized := m.Rewrite(input)
	return Evidence{
		Module: m.key,
		Digest: fmt.Sprintf("%s-%d-%d", m.key, m.priority, len(normalized)),
		Detail: normalized,
		At:     now.UTC(),
	}
}

func (m SearchModule) Score(input string) int {
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
