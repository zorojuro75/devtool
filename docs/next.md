# devtool next

Scaffold a production-ready Next.js 16 fullstack project with a single command.

    devtool next myapp

Interactive prompts configure your entire stack. The generated project is fully wired and ready to run after npm install and filling in .env.local.

---

## Demo

    $ devtool next myapp

    Configuring myapp
    ----------------------------------------
    ? Database:
      1) Prisma + SQLite  (local, zero config)
      2) Prisma + PostgreSQL
      3) Drizzle + PostgreSQL
      4) Drizzle + MySQL
      5) None

    Enter choice [1-5]: 1
      + Prisma + SQLite

    ? Authentication:
      1) Better Auth -- email + password
      2) Better Auth -- email + password + GitHub OAuth
      3) None

    Enter choice [1-3]: 1
      + Better Auth -- email + password

    ? Include Zustand (state management)? [y/N]: y
      + Yes

    ? Include Docker + docker-compose? [y/N]: n
      - No

    + Created myapp
      26 files generated

    --------------------------------------------------
      Stack
    --------------------------------------------------
      Framework    Next.js 16.2.4 (App Router)
      Language     TypeScript
      Styling      Tailwind CSS + shadcn/ui
      Auth         Better Auth (email + password)
      Database     Prisma 7 + SQLite
      State        Zustand
    --------------------------------------------------

      cd myapp
      Open SETUP.md for the complete setup guide.

---

## Usage

### Interactive mode

    devtool next myapp

### Non-interactive mode

    devtool next myapp --db prisma-sqlite --auth better-auth-email --state --docker --no-prompt

### Flags

| Flag | Description |
|------|-------------|
| --db string | prisma-sqlite, prisma-pg, drizzle-pg, drizzle-mysql, none |
| --auth string | better-auth-email, better-auth-github, none |
| --state | Include Zustand state management |
| --docker | Include Dockerfile + docker-compose |
| --no-prompt | Skip interactive prompts, require all flags |

---

## What gets generated

    myapp/
    |-- app/
    |   |-- (auth)/
    |   |   |-- login/page.tsx          Login page
    |   |   `-- register/page.tsx       Register page
    |   |-- (dashboard)/
    |   |   |-- layout.tsx              Protected layout
    |   |   `-- page.tsx                Dashboard page
    |   |-- api/auth/[...all]/route.ts  Better Auth API handler
    |   |-- globals.css
    |   `-- layout.tsx
    |-- components/
    |   |-- ui/                         shadcn/ui components
    |   |-- layout/Navbar.tsx
    |   `-- shared/LoadingSpinner.tsx
    |-- lib/
    |   |-- auth.ts                     Better Auth server config
    |   |-- auth-client.ts              Better Auth React hooks
    |   |-- db.ts                       Database client singleton
    |   |-- env.ts                      Zod environment validation
    |   `-- utils.ts                    cn() Tailwind helper
    |-- store/                          Zustand (if selected)
    |   |-- index.ts
    |   `-- useUserStore.ts
    |-- hooks/useAuth.ts
    |-- types/index.ts
    |-- prisma/schema.prisma            (if Prisma selected)
    |-- drizzle/                        (if Drizzle selected)
    |   |-- schema.ts
    |   |-- migrate.ts
    |   `-- migrations/
    |-- Dockerfile                      (if Docker selected)
    |-- docker-compose.yml              (if Docker selected)
    |-- proxy.ts                        Route protection (Next.js 16)
    |-- .env.example
    |-- package.json
    |-- tsconfig.json
    |-- tailwind.config.ts
    |-- components.json
    `-- SETUP.md                        Step-by-step guide for your stack

---

## Stack details

### Framework
Next.js 16.2.4 with App Router, TypeScript, and Turbopack.

### Styling
Tailwind CSS 3.4 with CSS variables for theming and shadcn/ui for components.
Add components with: npx shadcn@latest add button card input label

### Authentication

| Option | What is included |
|--------|-----------------|
| better-auth-email | Email + password sign up and sign in |
| better-auth-github | Email + password + GitHub OAuth |

Pre-wired files:
- app/api/auth/[...all]/route.ts -- API handler
- proxy.ts -- route protection (Next.js 16 replaces middleware.ts with proxy.ts)
- lib/auth-client.ts -- React hooks

### Database options

| Option | ORM | Driver |
|--------|-----|--------|
| prisma-sqlite | Prisma 7 | SQLite, no server needed |
| prisma-pg | Prisma 7 | PostgreSQL |
| drizzle-pg | Drizzle 0.45 | PostgreSQL |
| drizzle-mysql | Drizzle 0.45 | MySQL |

### State management
Zustand 5 with persist middleware. Pre-configured useUserStore with setUser and clearUser.

### Environment validation
Zod v4 validates all required env vars at startup via lib/env.ts with clear error messages.

---

## Setup after generation

Open SETUP.md inside the generated project -- it has every command for your exact stack.

### Prisma + SQLite quick start

    cd myapp
    npm install
    cp .env.example .env.local
    # Add BETTER_AUTH_SECRET to .env.local
    npx prisma migrate dev --name init
    npx shadcn@latest add button card input label
    npm run dev

### Drizzle + MySQL quick start

    cd myapp
    npm install
    cp .env.example .env.local
    # Add DATABASE_URL and BETTER_AUTH_SECRET to .env.local
    npx drizzle-kit generate
    npx drizzle-kit migrate
    npx shadcn@latest add button card input label
    npm run dev

---

## Using Better Auth in components

### Sign up

    import { authClient } from "@/lib/auth-client";

    await authClient.signUp.email({
      email: "user@example.com",
      password: "password123",
      name: "John Doe",
    });

### Sign in

    await authClient.signIn.email({
      email: "user@example.com",
      password: "password123",
    });

### Sign out

    await authClient.signOut();

### Get session in a client component

    import { useSession } from "@/lib/auth-client";

    const { data: session, isPending } = useSession();

### Get session in a server component

    import { auth } from "@/lib/auth";
    import { headers } from "next/headers";

    const session = await auth.api.getSession({ headers: await headers() });

---

## Environment variables

### Always required

    NEXT_PUBLIC_APP_URL="http://localhost:3000"

### If auth is selected

    BETTER_AUTH_SECRET=""   # generate: openssl rand -base64 32
    BETTER_AUTH_URL="http://localhost:3000"

### If GitHub OAuth is selected

    GITHUB_CLIENT_ID=""
    GITHUB_CLIENT_SECRET=""

### Database

    # SQLite
    DATABASE_URL="file:./dev.db"

    # PostgreSQL
    DATABASE_URL="postgresql://postgres:postgres@localhost:5432/myapp"

    # MySQL
    DATABASE_URL="mysql://root:@127.0.0.1:3306/myapp"

---

Back to README: ../README.md