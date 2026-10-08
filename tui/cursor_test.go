package tui

import (
	"strings"
	"testing"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/muesli/termenv"
)

// press 驱动一次按键并返回更新后的 model
func press(m model, msg tea.Msg) model {
	next, _ := m.Update(msg)
	return next.(model)
}

func TestRenderInputCursor(t *testing.T) {
	lipgloss.SetColorProfile(termenv.TrueColor)
	t.Cleanup(func() { lipgloss.SetColorProfile(termenv.Ascii) })

	// 末尾：块光标
	if got := renderInputCursor("abc", 3, true); !strings.Contains(got, "▌") {
		t.Fatalf("cursor at end should render block cursor, got %q", got)
	}
	// 中间：反显光标处字符
	if got := renderInputCursor("abc", 1, true); !strings.Contains(got, "7m") {
		t.Fatalf("cursor mid-text should reverse the char, got %q", got)
	}
	// 闪烁隐藏：原样返回
	if got := renderInputCursor("abc", 1, false); got != "abc" {
		t.Fatalf("hidden cursor should return plain value, got %q", got)
	}
	// 越界位置自动收敛
	if got := renderInputCursor("abc", 99, true); !strings.Contains(got, "▌") {
		t.Fatalf("out-of-range cursor should clamp to end, got %q", got)
	}
}

func TestFormCursorMoveAndInsert(t *testing.T) {
	m := model{state: editCodex, formField: FieldBaseURL, formData: ServiceConfig{BaseURL: "abc"}, cursorVisible: true}
	m.resetFormCursor() // cursor -> 3

	m = press(m, tea.KeyMsg{Type: tea.KeyLeft})
	m = press(m, tea.KeyMsg{Type: tea.KeyLeft})
	if m.formCursor != 1 {
		t.Fatalf("cursor after two lefts = %d, want 1", m.formCursor)
	}

	// 光标处插入 X：abc -> aXbc
	m = press(m, tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune("X")})
	if m.formData.BaseURL != "aXbc" || m.formCursor != 2 {
		t.Fatalf("insert at cursor: value=%q cursor=%d, want aXbc/2", m.formData.BaseURL, m.formCursor)
	}

	// 退格删前：aXbc -> abc，cursor 2 -> 1
	m = press(m, tea.KeyMsg{Type: tea.KeyBackspace})
	if m.formData.BaseURL != "abc" || m.formCursor != 1 {
		t.Fatalf("backspace before cursor: value=%q cursor=%d, want abc/1", m.formData.BaseURL, m.formCursor)
	}

	// Delete 删后：abc -> ac，cursor 不动
	m = press(m, tea.KeyMsg{Type: tea.KeyDelete})
	if m.formData.BaseURL != "ac" || m.formCursor != 1 {
		t.Fatalf("delete after cursor: value=%q cursor=%d, want ac/1", m.formData.BaseURL, m.formCursor)
	}

	// Home/End
	m = press(m, tea.KeyMsg{Type: tea.KeyHome})
	if m.formCursor != 0 {
		t.Fatalf("cursor after home = %d, want 0", m.formCursor)
	}
	m = press(m, tea.KeyMsg{Type: tea.KeyEnd})
	if m.formCursor != 2 {
		t.Fatalf("cursor after end = %d, want 2", m.formCursor)
	}
}

func TestFormCursorResetsOnFieldSwitch(t *testing.T) {
	m := model{
		state:     editCodex,
		formField: FieldBaseURL,
		formData:  ServiceConfig{Name: "n", BaseURL: "abcd"},
	}
	m.resetFormCursor()
	m = press(m, tea.KeyMsg{Type: tea.KeyLeft}) // cursor 4 -> 3

	// Tab 切到下一字段后光标 snap 到新字段值末尾
	m = press(m, tea.KeyMsg{Type: tea.KeyTab})
	if m.formField != FieldAPIKey || m.formCursor != 0 {
		t.Fatalf("after tab: field=%d cursor=%d, want APIKey/0", m.formField, m.formCursor)
	}

	// 选择型字段不显示光标、不接受输入
	m2 := model{state: editCodex, formField: FieldWireAPI, formData: ServiceConfig{WireAPI: DefaultWireAPI}}
	m2 = press(m2, tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune("x")})
	if m2.formData.WireAPI != DefaultWireAPI {
		t.Fatalf("select field accepted text input: %q", m2.formData.WireAPI)
	}
}

func TestBlinkCursor(t *testing.T) {
	m := model{state: editCodex, formField: FieldBaseURL, formData: ServiceConfig{BaseURL: "abc"}, cursorVisible: true}
	m.resetFormCursor()

	next, cmd := m.Update(blinkMsg{})
	got := next.(model)
	if got.cursorVisible {
		t.Fatal("blink should hide cursor")
	}
	if cmd == nil {
		t.Fatal("blink should schedule next tick")
	}

	// 任意输入让光标立即恢复可见
	got = press(got, tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune("z")})
	if !got.cursorVisible {
		t.Fatal("typing should make cursor visible")
	}
}
