package remove

import (
	"bufio"
	"fmt"
	"io"
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/cloudformation"
	"github.com/fatih/color"
	"github.com/mana-sys/adhesive/internal/cli/command"
	"github.com/mana-sys/adhesive/internal/cli/config"
	"github.com/mana-sys/adhesive/internal/cli/util"
	"github.com/spf13/cobra"
)

func NewRemoveCommand(adhesiveCli *command.AdhesiveCli) *cobra.Command {
	opts := &adhesiveCli.Config.Remove

	cmd := &cobra.Command{
		Use:   "remove",
		Short: "Remove the current deployment of your CloudFormation template.",
		RunE: func(cmd *cobra.Command, args []string) error {
			return remove(adhesiveCli, opts)
		},
	}

	flags := cmd.Flags()
	flags.StringVar(&opts.StackName, "stack-name", opts.StackName, "The name of the CloudFormation stack to remove")

	return cmd
}

func remove(adhesiveCli *command.AdhesiveCli, opts *config.RemoveOptions) error {
	if err := adhesiveCli.InitializeClients(); err != nil {
		return err
	}

	sc := bufio.NewScanner(os.Stdin)

	cfn := adhesiveCli.CloudFormation()

	// Make sure the stack exists.
	fmt.Println("Retrieving stack information...")
	out, err := cfn.DescribeStacks(&cloudformation.DescribeStacksInput{
		StackName: aws.String(opts.StackName),
	})

	if err != nil {
		return err
	}

	stack := out.Stacks[0]

	// Prompt user for confirmation.
	confirm, err := util.ScannerPrompt(sc, `This action will delete all stack resources and CANNOT be undone.
    Do you wish to continue? Only "yes" will be accepted to continue.

    Enter a value: `, nil)
	if err == io.EOF {
		fmt.Println("Abort")
		return nil
	}
	if err != nil {
		return err
	}
	if confirm != "yes" {
		fmt.Println("Remove cancelled.")
		return nil
	}

	// Delete the stack.
	_, err = cfn.DeleteStack(&cloudformation.DeleteStackInput{
		StackName: stack.StackId,
	})
	if err != nil {
		return err
	}

	fmt.Printf("\nSent stack deletion request. You may track the stack status below:\n\n")

	// Wait for the stack to finish deleting and stream CloudFormation events to the user.
	_, err = util.MonitorStack(cfn, *stack.StackId, *stack.StackName, util.OpDelete)
	if err != nil {
		return err
	}

	fmt.Println()
	color.Green("Deletion complete!")

	return nil
}
