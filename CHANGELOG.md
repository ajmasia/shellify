# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.1.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

## [0.3.0] - 2026-01-08

### Added

- `SessionRepository` interface for session persistence abstraction
- `SessionService` with business logic and validation
- Session lookup methods in storage layer (GetSessionByName, SessionExists)
- Interactive TUI prompts for session management (Catppuccin theme):
  - SessionPrompt: full session creation flow with windows
  - WindowPrompt: window name and command input
  - SelectMultiplexer: tmux/zellij selector
  - SelectSession: session picker from list
  - SelectSessionFromAll: cross-project session picker
- CLI commands for session management (dual mode: interactive + direct flags):
  - `sfy session list [-p project]` - List sessions (all or by project)
  - `sfy session create [name]` - Create session with windows
  - `sfy session get [id|name]` - Get session details
  - `sfy session update [id|name]` - Update session properties
  - `sfy session delete [id|name]` - Delete session with confirmation
- Auto-generated session names as `{prefix}_{name}` format
- Cross-project session lookup by name
- `--force` flag for skipping confirmation on session delete

### Fixed

- Resolve errcheck linting issues in CLI helpers

## [0.2.0] - 2026-01-08

### Added

- `ProjectRepository` interface for project persistence abstraction
- `ProjectService` with business logic and validation
- Interactive TUI with `charmbracelet/huh` (Catppuccin theme)
- CLI commands for project management (dual mode: interactive + direct flags):
  - `sfy project list` - List all projects
  - `sfy project create [name]` - Create project (interactive if no name)
  - `sfy project get [id|name]` - Get details (interactive selection if no args)
  - `sfy project update [id|name]` - Update project (interactive field selection)
  - `sfy project delete [id|name]` - Delete project (interactive selection + confirmation)
- Input validation (empty name, duplicate name checks)
- `--force` flag for skipping confirmation on delete
- JSON output support (`--json` flag) for all project commands

## [0.1.2] - 2026-01-08

### Changed

- Reduce release targets to linux/amd64 and .deb only (save CI minutes)

## [0.1.1] - 2026-01-08

### Fixed

- Remove commit hash from changelog entries in releases

## [0.1.0] - 2026-01-08

### Added

- Project structure with Clean Architecture
- Domain entities: Project, Session, Window, Pane
- Filesystem storage with XDG compliance
- CLI skeleton with Cobra framework
- `sfy version` command with build metadata
- Makefile for common development tasks
- GitHub Actions CI/CD pipeline
- GoReleaser configuration for multi-platform releases
