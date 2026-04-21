package executor

import (
	"crypto/rand"
	"fmt"
	"net/url"
	"os"
	"path/filepath"
	"strings"
	"time"

	md "github.com/JohannesKaufmann/html-to-markdown"
	"github.com/stevendaniels/web_minion/go-webminion/config"
	"github.com/stevendaniels/web_minion/go-webminion/driver"
	"github.com/stevendaniels/web_minion/go-webminion/watcher"
)

func newStepRegistry() map[string]StepFunc {
	return map[string]StepFunc{
		"go":                 stepGo,
		"get_form":           stepGetForm,
		"fill_in_input":      stepFillInInput,
		"click":              stepClick,
		"submit":             stepSubmit,
		"select":             stepSelect,
		"save_value":         stepSaveValue,
		"format_saved_value": stepFormatSavedValue,
		"wait_for_download":  stepWaitForDownload,
		"wait":               stepWait,
		"execute_script":     stepExecuteScript,
		"save_page_html":     stepSavePageHTML,
		"write_file":         stepWriteFile,
		"html_to_markdown":   stepHTMLToMarkdown,
		"wait_for_reply":     stepWaitForReply,
	}
}

func stepGo(e *Executor, step config.Step) error {
	url, err := e.resolve(step.Value)
	if err != nil {
		return err
	}
	if url == "" {
		return fmt.Errorf("step 'go' requires a value (URL)")
	}
	return e.driver.Navigate(url)
}

func stepGetForm(e *Executor, step config.Step) error {
	// HTTP driver uses this to select a form; Chrome driver navigates to it directly.
	// For now this is a no-op — form selection is implicit via fill_in_input.
	return nil
}

func stepFillInInput(e *Executor, step config.Step) error {
	if step.Target == nil {
		return fmt.Errorf("step 'fill_in_input' requires a target selector")
	}
	sel, err := driver.Resolve(step.Target)
	if err != nil {
		return err
	}

	value, err := e.resolve(step.Value)
	if err != nil {
		return err
	}

	// Resolve credential placeholders via vault if still unresolved.
	if value == "{{username}}" || value == "{{password}}" || value == "{{totp}}" {
		if e.vault == nil {
			return fmt.Errorf("credential %q referenced but no vault configured", value)
		}
		key := value[2 : len(value)-2] // strip {{ }}
		value, err = e.vault.Get(e.inst.ID, key)
		if err != nil {
			return err
		}
	}

	return e.driver.Fill(sel, value)
}

func stepClick(e *Executor, step config.Step) error {
	if step.Target == nil {
		return fmt.Errorf("step 'click' requires a target selector")
	}
	sel, err := driver.Resolve(step.Target)
	if err != nil {
		return err
	}
	return e.driver.Click(sel)
}

func stepSubmit(e *Executor, step config.Step) error {
	if step.Target == nil {
		return fmt.Errorf("step 'submit' requires a target selector")
	}
	sel, err := driver.Resolve(step.Target)
	if err != nil {
		return err
	}
	return e.driver.Submit(sel)
}

func stepSelect(e *Executor, step config.Step) error {
	if step.Target == nil {
		return fmt.Errorf("step 'select' requires a target selector")
	}
	sel, err := driver.Resolve(step.Target)
	if err != nil {
		return err
	}
	value, err := e.resolve(step.Value)
	if err != nil {
		return err
	}
	return e.driver.Select(sel, value)
}

func stepSaveValue(e *Executor, step config.Step) error {
	if step.Target == nil {
		return fmt.Errorf("step 'save_value' requires a target selector")
	}
	sel, err := driver.Resolve(step.Target)
	if err != nil {
		return err
	}
	text, err := e.driver.GetText(sel)
	if err != nil {
		return err
	}
	key, err := e.resolve(step.Value)
	if err != nil {
		return err
	}
	return e.vars.Set(key, text)
}

func stepFormatSavedValue(e *Executor, step config.Step) error {
	// step.Values should be map[string]string{"key": "varname", "format": "..."}
	// For now accept value as "varname" and pattern as format string.
	key := step.Value
	if key == "" {
		return fmt.Errorf("step 'format_saved_value' requires value (variable name)")
	}
	formatted, err := e.resolve(step.Pattern)
	if err != nil {
		return err
	}
	return e.vars.Set(key, formatted)
}

func stepWaitForDownload(e *Executor, step config.Step) error {
	timeout := time.Duration(step.Timeout) * time.Second
	if timeout == 0 {
		timeout = 30 * time.Second
	}
	pattern := step.Pattern
	if pattern == "" {
		pattern = "*"
	}
	_, err := watcher.Watch(e.inst.DownloadDir, pattern, timeout)
	return err
}

func stepWait(e *Executor, step config.Step) error {
	if step.Target == nil {
		return fmt.Errorf("step 'wait' requires a target selector")
	}
	sel, err := driver.Resolve(step.Target)
	if err != nil {
		return err
	}
	timeout := step.Timeout
	if timeout == 0 {
		timeout = 30
	}
	return e.driver.WaitForSelector(sel, timeout)
}

func stepExecuteScript(e *Executor, step config.Step) error {
	script, err := e.resolve(step.Script)
	if err != nil {
		return err
	}
	if script == "" {
		return fmt.Errorf("step 'execute_script' requires a script")
	}
	return e.driver.Eval(script)
}

func stepSavePageHTML(e *Executor, step config.Step) error {
	key := step.Value
	if key == "" {
		return fmt.Errorf("step 'save_page_html' requires value (variable name)")
	}
	html, err := e.driver.PageHTML()
	if err != nil {
		return err
	}
	return e.vars.Set(key, html)
}

func stepWriteFile(e *Executor, step config.Step) error {
	if step.Pattern == "" {
		return fmt.Errorf("step 'write_file' requires pattern (file path)")
	}
	if step.Value == "" {
		return fmt.Errorf("step 'write_file' requires value (content)")
	}

	uuid, err := newUUID()
	if err != nil {
		return fmt.Errorf("step 'write_file': %w", err)
	}
	extras := map[string]string{
		"uuid": uuid,
		"now":  time.Now().UTC().Format(time.RFC3339),
	}

	path, err := e.resolveWith(step.Pattern, extras)
	if err != nil {
		return err
	}
	if path == "" {
		return fmt.Errorf("step 'write_file': pattern resolved to empty path")
	}

	content, err := e.resolveWith(step.Value, extras)
	if err != nil {
		return err
	}

	if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
		return fmt.Errorf("step 'write_file': mkdir: %w", err)
	}
	return os.WriteFile(path, []byte(content), 0644)
}

// stepHTMLToMarkdown converts an HTML variable to Markdown in-place.
// step.Value is the variable name holding the HTML; it is overwritten with Markdown.
func stepHTMLToMarkdown(e *Executor, step config.Step) error {
	varName := step.Value
	if varName == "" {
		return fmt.Errorf("step 'html_to_markdown' requires value (variable name)")
	}
	html, err := e.resolve("{{" + varName + "}}")
	if err != nil {
		return fmt.Errorf("step 'html_to_markdown': variable %q not set: %w", varName, err)
	}
	domain := e.inst.BaseURL
	if u, err := url.Parse(e.inst.BaseURL); err == nil && u.Host != "" {
		domain = u.Host
	}
	converter := md.NewConverter(domain, true, nil)
	markdown, err := converter.ConvertString(html)
	if err != nil {
		return fmt.Errorf("step 'html_to_markdown': %w", err)
	}
	return e.vars.Set(varName, markdown)
}

// stepWaitForReply polls for a file at the path in step.Pattern every 500ms.
// When the file appears, its content is trimmed and saved as the variable named
// by step.Value, then the file is deleted. Returns an error on timeout.
func stepWaitForReply(e *Executor, step config.Step) error {
	path, err := e.resolve(step.Pattern)
	if err != nil {
		return err
	}
	if path == "" {
		return fmt.Errorf("step 'wait_for_reply' requires pattern (file path)")
	}
	varName := step.Value
	if varName == "" {
		return fmt.Errorf("step 'wait_for_reply' requires value (variable name)")
	}
	timeout := time.Duration(step.Timeout) * time.Second
	if timeout == 0 {
		timeout = 300 * time.Second
	}

	timer := time.NewTimer(timeout)
	defer timer.Stop()
	ticker := time.NewTicker(500 * time.Millisecond)
	defer ticker.Stop()

	for {
		if data, readErr := os.ReadFile(path); readErr == nil {
			content := strings.TrimSpace(string(data))
			if setErr := e.vars.Set(varName, content); setErr != nil {
				return fmt.Errorf("step 'wait_for_reply': save variable: %w", setErr)
			}
			_ = os.Remove(path)
			return nil
		}
		select {
		case <-timer.C:
			return fmt.Errorf("step 'wait_for_reply': timeout waiting for %s", path)
		case <-ticker.C:
		}
	}
}

func newUUID() (string, error) {
	b := make([]byte, 16)
	if _, err := rand.Read(b); err != nil {
		return "", fmt.Errorf("newUUID: %w", err)
	}
	b[6] = (b[6] & 0x0f) | 0x40 // version 4
	b[8] = (b[8] & 0x3f) | 0x80 // variant bits
	return fmt.Sprintf("%x-%x-%x-%x-%x", b[0:4], b[4:6], b[6:8], b[8:10], b[10:16]), nil
}
