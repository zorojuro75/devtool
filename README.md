# devtool

A CLI tool for scaffolding production-ready Next.js 16 fullstack projects — with Better Auth, Prisma or Drizzle ORM, Tailwind CSS, shadcn/ui, and optional Zustand and Docker — all configured and wired together in under a minute.

Built in Go as a portfolio project demonstrating embed.FS, text/template, interactive CLI prompts, programmatic JSON generation, and cross-compilation.

![CI](https://github.com/zorojuro75/devtool/actions/workflows/ci.yml/badge.svg)
![Go Version](https://img.shields.io/badge/go-1.22+-00ADD8?logo=go)
![License](https://img.shields.io/badge/license-MIT-green)

---

## Demo

```
$ devtool next myapp

Configuring myapp
────────────────────────────────────────
? Database:
  1) Prisma + SQLite  (local, zero config)
  2) Prisma + PostgreSQL
  3) Drizzle + PostgreSQL
  4) Drizzle + MySQL
  5) None

Enter choice [1-5]: 1
  ✓ Prisma + SQLite  (local, zero config)

? Authentication:
  1) Better Auth — email + password
  2) Better Auth — email + password + GitHub OAuth
  3) None

Enter choice [1-3]: 1
  ✓ Better Auth — email + password

? Include Zustand (state management)? [y/N]: y
  ✓ Yes

? Include Docker + docker-compose? [y/N]: n
  ✗ No

✓ Created myapp
  26 files generated

──────────────────────────────────────────────────
  Stack
──────────────────────────────────────────────────
  Framework    Next.js 16.2.4 (App Router)
  Language     TypeScript
  Styling      Tailwind CSS + shadcn/ui
  Auth         Better Auth (email + password)
  Database     Prisma 7 + SQLite
  State        Zustand
──────────────────────────────────────────────────

  Next steps

  cd myapp

  Open SETUP.md for the complete setup guide.
```

---

## What gets generated

Running `devtool next myapp` produces a fully wired project:

```
myapp/
├── app/
│   ├── (auth)/
│   │   ├── login/page.tsx
│   │   └── register/page.tsx
│   ├── (dashboard)/
│   │   ├── layout.tsx
│   │   └── page.tsx
│   ├── api/auth/[...all]/route.ts
│   ├── globals.css
│   └── layout.tsx
├── components/
│   ├── ui/                     shadcn/ui components go here
│   ├── layout/Navbar.tsx
│   └── shared/LoadingSpinner.tsx
├── lib/
│   ├── auth.ts                 Better Auth server config
│   ├── auth-client.ts          Better Auth React hooks
│   ├── db.ts                   Database client singleton
│   ├── env.ts                  Zod environment validation
│   └── utils.ts                cn() Tailwind helper
├── store/                      Zustand stores (if selected)
│   ├── index.ts
│   └── useUserStore.ts
├── hooks/useAuth.ts
├── types/index.ts
├── prisma/schema.prisma        (if Prisma selected)
├── drizzle/                    (if Drizzle selected)
│   ├── schema.ts
│   ├── migrate.ts
│   └── migrations/
├── Dockerfile                  (if Docker selected)
├── docker-compose.yml          (if Docker selected)
├── proxy.ts                    Route protection (Next.js 16)
├── .env.example
├── package.json                Built programmatically for your stack
├── tsconfig.json
├── tailwind.config.ts
├── components.json
└── SETUP.md                    Step-by-step setup guide for your stack
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

**Windows:** Rename to `devtool.exe` and add it to a folder on your PATH.

### Option 2 — Install with Go

```bash
go install github.com/zorojuro75/devtool@latest
```

---

## Usage

### Interactive mode

Answer prompts to configure your stack:

```bash
devtool next myapp
```

### Non-interactive mode

Skip prompts using flags — useful for scripts and dotfiles:

```bash
devtool next myapp \
  --db prisma-sqlite \
  --auth better-auth-email \
  --state \
  --docker \
  --no-prompt
```

### Flags

| Flag | Description |
|------|-------------|
| `--db string` | Database: `prisma-sqlite` \| `prisma-pg` \| `drizzle-pg` \| `drizzle-mysql` \| `none` |
| `--auth string` | Auth: `better-auth-email` \| `better-auth-github` \| `none` |
| `--state` | Include Zustand state management |
| `--docker` | Include Dockerfile + docker-compose |
| `--no-prompt` | Skip interactive prompts, require all flags |

---

## After generation

Every generated project includes a `SETUP.md` tailored to your exact stack. It covers:

- Installing dependencies
- Setting up `.env.local` with all required variables
- Generating and running database migrations
- Installing shadcn/ui components
- Starting the dev server
- Using Better Auth in your components
- Docker setup (if selected)

### Environment validation

`lib/env.ts` uses Zod to validate all required environment variables at startup. If anything is missing or invalid, you see a clear error immediately:

```
🚨 Invalid environment variables:

  ✗ BETTER_AUTH_SECRET: must be at least 32 characters.
    Generate one with: openssl rand -base64 32
  ✗ DATABASE_URL: DATABASE_URL is required.
    Example: file:./dev.db

Fix the above errors in your .env.local file and restart the server.
```

No more confusing runtime crashes from missing env vars.

---

## Stack details

### Framework
- **Next.js 16.2.4** — App Router, TypeScript, Turbopack

### Styling
- **Tailwind CSS 3.4** with CSS variables for theming
- **shadcn/ui** — component library (added via `npx shadcn@latest add`)

### Authentication
- **Better Auth 1.6.7** — email/password and optional GitHub OAuth
- Pre-wired API route at `app/api/auth/[...all]/route.ts`
- Route protection via `proxy.ts` (Next.js 16 middleware replacement)
- Ready-to-use hooks via `lib/auth-client.ts`

### Database options

| Option | ORM | Driver |
|--------|-----|--------|
| `prisma-sqlite` | Prisma 7 | SQLite (local file, no server) |
| `prisma-pg` | Prisma 7 | PostgreSQL |
| `drizzle-pg` | Drizzle 0.45 | PostgreSQL |
| `drizzle-mysql` | Drizzle 0.45 | MySQL |

### State management
- **Zustand 5** with `persist` middleware
- Pre-configured `useUserStore` for auth state

### Docker
- Multi-stage Dockerfile (Node 20 Alpine)
- `docker-compose.yml` with app + database services
- Health checks on database before app starts

---

## Version

```bash
devtool version
devtool version --json
```

```
version:    v1.1.0
build date: 2026-04-25T10:30:00Z
commit:     a1b2c3d
```

---

## Architecture

```
devtool/
├── cmd/
│   ├── root.go         Root command
│   ├── next.go         devtool next — fullstack scaffold command
│   └── version.go      devtool version — build metadata
├── internal/
│   └── next/
│       ├── next.go     Orchestrator — generates all files
│       ├── options.go  NextOptions struct + helper methods
│       ├── prompt.go   Interactive CLI prompts
│       ├── package.go  Programmatic package.json builder
│       ├── writer.go   embed.FS template renderer
│       └── templates/
│           ├── base/   Always generated
│           ├── db/     prisma-sqlite, prisma-pg, drizzle-pg, drizzle-mysql
│           ├── auth/   better-auth-email, better-auth-github
│           ├── state/  Zustand stores
│           └── docker/ Dockerfile + docker-compose
├── main.go
├── Makefile
└── go.mod
```

### Key Go concepts demonstrated

**embed.FS — zero runtime dependencies**
All 30+ template files are compiled into the binary at build time:
```go
//go:embed templates
var templateFS embed.FS
```

**Programmatic package.json generation**
Instead of fragile template conditionals, `package.json` is built as a Go struct and marshalled to JSON. Each stack option adds its own dependencies:
```go
func buildPackageJSON(opts *NextOptions) packageJSON {
    deps := map[string]string{"next": "16.2.4", ...}
    if opts.IsDrizzle() {
        deps["drizzle-orm"] = "^0.45.2"
        deps["drizzle-kit"] = "^0.31.4"
    }
    // ...
}
```

**Separate template files per adapter**
Auth templates are split by DB adapter (`auth-prisma.ts.tmpl`, `auth-drizzle.ts.tmpl`) rather than using `{{if}}` chains inside one template. This makes each file clean and independently maintainable.

**Cross-compilation**
A single `make release` produces binaries for 5 platforms:
```makefile
release:
    GOOS=linux   GOARCH=amd64 go build -o dist/devtool-linux-amd64 .
    GOOS=darwin  GOARCH=arm64 go build -o dist/devtool-darwin-arm64 .
    GOOS=windows GOARCH=amd64 go build -o dist/devtool-windows-amd64.exe .
    # ...
```

---

## Development

### Prerequisites

- Go 1.22+
- Git 2.0+
- `make` (Windows: `choco install make`)

### Commands

```bash
make build      # build for current platform → bin/devtool
make install    # install to $GOPATH/bin with version info
make release    # cross-compile for all 5 platforms → dist/
make test       # run all tests
make lint       # go vet
make clean      # remove bin/ and dist/
```

---

## Roadmap

- [ ] `devtool next` — SvelteKit and Nuxt support
- [ ] `devtool next` — tRPC option
- [ ] `devtool next` — Resend email integration
- [ ] `devtool next` — Stripe payments integration
- [ ] `devtool config` — manage config without editing YAML
- [ ] Homebrew tap — `brew install devtool`
- [ ] Shell completion — `devtool completion zsh`

---

## License

MIT — see [LICENSE](LICENSE) for details.