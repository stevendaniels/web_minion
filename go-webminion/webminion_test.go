package webminion

import (
	"os"
	"path/filepath"
	"testing"
)

func TestLoadConfig(t *testing.T) {
	tmpDir := t.TempDir()
	configPath := filepath.Join(tmpDir, "test.yaml")

	content := `id: test
name: Test
flow:
  actions:
    - key: start
      starting: true
      steps:
        - method: go
          value: https://example.com
`
	if err := os.WriteFile(configPath, []byte(content), 0644); err != nil {
		t.Fatal(err)
	}

	inst, err := LoadConfig(configPath)
	if err != nil {
		t.Fatalf("LoadConfig error: %v", err)
	}
	if inst.ID != "test" {
		t.Errorf("expected ID 'test', got '%s'", inst.ID)
	}
}

func TestNewFlow(t *testing.T) {
	tmpDir := t.TempDir()
	configPath := filepath.Join(tmpDir, "test.yaml")

	content := `id: test
flow:
  actions:
    - key: start
      starting: true
      steps: []
`
	if err := os.WriteFile(configPath, []byte(content), 0644); err != nil {
		t.Fatal(err)
	}

	inst, err := LoadConfig(configPath)
	if err != nil {
		t.Fatal(err)
	}

	flow := NewFlow(inst)
	if flow == nil {
		t.Error("NewFlow returned nil")
	}
}
