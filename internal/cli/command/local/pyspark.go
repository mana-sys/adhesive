package local

import "github.com/spf13/cobra"

func NewPySparkCommand(opts *dockerOptions) *cobra.Command {
	cmd := &cobra.Command{
		Use:                "pyspark",
		Short:              "Run PySpark locally",
		DisableFlagParsing: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			return buildAndRunDockerCommand("pyspark", opts, args)
		},
	}

	return cmd
}
