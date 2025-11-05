package tui

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"
)

// Styles
var (
	titleStyle          = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("12")).MarginTop(1)
	headerTitle         = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("#111827")).Background(lipgloss.Color("153")).Padding(0, 1)
	headerLine          = lipgloss.NewStyle().Foreground(lipgloss.Color("63"))
	cursorStyle         = lipgloss.NewStyle().Foreground(lipgloss.Color("109")).Bold(true)
	activeStyle         = lipgloss.NewStyle().Foreground(lipgloss.Color("#3b82f6")).Bold(true)
	editStyle           = lipgloss.NewStyle().Foreground(lipgloss.Color("5")).Bold(true)
	errorStyle          = lipgloss.NewStyle().Foreground(lipgloss.Color("1")).Bold(true)
	helpStyle           = lipgloss.NewStyle().Foreground(lipgloss.Color("8")).Italic(true)
	fieldHighlightStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("3")).Bold(true)
	boxStyle            = lipgloss.NewStyle().Border(lipgloss.RoundedBorder()).BorderForeground(lipgloss.Color("63")).Padding(1, 2)
	menuItemStyle       = lipgloss.NewStyle().Padding(0, 1)
	menuItemSelStyle    = lipgloss.NewStyle().Padding(0, 1).Foreground(lipgloss.Color("#111827")).Background(lipgloss.Color("#e5e7eb")).Bold(true)
	listRowSelStyle     = lipgloss.NewStyle().Background(lipgloss.Color("#e5e7eb")).Foreground(lipgloss.Color("#111827")).Bold(true)
	listRowStyle        = lipgloss.NewStyle().Foreground(lipgloss.Color("252"))
	formRowStyle        = lipgloss.NewStyle().Foreground(lipgloss.Color("252"))
	dividerStyle        = lipgloss.NewStyle().Foreground(lipgloss.Color("240"))
	itemBoxStyle        = lipgloss.NewStyle().Border(lipgloss.NormalBorder()).BorderForeground(lipgloss.Color("60")).Padding(0, 1).MarginBottom(0)
	itemBoxSelStyle     = lipgloss.NewStyle().Border(lipgloss.DoubleBorder()).BorderForeground(lipgloss.Color("66")).Padding(0, 1).MarginBottom(0).Bold(true)
)

// Header and UI helpers
func headerView(title string) string {
	line := strings.Repeat("─", max(8, len(title)+6))
	return lipgloss.JoinHorizontal(lipgloss.Top,
		headerTitle.Render(" "+title+" "),
		headerLine.Render(" "+line),
	)
}

func statusBarView(seg1, seg2, seg3, seg4 string) string {
	parts := []string{}
	for _, s := range []string{seg1, seg2, seg3, seg4} {
		if strings.TrimSpace(s) != "" {
			parts = append(parts, helpStyle.Render(" "+s+" "))
		}
	}
	return lipgloss.NewStyle().MarginTop(1).Render(strings.Join(parts, "｜"))
}

func menuItemView(text string, selected bool) string {
	if selected {
		return menuItemSelStyle.Render("▶ " + text)
	}
	return menuItemStyle.Render("  " + text)
}

func listRowView(cfg ServiceConfig, selected, active bool, compact bool) string {
	name := cfg.Name
	if active {
		name = name + " " + activeStyle.Render("[当前使用]")
	}
	badge := providerBadge(cfg.Provider)
	var text string
	if compact {
		text = fmt.Sprintf("%s  · %s", name, badge)
	} else {
		lines := []string{
			fmt.Sprintf("配置名称: %s", name),
			fmt.Sprintf("Provider: %s", badge),
			fmt.Sprintf("Base URL: %s", cfg.BaseURL),
		}
		text = strings.Join(lines, "\n")
	}
	if selected {
		text = "┃ " + text
		return listRowSelStyle.Render(text)
	}
	return listRowStyle.Render(text)
}

// Provider 彩色徽章（低饱和莫兰迪色系，稍偏彩）
func providerBadge(p string) string {
	c := providerColor(p)
	st := lipgloss.NewStyle().Foreground(lipgloss.Color("#111827")).Background(lipgloss.Color(c)).Padding(0, 1).Bold(true)
	return st.Render(p)
}

func providerColor(p string) string {
	_ = strings.ToLower(strings.TrimSpace(p))
	return "109"
}
