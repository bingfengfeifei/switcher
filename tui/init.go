package tui

import tea "github.com/charmbracelet/bubbletea"

func InitialModel(config *Config) model {
	return model{
		config:           config,
		state:            mainMenu,
		cursor:           0,
		compact:          false,
		sortedClaudeCode: nil,
		sortedCodex:      nil,
		sortedDroid:      nil,
	}
}

func (m model) Init() tea.Cmd {
	return nil
}
