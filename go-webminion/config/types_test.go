package config

import "testing"

func TestConfigTypes(t *testing.T) {
	inst := &Config{
		ID: "test",
		Flow: Flow{
			Actions: []Action{
				{
					Key:      "test_action",
					Starting: true,
					Steps:    []Step{{Method: "go"}},
				},
			},
		},
	}
	if inst.ID != "test" {
		t.Errorf("expected ID 'test', got '%s'", inst.ID)
	}
	if len(inst.Flow.Actions) != 1 {
		t.Errorf("expected 1 action, got %d", len(inst.Flow.Actions))
	}
}
