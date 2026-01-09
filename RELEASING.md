# Releasing Guide

This document describes the release process for Shellify.

## Versioning

Shellify follows [Semantic Versioning](https://semver.org/):

```
MAJOR.MINOR.PATCH
```

- **MAJOR**: Breaking changes (incompatible API changes)
- **MINOR**: New features (backwards compatible)
- **PATCH**: Bug fixes (backwards compatible)

## Release Process

### 1. Pre-release Checklist

Before creating a release, ensure:

- [ ] All tests pass: `make test`
- [ ] Linting passes: `make lint`
- [ ] GUI builds successfully: `make gui-build`
- [ ] Manual testing completed
- [ ] CHANGELOG.md is updated

### 2. Update Version

Update the version in the `VERSION` file:

```bash
echo "X.Y.Z" > VERSION
```

### 3. Update Changelog

Edit `CHANGELOG.md` to move items from `[Unreleased]` to the new version:

```markdown
## [X.Y.Z] - YYYY-MM-DD

### Added
- New features...

### Changed
- Changes...

### Fixed
- Bug fixes...
```

### 4. Create Release Commit

```bash
git add VERSION CHANGELOG.md
git commit -m "chore: release vX.Y.Z"
```

### 5. Tag and Push

```bash
git tag vX.Y.Z
git push origin main
git push origin vX.Y.Z
```

### 6. Automated Release

Once the tag is pushed, GitHub Actions will automatically:

1. Run tests and linting
2. Build the GUI
3. Build binaries for all platforms
4. Create a GitHub Release with artifacts
5. Generate changelog from commits

## Release Artifacts

### Platforms

| OS | Architecture |
|----|--------------|
| Linux | amd64, arm64 |
| macOS (Darwin) | amd64, arm64 |

### Package Formats

- `tar.gz`: Universal archive
- `.deb`: Debian/Ubuntu package
- `.rpm`: Fedora/RHEL package

### Binary Naming

```
shellify_<version>_<os>_<arch>.tar.gz

Examples:
shellify_0.8.0_linux_amd64.tar.gz
shellify_0.8.0_darwin_arm64.tar.gz
```

## Hotfix Process

For critical bug fixes:

1. Create branch from tag: `git checkout -b hotfix/X.Y.Z vX.Y.Z`
2. Apply the fix
3. Update CHANGELOG.md
4. Bump PATCH version
5. Create PR to main
6. After merge, follow normal release process

## Local Testing

### Test Release Build

```bash
# Install goreleaser (if not installed)
go install github.com/goreleaser/goreleaser@latest

# Test build without releasing
goreleaser build --snapshot --clean

# Test full release process (no publish)
goreleaser release --snapshot --clean
```

### Verify Binary

```bash
./dist/sfy_linux_amd64/sfy version
```

## Changelog Format

Follow [Keep a Changelog](https://keepachangelog.com/):

### Categories

- **Added**: New features
- **Changed**: Changes in existing functionality
- **Deprecated**: Soon-to-be removed features
- **Removed**: Removed features
- **Fixed**: Bug fixes
- **Security**: Security fixes

### Example

```markdown
# Changelog

## [Unreleased]

### Added
- New feature in development

## [0.8.0] - 2024-01-15

### Added
- Visual session editor (GUI)
- Session cloning

### Fixed
- Shell completion for project flag
```

## Questions?

If you have questions about the release process, please open an issue.
