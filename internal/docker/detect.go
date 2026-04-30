package docker

import (
	"fmt"
	"os"
	"path/filepath"
)

type StackInfo struct {
	Name        string // "nextjs"
	DisplayName string // "Next.js"
	DefaultPort int
}

var knownStacks = map[string]StackInfo{
	"nextjs": {Name: "nextjs", DisplayName: "Next.js", DefaultPort: 3000},
}

// DetectStack looks at the current directory and returns the detected stack.
// Returns an error if no known stack is found.
func DetectStack(dir string) (*StackInfo, error) {
	// Next.js detection — look for next.config.ts or next.config.js
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

// DetectOrConfirm auto-detects the stack and asks the user to confirm.
// If detection fails, returns the error — the caller handles the fallback.
func DetectOrConfirm(dir string) (*StackInfo, error) {
	info, err := DetectStack(dir)
	if err != nil {
		return nil, err
	}
	return info, nil
}

func fileExists(dir, name string) bool {
	_, err := os.Stat(filepath.Join(dir, name))
	return err == nil
}