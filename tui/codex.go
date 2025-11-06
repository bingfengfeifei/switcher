package tui

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/charmbracelet/lipgloss"
)

func (m model) codexListView() string {
	header := headerView(t("header_codex"))

	var rows []string
	if ok, actualBase, _ := checkAppliedCodexLocal(m.config); !ok && actualBase != "" {
		warn := errorStyle.Render(fmt.Sprintf(t("warn_mismatch"), "Codex") + actualBase)
		rows = append(rows, itemBoxStyle.Render(warn))
	}

	// 显示排序后的配置
	for i, cfg := range m.sortedCodex {
		active := m.config.Active.Codex == findConfigIndex(m.config.Codex, cfg)
		r := listRowView(cfg, i == m.cursor, active, m.compact)
		if i == m.cursor {
			r = itemBoxSelStyle.Render(r)
		} else {
			r = itemBoxStyle.Render(r)
		}
		rows = append(rows, r)
	}
	{
		backSel := m.cursor == len(m.config.Codex)
		back := menuItemView(t("back_to_menu"), backSel)
		if backSel {
			rows = append(rows, itemBoxSelStyle.Render(back))
		} else {
			rows = append(rows, itemBoxStyle.Render(back))
		}
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

// getSortedCodexConfigs 获取排序后的 Codex 配置列表
func (m model) getSortedCodexConfigs() []ServiceConfig {
	sortedConfigs := make([]ServiceConfig, 0, len(m.config.Codex))
	var activeConfig ServiceConfig
	var otherConfigs []ServiceConfig

	for i, cfg := range m.config.Codex {
		if m.config.Active.Codex == i {
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

// getOriginalCodexIndex 根据排序列表中的位置获取原始索引
func (m model) getOriginalCodexIndex(sortedIndex int) int {
	if sortedIndex < 0 || sortedIndex >= len(m.getSortedCodexConfigs()) {
		return -1
	}
	sortedConfig := m.getSortedCodexConfigs()[sortedIndex]
	return findConfigIndex(m.config.Codex, sortedConfig)
}

// sortCodexConfigs 对 Codex 配置进行排序
func (m *model) sortCodexConfigs() {
	if m.sortedCodex == nil {
		m.sortedCodex = m.getSortedCodexConfigs()
	}
}

func checkAppliedCodexLocal(c *Config) (bool, string, error) {
	active := c.GetActiveCodex()
	if active == nil {
		return true, "", nil
	}
	home, err := os.UserHomeDir()
	if err != nil {
		return true, "", err
	}
	a, err := os.ReadFile(filepath.Join(home, ".codex", "auth.json"))
	if err != nil {
		return true, "", nil
	}
	var au CodexAuth
	if err := json.Unmarshal(a, &au); err != nil {
		return true, "", err
	}
	if strings.TrimSpace(au.OPENAI_API_KEY) != strings.TrimSpace(active.APIKey) {
		return false, "", nil
	}
	b, err := os.ReadFile(filepath.Join(home, ".codex", "config.toml"))
	if err != nil {
		return true, "", nil
	}
	cfg := string(b)
	wantProv := fmt.Sprintf("model_provider = \"%s\"", strings.TrimSpace(active.Provider))
	wantBase := fmt.Sprintf("base_url = \"%s\"", strings.TrimSpace(active.BaseURL))
	if !strings.Contains(cfg, wantProv) {
		return false, "", nil
	}
	if !strings.Contains(cfg, wantBase) {
		return false, "", nil
	}
	return true, strings.TrimSpace(active.BaseURL), nil
}
