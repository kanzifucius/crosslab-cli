package provider

// Provider represents a Crossplane provider configuration
type Provider struct {
	Name    string `yaml:"name"`
	Package string `yaml:"package"`
	Version string `yaml:"version"`
}

// Config represents the configuration for all providers
type Config struct {
	AWS            AWSConfig  `yaml:"aws"`
	OtherProviders []Provider `yaml:"otherProviders"`
}

// AWSConfig represents AWS-specific provider configuration
type AWSConfig struct {
	Family   Provider   `yaml:"family"`
	Services []Provider `yaml:"services"`
}
