package docker

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// AddStandaloneOutput adds output: "standalone" to next.config.ts
// Returns true if the file was modified, false if it was already set
func AddStandaloneOutput(dir string) (bool, error) {
	// Try next.config.ts first, then next.config.js
	candidates := []string{"next.config.ts", "next.config.js"}
	var configPath string

	for _, name := range candidates {
		p := filepath.Join(dir, name)
		if fileExists(dir, name) {
			configPath = p
			break
		}
	}

	if configPath == "" {
		return false, fmt.Errorf("next.config.ts not found in %s", dir)
	}

	data, err := os.ReadFile(configPath)
	if err != nil {
		return false, fmt.Errorf("cannot read %s: %w", configPath, err)
	}

	content := string(data)

	// Already has standalone — skip
	if strings.Contains(content, "standalone") {
		return false, nil
	}

	// Try to insert output: "standalone" into the nextConfig object
	// Target pattern: const nextConfig: NextConfig = {
	// or:             const nextConfig = {
	targets := []string{
		"const nextConfig: NextConfig = {",
		"const nextConfig = {",
	}

	modified := false
	for _, target := range targets {
		if strings.Contains(content, target) {
			content = strings.Replace(
				content,
				target,
				target+"\n  output: \"standalone\",\n",
				1,
			)
			modified = true
			break
		}
	}

	if !modified {
		// Could not find the pattern — append a comment and warn
		return false, fmt.Errorf(
			"could not automatically add standalone output to %s\n"+
				"  Please add manually: output: \"standalone\" inside your nextConfig object",
			configPath,
		)
	}

	if err := os.WriteFile(configPath, []byte(content), 0644); err != nil {
		return false, fmt.Errorf("cannot write %s: %w", configPath, err)
	}

	return true, nil
}