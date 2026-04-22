package executor

import (
	"testing"
	"time"

	"github.com/stevendaniels/web_minion/go-webminion/config"
	"github.com/stevendaniels/web_minion/go-webminion/driver"
)

// fakeDriver implements driver.Driver for testing without a real browser.
type fakeDriver struct {
	navigated  []string
	clicked    []driver.ResolvedSelector
	filled     []fillCall
	body       string
	currentURL string
	pageHTML   string
}

type fillCall struct {
	sel   driver.ResolvedSelector
	value string
}

func (f *fakeDriver) Navigate(url string) error { f.navigated = append(f.navigated, url); return nil }
func (f *fakeDriver) Click(sel driver.ResolvedSelector) error {
	f.clicked = append(f.clicked, sel)
	return nil
}
func (f *fakeDriver) Fill(sel driver.ResolvedSelector, v string) error {
	f.filled = append(f.filled, fillCall{sel, v})
	return nil
}
func (f *fakeDriver) Submit(sel driver.ResolvedSelector) error                 { return nil }
func (f *fakeDriver) Select(sel driver.ResolvedSelector, v string) error       { return nil }
func (f *fakeDriver) GetText(sel driver.ResolvedSelector) (string, error)      { return "text", nil }
func (f *fakeDriver) WaitForSelector(sel driver.ResolvedSelector, t int) error { return nil }
func (f *fakeDriver) WaitForURL(pattern string, t int) error                   { return nil }
func (f *fakeDriver) BodyContains(text string) (bool, error)                   { return f.body == text, nil }
func (f *fakeDriver) CurrentURL() (string, error)                              { return f.currentURL, nil }
func (f *fakeDriver) PageHTML() (string, error) {
	if f.pageHTML != "" {
		return f.pageHTML, nil
	}
	return "<html/>", nil
}
func (f *fakeDriver) Eval(script string) error                                 { return nil }
func (f *fakeDriver) Close() error                                             { return nil }

func makeConfig(actions ...config.Action) *config.Config {
	return &config.Config{
		ID:          "test",
		DownloadDir: "/tmp",
		Flow:        config.Flow{Actions: actions},
	}
}

func newTestExecutor(cfg *config.Config, d driver.Driver) *Executor {
	return New(cfg, d, nil, time.Now(), time.Now(), nil)
}

func TestExecutor_Run_PlaceholderArg(t *testing.T) {
	fd := &fakeDriver{}
	cfg := &config.Config{
		ID:   "test",
		Vars: map[string]string{"url": ""},
		Flow: config.Flow{Actions: []config.Action{{
			Key:      "start",
			Starting: true,
			Steps:    []config.Step{{Method: "go", Value: "{{url}}"}},
		}}},
	}

	ex := New(cfg, fd, nil, time.Now(), time.Now(), map[string]string{"url": "https://example.com"})
	if err := ex.Run(); err != nil {
		t.Fatalf("Run error: %v", err)
	}
	if len(fd.navigated) != 1 || fd.navigated[0] != "https://example.com" {
		t.Errorf("expected Navigate(https://example.com), got %v", fd.navigated)
	}
}

func TestExecutor_Run_SingleAction(t *testing.T) {
	fd := &fakeDriver{}
	cfg := makeConfig(config.Action{
		Key:      "start",
		Starting: true,
		Steps:    []config.Step{{Method: "go", Value: "https://example.com"}},
	})

	ex := newTestExecutor(cfg, fd)
	if err := ex.Run(); err != nil {
		t.Fatalf("Run error: %v", err)
	}
	if len(fd.navigated) != 1 || fd.navigated[0] != "https://example.com" {
		t.Errorf("expected Navigate to be called with example.com, got %v", fd.navigated)
	}
}

func TestExecutor_Run_ChainedActions(t *testing.T) {
	fd := &fakeDriver{}
	cfg := makeConfig(
		config.Action{
			Key:       "start",
			Starting:  true,
			OnSuccess: "second",
			Steps:     []config.Step{{Method: "go", Value: "https://a.com"}},
		},
		config.Action{
			Key:   "second",
			Steps: []config.Step{{Method: "go", Value: "https://b.com"}},
		},
	)

	ex := newTestExecutor(cfg, fd)
	if err := ex.Run(); err != nil {
		t.Fatalf("Run error: %v", err)
	}
	if len(fd.navigated) != 2 {
		t.Errorf("expected 2 navigations, got %d", len(fd.navigated))
	}
}

func TestExecutor_Run_CycleDetection(t *testing.T) {
	fd := &fakeDriver{}
	cfg := makeConfig(
		config.Action{
			Key:       "start",
			Starting:  true,
			OnSuccess: "start", // cycle back to self
			Steps:     []config.Step{{Method: "go", Value: "https://example.com"}},
		},
	)

	ex := newTestExecutor(cfg, fd)
	err := ex.Run()
	if err == nil {
		t.Error("expected cycle detection error, got nil")
	}
}

func TestExecutor_Run_UnknownStepMethod(t *testing.T) {
	fd := &fakeDriver{}
	cfg := makeConfig(config.Action{
		Key:      "start",
		Starting: true,
		Steps:    []config.Step{{Method: "no_such_method"}},
	})

	ex := newTestExecutor(cfg, fd)
	if err := ex.Run(); err == nil {
		t.Error("expected error for unknown step method, got nil")
	}
}

func TestStepRegistry_HasAllMethods(t *testing.T) {
	registry := newStepRegistry()
	expected := []string{
		"go", "get_form", "fill_in_input", "click", "submit",
		"select", "save_value", "format_saved_value", "wait_for_download",
		"wait", "execute_script", "save_page_html",
		"write_file", "html_to_markdown", "wait_for_reply", "extract_readable",
	}
	for _, method := range expected {
		if _, ok := registry[method]; !ok {
			t.Errorf("missing step handler for method: %s", method)
		}
	}
}
