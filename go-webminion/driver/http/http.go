package http

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/PuerkitoBio/goquery"

	"github.com/stevendaniels/web_minion/go-webminion/driver"
)

// HTTPDriver implements driver.Driver using net/http and goquery.
// It handles static pages and simple form submission; JavaScript is not supported.
type HTTPDriver struct {
	client  *http.Client
	doc     *goquery.Document
	lastURL string
}

func New() *HTTPDriver {
	return &HTTPDriver{
		client: &http.Client{},
	}
}

func (h *HTTPDriver) Navigate(url string) error {
	resp, err := h.client.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		return err
	}

	h.doc = doc
	h.lastURL = url
	return nil
}

func (h *HTTPDriver) ensureDoc() error {
	if h.doc == nil {
		return fmt.Errorf("no page loaded; call Navigate first")
	}
	return nil
}

func (h *HTTPDriver) Click(sel driver.ResolvedSelector) error {
	return fmt.Errorf("click not supported by HTTP driver (no JavaScript)")
}

func (h *HTTPDriver) Fill(sel driver.ResolvedSelector, value string) error {
	if err := h.ensureDoc(); err != nil {
		return err
	}
	found := h.doc.Find(sel.Value)
	if found.Length() == 0 {
		return &driver.NotFoundError{Selector: sel}
	}
	found.SetAttr("value", value)
	return nil
}

func (h *HTTPDriver) Submit(sel driver.ResolvedSelector) error {
	return fmt.Errorf("submit not fully implemented in HTTP driver")
}

func (h *HTTPDriver) Select(sel driver.ResolvedSelector, value string) error {
	if err := h.ensureDoc(); err != nil {
		return err
	}
	h.doc.Find(sel.Value + " option").Each(func(_ int, s *goquery.Selection) {
		if s.Text() == value || s.AttrOr("value", "") == value {
			s.SetAttr("selected", "selected")
		} else {
			s.RemoveAttr("selected")
		}
	})
	return nil
}

func (h *HTTPDriver) GetText(sel driver.ResolvedSelector) (string, error) {
	if err := h.ensureDoc(); err != nil {
		return "", err
	}
	return h.doc.Find(sel.Value).Text(), nil
}

func (h *HTTPDriver) WaitForSelector(sel driver.ResolvedSelector, timeoutSecs int) error {
	if err := h.ensureDoc(); err != nil {
		return err
	}
	if h.doc.Find(sel.Value).Length() == 0 {
		return &driver.NotFoundError{Selector: sel}
	}
	return nil
}

func (h *HTTPDriver) WaitForURL(pattern string, timeoutSecs int) error {
	for _, part := range strings.Split(pattern, "|") {
		if strings.Contains(h.lastURL, part) {
			return nil
		}
	}
	return fmt.Errorf("current URL %q does not match pattern %q", h.lastURL, pattern)
}

func (h *HTTPDriver) BodyContains(text string) (bool, error) {
	if err := h.ensureDoc(); err != nil {
		return false, err
	}
	return strings.Contains(h.doc.Text(), text), nil
}

func (h *HTTPDriver) CurrentURL() (string, error) {
	return h.lastURL, nil
}

func (h *HTTPDriver) PageHTML() (string, error) {
	if err := h.ensureDoc(); err != nil {
		return "", err
	}
	html, err := h.doc.Html()
	if err != nil {
		return "", err
	}
	return html, nil
}

func (h *HTTPDriver) Eval(script string) error {
	return fmt.Errorf("Eval not supported by HTTP driver (no JavaScript)")
}

func (h *HTTPDriver) Close() error {
	return nil
}
