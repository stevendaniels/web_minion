package executor

import (
	"testing"
	"time"

	"github.com/stevendaniels/web_minion/go-webminion/config"
)

func runStep(t *testing.T, fd *fakeDriver, step config.Step) error {
	t.Helper()
	inst := makeInst(config.Action{
		Key:      "start",
		Starting: true,
		Steps:    []config.Step{step},
	})
	return newTestExecutor(inst, fd).Run()
}

func TestStepGo(t *testing.T) {
	fd := &fakeDriver{}
	if err := runStep(t, fd, config.Step{Method: "go", Value: "https://example.com"}); err != nil {
		t.Fatal(err)
	}
	if len(fd.navigated) == 0 || fd.navigated[0] != "https://example.com" {
		t.Errorf("expected Navigate('https://example.com'), got %v", fd.navigated)
	}
}

func TestStepGo_MissingValue(t *testing.T) {
	fd := &fakeDriver{}
	if err := runStep(t, fd, config.Step{Method: "go", Value: ""}); err == nil {
		t.Error("expected error for empty URL, got nil")
	}
}

func TestStepClick(t *testing.T) {
	fd := &fakeDriver{}
	err := runStep(t, fd, config.Step{
		Method: "click",
		Target: &config.Selector{CSSPath: ".btn"},
	})
	if err != nil {
		t.Fatal(err)
	}
	if len(fd.clicked) == 0 {
		t.Error("expected Click to be called")
	}
}

func TestStepClick_NoTarget(t *testing.T) {
	fd := &fakeDriver{}
	if err := runStep(t, fd, config.Step{Method: "click"}); err == nil {
		t.Error("expected error for missing target, got nil")
	}
}

func TestStepFillInInput(t *testing.T) {
	fd := &fakeDriver{}
	err := runStep(t, fd, config.Step{
		Method: "fill_in_input",
		Target: &config.Selector{CSSPath: "#email"},
		Value:  "user@example.com",
	})
	if err != nil {
		t.Fatal(err)
	}
	if len(fd.filled) == 0 || fd.filled[0].value != "user@example.com" {
		t.Errorf("expected fill with 'user@example.com', got %v", fd.filled)
	}
}

func TestStepFillInInput_WithInterp(t *testing.T) {
	fd := &fakeDriver{}
	inst := makeInst(config.Action{
		Key:      "start",
		Starting: true,
		Steps: []config.Step{{
			Method: "fill_in_input",
			Target: &config.Selector{CSSPath: "#date"},
			Value:  "{{today}}",
		}},
	})
	ex := newTestExecutor(inst, fd)
	if err := ex.Run(); err != nil {
		t.Fatal(err)
	}
	today := time.Now().Format("2006-01-02")
	if len(fd.filled) == 0 || fd.filled[0].value != today {
		t.Errorf("expected fill with today's date '%s', got %v", today, fd.filled)
	}
}

func TestStepSubmit(t *testing.T) {
	fd := &fakeDriver{}
	err := runStep(t, fd, config.Step{
		Method: "submit",
		Target: &config.Selector{CSSPath: "form"},
	})
	if err != nil {
		t.Fatal(err)
	}
}

func TestStepSelect(t *testing.T) {
	fd := &fakeDriver{}
	err := runStep(t, fd, config.Step{
		Method: "select",
		Target: &config.Selector{CSSPath: "select#month"},
		Value:  "January",
	})
	if err != nil {
		t.Fatal(err)
	}
}

func TestStepSaveValue(t *testing.T) {
	fd := &fakeDriver{}
	// GetText returns "text" in fakeDriver
	err := runStep(t, fd, config.Step{
		Method: "save_value",
		Target: &config.Selector{CSSPath: "#result"},
		Value:  "my_result",
	})
	if err != nil {
		t.Fatal(err)
	}
}

func TestStepGetForm(t *testing.T) {
	fd := &fakeDriver{}
	if err := runStep(t, fd, config.Step{Method: "get_form"}); err != nil {
		t.Fatal(err)
	}
}

func TestStepExecuteScript(t *testing.T) {
	fd := &fakeDriver{}
	err := runStep(t, fd, config.Step{
		Method: "execute_script",
		Script: "window.scrollTo(0,0)",
	})
	if err != nil {
		t.Fatal(err)
	}
}

func TestStepExecuteScript_Empty(t *testing.T) {
	fd := &fakeDriver{}
	if err := runStep(t, fd, config.Step{Method: "execute_script"}); err == nil {
		t.Error("expected error for empty script, got nil")
	}
}

func TestStepSavePageHTML(t *testing.T) {
	fd := &fakeDriver{}
	err := runStep(t, fd, config.Step{
		Method: "save_page_html",
		Value:  "page_content",
	})
	if err != nil {
		t.Fatal(err)
	}
}

func TestStepWait(t *testing.T) {
	fd := &fakeDriver{}
	err := runStep(t, fd, config.Step{
		Method:  "wait",
		Target:  &config.Selector{CSSPath: "#loaded"},
		Timeout: 1,
	})
	if err != nil {
		t.Fatal(err)
	}
}

func TestExecutor_Expects_BodyIncludes_Pass(t *testing.T) {
	fd := &fakeDriver{body: "Welcome"}
	inst := makeInst(
		config.Action{
			Key:       "start",
			Starting:  true,
			OnSuccess: "end",
			OnFailure: "fail",
			Expects:   &config.Expects{Type: "body_includes", Value: "Welcome"},
			Steps:     []config.Step{{Method: "go", Value: "https://example.com"}},
		},
		config.Action{Key: "end", Steps: []config.Step{}},
		config.Action{Key: "fail", Steps: []config.Step{}},
	)
	ex := newTestExecutor(inst, fd)
	if err := ex.Run(); err != nil {
		t.Fatal(err)
	}
	// Verify "end" was visited (not "fail")
	if !ex.visited["end"] {
		t.Error("expected 'end' action to be visited (expects passed)")
	}
	if ex.visited["fail"] {
		t.Error("'fail' action should not have been visited")
	}
}

func TestExecutor_Expects_URLIncludes(t *testing.T) {
	fd := &fakeDriver{currentURL: "https://example.com/dashboard"}
	inst := makeInst(
		config.Action{
			Key:       "start",
			Starting:  true,
			OnSuccess: "done",
			OnFailure: "fail",
			Expects:   &config.Expects{Type: "url_includes", Value: "dashboard"},
			Steps:     []config.Step{{Method: "go", Value: "https://example.com"}},
		},
		config.Action{Key: "done", Steps: []config.Step{}},
		config.Action{Key: "fail", Steps: []config.Step{}},
	)
	ex := newTestExecutor(inst, fd)
	if err := ex.Run(); err != nil {
		t.Fatal(err)
	}
	if !ex.visited["done"] {
		t.Error("expected 'done' action to be visited")
	}
}

func TestExecutor_Expects_ElementPresent(t *testing.T) {
	fd := &fakeDriver{}
	inst := makeInst(
		config.Action{
			Key:       "start",
			Starting:  true,
			OnSuccess: "done",
			Expects: &config.Expects{
				Type:   "element_present",
				Target: &config.Selector{CSSPath: "#ok"},
			},
			Steps: []config.Step{{Method: "go", Value: "https://example.com"}},
		},
		config.Action{Key: "done", Steps: []config.Step{}},
	)
	ex := newTestExecutor(inst, fd)
	if err := ex.Run(); err != nil {
		t.Fatal(err)
	}
}

func TestStepFormatSavedValue(t *testing.T) {
	fd := &fakeDriver{}
	// Save a variable first, then format it.
	inst := makeInst(config.Action{
		Key:      "start",
		Starting: true,
		Steps: []config.Step{
			{Method: "save_value", Target: &config.Selector{CSSPath: "#el"}, Value: "raw"},
			{Method: "format_saved_value", Value: "formatted", Pattern: "prefix-{{raw}}"},
		},
	})
	ex := newTestExecutor(inst, fd)
	if err := ex.Run(); err != nil {
		t.Fatal(err)
	}
}

func TestStepFormatSavedValue_NoKey(t *testing.T) {
	fd := &fakeDriver{}
	if err := runStep(t, fd, config.Step{Method: "format_saved_value"}); err == nil {
		t.Error("expected error for missing value (key), got nil")
	}
}

func TestStepWaitForDownload_Timeout(t *testing.T) {
	fd := &fakeDriver{}
	inst := makeInst(config.Action{
		Key:      "start",
		Starting: true,
		Steps: []config.Step{{
			Method:  "wait_for_download",
			Pattern: "*.csv",
			Timeout: 1, // 1 second, /tmp won't have a .csv
		}},
	})
	// DownloadDir defaults to "" which means the watcher will try to scan ""
	// Override with a known-empty temp dir
	inst.DownloadDir = t.TempDir()
	ex := newTestExecutor(inst, fd)
	// Should fail (timeout) since no file appears
	if err := ex.Run(); err == nil {
		t.Error("expected timeout error from wait_for_download, got nil")
	}
}

func TestStepSubmit_NoTarget(t *testing.T) {
	fd := &fakeDriver{}
	if err := runStep(t, fd, config.Step{Method: "submit"}); err == nil {
		t.Error("expected error for missing target, got nil")
	}
}

func TestStepSelect_NoTarget(t *testing.T) {
	fd := &fakeDriver{}
	if err := runStep(t, fd, config.Step{Method: "select"}); err == nil {
		t.Error("expected error for missing target, got nil")
	}
}

func TestStepSaveValue_NoTarget(t *testing.T) {
	fd := &fakeDriver{}
	if err := runStep(t, fd, config.Step{Method: "save_value"}); err == nil {
		t.Error("expected error for missing target, got nil")
	}
}

func TestStepWait_NoTarget(t *testing.T) {
	fd := &fakeDriver{}
	if err := runStep(t, fd, config.Step{Method: "wait"}); err == nil {
		t.Error("expected error for missing target, got nil")
	}
}

func TestStepSavePageHTML_NoKey(t *testing.T) {
	fd := &fakeDriver{}
	if err := runStep(t, fd, config.Step{Method: "save_page_html"}); err == nil {
		t.Error("expected error for missing value key, got nil")
	}
}

func TestStepFillInInput_NoTarget(t *testing.T) {
	fd := &fakeDriver{}
	if err := runStep(t, fd, config.Step{Method: "fill_in_input"}); err == nil {
		t.Error("expected error for missing target, got nil")
	}
}

func TestExecutor_NoStartingAction(t *testing.T) {
	fd := &fakeDriver{}
	inst := makeInst() // no actions
	ex := newTestExecutor(inst, fd)
	if err := ex.Run(); err == nil {
		t.Error("expected error for missing starting action, got nil")
	}
}

func TestExecutor_Expects_UnknownType(t *testing.T) {
	fd := &fakeDriver{}
	inst := makeInst(
		config.Action{
			Key:       "start",
			Starting:  true,
			OnSuccess: "pass",
			OnFailure: "fail",
			Expects:   &config.Expects{Type: "unknown_type"},
			Steps:     []config.Step{{Method: "go", Value: "https://example.com"}},
		},
		config.Action{Key: "pass", Steps: []config.Step{}},
		config.Action{Key: "fail", Steps: []config.Step{}},
	)
	ex := newTestExecutor(inst, fd)
	if err := ex.Run(); err != nil {
		t.Fatal(err)
	}
	if ex.visited["pass"] {
		t.Error("unknown expects type should evaluate to false → fail branch")
	}
}

func TestExecutor_Expects_ElementPresent_NilTarget(t *testing.T) {
	fd := &fakeDriver{}
	inst := makeInst(
		config.Action{
			Key:       "start",
			Starting:  true,
			OnSuccess: "pass",
			OnFailure: "fail",
			Expects:   &config.Expects{Type: "element_present", Target: nil},
			Steps:     []config.Step{{Method: "go", Value: "https://example.com"}},
		},
		config.Action{Key: "pass"},
		config.Action{Key: "fail"},
	)
	ex := newTestExecutor(inst, fd)
	if err := ex.Run(); err != nil {
		t.Fatal(err)
	}
	if ex.visited["pass"] {
		t.Error("nil target should evaluate element_present as false")
	}
}

func TestExecutor_Expects_DownloadExists_Timeout(t *testing.T) {
	fd := &fakeDriver{}
	inst := makeInst(
		config.Action{
			Key:       "start",
			Starting:  true,
			OnSuccess: "pass",
			OnFailure: "fail",
			Expects:   &config.Expects{Type: "download_exists", Pattern: "*.csv"},
			Steps:     []config.Step{{Method: "go", Value: "https://example.com"}},
		},
		config.Action{Key: "pass"},
		config.Action{Key: "fail"},
	)
	inst.DownloadDir = t.TempDir() // empty dir, no file
	ex := newTestExecutor(inst, fd)
	if err := ex.Run(); err != nil {
		t.Fatal(err)
	}
	if ex.visited["pass"] {
		t.Error("empty dir should make download_exists false → fail branch")
	}
}

func TestExecutor_Expects_BodyIncludes_NoMatch(t *testing.T) {
	fd := &fakeDriver{body: "Welcome"}
	inst := makeInst(
		config.Action{
			Key:       "start",
			Starting:  true,
			OnSuccess: "pass",
			OnFailure: "fail",
			Expects:   &config.Expects{Type: "body_includes", Value: "NotPresent"},
			Steps:     []config.Step{{Method: "go", Value: "https://example.com"}},
		},
		config.Action{Key: "pass"},
		config.Action{Key: "fail"},
	)
	ex := newTestExecutor(inst, fd)
	if err := ex.Run(); err != nil {
		t.Fatal(err)
	}
	if ex.visited["pass"] {
		t.Error("non-matching body should route to fail")
	}
}

func TestExecutor_OnFailure_Routing(t *testing.T) {
	fd := &fakeDriver{}
	inst := makeInst(
		config.Action{
			Key:       "start",
			Starting:  true,
			OnFailure: "recover",
			Steps:     []config.Step{{Method: "go", Value: ""}}, // empty URL → error
		},
		config.Action{Key: "recover", Steps: []config.Step{}},
	)
	ex := newTestExecutor(inst, fd)
	if err := ex.Run(); err != nil {
		t.Fatalf("expected OnFailure routing, got error: %v", err)
	}
	if !ex.visited["recover"] {
		t.Error("expected 'recover' action to be visited after failure")
	}
}
