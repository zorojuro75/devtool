package docker

import (
	"fmt"
	"path/filepath"
	"strings"

	"github.com/fatih/color"
)

// Remove removes one or more services from an existing docker-compose.yml
func Remove(services []string, dir string) error {
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
	if !fileExists(dir, "docker-compose.yml") {
		return fmt.Errorf("no docker-compose.yml found in %s\n  Nothing to remove from.", dir)
	}

	raw, err := ReadCompose(composePath)
	if err != nil {
		return fmt.Errorf("cannot read docker-compose.yml: %w", err)
	}

	fmt.Println()

	var removed []string
	var skipped []string

	for _, name := range services {
		if !HasServiceInRaw(raw, name) {
			yellow.Printf("  %s not found in docker-compose.yml -- skipped.\n", name)
			skipped = append(skipped, name)
			continue
		}

		// Warn if app depends on this service
		warnIfDepended(raw, name)

		// Remove the service
		RemoveServiceFromRaw(raw, name)
		removed = append(removed, name)
		green.Printf("+ Removed %s from docker-compose.yml\n", name)

		// Remove its volume if it has one
		if volName, ok := VolumeNames[name]; ok {
			if RemoveVolumeFromRaw(raw, volName) {
				green.Printf("+ Removed %s volume from docker-compose.yml\n", volName)
			}
		}
	}

	if len(removed) == 0 {
		fmt.Println()
		dim.Println("  Nothing was removed.")
		return nil
	}

	// Write back
	if err := WriteCompose(composePath, raw); err != nil {
		return fmt.Errorf("cannot write docker-compose.yml: %w", err)
	}

	// Print manual steps
	printRemoveHint(removed)

	return nil
}

// RemoveServiceFromRaw removes a service and its depends_on reference from raw compose
func RemoveServiceFromRaw(raw map[string]interface{}, name string) {
	services, ok := raw["services"].(map[string]interface{})
	if !ok {
		return
	}

	// Remove the service itself
	delete(services, name)

	// Remove it from any depends_on entries in other services
	for _, svcRaw := range services {
		svc, ok := svcRaw.(map[string]interface{})
		if !ok {
			continue
		}
		dependsOn, ok := svc["depends_on"].(map[string]interface{})
		if !ok {
			continue
		}
		delete(dependsOn, name)
		// Clean up empty depends_on
		if len(dependsOn) == 0 {
			delete(svc, "depends_on")
		}
	}
}

// RemoveVolumeFromRaw removes a named volume from the top-level volumes section
// Returns true if the volume was found and removed
func RemoveVolumeFromRaw(raw map[string]interface{}, volumeName string) bool {
	volumes, ok := raw["volumes"].(map[string]interface{})
	if !ok {
		return false
	}
	if _, exists := volumes[volumeName]; !exists {
		return false
	}
	delete(volumes, volumeName)
	// Clean up empty volumes section
	if len(volumes) == 0 {
		delete(raw, "volumes")
	}
	return true
}

// warnIfDepended warns the user if the app service depends on this service
func warnIfDepended(raw map[string]interface{}, name string) {
	yellow := color.New(color.FgYellow)

	services, ok := raw["services"].(map[string]interface{})
	if !ok {
		return
	}

	appSvc, ok := services["app"].(map[string]interface{})
	if !ok {
		return
	}

	dependsOn, ok := appSvc["depends_on"].(map[string]interface{})
	if !ok {
		return
	}

	if _, depended := dependsOn[name]; depended {
		yellow.Printf("  Warning: your app service depends on %s.\n", name)
		yellow.Println("  Removing it may cause your app container to fail to start.")
		yellow.Println("  The depends_on entry will be removed automatically.")
		fmt.Println()
	}
}

func printRemoveHint(removed []string) {
	cyan := color.New(color.FgCyan)
	dim := color.New(color.FgHiBlack)

	fmt.Println()
	dim.Println("  Note: running containers are not stopped automatically.")
	dim.Println("  To stop and remove them:")
	for _, name := range removed {
		cyan.Printf("  docker-compose stop %s && docker-compose rm %s\n", name, name)
	}

	fmt.Println()
	dim.Println("  To remove volume data:")
	for _, name := range removed {
		if volName, ok := VolumeNames[name]; ok {
			cyan.Printf("  docker volume rm $(basename $PWD)_%s\n", volName)
		}
	}
	fmt.Println()
}