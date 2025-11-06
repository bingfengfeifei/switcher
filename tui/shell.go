package tui

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
)

type ShellManager interface {
	SetEnvVar(key, value string) error
}

func NewShellManager() ShellManager {
	if runtime.GOOS == "windows" {
		return &windowsShellManager{}
	}
	return &unixShellManager{}
}

// Unix shell manager (bash, zsh, fish)
type unixShellManager struct{}

func (m *unixShellManager) SetEnvVar(key, value string) error {
	home, err := os.UserHomeDir()
	if err != nil {
		return fmt.Errorf("failed to get home directory: %w", err)
	}

	// Detect current shell
	shell := os.Getenv("SHELL")

	// Update bash/zsh
	if strings.Contains(shell, "bash") || strings.Contains(shell, "zsh") || shell == "" {
		bashrcPath := filepath.Join(home, ".bashrc")
		if err := updateBashConfig(bashrcPath, key, value); err != nil {
			return fmt.Errorf("failed to update .bashrc: %w", err)
		}
	}

	// Update fish if exists
	fishConfigDir := filepath.Join(home, ".config", "fish")
	if _, err := os.Stat(fishConfigDir); err == nil {
		fishConfigPath := filepath.Join(fishConfigDir, "config.fish")
		if err := updateFishConfig(fishConfigPath, key, value); err != nil {
			return fmt.Errorf("failed to update fish config: %w", err)
		}
	}

	return nil
}

func updateBashConfig(configPath, key, value string) error {
	content, err := os.ReadFile(configPath)
	if err != nil && !os.IsNotExist(err) {
		return err
	}

	lines := strings.Split(string(content), "\n")
	keyPattern := fmt.Sprintf("export %s=", key)

	var newLines []string
	for _, line := range lines {
		if !strings.Contains(line, keyPattern) {
			newLines = append(newLines, line)
		}
	}

	newLines = append(newLines, fmt.Sprintf("export %s=\"%s\"", key, value))
	return writeFileWithPerms(configPath, []byte(strings.Join(newLines, "\n")), 0644)
}

func updateFishConfig(configPath, key, value string) error {
	content, err := os.ReadFile(configPath)
	if err != nil && !os.IsNotExist(err) {
		return err
	}

	lines := strings.Split(string(content), "\n")
	keyPattern := fmt.Sprintf("set -x %s ", key)

	var newLines []string
	for _, line := range lines {
		if !strings.Contains(line, keyPattern) {
			newLines = append(newLines, line)
		}
	}

	newLines = append(newLines, fmt.Sprintf("set -x %s \"%s\"", key, value))
	return writeFileWithPerms(configPath, []byte(strings.Join(newLines, "\n")), 0644)
}

// Windows shell manager (PowerShell)
type windowsShellManager struct{}

func (m *windowsShellManager) SetEnvVar(key, value string) error {
	// Set user environment variable permanently
	cmd := exec.Command("powershell", "-Command",
		fmt.Sprintf("[Environment]::SetEnvironmentVariable('%s', '%s', 'User')", key, value))

	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to set environment variable: %w", err)
	}

	return nil
}
