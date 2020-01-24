package deploy

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/mana-sys/adhesive/pkg/watchstack"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/cloudformation"
)

type stackOp int

type StackOpError string

func (err StackOpError) Error() string {
	return string(err)
}

const (
	OpCreate stackOp = iota
	OpUpdate
	OpDelete
)

type stackMonitorState struct {
	name       string
	op         stackOp
	start      time.Time
	uiInterval time.Duration
}

func isSuccess(status string) bool {
	return status == cloudformation.StackStatusCreateComplete ||
		status == cloudformation.StackStatusDeleteComplete ||
		status == cloudformation.StackStatusUpdateComplete
}

func isFailed(status string) bool {
	return status == cloudformation.StackStatusDeleteFailed ||
		strings.HasSuffix(status, "ROLLBACK_COMPLETE") ||
		strings.HasSuffix(status, "ROLLBACK_FAILED")
}

func monitorStack(cfn *cloudformation.CloudFormation, stackId, stackName string, op stackOp) (*cloudformation.DescribeStacksOutput, error) {
	var (
		err error
		out *cloudformation.DescribeStacksOutput
	)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	state := stackMonitorState{
		name:       stackName,
		op:         OpCreate,
		start:      time.Now().Round(time.Second),
		uiInterval: 5 * time.Second,
	}

	w := watchstack.New(cfn)

	go w.StreamEvents(ctx, stackId)
	for range time.Tick(1 * time.Second) {
		out, err = cfn.DescribeStacks(&cloudformation.DescribeStacksInput{
			StackName: aws.String(stackId),
		})
		if err != nil {
			return nil, err
		}

		stack := out.Stacks[0]
		if isSuccess(*stack.StackStatus) {
			break
		} else if isFailed(*stack.StackStatus) {
			return nil, StackOpError(fmt.Sprintf("operation failed with status %s: %s",
				*stack.StackStatus, *stack.StackStatusReason))
		}
	}

	fmt.Printf("\n%s: operation complete after %v\n", stackName,
		time.Now().Round(1*time.Second).Sub(state.start))

	return out, nil
}

func consoleMonitorStack(ctx context.Context, state stackMonitorState) {
	ticker := time.NewTicker(state.uiInterval)

	var message string
	switch state.op {
	case OpCreate:
		message = "Still creating..."
	case OpDelete:
		message = "Still deleting..."
	case OpUpdate:
		message = "Still updating..."
	default:
		panic("unreachable")
	}

	for {
		select {
		case <-ctx.Done():
			return
		case t := <-ticker.C:
			fmt.Printf("%s: %s (%v elapsed)\n", state.name, message, t.Round(time.Second).Sub(state.start))
		}
	}
}
