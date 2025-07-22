package internal

import (
	"github.com/wcy-dt/ponghub/protos/defaultConfig"
	"log"
	"os"

	"gopkg.in/yaml.v3"
)

// ServiceConfig defines the configuration for a service, including its health and API ports
type ServiceConfig struct {
	Name    string       `yaml:"name"`
	Health  []PortConfig `yaml:"health"`
	API     []PortConfig `yaml:"api"`
	Timeout int          `yaml:"timeout,omitempty"`
	Retry   int          `yaml:"retry,omitempty"`
}

// PortConfig defines the configuration for a port
type PortConfig struct {
	URL           string `yaml:"url"`
	Method        string `yaml:"method,omitempty"`
	Body          string `yaml:"body,omitempty"`
	StatusCode    int    `yaml:"status_code,omitempty"`
	ResponseRegex string `yaml:"response_regex,omitempty"`
}

// Config defines the overall configuration structure for the application
type Config struct {
	Services   []ServiceConfig `yaml:"services"`
	Timeout    int             `yaml:"timeout,omitempty"`
	Retry      int             `yaml:"retry,omitempty"`
	MaxLogDays int             `yaml:"max_log_days,omitempty"`
}

// SetDefaultFields sets default values for the configuration fields
func SetDefaultFields(cfg *Config) {
	defaultConfig.SetDefaultTimeout(&cfg.Timeout)
	defaultConfig.SetDefaultRetry(&cfg.Retry)
	defaultConfig.SetDefaultMaxLogDays(&cfg.MaxLogDays)

	for i := range cfg.Services {
		defaultConfig.SetDefaultTimeout(&cfg.Services[i].Timeout)
		defaultConfig.SetDefaultRetry(&cfg.Services[i].Retry)
	}
}

// LoadConfig loads the configuration from a YAML file at the specified path
func LoadConfig(path string) (*Config, error) {
	// Read the configuration file
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer func(f *os.File) {
		if err := f.Close(); err != nil {
			log.Println("Error closing config file:", err)
		}
	}(f)

	// Decode the YAML configuration
	cfg := new(Config)
	decoder := yaml.NewDecoder(f)
	if err := decoder.Decode(cfg); err != nil {
		log.Fatalln("Failed to decode YAML config:", err)
	}
	// Set default values for the configuration
	SetDefaultFields(cfg)

	if len(cfg.Services) == 0 {
		log.Fatalln("No services defined in the configuration file")
	}
	return cfg, nil
}
