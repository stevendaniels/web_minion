package config

import "fmt"

// ValidateConfig checks a parsed Institution for logical errors.
func ValidateConfig(inst *Institution) []error {
	var errs []error

	startingCount := 0
	seen := make(map[string]bool)
	for _, action := range inst.Flow.Actions {
		if action.Starting {
			startingCount++
		}
		if seen[action.Key] {
			errs = append(errs, fmt.Errorf("duplicate action key: %s", action.Key))
		}
		seen[action.Key] = true
	}

	if startingCount == 0 {
		errs = append(errs, fmt.Errorf("no starting action defined"))
	} else if startingCount > 1 {
		errs = append(errs, fmt.Errorf("multiple starting actions defined"))
	}

	return errs
}
