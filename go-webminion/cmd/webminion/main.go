package main

import (
	"flag"
	"fmt"
	"os"
	"time"

	webminion "github.com/stevendaniels/web_minion/go-webminion"
	chromedv "github.com/stevendaniels/web_minion/go-webminion/driver/chrome"
	httpdv "github.com/stevendaniels/web_minion/go-webminion/driver/http"
	"github.com/stevendaniels/web_minion/go-webminion/vault"
)

func main() {
	if err := run(os.Args[1:]); err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}
}

func run(args []string) error {
	fs := flag.NewFlagSet("webminion", flag.ContinueOnError)

	command := fs.String("command", "", "Command to run: run | validate")
	configPath := fs.String("config", "", "Path to config file (YAML or JSON)")
	headless := fs.Bool("headless", true, "Run Chrome in headless mode")
	driverName := fs.String("driver", "chrome", "Driver to use: chrome | http")
	vaultType := fs.String("vault", "env", "Credential vault: env | keyring")
	startDateStr := fs.String("start-date", "", "Start date for date range (YYYY-MM-DD)")
	endDateStr := fs.String("end-date", "", "End date for date range (YYYY-MM-DD)")

	if err := fs.Parse(args); err != nil {
		return err
	}

	switch *command {
	case "run":
		return cmdRun(*configPath, *driverName, *vaultType, *headless, *startDateStr, *endDateStr)
	case "validate":
		return cmdValidate(*configPath)
	case "":
		fs.Usage()
		return fmt.Errorf("--command is required")
	default:
		return fmt.Errorf("unknown command: %q (valid: run, validate)", *command)
	}
}

func cmdRun(configPath, driverName, vaultType string, headless bool, startDateStr, endDateStr string) error {
	if configPath == "" {
		return fmt.Errorf("--config is required")
	}

	flow, err := webminion.NewFlow(configPath)
	if err != nil {
		return fmt.Errorf("loading config: %w", err)
	}

	if startDateStr != "" || endDateStr != "" {
		start, end, err := parseDateRange(startDateStr, endDateStr)
		if err != nil {
			return err
		}
		flow = flow.WithDateRange(start, end)
	}

	var v vault.Vault
	switch vaultType {
	case "keyring":
		v = vault.NewKeyringVault()
	default:
		v = vault.NewEnvVault()
	}

	switch driverName {
	case "http":
		d := httpdv.New()
		defer d.Close()
		return flow.Run(d, v)
	default:
		d := chromedv.New(headless)
		defer d.Close()
		return flow.Run(d, v)
	}
}

func cmdValidate(configPath string) error {
	if configPath == "" {
		return fmt.Errorf("--config is required")
	}
	if _, err := webminion.NewFlow(configPath); err != nil {
		return err
	}
	fmt.Println("config is valid")
	return nil
}

func parseDateRange(startStr, endStr string) (time.Time, time.Time, error) {
	now := time.Now()
	start, end := now, now

	if startStr != "" {
		t, err := time.Parse("2006-01-02", startStr)
		if err != nil {
			return time.Time{}, time.Time{}, fmt.Errorf("invalid --start-date %q: %w", startStr, err)
		}
		start = t
	}
	if endStr != "" {
		t, err := time.Parse("2006-01-02", endStr)
		if err != nil {
			return time.Time{}, time.Time{}, fmt.Errorf("invalid --end-date %q: %w", endStr, err)
		}
		end = t
	}
	return start, end, nil
}
