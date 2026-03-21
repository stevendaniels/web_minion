package config

// Institution is the top-level config for a single automation target.
type Institution struct {
	ID                  string `yaml:"id"                    json:"id"`
	Name                string `yaml:"name"                  json:"name"`
	BaseURL             string `yaml:"base_url"              json:"base_url"`
	ProfileDir          string `yaml:"profile_dir"           json:"profile_dir"`
	DownloadDir         string `yaml:"download_dir"          json:"download_dir"`
	ScreenshotOnFailure bool   `yaml:"screenshot_on_failure" json:"screenshot_on_failure"`
	MaxRetries          int    `yaml:"max_retries"           json:"max_retries"`
	Flow                Flow   `yaml:"flow"                  json:"flow"`
	Output              Output `yaml:"output"                json:"output"`
}

// Flow holds the ordered list of actions that make up an automation flow.
type Flow struct {
	Actions []Action `yaml:"actions" json:"actions"`
}

// Action is a named node in the flow graph.
type Action struct {
	Key       string   `yaml:"key"        json:"key"`
	Starting  bool     `yaml:"starting"   json:"starting"`
	Driver    string   `yaml:"driver"     json:"driver"`
	OnSuccess string   `yaml:"on_success" json:"on_success"`
	OnFailure string   `yaml:"on_failure" json:"on_failure"`
	Expects   *Expects `yaml:"expects"    json:"expects"`
	Steps     []Step   `yaml:"steps"      json:"steps"`
}

// Step is a single automation operation within an action.
type Step struct {
	Method        string     `yaml:"method"         json:"method"`
	Target        *Selector  `yaml:"target"         json:"target"`
	Targets       []Selector `yaml:"targets"        json:"targets"`
	Value         string     `yaml:"value"          json:"value"`
	Values        any        `yaml:"values"         json:"values"`
	Script        string     `yaml:"script"         json:"script"`
	Pattern       string     `yaml:"pattern"        json:"pattern"`
	Timeout       int        `yaml:"timeout"        json:"timeout"`
	RetainElement bool       `yaml:"retain_element" json:"retain_element"`
}

// Selector identifies a DOM element via one of several strategies.
type Selector struct {
	AriaLabel  string `yaml:"aria_label"  json:"aria_label"`
	AriaRole   string `yaml:"aria_role"   json:"aria_role"`
	DataTestID string `yaml:"data_testid" json:"data_testid"`
	ID         string `yaml:"id"          json:"id"`
	Name       string `yaml:"name"        json:"name"`
	Text       string `yaml:"text"        json:"text"`
	CSSPath    string `yaml:"css_path"    json:"css_path"`
}

// Expects defines the condition evaluated after an action's steps complete.
type Expects struct {
	Type    string    `yaml:"type"    json:"type"`
	Value   string    `yaml:"value"   json:"value"`
	Pattern string    `yaml:"pattern" json:"pattern"`
	Target  *Selector `yaml:"target"  json:"target"`
}

// Output configures how extracted data is mapped to output columns.
type Output struct {
	ColumnMap map[string]string `yaml:"column_map" json:"column_map"`
}
