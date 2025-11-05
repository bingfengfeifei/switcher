package tui

import (
	"fmt"
	"strings"
)

// 字段索引常量
const (
	FieldName = iota
	FieldProvider
	FieldBaseURL
	FieldAPIKey
	FieldModel
	FieldWireAPI
	FieldModelReasoningEffort
)

// 配置类型字段数量
const (
	ClaudeCodeFieldCount = 4
	CodexFieldCount      = 7
	DroidFieldCount      = 4
)

type state int

const (
	mainMenu state = iota
	claudeCodeList
	codexList
	droidList
	// 已移除操作菜单，使用简化流程：Enter 切换、Tab 编辑
	addClaudeCode
	addCodex
	addDroid
	editClaudeCode
	editCodex
	editDroid
	confirmDeleteClaudeCode
	confirmDeleteCodex
	confirmDeleteDroid
	confirmExitAddClaudeCode
	confirmExitAddCodex
	confirmExitAddDroid
)

type model struct {
	compact          bool
	config           *Config
	state            state
	cursor           int
	selected         int
	formData         ServiceConfig
	droidFormData    DroidConfig
	formField        int
	error            string
	editIndex        int
	actionType       string          // 新增：操作类型 ("switch" 或 "edit")
	deleteIndex      int             // 要删除的配置索引
	sortedClaudeCode []ServiceConfig // 排序后的 Claude Code 配置列表
	sortedCodex      []ServiceConfig // 排序后的 Codex 配置列表
	sortedDroid      []DroidConfig   // 排序后的 Droid 配置列表
}

func (m model) hasFormContent() bool {
	return m.formData.Name != "" || m.formData.Provider != "" || m.formData.BaseURL != "" || m.formData.APIKey != "" || m.formData.Model != "" || m.formData.WireAPI != "" || m.formData.EnvKey != "" || m.formData.ModelReasoningEffort != ""
}

func (m model) hasDroidFormContent() bool {
	return m.droidFormData.ModelDisplayName != "" || m.droidFormData.Model != "" || m.droidFormData.BaseURL != "" || m.droidFormData.APIKey != ""
}

func (m model) View() string {
	var content string

	switch m.state {
	case mainMenu:
		content = m.mainMenuView()
	case claudeCodeList:
		content = m.claudeCodeListView()
	case codexList:
		content = m.codexListView()
	case droidList:
		content = m.droidListView()
		// 操作菜单已移除
	case addClaudeCode:
		content = m.addConfigView("Claude Code")
	case addCodex:
		content = m.addConfigView("Codex")
	case addDroid:
		content = m.addDroidConfigView()
	case editClaudeCode:
		content = m.editConfigView("Claude Code")
	case editCodex:
		content = m.editConfigView("Codex")
	case editDroid:
		content = m.editDroidConfigView()
	case confirmDeleteClaudeCode:
		content = m.confirmDeleteView("Claude Code")
	case confirmDeleteCodex:
		content = m.confirmDeleteView("Codex")
	case confirmDeleteDroid:
		content = m.confirmDeleteView("Droid")
	case confirmExitAddClaudeCode:
		content = m.confirmExitAddView("Claude Code")
	case confirmExitAddCodex:
		content = m.confirmExitAddView("Codex")
	case confirmExitAddDroid:
		content = m.confirmExitAddView("Droid")
	}

	if m.error != "" {
		content += "\n\n" + errorStyle.Render(m.error)
	}

	return content
}

func (m model) mainMenuView() string {
	title := headerView("Codex/Claude Code/Droid配置切换器")

	activeClaude := "无"
	if active := m.config.GetActiveClaudeCode(); active != nil {
		activeClaude = active.Name
	}

	activeCodex := "无"
	if active := m.config.GetActiveCodex(); active != nil {
		activeCodex = active.Name
	}

	activeDroid := "无"
	if active := m.config.GetActiveDroid(); active != nil {
		activeDroid = active.ModelDisplayName
	}

	items := []string{
		fmt.Sprintf("🤖 Claude Code 配置 (当前: %s)", activeClaude),
		fmt.Sprintf("💻 Codex 配置 (当前: %s)", activeCodex),
		fmt.Sprintf("🔧 Droid 配置 (当前: %s)", activeDroid),
		"➕ 添加 Claude Code 配置",
		"➕ 添加 Codex 配置",
		"➕ 添加 Droid 配置",
		"🚪 退出程序",
	}

	var content strings.Builder
	content.WriteString(title)
	content.WriteString("\n\n")

	for i, item := range items {
		content.WriteString(menuItemView(item, m.cursor == i))
		if i < len(items)-1 {
			content.WriteString("\n")
		}
	}

	content.WriteString("\n")
	content.WriteString(statusBarView("↑/↓ 选择", "Enter 确认", "Esc 返回", ""))

	return content.String()
}

func (m model) addConfigView(serviceType string) string {
	title := headerView(fmt.Sprintf("添加 %s 配置", serviceType))

	var fields []struct {
		label string
		value string
	}

	if serviceType == "Codex" {
		// 设置默认值
		if m.formData.Model == "" {
			m.formData.Model = "gpt-5"
		}
		if m.formData.WireAPI == "" {
			m.formData.WireAPI = "responses"
		}
		if m.formData.ModelReasoningEffort == "" {
			m.formData.ModelReasoningEffort = "medium"
		}

		fields = []struct {
			label string
			value string
		}{
			{"配置名称", m.formData.Name},
			{"Provider", m.formData.Provider},
			{"Base URL", m.formData.BaseURL},
			{"API Key", m.formData.APIKey},
			{"Model", m.formData.Model},
			{"Wire API", m.formData.WireAPI},
			{"推理强度", m.formData.ModelReasoningEffort},
		}
	} else {
		fields = []struct {
			label string
			value string
		}{
			{"配置名称", m.formData.Name},
			{"Provider", m.formData.Provider},
			{"Base URL", m.formData.BaseURL},
			{"API Key", m.formData.APIKey},
		}
	}

	var content strings.Builder
	content.WriteString(title)
	content.WriteString("\n\n")

	var inner strings.Builder
	for i, field := range fields {
		prefix := "  "
		if m.formField == i {
			prefix = cursorStyle.Render(">")
		}

		highlight := ""
		if m.formField == i {
			highlight = fieldHighlightStyle.Render(" ◀")
		}

		// 对于Wire API字段，显示选择选项
		displayValue := field.value
		if serviceType == "Codex" && i == FieldWireAPI { // Wire API字段
			if m.formField == i {
				displayValue = field.value + " (←/→选择)"
			} else {
				displayValue = field.value
			}
		}
		// 对于推理强度字段，显示选择选项
		if serviceType == "Codex" && i == FieldModelReasoningEffort { // 推理强度字段
			if m.formField == i {
				displayValue = field.value + " (←/→选择)"
			} else {
				displayValue = field.value
			}
		}

		inner.WriteString(fmt.Sprintf("%s %s:%s %s\n", prefix, field.label, highlight, displayValue))
	}

	content.WriteString(boxStyle.Render(inner.String()))
	content.WriteString("\n")
	content.WriteString(statusBarView("Tab/↑/↓ 切字段", "Enter 保存", "Esc 取消", ""))

	return content.String()
}

func (m model) editConfigView(serviceType string) string {
	title := headerView(fmt.Sprintf("编辑 %s 配置", serviceType))

	var content strings.Builder
	content.WriteString(title)
	content.WriteString("\n\n")

	var fields []struct {
		label string
		value string
	}

	if serviceType == "Codex" {
		fields = []struct {
			label string
			value string
		}{
			{"配置名称", m.formData.Name},
			{"Provider", m.formData.Provider},
			{"Base URL", m.formData.BaseURL},
			{"API Key", m.formData.APIKey},
			{"Model", m.formData.Model},
			{"Wire API", m.formData.WireAPI},
			{"推理强度", m.formData.ModelReasoningEffort},
		}
	} else {
		fields = []struct {
			label string
			value string
		}{
			{"配置名称", m.formData.Name},
			{"Provider", m.formData.Provider},
			{"Base URL", m.formData.BaseURL},
			{"API Key", m.formData.APIKey},
		}
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

		// 对于Wire API字段，显示选择选项
		if serviceType == "Codex" && i == FieldWireAPI { // Wire API字段
			if m.formField == i {
				displayValue = field.value + " (←/→选择)"
			} else {
				displayValue = field.value
			}
		}

		// 对于推理强度字段，显示选择选项
		if serviceType == "Codex" && i == FieldModelReasoningEffort { // 推理强度字段
			if m.formField == i {
				displayValue = field.value + " (←/→选择)"
			} else {
				displayValue = field.value
			}
		}

		highlight := ""
		if m.formField == i {
			if serviceType == "Codex" && (i == FieldWireAPI || i == FieldModelReasoningEffort) { // Wire API和推理强度字段
				highlight = fieldHighlightStyle.Render(" ← 使用←/→选择")
			} else {
				highlight = fieldHighlightStyle.Render(" ← 正在编辑，请直接输入修改内容")
			}
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

// confirmDeleteView 显示删除确认对话框
func (m model) confirmDeleteView(serviceType string) string {
	var configName string
	if m.state == confirmDeleteClaudeCode && m.deleteIndex >= 0 && m.deleteIndex < len(m.config.ClaudeCode) {
		configName = m.config.ClaudeCode[m.deleteIndex].Name
	} else if m.state == confirmDeleteCodex && m.deleteIndex >= 0 && m.deleteIndex < len(m.config.Codex) {
		configName = m.config.Codex[m.deleteIndex].Name
	}

	title := headerView(fmt.Sprintf("删除 %s 配置", serviceType))

	var content strings.Builder
	content.WriteString(title)
	content.WriteString("\n\n")

	// 显示警告信息
	warning := fmt.Sprintf("⚠️  确定要删除配置 '%s' 吗？", configName)
	content.WriteString(errorStyle.Render(warning))
	content.WriteString("\n\n")
	content.WriteString("此操作无法撤销。")

	// 选项
	options := []string{
		"🗑️  确认删除",
		"❌ 取消",
	}

	content.WriteString("\n\n")
	for i, option := range options {
		prefix := "  "
		if m.cursor == i {
			prefix = cursorStyle.Render(">")
		}
		content.WriteString(fmt.Sprintf("%s %s\n", prefix, option))
	}

	content.WriteString("\n")
	content.WriteString(statusBarView("↑/↓/←/→ 选择", "Enter 确认", "Esc 取消", ""))

	return content.String()
}

// confirmExitAddView 显示退出添加配置确认对话框
func (m model) confirmExitAddView(serviceType string) string {
	title := headerView(fmt.Sprintf("退出添加 %s 配置", serviceType))

	var content strings.Builder
	content.WriteString(title)
	content.WriteString("\n\n")

	// 显示警告信息
	warning := "⚠️  确定要退出吗？表单中已填写的内容将被清空。"
	content.WriteString(errorStyle.Render(warning))
	content.WriteString("\n\n")
	content.WriteString("此操作无法撤销。")

	// 选项
	options := []string{
		"🚪 确认退出（清空内容）",
		"❌ 取消（继续编辑）",
	}

	content.WriteString("\n\n")
	for i, option := range options {
		prefix := "  "
		if m.cursor == i {
			prefix = cursorStyle.Render(">")
		}
		content.WriteString(fmt.Sprintf("%s %s\n", prefix, option))
	}

	content.WriteString("\n")
	content.WriteString(statusBarView("↑/↓/←/→ 选择", "Enter 确认", "Esc 返回编辑", ""))

	return content.String()
}
