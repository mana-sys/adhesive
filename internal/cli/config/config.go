package config

import "github.com/BurntSushi/toml"

type DeployOptions struct {
	ConfirmChangeSet   bool `toml:"confirm-change-set"`
	Guided             bool
	NoExecuteChangeSet bool   `toml:"no-execute-change-set"`
	StackName          string `toml:"stack-name"`
	TemplateFile       string `toml:"template-file"`
}

type LocalOptions struct {
}

type PackageOptions struct {
	TemplateFile       string `toml:"template-file"`
	s3Bucket           string
	s3Prefix           string
	kmsKeyID           string
	outputTemplateFile string
	useJSON            bool
	forceUpload        bool
}

type RemoveOptions struct {
	StackName string `toml:"stack-name"`
}

type Config struct {
	Deploy  DeployOptions
	Local   LocalOptions
	Package PackageOptions
	Remove  RemoveOptions

	// Root command options.
	Profile string
	Region  string
}

func defaultConfig() *Config {
	return &Config{
		Region: "us-west-2",
	}
}

func NewConfig() *Config {
	return defaultConfig()
}

func LoadConfigFileInto(config *Config, path string) error {
	_, err := toml.DecodeFile(path, &config)
	return err
}

// LoadConfigFile reads configuration from the adhesive.toml file.
func LoadConfigFile(path string) (*Config, error) {
	config := defaultConfig()
	if _, err := toml.DecodeFile(path, &config); err != nil {
		return nil, err
	}

	return config, nil
}
