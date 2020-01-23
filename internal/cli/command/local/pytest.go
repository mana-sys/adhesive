package local

import "github.com/spf13/cobra"

func NewPytestCommand(opts *dockerOptions) *cobra.Command {
	cmd := &cobra.Command{
		Use:                "pytest",
		Short:              "Run pytest locally.",
		DisableFlagParsing: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			return buildAndRunDockerCommand("pytest", opts, args)
		},
	}

	return cmd
}
