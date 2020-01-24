package watchstack

import (
	"context"
	"fmt"
	"io"
	"os"
	"strings"
	"time"

	"github.com/fatih/color"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/cloudformation"
	"github.com/aws/aws-sdk-go/service/cloudformation/cloudformationiface"
)

// WatchLog represents a CloudFormation events watcher. It will stream events
// to the configured destination until canceled. Adapted from Tyler Brock's
// github.com/TylerBrock/saw repository.
type WatchStack struct {
	cfn          cloudformationiface.CloudFormationAPI
	PollInterval time.Duration
	W            io.Writer
}

func New(cfn cloudformationiface.CloudFormationAPI) *WatchStack {
	return &WatchStack{
		cfn:          cfn,
		PollInterval: time.Second,
		W:            os.Stdout,
	}
}

func isFailedResource(status string) bool {
	return strings.HasSuffix(status, "FAILED")
}

func isSuccessResource(status string) bool {
	return strings.HasSuffix(status, "COMPLETE")
}

func (w *WatchStack) StreamEvents(ctx context.Context, stackId string) error {
	var (
		lastSeenTime time.Time
		seenEvents   map[string]struct{}
	)

	input := &cloudformation.DescribeStackEventsInput{
		StackName: aws.String(stackId),
	}
	pageHandler := func(output *cloudformation.DescribeStackEventsOutput, last bool) bool {
		for _, event := range output.StackEvents {
			// Filter out any events that occurred before that last seen timestamp.
			// If we see a new timestamp, then we set the last seen timestamp to that
			// new timestamp. We can also clear the set of seen events.
			if event.Timestamp.Before(lastSeenTime) {
				continue
			} else if event.Timestamp.After(lastSeenTime) {
				lastSeenTime = *event.Timestamp
				seenEvents = make(map[string]struct{})
			}

			// If we have already seen this event, don't output it. Otherwise, we
			// can output the event.
			if _, ok := seenEvents[*event.EventId]; ok {
				continue
			}

			c := color.YellowString
			if isSuccessResource(*event.ResourceStatus) {
				c = color.GreenString
			} else if isFailedResource(*event.ResourceStatus) {
				c = color.RedString
			}

			fmt.Fprintf(w.W, "[%s] %s (%s): %s\n",
				color.WhiteString(event.Timestamp.Format(time.RFC3339)),
				*event.LogicalResourceId,
				*event.ResourceType,
				c(*event.ResourceStatus),
			)
			seenEvents[*event.EventId] = struct{}{}
		}
		return !last
	}

	for {
		err := w.cfn.DescribeStackEventsPages(input, pageHandler)
		if err != nil {
			return err
		}

		time.Sleep(w.PollInterval)
	}
}
