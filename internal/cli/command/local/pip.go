package local

import (
	"github.com/spf13/cobra"
)

func NewPipCommand(opts *dockerOptions) *cobra.Command {
	cmd := &cobra.Command{
		Use:                "pip",
		Short:              "Install Python dependencies for local job runs",
		DisableFlagParsing: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			return buildAndRunDockerCommand("pip", opts, args)
		},
	}

	return cmd
}
