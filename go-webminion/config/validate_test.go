package config

import "testing"

func TestValidateConfig_Valid(t *testing.T) {
	inst := &Institution{
		ID: "test",
		Flow: Flow{
			Actions: []Action{
				{Key: "start", Starting: true, Steps: []Step{{Method: "go"}}},
			},
		},
	}
	if errs := ValidateConfig(inst); len(errs) != 0 {
		t.Errorf("expected no errors, got %v", errs)
	}
}

func TestValidateConfig_NoStarting(t *testing.T) {
	inst := &Institution{
		ID:   "test",
		Flow: Flow{Actions: []Action{{Key: "no_start"}}},
	}
	errs := ValidateConfig(inst)
	if len(errs) != 1 {
		t.Errorf("expected 1 error (no starting action), got %d: %v", len(errs), errs)
	}
}

func TestValidateConfig_MultipleStarting(t *testing.T) {
	inst := &Institution{
		ID: "test",
		Flow: Flow{
			Actions: []Action{
				{Key: "a", Starting: true},
				{Key: "b", Starting: true},
			},
		},
	}
	errs := ValidateConfig(inst)
	if len(errs) != 1 {
		t.Errorf("expected 1 error (multiple starting), got %d: %v", len(errs), errs)
	}
}

func TestValidateConfig_DuplicateKeys(t *testing.T) {
	inst := &Institution{
		ID: "test",
		Flow: Flow{
			Actions: []Action{
				{Key: "dup", Starting: true},
				{Key: "dup"},
			},
		},
	}
	errs := ValidateConfig(inst)
	if len(errs) != 1 {
		t.Errorf("expected 1 error (duplicate key), got %d: %v", len(errs), errs)
	}
}

func TestValidateConfig_NoActions(t *testing.T) {
	inst := &Institution{ID: "test"}
	errs := ValidateConfig(inst)
	if len(errs) == 0 {
		t.Error("expected error for empty actions, got none")
	}
}
