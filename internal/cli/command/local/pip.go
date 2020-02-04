package local

import (
	"github.com/mana-sys/adhesive/internal/cli/command"
	"github.com/spf13/cobra"
)

func NewPipCommand(adhesiveCli *command.AdhesiveCli, opts *dockerOptions) *cobra.Command {
	cmd := &cobra.Command{
		Use:                "pip",
		Short:              "Install Python dependencies for local job runs",
		DisableFlagParsing: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			return buildAndRunDockerCommand(adhesiveCli, "pip", opts, args)
		},
	}

	return cmd
}
