package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"strings"
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
	if len(args) == 0 {
		printUsage()
		return fmt.Errorf("subcommand required")
	}

	sub := args[0]
	rest := args[1:]

	switch sub {
	case "run":
		return cmdRun(rest)
	case "validate":
		return cmdValidate(rest)
	case "help":
		return cmdHelp(rest)
	default:
		printUsage()
		return fmt.Errorf("unknown subcommand %q (valid: run, validate, help)", sub)
	}
}

func printUsage() {
	fmt.Fprintln(os.Stderr, "Usage:")
	fmt.Fprintln(os.Stderr, "  wm run      --config=<path> [--arg key=val]... [--env-file=<path>]")
	fmt.Fprintln(os.Stderr, "  wm validate --config=<path>")
	fmt.Fprintln(os.Stderr, "  wm help     --config=<path>")
}

// argMap is a repeatable --arg flag collector.
type argMap map[string]string

func (a argMap) String() string { return fmt.Sprintf("%v", map[string]string(a)) }
func (a argMap) Set(s string) error {
	k, v, ok := strings.Cut(s, "=")
	if !ok {
		return fmt.Errorf("--arg must be key=value, got %q", s)
	}
	a[k] = v
	return nil
}

func cmdRun(args []string) error {
	fs := flag.NewFlagSet("run", flag.ContinueOnError)
	configPath := fs.String("config", "", "Path to config file (YAML or JSON)")
	headless := fs.Bool("headless", true, "Run Chrome in headless mode")
	driverName := fs.String("driver", "chrome", "Driver: chrome | http")
	vaultType := fs.String("vault", "env", "Credential vault: env | keyring")
	startDateStr := fs.String("start-date", "", "Start date (YYYY-MM-DD)")
	endDateStr := fs.String("end-date", "", "End date (YYYY-MM-DD)")
	envFile := fs.String("env-file", "", "Path to .env file with KEY=VALUE pairs")

	runtimeArgs := make(argMap)
	fs.Var(runtimeArgs, "arg", "Runtime arg as key=value (repeatable)")

	if err := fs.Parse(args); err != nil {
		return err
	}
	if *configPath == "" {
		return fmt.Errorf("--config is required")
	}

	flow, err := webminion.NewFlow(*configPath)
	if err != nil {
		return fmt.Errorf("loading config: %w", err)
	}

	if *startDateStr != "" || *endDateStr != "" {
		start, end, err := parseDateRange(*startDateStr, *endDateStr)
		if err != nil {
			return err
		}
		flow = flow.WithDateRange(start, end)
	}

	merged := make(map[string]string)
	if *envFile != "" {
		pairs, err := parseEnvFile(*envFile)
		if err != nil {
			return err
		}
		for k, v := range pairs {
			merged[k] = v
		}
	}
	for k, v := range runtimeArgs {
		merged[k] = v
	}
	if len(merged) > 0 {
		flow = flow.WithArgs(merged)
	}

	var v vault.Vault
	switch *vaultType {
	case "keyring":
		v = vault.NewKeyringVault()
	default:
		v = vault.NewEnvVault()
	}

	switch *driverName {
	case "http":
		d := httpdv.New()
		defer d.Close()
		return flow.Run(d, v)
	default:
		d := chromedv.New(*headless)
		defer d.Close()
		return flow.Run(d, v)
	}
}

func cmdValidate(args []string) error {
	fs := flag.NewFlagSet("validate", flag.ContinueOnError)
	configPath := fs.String("config", "", "Path to config file (YAML or JSON)")
	if err := fs.Parse(args); err != nil {
		return err
	}
	if *configPath == "" {
		return fmt.Errorf("--config is required")
	}
	if _, err := webminion.NewFlow(*configPath); err != nil {
		return err
	}
	fmt.Println("config is valid")
	return nil
}

func cmdHelp(args []string) error {
	fs := flag.NewFlagSet("help", flag.ContinueOnError)
	configPath := fs.String("config", "", "Path to config file (YAML or JSON)")
	if err := fs.Parse(args); err != nil {
		return err
	}
	if *configPath == "" {
		return fmt.Errorf("--config is required")
	}

	flow, err := webminion.NewFlow(*configPath)
	if err != nil {
		return fmt.Errorf("loading config: %w", err)
	}

	required := flow.RequiredArgs()
	if len(required) == 0 {
		fmt.Println("No required args.")
		return nil
	}
	fmt.Println("Required args:")
	for _, name := range required {
		fmt.Printf("  %s\n", name)
	}
	return nil
}

func parseEnvFile(path string) (map[string]string, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("opening env file: %w", err)
	}
	defer f.Close()

	result := make(map[string]string)
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		k, v, ok := strings.Cut(line, "=")
		if !ok {
			return nil, fmt.Errorf("env file line %q is not KEY=VALUE", line)
		}
		result[strings.TrimSpace(k)] = strings.TrimSpace(v)
	}
	return result, scanner.Err()
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
