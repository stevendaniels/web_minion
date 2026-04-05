package http

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stevendaniels/web_minion/go-webminion/driver"
)

func newTestServer(body string) *httptest.Server {
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(body))
	}))
}

func TestHTTPDriver_Navigate_BodyContains(t *testing.T) {
	server := newTestServer(`<html><body><p>Test Link</p></body></html>`)
	defer server.Close()

	h := New()
	if err := h.Navigate(server.URL); err != nil {
		t.Fatalf("Navigate error: %v", err)
	}

	ok, err := h.BodyContains("Test Link")
	if err != nil {
		t.Fatalf("BodyContains error: %v", err)
	}
	if !ok {
		t.Error("expected body to contain 'Test Link'")
	}
}

func TestHTTPDriver_GetText(t *testing.T) {
	server := newTestServer(`<html><body><h1 id="title">Hello</h1></body></html>`)
	defer server.Close()

	h := New()
	if err := h.Navigate(server.URL); err != nil {
		t.Fatal(err)
	}

	text, err := h.GetText(driver.ResolvedSelector{Strategy: "css", Value: "#title"})
	if err != nil {
		t.Fatal(err)
	}
	if text != "Hello" {
		t.Errorf("expected 'Hello', got '%s'", text)
	}
}

func TestHTTPDriver_CurrentURL(t *testing.T) {
	server := newTestServer(`<html><body></body></html>`)
	defer server.Close()

	h := New()
	if err := h.Navigate(server.URL); err != nil {
		t.Fatal(err)
	}

	url, err := h.CurrentURL()
	if err != nil {
		t.Fatal(err)
	}
	if url != server.URL {
		t.Errorf("expected '%s', got '%s'", server.URL, url)
	}
}

func TestHTTPDriver_EvalNotSupported(t *testing.T) {
	h := New()
	if err := h.Eval("console.log('test')"); err == nil {
		t.Error("expected error for Eval, got nil")
	}
}

func TestHTTPDriver_ImplementsInterface(t *testing.T) {
	var _ driver.Driver = New()
}
