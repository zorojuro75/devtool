package config

import (
	"os"
	"path/filepath"
	"testing"
)

func TestLoad(t *testing.T) {
	tests := []struct {
		name      string
		yaml      string
		envKey    string
		wantKey   string
		wantModel string
		wantErr   bool
	}{
		{
			name: "valid config",
			yaml: `provider: openrouter
api_key: test-key-123
timeout: 30
model: mistralai/mistral-7b-instruct:free`,
			wantKey:   "test-key-123",
			wantModel: "mistralai/mistral-7b-instruct:free",
		},
		{
			name: "env var overrides config key",
			yaml: `provider: openrouter
api_key: from-file
timeout: 30
model: some-model`,
			envKey:  "from-env",
			wantKey: "from-env",
		},
		{
			name:    "invalid yaml",
			yaml:    "this: is: not: valid: yaml:::",
			wantErr: true,
		},
		{
			name: "missing api key",
			yaml: `provider: openrouter
timeout: 30
model: some-model`,
			wantKey: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Write temp config file
			dir := t.TempDir()
			path := filepath.Join(dir, ".devtool.yaml")
			if err := os.WriteFile(path, []byte(tt.yaml), 0600); err != nil {
				t.Fatal(err)
			}

			// Set env var if needed
			if tt.envKey != "" {
				t.Setenv("DEVTOOL_API_KEY", tt.envKey)
			}

			cfg, err := Load(path)
			if tt.wantErr {
				if err == nil {
					t.Error("expected error, got nil")
				}
				return
			}
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if cfg.APIKey != tt.wantKey {
				t.Errorf("APIKey = %q, want %q", cfg.APIKey, tt.wantKey)
			}
			if tt.wantModel != "" && cfg.Model != tt.wantModel {
				t.Errorf("Model = %q, want %q", cfg.Model, tt.wantModel)
			}
		})
	}
}

func TestLoadDefaultTimeout(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, ".devtool.yaml")
	os.WriteFile(path, []byte("provider: openrouter\n"), 0600)

	cfg, err := Load(path)
	if err != nil {
		t.Fatal(err)
	}
	if cfg.Timeout != 30 {
		t.Errorf("default timeout = %d, want 30", cfg.Timeout)
	}
}