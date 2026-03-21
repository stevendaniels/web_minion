package executor

import (
	"fmt"
	"time"

	"github.com/stevendaniels/web_minion/go-webminion/config"
	"github.com/stevendaniels/web_minion/go-webminion/driver"
	"github.com/stevendaniels/web_minion/go-webminion/watcher"
)

func newStepRegistry() map[string]StepFunc {
	return map[string]StepFunc{
		"go":                  stepGo,
		"get_form":            stepGetForm,
		"fill_in_input":       stepFillInInput,
		"click":               stepClick,
		"submit":              stepSubmit,
		"select":              stepSelect,
		"save_value":          stepSaveValue,
		"format_saved_value":  stepFormatSavedValue,
		"wait_for_download":   stepWaitForDownload,
		"wait":                stepWait,
		"execute_script":      stepExecuteScript,
		"save_page_html":      stepSavePageHTML,
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
