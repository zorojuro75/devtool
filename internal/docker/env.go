package docker

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// AppendEnvExample appends service connection strings to .env.example
func AppendEnvExample(dir string, opts *DockerOptions) ([]EnvEntry, error) {
	envPath := filepath.Join(dir, ".env.example")

	existing := ""
	data, err := os.ReadFile(envPath)
	if err == nil {
		existing = string(data)
	}

	var added []EnvEntry
	var sb strings.Builder
	sb.WriteString(existing)

	if len(existing) > 0 && !strings.HasSuffix(existing, "\n") {
		sb.WriteString("\n")
	}

	for _, serviceName := range opts.Services {
		entries, ok := EnvEntries[serviceName]
		if !ok {
			continue
		}
		for _, entry := range entries {
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

	if err := os.WriteFile(envPath, []byte(sb.String()), 0644); err != nil {
		return nil, fmt.Errorf("cannot write .env.example: %w", err)
	}

	return added, nil
}

// AppendEnvEntries appends a specific list of EnvEntry to .env.example
func AppendEnvEntries(dir string, entries []EnvEntry) error {
	envPath := filepath.Join(dir, ".env.example")

	existing := ""
	data, err := os.ReadFile(envPath)
	if err == nil {
		existing = string(data)
	}

	var sb strings.Builder
	sb.WriteString(existing)

	if len(existing) > 0 && !strings.HasSuffix(existing, "\n") {
		sb.WriteString("\n")
	}

	for _, entry := range entries {
		if strings.Contains(existing, entry.Key+"=") ||
			strings.Contains(existing, entry.Key+"=\"") {
			continue
		}
		if entry.Comment != "" {
			sb.WriteString(fmt.Sprintf("\n# %s\n", entry.Comment))
		}
		sb.WriteString(fmt.Sprintf("%s=\"%s\"\n", entry.Key, entry.Value))
	}

	return os.WriteFile(envPath, []byte(sb.String()), 0644)
}

// PrintConnectionStrings prints the added env entries to the terminal
func PrintConnectionStrings(entries []EnvEntry) {
	if len(entries) == 0 {
		return
	}
	fmt.Println("\n  Connection strings added to .env.example:")
	for _, e := range entries {
		fmt.Printf("  %s=\"%s\"\n", e.Key, e.Value)
	}
	fmt.Println()
}

// ReadEnvLocal reads .env.local and returns a map of key=value pairs
func ReadEnvLocal(dir string) map[string]string {
	envPath := filepath.Join(dir, ".env.local")
	data, err := os.ReadFile(envPath)
	if err != nil {
		return nil
	}

	result := make(map[string]string)
	for _, line := range strings.Split(string(data), "\n") {
		line = strings.TrimSpace(line)
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		parts := strings.SplitN(line, "=", 2)
		if len(parts) != 2 {
			continue
		}
		key := strings.TrimSpace(parts[0])
		val := strings.Trim(strings.TrimSpace(parts[1]), `"'`)
		result[key] = val
	}
	return result
}

// ParseMySQLURL extracts credentials from mysql://user:password@host:port/dbname
func ParseMySQLURL(url string) (user, password, host, port, dbname string) {
	url = strings.TrimPrefix(url, "mysql://")

	atIdx := strings.LastIndex(url, "@")
	if atIdx == -1 {
		return
	}

	userInfo := url[:atIdx]
	hostInfo := url[atIdx+1:]

	if colonIdx := strings.Index(userInfo, ":"); colonIdx != -1 {
		user = userInfo[:colonIdx]
		password = userInfo[colonIdx+1:]
	} else {
		user = userInfo
	}

	if slashIdx := strings.Index(hostInfo, "/"); slashIdx != -1 {
		dbname = hostInfo[slashIdx+1:]
		hostPort := hostInfo[:slashIdx]
		if colonIdx := strings.LastIndex(hostPort, ":"); colonIdx != -1 {
			host = hostPort[:colonIdx]
			port = hostPort[colonIdx+1:]
		} else {
			host = hostPort
		}
	}
	return
}

// ParsePostgresURL extracts credentials from postgresql://user:password@host:port/dbname
func ParsePostgresURL(url string) (user, password, host, port, dbname string) {
	url = strings.TrimPrefix(url, "postgresql://")
	url = strings.TrimPrefix(url, "postgres://")

	atIdx := strings.LastIndex(url, "@")
	if atIdx == -1 {
		return
	}

	userInfo := url[:atIdx]
	hostInfo := url[atIdx+1:]

	if colonIdx := strings.Index(userInfo, ":"); colonIdx != -1 {
		user = userInfo[:colonIdx]
		password = userInfo[colonIdx+1:]
	} else {
		user = userInfo
	}

	if slashIdx := strings.Index(hostInfo, "/"); slashIdx != -1 {
		dbname = hostInfo[slashIdx+1:]
		hostPort := hostInfo[:slashIdx]
		if colonIdx := strings.LastIndex(hostPort, ":"); colonIdx != -1 {
			host = hostPort[:colonIdx]
			port = hostPort[colonIdx+1:]
		} else {
			host = hostPort
		}
	}
	return
}