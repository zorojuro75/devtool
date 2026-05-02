package docker

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

// ReadCompose reads a docker-compose.yml into a generic map
func ReadCompose(path string) (map[string]interface{}, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("cannot read %s: %w", path, err)
	}

	var raw map[string]interface{}
	if err := yaml.Unmarshal(data, &raw); err != nil {
		return nil, fmt.Errorf("cannot parse %s: %w", path, err)
	}

	// Ensure services key exists
	if raw == nil {
		raw = make(map[string]interface{})
	}
	if _, ok := raw["services"]; !ok {
		raw["services"] = make(map[string]interface{})
	}

	return raw, nil
}

// WriteCompose marshals the map back to YAML and writes it
func WriteCompose(path string, raw map[string]interface{}) error {
	data, err := yaml.Marshal(raw)
	if err != nil {
		return fmt.Errorf("cannot marshal docker-compose.yml: %w", err)
	}
	return os.WriteFile(path, data, 0644)
}

// NewRawCompose creates a minimal empty compose structure
func NewRawCompose() map[string]interface{} {
	return map[string]interface{}{
		"services": map[string]interface{}{},
	}
}

// HasServiceInRaw checks if a service name already exists in the compose file
func HasServiceInRaw(raw map[string]interface{}, name string) bool {
	services, ok := raw["services"].(map[string]interface{})
	if !ok {
		return false
	}
	_, exists := services[name]
	return exists
}

// AddServiceToRaw adds a service definition to the raw compose map
func AddServiceToRaw(raw map[string]interface{}, name string) {
	def, ok := ServiceRegistry[name]
	if !ok {
		return
	}

	// Build service map
	svc := map[string]interface{}{
		"image":   def.Image,
		"restart": def.Restart,
	}

	if len(def.Ports) > 0 {
		svc["ports"] = def.Ports
	}

	if len(def.Environment) > 0 {
		svc["environment"] = def.Environment
	}

	if len(def.Volumes) > 0 {
		svc["volumes"] = def.Volumes
		// Add to top-level volumes section
		if volName, ok := VolumeNames[name]; ok {
			AddVolumeToRaw(raw, volName)
		}
	}

	if def.HealthCheck != nil {
		svc["healthcheck"] = map[string]interface{}{
			"test":     def.HealthCheck.Test,
			"interval": def.HealthCheck.Interval,
			"timeout":  def.HealthCheck.Timeout,
			"retries":  def.HealthCheck.Retries,
		}
	}

	// Add to services
	services := raw["services"].(map[string]interface{})
	services[name] = svc
}

// AddVolumeToRaw adds a named volume to the top-level volumes section
func AddVolumeToRaw(raw map[string]interface{}, volumeName string) {
	if _, ok := raw["volumes"]; !ok {
		raw["volumes"] = map[string]interface{}{}
	}

	volumes, ok := raw["volumes"].(map[string]interface{})
	if !ok {
		return
	}

	// Only add if not already present
	if _, exists := volumes[volumeName]; !exists {
		volumes[volumeName] = nil
	}
}