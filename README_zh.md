# 🔄 Switcher

<div align="center">

![Go 版本](https://img.shields.io/badge/Go-1.24+-00ADD8?style=for-the-badge&logo=go)
![开源协议](https://img.shields.io/badge/License-Apache%202.0-green?style=for-the-badge)
![支持平台](https://img.shields.io/badge/Platform-Linux%20%7C%20macOS%20%7C%20Windows-lightgrey?style=for-the-badge)

[English](README.md) | **中文**

*一款精美的基于 TUI 的命令行工具，用于管理和切换 Claude Code、Codex 与 Droid 配置*

[![演示](https://img.shields.io/badge/演示-🎬-ff69b4?style=for-the-badge)](#-演示)
[![安装](https://img.shields.io/badge/安装-📦-4285f4?style=for-the-badge)](#-安装)
[![使用](https://img.shields.io/badge/使用-🚀-f39c12?style=for-the-badge)](#-使用方法)

</div>

## ✨ 功能特性

- 🎨 **精美 TUI 界面** - 基于 [Bubble Tea](https://github.com/charmbracelet/bubbletea) 构建的优雅终端体验
- 🌐 **双语支持** - 按 'L' 键即可在中英文之间无缝切换
- ⚡ **快速切换** - 即时切换不同的 API 配置
- 🔒 **安全管理** - API 密钥在显示时会被遮蔽，确保安全
- 📝 **配置 CRUD** - 轻松添加、编辑、删除和管理配置
- 🎯 **三服务支持** - 同时管理 Claude Code、Codex 和 Droid 配置
- 💻 **命令行模式** - 支持非交互式命令行切换
- 📂 **自动导入** - 首次运行时自动导入现有配置
- 🔄 **实时更新** - 更改立即应用到您的配置文件

## 🎬 演示

### 主菜单
<div align="center">
<img src="docs/screenshots/main-menu.png" alt="主菜单" width="600"/>

*支持双语的主菜单 - 轻松切换服务或更改语言*
</div>

### 配置列表
<div align="center">
<img src="docs/screenshots/config-list.png" alt="配置列表" width="600"/>

*浏览和管理您的配置，带有活动状态指示器*
</div>

### 编辑配置
<div align="center">
<img src="docs/screenshots/edit-config.png" alt="编辑配置" width="600"/>

*直观的表单界面，支持逐字段编辑和 API 密钥遮蔽*
</div>

### 快速命令

```bash
# 启动交互式 TUI
switcher

# 或通过 CLI 直接切换
switcher -switch-claude "OpenAI GPT-4"
switcher -switch-codex "Anthropic Claude"
switcher -switch-droid "Droid Model"
```

## 📦 安装

### Linux / macOS

**从源码安装**

```bash
# 克隆仓库
git clone https://github.com/bingfengfeifei/switcher.git
cd switcher

# 构建并安装
make build
sudo make install
```

**使用 Go 安装**

```bash
# 直接安装
go install github.com/bingfengfeifei/switcher@latest

# 或克隆后构建
git clone https://github.com/bingfengfeifei/switcher.git
cd switcher
go build -o switcher .
```

### Windows

**使用 PowerShell**

```powershell
# 克隆仓库
git clone https://github.com/bingfengfeifei/switcher.git
cd switcher

# 构建
.\build.ps1

# 构建并安装到系统
.\build.ps1 -Install
```

**使用 Go**

```powershell
# 直接安装
go install github.com/bingfengfeifei/switcher@latest

# 或克隆后构建
git clone https://github.com/bingfengfeifei/switcher.git
cd switcher
go build -o switcher.exe .
```

## 🚀 使用方法

### 交互模式（默认）

```bash
switcher
```

使用以下按键导航精美的 TUI 界面：
- **↑/↓** 或 **j/k** - 导航菜单项
- **Enter** - 选择/确认操作
- **Tab** - 在表单字段间切换
- **C** - 复制选中的配置
- **Esc** - 返回/退出
- **q** - 退出应用程序

### 命令行模式

```bash
# 切换 Claude Code 配置
switcher -switch-claude "配置名称"

# 切换 Codex 配置
switcher -switch-codex "配置名称"

# 切换 Droid 配置
switcher -switch-droid "配置名称"
```

## 📁 文件位置

### Linux

| 文件 | 位置 | 用途 |
|------|----------|---------|
| **可执行文件** | `/usr/bin/switcher` | 系统可执行文件 |
| **应用配置** | `~/.config/switcher/config.json` | 存储的配置 |
| **Claude Code** | `~/.claude/settings.json` | Claude Code 设置 |
| **Codex 认证** | `~/.codex/auth.json` | Codex 身份验证 |
| **Codex 配置** | `~/.codex/config.toml` | Codex 配置 |
| **Droid 配置** | `~/.factory/config.json` | Droid 配置 |

### macOS

| 文件 | 位置 | 用途 |
|------|----------|---------|
| **可执行文件** | `/usr/bin/switcher` | 系统可执行文件 |
| **应用配置** | `~/Library/Application Support/switcher/config.json` | 存储的配置 |
| **Claude Code** | `~/.claude/settings.json` | Claude Code 设置 |
| **Codex 认证** | `~/.codex/auth.json` | Codex 身份验证 |
| **Codex 配置** | `~/.codex/config.toml` | Codex 配置 |
| **Droid 配置** | `~/.factory/config.json` | Droid 配置 |

### Windows

| 文件 | 位置 | 用途 |
|------|----------|---------|
| **可执行文件** | `%LOCALAPPDATA%\Programs\switcher\switcher.exe` | 系统可执行文件 |
| **应用配置** | `%APPDATA%\switcher\config.json` | 存储的配置 |
| **Claude Code** | `%USERPROFILE%\.claude\settings.json` | Claude Code 设置 |
| **Codex 认证** | `%USERPROFILE%\.codex\auth.json` | Codex 身份验证 |
| **Codex 配置** | `%USERPROFILE%\.codex\config.toml` | Codex 配置 |
| **Droid 配置** | `%USERPROFILE%\.factory\config.json` | Droid 配置 |

## 🛠️ 配置结构

### Claude Code 和 Droid

```json
{
  "name": "我的 API 配置",
  "provider": "openai",
  "base_url": "https://api.openai.com/v1",
  "api_key": "sk-..."
}
```

### Codex（支持认证方式选择）

```json
{
  "name": "我的 Codex 配置",
  "provider": "openai",
  "base_url": "https://api.openai.com/v1",
  "api_key": "sk-...",
  "model": "gpt-5.1-codex",
  "wire_api": "responses",
  "auth_method": "auth.json",
  "model_reasoning_effort": "medium"
}
```

**Codex 认证方式：**
- **`auth.json`**（默认）- 同步 `~/.codex/auth.json`，并在 `~/.codex/config.toml` 的 provider 中设置 `experimental_bearer_token`，使切换后的请求使用所选 key
- **`env`** - 使用环境变量 `CODEX_KEY`（自动设置到 shell 配置文件中）

## 🎯 支持的提供商

- **OpenAI** - GPT 模型和 API
- **Anthropic** - Claude 模型
- **自定义** - 任何兼容 OpenAI 的 API 端点

## 🏗️ 架构

```
switcher/
├── main.go            # 入口点和 CLI 参数
├── tui/
│   ├── config.go      # 配置管理
│   ├── platform.go    # 跨平台路径抽象
│   ├── shell.go       # Shell 环境变量管理
│   ├── controller.go  # 事件处理和状态机
│   ├── menu.go        # 状态定义和视图路由
│   ├── init.go        # 模型初始化
│   ├── style.go       # 样式和UI组件
│   ├── util.go        # 工具函数
│   ├── claudecode.go  # Claude Code 服务组件
│   ├── codex.go       # Codex 服务组件
│   └── droid.go       # Droid 服务组件
├── Makefile           # 构建自动化 (Linux/macOS)
├── build.ps1          # 构建脚本 (Windows)
└── README.md          # 本文件
```

### 核心组件

- **配置引擎** (`tui/config.go`) - 处理配置的加载、保存和应用，支持 Claude Code、Codex 和 Droid
- **平台抽象层** (`tui/platform.go`) - 跨平台路径管理，支持 Linux、macOS 和 Windows
- **Shell 管理器** (`tui/shell.go`) - 跨平台环境变量管理（bash/zsh/fish/PowerShell）
- **TUI 控制器** (`tui/controller.go`) - 中央事件处理、状态转换和键盘输入处理
- **TUI 菜单系统** (`tui/menu.go`) - 状态管理、模型结构和视图路由
- **服务组件** (`tui/*code*.go`) - 各服务的列表视图和专用逻辑
- **样式系统** (`tui/style.go`) - 使用 Lipgloss 的样式库
- **CLI 接口** (`main.go`) - 命令行切换功能和 TUI 初始化

## 🔧 开发

### 环境要求

- Go 1.24.0 或更高版本
- 支持 Linux、macOS 和 Windows 操作系统
- Make（可选，用于 Linux/macOS 构建自动化）
- PowerShell（Windows 构建）

### 构建

**Linux / macOS**

```bash
# 构建二进制文件
make build

# 构建所有平台版本
make build-all

# 安装到系统
sudo make install

# 清理构建产物
make clean
```

**Windows**

```powershell
# 构建二进制文件
.\build.ps1

# 构建并安装
.\build.ps1 -Install
```

### 本地运行

```bash
# 从源码运行
go run .

# 或构建后运行
go build -o switcher .
./switcher
```

## 🎨 自定义

TUI 支持高级用户的键盘快捷键：

- **Vim 风格导航** 使用 `j` 和 `k`
- **快速操作** 单按键操作
- **表单导航** 使用 Tab 在字段间切换
- **转义序列** 直观的导航体验

## 🔒 安全性

- API 密钥在 TUI 显示中会被**遮蔽**（`sk-****...`）
- 配置文件具有**适当的权限**
- 命令输出中不会记录或暴露 API 密钥

## 🤝 贡献

1. Fork 本仓库
2. 创建您的功能分支 (`git checkout -b feature/amazing-feature`)
3. 提交您的更改 (`git commit -m '添加某个很棒的功能'`)
4. 推送到分支 (`git push origin feature/amazing-feature`)
5. 创建 Pull Request

## 📄 开源协议

本项目基于 Apache 2.0 协议开源 - 详情请参阅 [LICENSE](LICENSE) 文件。

## 🙏 致谢

- [Bubble Tea](https://github.com/charmbracelet/bubbletea) - 超赞的 TUI 框架
- [Lipgloss](https://github.com/charmbracelet/lipgloss) - 精美的样式库
- Go 社区的优秀生态系统

## 📞 支持

如果您遇到任何问题或有功能请求：

- 🐛 [报告错误](https://github.com/bingfengfeifei/switcher/issues/new?template=bug_report.md)
- 💡 [请求功能](https://github.com/bingfengfeifei/switcher/issues/new?template=feature_request.md)
- 💬 [开始讨论](https://github.com/bingfengfeifei/switcher/discussions)

---

<div align="center">

**⭐ 如果这个项目对您有帮助，请给我们一个 Star！**

由开源社区用 ❤️ 制作

</div>
