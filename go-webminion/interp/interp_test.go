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

func TestVars_ResolveWith_OverrideTakesPrecedence(t *testing.T) {
	v := New("test_flow", time.Now(), time.Now())
	if err := v.Set("key", "from_runtime"); err != nil {
		t.Fatal(err)
	}
	result, err := v.ResolveWith("{{key}}", map[string]string{"key": "from_override"})
	if err != nil {
		t.Fatalf("ResolveWith error: %v", err)
	}
	if result != "from_override" {
		t.Errorf("expected override to win, got %q", result)
	}
	// Confirm runtime was not mutated.
	runtime, _ := v.Resolve("{{key}}")
	if runtime != "from_runtime" {
		t.Errorf("ResolveWith should not mutate runtime, got %q", runtime)
	}
}

func TestVars_ResolveWith_FallsBackToRuntime(t *testing.T) {
	v := New("test_flow", time.Now(), time.Now())
	if err := v.Set("key", "runtime_val"); err != nil {
		t.Fatal(err)
	}
	result, err := v.ResolveWith("{{key}}", map[string]string{"other": "x"})
	if err != nil {
		t.Fatalf("ResolveWith error: %v", err)
	}
	if result != "runtime_val" {
		t.Errorf("expected fallback to runtime, got %q", result)
	}
}

func TestVars_ResolveWith_NilOverrides(t *testing.T) {
	v := New("test_flow", time.Now(), time.Now())
	if err := v.Set("key", "val"); err != nil {
		t.Fatal(err)
	}
	result, err := v.ResolveWith("{{key}}", nil)
	if err != nil {
		t.Fatalf("ResolveWith(nil) error: %v", err)
	}
	if result != "val" {
		t.Errorf("expected 'val', got %q", result)
	}
}
