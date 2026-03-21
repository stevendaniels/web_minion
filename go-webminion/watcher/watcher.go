package watcher

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"
)

// Watch blocks until a file matching pattern appears in dir, then returns its path.
// Files with a .crdownload suffix are ignored (Chrome partial downloads).
// Returns an error if timeout elapses before a match is found.
func Watch(dir, pattern string, timeout time.Duration) (string, error) {
	ticker := time.NewTicker(500 * time.Millisecond)
	defer ticker.Stop()
	deadline := time.Now().Add(timeout)

	for {
		if match, err := scanDir(dir, pattern); err != nil {
			return "", err
		} else if match != "" {
			return match, nil
		}

		if time.Now().After(deadline) {
			return "", fmt.Errorf("timeout waiting for file matching %q in %s", pattern, dir)
		}

		<-ticker.C
	}
}

func scanDir(dir, pattern string) (string, error) {
	entries, err := os.ReadDir(dir)
	if err != nil {
		return "", fmt.Errorf("failed to read directory: %w", err)
	}
	for _, e := range entries {
		if e.IsDir() {
			continue
		}
		name := e.Name()
		if strings.HasSuffix(name, ".crdownload") {
			continue
		}
		if matched, _ := filepath.Match(pattern, name); matched {
			return filepath.Join(dir, name), nil
		}
	}
	return "", nil
}
