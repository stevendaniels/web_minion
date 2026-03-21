package config

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"gopkg.in/yaml.v3"
)

// LoadConfig reads and parses a config file, auto-detecting YAML vs JSON.
func LoadConfig(path string) (*Institution, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var inst Institution
	if looksLikeJSON(data) {
		if err := json.Unmarshal(data, &inst); err != nil {
			return nil, err
		}
	} else {
		if err := yaml.Unmarshal(data, &inst); err != nil {
			return nil, err
		}
	}

	if errs := ValidateConfig(&inst); len(errs) > 0 {
		msgs := make([]string, len(errs))
		for i, e := range errs {
			msgs[i] = e.Error()
		}
		return nil, fmt.Errorf("invalid config: %s", strings.Join(msgs, "; "))
	}

	return &inst, nil
}

func looksLikeJSON(data []byte) bool {
	var js json.RawMessage
	return json.Unmarshal(bytes.TrimSpace(data), &js) == nil
}
