# AGENTS.md - Developer Guidelines

## Project Overview

This is a Laravel Dev Tools CLI application written in Go. It provides an interactive installer for PHP, MariaDB, Node.js, Composer, and Laravel Valet.

---

## Build, Run & Development Commands

### Build
```bash
# Build the application
go build .

# Build with custom output name
go build -o laravel-dev .

# Run the application
go run .

# Run with arguments
go run . --help
```

### Testing
```bash
# Run all tests
go test ./...

# Run tests with verbose output
go test -v ./...

# Run tests for specific package
go test -v ./pkg

# Run tests matching pattern
go test -v -run TestFunctionName ./...

# Run tests with coverage
go test -v -cover ./...
```

### Linting & Formatting
```bash
# Format code (go fmt)
go fmt ./...

# Run go vet (static analysis)
go vet ./...

# Run all checks (vet + tests)
go vet ./... && go test ./...

# Show what go vet would change
go vet -n ./...

# Check for unused imports (requires golangci-lint)
golangci-lint run ./...
```

### Dependencies
```bash
# Download dependencies
go mod download

# Tidy dependencies (clean go.mod/go.sum)
go mod tidy

# Show dependencies
go list -m all

# Update a dependency
go get github.com/charmbracelet/huh@v0.1.0
```

---

## Code Style Guidelines

### Formatting

- Use `gofmt` for code formatting: `go fmt ./...`
- Use 4 spaces for indentation (Go standard)
- Maximum line length: 100 characters (soft limit)
- Group imports: standard library, then third-party packages
- Use blank lines between groups of imports

### Import Order
```go
import (
    // Standard library
    "fmt"
    "os"
    "strings"
    
    // Third-party packages (alphabetical)
    "github.com/charmbracelet/huh"
    "github.com/charmbracelet/lipgloss"
)
```

### Naming Conventions

- **Packages**: lowercase, short names (e.g., `pkg`, `cmd`)
- **Functions**: PascalCase for exported, camelCase for unexported
- **Variables**: camelCase, use meaningful names
- **Constants**: PascalCase for exported, camelCase for unexported
- **Interfaces**: end with `-er` suffix when simple (e.g., `Reader`, `Writer`)

### Types

- Use explicit types for clarity
- Prefer `int` over `int64` unless needed
- Use pointers (`*Type`) only when necessary (for mutation or nil values)
- Use interfaces for dependency injection

### Error Handling

- Always handle errors with proper context
- Return errors with descriptive messages
- Use `fmt.Errorf` with `%w` for wrapped errors
- Check errors immediately after function calls
- Never ignore errors with `_`

```go
// Good
if err != nil {
    return fmt.Errorf("failed to install PHP: %w", err)
}

// Bad
_ = someFunction()
```

### Functions

- Keep functions small and focused
- Use early returns to reduce nesting
- Document exported functions with comments

### Lip Gloss Styles

- Define styles as package-level variables in `var` blocks
- Chain methods for readability:
```go
var titleStyle = lipgloss.NewStyle().
    Bold(true).
    Foreground(lipgloss.Color("9")).
    Align(lipgloss.Center)
```

---

## Project Structure

```
laravel-dev/
├── Main.go              # Entry point
├── cmd/
│   └── Welcome.go       # Welcome screen UI
├── pkg/
│   ├── php.go           # PHP installation
│   ├── mariadb.go       # MariaDB installation
│   ├── nodejs.go       # Node.js installation
│   ├── composer.go      # Composer installation
│   └── valet.go         # Laravel Valet installation
├── go.mod               # Module definition
└── go.sum               # Dependency checksums
```

---

## Key Libraries

- **github.com/charmbracelet/lipgloss** - Terminal styling (v0.9.1)
- **github.com/charmbracelet/huh** - Interactive forms/selects (v0.1.0)
- **github.com/charmbracelet/bubbles** - TUI components (v0.16.1)
- **github.com/charmbracelet/bubbletea** - TUI framework (v0.25.0)

---

## Common Patterns

### Terminal Detection
```go
func isInteractiveTerminal() bool {
    stat, err := os.Stdin.Stat()
    if err != nil {
        return false
    }
    if (stat.Mode() & os.ModeCharDevice) == 0 {
        return false
    }
    if _, err := os.OpenFile("/dev/tty", os.O_RDWR, 0); err != nil {
        return false
    }
    return true
}
```

### Running Shell Commands
```go
func runCommand(name string, args ...string) error {
    cmd := exec.Command(name, args...)
    cmd.Stdout = os.Stdout
    cmd.Stderr = os.Stderr
    cmd.Stdin = os.Stdin

    if err := cmd.Run(); err != nil {
        return fmt.Errorf("error running %s: %w", name, err)
    }
    return nil
}
```

### Version Extraction
```go
func extractVersion(output string) string {
    output = strings.TrimSpace(output)
    // Handle specific formats
    if strings.Contains(output, "PHP") {
        parts := strings.Fields(output)
        if len(parts) >= 2 {
            return "v" + parts[1]
        }
    }
    // Fallback to regex
    re := regexp.MustCompile(`(\d+\.\d+(\.\d+)?)`)
    matches := re.FindStringSubmatch(output)
    if len(matches) > 1 {
        return "v" + matches[1]
    }
    return output
}
```

---

## Working with huh (Interactive Forms)

```go
// Single select
form := huh.NewForm(
    huh.NewGroup(
        huh.NewSelect[string]().
            Title("Select option").
            Options(
                huh.NewOption("Option 1", "value1"),
                huh.NewOption("Option 2", "value2"),
            ).
            Value(&selectedValue),
    ),
)

// Multi-select with custom keymap
form := huh.NewForm(
    huh.NewGroup(
        huh.NewMultiSelect[string]().
            Title("Select multiple").
            Options(
                huh.NewOption("Option 1", "opt1").Selected(true),
            ).
            Value(&selectedOptions),
    ),
).WithWidth(60).WithKeyMap(&huh.KeyMap{
    Quit: key.NewBinding(key.WithKeys("q", "ctrl+c"), key.WithHelp("q", "cancelar")),
})
```

---

## Working with Lip Gloss (Styling)

```go
// Style definitions
var titleStyle = lipgloss.NewStyle().
    Bold(true).
    Foreground(lipgloss.Color("9")).
    Width(50).
    Align(lipgloss.Center)

// Border styles
var boxStyle = lipgloss.NewStyle().
    Border(lipgloss.RoundedBorder()).
    BorderForeground(lipgloss.Color("9")).
    Padding(1, 2)

// Color palette
// 9 = Red
// 10 = Green  
// 12 = Blue
// 14 = Cyan
// 245 = Gray
// 236 = Dark background
```

---

## Notes for Agents

1. **Always run `go mod tidy`** after adding new dependencies
2. **Test in interactive terminal** - non-TTY environments hide TUI elements
3. **Handle errors gracefully** - show user-friendly messages
4. **Verify commands exist** before running system commands
5. **Use lipgloss styles consistently** - avoid inline styling
6. **Keep backwards compatibility** when adding new features
7. **Add proper error context** - wrap errors with descriptive messages

---

## Version Info

- Go: 1.24.0
- Module: laravel-dev
