package config

import (
	"fmt"
	"regexp"
)

var varPattern = regexp.MustCompile(`\{\{(\w+)\}\}`)

var validMethods = map[string]bool{
	"go": true, "get_field": true, "get_form": true, "select": true,
	"click": true, "click_button_in_form": true, "submit": true,
	"fill_in_input": true, "url_equals": true, "value_equals": true,
	"body_includes": true, "save_page_html": true, "save_value": true,
	"wait": true, "format_saved_value": true, "write_file": true,
	"html_to_markdown": true, "wait_for_reply": true, "execute_script": true,
	"wait_for_download": true,
}

// ValidateConfig checks a parsed Config for logical errors.
func ValidateConfig(cfg *Config) []error {
	var errs []error

	startingCount := 0
	seen := make(map[string]bool)
	for _, action := range cfg.Flow.Actions {
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

	for _, action := range cfg.Flow.Actions {
		if len(action.Steps) == 0 {
			errs = append(errs, fmt.Errorf("action %q has no steps", action.Key))
			continue
		}

		if !action.Steps[len(action.Steps)-1].IsValidator {
			errs = append(errs, fmt.Errorf("action %q's final step is not a validation step", action.Key))
		}

		for i, step := range action.Steps {
			errs = append(errs, validateStep(action.Key, i, step)...)
		}
	}

	errs = append(errs, validateActionReferences(cfg.Flow.Actions)...)
	if len(errs) == 0 {
		errs = append(errs, detectCycles(cfg.Flow.Actions)...)
	}

	if len(errs) == 0 {
		errs = append(errs, validateVariables(cfg)...)
	}

	return errs
}

func validateStep(actionKey string, idx int, step Step) []error {
	var errs []error
	loc := fmt.Sprintf("action %q step %d (%s)", actionKey, idx, step.Method)
	if !validMethods[step.Method] {
		errs = append(errs, fmt.Errorf("%s: invalid method %q", loc, step.Method))
		return errs
	}
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
	case "extract_readable":
		if step.Value == "" {
			errs = append(errs, fmt.Errorf("%s: value (variable name) is required", loc))
		}
	}
	return errs
}

func validateVariables(inst *Config) []error {
	var errs []error

	knownVars := make(map[string]bool)
	if inst.Vars != nil {
		for k := range inst.Vars {
			knownVars[k] = true
		}
	}

	for _, action := range inst.Flow.Actions {
		for _, step := range action.Steps {
			for _, text := range []string{step.Value, step.Pattern} {
				for _, match := range varPattern.FindAllStringSubmatch(text, -1) {
					if len(match) > 1 {
						varName := match[1]
						if !knownVars[varName] {
							errs = append(errs, fmt.Errorf("action %q step %d (%s): variable %q not defined in config.vars", action.Key, 0, step.Method, varName))
						}
					}
				}
			}
			if step.Target != nil {
				for _, field := range []string{step.Target.AriaLabel, step.Target.ID, step.Target.Name, step.Target.Text, step.Target.CSSPath} {
					for _, match := range varPattern.FindAllStringSubmatch(field, -1) {
						if len(match) > 1 {
							varName := match[1]
							if !knownVars[varName] {
								errs = append(errs, fmt.Errorf("action %q step %d (%s): variable %q not defined in config.vars", action.Key, 0, step.Method, varName))
							}
						}
					}
				}
			}
			for _, target := range step.Targets {
				for _, field := range []string{target.AriaLabel, target.ID, target.Name, target.Text, target.CSSPath} {
					for _, match := range varPattern.FindAllStringSubmatch(field, -1) {
						if len(match) > 1 {
							varName := match[1]
							if !knownVars[varName] {
								errs = append(errs, fmt.Errorf("action %q step %d (%s): variable %q not defined in config.vars", action.Key, 0, step.Method, varName))
							}
						}
					}
				}
			}
		}
	}

	return errs
}

func validateActionReferences(actions []Action) []error {
	var errs []error

	actionKeys := make(map[string]bool)
	for _, action := range actions {
		actionKeys[action.Key] = true
	}

	for _, action := range actions {
		if action.OnSuccess != "" && action.OnSuccess != "WebMinion_Finalizer" && action.OnSuccess != "WebMinion_FlowError" {
			if !actionKeys[action.OnSuccess] {
				errs = append(errs, fmt.Errorf("action %q references non-existent success action %q", action.Key, action.OnSuccess))
			}
		}
		if action.OnFailure != "" && action.OnFailure != "WebMinion_Finalizer" && action.OnFailure != "WebMinion_FlowError" {
			if !actionKeys[action.OnFailure] {
				errs = append(errs, fmt.Errorf("action %q references non-existent failure action %q", action.Key, action.OnFailure))
			}
		}
	}

	return errs
}

func detectCycles(actions []Action) []error {
	var errs []error

	actionMap := make(map[string]*Action)
	for i := range actions {
		actionMap[actions[i].Key] = &actions[i]
	}

	var startingAction *Action
	for i := range actions {
		if actions[i].Starting {
			startingAction = &actions[i]
			break
		}
	}

	if startingAction == nil {
		return errs
	}

	visited := make(map[string]bool)
	onStack := make(map[string]bool)

	var hasCycle func(action *Action) bool
	hasCycle = func(action *Action) bool {
		if onStack[action.Key] {
			return true
		}
		if visited[action.Key] {
			return false
		}

		visited[action.Key] = true
		onStack[action.Key] = true

		traverse := func(key string) {
			if key != "" && key != "WebMinion_Finalizer" && key != "WebMinion_FlowError" {
				if nextAction, exists := actionMap[key]; exists {
					if hasCycle(nextAction) {
						errs = append(errs, fmt.Errorf("flow contains a cycle involving action %q", nextAction.Key))
					}
				}
			}
		}

		traverse(action.OnSuccess)
		traverse(action.OnFailure)

		onStack[action.Key] = false
		return false
	}

	hasCycle(startingAction)

	return errs
}
