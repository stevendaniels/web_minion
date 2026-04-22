package main

import (
	"os"
	"path/filepath"
	"testing"
)

func writeConfig(t *testing.T, content string) string {
	t.Helper()
	tmpDir := t.TempDir()
	p := filepath.Join(tmpDir, "test.yaml")
	if err := os.WriteFile(p, []byte(content), 0644); err != nil {
		t.Fatal(err)
	}
	return p
}

func TestCommandValidate(t *testing.T) {
	p := writeConfig(t, `id: test
flow:
  actions:
    - key: start
      starting: true
      steps:
        - method: go
          value: https://example.com
        - method: body_includes
          is_validator: true
`)
	err := run([]string{"--command", "validate", "--config", p})
	if err != nil {
		t.Fatalf("validate error: %v", err)
	}
}

func TestCommandValidate_BadConfig(t *testing.T) {
	p := writeConfig(t, `id: test
flow:
  actions: []
`)
	err := run([]string{"--command", "validate", "--config", p})
	if err == nil {
		t.Error("expected error for invalid config (no starting action), got nil")
	}
}

func TestCommandUnknown(t *testing.T) {
	if err := run([]string{"--command", "foobar"}); err == nil {
		t.Error("expected error for unknown command, got nil")
	}
}

func TestCommandMissingCommand(t *testing.T) {
	if err := run([]string{}); err == nil {
		t.Error("expected error when --command is missing, got nil")
	}
}

func TestParseDateRange(t *testing.T) {
	start, end, err := parseDateRange("2024-01-01", "2024-12-31")
	if err != nil {
		t.Fatal(err)
	}
	if start.Year() != 2024 || start.Month() != 1 || start.Day() != 1 {
		t.Errorf("unexpected start date: %v", start)
	}
	if end.Year() != 2024 || end.Month() != 12 || end.Day() != 31 {
		t.Errorf("unexpected end date: %v", end)
	}
}

func TestParseDateRange_Invalid(t *testing.T) {
	if _, _, err := parseDateRange("not-a-date", ""); err == nil {
		t.Error("expected error for invalid date, got nil")
	}
}
