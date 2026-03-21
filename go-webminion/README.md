# go-webminion

Go implementation of WebMinion — a metadata-driven browser automation library. Define your automation as a YAML or JSON config; the library handles execution, retries, credential injection, and download detection.

## Features

- Chrome automation via [go-rod](https://github.com/go-rod/rod)
- Lightweight HTTP driver (no JavaScript) via [goquery](https://github.com/PuerkitoBio/goquery)
- YAML and JSON config auto-detection
- Variable interpolation with built-in (`today`, `start_date`, `end_date`, `flow_name`) and runtime variables
- Credential vault: environment variables or OS keyring
- Download watcher for automated completion detection
- Flow graph execution with `on_success` / `on_failure` routing and cycle detection

## Requirements

- Go 1.21+
- Chrome (for the Chrome driver)

## Installation

```bash
go get github.com/stevendaniels/web_minion/go-webminion
```

## Library usage

```go
import (
    webminion "github.com/stevendaniels/web_minion/go-webminion"
    chromedv  "github.com/stevendaniels/web_minion/go-webminion/driver/chrome"
    "github.com/stevendaniels/web_minion/go-webminion/vault"
)

inst, err := webminion.LoadConfig("my_site.yaml")
if err != nil {
    log.Fatal(err)
}

flow := webminion.NewFlow(inst)

d := chromedv.New(true /* headless */)
defer d.Close()

v := vault.NewEnvVault()

if err := flow.Run(d, v); err != nil {
    log.Fatal(err)
}
```

## CLI usage

```bash
go build -o webminion ./cmd/webminion/

# Validate a config file
./webminion --command validate --config my_site.yaml

# Run a flow
./webminion --command run --config my_site.yaml

# Run with options
./webminion --command run --config my_site.yaml \
  --driver chrome \
  --headless \
  --vault env \
  --start-date 2024-01-01 \
  --end-date 2024-12-31
```

### Flags

| Flag | Default | Description |
|------|---------|-------------|
| `--command` | *(required)* | `run` or `validate` |
| `--config` | *(required)* | Path to YAML or JSON config file |
| `--driver` | `chrome` | `chrome` or `http` |
| `--headless` | `true` | Run Chrome in headless mode |
| `--vault` | `env` | `env` or `keyring` |
| `--start-date` | today | Date range start (`YYYY-MM-DD`) |
| `--end-date` | today | Date range end (`YYYY-MM-DD`) |

## Config format

```yaml
id: my_site
name: My Site
base_url: https://example.com
download_dir: ~/downloads

flow:
  actions:
    - key: login
      starting: true
      on_success: download
      on_failure: WebMinion_FlowError
      steps:
        - method: go
          value: https://example.com/login
        - method: fill_in_input
          target:
            css_path: "#email"
          value: "{{username}}"
        - method: fill_in_input
          target:
            css_path: "#password"
          value: "{{password}}"
        - method: click
          target:
            css_path: "button[type=submit]"
      expects:
        type: url_includes
        value: dashboard

    - key: download
      on_success: WebMinion_Finalizer
      steps:
        - method: click
          target:
            aria_label: Export CSV
        - method: wait_for_download
          pattern: "*.csv"
          timeout: 30
```

## Credential vault

### Environment variables (`--vault env`)

Variables are resolved as `MINION_<INSTITUTION_ID>_<KEY>` (uppercased):

```bash
export MINION_MY_SITE_USERNAME="user@example.com"
export MINION_MY_SITE_PASSWORD="s3cr3t"
```

### OS keyring (`--vault keyring`)

Secrets are stored under service `minion` with key `<institution_id>/<key>`:

```bash
# macOS Keychain example
security add-generic-password -s minion -a my_site/username -w "user@example.com"
```

## Step reference

| Method | Required fields | Description |
|--------|-----------------|-------------|
| `go` | `value` (URL) | Navigate to URL |
| `click` | `target` | Click an element |
| `fill_in_input` | `target`, `value` | Fill an input field |
| `submit` | `target` | Submit a form element |
| `select` | `target`, `value` | Select a dropdown option |
| `save_value` | `target`, `value` (var name) | Save element text to a variable |
| `format_saved_value` | `value` (var name), `pattern` | Format a string with interpolation |
| `wait` | `target`, `timeout` | Wait for element to appear |
| `wait_for_download` | `pattern`, `timeout` | Wait for a file to appear in download dir |
| `execute_script` | `script` | Run JavaScript (Chrome only) |
| `save_page_html` | `value` (var name) | Save current page HTML to a variable |

## Selector strategies (priority order)

1. `css_path` — raw CSS selector
2. `id` — element ID (`#<id>`)
3. `name` — name attribute (`[name="..."]`)
4. `aria_label` — ARIA label attribute
5. `aria_role` — ARIA role attribute
6. `data_testid` — data-testid attribute
7. `text` — visible text content

## Running tests

```bash
go test ./...                        # unit tests
go test -race ./...                  # with race detector
go test -tags=integration ./...      # include Chrome integration tests
```
