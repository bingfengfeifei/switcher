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
	FieldAuthMethod
	FieldModelReasoningEffort
	FieldClaudeDefaultHaikuModel
	FieldClaudeDefaultOpusModel
	FieldClaudeDefaultSonnetModel
)

// 配置类型字段数量
const (
	ClaudeCodeFieldCount = 7 // Name, Provider, BaseURL, APIKey, HaikuModel, OpusModel, SonnetModel
	CodexFieldCount      = 8
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
	windowHeight     int             // 终端窗口高度
}

func (m model) hasFormContent() bool {
	return m.formData.Name != "" || m.formData.Provider != "" || m.formData.BaseURL != "" || m.formData.APIKey != "" || m.formData.Model != "" || m.formData.WireAPI != "" || m.formData.EnvKey != "" || m.formData.ModelReasoningEffort != "" || m.formData.ClaudeDefaultHaikuModel != "" || m.formData.ClaudeDefaultOpusModel != "" || m.formData.ClaudeDefaultSonnetModel != ""
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
	version := GetVersion()
	title := headerViewWithVersion(t("app_title"), version)

	activeClaude := t("none")
	if active := m.config.GetActiveClaudeCode(); active != nil {
		activeClaude = active.Name
	}

	activeCodex := t("none")
	if active := m.config.GetActiveCodex(); active != nil {
		activeCodex = active.Name
	}

	activeDroid := t("none")
	if active := m.config.GetActiveDroid(); active != nil {
		activeDroid = active.ModelDisplayName
	}

	items := []string{
		fmt.Sprintf(t("menu_claude"), activeClaude),
		fmt.Sprintf(t("menu_codex"), activeCodex),
		fmt.Sprintf(t("menu_droid"), activeDroid),
		t("menu_add_claude"),
		t("menu_add_codex"),
		t("menu_add_droid"),
		t("menu_switch_lang"),
		t("menu_exit"),
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
	content.WriteString(statusBarView(t("nav_select"), t("nav_confirm"), t("nav_back"), t("nav_lang")))

	return content.String()
}

func (m model) addConfigView(serviceType string) string {
	title := headerView(fmt.Sprintf(t("form_add"), serviceType))

	var fields []struct {
		label string
		value string
	}

	if serviceType == "Codex" {
		// 设置默认值
		if m.formData.Model == "" {
			m.formData.Model = DefaultCodexModel
		}
		if m.formData.WireAPI == "" {
			m.formData.WireAPI = DefaultWireAPI
		}
		if m.formData.AuthMethod == "" {
			m.formData.AuthMethod = "auth.json"
		}
		if m.formData.ModelReasoningEffort == "" {
			m.formData.ModelReasoningEffort = DefaultModelReasoningEffort
		}

		fields = []struct {
			label string
			value string
		}{
			{t("field_name"), m.formData.Name},
			{t("field_provider"), m.formData.Provider},
			{t("field_base_url"), m.formData.BaseURL},
			{t("field_api_key"), m.formData.APIKey},
			{t("field_model"), m.formData.Model},
			{t("field_wire_api"), m.formData.WireAPI},
			{t("field_auth_method"), m.formData.AuthMethod},
			{t("field_reasoning"), m.formData.ModelReasoningEffort},
		}
	} else {
		fields = []struct {
			label string
			value string
		}{
			{t("field_name"), m.formData.Name},
			{t("field_provider"), m.formData.Provider},
			{t("field_base_url"), m.formData.BaseURL},
			{t("field_api_key"), m.formData.APIKey},
			{t("field_haiku_model"), m.formData.ClaudeDefaultHaikuModel},
			{t("field_opus_model"), m.formData.ClaudeDefaultOpusModel},
			{t("field_sonnet_model"), m.formData.ClaudeDefaultSonnetModel},
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
				displayValue = field.value + " " + t("hint_select")
			} else {
				displayValue = field.value
			}
		}
		// 对于推理强度字段，显示选择选项
		if serviceType == "Codex" && i == FieldModelReasoningEffort { // 推理强度字段
			if m.formField == i {
				displayValue = field.value + " " + t("hint_select")
			} else {
				displayValue = field.value
			}
		}

		inner.WriteString(fmt.Sprintf("%s %s:%s %s\n", prefix, field.label, highlight, displayValue))
	}

	content.WriteString(boxStyle.Render(inner.String()))
	content.WriteString("\n")
	content.WriteString(statusBarView(t("form_nav_field"), t("form_nav_save"), t("form_nav_cancel"), ""))

	return content.String()
}

func (m model) editConfigView(serviceType string) string {
	title := headerView(fmt.Sprintf(t("form_edit"), serviceType))

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
			{t("field_name"), m.formData.Name},
			{t("field_provider"), m.formData.Provider},
			{t("field_base_url"), m.formData.BaseURL},
			{t("field_api_key"), m.formData.APIKey},
			{t("field_model"), m.formData.Model},
			{t("field_wire_api"), m.formData.WireAPI},
			{t("field_auth_method"), m.formData.AuthMethod},
			{t("field_reasoning"), m.formData.ModelReasoningEffort},
		}
	} else {
		fields = []struct {
			label string
			value string
		}{
			{t("field_name"), m.formData.Name},
			{t("field_provider"), m.formData.Provider},
			{t("field_base_url"), m.formData.BaseURL},
			{t("field_api_key"), m.formData.APIKey},
			{t("field_haiku_model"), m.formData.ClaudeDefaultHaikuModel},
			{t("field_opus_model"), m.formData.ClaudeDefaultOpusModel},
			{t("field_sonnet_model"), m.formData.ClaudeDefaultSonnetModel},
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
			displayValue = field.value + " " + t("hint_editing")
		}

		// 对于Wire API字段，显示选择选项
		if serviceType == "Codex" && i == FieldWireAPI { // Wire API字段
			if m.formField == i {
				displayValue = field.value + " " + t("hint_select")
			} else {
				displayValue = field.value
			}
		}

		// 对于推理强度字段，显示选择选项
		if serviceType == "Codex" && i == FieldModelReasoningEffort { // 推理强度字段
			if m.formField == i {
				displayValue = field.value + " " + t("hint_select")
			} else {
				displayValue = field.value
			}
		}

		highlight := ""
		if m.formField == i {
			if serviceType == "Codex" && (i == FieldWireAPI || i == FieldModelReasoningEffort) { // Wire API和推理强度字段
				highlight = fieldHighlightStyle.Render(" " + t("hint_use_arrows"))
			} else {
				highlight = fieldHighlightStyle.Render(" " + t("hint_input"))
			}
		}

		inner.WriteString(formRowStyle.Render(fmt.Sprintf("%s %s:%s %s", prefix, field.label, highlight, displayValue)) + "\n")
	}

	content.WriteString(boxStyle.Render(inner.String()))
	content.WriteString("\n")
	content.WriteString(statusBarView(t("form_nav_field"), t("form_nav_save"), t("form_nav_cancel"), ""))

	// 添加当前编辑状态提示
	if m.formField >= 0 && m.formField < DroidFieldCount {
		content.WriteString("\n" + fieldHighlightStyle.Render(t("hint_current_edit")) + fields[m.formField].label)
		if m.formField == FieldAPIKey {
			content.WriteString("\n" + fieldHighlightStyle.Render("   " + t("hint_apikey_visible")))
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

	title := headerView(fmt.Sprintf(t("confirm_delete_title"), serviceType))

	var content strings.Builder
	content.WriteString(title)
	content.WriteString("\n\n")

	// 显示警告信息
	warning := fmt.Sprintf(t("confirm_delete_warn"), configName)
	content.WriteString(errorStyle.Render(warning))
	content.WriteString("\n\n")
	content.WriteString(t("confirm_delete_msg"))

	// 选项
	options := []string{
		t("confirm_delete_yes"),
		t("confirm_delete_no"),
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
	content.WriteString(statusBarView(t("confirm_nav"), t("nav_confirm"), t("confirm_nav_back"), ""))

	return content.String()
}

// confirmExitAddView 显示退出添加配置确认对话框
func (m model) confirmExitAddView(serviceType string) string {
	title := headerView(fmt.Sprintf(t("confirm_exit_title"), serviceType))

	var content strings.Builder
	content.WriteString(title)
	content.WriteString("\n\n")

	// 显示警告信息
	warning := t("confirm_exit_warn")
	content.WriteString(errorStyle.Render(warning))
	content.WriteString("\n\n")
	content.WriteString(t("confirm_exit_msg"))

	// 选项
	options := []string{
		t("confirm_exit_yes"),
		t("confirm_exit_no"),
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
	content.WriteString(statusBarView(t("confirm_nav"), t("nav_confirm"), t("confirm_nav_back"), ""))

	return content.String()
}
