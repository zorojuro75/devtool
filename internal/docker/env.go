package docker

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// AppendEnvExample appends service connection strings to .env.example
// If the file does not exist, it creates it
// Returns the list of env keys that were added
func AppendEnvExample(dir string, opts *DockerOptions) ([]EnvEntry, error) {
	envPath := filepath.Join(dir, ".env.example")

	// Read existing content
	existing := ""
	data, err := os.ReadFile(envPath)
	if err == nil {
		existing = string(data)
	}

	var added []EnvEntry
	var sb strings.Builder

	// Start from existing content
	sb.WriteString(existing)

	// Add a newline separator if file has content and doesn't end with newline
	if len(existing) > 0 && !strings.HasSuffix(existing, "\n") {
		sb.WriteString("\n")
	}

	// Append entries for each selected service
	for _, serviceName := range opts.Services {
		entries, ok := EnvEntries[serviceName]
		if !ok {
			continue
		}

		for _, entry := range entries {
			// Skip if key already exists in file
			if strings.Contains(existing, entry.Key+"=") ||
				strings.Contains(existing, entry.Key+"=\"") {
				continue
			}

			if entry.Comment != "" {
				sb.WriteString(fmt.Sprintf("\n# %s\n", entry.Comment))
			}
			sb.WriteString(fmt.Sprintf("%s=\"%s\"\n", entry.Key, entry.Value))
			added = append(added, entry)
		}
	}

	// Write back
	if err := os.WriteFile(envPath, []byte(sb.String()), 0644); err != nil {
		return nil, fmt.Errorf("cannot write .env.example: %w", err)
	}

	return added, nil
}

// PrintConnectionStrings prints the added env entries to the terminal
func PrintConnectionStrings(entries []EnvEntry) {
	if len(entries) == 0 {
		return
	}

	fmt.Println("\n  Connection strings added to .env.example:\n")
	for _, e := range entries {
		if e.Comment != "" {
			fmt.Printf("  %s=\"%s\"\n", e.Key, e.Value)
		} else {
			fmt.Printf("  %s=\"%s\"\n", e.Key, e.Value)
		}
	}
	fmt.Println()
}