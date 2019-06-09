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
}

// GetConfig 获取配置数据
func GetConfig() *Config {
	return conf
}
