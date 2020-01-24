package config

type DeployOptions struct {
	ConfirmChangeSet   bool
	Guided             bool
	NoExecuteChangeSet bool
	StackName          string
	TemplateFile       string
}

type LocalOptions struct {
}

type PackageOptions struct {
	templateFile       string
	s3Bucket           string
	s3Prefix           string
	kmsKeyID           string
	outputTemplateFile string
	useJSON            bool
	forceUpload        bool
}

type RemoveOptions struct {
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
