package vault

import "testing"

func TestEnvVault(t *testing.T) {
	t.Setenv("MINION_TEST_USERNAME", "testuser@example.com")

	v := NewEnvVault()
	val, err := v.Get("test", "username")
	if err != nil {
		t.Fatalf("Get error: %v", err)
	}
	if val != "testuser@example.com" {
		t.Errorf("expected 'testuser@example.com', got '%s'", val)
	}
}

func TestEnvVault_Missing(t *testing.T) {
	v := NewEnvVault()
	_, err := v.Get("no_such_institution", "password")
	if err == nil {
		t.Error("expected error for missing credential, got nil")
	}
}

func TestVaultInterface(t *testing.T) {
	// Verify EnvVault satisfies the Vault interface.
	var _ Vault = NewEnvVault()
	var _ Vault = NewKeyringVault()
}
