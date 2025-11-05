# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

This is a Go-based CLI tool called "switcher" that provides a TUI (Terminal User Interface) for managing and switching between different configurations for Claude Code, Codex, and Droid. The application stores configurations in `/opt/switcher/config.json` and applies them to the respective configuration files in the user's home directory.

## Architecture

### Core Components

- **main.go**: Entry point with command-line flag support for non-interactive switching and TUI initialization
- **tui/config.go**: Configuration management, persistence, and file operations for Claude Code (`.claude/settings.json`), Codex (`.codex/auth.json` and `.codex/config.toml`), and Droid (`.factory/config.json`)
- **tui/controller.go**: Central event handling, state transitions, and keyboard input processing
- **tui/menu.go**: State definitions, model struct, view routing, and form rendering
- **tui/*.go**: Service-specific components (claudecode.go, codex.go, droid.go) with list views and specialized logic
- **tui/style.go**: Lipgloss styling, UI component rendering, and provider badges
- **tui/util.go**: Utility functions for API key masking and input sanitization
- **tui/init.go**: Model initialization and Bubble Tea interface implementation

### Key Data Structures

- `ServiceConfig`: Configuration with Name, Provider, BaseURL, APIKey, Model, WireAPI, EnvKey
- `DroidConfig`: Specialized configuration for Droid service with ModelDisplayName, Model, BaseURL, APIKey, Provider
- `Config`: Main container holding arrays of configurations and active indices for all three services
- `model`: TUI state machine managing navigation, forms, cursor positions, and user interactions across 16 distinct states

### State Management

The TUI uses a sophisticated state machine with these key states:
- **Navigation**: `mainMenu`, `claudeCodeList`, `codexList`, `droidList`
- **Forms**: `addClaudeCode`, `addCodex`, `addDroid`, `editClaudeCode`, `editCodex`, `editDroid`
- **Confirmations**: `confirmDeleteClaudeCode`, `confirmDeleteCodex`, `confirmDeleteDroid`, `confirmExitAdd*`

Navigation flow: Main menu → Service lists → Individual actions with hierarchical Esc navigation back to main menu.

## Development Commands

### Building and Running

```bash
# Build the application
make build
# or
go build -o switcher .

# Run locally
./switcher

# Install to system
sudo make install

# Clean build artifacts
make clean
```

### Command-line Interface

The tool supports both interactive TUI mode and non-interactive command-line switching:

```bash
# Interactive mode (default)
switcher

# Non-interactive switching
switcher -switch-claude "Configuration Name"
switcher -switch-codex "Configuration Name"
switcher -switch-droid "Configuration Name"
```

CLI error codes: 2=config not found, 3=switch failed, 4=set active failed

## File Locations

- **Source code**: `/root/git/switcher/`
- **System binary**: `/usr/bin/switcher`
- **Application config**: `/opt/switcher/config.json`
- **Claude Code settings**: `~/.claude/settings.json`
- **Codex auth**: `~/.codex/auth.json`
- **Codex config**: `~/.codex/config.toml`
- **Droid config**: `~/.factory/config.json`

## Key Features

- **TUI Navigation**: Arrow keys, j/k vim-style navigation, Tab for field switching, Ctrl+S for save, V for view toggle, A for quick add
- **Configuration Management**: Add, edit, delete, and switch between configurations for three services
- **Security**: API keys are masked in display unless being edited
- **Import**: Automatically imports existing configurations on first run
- **Validation**: Ensures configuration consistency between stored and applied settings
- **Sorted Display**: Active configurations appear first, then alphabetical order
- **Confirmation Dialogs**: Consistent pattern for destructive operations

## Dependencies

- `github.com/charmbracelet/bubbletea`: TUI framework
- `github.com/charmbracelet/lipgloss`: Styling and layout
- Go 1.24.0 or higher

## Testing

No formal test suite is currently present in the codebase. Testing is done manually through the TUI interface and command-line operations.

## Key Architectural Patterns

- **Service-Specific Segregation**: Each service has dedicated files for list views and specialized logic
- **Form State Management**: Separate form data structures with field-by-field editing and Tab navigation
- **Configuration Migration**: Automatic import and migration of existing configurations
- **Error Handling**: Centralized error display with styled messages throughout the TUI
- **State Transitions**: Large switch statement in `Update()` method handles all 16 state transitions
- **Provider Styling**: Color-coded provider badges and consistent styling across all views