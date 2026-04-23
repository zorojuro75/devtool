package next

import (
	"embed"
	"fmt"
	"os"
	"path/filepath"
	"text/template"
)

func renderTemplate(fs embed.FS, tmplPath string, destPath string, opts *NextOptions) error {
	content, err := fs.ReadFile(tmplPath)
	if err != nil {
		return fmt.Errorf("template not found %s: %w", tmplPath, err)
	}

	if err := os.MkdirAll(filepath.Dir(destPath), 0755); err != nil {
		return fmt.Errorf("cannot create directory for %s: %w", destPath, err)
	}

	tmpl, err := template.New(filepath.Base(tmplPath)).Parse(string(content))
	if err != nil {
		return fmt.Errorf("template parse error %s: %w", tmplPath, err)
	}

	f, err := os.Create(destPath)
	if err != nil {
		return fmt.Errorf("cannot create file %s: %w", destPath, err)
	}
	defer f.Close()

	if err := tmpl.Execute(f, opts); err != nil {
		return fmt.Errorf("template render error %s: %w", tmplPath, err)
	}

	return nil
}