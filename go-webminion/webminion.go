// Package webminion is the public API for the Go WebMinion automation library.
// Load a config file and call Run with a credential Vault.
package webminion

import (
	"time"

	"github.com/stevendaniels/web_minion/go-webminion/config"
	"github.com/stevendaniels/web_minion/go-webminion/driver"
	"github.com/stevendaniels/web_minion/go-webminion/executor"
	"github.com/stevendaniels/web_minion/go-webminion/vault"
)

// Flow wraps a validated config and is the primary entry point for running automations.
type Flow struct {
	cfg       *config.Config
	startDate time.Time
	endDate   time.Time
}

// NewFlow loads and validates a YAML or JSON config file, returning a runnable Flow.
func NewFlow(path string) (*Flow, error) {
	cfg, err := config.LoadConfig(path)
	if err != nil {
		return nil, err
	}
	now := time.Now()
	return &Flow{cfg: cfg, startDate: now, endDate: now}, nil
}

// WithDateRange sets the start and end dates used for built-in interpolation variables.
func (f *Flow) WithDateRange(start, end time.Time) *Flow {
	f.startDate = start
	f.endDate = end
	return f
}

// Run executes the automation flow using the provided driver and vault.
func (f *Flow) Run(d driver.Driver, v vault.Vault) error {
	ex := executor.New(f.cfg, d, v, f.startDate, f.endDate)
	return ex.Run()
}
