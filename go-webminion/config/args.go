package config

import "sort"

var builtInVars = map[string]bool{
	"flow_name": true, "today": true, "start_date": true, "end_date": true,
}

var credentialVars = map[string]bool{
	"username": true, "password": true, "totp": true,
}

// RequiredArgs returns the placeholder names that must be supplied at runtime.
// A placeholder is required when it is not a built-in, not a credential, and
// is not already satisfied by a non-empty entry in config.Vars.
func RequiredArgs(inst *Config) []string {
	satisfied := make(map[string]bool)
	for k, v := range inst.Vars {
		if v != "" {
			satisfied[k] = true
		}
	}

	seen := make(map[string]bool)
	for _, action := range inst.Flow.Actions {
		for _, step := range action.Steps {
			collectPlaceholders(step.Value, satisfied, seen)
			collectPlaceholders(step.Pattern, satisfied, seen)
			collectPlaceholders(step.Script, satisfied, seen)
			if step.Target != nil {
				collectSelectorPlaceholders(step.Target, satisfied, seen)
			}
			for i := range step.Targets {
				collectSelectorPlaceholders(&step.Targets[i], satisfied, seen)
			}
		}
	}

	result := make([]string, 0, len(seen))
	for name := range seen {
		result = append(result, name)
	}
	sort.Strings(result)
	return result
}

func collectPlaceholders(s string, satisfied, seen map[string]bool) {
	for _, match := range varPattern.FindAllStringSubmatch(s, -1) {
		if len(match) < 2 {
			continue
		}
		name := match[1]
		if builtInVars[name] || credentialVars[name] || satisfied[name] {
			continue
		}
		seen[name] = true
	}
}

func collectSelectorPlaceholders(sel *Selector, satisfied, seen map[string]bool) {
	for _, field := range []string{sel.AriaLabel, sel.ID, sel.Name, sel.Text, sel.CSSPath} {
		collectPlaceholders(field, satisfied, seen)
	}
}
