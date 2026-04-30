# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.1.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

## [0.9.7] - 2026-04-30

### Changed

- Save button in session editor is now disabled when there are no unsaved changes

## [0.9.6] - 2026-04-30

### Fixed

- tmux prefix matching causing sessions with shared name prefixes to interfere with each other (status reporting, launch, kill)
- Special characters in session names (e.g. parentheses from cloned sessions) stripped to avoid tmux target syntax issues

## [0.9.5] - 2026-01-16

### Fixed

- Correct npmDepsHash for Nix GUI build

## [0.9.4] - 2026-01-16

### Added

- Nix package build with embedded GUI and shell completions
- Home Manager module (`programs.shellify.enable`)
- Overlay for NixOS configurations
- NixOS installation documentation in README

### Changed

- Update devShell to use Node.js 22 (required by GUI)
- Update prerequisites documentation (Go 1.24+, Node.js 22+)

## [0.9.3] - 2026-01-16

### Fixed

- Sessions created via CLI with multiple panes now render correctly (was showing only one pane)

### Added

- Nix flake for reproducible development environment

### Changed

- Update golangci-lint to v1.62.2 for Go 1.24+ compatibility

## [0.9.2] - 2026-01-16

### Fixed

- Session launch fails with "no such file or directory" on NixOS and other systems where bash is not at /bin/bash

## [0.9.1] - 2026-01-09

### Changed

- Update README badges to markdown link format with CI/Release status

## [0.9.0] - 2026-01-09

### Added

- Shell completion for project and session names
  - Dynamic completion for `-p` flag (project selection)
  - Session name completion filtered by project when `-p` is provided
- Initial project documentation
  - Comprehensive README with badges, installation, usage, and configuration
  - CONTRIBUTING.md with contribution guidelines
  - RELEASING.md with release process documentation
  - Project logo in assets/

### Fixed

- `-p` flag now properly filters sessions by project in interactive mode
- `--gui` flag shows server messages only when actually starting the server
- Window tabs UX: tabs now fit text width and expand on hover to reveal delete icon

## [0.8.0] - 2026-01-09

### Added

- Visual Session Editor for designing pane layouts:
  - Interactive pane grid with drag-to-resize
  - Split panes horizontally/vertically
  - Window management (add, remove, rename, reorder tabs)
  - Pane configuration panel (name, command, working directory)
  - Session settings modal (name, description, multiplexer, environment)
  - Undo/Redo functionality with keyboard shortcuts (Ctrl+Z, Ctrl+Y)
  - Unsaved changes detection with confirmation modal
  - Maximum 10 windows per session limit
- CLI-GUI integration:
  - `sfy session edit [id|name]` - Edit session JSON in $EDITOR with warning
  - `sfy session edit --gui` - Open visual editor in browser
  - `sfy session create --gui` - Create session with GUI editor
  - "Create + Edit" button in CreateSessionModal
  - Auto-start server when using `--gui` flag
- Editor icons: NewWindow, SplitVertical, SplitHorizontal, Undo, Redo, Save
- Favicon theme detection (dark/light system theme)
- Project name display in editor header (project: session)
- Alt+N keyboard shortcut to create session from session list
- Multiplexer chip in session cards

### Fixed

- Preserve nested pane tree structure when saving sessions via API
- Expand tilde (~) to $HOME in working directory for generators
- Use `--new-session-with-layout` flag for zellij session launch
- Generate bash wrapper script for zellij sessions

### Changed

- Improved editor UI with better spacing and text sizes
- Separated window tabs bar from terminal title bar
- Unified header styles across components

## [0.7.1] - 2026-01-08

### Fixed

- Handle `file.Close` error in embedded GUI handler (errcheck linter)

### Changed

- Restructure CI pipeline to integrate GUI build into Go workflow
  - Go Lint and Go Test now depend on GUI Build
  - Go Build uses `-tags embed_gui` with embedded GUI artifact

## [0.7.0] - 2026-01-08

### Added

- React GUI application with Vite and TypeScript (`gui/`)
  - Catppuccin Mocha dark theme with CSS custom properties
  - MobX store for UI state management
  - React Router for navigation
  - Sonner for toast notifications
- Domain types mirroring Go entities:
  - Project, Session, Window, Pane types
  - Multiplexer type constants
- API client module:
  - Typed HTTP client wrapper
  - Full CRUD operations for projects and sessions
  - Health check endpoint
- Component library:
  - Button, IconButton with variants and sizes
  - Tooltip with Radix UI primitives
  - Chip for status/category display
  - FAB (Floating Action Button) circular component
  - Icon components (Plus, Delete, Edit, Play, Stop, Folder, Back, Close, Settings)
- Modal system with React Portals:
  - BaseModal with backdrop blur and dark theme
  - ConfirmModal for delete confirmations
  - CreateProjectModal for project creation
  - CreateSessionModal for session creation
  - ModalContext and useModal hook
- Layout components:
  - MainLayout with header and watermark
  - Header with navigation and settings
- Feature modules:
  - ProjectList with grid view and project cards
  - SessionList with grid view and session cards
  - useProjects and useSessions hooks for data fetching
- Conditional GUI embedding in Go binary:
  - `gui/embed.go` with `//go:build embed_gui` tag
  - `gui/embed_stub.go` stub for non-embedded builds
  - `setupEmbeddedHandler` in HTTP server
- Build automation:
  - Makefile targets: gui-install, gui-dev, gui-build, gui-lint, gui-check, build-with-gui
  - CI jobs: frontend-deps, frontend-lint, frontend-build
  - Release pipeline with GUI build before GoReleaser
  - embed_gui build tag in .goreleaser.yaml

## [0.6.0] - 2026-01-08

### Added

- HTTP API server using Chi router (`internal/interfaces/http/`)
  - Middleware stack: RequestID, RealIP, Logger, Recoverer, Timeout, CORS
  - Standard JSON response format: `{success, data, error}`
  - Domain error to HTTP status mapping
- Project API endpoints:
  - `GET /api/projects` - List all projects with session counts
  - `POST /api/projects` - Create project
  - `GET /api/projects/{id}` - Get project
  - `PUT /api/projects/{id}` - Update project
  - `DELETE /api/projects/{id}` - Delete project (with `?force=true` option)
  - `GET /api/projects/{id}/sessions` - List project sessions
  - `POST /api/projects/{id}/backup` - Backup project with sessions
  - `POST /api/projects/{id}/restore` - Restore sessions
- Session API endpoints:
  - `GET /api/sessions` - List all sessions
  - `POST /api/sessions` - Create session
  - `GET /api/sessions/running` - List running sessions
  - `GET /api/sessions/{id}` - Get session
  - `PUT /api/sessions/{id}` - Update session
  - `DELETE /api/sessions/{id}` - Delete session
  - `POST /api/sessions/{id}/clone` - Clone session
  - `POST /api/sessions/{id}/launch` - Launch session (detached)
  - `POST /api/sessions/{id}/attach` - Attach to running session
  - `POST /api/sessions/{id}/stop` - Stop running session
  - `GET /api/sessions/{id}/status` - Check session status
- Settings API endpoints:
  - `GET /api/settings` - Get application settings
  - `PUT /api/settings` - Update settings
- Utility endpoints:
  - `GET /api/health` - Health check
  - `GET /api/docs` - API documentation (HTML)
- CLI server commands:
  - `sfy server` - Start HTTP server (foreground)
  - `sfy server -d` - Start as daemon (background)
  - `sfy server stop` - Stop running daemon
  - `sfy server status` - Check daemon status
  - Flags: `--port`, `--host`, `--static`, `--open`
- `FindSessionByID` method in SessionService for cross-project lookup
- Static file serving with SPA fallback support

## [0.5.0] - 2026-01-08

### Added

- Multiplexer detection infrastructure (`internal/infrastructure/multiplexer/`)
  - `CheckMultiplexers()`: detect tmux/zellij availability
  - `GetInstallInstructions()`: user-friendly installation help
- Terminal emulator detection (alacritty, kitty, wezterm, gnome-terminal, konsole, xfce4-terminal, xterm)
- `Launcher` for session execution and management:
  - Launch sessions (replaces current process for CLI)
  - Launch detached sessions (new terminal window for API)
  - Attach to running sessions
  - Stop running sessions
  - Check session running/attached status
- `LauncherService` in application layer for session lifecycle management
- CLI commands for session control:
  - `sfy session launch [id|name]` - Launch or attach to session
  - `sfy session attach [id|name]` - Attach to running session
  - `sfy session stop [id|name]` - Stop running session
  - `sfy session status [id|name]` - Show session running status
  - `sfy session clone [id|name]` - Clone existing session
  - `sfy session kill` - Multi-select and kill running sessions
- Status column in `sfy session list` showing running/stopped state
- `SelectMultipleSessions` TUI prompt for batch operations
- `CloneSession` method in SessionService

## [0.4.0] - 2026-01-08

### Added

- Generator interface for multiplexer script/layout generation
- TmuxGenerator: generates executable bash scripts for tmux sessions
  - Recursive pane tree traversal with correct split percentages
  - Environment variables and pre/post commands support
  - Terminal size detection and pane-base-index handling
- ZellijGenerator: generates KDL layout files for zellij sessions
  - Direction mapping (horizontal/vertical swap for zellij)
  - Tab focus and default tab template support
- GeneratorService in application layer for session generation
- CLI command `sfy session generate [id|name]`:
  - Output to stdout (default) or file (`-o` flag)
  - Make tmux scripts executable automatically
  - Show usage hints after file generation
- Complete interactive session creation workflow:
  - Environment variables collection with key-value loop
  - Pre-session and post-session commands collection
  - Default window selection for multiple windows
  - Immediate feedback after each user decision
  - Full session summary before confirmation
  - Final confirmation prompt before saving
- Default pane names (`pane-N`) when not provided by user

### Changed

- Extended `CreateSessionInput` with Environment, PreCommands, PostCommands, DefaultWindowID fields

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
