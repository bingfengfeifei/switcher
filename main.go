package main

import (
	"flag"
	"fmt"
	"os"

	tui "switcher/tui"

	tea "github.com/charmbracelet/bubbletea"
)

// Version information, set by GoReleaser at build time
var (
	version = "dev"
	commit  = "none"
	date    = "unknown"
	builtBy = "unknown"
)

func init() {
	// 设置版本信息到tui包
	tui.AppVersion = version
}

func main() {
	// Non-interactive switching support
	var switchCodexName string
	var switchClaudeName string
	var switchDroidName string
	var showVersion bool
	flag.StringVar(&switchCodexName, "switch-codex", "", "Switch Codex to config by name")
	flag.StringVar(&switchClaudeName, "switch-claude", "", "Switch Claude Code to config by name")
	flag.StringVar(&switchDroidName, "switch-droid", "", "Switch Droid to config by name")
	flag.BoolVar(&showVersion, "version", false, "Show version information")
	flag.BoolVar(&showVersion, "v", false, "Show version information (shorthand)")
	flag.Parse()

	// Handle version flag
	if showVersion {
		fmt.Printf("switcher version %s\n", version)
		fmt.Printf("  commit: %s\n", commit)
		fmt.Printf("  built at: %s\n", date)
		fmt.Printf("  built by: %s\n", builtBy)
		return
	}

	config := &tui.Config{}
	if err := config.Load(); err != nil {
		fmt.Printf("Error loading config: %v\n", err)
		os.Exit(1)
	}

	if switchCodexName != "" {
		idx := -1
		for i, sc := range config.Codex {
			if sc.Name == switchCodexName {
				idx = i
				break
			}
		}
		if idx == -1 {
			fmt.Printf("Codex config not found: %s\n", switchCodexName)
			os.Exit(2)
		}
		sc := config.Codex[idx]
		if err := config.SwitchCodex(&sc); err != nil {
			fmt.Printf("Switch Codex failed: %v\n", err)
			os.Exit(3)
		}
		if err := config.SetActiveCodex(idx); err != nil {
			fmt.Printf("Set active Codex failed: %v\n", err)
			os.Exit(4)
		}
		fmt.Printf("Switched Codex to '%s' (provider=%s)\n", sc.Name, sc.Provider)
		return
	}

	if switchClaudeName != "" {
		idx := -1
		for i, sc := range config.ClaudeCode {
			if sc.Name == switchClaudeName {
				idx = i
				break
			}
		}
		if idx == -1 {
			fmt.Printf("Claude Code config not found: %s\n", switchClaudeName)
			os.Exit(2)
		}
		sc := config.ClaudeCode[idx]
		if err := config.SwitchClaudeCode(&sc); err != nil {
			fmt.Printf("Switch Claude Code failed: %v\n", err)
			os.Exit(3)
		}
		if err := config.SetActiveClaudeCode(idx); err != nil {
			fmt.Printf("Set active Claude Code failed: %v\n", err)
			os.Exit(4)
		}
		fmt.Printf("Switched Claude Code to '%s' (provider=%s)\n", sc.Name, sc.Provider)
		return
	}

	if switchDroidName != "" {
		idx := -1
		for i, dc := range config.Droid {
			if dc.ModelDisplayName == switchDroidName {
				idx = i
				break
			}
		}
		if idx == -1 {
			fmt.Printf("Droid config not found: %s\n", switchDroidName)
			os.Exit(2)
		}
		dc := config.Droid[idx]
		if err := config.SwitchDroid(&dc); err != nil {
			fmt.Printf("Switch Droid failed: %v\n", err)
			os.Exit(3)
		}
		if err := config.SetActiveDroid(idx); err != nil {
			fmt.Printf("Set active Droid failed: %v\n", err)
			os.Exit(4)
		}
		fmt.Printf("Switched Droid to '%s' (provider=%s)\n", dc.ModelDisplayName, dc.Provider)
		return
	}

	// Default to TUI
	p := tea.NewProgram(tui.InitialModel(config), tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		fmt.Printf("Error running TUI: %v\n", err)
		os.Exit(1)
	}
}
