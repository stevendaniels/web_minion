// Package webminion is the public API for the Go WebMinion automation library.
// Load a YAML/JSON config, create a Flow, and call Run with a credential Vault.
package webminion

import (
	"time"

	"github.com/stevendaniels/web_minion/go-webminion/config"
	"github.com/stevendaniels/web_minion/go-webminion/driver"
	"github.com/stevendaniels/web_minion/go-webminion/executor"
	"github.com/stevendaniels/web_minion/go-webminion/vault"
)

// LoadConfig reads and validates a YAML or JSON config file.
func LoadConfig(path string) (*config.Institution, error) {
	return config.LoadConfig(path)
}

// Flow wraps an Institution config and is the primary entry point for running automations.
type Flow struct {
	inst      *config.Institution
	startDate time.Time
	endDate   time.Time
}

// NewFlow creates a Flow for the given Institution. Date range defaults to today.
func NewFlow(inst *config.Institution) *Flow {
	now := time.Now()
	return &Flow{inst: inst, startDate: now, endDate: now}
}

// WithDateRange sets the start and end dates used for built-in interpolation variables.
func (f *Flow) WithDateRange(start, end time.Time) *Flow {
	f.startDate = start
	f.endDate = end
	return f
}

// Run executes the automation flow using the provided driver and vault.
func (f *Flow) Run(d driver.Driver, v vault.Vault) error {
	ex := executor.New(f.inst, d, v, f.startDate, f.endDate)
	return ex.Run()
}
