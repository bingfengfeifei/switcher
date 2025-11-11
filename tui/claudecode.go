package tui

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/charmbracelet/lipgloss"
)

func (m model) claudeCodeListView() string {
	header := headerView(t("header_claude"))

	var rows []string

	// 检查是否有警告信息
	hasWarning := false
	if ok, actualBase, _ := checkAppliedClaudeLocal(m.config); !ok && actualBase != "" {
		warn := errorStyle.Render(fmt.Sprintf(t("warn_mismatch"), "Claude") + actualBase)
		rows = append(rows, itemBoxStyle.Render(warn))
		hasWarning = true
	}

	// Calculate viewport size based on window height
	configCount := len(m.sortedClaudeCode)
	if configCount > 0 {
		viewportSize := calculateListViewportHeight(m.windowHeight, hasWarning, m.compact)

		// Calculate visible range based on cursor position
		start, end := updateCursorViewport(m.cursor, configCount, viewportSize)

		// Render visible configurations
		for i := start; i < end; i++ {
			cfg := m.sortedClaudeCode[i]
			active := m.config.Active.ClaudeCode == findConfigIndex(m.config.ClaudeCode, cfg)
			r := listRowView(cfg, i == m.cursor, active, m.compact)
			if i == m.cursor {
				r = itemBoxSelStyle.Render(r)
			} else {
				r = itemBoxStyle.Render(r)
			}
			rows = append(rows, r)
		}

		// Add scroll indicators if needed
		if start > 0 {
			rows = append([]string{scrollIndicatorStyle.Render("↑ More items above")}, rows...)
		}
		if end < configCount {
			rows = append(rows, scrollIndicatorStyle.Render("↓ More items below"))
		}
	}

	// Add "Back to menu" option
	backSel := m.cursor == len(m.config.ClaudeCode)
	back := menuItemView(t("back_to_menu"), backSel)
	if backSel {
		rows = append(rows, itemBoxSelStyle.Render(back))
	} else {
		rows = append(rows, itemBoxStyle.Render(back))
	}

	body := lipgloss.JoinVertical(lipgloss.Left, rows...)

	var content strings.Builder
	content.WriteString(header)
	content.WriteString("\n\n")
	content.WriteString(body)
	content.WriteString("\n")
	content.WriteString(statusBarView(t("nav_select"), t("nav_confirm"), t("nav_edit"), t("nav_add")+"  "+t("nav_view")+"  "+t("nav_back")+"  "+t("nav_switch_service")))
	return content.String()
}

// getSortedClaudeCodeConfigs 获取排序后的 Claude Code 配置列表
func (m model) getSortedClaudeCodeConfigs() []ServiceConfig {
	sortedConfigs := make([]ServiceConfig, 0, len(m.config.ClaudeCode))
	var activeConfig ServiceConfig
	var otherConfigs []ServiceConfig

	for i, cfg := range m.config.ClaudeCode {
		if m.config.Active.ClaudeCode == i {
			activeConfig = cfg
		} else {
			otherConfigs = append(otherConfigs, cfg)
		}
	}

	// 按名称对剩余配置进行排序
	for i := 0; i < len(otherConfigs)-1; i++ {
		for j := i + 1; j < len(otherConfigs); j++ {
			if otherConfigs[i].Name > otherConfigs[j].Name {
				otherConfigs[i], otherConfigs[j] = otherConfigs[j], otherConfigs[i]
			}
		}
	}

	// 构建排序后的列表：当前配置在前，其余在后
	if activeConfig.Name != "" {
		sortedConfigs = append(sortedConfigs, activeConfig)
	}
	sortedConfigs = append(sortedConfigs, otherConfigs...)
	return sortedConfigs
}

// getOriginalClaudeCodeIndex 根据排序列表中的位置获取原始索引
func (m model) getOriginalClaudeCodeIndex(sortedIndex int) int {
	if sortedIndex < 0 || sortedIndex >= len(m.getSortedClaudeCodeConfigs()) {
		return -1
	}
	sortedConfig := m.getSortedClaudeCodeConfigs()[sortedIndex]
	return findConfigIndex(m.config.ClaudeCode, sortedConfig)
}

// sortClaudeCodeConfigs 对 Claude Code 配置进行排序
func (m *model) sortClaudeCodeConfigs() {
	if m.sortedClaudeCode == nil {
		m.sortedClaudeCode = m.getSortedClaudeCodeConfigs()
	}
}

// Local checks for applied vs selected configs
func checkAppliedClaudeLocal(c *Config) (bool, string, error) {
	active := c.GetActiveClaudeCode()
	if active == nil {
		return true, "", nil
	}
	home, err := os.UserHomeDir()
	if err != nil {
		return true, "", err
	}
	data, err := os.ReadFile(filepath.Join(home, ".claude", "settings.json"))
	if err != nil {
		return true, "", nil
	}
	var se ClaudeSettings
	if err := json.Unmarshal(data, &se); err != nil {
		return true, "", err
	}
	ab := strings.TrimSpace(se.Env["ANTHROPIC_BASE_URL"])
	ak := strings.TrimSpace(se.Env["ANTHROPIC_AUTH_TOKEN"])
	ok := ab == strings.TrimSpace(active.BaseURL) && ak == strings.TrimSpace(active.APIKey)
	return ok, ab, nil
}
