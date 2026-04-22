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

const validFlowYAML = `id: test
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

func TestCommandValidate(t *testing.T) {
	p := writeConfig(t, validFlowYAML)
	if err := run([]string{"validate", "--config", p}); err != nil {
		t.Fatalf("validate error: %v", err)
	}
}

func TestCommandValidate_BadConfig(t *testing.T) {
	p := writeConfig(t, `id: test
flow:
  actions: []
`)
	if err := run([]string{"validate", "--config", p}); err == nil {
		t.Error("expected error for invalid config (no starting action), got nil")
	}
}

func TestCommandUnknown(t *testing.T) {
	if err := run([]string{"foobar"}); err == nil {
		t.Error("expected error for unknown subcommand, got nil")
	}
}

func TestCommandMissingSubcommand(t *testing.T) {
	if err := run([]string{}); err == nil {
		t.Error("expected error when subcommand is missing, got nil")
	}
}

func TestCommandHelp(t *testing.T) {
	p := writeConfig(t, `id: test
vars:
  url: ""
flow:
  actions:
    - key: start
      starting: true
      steps:
        - method: go
          value: "{{url}}"
        - method: body_includes
          is_validator: true
`)
	if err := run([]string{"help", "--config", p}); err != nil {
		t.Fatalf("help error: %v", err)
	}
}

func TestCommandRun_EnvFile(t *testing.T) {
	envPath := filepath.Join(t.TempDir(), ".env")
	if err := os.WriteFile(envPath, []byte("MY_VAR=hello\n# comment\n\nOTHER=world\n"), 0644); err != nil {
		t.Fatal(err)
	}
	pairs, err := parseEnvFile(envPath)
	if err != nil {
		t.Fatalf("parseEnvFile error: %v", err)
	}
	if pairs["MY_VAR"] != "hello" {
		t.Errorf("expected MY_VAR=hello, got %q", pairs["MY_VAR"])
	}
	if pairs["OTHER"] != "world" {
		t.Errorf("expected OTHER=world, got %q", pairs["OTHER"])
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
