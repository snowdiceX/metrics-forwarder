package config

import (
	"encoding/json"
	"io/ioutil"
)

var conf = &Config{}

// Config data
type Config struct {

	// LogConfigPath log config file path
	LogConfigPath string `json:"log,omitempty"`

	// ConfigPath config file path
	ConfigPath string `json:"config,omitempty"`

	// Push URL
	Push string `json:"push,omitempty"`

	// Jobs pushing the metrics data to pushgateway
	Jobs []*JobConfig `json:"jobs,omitempty"`
}

// JobConfig config of the push job
type JobConfig struct {

	// Name label
	Name string `json:"name,omitempty"`

	// Pull URL
	Pull string `json:"pull,omitempty"`

	// Zone label
	Zone string `json:"zone,omitempty"`

	// Host label
	Host string `json:"host,omitempty"`

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

// Load the config file
func (c *Config) Load() error {
	bytes, err := ioutil.ReadFile(c.ConfigPath)
	if err != nil {
		return err
	}
	return createConfig(bytes, c)
}

// createConfig parse config
func createConfig(bytes []byte, conf *Config) error {
	err := json.Unmarshal(bytes, conf)
	if err != nil {
		return err
	}
	return nil
}
