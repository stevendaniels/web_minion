package watcher

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"
)

func TestWatch(t *testing.T) {
	tmpDir := t.TempDir()
	pattern := "*.csv"
	timeout := 2 * time.Second

	done := make(chan string, 1)
	errCh := make(chan error, 1)
	go func() {
		path, err := Watch(tmpDir, pattern, timeout)
		if err != nil {
			errCh <- err
			return
		}
		done <- path
	}()

	time.Sleep(200 * time.Millisecond)

	csvPath := filepath.Join(tmpDir, "test.csv")
	if err := os.WriteFile(csvPath, []byte("test,data"), 0644); err != nil {
		t.Fatal(err)
	}

	select {
	case path := <-done:
		if !strings.HasSuffix(path, "test.csv") {
			t.Errorf("expected path ending with 'test.csv', got '%s'", path)
		}
	case err := <-errCh:
		t.Fatalf("Watch error: %v", err)
	case <-time.After(3 * time.Second):
		t.Error("watcher timed out")
	}
}

func TestWatch_Timeout(t *testing.T) {
	tmpDir := t.TempDir()
	_, err := Watch(tmpDir, "*.csv", 200*time.Millisecond)
	if err == nil {
		t.Error("expected timeout error, got nil")
	}
}

func TestWatch_SkipsCrdownload(t *testing.T) {
	tmpDir := t.TempDir()

	// Create a .crdownload file — should be ignored
	crPath := filepath.Join(tmpDir, "download.csv.crdownload")
	if err := os.WriteFile(crPath, []byte("partial"), 0644); err != nil {
		t.Fatal(err)
	}

	// Watcher should time out because no real .csv exists
	_, err := Watch(tmpDir, "*.csv", 300*time.Millisecond)
	if err == nil {
		t.Error("expected timeout (crdownload file should be ignored), got nil")
	}
}
