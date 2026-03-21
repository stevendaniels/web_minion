package vault

import (
	"fmt"

	"github.com/zalando/go-keyring"
)

// KeyringVault resolves credentials from the OS keyring.
// Secrets are stored under service "minion" with key "<institutionID>/<key>".
type KeyringVault struct{}

func NewKeyringVault() *KeyringVault {
	return &KeyringVault{}
}

func (v *KeyringVault) Get(institutionID, key string) (string, error) {
	val, err := keyring.Get("minion", fmt.Sprintf("%s/%s", institutionID, key))
	if err != nil {
		return "", fmt.Errorf("credential not found in keyring for %s/%s: %w", institutionID, key, err)
	}
	return val, nil
}
