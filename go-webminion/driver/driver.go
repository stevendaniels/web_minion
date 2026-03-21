package driver

import (
	"fmt"

	"github.com/stevendaniels/web_minion/go-webminion/config"
)

// ResolvedSelector is a selector after it has been converted to a concrete strategy + value.
type ResolvedSelector struct {
	Strategy string // "css" | "aria-label" | "aria-role" | "text" | "id" | "name" | "data-testid"
	Value    string
}

// NotFoundError is returned when a selector matches no element.
type NotFoundError struct {
	Selector ResolvedSelector
}

func (e *NotFoundError) Error() string {
	return fmt.Sprintf("element not found: %s=%q", e.Selector.Strategy, e.Selector.Value)
}

// TimeoutError is returned when a wait operation exceeds its deadline.
type TimeoutError struct {
	Op      string
	Timeout int
}

func (e *TimeoutError) Error() string {
	return fmt.Sprintf("timeout after %ds waiting for %s", e.Timeout, e.Op)
}

// Driver defines the browser automation operations that all driver implementations must support.
type Driver interface {
	Navigate(url string) error
	Click(sel ResolvedSelector) error
	Fill(sel ResolvedSelector, value string) error
	Submit(sel ResolvedSelector) error
	Select(sel ResolvedSelector, value string) error
	GetText(sel ResolvedSelector) (string, error)
	WaitForSelector(sel ResolvedSelector, timeoutSecs int) error
	WaitForURL(pattern string, timeoutSecs int) error
	BodyContains(text string) (bool, error)
	CurrentURL() (string, error)
	PageHTML() (string, error)
	Eval(script string) error
	Close() error
}

// Resolve converts a config.Selector into a ResolvedSelector using priority order:
// css_path > id > name > aria_label > aria_role > data_testid > text
func Resolve(sel *config.Selector) (ResolvedSelector, error) {
	if sel == nil {
		return ResolvedSelector{}, fmt.Errorf("selector is nil")
	}
	switch {
	case sel.CSSPath != "":
		return ResolvedSelector{Strategy: "css", Value: sel.CSSPath}, nil
	case sel.ID != "":
		return ResolvedSelector{Strategy: "css", Value: "#" + sel.ID}, nil
	case sel.Name != "":
		return ResolvedSelector{Strategy: "css", Value: fmt.Sprintf("[name=%q]", sel.Name)}, nil
	case sel.AriaLabel != "":
		return ResolvedSelector{Strategy: "aria-label", Value: sel.AriaLabel}, nil
	case sel.AriaRole != "":
		return ResolvedSelector{Strategy: "aria-role", Value: sel.AriaRole}, nil
	case sel.DataTestID != "":
		return ResolvedSelector{Strategy: "css", Value: fmt.Sprintf("[data-testid=%q]", sel.DataTestID)}, nil
	case sel.Text != "":
		return ResolvedSelector{Strategy: "text", Value: sel.Text}, nil
	}
	return ResolvedSelector{}, fmt.Errorf("selector has no matching strategy set")
}
