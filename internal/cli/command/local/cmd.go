package local

import (
	"github.com/mana-sys/adhesive/internal/cli/command"
	"github.com/spf13/cobra"
)

type dockerOptions struct {
	arg         []string
	credentials string
	env         []string
	volumes     []string
}

type localOptions struct {
	dockerOptions
}

func NewLocalCommand(adhesiveCli *command.AdhesiveCli) *cobra.Command {
	var opts localOptions
	cmd := &cobra.Command{
		Use:   "local",
		Short: "Run AWS Glue jobs and test suites locally",
	}

	flags := cmd.Flags()
	flags.StringSliceVarP(&opts.arg, "docker-env", "e", nil, "Set Docker environment variables.")
	flags.StringSliceVarP(&opts.volumes, "docker-volumes", "v", nil, "Mount Docker volumes.")
	flags.StringSliceVarP(&opts.arg, "docker-arg", "a", nil,
		"Pass additional arguments to the \"docker run\" command.")

	dockerOpts := &opts.dockerOptions
	cmd.AddCommand(
		NewPipCommand(adhesiveCli, dockerOpts),
		NewPySparkCommand(adhesiveCli, dockerOpts),
		NewPytestCommand(adhesiveCli, dockerOpts),
		NewSparkSubmitCommand(adhesiveCli, dockerOpts),
	)

	return cmd
}
