# Repository Guidelines

## Project Structure & Module Organization
- Source lives at the repo root.
- `main.go`: CLI entry; TUI startup; flags `-switch-claude` and `-switch-codex`.
- `config.go`: config persistence at `/opt/switcher/config.json`; applies settings to `~/.claude/settings.json` and `~/.codex/{auth.json,config.toml}`.
- `tui.go`: Bubble Tea + Lip Gloss TUI.
- Tooling: `Makefile` (build/install/clean), `switcher.sh` (local run), `README.md`, `CLAUDE.md`.

## Build, Test, and Development Commands
- `make build` — compile to `./switcher`.
- `sudo make install` — install binary to `/usr/bin/switcher`.
- `make clean` — remove build artifacts.
- `go run .` — run locally without installing.
- `gofmt -s -w .` and `go vet ./...` — format and basic checks.
- `go test ./...` — run tests (none present yet; see below).

## Coding Style & Naming Conventions
- Go 1.24+. Always run `gofmt -s` before committing.
- Exported identifiers: `CamelCase`. Unexported: `lowerCamel`.
- Filenames: lowercase words (underscores if needed), package `main` for CLI.
- Error handling: prefer returning errors; wrap with `fmt.Errorf("context: %w", err)`. Avoid panics in non-`main` paths.

## Testing Guidelines
- Use the standard `testing` package; name files `*_test.go` next to sources.
- Favor table-driven tests for pure helpers (e.g., `parseTomlStringArray`, `serverToTomlArray`).
- For filesystem-affecting code, use temp dirs and avoid touching real `~/.claude`, `~/.codex`, or `/opt/switcher`.
- Run `go test -v ./...` locally; keep tests deterministic and fast.

## Commit & Pull Request Guidelines
- Conventional Commits style (e.g., `feat:`, `fix:`, `docs:`, `refactor:`, `chore:`, `test:`). English or Chinese OK.
- Keep PRs focused and small; include a clear description, rationale, and screenshots/gifs if UI behavior changes.
- Link related issues; update `README.md`/`CLAUDE.md` when flags, paths, or behavior change.

## Security & Configuration Tips
- Never commit secrets or user configs. `.gitignore` already excludes `switcher` and `.claude/`.
- Keep API keys masked in logs/UI (align with existing TUI behavior).
- Installing writes to `/usr/bin` and `/opt/switcher`; use `sudo` and least privilege.

## Agent-Specific Instructions
- Do not commit the compiled `switcher` binary.
- Prefer `make build` and `go vet` before proposing changes.
- Keep edits minimal and consistent with current patterns and formatting.

