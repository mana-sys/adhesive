package local

import (
	"github.com/mana-sys/adhesive/internal/cli/command"
	"github.com/spf13/cobra"
)

func NewPySparkCommand(adhesiveCli *command.AdhesiveCli, opts *dockerOptions) *cobra.Command {
	cmd := &cobra.Command{
		Use:                "pyspark",
		Short:              "Run PySpark locally",
		DisableFlagParsing: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			return buildAndRunDockerCommand(adhesiveCli, "pyspark", opts, args)
		},
	}

	return cmd
}
