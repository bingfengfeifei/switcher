package tui

import (
	"testing"

	tea "github.com/charmbracelet/bubbletea"
)

func TestHandleInputMapsClaudeFieldsToVisibleRows(t *testing.T) {
	m := model{state: addClaudeCode}

	m.formField = 1
	updated, _ := m.handleInput("b")
	got := updated.(model)
	if got.formData.BaseURL != "b" {
		t.Fatalf("Claude BaseURL input mapped to %q, want %q", got.formData.BaseURL, "b")
	}
	if got.formData.Provider != "" {
		t.Fatalf("Claude BaseURL input should not touch Provider, got %q", got.formData.Provider)
	}

	m = model{state: addClaudeCode}
	m.formField = 2
	updated, _ = m.handleInput("k")
	got = updated.(model)
	if got.formData.APIKey != "k" {
		t.Fatalf("Claude APIKey input mapped to %q, want %q", got.formData.APIKey, "k")
	}

	m = model{state: addClaudeCode}
	m.formField = 5
	updated, _ = m.handleInput("s")
	got = updated.(model)
	if got.formData.ClaudeDefaultSonnetModel != "s" {
		t.Fatalf("Claude Sonnet input mapped to %q, want %q", got.formData.ClaudeDefaultSonnetModel, "s")
	}
}

func TestHandleInputMapsCodexFieldsToVisibleRows(t *testing.T) {
	m := model{state: addCodex}

	m.formField = 1
	updated, _ := m.handleInput("u")
	got := updated.(model)
	if got.formData.BaseURL != "u" {
		t.Fatalf("Codex BaseURL input mapped to %q, want %q", got.formData.BaseURL, "u")
	}
	if got.formData.Provider != "" {
		t.Fatalf("Codex BaseURL input should not touch Provider, got %q", got.formData.Provider)
	}

	m = model{state: addCodex}
	m.formField = 2
	updated, _ = m.handleInput("k")
	got = updated.(model)
	if got.formData.APIKey != "k" {
		t.Fatalf("Codex APIKey input mapped to %q, want %q", got.formData.APIKey, "k")
	}

	m = model{state: addCodex}
	m.formField = 6
	updated, _ = m.handleInput("x")
	got = updated.(model)
	if got.formData.ModelReasoningEffort != "" {
		t.Fatalf("Codex reasoning field should not accept direct input, got %q", got.formData.ModelReasoningEffort)
	}
}

func TestCodexArrowKeysUpdateSelectedChoiceField(t *testing.T) {
	m := model{
		state:     editCodex,
		formField: FieldAuthMethod,
		formData: ServiceConfig{
			AuthMethod:           "auth.json",
			ModelReasoningEffort: ModelReasoningEffortMedium,
		},
	}

	updatedModel, _ := m.Update(tea.KeyMsg{Type: tea.KeyRight})
	got := updatedModel.(model)
	if got.formData.AuthMethod != "env" {
		t.Fatalf("AuthMethod after right arrow = %q, want %q", got.formData.AuthMethod, "env")
	}
	if got.formData.ModelReasoningEffort != ModelReasoningEffortMedium {
		t.Fatalf("Reasoning changed while editing AuthMethod: got %q", got.formData.ModelReasoningEffort)
	}

	m = model{
		state:     editCodex,
		formField: FieldModelReasoningEffort,
		formData: ServiceConfig{
			AuthMethod:           "auth.json",
			ModelReasoningEffort: ModelReasoningEffortMedium,
		},
	}

	updatedModel, _ = m.Update(tea.KeyMsg{Type: tea.KeyRight})
	got = updatedModel.(model)
	if got.formData.ModelReasoningEffort != ModelReasoningEffortHigh {
		t.Fatalf("Reasoning after right arrow = %q, want %q", got.formData.ModelReasoningEffort, ModelReasoningEffortHigh)
	}
	if got.formData.AuthMethod != "auth.json" {
		t.Fatalf("AuthMethod changed while editing reasoning: got %q", got.formData.AuthMethod)
	}
}
