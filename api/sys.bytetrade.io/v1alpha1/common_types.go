package v1alpha1

// EnvVarSpec defines the common fields for environment variables
// This struct is embedded in SystemEnv, UserEnv, and AppEnvVar
type EnvVarSpec struct {
	EnvName     string `json:"envName" yaml:"envName" validate:"required"`
	Value       string `json:"value,omitempty" yaml:"value,omitempty"`
	Default     string `json:"default,omitempty" yaml:"default,omitempty"`
	Editable    bool   `json:"editable,omitempty" yaml:"editable,omitempty"`
	Type        string `json:"type,omitempty" yaml:"type,omitempty"`
	Required    bool   `json:"required,omitempty" yaml:"required,omitempty"`
	Title       string `json:"title,omitempty" yaml:"title,omitempty"`
	Description string `json:"description,omitempty" yaml:"description,omitempty"`
	// Options defines a finite set of allowed values for this env var
	Options []EnvValueOptionItem `json:"options,omitempty" yaml:"options,omitempty"`
	// +kubebuilder:validation:Pattern=`^https?://`
	// RemoteOptions provides a URL (http/https) returning a JSON-encoded string array of allowed values
	RemoteOptions string `json:"remoteOptions,omitempty" yaml:"remoteOptions,omitempty"`
	Regex         string `json:"regex,omitempty" yaml:"regex,omitempty"`
	// MultiSelect indicates the value may hold multiple items joined by Splitter.
	// When true, the value is split on Splitter and each item is validated
	// independently against Type/Options/RemoteOptions/Regex.
	MultiSelect bool `json:"multiSelect,omitempty" yaml:"multiSelect,omitempty"`
	// Splitter is the delimiter used to join/split multiple values when
	// MultiSelect is true. Defaults to "," (DefaultEnvSplitter) when empty.
	Splitter string `json:"splitter,omitempty" yaml:"splitter,omitempty"`
}

// DefaultEnvSplitter is the delimiter used to join/split multi-select env
// values when Splitter is not explicitly set.
const DefaultEnvSplitter = ","

// GetEffectiveValue returns the effective value of the environment variable.
// If Value is not empty, it returns Value; otherwise, it returns Default.
func (e *EnvVarSpec) GetEffectiveValue() string {
	if e.Value != "" {
		return e.Value
	}
	return e.Default
}

// GetSplitter returns the delimiter used to join/split multiple values for a
// multi-select env var, defaulting to DefaultEnvSplitter when unset.
func (e *EnvVarSpec) GetSplitter() string {
	if e.Splitter != "" {
		return e.Splitter
	}
	return DefaultEnvSplitter
}
