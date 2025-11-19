package tui

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
)

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		// 更新窗口高度
		m.windowHeight = msg.Height
		return m, nil
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyCtrlC:
			if m.state != mainMenu {
				m.state = mainMenu
				m.cursor = 0
				m.error = ""
				return m, nil
			}
			return m, tea.Quit
		case tea.KeyEsc:
			switch m.state {
			case mainMenu:
				return m, tea.Quit
			case addClaudeCode:
				if m.hasFormContent() {
					m.state = confirmExitAddClaudeCode
					m.cursor = 1 // 默认选择"否"
				} else {
					m.state = mainMenu
					m.cursor = 0
					m.error = ""
				}
			case addCodex:
				if m.hasFormContent() {
					m.state = confirmExitAddCodex
					m.cursor = 1 // 默认选择"否"
				} else {
					m.state = mainMenu
					m.cursor = 0
					m.error = ""
				}
			case addDroid:
				if m.hasDroidFormContent() {
					m.state = confirmExitAddDroid
					m.cursor = 1 // 默认选择"否"
				} else {
					m.state = mainMenu
					m.cursor = 0
					m.error = ""
				}
			case claudeCodeList, codexList, droidList:
				// 退出配置列表时清空排序后的列表，确保下次进入时重新排序
				m.sortedClaudeCode = nil
				m.sortedCodex = nil
				m.sortedDroid = nil
				m.state = mainMenu
				m.cursor = 0
				m.error = ""
			default:
				// 编辑状态应该回到对应的列表
				switch m.state {
				case editClaudeCode:
					m.state = claudeCodeList
				case editCodex:
					m.state = codexList
				case editDroid:
					m.state = droidList
				default:
					m.state = mainMenu
				}
				m.cursor = 0
				m.error = ""
			}
			return m, nil

		case tea.KeyUp:
			if m.state == addClaudeCode || m.state == editClaudeCode || m.state == addDroid || m.state == editDroid {
				// Claude Code 和 Droid: 4个字段，在字段之间导航
				if m.cursor > DroidFieldCount-1 {
					m.cursor = DroidFieldCount - 1
					m.formField = m.cursor
				} else if m.cursor >= 0 {
					// 在字段之间移动，确保cursor和formField同步
					m.cursor--
					m.formField = m.cursor
				} else {
					// 从第一个字段循环到最后一个字段
					m.cursor = DroidFieldCount - 1
					m.formField = m.cursor
				}

			} else if m.state == addCodex || m.state == editCodex {
				// Codex: 6个字段，在字段和按钮之间导航
				if m.cursor > CodexFieldCount-1 {
					m.cursor = CodexFieldCount - 1
					m.formField = m.cursor
				} else if m.cursor >= 0 {
					// 在字段之间移动，确保cursor和formField同步
					m.cursor--
					m.formField = m.cursor
				} else {
					// 从第一个字段循环到最后一个字段
					m.cursor = CodexFieldCount - 1
					m.formField = m.cursor
				}
			} else if m.state == confirmDeleteClaudeCode || m.state == confirmDeleteCodex || m.state == confirmDeleteDroid || m.state == confirmExitAddClaudeCode || m.state == confirmExitAddCodex || m.state == confirmExitAddDroid {
				m.cursor = 0 // 只能在两个选项间切换
			} else if m.cursor > 0 {
				m.cursor--
			}
		case tea.KeyDown:
			if m.state == addClaudeCode || m.state == editClaudeCode || m.state == addDroid || m.state == editDroid {
				// Claude Code 和 Droid: 4个字段，只在字段之间导航
				if m.cursor < DroidFieldCount-1 {
					// 在字段之间切换时，更新formField为当前位置
					m.formField = m.cursor + 1
					m.cursor++
				} else {
					// 从最后一个字段回到第一个字段
					m.cursor = 0
					m.formField = 0
				}

			} else if m.state == addCodex || m.state == editCodex {
				// Codex: 6个字段，只在字段之间导航
				if m.cursor < CodexFieldCount-1 {
					// 在字段之间切换时，更新formField为当前位置
					m.formField = m.cursor + 1
					m.cursor++
				} else {
					// 从最后一个字段回到第一个字段
					m.cursor = 0
					m.formField = 0
				}
			} else if m.state == confirmDeleteClaudeCode || m.state == confirmDeleteCodex || m.state == confirmDeleteDroid || m.state == confirmExitAddClaudeCode || m.state == confirmExitAddCodex || m.state == confirmExitAddDroid {
				m.cursor = 1 // 只能在两个选项间切换
			} else if m.cursor < m.getMaxCursor() {
				m.cursor++
			}
		case tea.KeyRunes:
			// 在添加/编辑模式下，将任意文本输入写入当前字段；
			// 否则在列表/菜单中用 j/k 导航。
			if m.state == addClaudeCode || m.state == addCodex || m.state == addDroid || m.state == editClaudeCode || m.state == editCodex || m.state == editDroid {
				return m.handleInput(msg.String())
			}
			switch msg.Runes[0] {
			case 'k', 'K':
				if m.cursor > 0 {
					m.cursor--
				}
			case 'j', 'J':
				if m.cursor < m.getMaxCursor() {
					m.cursor++
				}
			case 'v', 'V':
				// 切换紧凑/展开模式
				m.compact = !m.compact
			case 'l', 'L':
				// 切换语言
				ToggleLanguage()
				m.config.Language = GetLanguage()
				m.config.Save()
				m.error = t("success_lang_switch")
			case 'a', 'A':
				if m.state == claudeCodeList {
					m.state = addClaudeCode
					m.formData = ServiceConfig{}
					m.formField = 0
					m.cursor = 4
					m.error = ""
				} else if m.state == codexList {
					m.state = addCodex
					m.formData = ServiceConfig{}
					m.formField = 0
					m.cursor = 0
					m.error = ""
				} else if m.state == droidList {
					m.state = addDroid
					m.droidFormData = DroidConfig{}
					m.formField = 0
					m.cursor = 4
					m.error = ""
				}
			}
		case tea.KeyLeft:
			if m.state == codexList {
				m.state = claudeCodeList
				if m.cursor > len(m.config.ClaudeCode) {
					m.cursor = len(m.config.ClaudeCode)
				}
			} else if m.state == droidList {
				m.state = codexList
				if m.cursor > len(m.config.Codex) {
					m.cursor = len(m.config.Codex)
				}
			} else if m.state == confirmDeleteClaudeCode || m.state == confirmDeleteCodex || m.state == confirmDeleteDroid || m.state == confirmExitAddClaudeCode || m.state == confirmExitAddCodex || m.state == confirmExitAddDroid {
				m.cursor = 0 // 选择"确认删除"或"确认退出"
			} else if (m.state == addCodex || m.state == editCodex) && m.formField == FieldWireAPI {
				// Wire API字段：在chat和responses之间切换
				if m.formData.WireAPI == DefaultWireAPI {
					m.formData.WireAPI = "chat"
				} else {
					m.formData.WireAPI = DefaultWireAPI
				}
			} else if (m.state == addCodex || m.state == editCodex) && m.formField == FieldAuthMethod {
				// 认证方式字段：在auth.json和env之间切换
				if m.formData.AuthMethod == "auth.json" {
					m.formData.AuthMethod = "env"
				} else {
					m.formData.AuthMethod = "auth.json"
				}
			} else if (m.state == addCodex || m.state == editCodex) && m.formField == FieldModelReasoningEffort {
				// 推理强度字段：在low、medium、high之间切换
				if m.formData.ModelReasoningEffort == "low" {
					m.formData.ModelReasoningEffort = "high"
				} else if m.formData.ModelReasoningEffort == "high" {
					m.formData.ModelReasoningEffort = DefaultModelReasoningEffort
				} else {
					m.formData.ModelReasoningEffort = "low"
				}
			}
		case tea.KeyRight:
			if m.state == claudeCodeList {
				m.state = codexList
				if m.cursor > len(m.config.Codex) {
					m.cursor = len(m.config.Codex)
				}
			} else if m.state == codexList {
				m.state = droidList
				if m.cursor > len(m.config.Droid) {
					m.cursor = len(m.config.Droid)
				}
			} else if m.state == confirmDeleteClaudeCode || m.state == confirmDeleteCodex || m.state == confirmDeleteDroid || m.state == confirmExitAddClaudeCode || m.state == confirmExitAddCodex || m.state == confirmExitAddDroid {
				m.cursor = 1 // 选择"取消"
			} else if (m.state == addCodex || m.state == editCodex) && m.formField == FieldWireAPI {
				// Wire API字段：在chat和responses之间切换
				if m.formData.WireAPI == DefaultWireAPI {
					m.formData.WireAPI = "chat"
				} else {
					m.formData.WireAPI = DefaultWireAPI
				}
			} else if (m.state == addCodex || m.state == editCodex) && m.formField == FieldAuthMethod {
				// 认证方式字段：在auth.json和env之间切换
				if m.formData.AuthMethod == "auth.json" {
					m.formData.AuthMethod = "env"
				} else {
					m.formData.AuthMethod = "auth.json"
				}
			} else if (m.state == addCodex || m.state == editCodex) && m.formField == FieldModelReasoningEffort {
				// 推理强度字段：在low、medium、high之间切换
				if m.formData.ModelReasoningEffort == ModelReasoningEffortLow {
					m.formData.ModelReasoningEffort = ModelReasoningEffortMedium
				} else if m.formData.ModelReasoningEffort == ModelReasoningEffortMedium {
					m.formData.ModelReasoningEffort = ModelReasoningEffortHigh
				} else {
					m.formData.ModelReasoningEffort = ModelReasoningEffortLow
				}
			}
		case tea.KeyEnter:
			// 在添加/编辑状态下，Enter直接保存
			if m.state == addClaudeCode {
				if m.formData.Name != "" && m.formData.Provider != "" && m.formData.BaseURL != "" && m.formData.APIKey != "" {
					err := m.config.AddClaudeCodeConfig(m.formData)
					if err != nil {
						m.error = err.Error()
					} else {
						m.error = t("success_add_claude")
						m.state = claudeCodeList
						m.cursor = 0
					}
				} else {
					m.error = t("error_fill_all")
				}
			} else if m.state == addCodex {
				if m.formData.Name != "" && m.formData.Provider != "" && m.formData.BaseURL != "" && m.formData.APIKey != "" {
					// Set default values for Codex config
					if m.formData.Model == "" {
						m.formData.Model = DefaultCodexModel
					}
					if m.formData.WireAPI == "" {
						m.formData.WireAPI = DefaultWireAPI
					}
					if m.formData.AuthMethod == "" {
						m.formData.AuthMethod = "auth.json"
					}
					m.formData.EnvKey = DefaultEnvKey // Always set to CODEX_KEY

					if err := m.config.AddCodexConfig(m.formData); err != nil {
						m.error = err.Error()
					} else {
						m.error = t("success_add_codex")
						m.state = codexList
						m.cursor = 0
					}
				} else {
					m.error = t("error_fill_all")
				}
			} else if m.state == addDroid {
				if m.droidFormData.ModelDisplayName != "" && m.droidFormData.Model != "" && m.droidFormData.BaseURL != "" && m.droidFormData.APIKey != "" {
					if err := m.config.AddDroidConfig(m.droidFormData); err != nil {
						m.error = err.Error()
					} else {
						m.error = t("success_add_droid")
						m.state = droidList
						m.cursor = 0
					}
				} else {
					m.error = t("error_fill_all")
				}
			} else if m.state == editClaudeCode {
				if m.formData.Name != "" && m.formData.Provider != "" && m.formData.BaseURL != "" && m.formData.APIKey != "" {
					m.config.ClaudeCode[m.editIndex] = m.formData
					err := m.config.Save()
					if err != nil {
						m.error = err.Error()
					} else {
						m.error = t("success_update_claude")
						m.state = claudeCodeList
						m.cursor = 0
					}
				} else {
					m.error = t("error_fill_all")
				}
			} else if m.state == editCodex {
				if m.formData.Name != "" && m.formData.Provider != "" && m.formData.BaseURL != "" && m.formData.APIKey != "" {
					// Set default values for Codex config
					if m.formData.Model == "" {
						m.formData.Model = DefaultCodexModel
					}
					if m.formData.WireAPI == "" {
						m.formData.WireAPI = DefaultWireAPI
					}
					if m.formData.ModelReasoningEffort == "" {
						m.formData.ModelReasoningEffort = DefaultModelReasoningEffort
					}
					m.formData.EnvKey = DefaultEnvKey // Always set to CODEX_KEY

					m.config.Codex[m.editIndex] = m.formData
					err := m.config.Save()
					if err != nil {
						m.error = err.Error()
					} else {
						m.error = t("success_update_codex")
						m.state = codexList
						m.cursor = 0
					}
				} else {
					m.error = t("error_fill_all")
				}
			} else if m.state == editDroid {
				if m.droidFormData.ModelDisplayName != "" && m.droidFormData.Model != "" && m.droidFormData.BaseURL != "" && m.droidFormData.APIKey != "" {
					m.config.Droid[m.editIndex] = m.droidFormData
					err := m.config.Save()
					if err != nil {
						m.error = err.Error()
					} else {
						m.error = t("success_update_droid")
						m.state = droidList
						m.cursor = 0
					}
				} else {
					m.error = t("error_fill_all")
				}
			} else {
				return m.handleSelect()
			}
		case tea.KeySpace:
			// 编辑/新增模式下，空格应作为字符输入
			if m.state == addClaudeCode || m.state == addCodex || m.state == addDroid || m.state == editClaudeCode || m.state == editCodex || m.state == editDroid {
				return m.handleInput(" ")
			}
			return m.handleSelect()
		case tea.KeyTab:
			// 列表界面：Tab 进入编辑；编辑/新增界面：Tab 切换字段
			if m.state == claudeCodeList {
				if m.cursor < len(m.sortedClaudeCode) {
					originalIndex := findConfigIndex(m.config.ClaudeCode, m.sortedClaudeCode[m.cursor])
					if originalIndex == -1 {
						m.error = t("error_config_index")
						break
					}
					m.editIndex = originalIndex
					m.formData = m.config.ClaudeCode[m.editIndex]
					m.state = editClaudeCode
					m.cursor = 0
					m.formField = 0
					m.error = ""
				}
			} else if m.state == codexList {
				if m.cursor < len(m.sortedCodex) {
					originalIndex := findConfigIndex(m.config.Codex, m.sortedCodex[m.cursor])
					if originalIndex == -1 {
						m.error = t("error_config_index")
						break
					}
					m.editIndex = originalIndex
					m.formData = m.config.Codex[m.editIndex]
					m.state = editCodex
					m.cursor = 0 // Codex有6个字段，从第一个字段开始
					m.formField = 0
					m.error = ""
				}
			} else if m.state == droidList {
				if m.cursor < len(m.sortedDroid) {
					originalIndex := findDroidConfigIndex(m.config.Droid, m.sortedDroid[m.cursor])
					if originalIndex == -1 {
						m.error = t("error_config_index")
						break
					}
					m.editIndex = originalIndex
					m.droidFormData = m.config.Droid[m.editIndex]
					m.state = editDroid
					m.cursor = 4
					m.formField = 0
					m.error = ""
				}
			} else if m.state == addClaudeCode || m.state == editClaudeCode {
				// Claude Code有4个字段
				m.formField = (m.formField + 1) % ClaudeCodeFieldCount
			} else if m.state == addCodex || m.state == editCodex {
				// Codex有6个字段(不包含EnvKey)
				m.formField = (m.formField + 1) % CodexFieldCount
			} else if m.state == addDroid || m.state == editDroid {
				// Droid有4个字段
				m.formField = (m.formField + 1) % DroidFieldCount
			}
		case tea.KeyCtrlS:
			// Ctrl+S 直接保存编辑/新增
			if m.state == editClaudeCode {
				if m.formData.Name != "" && m.formData.Provider != "" && m.formData.BaseURL != "" && m.formData.APIKey != "" {
					m.config.ClaudeCode[m.editIndex] = m.formData
					if err := m.config.Save(); err != nil {
						m.error = err.Error()
					} else {
						m.error = t("success_update_claude")
						m.state = claudeCodeList
						m.cursor = 0
					}
				} else {
					m.error = t("error_fill_all")
				}
			} else if m.state == addClaudeCode {
				if m.formData.Name != "" && m.formData.Provider != "" && m.formData.BaseURL != "" && m.formData.APIKey != "" {
					if err := m.config.AddClaudeCodeConfig(m.formData); err != nil {
						m.error = err.Error()
					} else {
						m.error = t("success_add_claude")
						m.state = claudeCodeList
						m.cursor = 0
					}
				} else {
					m.error = t("error_fill_all")
				}
			} else if m.state == editCodex {
				if m.formData.Name != "" && m.formData.Provider != "" && m.formData.BaseURL != "" && m.formData.APIKey != "" {
					// Set default values for Codex config
					if m.formData.Model == "" {
						m.formData.Model = DefaultCodexModel
					}
					if m.formData.WireAPI == "" {
						m.formData.WireAPI = DefaultWireAPI
					}
					if m.formData.ModelReasoningEffort == "" {
						m.formData.ModelReasoningEffort = DefaultModelReasoningEffort
					}
					m.formData.EnvKey = DefaultEnvKey // Always set to CODEX_KEY

					m.config.Codex[m.editIndex] = m.formData
					if err := m.config.Save(); err != nil {
						m.error = err.Error()
					} else {
						m.error = t("success_update_codex")
						m.state = codexList
						m.cursor = 0
					}
				} else {
					m.error = t("error_fill_all")
				}
			} else if m.state == addCodex {
				if m.formData.Name != "" && m.formData.Provider != "" && m.formData.BaseURL != "" && m.formData.APIKey != "" {
					// Set default values for Codex config
					if m.formData.Model == "" {
						m.formData.Model = DefaultCodexModel
					}
					if m.formData.WireAPI == "" {
						m.formData.WireAPI = DefaultWireAPI
					}
					if m.formData.AuthMethod == "" {
						m.formData.AuthMethod = "auth.json"
					}
					m.formData.EnvKey = DefaultEnvKey // Always set to CODEX_KEY

					if err := m.config.AddCodexConfig(m.formData); err != nil {
						m.error = err.Error()
					} else {
						m.error = t("success_add_codex")
						m.state = codexList
						m.cursor = 0
					}
				} else {
					m.error = t("error_fill_all")
				}
			} else if m.state == editDroid {
				if m.droidFormData.ModelDisplayName != "" && m.droidFormData.Model != "" && m.droidFormData.BaseURL != "" && m.droidFormData.APIKey != "" {
					m.config.Droid[m.editIndex] = m.droidFormData
					if err := m.config.Save(); err != nil {
						m.error = err.Error()
					} else {
						m.error = t("success_update_droid")
						m.state = droidList
						m.cursor = 0
					}
				} else {
					m.error = t("error_fill_all")
				}
			} else if m.state == addDroid {
				if m.droidFormData.ModelDisplayName != "" && m.droidFormData.Model != "" && m.droidFormData.BaseURL != "" && m.droidFormData.APIKey != "" {
					if err := m.config.AddDroidConfig(m.droidFormData); err != nil {
						m.error = err.Error()
					} else {
						m.error = t("success_add_droid")
						m.state = droidList
						m.cursor = 0
					}
				} else {
					m.error = t("error_fill_all")
				}
			}
		case tea.KeyDelete:
			// 在配置列表中，Delete键删除选中的配置
			if m.state == claudeCodeList {
				if m.cursor < len(m.sortedClaudeCode) {
					originalIndex := findConfigIndex(m.config.ClaudeCode, m.sortedClaudeCode[m.cursor])
					if originalIndex == -1 {
						m.error = t("error_config_index")
						break
					}
					m.deleteIndex = originalIndex
					m.state = confirmDeleteClaudeCode
					m.cursor = 1 // 默认选择"否"
				}
			} else if m.state == codexList {
				if m.cursor < len(m.sortedCodex) {
					originalIndex := findConfigIndex(m.config.Codex, m.sortedCodex[m.cursor])
					if originalIndex == -1 {
						m.error = t("error_config_index")
						break
					}
					m.deleteIndex = originalIndex
					m.state = confirmDeleteCodex
					m.cursor = 1 // 默认选择"否"
				}
			} else if m.state == droidList {
				if m.cursor < len(m.sortedDroid) {
					originalIndex := findDroidConfigIndex(m.config.Droid, m.sortedDroid[m.cursor])
					if originalIndex == -1 {
						m.error = t("error_config_index")
						break
					}
					m.deleteIndex = originalIndex
					m.state = confirmDeleteDroid
					m.cursor = 1 // 默认选择"否"
				}
			}
		case tea.KeyBackspace, tea.KeyCtrlH:
			if m.state == addClaudeCode || m.state == addCodex || m.state == addDroid || m.state == editClaudeCode || m.state == editCodex || m.state == editDroid {
				return m.handleBackspace()
			}
		default:
			// 兜底：非 KeyRunes 情况下一般不接收字符输入。
		}
	}
	return m, nil
}

func (m model) getMaxCursor() int {
	switch m.state {
	case mainMenu:
		return 7
	case claudeCodeList:
		return len(m.sortedClaudeCode) // 返回排序后的配置数量
	case codexList:
		return len(m.sortedCodex) // 返回排序后的配置数量
	case droidList:
		return len(m.sortedDroid) // 返回排序后的配置数量
	// 操作菜单已移除，保存/取消按钮已移除
	case addClaudeCode, editClaudeCode, addDroid, editDroid:
		return DroidFieldCount - 1 // 字段 0..3
	case addCodex, editCodex:
		return CodexFieldCount - 1 // 字段 0..5
	case confirmDeleteClaudeCode, confirmDeleteCodex, confirmDeleteDroid, confirmExitAddClaudeCode, confirmExitAddCodex, confirmExitAddDroid:
		return 1 // 0=确认操作，1=取消
	default:
		return 0
	}
}

func (m model) handleSelect() (tea.Model, tea.Cmd) {
	switch m.state {
	case mainMenu:
		switch m.cursor {
		case 0:
			m.state = claudeCodeList
			m.cursor = 0
			m.sortClaudeCodeConfigs() // 进入时排序一次
		case 1:
			m.state = codexList
			m.cursor = 0
			m.sortCodexConfigs() // 进入时排序一次
		case 2:
			m.state = droidList
			m.cursor = 0
			m.sortDroidConfigs() // 进入时排序一次
		case 3:
			m.state = addClaudeCode
			m.cursor = 0
			m.formField = 0
			m.formData = ServiceConfig{}
		case 4:
			m.state = addCodex
			m.cursor = 0
			m.formField = 0
			m.formData = ServiceConfig{}
		case 5:
			m.state = addDroid
			m.cursor = 0
			m.formField = 0
			m.droidFormData = DroidConfig{}
		case 6:
			// 切换语言
			ToggleLanguage()
			m.config.Language = GetLanguage()
			m.config.Save()
			m.error = t("success_lang_switch")
		case 7:
			return m, tea.Quit
		}
	case claudeCodeList:
		if m.cursor == len(m.sortedClaudeCode) {
			// 返回主菜单时清空排序后的列表，确保下次进入时重新排序
			m.sortedClaudeCode = nil
			m.state = mainMenu
			m.cursor = 0
		} else {
			// 回车直接切换
			originalIndex := findConfigIndex(m.config.ClaudeCode, m.sortedClaudeCode[m.cursor])
			if originalIndex == -1 {
				m.error = t("error_config_index")
				break
			}
			config := &m.config.ClaudeCode[originalIndex]
			if err := m.config.SwitchClaudeCode(config); err != nil {
				m.error = fmt.Sprintf(t("error_switch_claude"), err)
			} else if err := m.config.SetActiveClaudeCode(originalIndex); err != nil {
				m.error = err.Error()
			} else {
				m.error = t("success_switch_claude")
			}
		}
		// 操作菜单已移除
	case codexList:
		if m.cursor == len(m.sortedCodex) {
			// 返回主菜单时清空排序后的列表，确保下次进入时重新排序
			m.sortedCodex = nil
			m.state = mainMenu
			m.cursor = 0
		} else {
			// 回车直接切换
			originalIndex := findConfigIndex(m.config.Codex, m.sortedCodex[m.cursor])
			if originalIndex == -1 {
				m.error = t("error_config_index")
				break
			}
			config := &m.config.Codex[originalIndex]
			if err := m.config.SwitchCodex(config); err != nil {
				m.error = fmt.Sprintf(t("error_switch_codex"), err)
			} else if err := m.config.SetActiveCodex(originalIndex); err != nil {
				m.error = err.Error()
			} else {
				m.error = t("success_switch_codex")
			}
		}
	case droidList:
		if m.cursor == len(m.sortedDroid) {
			// 返回主菜单时清空排序后的列表，确保下次进入时重新排序
			m.sortedDroid = nil
			m.state = mainMenu
			m.cursor = 0
		} else {
			// 回车直接切换
			originalIndex := findDroidConfigIndex(m.config.Droid, m.sortedDroid[m.cursor])
			if originalIndex == -1 {
				m.error = t("error_config_index")
				break
			}
			config := &m.config.Droid[originalIndex]
			if err := m.config.SwitchDroid(config); err != nil {
				m.error = fmt.Sprintf(t("error_switch_droid"), err)
			} else if err := m.config.SetActiveDroid(originalIndex); err != nil {
				m.error = err.Error()
			} else {
				m.error = t("success_switch_droid")
			}
		}
	// 操作菜单已移除
	case addClaudeCode:
		if m.cursor == 4 {
			if m.formData.Name != "" && m.formData.Provider != "" && m.formData.BaseURL != "" && m.formData.APIKey != "" {
				err := m.config.AddClaudeCodeConfig(m.formData)
				if err != nil {
					m.error = err.Error()
				} else {
					m.error = "✅ Claude Code 配置添加成功！"
					m.state = claudeCodeList
					m.cursor = 0
				}
			} else {
				m.error = "⚠️ 请填写所有字段"
			}
		} else if m.cursor == 5 {
			m.state = mainMenu
			m.cursor = 0
		}
	case addCodex:
		if m.cursor == 6 {
			if m.formData.Name != "" && m.formData.Provider != "" && m.formData.BaseURL != "" && m.formData.APIKey != "" {
				// Set default values for Codex config
				if m.formData.Model == "" {
					m.formData.Model = DefaultCodexModel
				}
				if m.formData.WireAPI == "" {
					m.formData.WireAPI = DefaultWireAPI
				}
				if m.formData.ModelReasoningEffort == "" {
					m.formData.ModelReasoningEffort = DefaultModelReasoningEffort
				}
				m.formData.EnvKey = DefaultEnvKey // Always set to CODEX_KEY

				err := m.config.AddCodexConfig(m.formData)
				if err != nil {
					m.error = err.Error()
				} else {
					m.error = "✅ Codex 配置添加成功！"
					m.state = codexList
					m.cursor = 0
				}
			} else {
				m.error = "⚠️ 请填写所有字段"
			}
		} else if m.cursor == 7 {
			m.state = mainMenu
			m.cursor = 0
		}
	case addDroid:
		if m.cursor == 4 {
			if m.droidFormData.ModelDisplayName != "" && m.droidFormData.Model != "" && m.droidFormData.BaseURL != "" && m.droidFormData.APIKey != "" {
				err := m.config.AddDroidConfig(m.droidFormData)
				if err != nil {
					m.error = err.Error()
				} else {
					m.error = "✅ Droid 配置添加成功！"
					m.state = droidList
					m.cursor = 0
				}
			} else {
				m.error = "⚠️ 请填写所有字段"
			}
		} else if m.cursor == 5 {
			m.state = mainMenu
			m.cursor = 0
		}
	case editClaudeCode:
		if m.cursor == 4 {
			if m.formData.Name != "" && m.formData.Provider != "" && m.formData.BaseURL != "" && m.formData.APIKey != "" {
				m.config.ClaudeCode[m.editIndex] = m.formData
				err := m.config.Save()
				if err != nil {
					m.error = err.Error()
				} else {
					m.error = "✅ Claude Code 配置更新成功！"
					m.state = claudeCodeList
					m.cursor = 0
				}
			} else {
				m.error = "⚠️ 请填写所有字段"
			}
		} else if m.cursor == 5 {
			m.state = claudeCodeList
			m.cursor = 0
		}
	case editCodex:
		if m.cursor == 6 {
			if m.formData.Name != "" && m.formData.Provider != "" && m.formData.BaseURL != "" && m.formData.APIKey != "" {
				// Set default values for Codex config
				if m.formData.Model == "" {
					m.formData.Model = DefaultCodexModel
				}
				if m.formData.WireAPI == "" {
					m.formData.WireAPI = DefaultWireAPI
				}
				if m.formData.ModelReasoningEffort == "" {
					m.formData.ModelReasoningEffort = DefaultModelReasoningEffort
				}
				m.formData.EnvKey = DefaultEnvKey // Always set to CODEX_KEY

				m.config.Codex[m.editIndex] = m.formData
				err := m.config.Save()
				if err != nil {
					m.error = err.Error()
				} else {
					m.error = "✅ Codex 配置更新成功！"
					m.state = codexList
					m.cursor = 0
				}
			} else {
				m.error = "⚠️ 请填写所有字段"
			}
		} else if m.cursor == 7 {
			m.state = codexList
			m.cursor = 0
		}
	case editDroid:
		if m.cursor == 4 {
			if m.droidFormData.ModelDisplayName != "" && m.droidFormData.Model != "" && m.droidFormData.BaseURL != "" && m.droidFormData.APIKey != "" {
				m.config.Droid[m.editIndex] = m.droidFormData
				err := m.config.Save()
				if err != nil {
					m.error = err.Error()
				} else {
					m.error = "✅ Droid 配置更新成功！"
					m.state = droidList
					m.cursor = 0
				}
			} else {
				m.error = "⚠️ 请填写所有字段"
			}
		} else if m.cursor == 5 {
			m.state = droidList
			m.cursor = 0
		}
	case confirmDeleteClaudeCode:
		if m.cursor == 0 {
			// 确认删除
			configName := m.config.ClaudeCode[m.deleteIndex].Name
			err := m.config.DeleteClaudeCodeConfig(m.deleteIndex)
			if err != nil {
				m.error = err.Error()
			} else {
				m.error = fmt.Sprintf(t("success_delete_claude"), configName)
				m.state = claudeCodeList
				m.cursor = 0
				// 调整光标位置
				if m.cursor >= len(m.config.ClaudeCode) && m.cursor > 0 {
					m.cursor = len(m.config.ClaudeCode) - 1
				}
			}
		} else {
			// 取消删除
			m.state = claudeCodeList
			m.cursor = m.deleteIndex
		}
	case confirmDeleteCodex:
		if m.cursor == 0 {
			// 确认删除
			configName := m.config.Codex[m.deleteIndex].Name
			err := m.config.DeleteCodexConfig(m.deleteIndex)
			if err != nil {
				m.error = err.Error()
			} else {
				m.error = fmt.Sprintf(t("success_delete_codex"), configName)
				m.state = codexList
				m.cursor = 0
				// 调整光标位置
				if m.cursor >= len(m.config.Codex) && m.cursor > 0 {
					m.cursor = len(m.config.Codex) - 1
				}
			}
		} else {
			// 取消删除
			m.state = codexList
			m.cursor = m.deleteIndex
		}
	case confirmDeleteDroid:
		if m.cursor == 0 {
			// 确认删除
			configName := m.config.Droid[m.deleteIndex].ModelDisplayName
			err := m.config.DeleteDroidConfig(m.deleteIndex)
			if err != nil {
				m.error = err.Error()
			} else {
				m.error = fmt.Sprintf(t("success_delete_droid"), configName)
				m.state = droidList
				m.cursor = 0
				// 调整光标位置
				if m.cursor >= len(m.config.Droid) && m.cursor > 0 {
					m.cursor = len(m.config.Droid) - 1
				}
			}
		} else {
			// 取消删除
			m.state = droidList
			m.cursor = m.deleteIndex
		}
	case confirmExitAddClaudeCode:
		if m.cursor == 0 {
			// 确认退出，清空表单内容
			m.state = mainMenu
			m.cursor = 0
			m.formData = ServiceConfig{}
			m.error = ""
		} else {
			// 取消退出，返回到表单
			m.state = addClaudeCode
			m.cursor = 4
		}
	case confirmExitAddCodex:
		if m.cursor == 0 {
			// 确认退出，清空表单内容
			m.state = mainMenu
			m.cursor = 0
			m.formData = ServiceConfig{}
			m.error = ""
		} else {
			// 取消退出，返回到表单
			m.state = addCodex
			m.cursor = 4
		}
	case confirmExitAddDroid:
		if m.cursor == 0 {
			// 确认退出，清空表单内容
			m.state = mainMenu
			m.cursor = 0
			m.droidFormData = DroidConfig{}
			m.error = ""
		} else {
			// 取消退出，返回到表单
			m.state = addDroid
			m.cursor = 4
		}
	}
	return m, nil
}

func (m model) handleInput(char string) (tea.Model, tea.Cmd) {
	s := sanitizeInput(char)
	if s == "" {
		return m, nil
	}

	// Handle Droid form inputs
	if m.state == addDroid || m.state == editDroid {
		switch m.formField {
		case 0:
			m.droidFormData.ModelDisplayName += s
		case 1:
			m.droidFormData.Model += s
		case 2:
			m.droidFormData.BaseURL += s
		case 3:
			m.droidFormData.APIKey += s
		}
		return m, nil
	}

	// Handle Claude Code and Codex form inputs
	switch m.formField {
	case 0:
		m.formData.Name += s
	case 1:
		m.formData.Provider += s
	case 2:
		m.formData.BaseURL += s
	case 3:
		m.formData.APIKey += s
	case 4:
		m.formData.Model += s
	case 5:
		// Wire API字段：不处理输入，使用左右键选择
	case 6:
		m.formData.EnvKey += s
	}
	return m, nil
}

func (m model) handleBackspace() (tea.Model, tea.Cmd) {
	// Handle Droid form backspace
	if m.state == addDroid || m.state == editDroid {
		switch m.formField {
		case 0:
			if len(m.droidFormData.ModelDisplayName) > 0 {
				r := []rune(m.droidFormData.ModelDisplayName)
				m.droidFormData.ModelDisplayName = string(r[:len(r)-1])
			}
		case 1:
			if len(m.droidFormData.Model) > 0 {
				r := []rune(m.droidFormData.Model)
				m.droidFormData.Model = string(r[:len(r)-1])
			}
		case 2:
			if len(m.droidFormData.BaseURL) > 0 {
				r := []rune(m.droidFormData.BaseURL)
				m.droidFormData.BaseURL = string(r[:len(r)-1])
			}
		case 3:
			if len(m.droidFormData.APIKey) > 0 {
				r := []rune(m.droidFormData.APIKey)
				m.droidFormData.APIKey = string(r[:len(r)-1])
			}
		}
		return m, nil
	}

	// Handle Claude Code and Codex form backspace
	switch m.formField {
	case 0:
		if len(m.formData.Name) > 0 {
			r := []rune(m.formData.Name)
			m.formData.Name = string(r[:len(r)-1])
		}
	case 1:
		if len(m.formData.Provider) > 0 {
			r := []rune(m.formData.Provider)
			m.formData.Provider = string(r[:len(r)-1])
		}
	case 2:
		if len(m.formData.BaseURL) > 0 {
			r := []rune(m.formData.BaseURL)
			m.formData.BaseURL = string(r[:len(r)-1])
		}
	case 3:
		if len(m.formData.APIKey) > 0 {
			r := []rune(m.formData.APIKey)
			m.formData.APIKey = string(r[:len(r)-1])
		}
	case 4:
		if len(m.formData.Model) > 0 {
			r := []rune(m.formData.Model)
			m.formData.Model = string(r[:len(r)-1])
		}
	case 5:
		// Wire API字段：不处理退格，使用左右键选择
	case 6:
		if len(m.formData.EnvKey) > 0 {
			r := []rune(m.formData.EnvKey)
			m.formData.EnvKey = string(r[:len(r)-1])
		}
	}
	return m, nil
}
