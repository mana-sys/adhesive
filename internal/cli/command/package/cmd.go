package package1

import (
	"errors"
	"fmt"
	"io/ioutil"

	"github.com/mana-sys/adhesive/internal/cli/config"

	"github.com/mana-sys/adhesive/internal/cli/command"
	"github.com/mana-sys/adhesive/pkg/packager"
	"github.com/spf13/cobra"
)

type packageOptions struct {
	templateFile       string
	s3Bucket           string
	s3Prefix           string
	kmsKeyID           string
	outputTemplateFile string
	useJSON            bool
	forceUpload        bool
}

func NewPackageCommand(adhesiveCli *command.AdhesiveCli, opts *config.PackageOptions) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "package",
		Short: "Packages the Glue jobs in your AWS CloudFormation template",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			return package1(adhesiveCli)
		},
	}

	flags := cmd.Flags()
	flags.StringVar(&opts.TemplateFile, "template-file", "template.yml", "The path where your AWS CloudFormation template is located")
	flags.StringVar(&opts.S3Bucket, "s3-bucket", "", "The S3 bucket where artifacts will be uploaded")
	flags.StringVar(&opts.S3Prefix, "s3-prefix", "", "The prefix added to the names of the artifacts uploaded to the S3 bucket")
	flags.StringVar(&opts.KmsKeyID, "kms-key-id", "", "The ID of the KMS key used to encrypt artifacts in the S3 bucket")
	flags.StringVar(&opts.OutputTemplateFile, "output-template-file", "", "The path to the file to which the packaged template will be written")
	flags.BoolVar(&opts.UseJSON, "use-json", false, "Use JSON for the template output format")
	flags.BoolVar(&opts.ForceUpload, "force-upload", false, "Override existing files in the the S3 bucket")

	return cmd
}

func package1(adhesiveCli *command.AdhesiveCli) error {
	opts := adhesiveCli.Config.Package
	var (
		s3Bucket string
		s3Prefix string
	)

	if err := adhesiveCli.InitializeClients(); err != nil {
		return err
	}

	// Determine packaging parameters. These may come from either the packageOptions
	// or from the CLI instance itself.
	if opts.S3Bucket != "" {
		s3Bucket = opts.S3Bucket
	} else {
		return errors.New("must specify an S3 bucket")
	}

	if opts.S3Prefix != "" {
		s3Prefix = opts.S3Prefix
	}

	// Initialize a packager.
	pack := packager.NewFromS3(adhesiveCli.S3())
	pack.S3Bucket = s3Bucket
	pack.S3Prefix = s3Prefix
	pack.KMSKeyID = opts.KmsKeyID
	if opts.UseJSON {
		pack.Format = packager.FormatJSON
	}

	// Package the CloudFormation template.
	b, err := pack.PackageTemplateFile(opts.TemplateFile)
	if err != nil {
		return err
	}

	// If an output file was specified, save the packaged template to
	// that file. Otherwise, output the result to standard output.
	if opts.OutputTemplateFile != "" {
		return ioutil.WriteFile(opts.OutputTemplateFile, b, 0644)
	}

	fmt.Printf("%s", b)
	return nil
}
