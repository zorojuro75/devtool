# devtool docker

Add Docker to any project with two commands.

    devtool docker init     # generate Dockerfile + docker-compose.yml
    devtool docker add      # add a service to existing docker-compose.yml
    devtool docker remove   # remove a service from docker-compose.yml

---

## How it works

devtool docker reads your existing .env.local file to get real credentials
before generating any Docker files. This means:

- The MySQL/PostgreSQL container uses the same database name and password as your app
- No manual editing of docker-compose.yml after generation
- Your app container reads .env.local directly at runtime

If .env.local is not found, devtool falls back to sensible defaults.

---

## Prerequisites

- Docker Desktop installed and running
- Download: https://www.docker.com/products/docker-desktop/

Windows: open Docker Desktop from the Start menu and wait for the whale
icon in the taskbar to stop animating before running any docker commands.

---

## devtool docker init

Detects your project type, reads .env.local for credentials, asks which
services you need, and generates a complete Docker setup.

### Usage

    devtool docker init

### Demo

    $ cd myapp
    $ devtool docker init

    Detecting project type...
    Detected: Next.js (found next.config.ts)

    ? Confirm stack: Next.js [Y/n]: y
      Reading credentials from .env.local

    ? Select services to include:
      1) PostgreSQL 16
      2) MySQL 8
      3) Redis 7
      4) Mailhog (local email)
      5) None

    Enter choices (comma separated, e.g. 1,3): 2

      + MySQL 8

    ? App port [3000]: 3000

    + Created Dockerfile
    + Created docker-compose.yml
    + Created .dockerignore
    + Updated next.config.ts (standalone output)
      Credentials sourced from .env.local

    --------------------------------------------------
      Next steps
    --------------------------------------------------
      docker-compose up -d

      First time - run migrations inside the container:
      docker-compose exec app npx drizzle-kit generate
      docker-compose exec app npx drizzle-kit migrate
    --------------------------------------------------

### What gets generated

Dockerfile -- multi-stage Node 20 Alpine build:

    Stage 1: install dependencies (npm ci)
    Stage 2: build the Next.js app (npm run build)
    Stage 3: minimal production runner using standalone output

docker-compose.yml -- built from your .env.local credentials:

    app service     your Next.js app, reads .env.local at runtime
    db service      MySQL or PostgreSQL using your actual DB name and password
    depends_on      app waits for DB health check before starting

.dockerignore -- excludes node_modules, .next, .env files, build artifacts

next.config.ts -- updated with output: "standalone" for minimal image size

### Credential sourcing from .env.local

If your .env.local contains:

    DATABASE_URL="mysql://root:@127.0.0.1:3306/testapp"

The generated docker-compose.yml will have:

    mysql:
      environment:
        MYSQL_DATABASE: testapp
        MYSQL_ALLOW_EMPTY_PASSWORD: yes

If your .env.local contains:

    DATABASE_URL="postgresql://john:secret@localhost:5432/mydb"

The generated docker-compose.yml will have:

    postgres:
      environment:
        POSTGRES_USER: john
        POSTGRES_PASSWORD: secret
        POSTGRES_DB: mydb

### Supported services

| Service    | Image              | Ports             |
|------------|--------------------|-------------------|
| PostgreSQL | postgres:16-alpine | 5432              |
| MySQL      | mysql:8            | 3306              |
| Redis      | redis:7-alpine     | 6379              |
| Mailhog    | mailhog/mailhog    | 1025 SMTP, 8025 UI|

### Flags

| Flag              | Description                                    |
|-------------------|------------------------------------------------|
| --stack string    | nextjs (skip auto-detection)                   |
| --port int        | app port, default 3000                         |
| --services string | comma-separated: postgres,mysql,redis,mailhog  |
| --no-prompt       | skip prompts, use flags only                   |

### Non-interactive mode

    devtool docker init --services mysql --no-prompt

### If files already exist

If Dockerfile or docker-compose.yml already exist, devtool asks before
overwriting each file individually. You can skip any file you want to keep.

---

## devtool docker add

Add one or more services to an existing docker-compose.yml without
touching anything else in the file.

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
    Parses the file and adds only the new services
    Skips any service that already exists with a clear message
    Preserves all existing services and configuration

If docker-compose.yml does not exist:
    Creates a new docker-compose.yml with just the requested services
    Prints a note to run devtool docker init to add the app service

### Multiple services at once

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

---

## devtool docker remove

Remove one or more services from docker-compose.yml.

### Usage

    devtool docker remove [service...]

### Examples

    devtool docker remove redis
    devtool docker remove redis mailhog

### Demo

    $ devtool docker remove redis

    + Removed redis from docker-compose.yml
    + Removed redisdata volume from docker-compose.yml

      Note: running containers are not stopped automatically.
      To stop and remove them:
      docker-compose stop redis && docker-compose rm redis

      To remove volume data:
      docker volume rm myapp_redisdata

### What gets removed

- The service block from services:
- Its named volume from the top-level volumes: section
- Its depends_on entry from the app service if present

### What does NOT get removed

- Running containers (you must stop them manually)
- Actual Docker volume data on disk (you must remove it manually)
- Env vars from .env.local or .env.example (kept as reference)

### If service depends on app

If the app service has depends_on for the service you are removing,
devtool warns you before proceeding and removes the depends_on entry.

---

## Running your project with Docker Desktop

### Step 1 -- Start Docker Desktop

Open Docker Desktop from the Start menu (Windows) or Applications (Mac).
Wait for the whale icon to stop animating -- this means Docker is ready.

### Step 2 -- Start your stack

Run from your project folder:

    docker-compose up -d

The -d flag runs containers in the background (detached mode).

Docker will:
1. Build your Next.js app image (takes 1-3 minutes on first run)
2. Pull database and service images
3. Start all containers
4. Wait for health checks to pass before starting the app

### Step 3 -- Run database migrations (first time only)

For Prisma:

    docker-compose exec app npx prisma migrate deploy

For Drizzle:

    docker-compose exec app npx drizzle-kit generate
    docker-compose exec app npx drizzle-kit migrate

### Step 4 -- Open your app

    http://localhost:3000

### Step 5 -- Check container status

    docker-compose ps

You should see all services with status "running" or "healthy".

### Rebuild after code changes

    docker-compose up -d --build

The --build flag rebuilds the app image with your latest code changes.

### View logs

    # All containers
    docker-compose logs -f

    # Just the app
    docker-compose logs -f app

    # Just the database
    docker-compose logs -f mysql

### Stop everything

    docker-compose down

### Stop and remove all data (fresh start)

    docker-compose down -v

The -v flag also removes named volumes, deleting all database data.
Use this when you want a completely clean state.

### Open a shell inside a container

    # Inside the app container
    docker-compose exec app sh

    # Inside MySQL
    docker-compose exec mysql mysql -u root -p

---

## Connection strings reference

### PostgreSQL

    DATABASE_URL="postgresql://postgres:postgres@localhost:5432/myapp"

    # Inside Docker network (from app container to db container)
    DATABASE_URL="postgresql://postgres:postgres@postgres:5432/myapp"

### MySQL

    # With password
    DATABASE_URL="mysql://root:password@127.0.0.1:3306/myapp"

    # Without password (XAMPP default)
    DATABASE_URL="mysql://root:@127.0.0.1:3306/myapp"

    # Inside Docker network
    DATABASE_URL="mysql://root:password@mysql:3306/myapp"

### Redis

    REDIS_URL="redis://localhost:6379"

    # Inside Docker network
    REDIS_URL="redis://redis:6379"

### Mailhog

    SMTP_HOST="localhost"
    SMTP_PORT="1025"
    SMTP_FROM="noreply@myapp.local"

    # Web UI to view sent emails
    http://localhost:8025

---

## Troubleshooting

### Docker daemon not running

    Error: Cannot connect to the Docker daemon

Solution: Open Docker Desktop and wait for it to fully start.

### Port already in use

    Error: Bind for 0.0.0.0:3306 failed: port is already allocated

Solution: XAMPP MySQL is running on the same port. Stop it in the
XAMPP Control Panel before starting Docker MySQL, or change the port
in docker-compose.yml:

    ports:
      - "3307:3306"   # use 3307 on host instead

Then update DATABASE_URL:

    DATABASE_URL="mysql://root:@127.0.0.1:3307/myapp"

### App cannot connect to database

    Error: connect ECONNREFUSED 127.0.0.1:3306

Inside Docker, services communicate using their service name as hostname,
not localhost or 127.0.0.1. Update DATABASE_URL in .env.local for Docker:

    # For MySQL inside Docker
    DATABASE_URL="mysql://root:@mysql:3306/testapp"

    # For PostgreSQL inside Docker
    DATABASE_URL="postgresql://postgres:postgres@postgres:5432/myapp"

Then rebuild:

    docker-compose up -d --build

### App starts before database is ready

The depends_on with condition: service_healthy ensures Docker waits
for the database health check to pass before starting the app. If you
still see connection errors, the health check may need more time.
Increase retries in docker-compose.yml:

    healthcheck:
      retries: 10

### Cannot remove volume -- in use

    Error: volume is in use

Stop the containers first:

    docker-compose down
    docker volume rm myapp_mysqldata

---

## Note on YAML round-trip

When devtool docker add or devtool docker remove rewrites docker-compose.yml,
gopkg.in/yaml.v3 does not preserve comments. If you have written comments
in your docker-compose.yml they will be removed. All service definitions
and configuration values are preserved correctly.

---

Back to README: ../README.md