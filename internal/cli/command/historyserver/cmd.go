package historyserver

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/mana-sys/adhesive/internal/cli/command"
	"github.com/mana-sys/adhesive/internal/cli/config"
	"github.com/spf13/cobra"
)

func NewHistoryServerCommand(adhesiveCli *command.AdhesiveCli) *cobra.Command {
	opts := &adhesiveCli.Config.HistoryServer

	cmd := &cobra.Command{
		Use:   "history-server",
		Short: "Launch the Spark history server",
		RunE: func(cmd *cobra.Command, args []string) error {
			return historyServer(adhesiveCli, opts)
		},
	}

	flags := cmd.Flags()
	flags.IntVarP(&opts.Port, "port", "p", opts.Port, "The port to listen on")
	flags.StringVar(&opts.LogDirectory, "log-directory", opts.LogDirectory, "The location of the Spark logs. Must be an s3a:// formatted path.")

	return cmd
}

// buildDockerCommand builds an exec.Cmd to run the history server Docker container with the provided options.
func buildDockerCommand(opts *config.HistoryServerOptions) (*exec.Cmd, error) {
	credsDir, err := os.UserHomeDir()
	if err != nil {
		return nil, err
	}

	credsDir = filepath.Join(credsDir, ".aws")

	dockerArgs := []string{"run", "--rm"}
	dockerArgs = append(dockerArgs, "-v", credsDir+":/root/.aws")
	dockerArgs = append(dockerArgs, "-p", strconv.FormatInt(int64(opts.Port), 10)+":18080")
	sparkHistoryOptsStringFormat := "SPARK_HISTORY_OPTS=-Dspark.hadoop.fs.s3a.aws.credentials.provider=com.amazonaws.auth.DefaultAWSCredentialsProviderChain " +
		"-Dspark.history.fs.logDirectory=%s"

	dockerArgs = append(dockerArgs, "-e", fmt.Sprintf(sparkHistoryOptsStringFormat, opts.LogDirectory))
	dockerArgs = append(dockerArgs, "sysmana/sparkui:latest",
		"/opt/spark/bin/spark-class org.apache.spark.deploy.history.HistoryServer")

	return exec.Command("docker", dockerArgs...), nil
}

func historyServer(adhesiveCli *command.AdhesiveCli, opts *config.HistoryServerOptions) error {
	if opts.LogDirectory == "" {
		return errors.New("option --log-directory is required")
	}

	if opts.Port == 0 {
		opts.Port = 18080
	}

	if !strings.HasPrefix(opts.LogDirectory, "s3a://") {
		return errors.New("option --log-directory must be an s3a:// formatted path")
	}

	cmd, err := buildDockerCommand(opts)
	if err != nil {
		return err
	}

	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err = cmd.Start(); err != nil {
		return err
	}

	return cmd.Wait()
}
