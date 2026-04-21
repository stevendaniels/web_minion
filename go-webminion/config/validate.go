package config

import "fmt"

// ValidateConfig checks a parsed Config for logical errors.
func ValidateConfig(inst *Config) []error {
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

	for _, action := range inst.Flow.Actions {
		for i, step := range action.Steps {
			errs = append(errs, validateStep(action.Key, i, step)...)
		}
	}

	return errs
}

func validateStep(actionKey string, idx int, step Step) []error {
	var errs []error
	loc := fmt.Sprintf("action %q step %d (%s)", actionKey, idx, step.Method)
	switch step.Method {
	case "write_file":
		if step.Pattern == "" {
			errs = append(errs, fmt.Errorf("%s: pattern (file path) is required", loc))
		}
		if step.Value == "" {
			errs = append(errs, fmt.Errorf("%s: value (content) is required", loc))
		}
	case "html_to_markdown":
		if step.Value == "" {
			errs = append(errs, fmt.Errorf("%s: value (variable name) is required", loc))
		}
	case "wait_for_reply":
		if step.Pattern == "" {
			errs = append(errs, fmt.Errorf("%s: pattern (file path) is required", loc))
		}
		if step.Value == "" {
			errs = append(errs, fmt.Errorf("%s: value (variable name) is required", loc))
		}
		if step.Timeout < 0 {
			errs = append(errs, fmt.Errorf("%s: timeout must be >= 0", loc))
		}
	}
	return errs
}
