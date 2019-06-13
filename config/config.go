package config

var conf = &Config{}

// Config data
type Config struct {

	// LogConfigPath log config file path
	LogConfigPath string `json:"log,omitempty"`

	// Pull URL
	Pull string `json:"pull,omitempty"`

	// Push url
	Push string `json:"push,omitempty"`

	// Zone label
	Zone string `json:"zone,omitempty"`

	// Host label
	Host string `json:"host,omitempty"`

	// Job label
	Job string `json:"job,omitempty"`

	// Group label
	Group string `json:"group,omitempty"`

	// Group value
	GroupValue string `json:"group_value,omitempty"`

	// Ticker of time for pull and push job
	Ticker uint32 `json:"ticker,omitempty"`
}

// GetConfig return a config
func GetConfig() *Config {
	return conf
}
