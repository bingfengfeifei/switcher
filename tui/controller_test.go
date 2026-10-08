package tui

import (
	"path/filepath"
	"testing"

	tea "github.com/charmbracelet/bubbletea"
)

type copyTestPaths struct {
	root string
}

func (paths copyTestPaths) GetAppConfigPath() string {
	return filepath.Join(paths.root, "config.json")
}

func (paths copyTestPaths) GetClaudeConfigDir() string {
	return filepath.Join(paths.root, "claude")
}

func (paths copyTestPaths) GetCodexConfigDir() string {
	return filepath.Join(paths.root, "codex")
}

func (paths copyTestPaths) GetDroidConfigDir() string {
	return filepath.Join(paths.root, "droid")
}

func useCopyTestPaths(t *testing.T) {
	t.Helper()
	oldPaths := platformPaths
	platformPaths = copyTestPaths{root: t.TempDir()}
	t.Cleanup(func() {
		platformPaths = oldPaths
	})
}

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
	m.formField = 6
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

func TestNextCopyNameAvoidsExistingNames(t *testing.T) {
	taken := map[string]bool{"prod (copy)": true, "prod (copy 2)": true}

	got := nextCopyName("prod", func(name string) bool {
		return taken[name]
	})
	if got != "prod (copy 3)" {
		t.Fatalf("nextCopyName() = %q, want %q", got, "prod (copy 3)")
	}
}

func TestCopyKeyDuplicatesSelectedClaudeConfig(t *testing.T) {
	useCopyTestPaths(t)
	config := &Config{
		ClaudeCode: []ServiceConfig{{
			Name:     "prod",
			Provider: "switcher",
			BaseURL:  "https://api.example.com",
			APIKey:   "secret",
			Model:    "model",
		}},
		Active: ActiveConfig{ClaudeCode: 0},
	}
	m := model{config: config, state: claudeCodeList}
	m.sortClaudeCodeConfigs()

	updated, _ := m.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune("c")})
	got := updated.(model)

	if len(got.config.ClaudeCode) != 2 {
		t.Fatalf("Claude config count after copy = %d, want 2", len(got.config.ClaudeCode))
	}
	copied := got.config.ClaudeCode[1]
	if copied.Name != "prod (copy)" || copied.BaseURL != "https://api.example.com" || copied.APIKey != "secret" || copied.Model != "model" {
		t.Fatalf("copied config = %+v", copied)
	}
	if got.config.Active.ClaudeCode != 0 {
		t.Fatalf("active Claude index changed to %d, want 0", got.config.Active.ClaudeCode)
	}
	if len(got.sortedClaudeCode) != 2 {
		t.Fatalf("sorted Claude config count = %d, want 2", len(got.sortedClaudeCode))
	}
	wantStatus := translations[currentLang]["success_copy"]
	if got.error != wantStatus {
		t.Fatalf("copy status = %q, want %q", got.error, wantStatus)
	}
}

func TestCopyConfigMethodsDuplicateCodexAndDroid(t *testing.T) {
	useCopyTestPaths(t)
	config := &Config{
		Codex: []ServiceConfig{{
			Name:    "codex",
			BaseURL: "https://codex.example.com",
			APIKey:  "secret",
		}},
		Droid: []DroidConfig{{
			ModelDisplayName: "droid",
			Model:            "model",
			BaseURL:          "https://droid.example.com",
			APIKey:           "secret",
		}},
		Active: ActiveConfig{Codex: 0, Droid: 0},
	}

	if err := config.CopyCodexConfig(0); err != nil {
		t.Fatalf("CopyCodexConfig() error = %v", err)
	}
	if err := config.CopyDroidConfig(0); err != nil {
		t.Fatalf("CopyDroidConfig() error = %v", err)
	}

	if config.Codex[1].Name != "codex (copy)" || config.Codex[1].Provider != "switcher" {
		t.Fatalf("copied Codex config = %+v", config.Codex[1])
	}
	if config.Droid[1].ModelDisplayName != "droid (copy)" || config.Droid[1].Model != "model" {
		t.Fatalf("copied Droid config = %+v", config.Droid[1])
	}
	if config.Active.Codex != 0 || config.Active.Droid != 0 {
		t.Fatalf("active indices changed: codex=%d droid=%d, want 0 and 0", config.Active.Codex, config.Active.Droid)
	}
}
