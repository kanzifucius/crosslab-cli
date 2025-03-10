package config

import (
	"fmt"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

// Provider represents a Crossplane provider configuration
type Provider struct {
	Name    string `yaml:"name"`
	Package string `yaml:"package"`
	Version string `yaml:"version"`
}

// AWSConfig represents AWS-specific provider configuration
type AWSConfig struct {
	Family   Provider   `yaml:"family"`
	Services []Provider `yaml:"services"`
}

// Config represents the complete provider configuration
type Config struct {
	AWS            AWSConfig  `yaml:"aws"`
	OtherProviders []Provider `yaml:"otherProviders"`
}

// LoadConfig loads provider configuration from a YAML file
func LoadConfig(configPath string) (*Config, error) {
	// If configPath is empty, try to find the default config file
	if configPath == "" {
		// Try common locations
		locations := []string{
			".crosslab/config/crosslab-config.yaml",
			"config/crosslab-config.yaml",
			"crosslab-config.yaml",
		}

		for _, loc := range locations {
			if _, err := os.Stat(loc); err == nil {
				configPath = loc
				break
			}
		}

		if configPath == "" {
			return nil, fmt.Errorf("no provider configuration file found")
		}
	}

	// Read the configuration file
	data, err := os.ReadFile(configPath)
	if err != nil {
		return nil, fmt.Errorf("error reading config file: %v", err)
	}

	// Parse the configuration
	config := &Config{}
	if err := yaml.Unmarshal(data, config); err != nil {
		return nil, fmt.Errorf("error parsing config file: %v", err)
	}

	return config, nil
}

// GetDefaultConfigPath returns the default configuration file path
func GetDefaultConfigPath() string {
	// Try common locations
	locations := []string{
		".crosslab/config/crosslab-config.yaml",
		"config/crosslab-config.yaml",
		"crosslab-config.yaml",
	}

	for _, loc := range locations {
		absPath, err := filepath.Abs(loc)
		if err != nil {
			continue
		}
		if _, err := os.Stat(absPath); err == nil {
			return absPath
		}
	}

	return ""
}

// Validate validates the provider configuration
func (c *Config) Validate() error {
	if c.AWS.Family.Name == "" || c.AWS.Family.Package == "" || c.AWS.Family.Version == "" {
		return fmt.Errorf("aws family provider configuration is incomplete")
	}

	for i, service := range c.AWS.Services {
		if service.Name == "" || service.Package == "" || service.Version == "" {
			return fmt.Errorf("aws service provider at index %d is incomplete", i)
		}
	}

	for i, provider := range c.OtherProviders {
		if provider.Name == "" || provider.Package == "" || provider.Version == "" {
			return fmt.Errorf("other provider at index %d is incomplete", i)
		}
	}

	return nil
}
