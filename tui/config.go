package tui

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// Default model constants
const (
	DefaultCodexModel  = "gpt-5.1-codex"
	DefaultClaudeModel = "claude-sonnet-4-5-20250929"
	DefaultDroidModel  = "gpt-5.1-codex"
)

// Default configuration constants
const (
	DefaultWireAPI              = "responses"
	DefaultModelReasoningEffort = "medium"
	DefaultEnvKey               = "CODEX_KEY"
)

// Model reasoning effort levels
const (
	ModelReasoningEffortLow    = "low"
	ModelReasoningEffortMedium = "medium"
	ModelReasoningEffortHigh   = "high"
)

var platformPaths PlatformPaths
var shellManager ShellManager

func init() {
	var err error
	platformPaths, err = NewPlatformPaths()
	if err != nil {
		panic(fmt.Sprintf("failed to initialize platform paths: %v", err))
	}
	shellManager = NewShellManager()
}

type ServiceConfig struct {
	Name                  string `json:"name"`
	Provider              string `json:"provider"`
	BaseURL               string `json:"base_url"`
	APIKey                string `json:"api_key"`
	Model                 string `json:"model,omitempty"`
	WireAPI               string `json:"wire_api,omitempty"`
	AuthMethod            string `json:"auth_method,omitempty"`
	EnvKey                string `json:"env_key,omitempty"`
	ModelReasoningEffort  string `json:"model_reasoning_effort,omitempty"`
	ClaudeDefaultModel    string `json:"claude_default_model,omitempty"`
}

type DroidConfig struct {
	ModelDisplayName string `json:"model_display_name"`
	Model            string `json:"model"`
	BaseURL          string `json:"base_url"`
	APIKey           string `json:"api_key"`
	Provider         string `json:"provider"`
}

type FactoryConfig struct {
	CustomModels []DroidConfig `json:"custom_models"`
}

type Config struct {
	ClaudeCode []ServiceConfig `json:"claude_code"`
	Codex      []ServiceConfig `json:"codex"`
	Droid      []DroidConfig   `json:"droid"`
	Active     ActiveConfig    `json:"active"`
	Language   string          `json:"language,omitempty"`
}

type ActiveConfig struct {
	ClaudeCode int `json:"claude_code"`
	Codex      int `json:"codex"`
	Droid      int `json:"droid"`
}

type ClaudeSettings struct {
	Env         map[string]string `json:"env"`
	Permissions struct {
		Allow []string `json:"allow"`
		Deny  []string `json:"deny"`
	} `json:"permissions"`
	AlwaysThinkingEnabled bool `json:"alwaysThinkingEnabled"`
}

type CodexConfig struct {
	ModelProvider          string                    `toml:"model_provider"`
	Model                  string                    `toml:"model"`
	ModelReasoningEffort   string                    `toml:"model_reasoning_effort"`
	DisableResponseStorage bool                      `toml:"disable_response_storage"`
	ModelProviders         map[string]CodexProvider  `toml:"model_providers"`
	Projects               map[string]CodexProject   `toml:"projects"`
	MCPServers             map[string]CodexMCPServer `toml:"mcp_servers"`
}

type CodexProvider struct {
	Name               string `toml:"name"`
	BaseURL            string `toml:"base_url"`
	WireAPI            string `toml:"wire_api"`
	EnvKey             string `toml:"env_key"`
	RequiresOpenAIAuth bool   `toml:"requires_openai_auth"`
}

type CodexProject struct {
	TrustLevel string `toml:"trust_level"`
}

type CodexMCPServer struct {
	Command string   `toml:"command"`
	Args    []string `toml:"args"`
}

type CodexAuth struct {
	OPENAI_API_KEY string `json:"OPENAI_API_KEY"`
}

func (c *Config) Save() error {
	configPath := c.getConfigPath()
	configDir := filepath.Dir(configPath)
	if err := mkdirWithPerms(configDir, 0755); err != nil {
		return fmt.Errorf("failed to create config directory: %w", err)
	}

	data, err := json.MarshalIndent(c, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal config: %w", err)
	}

	return writeFileWithPerms(configPath, data, 0644)
}

func (c *Config) Load() error {
	data, err := os.ReadFile(c.getConfigPath())
	if err != nil {
		if os.IsNotExist(err) {
			// Create default config if file doesn't exist
			c.ClaudeCode = []ServiceConfig{}
			c.Codex = []ServiceConfig{}
			c.Droid = []DroidConfig{}
			c.Active = ActiveConfig{ClaudeCode: -1, Codex: -1, Droid: -1}
			c.Language = "zh"
			SetLanguage(c.Language)

			// Import existing configurations
			c.importExistingConfigs()

			return c.Save()
		}
		return fmt.Errorf("failed to read config file: %w", err)
	}

	if err := json.Unmarshal(data, c); err != nil {
		return fmt.Errorf("failed to unmarshal config: %w", err)
	}

	// Migrate old configurations
	c.migrateCodexConfigs()

	// Validate active indices
	if c.Active.ClaudeCode >= len(c.ClaudeCode) {
		c.Active.ClaudeCode = -1
	}
	if c.Active.Codex >= len(c.Codex) {
		c.Active.Codex = -1
	}

	// Initialize language setting
	if c.Language == "" {
		c.Language = "zh"
	}
	SetLanguage(c.Language)

	return nil
}

func (c *Config) importExistingConfigs() {
	// Check if config already exists
	if _, err := os.Stat(c.getConfigPath()); err == nil {
		return // Config already exists, don't import
	}

	// Import Claude Code configuration
	claudeSettingsPath := filepath.Join(platformPaths.GetClaudeConfigDir(), "settings.json")
	if data, err := os.ReadFile(claudeSettingsPath); err == nil {
		var settings ClaudeSettings
		if err := json.Unmarshal(data, &settings); err == nil {
			if authToken, exists := settings.Env["ANTHROPIC_AUTH_TOKEN"]; exists {
				if baseURL, exists := settings.Env["ANTHROPIC_BASE_URL"]; exists {
					claudeConfig := ServiceConfig{
						Name:     "Current Claude",
						Provider: "Current",
						BaseURL:  baseURL,
						APIKey:   authToken,
					}
					c.ClaudeCode = append(c.ClaudeCode, claudeConfig)
					c.Active.ClaudeCode = 0
				}
			}
		}
	}

	// Import Codex configuration
	codexAuthPath := filepath.Join(platformPaths.GetCodexConfigDir(), "auth.json")
	codexConfigPath := filepath.Join(platformPaths.GetCodexConfigDir(), "config.toml")

	if authData, err := os.ReadFile(codexAuthPath); err == nil {
		var auth CodexAuth
		if err := json.Unmarshal(authData, &auth); err == nil && auth.OPENAI_API_KEY != "" {
			// Simple parsing to extract provider info from config.toml
			baseURL := "https://api.openai.com/v1" // default
			providerName := "openai"

			if configData, err := os.ReadFile(codexConfigPath); err == nil {
				content := string(configData)
				// Extract base_url from config
				if idx := strings.Index(content, "base_url = "); idx != -1 {
					start := idx + 13
					if end := strings.Index(content[start:], "\n"); end != -1 {
						baseURL = strings.Trim(strings.TrimSpace(content[start:start+end]), `"`)
					}
				}
				// Extract model_provider
				if idx := strings.Index(content, "model_provider = "); idx != -1 {
					start := idx + 20
					if end := strings.Index(content[start:], "\n"); end != -1 {
						providerName = strings.Trim(strings.TrimSpace(content[start:start+end]), `"`)
					}
				}
			}

			codexConfig := ServiceConfig{
				Name:     "Current Codex",
				Provider: providerName,
				BaseURL:  baseURL,
				APIKey:   auth.OPENAI_API_KEY,
			}
			c.Codex = append(c.Codex, codexConfig)
			c.Active.Codex = 0
		}
	}

	// Import Droid configuration from Factory
	factoryConfigPath := filepath.Join(platformPaths.GetDroidConfigDir(), "config.json")
	if factoryData, err := os.ReadFile(factoryConfigPath); err == nil {
		var factoryConfig FactoryConfig
		if err := json.Unmarshal(factoryData, &factoryConfig); err == nil && len(factoryConfig.CustomModels) > 0 {
			// Import the first custom model as current Droid configuration
			model := factoryConfig.CustomModels[0]
			droidConfig := DroidConfig{
				ModelDisplayName: model.ModelDisplayName,
				Model:            model.Model,
				BaseURL:          model.BaseURL,
				APIKey:           model.APIKey,
				Provider:         model.Provider,
			}
			c.Droid = append(c.Droid, droidConfig)
			c.Active.Droid = 0
		}
	}
}

func (c *Config) getConfigPath() string {
	return platformPaths.GetAppConfigPath()
}

// migrateCodexConfigs migrates old codex configurations to new format with default values
func (c *Config) migrateCodexConfigs() {
	migrated := false
	for i := range c.Codex {
		// Set default values for new fields if they are empty
		if c.Codex[i].Model == "" {
			c.Codex[i].Model = DefaultCodexModel
			migrated = true
		}
		if c.Codex[i].WireAPI == "" {
			c.Codex[i].WireAPI = DefaultWireAPI
			migrated = true
		}
		if c.Codex[i].AuthMethod == "" {
			c.Codex[i].AuthMethod = "auth.json"
			migrated = true
		}
		// Only set EnvKey for env auth method
		if c.Codex[i].AuthMethod == "env" && c.Codex[i].EnvKey == "" {
			c.Codex[i].EnvKey = DefaultEnvKey
			migrated = true
		}
		// Remove EnvKey for auth.json method
		if c.Codex[i].AuthMethod == "auth.json" && c.Codex[i].EnvKey != "" {
			c.Codex[i].EnvKey = ""
			migrated = true
		}
		if c.Codex[i].ModelReasoningEffort == "" {
			c.Codex[i].ModelReasoningEffort = DefaultModelReasoningEffort
			migrated = true
		}
	}

	// Save if any migrations were applied
	if migrated {
		c.Save()
	}
}

func (c *Config) AddClaudeCodeConfig(config ServiceConfig) error {
	c.ClaudeCode = append(c.ClaudeCode, config)
	return c.Save()
}

func (c *Config) AddCodexConfig(config ServiceConfig) error {
	c.Codex = append(c.Codex, config)
	return c.Save()
}

func (c *Config) DeleteClaudeCodeConfig(index int) error {
	if index < 0 || index >= len(c.ClaudeCode) {
		return fmt.Errorf("invalid Claude Code index")
	}

	// Remove the config at index
	c.ClaudeCode = append(c.ClaudeCode[:index], c.ClaudeCode[index+1:]...)

	// Adjust active index if needed
	if c.Active.ClaudeCode == index {
		c.Active.ClaudeCode = -1
	} else if c.Active.ClaudeCode > index {
		c.Active.ClaudeCode--
	}

	return c.Save()
}

func (c *Config) DeleteCodexConfig(index int) error {
	if index < 0 || index >= len(c.Codex) {
		return fmt.Errorf("invalid Codex index")
	}

	// Remove the config at index
	c.Codex = append(c.Codex[:index], c.Codex[index+1:]...)

	// Adjust active index if needed
	if c.Active.Codex == index {
		c.Active.Codex = -1
	} else if c.Active.Codex > index {
		c.Active.Codex--
	}

	return c.Save()
}

func (c *Config) SetActiveClaudeCode(index int) error {
	if index >= 0 && index < len(c.ClaudeCode) {
		c.Active.ClaudeCode = index
		return c.Save()
	}
	return fmt.Errorf("invalid Claude Code index")
}

func (c *Config) SetActiveCodex(index int) error {
	if index >= 0 && index < len(c.Codex) {
		c.Active.Codex = index
		return c.Save()
	}
	return fmt.Errorf("invalid Codex index")
}

func (c *Config) GetActiveClaudeCode() *ServiceConfig {
	if c.Active.ClaudeCode >= 0 && c.Active.ClaudeCode < len(c.ClaudeCode) {
		return &c.ClaudeCode[c.Active.ClaudeCode]
	}
	return nil
}

func (c *Config) GetActiveCodex() *ServiceConfig {
	if c.Active.Codex >= 0 && c.Active.Codex < len(c.Codex) {
		return &c.Codex[c.Active.Codex]
	}
	return nil
}

func (c *Config) SwitchClaudeCode(config *ServiceConfig) error {
	if config == nil {
		return fmt.Errorf("config cannot be nil")
	}

	settingsPath := filepath.Join(platformPaths.GetClaudeConfigDir(), "settings.json")

	settings := ClaudeSettings{
		Env: map[string]string{
			"ANTHROPIC_AUTH_TOKEN":                     config.APIKey,
			"ANTHROPIC_BASE_URL":                       config.BaseURL,
			"CLAUDE_CODE_DISABLE_NONESSENTIAL_TRAFFIC": "1",
		},
		Permissions: struct {
			Allow []string `json:"allow"`
			Deny  []string `json:"deny"`
		}{
			Allow: []string{},
			Deny:  []string{},
		},
		AlwaysThinkingEnabled: false,
	}

	// 如果设置了默认模型，添加三个环境变量
	if config.ClaudeDefaultModel != "" {
		settings.Env["ANTHROPIC_DEFAULT_OPUS_MODEL"] = config.ClaudeDefaultModel
		settings.Env["ANTHROPIC_DEFAULT_SONNET_MODEL"] = config.ClaudeDefaultModel
		settings.Env["ANTHROPIC_DEFAULT_HAIKU_MODEL"] = config.ClaudeDefaultModel
	}

	data, err := json.MarshalIndent(settings, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal Claude settings: %w", err)
	}

	claudeDir := filepath.Dir(settingsPath)
	if err := mkdirWithPerms(claudeDir, 0755); err != nil {
		return fmt.Errorf("failed to create .claude directory: %w", err)
	}

	return writeFileWithPerms(settingsPath, data, 0644)
}

func (c *Config) SwitchCodex(config *ServiceConfig) error {
	if config == nil {
		return fmt.Errorf("config cannot be nil")
	}

	codexDir := platformPaths.GetCodexConfigDir()
	if err := mkdirWithPerms(codexDir, 0755); err != nil {
		return fmt.Errorf("failed to create .codex directory: %w", err)
	}

	// Write auth.json
	auth := CodexAuth{
		OPENAI_API_KEY: config.APIKey,
	}

	authData, err := json.Marshal(auth)
	if err != nil {
		return fmt.Errorf("failed to marshal Codex auth: %w", err)
	}

	authPath := filepath.Join(codexDir, "auth.json")
	if err := writeFileWithPerms(authPath, authData, 0644); err != nil {
		return fmt.Errorf("failed to write auth.json: %w", err)
	}

	// Read existing config.toml to preserve other settings
	configPath := filepath.Join(codexDir, "config.toml")
	existingConfig := c.loadCodexConfig(configPath)

	// Update provider settings
	if existingConfig.ModelProviders == nil {
		existingConfig.ModelProviders = make(map[string]CodexProvider)
	}

	providerName := config.Provider

	// Set model if specified, otherwise use existing or default
	model := config.Model
	if model == "" {
		model = existingConfig.Model
		if model == "" {
			model = DefaultCodexModel
		}
	}
	existingConfig.Model = model

	// Set wire_api if specified, otherwise use existing or default
	wireAPI := config.WireAPI
	if wireAPI == "" {
		wireAPI = DefaultWireAPI
	}

	// Set model_reasoning_effort if specified, otherwise use existing or default
	modelReasoningEffort := config.ModelReasoningEffort
	if modelReasoningEffort == "" {
		modelReasoningEffort = existingConfig.ModelReasoningEffort
		if modelReasoningEffort == "" {
			modelReasoningEffort = DefaultModelReasoningEffort
		}
	}
	existingConfig.ModelReasoningEffort = modelReasoningEffort

	// Set env_key only for env auth method
	envKey := ""
	if config.AuthMethod == "env" {
		envKey = config.EnvKey
		if envKey == "" {
			envKey = DefaultEnvKey
		}
	}

	existingConfig.ModelProviders[providerName] = CodexProvider{
		Name:               providerName,
		BaseURL:            config.BaseURL,
		WireAPI:            wireAPI,
		EnvKey:             envKey,
		RequiresOpenAIAuth: true,
	}

	// Set as active provider
	existingConfig.ModelProvider = providerName

	// Write updated config.toml
	tomlContent := c.generateCodexConfigTOML(existingConfig)
	if err := writeFileWithPerms(configPath, []byte(tomlContent), 0644); err != nil {
		return err
	}

	// Set environment variable in shell configurations based on AuthMethod
	authMethod := config.AuthMethod
	if authMethod == "" {
		authMethod = "auth.json"
	}

	if authMethod == "env" {
		provider := existingConfig.ModelProviders[providerName]
		if provider.EnvKey == "" {
			provider.EnvKey = DefaultEnvKey
		}
		return shellManager.SetEnvVar(provider.EnvKey, config.APIKey)
	}

	return nil
}

func (c *Config) loadCodexConfig(path string) CodexConfig {
	config := CodexConfig{
		ModelProvider:          "openai",
		Model:                  DefaultCodexModel,
		ModelReasoningEffort:   DefaultModelReasoningEffort,
		DisableResponseStorage: false,
		ModelProviders:         make(map[string]CodexProvider),
		Projects:               make(map[string]CodexProject),
		MCPServers:             make(map[string]CodexMCPServer),
	}

	if data, err := os.ReadFile(path); err == nil {
		// Simple TOML parsing - preserve existing settings including MCP servers
		lines := strings.Split(string(data), "\n")
		currentMCP := ""
		for _, raw := range lines {
			line := strings.TrimSpace(raw)
			if line == "" || strings.HasPrefix(line, "#") {
				continue
			}
			// Handle section headers
			if strings.HasPrefix(line, "[") && strings.HasSuffix(line, "]") {
				currentMCP = ""
				if strings.HasPrefix(line, "[projects.") {
					if idx := strings.Index(line, "]"); idx != -1 {
						projectPath := strings.TrimSpace(line[10:idx])
						if projectPath != "" {
							if _, ok := config.Projects[projectPath]; !ok {
								config.Projects[projectPath] = CodexProject{TrustLevel: "trusted"}
							}
						}
					}
				} else if strings.HasPrefix(line, "[mcp_servers.") {
					if idx := strings.Index(line, "]"); idx != -1 {
						serverName := strings.TrimSpace(line[13:idx])
						if serverName != "" {
							currentMCP = serverName
							if _, ok := config.MCPServers[serverName]; !ok {
								config.MCPServers[serverName] = CodexMCPServer{Command: "npx", Args: []string{}}
							}
						}
					}
				}
				continue
			}

			// Top-level keys
			if strings.HasPrefix(line, "model = ") {
				config.Model = strings.Trim(strings.TrimSpace(line[len("model = "):]), `"`)
				continue
			}
			if strings.HasPrefix(line, "model_reasoning_effort = ") {
				config.ModelReasoningEffort = strings.Trim(strings.TrimSpace(line[len("model_reasoning_effort = "):]), `"`)
				continue
			}
			if strings.HasPrefix(line, "disable_response_storage = ") {
				config.DisableResponseStorage = strings.Contains(line, "true")
				continue
			}

			// Inside MCP server block
			if currentMCP != "" {
				srv := config.MCPServers[currentMCP]
				if strings.HasPrefix(line, "command = ") {
					srv.Command = strings.Trim(strings.TrimSpace(line[len("command = "):]), `"`)
					config.MCPServers[currentMCP] = srv
					continue
				}
				if strings.HasPrefix(line, "args = ") {
					argsPart := strings.TrimSpace(line[len("args = "):])
					srv.Args = parseTomlStringArray(argsPart)
					config.MCPServers[currentMCP] = srv
					continue
				}
			}
		}
	}

	return config
}

func (c *Config) generateCodexConfigTOML(config CodexConfig) string {
	toml := fmt.Sprintf(`model_provider = "%s"
model = "%s"
model_reasoning_effort = "%s"
disable_response_storage = %t

`, config.ModelProvider, config.Model, config.ModelReasoningEffort, config.DisableResponseStorage)

	for name, provider := range config.ModelProviders {
		toml += fmt.Sprintf(`
[model_providers.%s]
name = "%s"
base_url = "%s"
wire_api = "%s"
`, name, provider.Name, provider.BaseURL, provider.WireAPI)
		// Only add env_key if it's not empty
		if provider.EnvKey != "" {
			toml += fmt.Sprintf(`env_key = "%s"
`, provider.EnvKey)
		}
		toml += fmt.Sprintf(`requires_openai_auth = %t
`, provider.RequiresOpenAIAuth)
	}

	for path, project := range config.Projects {
		toml += fmt.Sprintf(`
[projects.%s]
trust_level = "%s"
`, path, project.TrustLevel)
	}

	for name, server := range config.MCPServers {
		toml += fmt.Sprintf(`
[mcp_servers.%s]
args = %s
command = "%s"
`, name, serverToTomlArray(server.Args), server.Command)
	}

	return toml
}

// parseTomlStringArray parses a TOML array of strings like ["a", "b"].
func parseTomlStringArray(s string) []string {
	s = strings.TrimSpace(s)
	start := strings.Index(s, "[")
	end := strings.LastIndex(s, "]")
	if start == -1 || end == -1 || end <= start {
		return []string{}
	}
	inner := strings.TrimSpace(s[start+1 : end])
	if inner == "" {
		return []string{}
	}
	parts := strings.Split(inner, ",")
	out := make([]string, 0, len(parts))
	for _, p := range parts {
		p = strings.TrimSpace(p)
		p = strings.Trim(p, `"`)
		if p != "" {
			out = append(out, p)
		}
	}
	return out
}

// serverToTomlArray formats []string as a TOML string array.
func serverToTomlArray(args []string) string {
	if len(args) == 0 {
		return "[]"
	}
	b := strings.Builder{}
	b.WriteString("[")
	for i, a := range args {
		if i > 0 {
			b.WriteString(", ")
		}
		a = strings.ReplaceAll(a, "\"", "\\\"")
		b.WriteString("\"")
		b.WriteString(a)
		b.WriteString("\"")
	}
	b.WriteString("]")
	return b.String()
}

// Droid configuration management methods
func (c *Config) AddDroidConfig(config DroidConfig) error {
	c.Droid = append(c.Droid, config)
	return c.Save()
}

func (c *Config) DeleteDroidConfig(index int) error {
	if index < 0 || index >= len(c.Droid) {
		return fmt.Errorf("invalid Droid index")
	}

	// Remove the config at index
	c.Droid = append(c.Droid[:index], c.Droid[index+1:]...)

	// Adjust active index if needed
	if c.Active.Droid == index {
		c.Active.Droid = -1
	} else if c.Active.Droid > index {
		c.Active.Droid--
	}

	return c.Save()
}

func (c *Config) SetActiveDroid(index int) error {
	if index >= 0 && index < len(c.Droid) {
		c.Active.Droid = index
		return c.Save()
	}
	return fmt.Errorf("invalid Droid index")
}

func (c *Config) GetActiveDroid() *DroidConfig {
	if c.Active.Droid >= 0 && c.Active.Droid < len(c.Droid) {
		return &c.Droid[c.Active.Droid]
	}
	return nil
}

func (c *Config) SwitchDroid(config *DroidConfig) error {
	if config == nil {
		return fmt.Errorf("config cannot be nil")
	}

	factoryDir := platformPaths.GetDroidConfigDir()
	if err := mkdirWithPerms(factoryDir, 0755); err != nil {
		return fmt.Errorf("failed to create .factory directory: %w", err)
	}

	// Create FactoryConfig with the custom model
	factoryConfig := FactoryConfig{
		CustomModels: []DroidConfig{*config},
	}

	data, err := json.MarshalIndent(factoryConfig, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal Factory config: %w", err)
	}

	configPath := filepath.Join(factoryDir, "config.json")
	if err := writeFileWithPerms(configPath, data, 0644); err != nil {
		return fmt.Errorf("failed to write config.json: %w", err)
	}

	// Update settings.json with the custom model
	settingsPath := filepath.Join(factoryDir, "settings.json")

	// Read existing settings.json to preserve other settings
	existingSettings := make(map[string]interface{})
	if data, err := os.ReadFile(settingsPath); err == nil {
		// Parse JSON with comments support by stripping comments first
		content := string(data)
		lines := strings.Split(content, "\n")
		var jsonLines []string
		for _, line := range lines {
			trimmed := strings.TrimSpace(line)
			if !strings.HasPrefix(trimmed, "//") && trimmed != "" {
				jsonLines = append(jsonLines, line)
			}
		}
		jsonContent := strings.Join(jsonLines, "\n")
		if jsonContent != "" {
			if err := json.Unmarshal([]byte(jsonContent), &existingSettings); err != nil {
				// If parsing fails, start with empty settings
				existingSettings = make(map[string]interface{})
			}
		}
	}

	// Update the model field with custom: prefix
	modelName := config.Model
	if modelName != "" {
		existingSettings["model"] = "custom:" + modelName
	}

	// Marshal back to JSON
	updatedData, err := json.MarshalIndent(existingSettings, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal updated settings: %w", err)
	}

	// Write the updated settings with comments header
	var finalContent strings.Builder
	finalContent.WriteString("// Factory CLI Settings\n")
	finalContent.WriteString("// This file contains your Factory CLI configuration.\n")
	finalContent.WriteString(string(updatedData))
	finalContent.WriteString("\n")

	return writeFileWithPerms(settingsPath, []byte(finalContent.String()), 0644)
}

// findConfigIndex 查找配置在原始列表中的索引
func findConfigIndex(configs []ServiceConfig, target ServiceConfig) int {
	for i, cfg := range configs {
		if cfg.Name == target.Name && cfg.Provider == target.Provider && cfg.BaseURL == target.BaseURL && cfg.APIKey == target.APIKey {
			return i
		}
	}
	return -1
}
