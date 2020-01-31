package main

//import (
//	"context"
//	"fmt"
//
//	"github.com/mana-sys/adhesive/internal/cli/command"
//
//	"github.com/aws/aws-sdk-go/aws"
//	"github.com/aws/aws-sdk-go/aws/session"
//	"github.com/aws/aws-sdk-go/service/cloudformation"
//	"github.com/aws/aws-sdk-go/service/cloudwatchlogs"
//	"github.com/aws/aws-sdk-go/service/glue"
//	"github.com/mana-sys/adhesive/pkg/watchlog"
//	"github.com/spf13/cobra"
//)
//
//type invokeOptions struct {
//	arguments []string
//	name      string
//	stackName string
//	watch     bool
//}
//
//var (
//	invokeOpts invokeOptions
//
//	invokeCommand = &cobra.Command{
//		Use:  "invoke",
//		RunE: invoke,
//	}
//)
//
//func resolveJobLogicalId(cfn *cloudformation.CloudFormation, stackName string, resource string) (string, error) {
//	out, err := cfn.DescribeStackResource(&cloudformation.DescribeStackResourceInput{
//		LogicalResourceId: aws.String(resource),
//		StackName:         aws.String(stackName),
//	})
//	if err != nil {
//		return "", err
//	}
//
//	return *out.StackResourceDetail.PhysicalResourceId, nil
//}
//
//func streamLogs(ctx context.Context, cwl *cloudwatchlogs.CloudWatchLogs, logGroupName string, logStreamName string) {
//	w := watchlog.NewWatchLog(cwl)
//	_ = w.WatchLogStream(ctx, &cloudwatchlogs.FilterLogEventsInput{
//		LogGroupName:   aws.String(logGroupName),
//		LogStreamNames: []*string{aws.String(logStreamName)},
//	})
//}
//
//func invoke(cmd *cobra.Command, args []string) error {
//	// Load configuration and open Glue client.
//	config, err := command.LoadConfigFile("configs/adhesive.toml")
//	if err != nil {
//		return err
//	}
//
//	sess, err := session.NewSession(&aws.Config{
//		Region: aws.String(config.Region),
//	})
//	if err != nil {
//		return err
//	}
//
//	svc := glue.New(sess)
//
//	// Find the Glue job to run. If a stack name is provided, interpret the
//	// provided name as a CloudFormation logical ID.
//	jobName := invokeOpts.name
//	if invokeOpts.stackName != "" {
//		if jobName, err = resolveJobLogicalId(cloudformation.New(sess), invokeOpts.stackName, jobName); err != nil {
//			return err
//		}
//	}
//
//	// Start the Glue job run.
//	out, err := svc.StartJobRun(&glue.StartJobRunInput{
//		JobName: aws.String(jobName),
//	})
//	if err != nil {
//		return err
//	}
//
//	fmt.Printf("Started Glue job with job run ID: %s\n", *out.JobRunId)
//	return nil
//
//	//
//	//// If the watch option is enabled, stream the CloudWatch logs for the
//	//// job to the user.
//	//out2, err := svc.GetJobRun(&glue.GetJobRunInput{
//	//	JobName: aws.String(jobName),
//	//	RunId:   out.JobRunId,
//	//})
//	//if err != nil {
//	//	return err
//	//}
//	//
//	//ctx, cancel := context.WithCancel(context.Background())
//	//defer cancel()
//	//go streamLogs(ctx, cloudwatchlogs.New(sess), *out2.JobRun.LogGroupName, *out.JobRunId)
//	//
//	//// Periodically poll the Glue job for completion.
//	//for {
//	//	out2, err := svc.GetJobRun(&glue.GetJobRunInput{
//	//		JobName: aws.String(jobName),
//	//		RunId:   out.JobRunId,
//	//	})
//	//
//	//	if err != nil {
//	//		return err
//	//	}
//	//
//	//	//out2.JobRun.CompletedOn
//	//}
//	//
//	//fmt.P
//	//return nil
//}
//
//func init() {
//	flags := invokeCommand.Flags()
//	flags.StringSliceVarP(&invokeOpts.arguments, "arguments", "a", nil, "Arguments for this job run.")
//	flags.StringVarP(&invokeOpts.name, "name", "n", "", "The name of the Glue job to invoke.")
//}
