package chrome

import (
	"fmt"
	"strings"
	"time"

	"github.com/go-rod/rod"
	"github.com/go-rod/rod/lib/input"
	"github.com/go-rod/rod/lib/proto"

	"github.com/stevendaniels/web_minion/go-webminion/driver"
)

// ChromeDriver implements driver.Driver using go-rod.
type ChromeDriver struct {
	browser  *rod.Browser
	page     *rod.Page
	headless bool
}

// New creates a ChromeDriver. Call Navigate to open the first page.
func New(headless bool) *ChromeDriver {
	return &ChromeDriver{headless: headless}
}

func (c *ChromeDriver) ensureBrowser() {
	if c.browser == nil {
		c.browser = rod.New().MustConnect()
	}
}

func (c *ChromeDriver) Navigate(url string) error {
	c.ensureBrowser()
	if c.page == nil {
		c.page = c.browser.MustPage("")
	}
	return c.page.Navigate(url)
}

func (c *ChromeDriver) element(sel driver.ResolvedSelector) (*rod.Element, error) {
	css, err := c.toCSS(sel)
	if err != nil {
		return nil, err
	}
	elem, err := c.page.Element(css)
	if err != nil {
		return nil, &driver.NotFoundError{Selector: sel}
	}
	return elem, nil
}

// toCSS converts a ResolvedSelector to a CSS selector string go-rod can use.
func (c *ChromeDriver) toCSS(sel driver.ResolvedSelector) (string, error) {
	switch sel.Strategy {
	case "css":
		return sel.Value, nil
	case "aria-label":
		return fmt.Sprintf("[aria-label=%q]", sel.Value), nil
	case "aria-role":
		return fmt.Sprintf("[role=%q]", sel.Value), nil
	case "text":
		// go-rod supports :contains() via Rod's Has method; use attribute fallback here
		return fmt.Sprintf("//*[contains(text(),%q)]", sel.Value), nil
	default:
		return "", fmt.Errorf("unsupported selector strategy: %s", sel.Strategy)
	}
}

func (c *ChromeDriver) Click(sel driver.ResolvedSelector) error {
	elem, err := c.element(sel)
	if err != nil {
		return err
	}
	return elem.Click(proto.InputMouseButtonLeft, 1)
}

func (c *ChromeDriver) Fill(sel driver.ResolvedSelector, value string) error {
	elem, err := c.element(sel)
	if err != nil {
		return err
	}
	return elem.Input(value)
}

func (c *ChromeDriver) Submit(sel driver.ResolvedSelector) error {
	elem, err := c.element(sel)
	if err != nil {
		return err
	}
	return elem.Type(input.Enter)
}

func (c *ChromeDriver) Select(sel driver.ResolvedSelector, value string) error {
	elem, err := c.element(sel)
	if err != nil {
		return err
	}
	return elem.Select([]string{value}, true, rod.SelectorTypeText)
}

func (c *ChromeDriver) GetText(sel driver.ResolvedSelector) (string, error) {
	elem, err := c.element(sel)
	if err != nil {
		return "", err
	}
	return elem.Text()
}

func (c *ChromeDriver) WaitForSelector(sel driver.ResolvedSelector, timeoutSecs int) error {
	css, err := c.toCSS(sel)
	if err != nil {
		return err
	}
	page := c.page.Timeout(time.Duration(timeoutSecs) * time.Second)
	_, err = page.Element(css)
	if err != nil {
		return &driver.TimeoutError{Op: css, Timeout: timeoutSecs}
	}
	return nil
}

func (c *ChromeDriver) WaitForURL(pattern string, timeoutSecs int) error {
	deadline := time.Now().Add(time.Duration(timeoutSecs) * time.Second)
	for time.Now().Before(deadline) {
		url, err := c.CurrentURL()
		if err != nil {
			return err
		}
		for _, part := range strings.Split(pattern, "|") {
			if strings.Contains(url, part) {
				return nil
			}
		}
		time.Sleep(500 * time.Millisecond)
	}
	return &driver.TimeoutError{Op: "url contains " + pattern, Timeout: timeoutSecs}
}

func (c *ChromeDriver) BodyContains(text string) (bool, error) {
	body, err := c.page.Element("body")
	if err != nil {
		return false, err
	}
	bodyText, err := body.Text()
	if err != nil {
		return false, err
	}
	return strings.Contains(bodyText, text), nil
}

func (c *ChromeDriver) CurrentURL() (string, error) {
	info, err := c.page.Info()
	if err != nil {
		return "", err
	}
	return info.URL, nil
}

func (c *ChromeDriver) PageHTML() (string, error) {
	return c.page.HTML()
}

func (c *ChromeDriver) Eval(script string) error {
	_, err := c.page.Eval(script)
	return err
}

func (c *ChromeDriver) Close() error {
	if c.browser != nil {
		return c.browser.Close()
	}
	return nil
}
