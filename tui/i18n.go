package tui

var currentLang = "zh"

var translations = map[string]map[string]string{
	"zh": {
		// Main menu
		"app_title":           "Codex/Claude Code/Droid配置切换器",
		"menu_claude":         "🤖 Claude Code 配置 (当前: %s)",
		"menu_codex":          "💻 Codex 配置 (当前: %s)",
		"menu_droid":          "🔧 Droid 配置 (当前: %s)",
		"menu_add_claude":     "➕ 添加 Claude Code 配置",
		"menu_add_codex":      "➕ 添加 Codex 配置",
		"menu_add_droid":      "➕ 添加 Droid 配置",
		"menu_switch_lang":    "🌐 切换语言 / Switch Language",
		"menu_exit":           "🚪 退出程序",
		"none":                "无",

		// Navigation
		"nav_select":          "↑/↓ 选择",
		"nav_confirm":         "Enter 确认",
		"nav_back":            "Esc 返回",
		"nav_edit":            "Tab 编辑 Del 删除",
		"nav_add":             "A 添加",
		"nav_view":            "V 紧凑/展开",
		"nav_switch_service":  "←/→ 切服务",
		"nav_lang":            "L 切换语言",
		"back_to_menu":        "← 返回主菜单",

		// Headers
		"header_claude":       "Claude Code 配置",
		"header_codex":        "Codex 配置",
		"header_droid":        "Droid 配置",

		// Warnings
		"warn_mismatch":       "⚠️ 实际应用的 %s 设置与所选配置不一致: ",

		// Form titles
		"form_add":            "添加 %s 配置",
		"form_edit":           "编辑 %s 配置",

		// Form fields
		"field_name":          "配置名称",
		"field_provider":      "Provider",
		"field_base_url":      "Base URL",
		"field_api_key":       "API Key",
		"field_model":         "Model",
		"field_wire_api":      "Wire API",
		"field_auth_method":   "认证方式",
		"field_reasoning":     "推理强度",
		"field_display_name":  "模型显示名称",
		"field_model_name":    "模型名称",

		// Form hints
		"hint_select":         "(←/→选择)",
		"hint_editing":        "(编辑中)",
		"hint_input":          "← 正在编辑，请直接输入修改内容",
		"hint_use_arrows":     "← 使用←/→选择",
		"hint_current_edit":   "✨ 当前正在编辑: ",
		"hint_apikey_visible": "(API Key 正在显示完整内容以便编辑)",

		// Form navigation
		"form_nav_field":      "Tab/↑/↓ 切字段",
		"form_nav_save":       "Enter 保存",
		"form_nav_cancel":     "Esc 取消",

		// Confirmation dialogs
		"confirm_delete_title":   "删除 %s 配置",
		"confirm_delete_warn":    "⚠️  确定要删除配置 '%s' 吗？",
		"confirm_delete_msg":     "此操作无法撤销。",
		"confirm_delete_yes":     "🗑️  确认删除",
		"confirm_delete_no":      "❌ 取消",

		"confirm_exit_title":     "退出添加 %s 配置",
		"confirm_exit_warn":      "⚠️  确定要退出吗？表单中已填写的内容将被清空。",
		"confirm_exit_msg":       "此操作无法撤销。",
		"confirm_exit_yes":       "🚪 确认退出（清空内容）",
		"confirm_exit_no":        "❌ 取消（继续编辑）",

		"confirm_nav":            "↑/↓/←/→ 选择",
		"confirm_nav_back":       "Esc 取消/返回编辑",

		// Success messages
		"success_add_claude":     "✅ Claude Code 配置添加成功！",
		"success_add_codex":      "✅ Codex 配置添加成功！",
		"success_add_droid":      "✅ Droid 配置添加成功！",
		"success_update_claude":  "✅ Claude Code 配置更新成功！",
		"success_update_codex":   "✅ Codex 配置更新成功！",
		"success_update_droid":   "✅ Droid 配置更新成功！",
		"success_switch_claude":  "✅ Claude Code 配置切换成功！",
		"success_switch_codex":   "✅ Codex 配置切换成功！",
		"success_switch_droid":   "✅ Droid 配置切换成功！",
		"success_delete_claude":  "✅ Claude Code 配置 '%s' 删除成功！",
		"success_delete_codex":   "✅ Codex 配置 '%s' 删除成功！",
		"success_delete_droid":   "✅ Droid 配置 '%s' 删除成功！",
		"success_lang_switch":    "✅ 语言已切换！",

		// Error messages
		"error_fill_all":         "⚠️ 请填写所有字段",
		"error_config_index":     "❌ 配置索引错误",
		"error_switch_claude":    "切换 Claude Code 配置失败: %v",
		"error_switch_codex":     "切换 Codex 配置失败: %v",
		"error_switch_droid":     "切换 Droid 配置失败: %v",

		// Display
		"display_active":         "[当前使用]",
		"display_name":           "配置名称: %s",
		"display_provider":       "Provider: %s",
		"display_base_url":       "Base URL: %s",

		// CLI messages
		"cli_error_load":         "Error loading config: %v\n",
		"cli_not_found_codex":    "Codex config not found: %s\n",
		"cli_switch_fail_codex":  "Switch Codex failed: %v\n",
		"cli_active_fail_codex":  "Set active Codex failed: %v\n",
		"cli_not_found_claude":   "Claude Code config not found: %s\n",
		"cli_switch_fail_claude": "Switch Claude Code failed: %v\n",
		"cli_active_fail_claude": "Set active Claude Code failed: %v\n",
		"cli_not_found_droid":    "Droid config not found: %s\n",
		"cli_switch_fail_droid":  "Switch Droid failed: %v\n",
		"cli_active_fail_droid":  "Set active Droid failed: %v\n",
		"cli_error_tui":          "Error running TUI: %v\n",
		"cli_switched_codex":     "Switched Codex to '%s' (provider=%s)\n",
		"cli_switched_claude":    "Switched Claude Code to '%s' (provider=%s)\n",
		"cli_switched_droid":     "Switched Droid to '%s' (provider=%s)\n",
	},
	"en": {
		// Main menu
		"app_title":           "Codex/Claude Code/Droid Configuration Switcher",
		"menu_claude":         "🤖 Claude Code Config (Current: %s)",
		"menu_codex":          "💻 Codex Config (Current: %s)",
		"menu_droid":          "🔧 Droid Config (Current: %s)",
		"menu_add_claude":     "➕ Add Claude Code Config",
		"menu_add_codex":      "➕ Add Codex Config",
		"menu_add_droid":      "➕ Add Droid Config",
		"menu_switch_lang":    "🌐 Switch Language / 切换语言",
		"menu_exit":           "🚪 Exit",
		"none":                "None",

		// Navigation
		"nav_select":          "↑/↓ Select",
		"nav_confirm":         "Enter Confirm",
		"nav_back":            "Esc Back",
		"nav_edit":            "Tab Edit Del Delete",
		"nav_add":             "A Add",
		"nav_view":            "V Compact/Expand",
		"nav_switch_service":  "←/→ Switch Service",
		"nav_lang":            "L Switch Language",
		"back_to_menu":        "← Back to Main Menu",

		// Headers
		"header_claude":       "Claude Code Configuration",
		"header_codex":        "Codex Configuration",
		"header_droid":        "Droid Configuration",

		// Warnings
		"warn_mismatch":       "⚠️ Applied %s settings don't match selected config: ",

		// Form titles
		"form_add":            "Add %s Configuration",
		"form_edit":           "Edit %s Configuration",

		// Form fields
		"field_name":          "Configuration Name",
		"field_provider":      "Provider",
		"field_base_url":      "Base URL",
		"field_api_key":       "API Key",
		"field_model":         "Model",
		"field_wire_api":      "Wire API",
		"field_auth_method":   "Auth Method",
		"field_reasoning":     "Reasoning Effort",
		"field_display_name":  "Model Display Name",
		"field_model_name":    "Model Name",

		// Form hints
		"hint_select":         "(←/→ select)",
		"hint_editing":        "(editing)",
		"hint_input":          "← Editing, input directly to modify",
		"hint_use_arrows":     "← Use ←/→ to select",
		"hint_current_edit":   "✨ Currently editing: ",
		"hint_apikey_visible": "(API Key is fully visible for editing)",

		// Form navigation
		"form_nav_field":      "Tab/↑/↓ Switch Field",
		"form_nav_save":       "Enter Save",
		"form_nav_cancel":     "Esc Cancel",

		// Confirmation dialogs
		"confirm_delete_title":   "Delete %s Configuration",
		"confirm_delete_warn":    "⚠️  Are you sure you want to delete config '%s'?",
		"confirm_delete_msg":     "This action cannot be undone.",
		"confirm_delete_yes":     "🗑️  Confirm Delete",
		"confirm_delete_no":      "❌ Cancel",

		"confirm_exit_title":     "Exit Add %s Configuration",
		"confirm_exit_warn":      "⚠️  Are you sure you want to exit? Form content will be cleared.",
		"confirm_exit_msg":       "This action cannot be undone.",
		"confirm_exit_yes":       "🚪 Confirm Exit (Clear Content)",
		"confirm_exit_no":        "❌ Cancel (Continue Editing)",

		"confirm_nav":            "↑/↓/←/→ Select",
		"confirm_nav_back":       "Esc Cancel/Back to Edit",

		// Success messages
		"success_add_claude":     "✅ Claude Code configuration added successfully!",
		"success_add_codex":      "✅ Codex configuration added successfully!",
		"success_add_droid":      "✅ Droid configuration added successfully!",
		"success_update_claude":  "✅ Claude Code configuration updated successfully!",
		"success_update_codex":   "✅ Codex configuration updated successfully!",
		"success_update_droid":   "✅ Droid configuration updated successfully!",
		"success_switch_claude":  "✅ Claude Code configuration switched successfully!",
		"success_switch_codex":   "✅ Codex configuration switched successfully!",
		"success_switch_droid":   "✅ Droid configuration switched successfully!",
		"success_delete_claude":  "✅ Claude Code configuration '%s' deleted successfully!",
		"success_delete_codex":   "✅ Codex configuration '%s' deleted successfully!",
		"success_delete_droid":   "✅ Droid configuration '%s' deleted successfully!",
		"success_lang_switch":    "✅ Language switched!",

		// Error messages
		"error_fill_all":         "⚠️ Please fill in all fields",
		"error_config_index":     "❌ Configuration index error",
		"error_switch_claude":    "Failed to switch Claude Code config: %v",
		"error_switch_codex":     "Failed to switch Codex config: %v",
		"error_switch_droid":     "Failed to switch Droid config: %v",

		// Display
		"display_active":         "[Active]",
		"display_name":           "Name: %s",
		"display_provider":       "Provider: %s",
		"display_base_url":       "Base URL: %s",

		// CLI messages
		"cli_error_load":         "Error loading config: %v\n",
		"cli_not_found_codex":    "Codex config not found: %s\n",
		"cli_switch_fail_codex":  "Switch Codex failed: %v\n",
		"cli_active_fail_codex":  "Set active Codex failed: %v\n",
		"cli_not_found_claude":   "Claude Code config not found: %s\n",
		"cli_switch_fail_claude": "Switch Claude Code failed: %v\n",
		"cli_active_fail_claude": "Set active Claude Code failed: %v\n",
		"cli_not_found_droid":    "Droid config not found: %s\n",
		"cli_switch_fail_droid":  "Switch Droid failed: %v\n",
		"cli_active_fail_droid":  "Set active Droid failed: %v\n",
		"cli_error_tui":          "Error running TUI: %v\n",
		"cli_switched_codex":     "Switched Codex to '%s' (provider=%s)\n",
		"cli_switched_claude":    "Switched Claude Code to '%s' (provider=%s)\n",
		"cli_switched_droid":     "Switched Droid to '%s' (provider=%s)\n",
	},
}

func t(key string) string {
	if trans, ok := translations[currentLang][key]; ok {
		return trans
	}
	return key
}

func SetLanguage(lang string) {
	if lang == "zh" || lang == "en" {
		currentLang = lang
	}
}

func GetLanguage() string {
	return currentLang
}

func ToggleLanguage() string {
	if currentLang == "zh" {
		currentLang = "en"
	} else {
		currentLang = "zh"
	}
	return currentLang
}
