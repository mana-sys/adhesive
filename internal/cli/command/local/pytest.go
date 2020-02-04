package local

import (
	"github.com/mana-sys/adhesive/internal/cli/command"
	"github.com/spf13/cobra"
)

func NewPytestCommand(adhesiveCli *command.AdhesiveCli, opts *dockerOptions) *cobra.Command {
	cmd := &cobra.Command{
		Use:                "pytest",
		Short:              "Run pytest locally.",
		DisableFlagParsing: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			return buildAndRunDockerCommand(adhesiveCli, "pytest", opts, args)
		},
	}

	return cmd
}
