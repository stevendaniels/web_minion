package config

import "testing"

func TestValidateConfig_Valid(t *testing.T) {
	cfg := &Config{
		ID: "test",
		Flow: Flow{
			Actions: []Action{
				{Key: "start", Starting: true, Steps: []Step{{Method: "go"}}},
			},
		},
	}
	if errs := ValidateConfig(cfg); len(errs) != 0 {
		t.Errorf("expected no errors, got %v", errs)
	}
}

func TestValidateConfig_NoStarting(t *testing.T) {
	cfg := &Config{
		ID:   "test",
		Flow: Flow{Actions: []Action{{Key: "no_start"}}},
	}
	errs := ValidateConfig(cfg)
	if len(errs) != 1 {
		t.Errorf("expected 1 error (no starting action), got %d: %v", len(errs), errs)
	}
}

func TestValidateConfig_MultipleStarting(t *testing.T) {
	cfg := &Config{
		ID: "test",
		Flow: Flow{
			Actions: []Action{
				{Key: "a", Starting: true},
				{Key: "b", Starting: true},
			},
		},
	}
	errs := ValidateConfig(cfg)
	if len(errs) != 1 {
		t.Errorf("expected 1 error (multiple starting), got %d: %v", len(errs), errs)
	}
}

func TestValidateConfig_DuplicateKeys(t *testing.T) {
	cfg := &Config{
		ID: "test",
		Flow: Flow{
			Actions: []Action{
				{Key: "dup", Starting: true},
				{Key: "dup"},
			},
		},
	}
	errs := ValidateConfig(cfg)
	if len(errs) != 1 {
		t.Errorf("expected 1 error (duplicate key), got %d: %v", len(errs), errs)
	}
}

func TestValidateConfig_NoActions(t *testing.T) {
	cfg := &Config{ID: "test"}
	errs := ValidateConfig(cfg)
	if len(errs) == 0 {
		t.Error("expected error for empty actions, got none")
	}
}

func TestValidateStep_WriteFile_MissingPattern(t *testing.T) {
	cfg := &Config{
		ID: "test",
		Flow: Flow{Actions: []Action{{
			Key:      "start",
			Starting: true,
			Steps:    []Step{{Method: "write_file", Value: "/tmp/reply.txt"}},
		}}},
	}
	errs := ValidateConfig(cfg)
	if len(errs) == 0 {
		t.Error("expected error for write_file missing pattern")
	}
}

func TestValidateStep_WriteFile_MissingValue(t *testing.T) {
	cfg := &Config{
		ID: "test",
		Flow: Flow{Actions: []Action{{
			Key:      "start",
			Starting: true,
			Steps:    []Step{{Method: "write_file", Pattern: "/dispatches/out"}},
		}}},
	}
	errs := ValidateConfig(cfg)
	if len(errs) == 0 {
		t.Error("expected error for write_file missing value")
	}
}

func TestValidateStep_WaitForReply_MissingPattern(t *testing.T) {
	cfg := &Config{
		ID: "test",
		Flow: Flow{Actions: []Action{{
			Key:      "start",
			Starting: true,
			Steps:    []Step{{Method: "wait_for_reply", Value: "mfa_code"}},
		}}},
	}
	errs := ValidateConfig(cfg)
	if len(errs) == 0 {
		t.Error("expected error for wait_for_reply missing pattern")
	}
}

func TestValidateStep_WaitForReply_MissingValue(t *testing.T) {
	cfg := &Config{
		ID: "test",
		Flow: Flow{Actions: []Action{{
			Key:      "start",
			Starting: true,
			Steps:    []Step{{Method: "wait_for_reply", Pattern: "/dispatches/web-minion-in/reply.txt"}},
		}}},
	}
	errs := ValidateConfig(cfg)
	if len(errs) == 0 {
		t.Error("expected error for wait_for_reply missing value")
	}
}

func TestValidateStep_WaitForReply_NegativeTimeout(t *testing.T) {
	cfg := &Config{
		ID: "test",
		Flow: Flow{Actions: []Action{{
			Key:      "start",
			Starting: true,
			Steps:    []Step{{Method: "wait_for_reply", Pattern: "/tmp/r.txt", Value: "v", Timeout: -1}},
		}}},
	}
	errs := ValidateConfig(cfg)
	if len(errs) == 0 {
		t.Error("expected error for negative timeout")
	}
}

func TestValidateStep_WriteFile_Valid(t *testing.T) {
	cfg := &Config{
		ID: "test",
		Flow: Flow{Actions: []Action{{
			Key:      "start",
			Starting: true,
			Steps:    []Step{{Method: "write_file", Pattern: "/downloads/out.md", Value: "{{page_html}}"}},
		}}},
	}
	errs := ValidateConfig(cfg)
	if len(errs) != 0 {
		t.Errorf("expected no errors for valid write_file, got: %v", errs)
	}
}

func TestValidateStep_HTMLToMarkdown_MissingValue(t *testing.T) {
	cfg := &Config{
		ID: "test",
		Flow: Flow{Actions: []Action{{
			Key:      "start",
			Starting: true,
			Steps:    []Step{{Method: "html_to_markdown"}},
		}}},
	}
	errs := ValidateConfig(cfg)
	if len(errs) == 0 {
		t.Error("expected error for html_to_markdown missing value")
	}
}
