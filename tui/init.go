package tui

import (
	"os"

	tea "github.com/charmbracelet/bubbletea"
	"golang.org/x/term"
)

func InitialModel(config *Config) model {
	return model{
		config:           config,
		state:            mainMenu,
		cursor:           0,
		compact:          false,
		sortedClaudeCode: nil,
		sortedCodex:      nil,
		sortedDroid:      nil,
		windowHeight:     detectWindowHeight(),
	}
}

func (m model) Init() tea.Cmd {
	return nil
}

func detectWindowHeight() int {
	fds := []uintptr{
		os.Stdout.Fd(),
		os.Stdin.Fd(),
		os.Stderr.Fd(),
	}

	for _, fd := range fds {
		if _, h, err := term.GetSize(int(fd)); err == nil && h > 0 {
			return h
		}
	}

	// Reasonable default that keeps the initial viewport compact on small terminals.
	return 24
}
