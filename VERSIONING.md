# Version System Documentation

This document explains how the version system works in `clim_cli` and how to use it for releases and installations.

## Overview

The version system provides:
- Semantic versioning support
- Git commit hash tracking
- Build date and Go version information
- Multiple output formats (short, full, JSON)
- Support for `go install @latest` and `go install @main`
- Automated releases with Goreleaser

## Version Information

The version system tracks:
- **Version**: Semantic version (e.g., `v1.2.3`, `v1.2.3-alpha.1`)
- **Git Commit**: Full commit hash
- **Build Date**: UTC timestamp of build
- **Go Version**: Go version used for compilation
- **Platform**: OS/Architecture (e.g., `linux/amd64`)
- **Build Type**: `release`, `development`, or `pre-release`

## Usage

### Command Line Options

```bash
# Show short version (same as --version flag)
clim_cli --version
clim_cli -v

# Show detailed version information
clim_cli version

# Show version in JSON format
clim_cli version --json

# Show only version number
clim_cli version --short
```

### Example Output

**Short version:**
```
clim_cli v1.2.3
```

**Full version:**
```
clim_cli v1.2.3
Git Commit: a1b2c3d4e5f6...
Build Date: 2025-10-02_12:14:59
Go Version: go1.25.1
Platform: linux/amd64
Build Type: release
```

**JSON version:**
```json
{
  "version": "v1.2.3",
  "git_commit": "a1b2c3d4e5f6...",
  "build_date": "2025-10-02_12:14:59",
  "go_version": "go1.25.1",
  "platform": "linux/amd64",
  "build_type": "release"
}
```

## Installation Methods

### From GitHub Releases (Recommended)

```bash
# Install latest release
go install github.com/romaingallez/clim_cli@latest

# Install specific version
go install github.com/romaingallez/clim_cli@v1.2.3

# Install from main branch
go install github.com/romaingallez/clim_cli@main
```

### Version Detection with go install

When installing with `go install`, the version system automatically detects:

- **Tagged versions** (`@v1.2.3`): Shows the exact tag version
- **@latest**: Shows the latest tag version
- **@main**: Shows commit-based version (e.g., `dev-a1b2c3d4`)
- **Local builds**: Shows commit-based version with build info

The version detection uses Go's `runtime/debug.BuildInfo` to extract:
- Git commit hash from VCS information
- Version from module path (when installing specific tags)
- Modified status (dirty working directory)

### From Source

```bash
# Clone and build
git clone https://github.com/romaingallez/clim_cli.git
cd clim_cli
make build
```

## Building and Releasing

### Local Development Build

```bash
# Build with development version info
make build-dev

# Build with proper version injection
make build
```

### Creating Releases

#### Using Makefile (Recommended)

```bash
# Create a snapshot release (for testing)
make snapshot

# Create a full release (requires git tag)
make release
```

#### Manual Process

1. **Tag a version:**
   ```bash
   git tag -a v1.2.3 -m "Release version 1.2.3"
   git push origin v1.2.3
   ```

2. **Create release with Goreleaser:**
   ```bash
   goreleaser release --clean
   ```

### Build Flags

The build process uses these ldflags to inject version information:

```bash
-ldflags "-s -w \
  -X 'github.com/romaingallez/clim_cli/internals/version.Version=$(VERSION)' \
  -X 'github.com/romaingallez/clim_cli/internals/version.GitCommit=$(GIT_COMMIT)' \
  -X 'github.com/romaingallez/clim_cli/internals/version.BuildDate=$(BUILD_DATE)'"
```

## Version Detection

The version system automatically detects:

- **Development builds**: When `Version` is "dev" or contains "dev"
- **Pre-release builds**: When version contains "alpha", "beta", or "rc"
- **Release builds**: Clean semantic versions without pre-release identifiers

### Version Sources

The version information comes from different sources depending on how the binary was built:

1. **Goreleaser builds** (releases): Uses ldflags to inject exact version info
2. **Makefile builds**: Uses git describe to get version from tags
3. **go install builds**: Uses runtime/debug.BuildInfo to detect version at runtime
4. **Local development**: Falls back to commit-based versioning

### Example Version Formats

- `v1.2.3` - Clean release version
- `v1.2.3-alpha.1` - Pre-release version
- `v1.2.3-8-g2e55078` - Tag with additional commits (git describe)
- `v1.2.3-8-g2e55078-dirty` - Tag with uncommitted changes
- `dev-a1b2c3d4` - Development build from commit
- `dev-a1b2c3d4-dirty` - Development build with uncommitted changes

## Goreleaser Configuration

The `.goreleaser.yml` file configures:

- **Multi-platform builds**: Linux, macOS, Windows (amd64, arm64)
- **GitHub releases**: Automatic release creation
- **Package managers**: Homebrew, Scoop, Snapcraft support
- **Checksums**: SHA256 checksums for all binaries
- **Archives**: Tar.gz archives with license and README

## Makefile Targets

| Target | Description |
|--------|-------------|
| `make build` | Build binary with version injection |
| `make build-dev` | Build for development (no version injection) |
| `make clean` | Clean build artifacts |
| `make cross-build` | Cross-compile for multiple platforms |
| `make release` | Create GitHub release with Goreleaser |
| `make snapshot` | Create snapshot release for testing |
| `make version` | Show current version information |

## Semantic Versioning

Follow [Semantic Versioning](https://semver.org/) principles:

- **MAJOR**: Incompatible API changes
- **MINOR**: Backward-compatible functionality additions
- **PATCH**: Backward-compatible bug fixes

Examples:
- `v1.0.0` - First stable release
- `v1.1.0` - New features, backward compatible
- `v1.1.1` - Bug fixes only
- `v2.0.0` - Breaking changes
- `v1.2.0-alpha.1` - Pre-release version

## Troubleshooting

### Version shows "dev"
- Ensure you're using `make build` or have proper ldflags
- Check that git is available and repository is initialized

### Goreleaser fails
- Ensure you have a valid GitHub token
- Check that the repository exists and you have push access
- Verify the `.goreleaser.yml` configuration

### Go install fails
- Ensure the repository is public or you have access
- Check that the version/tag exists
- Verify the module path in `go.mod` matches the repository

## Best Practices

1. **Always tag releases**: Use semantic versioning tags
2. **Test before release**: Use `make snapshot` to test releases
3. **Document changes**: Update CHANGELOG.md for each release
4. **Version consistency**: Ensure version in code matches git tag
5. **Build verification**: Test builds on multiple platforms
