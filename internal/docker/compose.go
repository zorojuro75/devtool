package docker

import (
	"fmt"
	"os"
	"strings"
)

// ComposeFile represents a docker-compose.yml file
type ComposeFile struct {
	Services map[string]*ComposeService
	Volumes  []string
}

// ComposeService represents a single service in docker-compose.yml
type ComposeService struct {
	Name        string
	Image       string
	Build       string
	Ports       []string
	EnvFile     []string
	Environment map[string]string
	Volumes     []string
	DependsOn   map[string]string // service -> condition
	HealthCheck *HealthCheck
	Restart     string
}

// NewComposeFile creates a new empty compose file
func NewComposeFile() *ComposeFile {
	return &ComposeFile{
		Services: make(map[string]*ComposeService),
		Volumes:  []string{},
	}
}

// AddService adds a service to the compose file
func (c *ComposeFile) AddService(name string, svc *ComposeService) {
	svc.Name = name
	c.Services[name] = svc
}

// AddVolume adds a named volume if not already present
func (c *ComposeFile) AddVolume(name string) {
	for _, v := range c.Volumes {
		if v == name {
			return
		}
	}
	c.Volumes = append(c.Volumes, name)
}

// HasService checks if a service already exists
func (c *ComposeFile) HasService(name string) bool {
	_, ok := c.Services[name]
	return ok
}

// Build builds the ComposeFile from DockerOptions
func BuildComposeFile(opts *DockerOptions) *ComposeFile {
	cf := NewComposeFile()

	// App service
	app := &ComposeService{
		Build:   ".",
		Ports:   []string{fmt.Sprintf("%d:%d", opts.Port, opts.Port)},
		EnvFile: []string{".env.local"},
		Restart: "unless-stopped",
	}

	// Add depends_on for DB services with health checks
	if opts.HasDB() {
		app.DependsOn = make(map[string]string)
		if opts.HasPostgres() {
			app.DependsOn["postgres"] = "service_healthy"
		}
		if opts.HasMySQL() {
			app.DependsOn["mysql"] = "service_healthy"
		}
	}

	cf.AddService("app", app)

	// Add selected services
	for _, serviceName := range opts.Services {
		def, ok := ServiceRegistry[serviceName]
		if !ok {
			continue
		}

		svc := &ComposeService{
			Image:       def.Image,
			Ports:       def.Ports,
			Environment: def.Environment,
			Volumes:     def.Volumes,
			HealthCheck: def.HealthCheck,
			Restart:     def.Restart,
		}

		cf.AddService(serviceName, svc)

		// Add named volume if needed
		if volName, ok := VolumeNames[serviceName]; ok {
			cf.AddVolume(volName)
		}
	}

	return cf
}

// Write writes the ComposeFile to a file at the given path
func (c *ComposeFile) Write(path string) error {
	content := c.render()
	return os.WriteFile(path, []byte(content), 0644)
}

// render produces the YAML string manually
// We write YAML manually instead of using a library to avoid
// extra dependencies and to have full control over formatting
func (c *ComposeFile) render() string {
	var sb strings.Builder

	sb.WriteString("services:\n")

	// Always write app first, then other services in order
	serviceOrder := []string{"app"}
	for _, name := range ServiceOrder {
		if _, ok := c.Services[name]; ok {
			serviceOrder = append(serviceOrder, name)
		}
	}

	for _, name := range serviceOrder {
		svc, ok := c.Services[name]
		if !ok {
			continue
		}
		sb.WriteString(fmt.Sprintf("  %s:\n", name))

		if svc.Build != "" {
			sb.WriteString(fmt.Sprintf("    build: %s\n", svc.Build))
		}
		if svc.Image != "" {
			sb.WriteString(fmt.Sprintf("    image: %s\n", svc.Image))
		}

		if len(svc.Ports) > 0 {
			sb.WriteString("    ports:\n")
			for _, p := range svc.Ports {
				sb.WriteString(fmt.Sprintf("      - \"%s\"\n", p))
			}
		}

		if len(svc.EnvFile) > 0 {
			sb.WriteString("    env_file:\n")
			for _, f := range svc.EnvFile {
				sb.WriteString(fmt.Sprintf("      - %s\n", f))
			}
		}

		if len(svc.Environment) > 0 {
			sb.WriteString("    environment:\n")
			// Write in a stable order
			envKeys := []string{}
			for k := range svc.Environment {
				envKeys = append(envKeys, k)
			}
			// Sort for deterministic output
			sortStrings(envKeys)
			for _, k := range envKeys {
				sb.WriteString(fmt.Sprintf("      %s: %s\n", k, svc.Environment[k]))
			}
		}

		if len(svc.Volumes) > 0 {
			sb.WriteString("    volumes:\n")
			for _, v := range svc.Volumes {
				sb.WriteString(fmt.Sprintf("      - %s\n", v))
			}
		}

		if len(svc.DependsOn) > 0 {
			sb.WriteString("    depends_on:\n")
			for dep, condition := range svc.DependsOn {
				sb.WriteString(fmt.Sprintf("      %s:\n", dep))
				sb.WriteString(fmt.Sprintf("        condition: %s\n", condition))
			}
		}

		if svc.HealthCheck != nil {
			sb.WriteString("    healthcheck:\n")
			sb.WriteString("      test:\n")
			for _, t := range svc.HealthCheck.Test {
				sb.WriteString(fmt.Sprintf("        - \"%s\"\n", t))
			}
			sb.WriteString(fmt.Sprintf("      interval: %s\n", svc.HealthCheck.Interval))
			sb.WriteString(fmt.Sprintf("      timeout: %s\n", svc.HealthCheck.Timeout))
			sb.WriteString(fmt.Sprintf("      retries: %d\n", svc.HealthCheck.Retries))
		}

		if svc.Restart != "" {
			sb.WriteString(fmt.Sprintf("    restart: %s\n", svc.Restart))
		}

		sb.WriteString("\n")
	}

	// Write volumes section if needed
	if len(c.Volumes) > 0 {
		sb.WriteString("volumes:\n")
		for _, v := range c.Volumes {
			sb.WriteString(fmt.Sprintf("  %s:\n", v))
		}
	}

	return sb.String()
}

// sortStrings sorts a string slice in place (simple insertion sort)
func sortStrings(s []string) {
	for i := 1; i < len(s); i++ {
		key := s[i]
		j := i - 1
		for j >= 0 && s[j] > key {
			s[j+1] = s[j]
			j--
		}
		s[j+1] = key
	}
}