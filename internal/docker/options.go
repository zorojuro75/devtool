package docker

type DockerOptions struct {
	ProjectName string
	Stack       string   // "nextjs"
	Port        int      // default 3000
	Services    []string // ["postgres", "redis"]
	NodeVersion string   // "20"
}

func NewDockerOptions() *DockerOptions {
	return &DockerOptions{
		Stack:       "nextjs",
		Port:        3000,
		NodeVersion: "20",
		Services:    []string{},
	}
}

func (o *DockerOptions) HasPostgres() bool {
	return o.hasService("postgres")
}

func (o *DockerOptions) HasMySQL() bool {
	return o.hasService("mysql")
}

func (o *DockerOptions) HasRedis() bool {
	return o.hasService("redis")
}

func (o *DockerOptions) HasMailhog() bool {
	return o.hasService("mailhog")
}

func (o *DockerOptions) HasDB() bool {
	return o.HasPostgres() || o.HasMySQL()
}

func (o *DockerOptions) HasServices() bool {
	return len(o.Services) > 0
}

// DBService returns the first selected DB service name
// Used for depends_on in the app service
func (o *DockerOptions) DBService() string {
	if o.HasPostgres() {
		return "postgres"
	}
	if o.HasMySQL() {
		return "mysql"
	}
	return ""
}

func (o *DockerOptions) hasService(name string) bool {
	for _, s := range o.Services {
		if s == name {
			return true
		}
	}
	return false
}