package vault

import (
	"fmt"
	"os"
	"strings"
)

// EnvVault resolves credentials from environment variables.
// It looks for MINION_<INSTITUTION_ID>_<KEY> (all uppercase).
type EnvVault struct{}

func NewEnvVault() *EnvVault {
	return &EnvVault{}
}

func (v *EnvVault) Get(institutionID, key string) (string, error) {
	envKey := fmt.Sprintf("MINION_%s_%s", strings.ToUpper(institutionID), strings.ToUpper(key))
	val := os.Getenv(envKey)
	if val == "" {
		return "", fmt.Errorf("credential not found: %s", envKey)
	}
	return val, nil
}
