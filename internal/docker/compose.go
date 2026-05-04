package docker

import (
	"fmt"
	"os"
	"strings"
)

type ComposeFile struct {
	Services map[string]*ComposeService
	Volumes  []string
}

type ComposeService struct {
	Name        string
	Image       string
	Build       string
	Ports       []string
	EnvFile     []string
	Environment map[string]string
	Volumes     []string
	DependsOn   map[string]string
	HealthCheck *HealthCheck
	Restart     string
}

func NewComposeFile() *ComposeFile {
	return &ComposeFile{
		Services: make(map[string]*ComposeService),
		Volumes:  []string{},
	}
}

func (c *ComposeFile) AddService(name string, svc *ComposeService) {
	svc.Name = name
	c.Services[name] = svc
}

func (c *ComposeFile) AddVolume(name string) {
	for _, v := range c.Volumes {
		if v == name {
			return
		}
	}
	c.Volumes = append(c.Volumes, name)
}

func (c *ComposeFile) HasService(name string) bool {
	_, ok := c.Services[name]
	return ok
}

// BuildComposeFile builds the ComposeFile from DockerOptions
// reads .env.local to get real DB credentials when available
func BuildComposeFile(opts *DockerOptions, envVars map[string]string) *ComposeFile {
	cf := NewComposeFile()

	// App service
	app := &ComposeService{
		Build:   ".",
		Ports:   []string{fmt.Sprintf("%d:%d", opts.Port, opts.Port)},
		EnvFile: []string{".env.local"},
		Restart: "unless-stopped",
	}

	// Inside Docker, services talk to each other via service name not 127.0.0.1
	// Override DATABASE_URL so drizzle-kit and the app can connect to the DB container
	if opts.HasDB() && envVars != nil {
		if url, ok := envVars["DATABASE_URL"]; ok {
			dockerURL := url
			dockerURL = strings.ReplaceAll(dockerURL, "127.0.0.1", opts.DBService())
			dockerURL = strings.ReplaceAll(dockerURL, "localhost", opts.DBService())
			app.Environment = map[string]string{
				"DATABASE_URL": dockerURL,
			}
		}
	}

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
			Environment: buildServiceEnv(serviceName, def.Environment, envVars),
			Volumes:     def.Volumes,
			HealthCheck: def.HealthCheck,
			Restart:     def.Restart,
		}

		cf.AddService(serviceName, svc)

		if volName, ok := VolumeNames[serviceName]; ok {
			cf.AddVolume(volName)
		}
	}

	return cf
}

// buildServiceEnv builds the environment map for a service
// using real credentials from .env.local when available
func buildServiceEnv(serviceName string, defaults map[string]string, envVars map[string]string) map[string]string {
	if envVars == nil {
		return defaults
	}

	dbURL, hasURL := envVars["DATABASE_URL"]
	if !hasURL {
		return defaults
	}

	env := make(map[string]string)

	switch serviceName {
	case "mysql":
		_, password, _, _, dbname := ParseMySQLURL(dbURL)
		if dbname != "" {
			env["MYSQL_DATABASE"] = dbname
		} else {
			env["MYSQL_DATABASE"] = "myapp"
		}
		if password == "" {
			env["MYSQL_ALLOW_EMPTY_PASSWORD"] = "yes"
		} else {
			env["MYSQL_ROOT_PASSWORD"] = password
		}

	case "postgres":
		user, password, _, _, dbname := ParsePostgresURL(dbURL)
		if user != "" {
			env["POSTGRES_USER"] = user
		} else {
			env["POSTGRES_USER"] = "postgres"
		}
		if password != "" {
			env["POSTGRES_PASSWORD"] = password
		} else {
			env["POSTGRES_PASSWORD"] = "postgres"
		}
		if dbname != "" {
			env["POSTGRES_DB"] = dbname
		} else {
			env["POSTGRES_DB"] = "myapp"
		}

	default:
		return defaults
	}

	return env
}

func (c *ComposeFile) Write(path string) error {
	content := c.render()
	return os.WriteFile(path, []byte(content), 0644)
}

func (c *ComposeFile) render() string {
	var sb strings.Builder

	sb.WriteString("services:\n")

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
			envKeys := []string{}
			for k := range svc.Environment {
				envKeys = append(envKeys, k)
			}
			sortStrings(envKeys)
			for _, k := range envKeys {
				sb.WriteString(fmt.Sprintf("      %s: \"%s\"\n", k, svc.Environment[k]))
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

	if len(c.Volumes) > 0 {
		sb.WriteString("volumes:\n")
		for _, v := range c.Volumes {
			sb.WriteString(fmt.Sprintf("  %s:\n", v))
		}
	}

	return sb.String()
}

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