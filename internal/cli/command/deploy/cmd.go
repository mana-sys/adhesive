package deploy

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/service/cloudformation"
	"github.com/fatih/color"
	"github.com/mana-sys/adhesive/internal/cli/command"
	"github.com/mana-sys/adhesive/internal/cli/config"
	"github.com/mana-sys/adhesive/internal/cli/util"
	"github.com/olekukonko/tablewriter"
	"github.com/spf13/cobra"
)

type deployOpts struct {
	templateFile string
	stackName    string
}

func NewDeployCommand(adhesiveCli *command.AdhesiveCli) *cobra.Command {
	var opts config.DeployOptions

	cmd := &cobra.Command{
		Use:   "deploy",
		Short: "Deploy your AWS Glue jobs with CloudFormation",
		RunE: func(cmd *cobra.Command, args []string) error {
			return deploy(adhesiveCli, &opts)
		},
	}

	flags := cmd.Flags()
	flags.BoolVarP(&opts.Guided, "guided", "g", false, "Allow Adhesive to guide you through the deployment")
	flags.StringVar(&opts.StackName, "stack-name", "", "The name of the CloudFormation stack being deployed to")
	flags.StringVar(&opts.TemplateFile, "template-file", "template.yml",
		"The path to your CloudFormation template")

	return cmd
}

func promptOptions(adhesiveCli *command.AdhesiveCli, opts *config.DeployOptions) (bool, error) {
	var err error
	sc := bufio.NewScanner(os.Stdin)

	if adhesiveCli.FoundConfigFile {
		fmt.Println("Defaults loaded from adhesive.toml")
	} else {
		fmt.Println("Could not find adhesive.toml")
	}

	// Prompt stack name.
	prompt := "stack name: "
	if opts.StackName != "" {
		prompt += "(" + opts.StackName + ") "
	}

	opts.StackName, err = util.ScannerPrompt(sc, prompt, nil)
	if err != nil {
		return false, err
	}

	// Prompt confirm change set before deployment.
	prompt = fmt.Sprintf("confirm changes before deployment: (%s) ", strconv.FormatBool(opts.ConfirmChangeSet))
	confirmChangeSet, err := util.ScannerPrompt(sc, prompt, []string{"true", "false"})
	if err != nil {
		return false, err
	}

	// Parsing here is infallible.
	opts.ConfirmChangeSet, _ = strconv.ParseBool(confirmChangeSet)

	return false, nil
}

func deploy(adhesiveCli *command.AdhesiveCli, opts *config.DeployOptions) error {
	if err := adhesiveCli.InitializeClients(); err != nil {
		return err
	}

	sc := bufio.NewScanner(os.Stdin)

	cfn := adhesiveCli.CloudFormation()

	// If the --guided flag was set, prompt for default options.
	if opts.Guided {
		if _, err := promptOptions(adhesiveCli, opts); err != nil {
			return err
		}
	}

	// Option validation.
	if opts.StackName == "" {
		return errors.New("--stack-name must be specified")
	}

	// Read the template file.
	b, err := ioutil.ReadFile(opts.TemplateFile)
	if err != nil {
		return err
	}

	// Check if the stack exists. If the stack does not exist, or if the stack
	// status is "REVIEW_IN_PROGRESS", then the change set type will be "CREATE."
	// Otherwise, the change set type will be "UPDATE."
	var changeSetType string
	changeSetName := fmt.Sprintf("adhesive-%d", time.Now().Unix())
	fmt.Println("Retrieving stack information...")
	out, err := cfn.DescribeStacks(&cloudformation.DescribeStacksInput{
		StackName: aws.String(opts.StackName),
	})
	if err != nil {
		if aerr, ok := err.(awserr.Error); ok && aerr.Code() == "ValidationError" {
			changeSetType = cloudformation.ChangeSetTypeCreate
		} else {
			return err
		}
	} else if *out.Stacks[0].StackStatus == cloudformation.StackStatusReviewInProgress {
		changeSetType = cloudformation.ChangeSetTypeCreate
	} else {
		changeSetType = cloudformation.ChangeSetTypeUpdate
	}

	// Create the change set.
	fmt.Println("Creating change set...")
	out2, err := cfn.CreateChangeSet(&cloudformation.CreateChangeSetInput{
		StackName:     aws.String(opts.StackName),
		ChangeSetName: aws.String(changeSetName),
		ChangeSetType: aws.String(changeSetType),
		TemplateBody:  aws.String(string(b)),
	})
	if err != nil {
		return err
	}

	// Wait for the change set to finish creating.
	var changeSetOutput *cloudformation.DescribeChangeSetOutput
	for {
		changeSetOutput, err = cfn.DescribeChangeSet(&cloudformation.DescribeChangeSetInput{
			ChangeSetName: out2.Id,
			StackName:     out2.StackId,
		})
		if err != nil {
			return err
		}

		status := *changeSetOutput.Status
		if status == cloudformation.ChangeSetStatusCreateComplete {
			break
		} else if status == cloudformation.ChangeSetStatusFailed {
			if strings.HasPrefix(*changeSetOutput.StatusReason, "The submitted information didn't contain changes.") {
				fmt.Printf("\nNo changes to make.\n")
				return nil
			}

			return fmt.Errorf("failed to create change set: %s", *changeSetOutput.StatusReason)
		}

		time.Sleep(time.Second)
	}

	fmt.Printf(`
Finished creating change set.
The following changes will be made as part of the deployment:

`)

	var (
		numAdd    int
		numRemove int
		numUpdate int
	)

	for _, change := range changeSetOutput.Changes {
		if *change.ResourceChange.Action == "Add" {
			numAdd += 1
		} else if *change.ResourceChange.Action == "Remove" {
			numRemove += 1
		} else if *change.ResourceChange.Action == "Update" {
			numUpdate += 1
		}
	}

	// Output the proposed changes.
	renderChangeSet(os.Stdout, changeSetOutput.Changes)

	// If the --no-execute-change-set flag is present, we are done.
	if opts.NoExecuteChangeSet {
		return nil
	}

	fmt.Printf("\nChange set: %s, %s, %s\n",
		color.GreenString("%d to add", numAdd),
		color.YellowString("%d to update", numUpdate),
		color.RedString("%d to remove", numRemove),
	)

	// If the --confirm-change-set flag is present, prompt for confirmation.
	confirm, err := util.ScannerPrompt(sc, `
Do you want to apply this change set?
    Only "yes" will be accepted to approve.

    Enter a value: `, nil)
	if err == io.EOF {
		fmt.Println("Abort.")
		return nil
	}
	if err != nil {
		return err
	}

	if confirm != "yes" {
		fmt.Printf("\nDeploy cancelled.\n")
		return nil
	}

	fmt.Println()

	_, err = cfn.ExecuteChangeSet(&cloudformation.ExecuteChangeSetInput{
		StackName:     changeSetOutput.StackId,
		ChangeSetName: changeSetOutput.ChangeSetName,
	})

	if err != nil {
		return err
	}

	fmt.Printf("Sent the change set execution request. You may track the stack status below:\n\n")

	// Wait until stack completion is done.
	stackOut, err := util.MonitorStack(cfn, *changeSetOutput.StackId, *changeSetOutput.StackName, util.OpUpdate)
	if err != nil {
		return err
	}

	// Print change set execution summary.
	fmt.Println()
	color.Green("Change set execution complete!")
	fmt.Printf("Summary of changes: %s, %s, %s\n\n",
		color.GreenString("%d added", numAdd),
		color.YellowString("%d updated", numUpdate),
		color.RedString("%d removed", numRemove),
	)

	// Print stack outputs.
	fmt.Println("Outputs: ")
	for _, output := range stackOut.Stacks[0].Outputs {
		fmt.Printf("%s = %s\n", *output.OutputKey, *output.OutputValue)
	}

	return nil
}

func stringElseDash(s *string) string {
	if s == nil {
		return "-"
	}
	return *s
}

func renderChangeSet(w io.Writer, changes []*cloudformation.Change) {
	data := make([][]string, 0, len(changes))

	// Translate from *cloudformation.Change to []string
	for _, change := range changes {
		row := []string{
			stringElseDash(change.ResourceChange.Action),
			stringElseDash(change.ResourceChange.LogicalResourceId),
			stringElseDash(change.ResourceChange.PhysicalResourceId),
			stringElseDash(change.ResourceChange.ResourceType),
			stringElseDash(change.ResourceChange.Replacement),
		}

		data = append(data, row)
	}

	table := tablewriter.NewWriter(w)
	table.SetHeader([]string{"Action", "Logical ID", "Physical ID", "Resource type", "Replacement"})
	table.SetAutoFormatHeaders(false)

	for _, v := range data {
		var color int
		if v[0] == "Add" {
			color = tablewriter.FgGreenColor
		} else if v[0] == "Remove" {
			color = tablewriter.FgRedColor
		} else if v[0] == "Modify" {
			color = tablewriter.FgYellowColor
		}

		table.Rich(v, []tablewriter.Colors{
			{color},
			{tablewriter.FgWhiteColor},
			{tablewriter.FgWhiteColor},
			{tablewriter.FgWhiteColor},
			{tablewriter.FgWhiteColor},
		})
	}
	table.Render()
}
