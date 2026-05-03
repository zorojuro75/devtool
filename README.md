# devtool

A CLI tool for scaffolding production-ready projects and adding Docker to existing ones.

Built in Go as a portfolio project demonstrating embed.FS, text/template, interactive CLI prompts, programmatic file generation, and cross-compilation.

![CI](https://github.com/zorojuro75/devtool/actions/workflows/ci.yml/badge.svg)
![Go Version](https://img.shields.io/badge/go-1.22+-00ADD8?logo=go)
![License](https://img.shields.io/badge/license-MIT-green)

---

## Commands

| Command | What it does |
|---------|-------------|
| [devtool next](docs/next.md) | Scaffold a production-ready Next.js 16 fullstack project |
| [devtool docker](docs/docker.md) | Add Docker to any project |
| devtool version | Print version and build metadata |

---

## Installation

### Option 1 — Download binary (no Go required)

Download the latest binary for your platform from the [Releases](https://github.com/zorojuro75/devtool/releases) page:

| Platform | File |
|----------|------|
| Linux (amd64) | devtool-linux-amd64 |
| Linux (arm64) | devtool-linux-arm64 |
| macOS (Apple Silicon) | devtool-darwin-arm64 |
| macOS (Intel) | devtool-darwin-amd64 |
| Windows | devtool-windows-amd64.exe |

Linux / macOS:

    chmod +x devtool-linux-amd64
    sudo mv devtool-linux-amd64 /usr/local/bin/devtool

Windows: Rename to devtool.exe and add it to a folder on your PATH.

### Option 2 — Install with Go

    go install github.com/zorojuro75/devtool@latest

---

## Quick start

    # Scaffold a Next.js fullstack project
    devtool next myapp

    # Add Docker to an existing project
    cd myapp
    devtool docker init

    # Add a specific service to docker-compose.yml
    devtool docker add redis

---

## Documentation

- [devtool next](docs/next.md) — Next.js 16 scaffold with Better Auth, Prisma/Drizzle, Tailwind, shadcn/ui
- [devtool docker](docs/docker.md) — Docker init and service management

---

## Architecture

    devtool/
    |-- cmd/
    |   |-- root.go         Root command
    |   |-- next.go         devtool next command
    |   |-- docker.go       devtool docker command
    |   `-- version.go      devtool version command
    |-- internal/
    |   |-- next/
    |   |   |-- next.go         Orchestrator
    |   |   |-- options.go      NextOptions struct + helpers
    |   |   |-- prompt.go       Interactive prompts
    |   |   |-- package.go      Programmatic package.json builder
    |   |   |-- writer.go       embed.FS template renderer
    |   |   `-- templates/      Embedded project templates
    |   `-- docker/
    |       |-- docker.go       Orchestrator
    |       |-- detect.go       Project type detection
    |       |-- options.go      DockerOptions struct + helpers
    |       |-- prompt.go       Interactive prompts
    |       |-- compose.go      docker-compose.yml builder
    |       |-- services.go     Service definitions
    |       |-- add.go          docker add orchestrator
    |       |-- yamlrw.go       YAML read/write/merge
    |       |-- config.go       next.config.ts modifier
    |       |-- env.go          .env.example appender
    |       `-- templates/      Embedded Dockerfile templates
    |-- main.go
    |-- Makefile
    |-- go.mod
    `-- docs/
        |-- next.md
        `-- docker.md

---

## Key Go concepts demonstrated

embed.FS compiles all template files into the binary at build time with zero runtime dependencies.

Programmatic file generation builds package.json and docker-compose.yml as Go structs marshalled to JSON/YAML instead of using fragile template conditionals.

Separate template files per adapter splits auth templates by DB adapter rather than using long if-else chains inside one template file.

YAML round-trip merging uses gopkg.in/yaml.v3 to parse existing compose files and merge new services without losing existing content.

Cross-compilation produces binaries for 5 platforms from one make release command with version metadata injected via ldflags.

---

## Development

Prerequisites: Go 1.22+, Git 2.0+, make (Windows: choco install make)

    make build      # build for current platform
    make install    # install to GOPATH/bin with version info
    make release    # cross-compile for all 5 platforms
    make test       # run all tests
    make lint       # go vet
    make clean      # remove bin/ and dist/

---

## Roadmap

- devtool next -- tRPC option
- devtool next -- Resend email integration
- devtool next -- Stripe payments option
- devtool next -- SvelteKit support
- devtool docker -- Go and Laravel stack support
- devtool config -- manage config without editing YAML
- Homebrew tap

---

## License

MIT -- see LICENSE for details.