package command

import (
	"fmt"

	"github.com/BurntSushi/toml"
)

// Config corresponds to the properties defined in an adhesive.yml file.
type Config struct {
	Capability   string
	Region       string
	S3Bucket     string `toml:"s3-bucket"`
	S3Prefix     string `toml:"s3-prefix"`
	StackName    string `toml:"stack-name"`
	TemplatePath string
}

type ErrMissingProperty struct {
	Property string
}

func (err *ErrMissingProperty) Error() string {
	return fmt.Sprintf("missing property: %s", err.Property)
}

type ErrBadPropertyValue struct {
	Property string
	Value    string
}

func (err *ErrBadPropertyValue) Error() string {
	return fmt.Sprintf("invalid value for property %q: %q", err.Property,
		err.Value)
}

func errMissingProperty(property string) error {
	return &ErrMissingProperty{property}
}

func errBadPropertyValue(property string, value string) error {
	return &ErrBadPropertyValue{property, value}
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
