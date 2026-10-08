package tui

import (
	"os"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"golang.org/x/term"
)

func InitialModel(config *Config) model {
	return model{
		config:           config,
		state:            mainMenu,
		cursor:           0,
		compact:          false,
		cursorVisible:    true,
		sortedClaudeCode: nil,
		sortedCodex:      nil,
		sortedDroid:      nil,
		windowHeight:     detectWindowHeight(),
	}
}

// blinkMsg 输入光标闪烁消息
type blinkMsg struct{}

// blinkCursor 每 530ms 切换一次输入光标可见性（收到 blinkMsg 后重新注册）
func blinkCursor() tea.Cmd {
	return tea.Tick(530*time.Millisecond, func(time.Time) tea.Msg {
		return blinkMsg{}
	})
}

func (m model) Init() tea.Cmd {
	return blinkCursor()
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
