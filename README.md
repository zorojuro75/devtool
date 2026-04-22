# devtool

A developer CLI with AI assistance — scaffold projects, explain errors, and summarise git history, all from your terminal.

Built in Go as a portfolio project demonstrating goroutines, streaming HTTP, interfaces, embed.FS, and cross-compilation.

![CI](https://github.com/zorojuro75/devtool/actions/workflows/ci.yml/badge.svg)
![Go Version](https://img.shields.io/badge/go-1.21+-00ADD8?logo=go)
![License](https://img.shields.io/badge/license-MIT-green)

---

## Demo

```
$ devtool explain "fatal error: all goroutines are asleep - deadlock!"

Explanation:
**Error type & language**: Go runtime panic

**Root cause**
Your program created goroutines that are all waiting on channel operations
that will never be satisfied because no goroutine is left to unblock them.

**Fixes**
1. Close or drain channels after all sends are complete
2. Ensure WaitGroup.Done() is called for every Add()
3. Add a timeout case to blocking select statements
```

```
$ devtool gitlog --since 7d --format changelog

Git summary (changelog):
**Added**
- feat: add streaming explain command
- feat: add scaffold templates for go, laravel, next

**Changed**
- refactor: move AI client behind Completer interface

**Fixed**
- fix: context cancellation on Ctrl+C now exits cleanly
```

---

## Installation

### Option 1 — Download binary (no Go required)

Download the latest binary for your platform from the [Releases](https://github.com/zorojuro75/devtool/releases) page:

| Platform | File |
|----------|------|
| Linux (amd64) | `devtool-linux-amd64` |
| Linux (arm64) | `devtool-linux-arm64` |
| macOS (Apple Silicon) | `devtool-darwin-arm64` |
| macOS (Intel) | `devtool-darwin-amd64` |
| Windows | `devtool-windows-amd64.exe` |

**Linux / macOS:**
```bash
chmod +x devtool-linux-amd64
sudo mv devtool-linux-amd64 /usr/local/bin/devtool
```

**Windows:** Rename to `devtool.exe` and move to a folder on your PATH.

### Option 2 — Install with Go

```bash
go install github.com/zorojuro75/devtool@latest
```

---

## Setup

On first run, devtool creates a config file at `~/.devtool.yaml` automatically.
Open it and add your API key:

```yaml
provider: openrouter
api_key: sk-or-v1-xxxxxxxxxxxxxxxx
timeout: 30
model: "nvidia/nemotron-3-super-120b-a12b:free"
```

Get a **free** API key at [openrouter.ai](https://openrouter.ai) — no credit card required.

You can also set the API key via environment variable:

```bash
export DEVTOOL_API_KEY=sk-or-v1-xxxxxxxxxxxxxxxx
```

### Recommended free models on OpenRouter

| Model | Speed | Quality |
|-------|-------|---------|
| `nvidia/nemotron-3-super-120b-a12b:free` | Fast | Good |
| `z-ai/glm-4.5-air:free` | Fast | Good |
| `google/gemma-7b-it:free` | Medium | Good |

---

## Commands

### `scaffold` — Generate a project skeleton

```bash
devtool scaffold [framework] [project-name]
```

**Supported frameworks:**

```bash
devtool scaffold go myapi          # Go module with cmd/, internal/, Makefile
devtool scaffold laravel blog      # Laravel skeleton with .env.example
devtool scaffold next my-app       # Next.js 14 App Router with TypeScript
```

**Flags:**

| Flag | Description |
|------|-------------|
| `--force` | Overwrite existing directory without prompting |

---

### `explain` — Explain any error message using AI

```bash
devtool explain [error message] [flags]
```

**Examples:**

```bash
# Pass error as argument
devtool explain "panic: runtime error: index out of range [0] with length 0"

# Pipe from a log file
cat build.log | devtool explain

# Hint the language for better context
devtool explain "Undefined variable $name" --lang php
```

**Flags:**

| Flag | Description |
|------|-------------|
| `--lang string` | Hint the programming language (e.g. go, php, typescript) |

The response streams token-by-token — you see the explanation as it is generated, not all at once.

---

### `gitlog` — Summarise recent git commits using AI

```bash
devtool gitlog [flags]
```

**Examples:**

```bash
# Daily standup notes
devtool gitlog --since 1d --format standup

# Weekly changelog
devtool gitlog --since 7d --format changelog

# General summary
devtool gitlog --since 30d --format summary
```

**Flags:**

| Flag | Default | Description |
|------|---------|-------------|
| `--since string` | `1d` | Time range: `1d`, `7d`, `30d`, or `YYYY-MM-DD` |
| `--format string` | `summary` | Output style: `summary`, `changelog`, `standup` |
| `--branch string` | current | Branch to summarise |
| `--max-commits int` | `50` | Max commits to include in context |

---

### `version` — Print build metadata

```bash
devtool version
devtool version --json
```

Output:
```
version:    v1.0.0
build date: 2026-04-22T10:30:00Z
commit:     a1b2c3d
```

---

## Global Flags

These flags work with every command:

| Flag | Description |
|------|-------------|
| `--config string` | Path to config file (default: `~/.devtool.yaml`) |
| `--no-color` | Disable ANSI colour output |
| `--verbose` | Enable verbose HTTP logging |

---

## Architecture

```
devtool/
├── cmd/                        # Cobra command definitions (one file per command)
│   ├── root.go                 # Root command, global flags, Execute()
│   ├── scaffold.go
│   ├── explain.go
│   ├── gitlog.go
│   └── version.go
├── internal/
│   ├── ai/
│   │   ├── client.go           # Completer interface + OpenRouter implementation
│   │   └── stream.go           # SSE stream parser (bufio.Scanner)
│   ├── config/
│   │   └── config.go           # YAML loader, env var override, default file creation
│   ├── git/
│   │   └── log.go              # os/exec wrapper for git log with context timeout
│   └── scaffold/
│       ├── scaffold.go         # Template engine using embed.FS + text/template
│       └── templates/          # Embedded project templates
│           ├── go/
│           ├── laravel/
│           └── next/
├── main.go                     # Entry point — calls cmd.Execute() only
├── Makefile
└── go.mod
```

### Key Go concepts demonstrated

**Interfaces for testability**
The AI client is defined as a single-method interface:
```go
type Completer interface {
    Stream(ctx context.Context, prompt, system string) (io.Reader, error)
}
```
Tests inject a mock that returns a `strings.NewReader` — no real HTTP calls needed.

**Streaming HTTP with SSE parsing**
```go
// Each token printed as received, not buffered
scanner := bufio.NewScanner(resp.Body)
for scanner.Scan() {
    line := scanner.Text()
    if strings.HasPrefix(line, "data: ") {
        // parse JSON delta, write to stdout immediately
    }
}
```

**Context & cancellation**
Every blocking operation — HTTP requests, `os/exec` calls — uses `context.WithTimeout`. Ctrl+C cancels the entire call chain cleanly via `context.WithCancel`.

**embed.FS**
Template files are compiled directly into the binary:
```go
//go:embed templates/go/main.go.tmpl
var goMainTmpl string
```
Zero runtime file dependencies — the binary is fully self-contained.

**Table-driven tests**
```go
tests := []struct {
    name    string
    input   string
    want    string
}{
    {"single token", `data: {"choices":[...]}`, "Hello\n"},
    {"handles [DONE]", "data: [DONE]\n", "\n"},
}
for _, tt := range tests {
    t.Run(tt.name, func(t *testing.T) { ... })
}
```

---

## Development

### Prerequisites

- Go 1.21+
- Git 2.0+
- `make` (install via `choco install make` on Windows)

### Build

```bash
make build       # build for current platform → bin/devtool
make release     # cross-compile for all 5 platforms → dist/
make test        # run all tests
make lint        # go vet
make install     # install to $GOPATH/bin
make clean       # remove bin/ and dist/
```

### Running tests

```bash
make test
```

Expected output:
```
ok   github.com/zorojuro75/devtool/internal/ai        0.18s
ok   github.com/zorojuro75/devtool/internal/config    0.24s
ok   github.com/zorojuro75/devtool/internal/git       5.18s
```

### Adding a new command

1. Create `cmd/yourcommand.go` with a `newYourCmd()` function
2. Register it in `cmd/root.go` inside `Execute()`
3. Add business logic under `internal/yourpackage/`
4. Write table-driven tests in `internal/yourpackage/yourpackage_test.go`

---

## Roadmap

- [ ] `review` command — AI code review of a git diff
- [ ] Conversation history — persistent multi-turn context per project
- [ ] Local LLM support via Ollama (offline, no API key needed)
- [ ] Homebrew tap for `brew install devtool`
- [ ] Plugin system — user-defined sub-commands

---

## License

MIT — see [LICENSE](LICENSE) for details.