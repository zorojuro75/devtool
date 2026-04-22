package scaffold

import (
	"bufio"
	_ "embed"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"text/template"
	"time"
)

// Templates are embedded in the binary at compile time.
// You will populate these files in the next step.
//go:embed templates/go/main.go.tmpl
var goMainTmpl string

//go:embed templates/go/README.md.tmpl
var goReadmeTmpl string

//go:embed templates/go/.gitignore.tmpl
var goGitignoreTmpl string

type templateFile struct {
	path    string
	content string
}

type templateVars struct {
	ProjectName string
	Year        int
	ModulePath  string
}

var frameworks = map[string]func(string) []templateFile{
	"go":      goTemplates,
	"laravel": laravelTemplates,
	"next":    nextTemplates,
}

func Generate(framework, name string, force bool) ([]string, error) {
	builder, ok := frameworks[framework]
	if !ok {
		keys := make([]string, 0, len(frameworks))
		for k := range frameworks {
			keys = append(keys, k)
		}
		return nil, fmt.Errorf("unknown framework %q — supported: %s", framework, strings.Join(keys, ", "))
	}

	if _, err := os.Stat(name); err == nil && !force {
		fmt.Printf("Directory %q already exists. Overwrite? [y/N] ", name)
		scanner := bufio.NewScanner(os.Stdin)
		scanner.Scan()
		if strings.ToLower(strings.TrimSpace(scanner.Text())) != "y" {
			return nil, fmt.Errorf("aborted")
		}
	}

	vars := templateVars{
		ProjectName: name,
		Year:        time.Now().Year(),
		ModulePath:  "github.com/zorojuro75/" + name,
	}

	files := builder(name)
	var created []string

	for _, f := range files {
		fullPath := filepath.Join(name, f.path)
		dir := filepath.Dir(fullPath)

		if err := os.MkdirAll(dir, 0755); err != nil {
			return nil, fmt.Errorf("cannot create directory %s: %w", dir, err)
		}

		tmpl, err := template.New(f.path).Parse(f.content)
		if err != nil {
			return nil, fmt.Errorf("template parse error in %s: %w", f.path, err)
		}

		file, err := os.Create(fullPath)
		if err != nil {
			return nil, fmt.Errorf("cannot create file %s: %w", fullPath, err)
		}

		if err := tmpl.Execute(file, vars); err != nil {
			file.Close()
			return nil, fmt.Errorf("template render error in %s: %w", f.path, err)
		}
		file.Close()
		created = append(created, filepath.Join(name, f.path))
	}

	return created, nil
}

func goTemplates(name string) []templateFile {
	return []templateFile{
		{path: "main.go", content: goMainTmpl},
		{path: "README.md", content: goReadmeTmpl},
		{path: ".gitignore", content: goGitignoreTmpl},
		{path: "cmd/.gitkeep", content: ""},
		{path: "internal/.gitkeep", content: ""},
	}
}

func laravelTemplates(name string) []templateFile {
	return []templateFile{
		{path: "README.md", content: "# {{.ProjectName}}\n\nA Laravel project.\n"},
		{path: ".env.example", content: "APP_NAME={{.ProjectName}}\nAPP_ENV=local\nAPP_KEY=\n"},
	}
}

func nextTemplates(name string) []templateFile {
	return []templateFile{
		{path: "README.md", content: "# {{.ProjectName}}\n\nA Next.js project.\n"},
		{path: ".gitignore", content: "node_modules/\n.next/\n.env.local\n"},
	}
}