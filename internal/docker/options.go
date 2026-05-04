package docker

type DockerOptions struct {
	ProjectName string
	Stack       string   // "nextjs"
	ORM         string   // "prisma" | "drizzle" | "none"
	IsTypeScript bool
	Port        int      // default 3000
	Services    []string // ["postgres", "redis"]
	NodeVersion string   // "20"
}

func NewDockerOptions() *DockerOptions {
	return &DockerOptions{
		Stack:        "nextjs",
		ORM:          "none",
		IsTypeScript: true,
		Port:         3000,
		NodeVersion:  "20",
		Services:     []string{},
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

func (o *DockerOptions) IsPrisma() bool {
	return o.ORM == "prisma"
}

func (o *DockerOptions) IsDrizzle() bool {
	return o.ORM == "drizzle"
}

func (o *DockerOptions) HasORM() bool {
	return o.ORM != "none" && o.ORM != ""
}

// MigrateCommand returns the correct migration command for the detected ORM
func (o *DockerOptions) MigrateCommand() string {
	switch o.ORM {
	case "prisma":
		return "docker-compose exec app npx prisma migrate deploy"
	case "drizzle":
		return "docker-compose exec app sh -c \"npx drizzle-kit generate && npx drizzle-kit migrate\""
	default:
		return ""
	}
}

// DBService returns the first selected DB service name
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