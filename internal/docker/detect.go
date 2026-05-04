package docker

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

type StackInfo struct {
	Name        string
	DisplayName string
	DefaultPort int
}

type ORMInfo struct {
	Name        string
	DisplayName string
}

var knownStacks = map[string]StackInfo{
	"nextjs": {Name: "nextjs", DisplayName: "Next.js", DefaultPort: 3000},
}

// DetectStack looks at the current directory and returns the detected stack
func DetectStack(dir string) (*StackInfo, error) {
	if fileExists(dir, "next.config.ts") || fileExists(dir, "next.config.js") {
		info := knownStacks["nextjs"]
		return &info, nil
	}
	return nil, fmt.Errorf(
		"could not detect project type in %s\n"+
			"  Looked for: next.config.ts, next.config.js\n"+
			"  Run this command from inside your project directory.",
		dir,
	)
}

// DetectOrConfirm auto-detects the stack and returns it
func DetectOrConfirm(dir string) (*StackInfo, error) {
	info, err := DetectStack(dir)
	if err != nil {
		return nil, err
	}
	return info, nil
}

// DetectORM reads package.json and detects whether Prisma or Drizzle is used
func DetectORM(dir string) (*ORMInfo, error) {
	pkgPath := filepath.Join(dir, "package.json")
	data, err := os.ReadFile(pkgPath)
	if err != nil {
		return nil, fmt.Errorf("cannot read package.json: %w", err)
	}

	var pkg struct {
		Dependencies    map[string]string `json:"dependencies"`
		DevDependencies map[string]string `json:"devDependencies"`
	}

	if err := json.Unmarshal(data, &pkg); err != nil {
		return nil, fmt.Errorf("cannot parse package.json: %w", err)
	}

	// Check dependencies and devDependencies
	allDeps := make(map[string]string)
	for k, v := range pkg.Dependencies {
		allDeps[k] = v
	}
	for k, v := range pkg.DevDependencies {
		allDeps[k] = v
	}

	if _, ok := allDeps["drizzle-orm"]; ok {
		return &ORMInfo{Name: "drizzle", DisplayName: "Drizzle ORM"}, nil
	}
	if _, ok := allDeps["@prisma/client"]; ok {
		return &ORMInfo{Name: "prisma", DisplayName: "Prisma"}, nil
	}

	return &ORMInfo{Name: "none", DisplayName: "None"}, nil
}

// DetectTypeScript checks if the project uses TypeScript
func DetectTypeScript(dir string) bool {
	return fileExists(dir, "tsconfig.json")
}

func fileExists(dir, name string) bool {
	_, err := os.Stat(filepath.Join(dir, name))
	return err == nil
}