package docker

import (
	"fmt"
	"path/filepath"
	"strings"

	"github.com/fatih/color"
)

// Add adds one or more services to an existing or new docker-compose.yml
func Add(services []string, dir string) error {
	green := color.New(color.FgGreen, color.Bold)
	yellow := color.New(color.FgYellow)
	dim := color.New(color.FgHiBlack)

	// Validate all service names first
	for _, name := range services {
		if _, ok := ServiceRegistry[name]; !ok {
			supported := strings.Join(ServiceOrder, ", ")
			return fmt.Errorf("unknown service %q\n  Supported services: %s", name, supported)
		}
	}

	composePath := filepath.Join(dir, "docker-compose.yml")
	composeExists := fileExists(dir, "docker-compose.yml")

	var added []string
	var skipped []string
	var allEnvEntries []EnvEntry

	fmt.Println()

	if composeExists {
		// Parse existing compose file
		raw, err := ReadCompose(composePath)
		if err != nil {
			return fmt.Errorf("cannot read docker-compose.yml: %w", err)
		}

		for _, name := range services {
			if HasServiceInRaw(raw, name) {
				yellow.Printf("  %s already exists in docker-compose.yml — skipped.\n", name)
				skipped = append(skipped, name)
				continue
			}
			AddServiceToRaw(raw, name)
			added = append(added, name)
			if entries, ok := EnvEntries[name]; ok {
				allEnvEntries = append(allEnvEntries, entries...)
			}
		}

		if len(added) > 0 {
			if err := WriteCompose(composePath, raw); err != nil {
				return fmt.Errorf("cannot write docker-compose.yml: %w", err)
			}
			for _, name := range added {
				green.Printf("+ Added %s to docker-compose.yml\n", name)
			}
		}

	} else {
		// No compose file — create one with just the requested services
		dim.Println("  No docker-compose.yml found.")
		dim.Println("  Creating docker-compose.yml with requested services...")
		fmt.Println()

		raw := NewRawCompose()
		for _, name := range services {
			AddServiceToRaw(raw, name)
			added = append(added, name)
			if entries, ok := EnvEntries[name]; ok {
				allEnvEntries = append(allEnvEntries, entries...)
			}
		}

		if err := WriteCompose(composePath, raw); err != nil {
			return fmt.Errorf("cannot write docker-compose.yml: %w", err)
		}
		green.Println("+ Created docker-compose.yml")
	}

	// Update .env.example
	if len(allEnvEntries) > 0 {
		if err := AppendEnvEntries(dir, allEnvEntries); err != nil {
			yellow.Printf("  Warning: could not update .env.example: %s\n", err)
		} else {
			green.Println("+ Updated .env.example")
		}
		PrintConnectionStrings(allEnvEntries)
	}

	// Print restart hint
	if len(added) > 0 {
		printRestartHint(added, composeExists)
	}

	return nil
}

func printRestartHint(added []string, composeExisted bool) {
	cyan := color.New(color.FgCyan)
	dim := color.New(color.FgHiBlack)

	fmt.Println()
	if composeExisted {
		dim.Println("  Restart your stack to apply:")
		cyan.Printf("  docker-compose up -d %s\n", strings.Join(added, " "))
	} else {
		dim.Println("  Note: run devtool docker init to add your app service.")
		dim.Println("  Then start your stack:")
		cyan.Println("  docker-compose up -d")
	}
	fmt.Println()
}