# Contributing to Shellify

Thank you for your interest in contributing to Shellify! This document provides guidelines and instructions for contributing.

## Code of Conduct

By participating in this project, you agree to maintain a respectful and inclusive environment for everyone.

## How to Contribute

### Reporting Bugs

Before submitting a bug report:

1. Check the [existing issues](https://github.com/ajmasia/shellify/issues) to avoid duplicates
2. Use the latest version of Shellify
3. Collect relevant information (OS, version, steps to reproduce)

When submitting a bug report, include:

- Clear and descriptive title
- Steps to reproduce the issue
- Expected behavior vs actual behavior
- Screenshots or terminal output if applicable
- Environment details (OS, shell, tmux/zellij version)

### Suggesting Features

Feature requests are welcome! Please:

1. Check existing issues and discussions first
2. Describe the problem your feature would solve
3. Explain your proposed solution
4. Consider if it fits the project scope

### Pull Requests

1. Fork the repository
2. Create a feature branch: `git checkout -b feature/my-feature`
3. Make your changes
4. Ensure tests pass: `make test && make lint`
5. Commit with descriptive messages (see below)
6. Push to your fork
7. Open a Pull Request

## Development Setup

### Prerequisites

- Go 1.23+
- Node.js 20+
- tmux or zellij (for testing)

### Getting Started

```bash
# Clone your fork
git clone https://github.com/YOUR_USERNAME/shellify.git
cd shellify

# Install dependencies
go mod download
cd gui && npm ci && cd ..

# Build
make build

# Run tests
make test
make lint
```

### Project Structure

```
shellify/
├── cmd/sfy/          # CLI entry point
├── internal/         # Private application code
│   ├── domain/       # Business entities
│   ├── application/  # Use cases
│   ├── infrastructure/  # External implementations
│   └── interfaces/   # CLI and HTTP handlers
├── gui/              # React frontend
└── assets/           # Static assets
```

## Coding Standards

### Go

- Follow [Effective Go](https://golang.org/doc/effective_go) guidelines
- Use `gofmt` for formatting
- Run `golangci-lint` before committing
- Write tests for new functionality
- Handle errors explicitly

### TypeScript/React

- Use TypeScript strict mode
- Follow existing code patterns
- Use CSS Modules for styling
- Write component tests

### Commits

Use [Conventional Commits](https://www.conventionalcommits.org/):

```
<type>(<scope>): <description>

[optional body]
```

**Types:**
- `feat`: New feature
- `fix`: Bug fix
- `docs`: Documentation
- `style`: Formatting (no code change)
- `refactor`: Code restructuring
- `test`: Adding tests
- `chore`: Maintenance

**Examples:**
```
feat(cli): add session clone command
fix(gui): resolve button alignment issue
docs(readme): update installation instructions
```

### Branch Naming

```
<type>/<short-description>

Examples:
feature/session-export
fix/login-crash
docs/api-reference
```

## Testing

### Go Tests

```bash
make test           # Run all tests
make test-coverage  # With coverage report
make test-race      # With race detector
```

### GUI Tests

```bash
cd gui
npm run test        # Run tests
npm run lint        # Run linter
npm run typecheck   # Type checking
```

## Pull Request Process

1. Update documentation if needed
2. Add tests for new functionality
3. Ensure all tests pass
4. Update CHANGELOG.md if applicable
5. Request review from maintainers

### PR Checklist

- [ ] Code compiles without errors
- [ ] Tests pass
- [ ] Linting passes
- [ ] Documentation updated (if needed)
- [ ] Commit messages follow conventions
- [ ] No secrets or credentials included

## Questions?

Feel free to open an issue for any questions about contributing.

Thank you for helping improve Shellify!
