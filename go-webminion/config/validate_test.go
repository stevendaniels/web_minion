package config

import "testing"

func TestValidateConfig_Valid(t *testing.T) {
	cfg := &Config{
		ID: "test",
		Flow: Flow{
			Actions: []Action{
				{Key: "start", Starting: true, Steps: []Step{{Method: "go"}, {Method: "body_includes", IsValidator: true}}},
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
	if len(errs) == 0 {
		t.Error("expected error for no starting action")
	}
}

func TestValidateConfig_MultipleStarting(t *testing.T) {
	cfg := &Config{
		ID: "test",
		Flow: Flow{
			Actions: []Action{
				{Key: "a", Starting: true, Steps: []Step{{Method: "go"}, {Method: "body_includes", IsValidator: true}}},
				{Key: "b", Starting: true, Steps: []Step{{Method: "go"}, {Method: "body_includes", IsValidator: true}}},
			},
		},
	}
	errs := ValidateConfig(cfg)
	if len(errs) == 0 {
		t.Error("expected error for multiple starting actions")
	}
}

func TestValidateConfig_DuplicateKeys(t *testing.T) {
	cfg := &Config{
		ID: "test",
		Flow: Flow{
			Actions: []Action{
				{Key: "start", Starting: true, Steps: []Step{{Method: "go"}, {Method: "body_includes", IsValidator: true}}},
				{Key: "start", Steps: []Step{{Method: "go"}, {Method: "body_includes", IsValidator: true}}},
			},
		},
	}
	errs := ValidateConfig(cfg)
	if len(errs) == 0 {
		t.Error("expected error for duplicate action keys")
	}
}

func TestValidateConfig_NoActions(t *testing.T) {
	cfg := &Config{ID: "test"}
	errs := ValidateConfig(cfg)
	if len(errs) == 0 {
		t.Error("expected error for empty actions, got none")
	}
}

func TestValidateConfig_ActionMissingValidation(t *testing.T) {
	inst := &Config{
		ID: "test",
		Flow: Flow{
			Actions: []Action{{
				Key:      "start",
				Starting: true,
				Steps:    []Step{{Method: "go"}, {Method: "click"}},
			}},
		},
	}
	errs := ValidateConfig(inst)
	if len(errs) == 0 {
		t.Error("expected error for action missing validation step")
	}
}

func TestValidateConfig_ActionNoSteps(t *testing.T) {
	inst := &Config{
		ID: "test",
		Flow: Flow{
			Actions: []Action{{
				Key:      "start",
				Starting: true,
				Steps:    []Step{},
			}},
		},
	}
	errs := ValidateConfig(inst)
	if len(errs) == 0 {
		t.Error("expected error for action with no steps")
	}
}

func TestValidateConfig_DuplicateActionKeysError(t *testing.T) {
	inst := &Config{
		ID: "test",
		Flow: Flow{
			Actions: []Action{
				{Key: "dup", Starting: true, Steps: []Step{{Method: "go"}, {Method: "body_includes", IsValidator: true}}},
				{Key: "dup", Steps: []Step{{Method: "go"}, {Method: "body_includes", IsValidator: true}}},
			},
		},
	}
	errs := ValidateConfig(inst)
	if len(errs) == 0 {
		t.Error("expected error for duplicate action keys")
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
			Steps:    []Step{{Method: "write_file", Pattern: "/downloads/out.md", Value: "hello world"}, {Method: "body_includes", IsValidator: true}},
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

func TestValidateConfig_ValidMultipleActions(t *testing.T) {
	inst := &Config{
		ID: "test",
		Flow: Flow{
			Actions: []Action{
				{Key: "start", Starting: true, Steps: []Step{{Method: "go"}, {Method: "body_includes", IsValidator: true}}},
				{Key: "next", OnSuccess: "complete", Steps: []Step{{Method: "get_field"}, {Method: "value_equals", IsValidator: true}}},
				{Key: "complete", Steps: []Step{{Method: "save_value"}, {Method: "body_includes", IsValidator: true}}},
			},
		},
	}
	errs := ValidateConfig(inst)
	if len(errs) != 0 {
		t.Errorf("expected no errors for valid multi-action flow, got: %v", errs)
	}
}

func TestValidateStep_InvalidMethod(t *testing.T) {
	inst := &Config{
		ID: "test",
		Flow: Flow{Actions: []Action{{
			Key:      "start",
			Starting: true,
			Steps:    []Step{{Method: "invalid_method"}, {Method: "body_includes", IsValidator: true}},
		}}},
	}
	errs := ValidateConfig(inst)
	if len(errs) == 0 {
		t.Error("expected error for invalid method")
	}
}

func TestValidateConfig_CycleDetected(t *testing.T) {
	inst := &Config{
		ID: "test",
		Flow: Flow{
			Actions: []Action{
				{Key: "start", Starting: true, OnSuccess: "action_two", Steps: []Step{{Method: "go"}, {Method: "body_includes", IsValidator: true}}},
				{Key: "action_two", OnSuccess: "action_three", Steps: []Step{{Method: "body_includes", IsValidator: true}}},
				{Key: "action_three", OnSuccess: "start", Steps: []Step{{Method: "body_includes", IsValidator: true}}},
			},
		},
	}
	errs := ValidateConfig(inst)
	if len(errs) == 0 {
		t.Error("expected error for cyclical flow")
	}
}

func TestValidateStep_MissingVariable(t *testing.T) {
	inst := &Config{
		ID: "test",
		Flow: Flow{Actions: []Action{{
			Key:      "start",
			Starting: true,
			Steps:    []Step{{Method: "go", Value: "{{missing_var}}"}, {Method: "body_includes", IsValidator: true}},
		}}},
	}
	errs := ValidateConfig(inst)
	if len(errs) == 0 {
		t.Error("expected error for missing variable")
	}
}

func TestValidateStep_ValidVariable(t *testing.T) {
	inst := &Config{
		ID: "test",
		Flow: Flow{Actions: []Action{{
			Key:      "start",
			Starting: true,
			Steps:    []Step{{Method: "go", Value: "{{my_url}}"}, {Method: "body_includes", IsValidator: true}},
		}}},
		Vars: map[string]string{"my_url": "https://example.com"},
	}
	errs := ValidateConfig(inst)
	if len(errs) != 0 {
		t.Errorf("expected no errors for valid variable, got: %v", errs)
	}
}

func TestValidateConfig_InvalidActionReferences(t *testing.T) {
	inst := &Config{
		ID: "test",
		Flow: Flow{Actions: []Action{
			{Key: "start", Starting: true, OnSuccess: "nonexistent", Steps: []Step{{Method: "go"}, {Method: "body_includes", IsValidator: true}}},
			{Key: "other", OnFailure: "also_missing", Steps: []Step{{Method: "go"}, {Method: "body_includes", IsValidator: true}}},
		}},
	}
	errs := ValidateConfig(inst)
	if len(errs) < 2 {
		t.Errorf("expected at least 2 errors for invalid action references, got %d: %v", len(errs), errs)
	}
}

func TestValidateConfig_NoCycle(t *testing.T) {
	inst := &Config{
		ID: "test",
		Flow: Flow{
			Actions: []Action{
				{Key: "start", Starting: true, OnSuccess: "action_two", Steps: []Step{{Method: "go"}, {Method: "body_includes", IsValidator: true}}},
				{Key: "action_two", Steps: []Step{{Method: "body_includes", IsValidator: true}}},
			},
		},
	}
	errs := ValidateConfig(inst)
	if len(errs) != 0 {
		t.Errorf("expected no errors for non-cyclical flow, got: %v", errs)
	}
}
