package config

import (
	"fmt"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

type Config struct {
	Provider string `yaml:"provider"`
	APIKey   string `yaml:"api_key"`
	Timeout  int    `yaml:"timeout"`
	Model    string `yaml:"model"`
}

const defaultConfig = `provider: openrouter
api_key: ""        # or set DEVTOOL_API_KEY env var
timeout: 30
model: "mistralai/mistral-7b-instruct:free"
`

func Load(path string) (*Config, error) {
	if path == "" {
		home, err := os.UserHomeDir()
		if err != nil {
			return nil, fmt.Errorf("cannot find home directory: %w", err)
		}
		path = filepath.Join(home, ".devtool.yaml")
	}

	// Create default config if it doesn't exist
	if _, err := os.Stat(path); os.IsNotExist(err) {
		if err := os.WriteFile(path, []byte(defaultConfig), 0600); err != nil {
			return nil, fmt.Errorf("could not create config file: %w", err)
		}
		fmt.Printf("Created default config at %s — add your OpenRouter API key to get started.\n\n", path)
	}

	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("cannot read config: %w", err)
	}

	var cfg Config
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return nil, fmt.Errorf("invalid config YAML: %w", err)
	}

	// Env var overrides config file
	if key := os.Getenv("DEVTOOL_API_KEY"); key != "" {
		cfg.APIKey = key
	}

	if cfg.Timeout == 0 {
		cfg.Timeout = 30
	}

	return &cfg, nil
}