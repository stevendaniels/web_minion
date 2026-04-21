package executor

import (
	"fmt"
	"strings"
	"time"

	"github.com/stevendaniels/web_minion/go-webminion/config"
	"github.com/stevendaniels/web_minion/go-webminion/driver"
	"github.com/stevendaniels/web_minion/go-webminion/interp"
	"github.com/stevendaniels/web_minion/go-webminion/vault"
	"github.com/stevendaniels/web_minion/go-webminion/watcher"
)

// StepFunc is a handler for a single step method.
type StepFunc func(e *Executor, step config.Step) error

// Executor runs the flow graph defined in a Config.
type Executor struct {
	cfg      *config.Config
	driver   driver.Driver
	vault    vault.Vault
	vars     *interp.Vars
	registry map[string]StepFunc
	visited  map[string]bool
}

// New creates an Executor. startDate and endDate are used for built-in interpolation vars.
func New(cfg *config.Config, d driver.Driver, v vault.Vault, startDate, endDate time.Time) *Executor {
	return &Executor{
		cfg:      cfg,
		driver:   d,
		vault:    v,
		vars:     interp.New(cfg.ID, startDate, endDate),
		registry: newStepRegistry(),
		visited:  make(map[string]bool),
	}
}

// Run executes the flow starting from the starting action.
func (e *Executor) Run() error {
	currentKey := findStartingAction(e.cfg)
	if currentKey == "" {
		return fmt.Errorf("no starting action found")
	}

	for currentKey != "" {
		if e.visited[currentKey] {
			return fmt.Errorf("cycle detected at action: %s", currentKey)
		}
		e.visited[currentKey] = true

		action := findAction(e.cfg, currentKey)
		if action == nil {
			return fmt.Errorf("action not found: %s", currentKey)
		}

		stepErr := e.executeSteps(action)

		var next string
		if stepErr != nil {
			next = action.OnFailure
			if next == "" {
				return stepErr
			}
		} else if action.Expects != nil {
			if e.evaluateExpects(action) {
				next = action.OnSuccess
			} else {
				next = action.OnFailure
			}
		} else {
			next = action.OnSuccess
		}

		// Terminal sentinels exit cleanly.
		if next == "WebMinion_Finalizer" || next == "WebMinion_FlowError" {
			return stepErr
		}

		currentKey = next
	}

	return nil
}

func findStartingAction(cfg *config.Config) string {
	for _, a := range cfg.Flow.Actions {
		if a.Starting {
			return a.Key
		}
	}
	return ""
}

func findAction(cfg *config.Config, key string) *config.Action {
	for i := range cfg.Flow.Actions {
		if cfg.Flow.Actions[i].Key == key {
			return &cfg.Flow.Actions[i]
		}
	}
	return nil
}

func (e *Executor) executeSteps(action *config.Action) error {
	for _, step := range action.Steps {
		handler, ok := e.registry[step.Method]
		if !ok {
			return fmt.Errorf("unknown step method: %s", step.Method)
		}
		if err := handler(e, step); err != nil {
			return fmt.Errorf("step %q failed: %w", step.Method, err)
		}
	}
	return nil
}

func (e *Executor) evaluateExpects(action *config.Action) bool {
	exp := action.Expects
	if exp == nil {
		return true
	}
	switch exp.Type {
	case "body_includes":
		ok, err := e.driver.BodyContains(exp.Value)
		return err == nil && ok
	case "url_includes":
		url, err := e.driver.CurrentURL()
		if err != nil {
			return false
		}
		for _, part := range strings.Split(exp.Value, "|") {
			if strings.Contains(url, part) {
				return true
			}
		}
		return false
	case "element_present":
		if exp.Target == nil {
			return false
		}
		sel, err := driver.Resolve(exp.Target)
		if err != nil {
			return false
		}
		return e.driver.WaitForSelector(sel, 1) == nil
	case "download_exists":
		_, err := watcher.Watch(e.cfg.DownloadDir, exp.Pattern, 1*time.Second)
		return err == nil
	}
	return false
}

// resolve interpolates a string value using the executor's Vars.
func (e *Executor) resolve(s string) (string, error) {
	return e.vars.Resolve(s)
}

func (e *Executor) resolveWith(s string, overrides map[string]string) (string, error) {
	return e.vars.ResolveWith(s, overrides)
}
