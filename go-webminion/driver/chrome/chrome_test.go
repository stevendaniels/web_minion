//go:build integration

package chrome

import (
	"testing"

	"github.com/stevendaniels/web_minion/go-webminion/driver"
)

func TestChromeDriver_Navigate(t *testing.T) {
	ch := New(true)
	defer ch.Close()

	if err := ch.Navigate("https://example.com"); err != nil {
		t.Fatalf("Navigate error: %v", err)
	}

	ok, err := ch.BodyContains("Example Domain")
	if err != nil {
		t.Fatalf("BodyContains error: %v", err)
	}
	if !ok {
		t.Error("expected page to contain 'Example Domain'")
	}
}

func TestChromeDriver_ImplementsInterface(t *testing.T) {
	var _ driver.Driver = New(true)
}
