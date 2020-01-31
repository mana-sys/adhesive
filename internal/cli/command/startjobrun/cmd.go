package startjobrun

import (
	"errors"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/cloudformation"
	"github.com/aws/aws-sdk-go/service/glue"
	"github.com/mana-sys/adhesive/internal/cli/command"
	"github.com/mana-sys/adhesive/internal/cli/config"
	"github.com/spf13/cobra"
)

func NewStartJobRunCommand(adhesiveCli *command.AdhesiveCli, opts *config.StartJobRunOptions) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "start-job-run",
		Short: "Remove the current deployment of your CloudFormation template.",
		RunE: func(cmd *cobra.Command, args []string) error {
			return startJobRun(adhesiveCli, opts)
		},
	}

	flags := cmd.Flags()
	flags.StringVar(&opts.StackName, "stack-name", "", "The name of the CloudFormation stack to remove")
	flags.StringVar(&opts.JobName, "job-name", "", "The name of the Glue job to run. This can also be the logical resource ID.")

	return cmd
}

func startJobRun(adhesiveCli *command.AdhesiveCli, opts *config.StartJobRunOptions) error {

	if err := adhesiveCli.InitializeClients(); err != nil {
		return err
	}

	glu := adhesiveCli.Glue()
	cfn := adhesiveCli.CloudFormation()

	name := opts.JobName

	if name == "" {
		return errors.New("option --job-name is required")
	}

	// If stack-name is specified, then interpret job-name as the logical resource ID
	// of the job to be run.
	if opts.StackName != "" {
		out, err := cfn.DescribeStackResource(&cloudformation.DescribeStackResourceInput{
			StackName:         aws.String(opts.StackName),
			LogicalResourceId: aws.String(opts.JobName),
		})
		if err != nil {
			return err
		}

		if *out.StackResourceDetail.ResourceType != "AWS::Glue::Job" {
			return errors.New(name + " is not of type AWS::Glue::Job")
		}

		name = *out.StackResourceDetail.PhysicalResourceId
	}

	// Start the job.
	_, err := glu.StartJobRun(&glue.StartJobRunInput{
		JobName: aws.String(name),
	})

	if err != nil {
		return err
	}

	// TODO: If the --tail-logs option is enabled, then stream the logs to the console.
	return nil
}
