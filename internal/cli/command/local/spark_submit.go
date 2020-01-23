package local

import "github.com/spf13/cobra"

func NewSparkSubmitCommand(opts *dockerOptions) *cobra.Command {
	cmd := &cobra.Command{
		Use:                "spark-submit",
		Short:              "Submit a Glue job locally",
		DisableFlagParsing: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			return buildAndRunDockerCommand("spark-submit", opts, args)
		},
	}

	return cmd
}
