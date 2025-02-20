package config

import (
	"fmt"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

// Initializer handles configuration initialization
type Initializer struct {
	OutputDir string
}

// NewInitializer creates a new configuration initializer
func NewInitializer(outputDir string) *Initializer {
	return &Initializer{
		OutputDir: outputDir,
	}
}

// Initialize creates the default configuration files
func (i *Initializer) Initialize() error {
	// Create output directory if it doesn't exist
	if err := os.MkdirAll(i.OutputDir, 0755); err != nil {
		return fmt.Errorf("failed to create output directory: %v", err)
	}

	// Create Kind configuration
	if err := i.createKindConfig(); err != nil {
		return err
	}

	// Create providers configuration
	if err := i.createProvidersConfig(); err != nil {
		return err
	}

	return nil
}

// createKindConfig creates the Kind cluster configuration file
func (i *Initializer) createKindConfig() error {
	kindConfig := DefaultKindConfig()
	kindConfigPath := filepath.Join(i.OutputDir, "kind-config.yaml")

	if err := writeYAMLFile(kindConfigPath, kindConfig); err != nil {
		return fmt.Errorf("failed to write Kind configuration: %v", err)
	}

	fmt.Printf("Created Kind configuration at: %s\n", kindConfigPath)
	return nil
}

// createProvidersConfig creates the providers configuration file
func (i *Initializer) createProvidersConfig() error {
	configDir := filepath.Join(i.OutputDir, "config")
	if err := os.MkdirAll(configDir, 0755); err != nil {
		return fmt.Errorf("failed to create config directory: %v", err)
	}

	providersConfig := DefaultProvidersConfig()
	providersConfigPath := filepath.Join(configDir, "crosslab-config.yaml")

	if err := writeYAMLFile(providersConfigPath, providersConfig); err != nil {
		return fmt.Errorf("failed to write providers configuration: %v", err)
	}

	fmt.Printf("Created providers configuration at: %s\n", providersConfigPath)
	return nil
}

// writeYAMLFile writes data to a YAML file
func writeYAMLFile(path string, data interface{}) error {
	file, err := os.Create(path)
	if err != nil {
		return err
	}
	defer file.Close()

	encoder := yaml.NewEncoder(file)
	encoder.SetIndent(2)
	return encoder.Encode(data)
}

// DefaultProvidersConfig returns a default providers configuration
func DefaultProvidersConfig() map[string]interface{} {
	return map[string]interface{}{
		"aws": map[string]interface{}{
			"family": map[string]string{
				"name":    "upbound-provider-aws",
				"package": "xpkg.upbound.io/upbound/provider-family-aws",
				"version": "v1",
			},
			"services": []map[string]string{
				{
					"name":    "provider-aws-iam",
					"package": "xpkg.upbound.io/upbound/provider-aws-iam",
					"version": "v1",
				},
				{
					"name":    "provider-aws-s3",
					"package": "xpkg.upbound.io/upbound/provider-aws-s3",
					"version": "v1",
				},
				{
					"name":    "provider-aws-rds",
					"package": "xpkg.upbound.io/upbound/provider-aws-rds",
					"version": "v1",
				},
			},
		},
		"otherProviders": []map[string]string{
			{
				"name":    "provider-helm",
				"package": "xpkg.upbound.io/upbound/provider-helm",
				"version": "v0.20.4",
			},
			{
				"name":    "provider-kubernetes",
				"package": "xpkg.upbound.io/upbound/provider-kubernetes",
				"version": "v0.16.3",
			},
		},
	}
}

// DefaultConfigPaths contains the paths that should exist after initialization
var DefaultConfigPaths = []string{
	".crosslab/config/crosslab-config.yaml",
	".crosslab/kind-config.yaml",
}

// FileExists checks if a file exists at the given path
func FileExists(path string) bool {
	_, err := os.Stat(path)
	return !os.IsNotExist(err)
}

// CheckInitialized verifies that the necessary configuration files exist
// indicating that initialization has been completed
func CheckInitialized() error {
	for _, path := range DefaultConfigPaths {
		if !FileExists(path) {
			return fmt.Errorf("required configuration file not found at %s. Please run 'crosslab init' first", path)
		}
	}
	return nil
}

// CheckConfigFile verifies that a specific configuration file exists
func CheckConfigFile(path string) error {
	if !FileExists(path) {
		dir := filepath.Dir(path)
		return fmt.Errorf("configuration file not found at %s. Please run 'crosslab init' to create the %s directory and configuration files", path, dir)
	}
	return nil
}
func (i *Initializer) GetKindConfig() string {
	return filepath.Join(i.OutputDir, "kind-config.yaml")
}
func (i *Initializer) GetProvidersConfig() string {
	return filepath.Join(i.OutputDir, "config", "crosslab-config.yaml")
}
