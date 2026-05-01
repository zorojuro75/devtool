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

    fmt.Println()

    // Dockerfile
    dockerfilePath := filepath.Join(dir, "Dockerfile")
    if fileExists(dir, "Dockerfile") {
        if !AskOverwrite("Dockerfile") {
            dim.Println("  Skipped Dockerfile")
        } else {
            if err := renderTemplate("templates/nextjs/Dockerfile.tmpl", dockerfilePath, opts); err != nil {
                return fmt.Errorf("Dockerfile: %w", err)
            }
            green.Println("+ Created Dockerfile")
        }
    } else {
        if err := renderTemplate("templates/nextjs/Dockerfile.tmpl", dockerfilePath, opts); err != nil {
            return fmt.Errorf("Dockerfile: %w", err)
        }
        green.Println("+ Created Dockerfile")
    }

    // docker-compose.yml
    composePath := filepath.Join(dir, "docker-compose.yml")
    if fileExists(dir, "docker-compose.yml") {
        if !AskOverwrite("docker-compose.yml") {
            dim.Println("  Skipped docker-compose.yml")
        } else {
            if err := writeCompose(opts, composePath); err != nil {
                return fmt.Errorf("docker-compose.yml: %w", err)
            }
            green.Println("+ Created docker-compose.yml")
        }
    } else {
        if err := writeCompose(opts, composePath); err != nil {
            return fmt.Errorf("docker-compose.yml: %w", err)
        }
        green.Println("+ Created docker-compose.yml")
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
        color.New(color.FgYellow).Printf("  Warning: %s\n", err)
    } else if modified {
        green.Println("+ Updated next.config.ts (standalone output)")
    } else {
        dim.Println("  next.config.ts already has standalone output")
    }

    // .env.example
    if opts.HasServices() {
        added, err := AppendEnvExample(dir, opts)
        if err != nil {
            color.New(color.FgYellow).Printf("  Warning: could not update .env.example: %s\n", err)
        } else if len(added) > 0 {
            green.Println("+ Updated .env.example")
            PrintConnectionStrings(added)
        }
    }

    printNextSteps(opts)
    return nil
}

func writeCompose(opts *DockerOptions, path string) error {
    cf := BuildComposeFile(opts)
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
    cyan.Println("  docker-compose up -d")
    fmt.Println()

    if opts.HasDB() {
        dim.Println("  First time - run migrations inside the container:")
        cyan.Println("  docker-compose exec app npx prisma migrate deploy")
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
