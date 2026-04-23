package next

import "time"

type NextOptions struct {
	ProjectName string
	Year        int
	DB          string // "prisma-sqlite"|"prisma-pg"|"drizzle-pg"|"drizzle-mysql"|"none"
	Auth        string // "better-auth-email"|"better-auth-github"|"none"
	State       bool
	Docker      bool
}

func (o *NextOptions) DBLabel() string {
	switch o.DB {
	case "prisma-sqlite":
		return "Prisma 7 + SQLite"
	case "prisma-pg":
		return "Prisma 7 + PostgreSQL"
	case "drizzle-pg":
		return "Drizzle + PostgreSQL"
	case "drizzle-mysql":
		return "Drizzle + MySQL"
	default:
		return "None"
	}
}

func (o *NextOptions) AuthLabel() string {
	switch o.Auth {
	case "better-auth-email":
		return "Better Auth (email + password)"
	case "better-auth-github":
		return "Better Auth (email + password + GitHub OAuth)"
	default:
		return "None"
	}
}

func (o *NextOptions) ORMLabel() string {
	switch o.DB {
	case "prisma-sqlite", "prisma-pg":
		return "Prisma 7"
	case "drizzle-pg", "drizzle-mysql":
		return "Drizzle ORM"
	default:
		return "None"
	}
}

func (o *NextOptions) IsPrisma() bool {
	return o.DB == "prisma-sqlite" || o.DB == "prisma-pg"
}

func (o *NextOptions) IsDrizzle() bool {
	return o.DB == "drizzle-pg" || o.DB == "drizzle-mysql"
}

func (o *NextOptions) IsPostgres() bool {
	return o.DB == "prisma-pg" || o.DB == "drizzle-pg"
}

func (o *NextOptions) IsMySQL() bool {
	return o.DB == "drizzle-mysql"
}

func (o *NextOptions) IsSQLite() bool {
	return o.DB == "prisma-sqlite"
}

func (o *NextOptions) HasGitHub() bool {
	return o.Auth == "better-auth-github"
}

func (o *NextOptions) HasAuth() bool {
	return o.Auth != "none" && o.Auth != ""
}

func (o *NextOptions) HasDB() bool {
	return o.DB != "none" && o.DB != ""
}

func (o *NextOptions) DBProvider() string {
	switch o.DB {
	case "prisma-sqlite":
		return "sqlite"
	case "prisma-pg":
		return "postgresql"
	case "drizzle-pg":
		return "pg"
	case "drizzle-mysql":
		return "mysql2"
	default:
		return ""
	}
}
// DrizzleDialect returns the dialect name for drizzle.config.ts
func (o *NextOptions) DrizzleDialect() string {
	switch o.DB {
	case "drizzle-pg":
		return "postgresql"
	case "drizzle-mysql":
		return "mysql"
	default:
		return "sqlite"
	}
}

// BetterAuthProvider returns the provider string for Better Auth adapter
func (o *NextOptions) BetterAuthProvider() string {
	switch o.DB {
	case "prisma-sqlite":
		return "sqlite"
	case "prisma-pg":
		return "postgresql"
	case "drizzle-pg":
		return "postgresql"
	case "drizzle-mysql":
		return "mysql"
	default:
		return "sqlite"
	}
}

func (o *NextOptions) SetDefaults() {
	if o.Year == 0 {
		o.Year = time.Now().Year()
	}
}