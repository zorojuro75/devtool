package docker

// ServiceDef defines a complete docker service configuration
type ServiceDef struct {
	Image       string
	Ports       []string
	Environment map[string]string
	Volumes     []string
	HealthCheck *HealthCheck
	Restart     string
}

type HealthCheck struct {
	Test     []string
	Interval string
	Timeout  string
	Retries  int
}

// EnvEntry is a connection string to append to .env.example
type EnvEntry struct {
	Comment string
	Key     string
	Value   string
}

// ServiceRegistry holds all supported service definitions
var ServiceRegistry = map[string]ServiceDef{
	"postgres": {
		Image: "postgres:16-alpine",
		Ports: []string{"5432:5432"},
		Environment: map[string]string{
			"POSTGRES_USER":     "postgres",
			"POSTGRES_PASSWORD": "postgres",
			"POSTGRES_DB":       "${POSTGRES_DB:-myapp}",
		},
		Volumes: []string{"pgdata:/var/lib/postgresql/data"},
		HealthCheck: &HealthCheck{
			Test:     []string{"CMD-SHELL", "pg_isready -U postgres"},
			Interval: "5s",
			Timeout:  "5s",
			Retries:  5,
		},
		Restart: "unless-stopped",
	},

	"mysql": {
		Image: "mysql:8",
		Ports: []string{"3306:3306"},
		Environment: map[string]string{
			"MYSQL_ROOT_PASSWORD": "${MYSQL_ROOT_PASSWORD:-password}",
			"MYSQL_DATABASE":      "${MYSQL_DATABASE:-myapp}",
		},
		Volumes: []string{"mysqldata:/var/lib/mysql"},
		HealthCheck: &HealthCheck{
			Test:     []string{"CMD", "mysqladmin", "ping", "-h", "localhost"},
			Interval: "5s",
			Timeout:  "5s",
			Retries:  5,
		},
		Restart: "unless-stopped",
	},

	"redis": {
		Image:   "redis:7-alpine",
		Ports:   []string{"6379:6379"},
		Volumes: []string{"redisdata:/data"},
		HealthCheck: &HealthCheck{
			Test:     []string{"CMD", "redis-cli", "ping"},
			Interval: "5s",
			Timeout:  "5s",
			Retries:  5,
		},
		Restart: "unless-stopped",
	},

	"mailhog": {
		Image:   "mailhog/mailhog",
		Ports:   []string{"1025:1025", "8025:8025"},
		Restart: "unless-stopped",
	},
}

// VolumeNames maps service names to their volume names
var VolumeNames = map[string]string{
	"postgres": "pgdata",
	"mysql":    "mysqldata",
	"redis":    "redisdata",
}

// EnvEntries maps service names to their .env.example entries
var EnvEntries = map[string][]EnvEntry{
	"postgres": {
		{
			Comment: "Docker — PostgreSQL",
			Key:     "DATABASE_URL",
			Value:   "postgresql://postgres:postgres@localhost:5432/myapp",
		},
	},
	"mysql": {
		{
			Comment: "Docker — MySQL",
			Key:     "DATABASE_URL",
			Value:   "mysql://root:password@127.0.0.1:3306/myapp",
		},
	},
	"redis": {
		{
			Comment: "Docker — Redis",
			Key:     "REDIS_URL",
			Value:   "redis://localhost:6379",
		},
	},
	"mailhog": {
		{
			Comment: "Docker — Mailhog SMTP",
			Key:     "SMTP_HOST",
			Value:   "localhost",
		},
		{
			Comment: "",
			Key:     "SMTP_PORT",
			Value:   "1025",
		},
		{
			Comment: "",
			Key:     "SMTP_FROM",
			Value:   "noreply@myapp.local",
		},
	},
}

// ServiceOrder defines the display order in prompts
var ServiceOrder = []string{"postgres", "mysql", "redis", "mailhog"}

// ServiceDisplayNames maps service keys to human-readable names
var ServiceDisplayNames = map[string]string{
	"postgres": "PostgreSQL 16",
	"mysql":    "MySQL 8",
	"redis":    "Redis 7",
	"mailhog":  "Mailhog (local email)",
}