package local

import "github.com/spf13/cobra"

type dockerOptions struct {
	arg         []string
	credentials string
	env         []string
	volumes     []string
}

type localOptions struct {
	dockerOptions
}

func NewLocalCommand() *cobra.Command {
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
		NewPipCommand(dockerOpts),
		NewPySparkCommand(dockerOpts),
		NewPytestCommand(dockerOpts),
		NewSparkSubmitCommand(dockerOpts),
	)

	return cmd
}
