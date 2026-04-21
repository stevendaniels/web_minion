package config

import "testing"

func TestConfigTypes(t *testing.T) {
	cfg := &Config{
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
	if cfg.ID != "test" {
		t.Errorf("expected ID 'test', got '%s'", cfg.ID)
	}
	if len(cfg.Flow.Actions) != 1 {
		t.Errorf("expected 1 action, got %d", len(cfg.Flow.Actions))
	}
}
