package driver

import (
	"testing"

	"github.com/stevendaniels/web_minion/go-webminion/config"
)

func TestResolve_CSSPath(t *testing.T) {
	sel, err := Resolve(&config.Selector{CSSPath: ".my-class"})
	if err != nil {
		t.Fatal(err)
	}
	if sel.Strategy != "css" || sel.Value != ".my-class" {
		t.Errorf("unexpected: %+v", sel)
	}
}

func TestResolve_ID(t *testing.T) {
	sel, err := Resolve(&config.Selector{ID: "submit-btn"})
	if err != nil {
		t.Fatal(err)
	}
	if sel.Strategy != "css" || sel.Value != "#submit-btn" {
		t.Errorf("unexpected: %+v", sel)
	}
}

func TestResolve_AriaLabel(t *testing.T) {
	sel, err := Resolve(&config.Selector{AriaLabel: "Search"})
	if err != nil {
		t.Fatal(err)
	}
	if sel.Strategy != "aria-label" || sel.Value != "Search" {
		t.Errorf("unexpected: %+v", sel)
	}
}

func TestResolve_Nil(t *testing.T) {
	_, err := Resolve(nil)
	if err == nil {
		t.Error("expected error for nil selector")
	}
}

func TestResolve_Empty(t *testing.T) {
	_, err := Resolve(&config.Selector{})
	if err == nil {
		t.Error("expected error for empty selector")
	}
}

func TestDriverInterface(t *testing.T) {
	// Compile-time check that Driver interface is well-formed.
	var _ Driver = (*mockDriver)(nil)
}

// mockDriver satisfies the Driver interface for compile-time verification.
type mockDriver struct{}

func (m *mockDriver) Navigate(url string) error                             { return nil }
func (m *mockDriver) Click(sel ResolvedSelector) error                      { return nil }
func (m *mockDriver) Fill(sel ResolvedSelector, value string) error         { return nil }
func (m *mockDriver) Submit(sel ResolvedSelector) error                     { return nil }
func (m *mockDriver) Select(sel ResolvedSelector, value string) error       { return nil }
func (m *mockDriver) GetText(sel ResolvedSelector) (string, error)          { return "", nil }
func (m *mockDriver) WaitForSelector(sel ResolvedSelector, t int) error     { return nil }
func (m *mockDriver) WaitForURL(pattern string, t int) error                { return nil }
func (m *mockDriver) BodyContains(text string) (bool, error)                { return false, nil }
func (m *mockDriver) CurrentURL() (string, error)                           { return "", nil }
func (m *mockDriver) PageHTML() (string, error)                             { return "", nil }
func (m *mockDriver) Eval(script string) error                              { return nil }
func (m *mockDriver) Close() error                                          { return nil }
