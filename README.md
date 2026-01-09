<p align="center">
  <img src="assets/logo.svg" alt="Shellify Logo" width="128" height="128">
</p>

<h1 align="center">Shellify</h1>

<p align="center">
  <strong>Terminal multiplexer session manager for tmux and zellij</strong>
</p>

<p align="center">
  <a href="https://github.com/ajmasia/shellify/releases/latest">
    <img src="https://img.shields.io/github/v/release/ajmasia/shellify?style=flat-square&color=blue" alt="Latest Release">
  </a>
  <a href="https://github.com/ajmasia/shellify/blob/main/LICENSE">
    <img src="https://img.shields.io/github/license/ajmasia/shellify?style=flat-square&color=green" alt="License">
  </a>
  <a href="https://github.com/ajmasia/shellify/stargazers">
    <img src="https://img.shields.io/github/stars/ajmasia/shellify?style=flat-square&color=yellow" alt="Stars">
  </a>
  <img src="https://img.shields.io/badge/Go-1.23+-00ADD8?style=flat-square&logo=go" alt="Go Version">
  <img src="https://img.shields.io/badge/React-19-61DAFB?style=flat-square&logo=react" alt="React Version">
  <img src="https://img.shields.io/badge/Platform-Linux%20%7C%20macOS-lightgrey?style=flat-square" alt="Platform">
</p>

---

## About

Shellify is a session manager for terminal multiplexers (tmux and zellij). It allows you to define, save, and launch complex terminal layouts with multiple windows and panes, organized by projects.

**Key Features:**

- Organize sessions by projects
- Support for tmux and zellij
- Visual session editor (GUI)
- Interactive CLI with shell completions
- Launch, attach, and stop sessions
- Clone and duplicate sessions
- JSON-based portable configuration

## Installation

### From Releases (Recommended)

Download the latest release for your platform from [GitHub Releases](https://github.com/ajmasia/shellify/releases/latest).

#### Debian/Ubuntu (.deb)

```bash
# Download and install (amd64)
curl -LO https://github.com/ajmasia/shellify/releases/latest/download/shellify_amd64.deb
sudo dpkg -i shellify_amd64.deb

# Or for arm64
curl -LO https://github.com/ajmasia/shellify/releases/latest/download/shellify_arm64.deb
sudo dpkg -i shellify_arm64.deb
```

#### Manual Installation (tar.gz)

```bash
# Linux (amd64)
curl -LO https://github.com/ajmasia/shellify/releases/latest/download/shellify_linux_amd64.tar.gz
tar -xzf shellify_linux_amd64.tar.gz
sudo mv sfy /usr/local/bin/

# Linux (arm64)
curl -LO https://github.com/ajmasia/shellify/releases/latest/download/shellify_linux_arm64.tar.gz
tar -xzf shellify_linux_arm64.tar.gz
sudo mv sfy /usr/local/bin/

# macOS (Intel)
curl -LO https://github.com/ajmasia/shellify/releases/latest/download/shellify_darwin_amd64.tar.gz
tar -xzf shellify_darwin_amd64.tar.gz
sudo mv sfy /usr/local/bin/

# macOS (Apple Silicon)
curl -LO https://github.com/ajmasia/shellify/releases/latest/download/shellify_darwin_arm64.tar.gz
tar -xzf shellify_darwin_arm64.tar.gz
sudo mv sfy /usr/local/bin/
```

### From Source

```bash
git clone https://github.com/ajmasia/shellify.git
cd shellify
make build
sudo mv bin/sfy /usr/local/bin/
```

### Uninstall

```bash
# If installed via .deb
sudo apt remove shellify

# If installed manually
sudo rm /usr/local/bin/sfy

# Optional: remove configuration
rm -rf ~/.config/shellify
```

## Quick Start

```bash
# Create a project
sfy project create my-project

# Create a session
sfy session create dev -p my-project

# Edit session layout (GUI)
sfy session edit dev --gui

# Launch session
sfy session launch dev
```

## Usage

### Projects

```bash
sfy project list                    # List all projects
sfy project create [name]           # Create project (interactive if no name)
sfy project get <name>              # Show project details
sfy project update <name>           # Update project
sfy project delete <name>           # Delete project
```

### Sessions

```bash
sfy session list [-p project]       # List sessions
sfy session create [-p project]     # Create session
sfy session get <name>              # Show session details
sfy session edit <name> --gui       # Visual editor
sfy session clone <name>            # Clone session
sfy session delete <name>           # Delete session

sfy session launch <name>           # Launch session
sfy session attach <name>           # Attach to running session
sfy session stop <name>             # Stop session
sfy session status <name>           # Check session status
```

### Server & GUI

```bash
sfy server                          # Start HTTP server
sfy server stop                     # Stop server
sfy server status                   # Check server status
```

### Shell Completions

```bash
# Bash
sfy completion bash > /etc/bash_completion.d/sfy

# Zsh
sfy completion zsh > "${fpath[1]}/_sfy"

# Fish
sfy completion fish > ~/.config/fish/completions/sfy.fish
```

For more details, run `sfy --help` or `sfy <command> --help`.

## Configuration

Shellify follows the [XDG Base Directory Specification](https://specifications.freedesktop.org/basedir-spec/basedir-spec-latest.html).

### Data Location

```
~/.config/shellify/
├── config.json                    # Application settings
└── projects/
    └── <project-id>/
        ├── project.json           # Project metadata
        └── sessions/
            └── <session-id>.json  # Session configuration
```

### Configuration File

The `config.json` file stores application settings:

```json
{
  "defaultMultiplexer": "tmux",
  "defaultShell": "/bin/bash"
}
```

## Development

### Prerequisites

- Go 1.23+
- Node.js 20+
- tmux or zellij (for testing)

### Setup

```bash
# Clone repository
git clone https://github.com/ajmasia/shellify.git
cd shellify

# Install dependencies
go mod download
cd gui && npm ci && cd ..

# Build
make build           # CLI only
make build-with-gui  # CLI with embedded GUI
```

### Development Mode

```bash
# Terminal 1: Start API server
make run ARGS="server --static gui/dist"

# Terminal 2: Start GUI dev server
make gui-dev
```

### Testing

```bash
make test       # Run Go tests
make lint       # Run linter
make gui-check  # Run GUI lint and typecheck
```

## Contributing

Contributions are welcome! Please read our [Contributing Guide](CONTRIBUTING.md) for details on our code of conduct and the process for submitting pull requests.

## License

This project is licensed under the GNU General Public License v3.0 - see the [LICENSE](LICENSE) file for details.

---

<p align="center">
  Made with <span style="color: #e25555;">&#9829;</span> using Go and React
</p>
