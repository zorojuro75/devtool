# devtool docker

Add Docker to any project with two commands.

    devtool docker init     # generate Dockerfile + docker-compose.yml
    devtool docker add      # add a service to existing docker-compose.yml

---

## devtool docker init

Detects your project type, asks which services you need, and generates a production-ready Docker setup.

### Demo

    $ cd myapp
    $ devtool docker init

    Detecting project type...
    Detected: Next.js (found next.config.ts)

    ? Confirm stack: Next.js [Y/n]: y

    ? Select services to include:
      1) PostgreSQL 16
      2) MySQL 8
      3) Redis 7
      4) Mailhog (local email)
      5) None

    Enter choices (comma separated, e.g. 1,3): 1,3
      + PostgreSQL 16
      + Redis 7

    ? App port [3000]: 3000

    + Created Dockerfile
    + Created docker-compose.yml
    + Created .dockerignore
    + Updated next.config.ts (standalone output)
    + Updated .env.example

      Connection strings added to .env.example:
      DATABASE_URL="postgresql://postgres:postgres@localhost:5432/myapp"
      REDIS_URL="redis://localhost:6379"

    --------------------------------------------------
      Next steps
    --------------------------------------------------
      docker-compose up -d

      First time - run migrations inside the container:
      docker-compose exec app npx prisma migrate deploy
    --------------------------------------------------

### What gets generated

Dockerfile -- multi-stage Node 20 Alpine build:
- Stage 1: install dependencies
- Stage 2: build the app
- Stage 3: minimal production runner (standalone output)

docker-compose.yml -- built programmatically for your selected services:
- app service wired to your selected database with health check dependency
- database service with health check
- named volumes for persistence

.dockerignore -- sensible defaults for Next.js projects

next.config.ts -- updated with output: "standalone" for smaller images

.env.example -- connection strings appended for selected services

### Supported services

| Service | Image | Ports |
|---------|-------|-------|
| PostgreSQL | postgres:16-alpine | 5432 |
| MySQL | mysql:8 | 3306 |
| Redis | redis:7-alpine | 6379 |
| Mailhog | mailhog/mailhog | 1025 (SMTP), 8025 (web UI) |

### Flags

| Flag | Description |
|------|-------------|
| --stack string | nextjs (skip auto-detection) |
| --port int | app port, default 3000 |
| --services string | comma-separated: postgres,mysql,redis,mailhog |
| --no-prompt | skip prompts, use flags only |

### Non-interactive mode

    devtool docker init --services postgres,redis --no-prompt

### If files already exist

If Dockerfile or docker-compose.yml already exist, the command asks before overwriting each one. You can skip individual files.

---

## devtool docker add

Add one or more services to an existing docker-compose.yml without touching anything else.

### Usage

    devtool docker add [service...]

### Examples

    devtool docker add redis
    devtool docker add postgres
    devtool docker add postgres redis
    devtool docker add postgres redis mailhog

### Demo

    $ devtool docker add redis

    + Added redis to docker-compose.yml
    + Updated .env.example

      Connection strings added to .env.example:
      REDIS_URL="redis://localhost:6379"

      Restart your stack to apply:
      docker-compose up -d redis

### Behaviour

If docker-compose.yml exists:
- Parses the existing file using gopkg.in/yaml.v3
- Adds only the new services, leaves everything else untouched
- Skips any service that already exists with a clear message

If docker-compose.yml does not exist:
- Creates a new docker-compose.yml with just the requested services
- Prints a note to run devtool docker init to add the app service

### Multiple services

    $ devtool docker add postgres redis mailhog

    + Added postgres to docker-compose.yml
    + Added redis to docker-compose.yml
    + Added mailhog to docker-compose.yml
    + Updated .env.example

      Connection strings:
      DATABASE_URL="postgresql://postgres:postgres@localhost:5432/myapp"
      REDIS_URL="redis://localhost:6379"
      SMTP_HOST="localhost"
      SMTP_PORT="1025"
      SMTP_FROM="noreply@myapp.local"

      Restart your stack to apply:
      docker-compose up -d postgres redis mailhog

### If service already exists

    $ devtool docker add redis

      redis already exists in docker-compose.yml -- skipped.

### Unknown service

    $ devtool docker add mongodb

    Error: unknown service "mongodb"
      Supported services: postgres, mysql, redis, mailhog

### Note on YAML round-trip

gopkg.in/yaml.v3 does not preserve comments when rewriting a file.
If you have written comments in your docker-compose.yml, they will be
removed when devtool docker add rewrites the file. All service
definitions and configuration are preserved correctly.

---

## Connection strings reference

### PostgreSQL

    DATABASE_URL="postgresql://postgres:postgres@localhost:5432/myapp"

### MySQL

    DATABASE_URL="mysql://root:@127.0.0.1:3306/myapp"

### Redis

    REDIS_URL="redis://localhost:6379"

### Mailhog SMTP

    SMTP_HOST="localhost"
    SMTP_PORT="1025"
    SMTP_FROM="noreply@myapp.local"

    # Web UI available at http://localhost:8025

---

## Running your stack

### Start everything

    docker-compose up -d

### Run migrations inside the container (first time)

    # Prisma
    docker-compose exec app npx prisma migrate deploy

    # Drizzle
    docker-compose exec app npx drizzle-kit generate
    docker-compose exec app npx drizzle-kit migrate

### Rebuild after code changes

    docker-compose up -d --build

### Stop everything

    docker-compose down

### Stop and remove volumes (fresh start)

    docker-compose down -v

---

Back to README: ../README.md