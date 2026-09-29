package tui

import (
	"encoding/json"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"testing"
)

type recordingShellManager struct {
	key   string
	value string
}

func (m *recordingShellManager) SetEnvVar(key, value string) error {
	m.key = key
	m.value = value
	return nil
}

func TestSwitchCodexKeepsEndpointAndKeyTogether(t *testing.T) {
	useCopyTestPaths(t)
	codexDir := platformPaths.GetCodexConfigDir()
	if err := os.MkdirAll(codexDir, 0755); err != nil {
		t.Fatal(err)
	}
	configPath := filepath.Join(codexDir, "config.toml")
	if err := os.WriteFile(configPath, []byte("[projects.\"/tmp\"]\ntrust_level = \"trusted\"\n"), 0644); err != nil {
		t.Fatal(err)
	}

	configs := []ServiceConfig{
		{Provider: "switcher", BaseURL: "https://old.example/v1", APIKey: "sk-old", Model: "gpt-5.5"},
		{Provider: "switcher", BaseURL: "https://new.example/v1", APIKey: "sk-new", Model: "gpt-5.5"},
	}
	for i := range configs {
		if err := (&Config{}).SwitchCodex(&configs[i]); err != nil {
			t.Fatalf("SwitchCodex(%d): %v", i, err)
		}
	}

	content, err := os.ReadFile(configPath)
	if err != nil {
		t.Fatal(err)
	}
	for _, want := range []string{
		`base_url = "https://new.example/v1"`,
		`experimental_bearer_token = "sk-new"`,
		"requires_openai_auth = false",
		`[projects."/tmp"]`,
	} {
		if !strings.Contains(string(content), want) {
			t.Errorf("config.toml does not contain %q", want)
		}
	}
	if strings.Contains(string(content), "sk-old") || strings.Contains(string(content), "https://old.example/v1") {
		t.Error("config.toml still contains the previous provider credentials")
	}

	authPath := filepath.Join(codexDir, "auth.json")
	authData, err := os.ReadFile(authPath)
	if err != nil {
		t.Fatal(err)
	}
	var auth CodexAuth
	if err := json.Unmarshal(authData, &auth); err != nil {
		t.Fatal(err)
	}
	if auth.AuthMode != "apikey" || auth.OPENAI_API_KEY != "sk-new" {
		t.Error("auth.json does not contain the selected API key and auth mode")
	}
	if runtime.GOOS != "windows" {
		for _, path := range []string{configPath, authPath} {
			info, err := os.Stat(path)
			if err != nil {
				t.Fatal(err)
			}
			if got := info.Mode().Perm(); got != 0600 {
				t.Errorf("%s permissions = %o, want 600", path, got)
			}
		}
	}
}

func TestSwitchCodexEnvAuthDoesNotWriteBearerToken(t *testing.T) {
	useCopyTestPaths(t)
	manager := &recordingShellManager{}
	oldManager := shellManager
	shellManager = manager
	t.Cleanup(func() { shellManager = oldManager })

	config := ServiceConfig{
		Provider:   "switcher",
		BaseURL:    "https://example.com/v1",
		APIKey:     "sk-env",
		AuthMethod: "env",
		EnvKey:     "CODEX_KEY",
	}
	if err := (&Config{}).SwitchCodex(&config); err != nil {
		t.Fatal(err)
	}
	content, err := os.ReadFile(filepath.Join(platformPaths.GetCodexConfigDir(), "config.toml"))
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(string(content), `env_key = "CODEX_KEY"`) || strings.Contains(string(content), "experimental_bearer_token") {
		t.Error("env authentication did not select only env_key")
	}
	if !strings.Contains(string(content), "requires_openai_auth = false") {
		t.Error("env authentication still requires global OpenAI auth")
	}
	if manager.key != "CODEX_KEY" || manager.value != "sk-env" {
		t.Error("selected key was not saved to the shell environment")
	}
}
