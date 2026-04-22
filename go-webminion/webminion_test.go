package webminion

import (
	"os"
	"path/filepath"
	"testing"
)

const validYAML = `id: test
name: Test
flow:
  actions:
    - key: start
      starting: true
      steps:
        - method: go
          value: https://example.com
        - method: body_includes
          is_validator: true
`

func writeConfig(t *testing.T, content string) string {
	t.Helper()
	path := filepath.Join(t.TempDir(), "test.yaml")
	if err := os.WriteFile(path, []byte(content), 0644); err != nil {
		t.Fatal(err)
	}
	return path
}

func TestNewFlow(t *testing.T) {
	flow, err := NewFlow(writeConfig(t, validYAML))
	if err != nil {
		t.Fatalf("NewFlow error: %v", err)
	}
	if flow == nil {
		t.Error("NewFlow returned nil")
	}
}

func TestNewFlow_InvalidConfig(t *testing.T) {
	yaml := `id: test
flow:
  actions: []
`
	_, err := NewFlow(writeConfig(t, yaml))
	if err == nil {
		t.Error("expected error for invalid config (no starting action), got nil")
	}
}

func TestNewFlow_MissingFile(t *testing.T) {
	_, err := NewFlow("/nonexistent/path/config.yaml")
	if err == nil {
		t.Error("expected error for missing file, got nil")
	}
}
