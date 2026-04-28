package next

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

type packageJSON struct {
	Name            string            `json:"name"`
	Version         string            `json:"version"`
	Private         bool              `json:"private"`
	Scripts         map[string]string `json:"scripts"`
	Dependencies    map[string]string `json:"dependencies"`
	DevDependencies map[string]string `json:"devDependencies"`
}

func buildPackageJSON(opts *NextOptions) packageJSON {
	deps := map[string]string{
		"next":                     "16.2.4",
		"react":                    "^19",
		"react-dom":                "^19",
		"clsx":                     "^2",
		"tailwind-merge":           "^2",
		"class-variance-authority": "^0.7",
		"lucide-react":             "^1.8.0",
		"zod":                      "^4",
		"dotenv":                   "^16",
	}

	devDeps := map[string]string{
		"typescript":          "^5",
		"@types/node":         "^20",
		"@types/react":        "^19",
		"@types/react-dom":    "^19",
		"tailwindcss":         "^3.4",
		"tailwindcss-animate": "^1",
		"autoprefixer":        "^10",
		"postcss":             "^8",
		"eslint":              "^9",
		"eslint-config-next":  "16.2.4",
	}

	scripts := map[string]string{
		"dev":   "next dev",
		"build": "next build",
		"start": "next start",
		"lint":  "next lint",
	}

	if opts.HasAuth() {
		deps["better-auth"] = "^1.6.9"
	}

	if opts.State {
		deps["zustand"] = "^5"
	}

	switch opts.DB {
	case "prisma-sqlite", "prisma-pg":
		deps["@prisma/client"] = "^7.7.0"
		devDeps["prisma"] = "^7.7.0"
		scripts["db:migrate"] = "prisma migrate dev"
		scripts["db:studio"] = "prisma studio"
		scripts["db:generate"] = "prisma generate"

	case "drizzle-pg":
		deps["drizzle-orm"] = "^0.45.2"
		deps["pg"] = "^8"
		devDeps["drizzle-kit"] = "^0.31.4"
		devDeps["@types/pg"] = "^8"
		scripts["db:migrate"] = "drizzle-kit migrate"
		scripts["db:studio"] = "drizzle-kit studio"
		scripts["db:generate"] = "drizzle-kit generate"

	case "drizzle-mysql":
		deps["drizzle-orm"] = "^0.45.2"
		deps["mysql2"] = "^3"
		devDeps["drizzle-kit"] = "^0.31.4"
		scripts["db:migrate"] = "drizzle-kit migrate"
		scripts["db:studio"] = "drizzle-kit studio"
		scripts["db:generate"] = "drizzle-kit generate"
	}

	if opts.HasAuth() && opts.HasDB() {
		if opts.IsPrisma() {
			deps["@better-auth/prisma-adapter"] = "^1.6.9"
		} else if opts.IsDrizzle() {
			deps["@better-auth/drizzle-adapter"] = "^1.6.9"
		}
	}

	return packageJSON{
		Name:            opts.ProjectName,
		Version:         "0.1.0",
		Private:         true,
		Scripts:         scripts,
		Dependencies:    deps,
		DevDependencies: devDeps,
	}
}

func writePackageJSON(opts *NextOptions, dir string) error {
	pkg := buildPackageJSON(opts)
	data, err := json.MarshalIndent(pkg, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to build package.json: %w", err)
	}
	return os.WriteFile(filepath.Join(dir, "package.json"), data, 0644)
}