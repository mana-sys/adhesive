package config

import "github.com/BurntSushi/toml"

type DeployOptions struct {
	ConfirmChangeSet   bool   `toml:"confirm-change-set"`
	Guided             bool   `toml:"-"`
	NoExecuteChangeSet bool   `toml:"no-execute-change-set"`
	StackName          string `toml:"stack-name"`
	TemplateFile       string `toml:"template-file"`
}

type LocalOptions struct {
}

type PackageOptions struct {
	TemplateFile       string `toml:"template-file"`
	S3Bucket           string `toml:"s3-bucket"`
	S3Prefix           string `toml:"s3-prefix"`
	KmsKeyID           string `toml:"kms-key-id"`
	OutputTemplateFile string `toml:"output-template-file"`
	UseJSON            bool   `toml:"use-json"`
	ForceUpload        bool   `toml:"force-upload"`
}

type HistoryServerOptions struct {
	Port         int
	LogDirectory string `toml:"log-directory"`
}

type RemoveOptions struct {
	StackName string `toml:"stack-name"`
}

type StartJobRunOptions struct {
	JobName   string `toml:"job-name"`
	JobRunID  string `toml:"job-run-id"`
	StackName string `toml:"stack-name"`
	TailLogs  bool   `toml:"tail-logs"`
}

type Config struct {
	Deploy        DeployOptions
	StartJobRun   StartJobRunOptions `toml:"start-job-run"`
	Local         LocalOptions
	Package       PackageOptions
	Remove        RemoveOptions
	HistoryServer HistoryServerOptions `toml:"history-server"`

	// Root command options.
	ConfigFile string `toml:"-"`
	Profile    string
	Region     string
	Debug      bool `toml:"-"`
}

func defaultConfig() *Config {
	return &Config{
		Region: "us-west-2",
	}
}

func NewConfig() *Config {
	return defaultConfig()
}

func (conf *Config) MergeConfig(other *Config) {
	// Merge root command options.
	if other.ConfigFile != "" {
		conf.ConfigFile = other.ConfigFile
	}

	if other.Profile != "" {
		conf.Profile = other.Profile
	}

	if other.Region != "" {
		conf.Region = other.Region
	}

	if other.Debug {
		conf.Debug = other.Debug
	}

	// Merge "local" options.
	conf.Local.mergeLocalOptions(&other.Local)

	// Merge "package" options.
	conf.Package.mergePackageOptions(&other.Package)

	// Merge "remove" options.
	conf.Remove.mergeRemoveOptions(&other.Remove)

	// Merge "history-server" options.
	conf.HistoryServer.mergeHistoryServerOptions(&other.HistoryServer)

	// Merge "start-job-run" options.
	conf.StartJobRun.mergeStartJobRunOptions(&other.StartJobRun)

	// Merge "deploy" options.
	conf.Deploy.mergeDeployOptions(&other.Deploy)
}

func (opts *LocalOptions) mergeLocalOptions(other *LocalOptions) {
}

func (opts *PackageOptions) mergePackageOptions(other *PackageOptions) {
	if other.TemplateFile != "" {
		opts.TemplateFile = other.TemplateFile
	}

	if other.KmsKeyID != "" {
		opts.KmsKeyID = other.KmsKeyID
	}

	if other.OutputTemplateFile != "" {
		opts.OutputTemplateFile = other.OutputTemplateFile
	}

	if other.S3Bucket != "" {
		opts.S3Bucket = other.S3Bucket
	}

	if other.S3Prefix != "" {
		opts.S3Prefix = other.S3Prefix
	}

	if other.UseJSON {
		opts.UseJSON = other.UseJSON
	}

	if other.ForceUpload {
		opts.ForceUpload = other.ForceUpload
	}
}

func (opts *HistoryServerOptions) mergeHistoryServerOptions(other *HistoryServerOptions) {
	if other.Port != 0 {
		opts.Port = other.Port
	}

	if other.LogDirectory != "" {
		opts.LogDirectory = other.LogDirectory
	}
}

func (opts *RemoveOptions) mergeRemoveOptions(other *RemoveOptions) {
	if other.StackName != "" {
		opts.StackName = other.StackName
	}
}

func (opts *StartJobRunOptions) mergeStartJobRunOptions(other *StartJobRunOptions) {
	if other.JobName != "" {
		opts.JobName = other.JobName
	}

	if other.StackName != "" {
		opts.StackName = other.StackName
	}

	if other.JobRunID != "" {
		opts.JobRunID = other.JobRunID
	}

	if other.TailLogs {
		opts.TailLogs = other.TailLogs
	}
}

func (opts *DeployOptions) mergeDeployOptions(other *DeployOptions) {
	if other.StackName != "" {
		opts.StackName = other.StackName
	}

	if other.TemplateFile != "" {
		opts.TemplateFile = other.TemplateFile
	}

	if other.ConfirmChangeSet {
		opts.ConfirmChangeSet = other.ConfirmChangeSet
	}

	if other.Guided {
		opts.Guided = other.Guided
	}

	if other.NoExecuteChangeSet {
		opts.NoExecuteChangeSet = other.NoExecuteChangeSet
	}
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
