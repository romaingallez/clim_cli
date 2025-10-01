# Agent Guidelines for clim_cli

## Build Commands
- Build to tmp folder: `go build -o ./tmp/clim_cli .`
- Standard build: `go build .`
- Cross-platform build: `goreleaser release --clean`

## Test Commands
- No test files present in codebase

## Lint Commands
- No linting configuration found

## Code Style Guidelines

### Imports
- Group imports: standard library, third-party, local packages
- Use blank lines between import groups

### Naming Conventions
- Functions: PascalCase for exported, camelCase for unexported
- Structs: PascalCase for exported fields
- Variables: camelCase
- Constants: PascalCase

### Error Handling
- Return nil on HTTP/network errors
- Use fmt.Printf for error output in CLI commands
- Use log package for structured logging with prefix

### Formatting
- Use gofmt for consistent formatting
- No inline comments in implementation code
- Copyright headers on all .go files

### Architecture
- CLI commands in `cmd/` package using Cobra
- Business logic in `internals/` packages
- API layer in `internals/api/`
- Command handlers in `internals/commands/`

## Cursor Rules
- Build binary to `./tmp` folder (see `.cursor/rules/build-tmp.mdc`)

## Dependencies
- Go 1.23.0 minimum
- Uses Cobra for CLI, Viper for config (though config not currently used)