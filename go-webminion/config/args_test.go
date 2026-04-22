package config

import (
	"reflect"
	"testing"
)

func TestRequiredArgs(t *testing.T) {
	tests := []struct {
		name string
		inst *Config
		want []string
	}{
		{
			name: "no placeholders",
			inst: &Config{
				Flow: Flow{Actions: []Action{{
					Key: "a", Starting: true,
					Steps: []Step{{Method: "go", Value: "https://example.com", IsValidator: true}},
				}}},
			},
			want: []string{},
		},
		{
			name: "url placeholder with empty var declared",
			inst: &Config{
				Vars: map[string]string{"url": ""},
				Flow: Flow{Actions: []Action{{
					Key: "a", Starting: true,
					Steps: []Step{{Method: "go", Value: "{{url}}", IsValidator: true}},
				}}},
			},
			want: []string{"url"},
		},
		{
			name: "satisfied by non-empty var",
			inst: &Config{
				Vars: map[string]string{"url": "https://example.com"},
				Flow: Flow{Actions: []Action{{
					Key: "a", Starting: true,
					Steps: []Step{{Method: "go", Value: "{{url}}", IsValidator: true}},
				}}},
			},
			want: []string{},
		},
		{
			name: "built-ins and credentials excluded",
			inst: &Config{
				Flow: Flow{Actions: []Action{{
					Key: "a", Starting: true,
					Steps: []Step{{
						Method:      "go",
						Value:       "{{today}}/{{username}}",
						IsValidator: true,
					}},
				}}},
			},
			want: []string{},
		},
		{
			name: "multiple required args sorted",
			inst: &Config{
				Vars: map[string]string{"b": "", "a": ""},
				Flow: Flow{Actions: []Action{{
					Key: "a", Starting: true,
					Steps: []Step{{Method: "go", Value: "{{b}}/{{a}}", IsValidator: true}},
				}}},
			},
			want: []string{"a", "b"},
		},
		{
			name: "placeholder in selector",
			inst: &Config{
				Vars: map[string]string{"label": ""},
				Flow: Flow{Actions: []Action{{
					Key: "a", Starting: true,
					Steps: []Step{{
						Method:      "click",
						Target:      &Selector{AriaLabel: "{{label}}"},
						IsValidator: true,
					}},
				}}},
			},
			want: []string{"label"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := RequiredArgs(tt.inst)
			if got == nil {
				got = []string{}
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("RequiredArgs() = %v, want %v", got, tt.want)
			}
		})
	}
}
