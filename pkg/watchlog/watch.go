package watchlog

import (
	"context"
	"fmt"
	"io"
	"os"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go/aws/request"

	"github.com/aws/aws-sdk-go/aws/awserr"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/cloudwatchlogs"
	"github.com/aws/aws-sdk-go/service/cloudwatchlogs/cloudwatchlogsiface"
	"github.com/fatih/color"
)

// WatchLog represents a CloudWatch Logs watcher. It will stream logs to the
// configured destination until canceled. Adapted from Tyler Brock's
// github.com/TylerBrock/saw repository.
type WatchLog struct {
	cwl cloudwatchlogsiface.CloudWatchLogsAPI

	Destination  io.Writer
	Formatter    Formatter
	PollInterval time.Duration
}

// A mechanism to format a CloudWatch Logs filtered log event.
type Formatter interface {
	FormatFilteredLogEvent(event *cloudwatchlogs.FilteredLogEvent) string
}

// An adapter type to allow the use of ordinary functions as Formatters.
type FormatterFunc func(event *cloudwatchlogs.FilteredLogEvent) string

func (f FormatterFunc) FormatFilteredLogEvent(event *cloudwatchlogs.FilteredLogEvent) string {
	return f(event)
}

var (
	// Default formatter. Formats the event with the timestamp, log stream
	// name, and message.
	defaultFormatter = FormatterFunc(func(event *cloudwatchlogs.FilteredLogEvent) string {
		white := color.New(color.FgWhite).SprintFunc()

		date := aws.MillisecondsTimeValue(event.Timestamp)
		dateStr := date.Format(time.RFC3339)
		return fmt.Sprintf("[%s] - %s", white(dateStr), *event.Message)
	})

	// Returns the message only.
	rawFormatter = FormatterFunc(func(event *cloudwatchlogs.FilteredLogEvent) string {
		return *event.Message
	})
)

// NewWatchLog returns a WatchLog instance with the default configuration. By
// default, the output destination is os.Stdout, and the poll interval is 1
// second.
func NewWatchLog(cwl cloudwatchlogsiface.CloudWatchLogsAPI) *WatchLog {
	return &WatchLog{
		cwl:          cwl,
		Destination:  os.Stdout,
		Formatter:    defaultFormatter,
		PollInterval: time.Second,
	}
}

// A configuration option for initializing a WatchLog instance.
type Opt func(w *WatchLog)

func NewWatchLogWithOpts(opts ...Opt) *WatchLog {
	w := new(WatchLog)
	for _, opt := range opts {
		opt(w)
	}

	return w
}

func WithCloudWatchLogs(cwl *cloudwatchlogs.CloudWatchLogs) Opt {
	return func(w *WatchLog) { w.cwl = cwl }
}

func WithDestination(dest io.Writer) Opt {
	return func(w *WatchLog) { w.Destination = dest }
}

func WithFormatter(formatter Formatter) Opt {
	return func(w *WatchLog) { w.Formatter = formatter }
}

func WithPollInterval(interval time.Duration) Opt {
	return func(w *WatchLog) { w.PollInterval = interval }
}

// WatchLogStream watches the specified log stream using the configured polling
// interval. Adapted from the Blade Stream() method.
func (w *WatchLog) WatchLogStream(ctx context.Context, input *cloudwatchlogs.FilterLogEventsInput) error {
	var (
		// Need two errors so that we can keep track of errors in the handler.
		err, wrErr   error
		lastSeenTime *int64
		seenEventIDs map[string]bool
	)

	// Page handler function. Called for every page returned by
	// FilterLogEventsPages().
	handlePage := func(page *cloudwatchlogs.FilterLogEventsOutput, lastPage bool) bool {
		for _, event := range page.Events {

			// Update lastSeenTime, which will be used as the start time for
			// the next call to FilterLogEventsPages().
			if lastSeenTime == nil || *event.Timestamp > *lastSeenTime {
				lastSeenTime = event.Timestamp
				seenEventIDs = make(map[string]bool, 0)
			}

			// If we have not yet seen this event, write it to the destination.
			if _, seen := seenEventIDs[*event.EventId]; !seen {
				// Format the event with the formatter.
				message := w.Formatter.FormatFilteredLogEvent(event)
				message = strings.TrimRight(message, "\n")

				// Write the message to the destination.
				_, wrErr = fmt.Fprintln(w.Destination, message)
				if wrErr != nil {
					return false
				}

				// Don't process this event if we see it again.
				seenEventIDs[*event.EventId] = true
			}
		}
		return !lastPage
	}

	for {
		err = w.cwl.FilterLogEventsPagesWithContext(ctx, input, handlePage)
		if err != nil {
			// If the request was canceled, return the cancellation error.
			if aerr, ok := err.(awserr.Error); ok && aerr.Code() == request.CanceledErrorCode {
				_, _ = fmt.Fprintf(w.Destination, "request was canceled: %v", err)
				return ctx.Err()
			}

			// Otherwise, break the loop and return the error.
			_, _ = fmt.Fprintf(w.Destination, "failed to retrieve logs: %v", err)
			break
		}

		// Check possible errors from the log events handler.
		if wrErr != nil {
			return wrErr
		}

		// Set StartTime filter for the next query.
		if lastSeenTime != nil {
			input.SetStartTime(*lastSeenTime)
		}

		time.Sleep(w.PollInterval)
	}

	return err
}
