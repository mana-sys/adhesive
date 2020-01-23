package remove

import (
	"github.com/mana-sys/adhesive/internal/cli/command"
	"github.com/spf13/cobra"
)

func NewRemoveCommand(adhesiveCli *command.AdhesiveCli) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "remove",
		Short: "Remove the current deployment of your Glue jobs.",
		RunE: func(cmd *cobra.Command, args []string) error {
			return nil
		},
	}

	return cmd
}
