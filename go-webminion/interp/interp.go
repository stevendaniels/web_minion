package interp

import (
	"fmt"
	"regexp"
	"strings"
	"time"
)

var placeholder = regexp.MustCompile(`\{\{(\w+)\}\}`)

// Vars holds runtime variables for interpolation.
type Vars struct {
	builtIns map[string]string
	runtime  map[string]string
}

// New creates a Vars instance with built-in variables pre-populated.
func New(flowName string, startDate, endDate time.Time) *Vars {
	v := &Vars{
		builtIns: map[string]string{
			"flow_name":  flowName,
			"today":      time.Now().Format("2006-01-02"),
			"start_date": startDate.Format("2006-01-02"),
			"end_date":   endDate.Format("2006-01-02"),
		},
		runtime: make(map[string]string),
	}
	return v
}

// Set stores a runtime variable. Credential keys are rejected.
func (v *Vars) Set(key, value string) error {
	switch key {
	case "username", "password", "totp":
		return fmt.Errorf("cannot store credential %q in Vars; use Vault instead", key)
	}
	v.runtime[key] = value
	return nil
}

// Resolve replaces all {{key}} occurrences in s with their values.
// Returns an error if an unknown (non-credential) variable is referenced.
func (v *Vars) Resolve(s string) (string, error) {
	var resolveErr error
	result := placeholder.ReplaceAllStringFunc(s, func(match string) string {
		if resolveErr != nil {
			return match
		}
		key := strings.TrimPrefix(strings.TrimSuffix(match, "}}"), "{{")

		if val, ok := v.runtime[key]; ok {
			return val
		}
		if val, ok := v.builtIns[key]; ok {
			return val
		}
		// Credential placeholders are left for the vault to handle.
		if key == "username" || key == "password" || key == "totp" {
			return match
		}
		resolveErr = fmt.Errorf("unknown variable: %s", key)
		return match
	})
	return result, resolveErr
}
