package interp

import (
	"testing"
	"time"
)

func TestVars_Resolve(t *testing.T) {
	v := New("test_flow", time.Now(), time.Now())

	if err := v.Set("my_var", "hello"); err != nil {
		t.Fatal(err)
	}

	result, err := v.Resolve("{{my_var}} world")
	if err != nil {
		t.Fatalf("Resolve error: %v", err)
	}
	if result != "hello world" {
		t.Errorf("expected 'hello world', got '%s'", result)
	}
}

func TestVars_BuiltInVars(t *testing.T) {
	now := time.Now()
	v := New("test_flow", now, now)

	result, err := v.Resolve("{{today}}")
	if err != nil {
		t.Fatalf("Resolve error: %v", err)
	}
	if result != now.Format("2006-01-02") {
		t.Errorf("expected '%s', got '%s'", now.Format("2006-01-02"), result)
	}
}

func TestVars_CredentialKeysRejected(t *testing.T) {
	v := New("test_flow", time.Now(), time.Now())
	if err := v.Set("username", "foo"); err == nil {
		t.Error("expected error when setting credential key, got nil")
	}
}

func TestVars_CredentialPlaceholdersPassThrough(t *testing.T) {
	v := New("test_flow", time.Now(), time.Now())
	result, err := v.Resolve("{{username}}")
	if err != nil {
		t.Fatalf("Resolve error: %v", err)
	}
	if result != "{{username}}" {
		t.Errorf("expected credential placeholder to pass through, got '%s'", result)
	}
}

func TestVars_UnknownVarReturnsError(t *testing.T) {
	v := New("test_flow", time.Now(), time.Now())
	_, err := v.Resolve("{{no_such_var}}")
	if err == nil {
		t.Error("expected error for unknown variable, got nil")
	}
}
