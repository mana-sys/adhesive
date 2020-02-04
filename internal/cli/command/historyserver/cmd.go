package historyserver

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/aws/aws-sdk-go/aws/credentials/processcreds"
	"github.com/mana-sys/adhesive/internal/cli/command"
	"github.com/mana-sys/adhesive/internal/cli/config"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

func NewHistoryServerCommand(adhesiveCli *command.AdhesiveCli, opts *config.HistoryServerOptions) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "history-server",
		Short: "Launch the Spark history server",
		RunE: func(cmd *cobra.Command, args []string) error {
			return historyServer(adhesiveCli, opts)
		},
	}

	flags := cmd.Flags()
	flags.IntVarP(&opts.Port, "port", "p", 0, "The port to listen on")
	flags.StringVar(&opts.LogDirectory, "log-directory", "", "The location of the Spark logs. Must be an s3a:// formatted path.")

	return cmd
}

// buildDockerCommand builds an exec.Cmd to run the history server Docker container with the provided options.
func buildDockerCommand(adhesiveCli *command.AdhesiveCli, opts *config.HistoryServerOptions) (*exec.Cmd, error) {
	credsDir, err := os.UserHomeDir()
	if err != nil {
		return nil, err
	}

	credsDir = filepath.Join(credsDir, ".aws")

	dockerArgs := []string{"run", "--rm"}

	dockerArgs = append(dockerArgs, "-v", credsDir+":/root/.aws")
	dockerArgs = append(dockerArgs, "-p", strconv.FormatInt(int64(opts.Port), 10)+":18080")

	// Super hack: If we used the ProcessProvider, then we pass the credentials via environment variables to the
	// Docker container.
	value, err := adhesiveCli.Session().Config.Credentials.Get()
	if err != nil {
		return nil, err
	}

	if value.ProviderName == processcreds.ProviderName {
		dockerArgs = append(dockerArgs,
			"-e", "AWS_ACCESS_KEY_ID="+value.AccessKeyID,
			"-e", "AWS_SECRET_ACCESS_KEY="+value.SecretAccessKey,
			"-e", "AWS_SESSION_TOKEN="+value.SessionToken,
		)
	}

	// Environment variable for Spark history server options.
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

	// Super hack: initialize the clients to retrieve the credentials. This is needed I couldn't
	// figure out how to get credential_process to work for Java.
	if err := adhesiveCli.InitializeClients(); err != nil {
		return err
	}

	cmd, err := buildDockerCommand(adhesiveCli, opts)
	if err != nil {
		return err
	}

	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	log.Debug("Running Docker command: ", cmd.Args)

	if err = cmd.Start(); err != nil {
		return err
	}

	return cmd.Wait()
}
