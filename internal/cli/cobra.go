package cli

import (
	"github.com/mana-sys/adhesive/internal/cli/command"
	"github.com/mana-sys/adhesive/internal/cli/command/deploy"
	"github.com/mana-sys/adhesive/internal/cli/command/historyserver"
	"github.com/mana-sys/adhesive/internal/cli/command/local"
	package1 "github.com/mana-sys/adhesive/internal/cli/command/package"
	"github.com/mana-sys/adhesive/internal/cli/command/remove"
	"github.com/mana-sys/adhesive/internal/cli/command/startjobrun"
	"github.com/mana-sys/adhesive/internal/cli/version"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

type rootOptions struct {
	configFile string
	debug      bool
}

func NewRootCommand(adhesiveCli *command.AdhesiveCli) *cobra.Command {
	var opts rootOptions

	cmd := &cobra.Command{
		Use: "adhesive",
		PersistentPreRun: func(cmd *cobra.Command, args []string) {
			// Set debug mode.
			if opts.debug {
				logrus.SetLevel(logrus.DebugLevel)
			}
		},
		SilenceErrors:    true,
		SilenceUsage:     true,
		TraverseChildren: true,
		Version:          version.Version,
	}

	flags := cmd.Flags()
	flags.BoolVarP(&opts.debug, "debug", "d", false, "Enable debug mode")
	flags.StringVarP(&opts.configFile, "config", "c", "adhesive.toml",
		"Path to Adhesive configuration file")

	cmd.AddCommand(
		deploy.NewDeployCommand(adhesiveCli),
		local.NewLocalCommand(),
		package1.NewPackageCommand(adhesiveCli),
		remove.NewRemoveCommand(adhesiveCli),
		historyserver.NewHistoryServerCommand(adhesiveCli),
		startjobrun.NewStartJobRunCommand(adhesiveCli),
	)

	return cmd
}
