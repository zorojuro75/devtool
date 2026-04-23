package next

import (
	"embed"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/fatih/color"
)

//go:embed templates
var templateFS embed.FS

func Generate(opts *NextOptions) ([]string, error) {
	opts.SetDefaults()

	absDir, err := filepath.Abs(opts.ProjectName)
	if err != nil {
		return nil, fmt.Errorf("cannot resolve project path: %w", err)
	}
	dir := absDir

	if _, err := os.Stat(dir); err == nil {
		return nil, fmt.Errorf("directory %q already exists — delete it first or use a different name", opts.ProjectName)
	}

	if err := os.MkdirAll(dir, 0755); err != nil {
		return nil, fmt.Errorf("cannot create project directory: %w", err)
	}

	var created []string

	add := func(tmplPath, destPath string) error {
		full := filepath.Join(dir, destPath)
		if err := renderTemplate(templateFS, tmplPath, full, opts); err != nil {
			return err
		}
		created = append(created, destPath)
		return nil
	}

	// Base files
	// embed paths use safe names (no brackets/parens)
	// dest paths use real Next.js names
	baseFiles := [][2]string{
		{"templates/base/app/layout.tsx.tmpl", "app/layout.tsx"},
		{"templates/base/app/page.tsx.tmpl", "app/page.tsx"},
		{"templates/base/app/globals.css.tmpl", "app/globals.css"},
		{"templates/base/app/auth-login/page.tsx.tmpl", "app/(auth)/login/page.tsx"},
		{"templates/base/app/auth-register/page.tsx.tmpl", "app/(auth)/register/page.tsx"},
		{"templates/base/app/dashboard/page.tsx.tmpl", "app/(dashboard)/page.tsx"},
		{"templates/base/app/dashboard/layout.tsx.tmpl", "app/(dashboard)/layout.tsx"},
		{"templates/base/app/api/auth/catchall/route.ts.tmpl", "app/api/auth/[...all]/route.ts"},
		{"templates/base/components/layout/Navbar.tsx.tmpl", "components/layout/Navbar.tsx"},
		{"templates/base/components/shared/LoadingSpinner.tsx.tmpl", "components/shared/LoadingSpinner.tsx"},
		{"templates/base/hooks/useAuth.ts.tmpl", "hooks/useAuth.ts"},
		{"templates/base/types/index.ts.tmpl", "types/index.ts"},
		{"templates/base/lib/utils.ts.tmpl", "lib/utils.ts"},
		{"templates/base/next.config.ts.tmpl", "next.config.ts"},
		{"templates/base/tsconfig.json.tmpl", "tsconfig.json"},
		{"templates/base/tailwind.config.ts.tmpl", "tailwind.config.ts"},
		{"templates/base/postcss.config.js.tmpl", "postcss.config.js"},
		{"templates/base/components.json.tmpl", "components.json"},
	}

	for _, f := range baseFiles {
		if err := add(f[0], f[1]); err != nil {
			return nil, err
		}
	}

	// Database files
	if opts.HasDB() {
		for _, f := range dbTemplateFiles(opts) {
			if err := add(f[0], f[1]); err != nil {
				return nil, err
			}
		}
	}

	// Auth files
	if opts.HasAuth() {
		for _, f := range authTemplateFiles(opts) {
			if err := add(f[0], f[1]); err != nil {
				return nil, err
			}
		}
	}

	// State files
	if opts.State {
		stateFiles := [][2]string{
			{"templates/state/zustand/store/useUserStore.ts.tmpl", "store/useUserStore.ts"},
			{"templates/state/zustand/store/index.ts.tmpl", "store/index.ts"},
		}
		for _, f := range stateFiles {
			if err := add(f[0], f[1]); err != nil {
				return nil, err
			}
		}
	}

	// Docker files
	if opts.Docker {
		dockerFiles := [][2]string{
			{"templates/docker/Dockerfile.tmpl", "Dockerfile"},
			{"templates/docker/docker-compose.yml.tmpl", "docker-compose.yml"},
			{"templates/docker/dockerignore.tmpl", ".dockerignore"},
		}
		for _, f := range dockerFiles {
			if err := add(f[0], f[1]); err != nil {
				return nil, err
			}
		}
	}

	// package.json — built programmatically
	if err := writePackageJSON(opts, dir); err != nil {
		return nil, fmt.Errorf("package.json: %w", err)
	}
	created = append(created, "package.json")

	// .env.example — built programmatically
	if err := writeEnvExample(opts, dir); err != nil {
		return nil, fmt.Errorf(".env.example: %w", err)
	}
	created = append(created, ".env.example")

	// .gitignore
	if err := add("templates/base/gitignore.tmpl", ".gitignore"); err != nil {
		return nil, err
	}

	// SETUP.md
	if err := add("templates/base/SETUP.md.tmpl", "SETUP.md"); err != nil {
		return nil, err
	}

	// Static directories
	os.MkdirAll(filepath.Join(dir, "public"), 0755)
	os.WriteFile(filepath.Join(dir, "public", ".gitkeep"), []byte(""), 0644)
	os.MkdirAll(filepath.Join(dir, "components", "ui"), 0755)
	os.WriteFile(filepath.Join(dir, "components", "ui", ".gitkeep"), []byte(""), 0644)

	printSuccess(opts, created)
	return created, nil
}

func dbTemplateFiles(opts *NextOptions) [][2]string {
	base := "templates/db/" + opts.DB + "/"
	if opts.IsPrisma() {
		return [][2]string{
			{base + "prisma/schema.prisma.tmpl", "prisma/schema.prisma"},
			{base + "lib/db.ts.tmpl", "lib/db.ts"},
		}
	}
	if opts.IsDrizzle() {
		return [][2]string{
			{base + "drizzle/schema.ts.tmpl", "drizzle/schema.ts"},
			{base + "drizzle/migrate.ts.tmpl", "drizzle/migrate.ts"},
			{base + "drizzle.config.ts.tmpl", "drizzle.config.ts"},
			{base + "lib/db.ts.tmpl", "lib/db.ts"},
		}
	}
	return nil
}

func authTemplateFiles(opts *NextOptions) [][2]string {
	base := "templates/auth/" + opts.Auth + "/"

	authTmpl := base + "lib/auth-none.ts.tmpl"
	if opts.IsPrisma() {
		authTmpl = base + "lib/auth-prisma.ts.tmpl"
	} else if opts.IsDrizzle() {
		authTmpl = base + "lib/auth-drizzle.ts.tmpl"
	}

	return [][2]string{
		{authTmpl, "lib/auth.ts"},
		{base + "lib/auth-client.ts.tmpl", "lib/auth-client.ts"},
		{base + "app/api/auth/catchall/route.ts.tmpl", "app/api/auth/[...all]/route.ts"},
		{base + "proxy.ts.tmpl", "proxy.ts"},
	}
}

func writeEnvExample(opts *NextOptions, dir string) error {
	var sb strings.Builder

	switch opts.DB {
	case "prisma-sqlite":
		sb.WriteString("# Database - SQLite (local file, no server needed)\n")
		sb.WriteString("DATABASE_URL=\"file:./dev.db\"\n")
	case "prisma-pg", "drizzle-pg":
		sb.WriteString("# Database - PostgreSQL\n")
		sb.WriteString(fmt.Sprintf("DATABASE_URL=\"postgresql://postgres:postgres@localhost:5432/%s\"\n", opts.ProjectName))
	case "drizzle-mysql":
		sb.WriteString("# Database - MySQL\n")
		sb.WriteString(fmt.Sprintf("DATABASE_URL=\"mysql://root:password@localhost:3306/%s\"\n", opts.ProjectName))
	}

	if opts.HasAuth() {
		sb.WriteString("\n# Better Auth\n")
		sb.WriteString("# Generate secret: openssl rand -base64 32\n")
		sb.WriteString("BETTER_AUTH_SECRET=\"\"\n")
		sb.WriteString("BETTER_AUTH_URL=\"http://localhost:3000\"\n")
		sb.WriteString("NEXT_PUBLIC_APP_URL=\"http://localhost:3000\"\n")
	}

	if opts.HasGitHub() {
		sb.WriteString("\n# GitHub OAuth\n")
		sb.WriteString("# Create at: https://github.com/settings/applications/new\n")
		sb.WriteString("GITHUB_CLIENT_ID=\"\"\n")
		sb.WriteString("GITHUB_CLIENT_SECRET=\"\"\n")
	}

	return os.WriteFile(filepath.Join(dir, ".env.example"), []byte(sb.String()), 0644)
}

func printSuccess(opts *NextOptions, created []string) {
	green := color.New(color.FgGreen, color.Bold)
	cyan := color.New(color.FgCyan)
	dim := color.New(color.FgHiBlack)

	fmt.Println()
	green.Printf("Created %s\n", opts.ProjectName)
	dim.Printf("  %d files generated\n\n", len(created))

	fmt.Println(strings.Repeat("-", 50))
	green.Println("  Stack")
	fmt.Println(strings.Repeat("-", 50))
	printRow("Framework", "Next.js 16.2.4 (App Router)")
	printRow("Language", "TypeScript")
	printRow("Styling", "Tailwind CSS + shadcn/ui")
	if opts.HasAuth() {
		printRow("Auth", opts.AuthLabel())
	}
	if opts.HasDB() {
		printRow("Database", opts.DBLabel())
	}
	if opts.State {
		printRow("State", "Zustand")
	}
	if opts.Docker {
		printRow("Docker", "Dockerfile + docker-compose.yml")
	}
	fmt.Println(strings.Repeat("-", 50))
	fmt.Println()
	green.Println("  Next steps")
	fmt.Println()
	cyan.Printf("  cd %s\n", opts.ProjectName)
	fmt.Println()
	dim.Println("  Open SETUP.md for the complete setup guide.")
	fmt.Println()
}

func printRow(label, value string) {
	fmt.Printf("  %-12s %s\n", label, value)
}