package local

import (
	"github.com/mana-sys/adhesive/internal/cli/command"
	"github.com/spf13/cobra"
)

func NewSparkSubmitCommand(adhesiveCli *command.AdhesiveCli, opts *dockerOptions) *cobra.Command {
	cmd := &cobra.Command{
		Use:                "spark-submit",
		Short:              "Submit a Glue job locally",
		DisableFlagParsing: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			return buildAndRunDockerCommand(adhesiveCli, "spark-submit", opts, args)
		},
	}

	return cmd
}
