package tui

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"
)

func (m model) droidListView() string {
	header := headerView("Droid 配置")

	var rows []string

	// 显示排序后的配置
	for i, cfg := range m.sortedDroid {
		active := m.config.Active.Droid == findDroidConfigIndex(m.config.Droid, cfg)
		r := droidListRowView(cfg, i == m.cursor, active, m.compact)
		if i == m.cursor {
			r = itemBoxSelStyle.Render(r)
		} else {
			r = itemBoxStyle.Render(r)
		}
		rows = append(rows, r)
	}
	{
		backSel := m.cursor == len(m.config.Droid)
		back := menuItemView("← 返回主菜单", backSel)
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
	content.WriteString(statusBarView("↑/↓ 选择", "Enter 切换", "Tab 编辑 Del 删除", "A 添加  V 紧凑/展开  Esc 返回  ←/→ 切服务"))
	return content.String()
}

func (m model) addDroidConfigView() string {
	title := headerView("添加 Droid 配置")

	var content strings.Builder
	content.WriteString(title)
	content.WriteString("\n\n")

	fields := []struct {
		label string
		value string
	}{
		{"模型显示名称", m.droidFormData.ModelDisplayName},
		{"模型名称", m.droidFormData.Model},
		{"Base URL", m.droidFormData.BaseURL},
		{"API Key", m.droidFormData.APIKey},
	}

	var inner strings.Builder
	for i, field := range fields {
		prefix := "  "
		if m.formField == i {
			prefix = cursorStyle.Render(">")
		}

		// 对于API密钥字段，如果正在编辑，显示完整内容，否则显示遮蔽内容
		displayValue := field.value
		if i == FieldAPIKey && m.formField != FieldAPIKey { // API密钥字段且不在编辑状态
			displayValue = maskAPIKey(field.value)
		} else if i == FieldAPIKey && m.formField == FieldAPIKey {
			// 如果正在编辑API密钥字段，显示完整内容但添加提示
			displayValue = field.value + " (编辑中)"
		}

		highlight := ""
		if m.formField == i {
			highlight = fieldHighlightStyle.Render(" ← 正在编辑，请直接输入修改内容")
		}

		inner.WriteString(formRowStyle.Render(fmt.Sprintf("%s %s:%s %s", prefix, field.label, highlight, displayValue)) + "\n")
	}

	content.WriteString(boxStyle.Render(inner.String()))
	content.WriteString("\n")
	content.WriteString(statusBarView("Tab/↑/↓ 切字段", "Enter 保存", "Esc 取消", ""))

	// 添加当前编辑状态提示
	if m.formField >= 0 && m.formField < DroidFieldCount {
		content.WriteString("\n" + fieldHighlightStyle.Render("✨ 当前正在编辑: ") + fields[m.formField].label)
		if m.formField == FieldAPIKey {
			content.WriteString("\n" + fieldHighlightStyle.Render("   (API Key 正在显示完整内容以便编辑)"))
		}
	}

	return content.String()
}

func (m model) editDroidConfigView() string {
	title := headerView("编辑 Droid 配置")

	var content strings.Builder
	content.WriteString(title)
	content.WriteString("\n\n")

	fields := []struct {
		label string
		value string
	}{
		{"模型显示名称", m.droidFormData.ModelDisplayName},
		{"模型名称", m.droidFormData.Model},
		{"Base URL", m.droidFormData.BaseURL},
		{"API Key", m.droidFormData.APIKey},
	}

	var inner strings.Builder
	for i, field := range fields {
		prefix := "  "
		if m.formField == i {
			prefix = cursorStyle.Render(">")
		}

		// 对于API密钥字段，如果正在编辑，显示完整内容，否则显示遮蔽内容
		displayValue := field.value
		if i == FieldAPIKey && m.formField != FieldAPIKey { // API密钥字段且不在编辑状态
			displayValue = maskAPIKey(field.value)
		} else if i == FieldAPIKey && m.formField == FieldAPIKey {
			// 如果正在编辑API密钥字段，显示完整内容但添加提示
			displayValue = field.value + " (编辑中)"
		}

		highlight := ""
		if m.formField == i {
			highlight = fieldHighlightStyle.Render(" ← 正在编辑，请直接输入修改内容")
		}

		inner.WriteString(formRowStyle.Render(fmt.Sprintf("%s %s:%s %s", prefix, field.label, highlight, displayValue)) + "\n")
	}

	content.WriteString(boxStyle.Render(inner.String()))
	content.WriteString("\n")
	content.WriteString(statusBarView("Tab/↑/↓ 切字段", "Enter 保存", "Esc 取消", ""))

	// 添加当前编辑状态提示
	if m.formField >= 0 && m.formField < DroidFieldCount {
		content.WriteString("\n" + fieldHighlightStyle.Render("✨ 当前正在编辑: ") + fields[m.formField].label)
		if m.formField == FieldAPIKey {
			content.WriteString("\n" + fieldHighlightStyle.Render("   (API Key 正在显示完整内容以便编辑)"))
		}
	}

	return content.String()
}

func droidListRowView(cfg DroidConfig, selected bool, active bool, compact bool) string {
	var name string
	if cfg.ModelDisplayName != "" {
		name = cfg.ModelDisplayName
	} else {
		name = cfg.Model
	}

	status := "  "
	if active {
		status = activeStyle.Render("✓ ")
	}

	key := cfg.APIKey
	if !compact && !selected {
		key = maskAPIKey(key)
	}

	if compact {
		return fmt.Sprintf("%s%s (%s)", status, name, cfg.Provider)
	}
	return fmt.Sprintf("%s%s\n    Provider: %s\n    BaseURL: %s\n    APIKey: %s",
		status, name, cfg.Provider, cfg.BaseURL, key)
}

// findDroidConfigIndex 查找 Droid 配置在原始列表中的索引
func findDroidConfigIndex(configs []DroidConfig, target DroidConfig) int {
	for i, cfg := range configs {
		if cfg.ModelDisplayName == target.ModelDisplayName && cfg.Model == target.Model && cfg.BaseURL == target.BaseURL && cfg.APIKey == target.APIKey {
			return i
		}
	}
	return -1
}

// getSortedDroidConfigs 获取排序后的 Droid 配置列表
func (m model) getSortedDroidConfigs() []DroidConfig {
	sortedConfigs := make([]DroidConfig, 0, len(m.config.Droid))
	var activeConfig DroidConfig
	var otherConfigs []DroidConfig

	for i, cfg := range m.config.Droid {
		if m.config.Active.Droid == i {
			activeConfig = cfg
		} else {
			otherConfigs = append(otherConfigs, cfg)
		}
	}

	// 按名称对剩余配置进行排序
	for i := 0; i < len(otherConfigs)-1; i++ {
		for j := i + 1; j < len(otherConfigs); j++ {
			if otherConfigs[i].ModelDisplayName > otherConfigs[j].ModelDisplayName {
				otherConfigs[i], otherConfigs[j] = otherConfigs[j], otherConfigs[i]
			}
		}
	}

	// 构建排序后的列表：当前配置在前，其余在后
	if activeConfig.ModelDisplayName != "" {
		sortedConfigs = append(sortedConfigs, activeConfig)
	}
	sortedConfigs = append(sortedConfigs, otherConfigs...)
	return sortedConfigs
}

// getOriginalDroidIndex 根据排序列表中的位置获取原始索引
func (m model) getOriginalDroidIndex(sortedIndex int) int {
	if sortedIndex < 0 || sortedIndex >= len(m.getSortedDroidConfigs()) {
		return -1
	}
	sortedConfig := m.getSortedDroidConfigs()[sortedIndex]
	return findDroidConfigIndex(m.config.Droid, sortedConfig)
}

// sortDroidConfigs 对 Droid 配置进行排序
func (m *model) sortDroidConfigs() {
	if m.sortedDroid == nil {
		m.sortedDroid = m.getSortedDroidConfigs()
	}
}
