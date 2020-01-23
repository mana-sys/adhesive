package deploy

import (
	"fmt"

	"github.com/spf13/cobra"
)

type deployOpts struct {
	templateFile string
	stackName    string
}

func NewDeployCommand() *cobra.Command {
	var opts deployOpts

	cmd := &cobra.Command{
		Use:   "deploy",
		Short: "Deploy your AWS Glue jobs with CloudFormation",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println("deploy")
		},
	}

	flags := cmd.Flags()
	flags.StringVar(&opts.stackName, "stack-name", "", "The name of the CloudFormation stack being deployed to")
	flags.StringVar(&opts.templateFile, "template-file", "template.yml",
		"The path to your CloudFormation template")

	cmd.MarkFlagRequired("stack-name")
	return cmd
}
