package tui

import (
	"os"
	"path/filepath"
	"runtime"
)

type PlatformPaths interface {
	GetAppConfigPath() string
	GetClaudeConfigDir() string
	GetCodexConfigDir() string
	GetDroidConfigDir() string
}

type platformPathsImpl struct {
	home string
}

func NewPlatformPaths() (PlatformPaths, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return nil, err
	}

	switch runtime.GOOS {
	case "windows":
		return &windowsPaths{home: home}, nil
	case "darwin":
		return &darwinPaths{home: home}, nil
	default: // linux and others
		return &linuxPaths{home: home}, nil
	}
}

// Linux paths
type linuxPaths struct {
	home string
}

func (p *linuxPaths) GetAppConfigPath() string {
	return filepath.Join(p.home, ".config", "switcher", "config.json")
}

func (p *linuxPaths) GetClaudeConfigDir() string {
	return filepath.Join(p.home, ".claude")
}

func (p *linuxPaths) GetCodexConfigDir() string {
	return filepath.Join(p.home, ".codex")
}

func (p *linuxPaths) GetDroidConfigDir() string {
	return filepath.Join(p.home, ".factory")
}

// macOS paths
type darwinPaths struct {
	home string
}

func (p *darwinPaths) GetAppConfigPath() string {
	return filepath.Join(p.home, "Library", "Application Support", "switcher", "config.json")
}

func (p *darwinPaths) GetClaudeConfigDir() string {
	return filepath.Join(p.home, ".claude")
}

func (p *darwinPaths) GetCodexConfigDir() string {
	return filepath.Join(p.home, ".codex")
}

func (p *darwinPaths) GetDroidConfigDir() string {
	return filepath.Join(p.home, ".factory")
}

// Windows paths
type windowsPaths struct {
	home string
}

func (p *windowsPaths) GetAppConfigPath() string {
	appData := os.Getenv("APPDATA")
	if appData == "" {
		appData = filepath.Join(p.home, "AppData", "Roaming")
	}
	return filepath.Join(appData, "switcher", "config.json")
}

func (p *windowsPaths) GetClaudeConfigDir() string {
	return filepath.Join(p.home, ".claude")
}

func (p *windowsPaths) GetCodexConfigDir() string {
	return filepath.Join(p.home, ".codex")
}

func (p *windowsPaths) GetDroidConfigDir() string {
	return filepath.Join(p.home, ".factory")
}

// mkdirWithPerms creates directory with permissions on Unix, ignores perms on Windows
func mkdirWithPerms(path string, perm os.FileMode) error {
	if runtime.GOOS == "windows" {
		return os.MkdirAll(path, 0)
	}
	return os.MkdirAll(path, perm)
}

// writeFileWithPerms writes file with permissions on Unix, ignores perms on Windows
func writeFileWithPerms(path string, data []byte, perm os.FileMode) error {
	if runtime.GOOS == "windows" {
		return os.WriteFile(path, data, 0)
	}
	return os.WriteFile(path, data, perm)
}
