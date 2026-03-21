package vault

// Vault resolves credential variables at runtime.
type Vault interface {
	Get(institutionID, key string) (string, error)
}
