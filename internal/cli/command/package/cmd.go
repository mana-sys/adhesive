package package1

import (
	"errors"
	"fmt"
	"io/ioutil"

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

func NewPackageCommand(adhesive *command.AdhesiveCli) *cobra.Command {
	var opts packageOptions

	cmd := &cobra.Command{
		Use:   "package",
		Short: "Packages the Glue jobs in your AWS CloudFormation template",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			return package1(adhesive, &opts)
		},
	}

	flags := cmd.Flags()
	flags.StringVar(&opts.templateFile, "template-file", "template.yml", "The path where your AWS CloudFormation template is located")
	flags.StringVar(&opts.s3Bucket, "s3-bucket", "", "The S3 bucket where artifacts will be uploaded")
	flags.StringVar(&opts.s3Prefix, "s3-prefix", "", "The prefix added to the names of the artifacts uploaded to the S3 bucket")
	flags.StringVar(&opts.kmsKeyID, "kms-key-id", "", "The ID of the KMS key used to encrypt artifacts in the S3 bucket")
	flags.StringVar(&opts.outputTemplateFile, "output-template-file", "", "The path to the file to which the packaged template will be written")
	flags.BoolVar(&opts.useJSON, "use-json", false, "Use JSON for the template output format")
	flags.BoolVar(&opts.forceUpload, "force-upload", false, "Override existing files in the the S3 bucket")

	return cmd
}

func package1(adhesiveCli *command.AdhesiveCli, opts *packageOptions) error {
	var (
		s3Bucket string
		s3Prefix string
	)

	if err := adhesiveCli.InitializeClients(); err != nil {
		return err
	}

	// Determine packaging parameters. These may come from either the packageOptions
	// or from the CLI instance itself.
	if opts.s3Bucket != "" {
		s3Bucket = opts.s3Bucket
	} else {
		return errors.New("must specify an S3 bucket")
	}

	if opts.s3Prefix != "" {
		s3Prefix = opts.s3Prefix
	}

	// Initialize a packager.
	pack := packager.NewFromS3(adhesiveCli.S3())
	pack.S3Bucket = s3Bucket
	pack.S3Prefix = s3Prefix
	pack.KMSKeyID = opts.kmsKeyID
	if opts.useJSON {
		pack.Format = packager.FormatJSON
	}

	// Package the CloudFormation template.
	b, err := pack.PackageTemplateFile(opts.templateFile)
	if err != nil {
		return err
	}

	// If an output file was specified, save the packaged template to
	// that file. Otherwise, output the result to standard output.
	if opts.outputTemplateFile != "" {
		return ioutil.WriteFile(opts.outputTemplateFile, b, 0644)
	}

	fmt.Printf("%s", b)
	return nil
}
