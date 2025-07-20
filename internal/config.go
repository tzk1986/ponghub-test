package internal

import (
	"os"

	"gopkg.in/yaml.v3"
)

type ServiceConfig struct {
	Name    string       `yaml:"name"`
	Health  []PortConfig `yaml:"health"`
	API     []PortConfig `yaml:"api"`
	Timeout int          `yaml:"timeout,omitempty"` // Units: seconds (service-level default)
	Retry   int          `yaml:"retry,omitempty"`
}

type PortConfig struct {
	URL           string `yaml:"url"`
	Method        string `yaml:"method,omitempty"`
	Body          string `yaml:"body,omitempty"`
	StatusCode    int    `yaml:"status_code,omitempty"`
	ResponseRegex string `yaml:"response_regex,omitempty"`
}

type Config struct {
	Services   []ServiceConfig `yaml:"services"`
	Timeout    int             `yaml:"timeout,omitempty"`      // Units: seconds (global default)
	Retry      int             `yaml:"retry,omitempty"`        // Global default retry count
	MaxLogDays int             `yaml:"max_log_days,omitempty"` // Maximum days to keep logs
}

func LoadConfig(path string) (*Config, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	var cfg Config
	decoder := yaml.NewDecoder(f)
	if err := decoder.Decode(&cfg); err != nil {
		return nil, err
	}
	// Default values for global settings
	if cfg.Timeout <= 0 {
		cfg.Timeout = 5 // Default timeout 5 seconds
	}
	if cfg.Retry <= 0 {
		cfg.Retry = 2 // Default retry count 2
	}
	if cfg.MaxLogDays <= 0 {
		cfg.MaxLogDays = 30 // Default max log days 30
	}
	for i := range cfg.Services {
		if cfg.Services[i].Timeout <= 0 {
			cfg.Services[i].Timeout = cfg.Timeout
		}
		if cfg.Services[i].Retry <= 0 {
			cfg.Services[i].Retry = cfg.Retry
		}
	}
	return &cfg, nil
}
