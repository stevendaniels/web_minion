package config

import (
	"os"
	"path/filepath"
	"testing"
)

func TestLoadConfigYAML(t *testing.T) {
	tmpDir := t.TempDir()
	configPath := filepath.Join(tmpDir, "test.yaml")

	yamlContent := `id: test
name: Test
base_url: https://example.com
download_dir: /tmp/test
flow:
  actions:
    - key: test_action
      starting: true
      steps:
        - method: go
`
	if err := os.WriteFile(configPath, []byte(yamlContent), 0644); err != nil {
		t.Fatal(err)
	}

	cfg, err := LoadConfig(configPath)
	if err != nil {
		t.Fatalf("LoadConfig error: %v", err)
	}
	if cfg.ID != "test" {
		t.Errorf("expected ID 'test', got '%s'", cfg.ID)
	}
}

func TestLoadConfigJSON(t *testing.T) {
	tmpDir := t.TempDir()
	configPath := filepath.Join(tmpDir, "test.json")

	jsonContent := `{
		"id": "test",
		"name": "Test",
		"base_url": "https://example.com",
		"download_dir": "/tmp/test",
		"flow": {
			"actions": [{
				"key": "test_action",
				"starting": true,
				"steps": [{"method": "go"}]
			}]
		}
	}`
	if err := os.WriteFile(configPath, []byte(jsonContent), 0644); err != nil {
		t.Fatal(err)
	}

	cfg, err := LoadConfig(configPath)
	if err != nil {
		t.Fatalf("LoadConfig error: %v", err)
	}
	if cfg.ID != "test" {
		t.Errorf("expected ID 'test', got '%s'", cfg.ID)
	}
}

func TestLoadConfigMissingFile(t *testing.T) {
	_, err := LoadConfig("/nonexistent/path/config.yaml")
	if err == nil {
		t.Error("expected error for missing file, got nil")
	}
}
