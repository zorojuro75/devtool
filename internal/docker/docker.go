package docker

import (
	"embed"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"text/template"

	"github.com/fatih/color"
)

//go:embed templates
var templateFS embed.FS

func Generate(opts *DockerOptions, dir string) error {
	green := color.New(color.FgGreen, color.Bold)
	dim := color.New(color.FgHiBlack)
	yellow := color.New(color.FgYellow)

	// Read .env.local for real credentials
	envVars := ReadEnvLocal(dir)
	if envVars != nil {
		dim.Println("  Reading credentials from .env.local")
	} else {
		yellow.Println("  Warning: .env.local not found, using default credentials")
	}

	fmt.Println()

	// Pick Dockerfile template based on detected ORM
	dockerTemplatePath := "templates/nextjs/Dockerfile.none.tmpl"
	switch opts.ORM {
	case "drizzle":
		dockerTemplatePath = "templates/nextjs/Dockerfile.drizzle.tmpl"
	case "prisma":
		dockerTemplatePath = "templates/nextjs/Dockerfile.prisma.tmpl"
	}

	// Dockerfile
	dockerfilePath := filepath.Join(dir, "Dockerfile")
	if fileExists(dir, "Dockerfile") {
		if !AskOverwrite("Dockerfile") {
			dim.Println("  Skipped Dockerfile")
		} else {
			if err := renderTemplate(dockerTemplatePath, dockerfilePath, opts); err != nil {
				return fmt.Errorf("Dockerfile: %w", err)
			}
			green.Println("+ Created Dockerfile")
		}
	} else {
		if err := renderTemplate(dockerTemplatePath, dockerfilePath, opts); err != nil {
			return fmt.Errorf("Dockerfile: %w", err)
		}
		green.Println("+ Created Dockerfile")
	}

	// docker-compose.yml
	composePath := filepath.Join(dir, "docker-compose.yml")
	writeComposeFn := func() error {
		if err := writeCompose(opts, composePath, envVars); err != nil {
			return fmt.Errorf("docker-compose.yml: %w", err)
		}
		green.Println("+ Created docker-compose.yml")
		return nil
	}
	if fileExists(dir, "docker-compose.yml") {
		if !AskOverwrite("docker-compose.yml") {
			dim.Println("  Skipped docker-compose.yml")
		} else {
			if err := writeComposeFn(); err != nil {
				return err
			}
		}
	} else {
		if err := writeComposeFn(); err != nil {
			return err
		}
	}

	// .dockerignore
	ignorePath := filepath.Join(dir, ".dockerignore")
	if fileExists(dir, ".dockerignore") {
		if !AskOverwrite(".dockerignore") {
			dim.Println("  Skipped .dockerignore")
		} else {
			if err := renderTemplate("templates/nextjs/dockerignore.tmpl", ignorePath, opts); err != nil {
				return fmt.Errorf(".dockerignore: %w", err)
			}
			green.Println("+ Created .dockerignore")
		}
	} else {
		if err := renderTemplate("templates/nextjs/dockerignore.tmpl", ignorePath, opts); err != nil {
			return fmt.Errorf(".dockerignore: %w", err)
		}
		green.Println("+ Created .dockerignore")
	}

	// next.config.ts standalone output
	modified, err := AddStandaloneOutput(dir)
	if err != nil {
		yellow.Printf("  Warning: %s\n", err)
	} else if modified {
		green.Println("+ Updated next.config.ts (standalone output)")
	} else {
		dim.Println("  next.config.ts already has standalone output")
	}

	dim.Println("  Credentials sourced from .env.local")

	printNextSteps(opts)
	return nil
}

func writeCompose(opts *DockerOptions, path string, envVars map[string]string) error {
	cf := BuildComposeFile(opts, envVars)
	return cf.Write(path)
}

func renderTemplate(tmplPath, destPath string, opts *DockerOptions) error {
	content, err := templateFS.ReadFile(tmplPath)
	if err != nil {
		return fmt.Errorf("template not found %s: %w", tmplPath, err)
	}

	if err := os.MkdirAll(filepath.Dir(destPath), 0755); err != nil {
		return fmt.Errorf("cannot create directory: %w", err)
	}

	tmpl, err := template.New(filepath.Base(tmplPath)).Parse(string(content))
	if err != nil {
		return fmt.Errorf("template parse error: %w", err)
	}

	f, err := os.Create(destPath)
	if err != nil {
		return fmt.Errorf("cannot create file: %w", err)
	}
	defer f.Close()

	return tmpl.Execute(f, opts)
}

func printNextSteps(opts *DockerOptions) {
	green := color.New(color.FgGreen, color.Bold)
	cyan := color.New(color.FgCyan)
	dim := color.New(color.FgHiBlack)

	fmt.Println()
	fmt.Println(strings.Repeat("-", 50))
	green.Println("  Next steps")
	fmt.Println(strings.Repeat("-", 50))
	fmt.Println()

	dim.Println("  Start your stack:")
	cyan.Println("  docker-compose up -d")
	fmt.Println()

	if opts.HasDB() && opts.HasORM() {
		dim.Println("  First time - run migrations inside the container:")
		migrateCmd := opts.MigrateCommand()
		if migrateCmd != "" {
			cyan.Printf("  %s\n", migrateCmd)
		}
		fmt.Println()
	}

	if opts.HasMailhog() {
		dim.Println("  Mailhog web UI:")
		cyan.Println("  http://localhost:8025")
		fmt.Println()
	}

	fmt.Println(strings.Repeat("-", 50))
	fmt.Println()
}