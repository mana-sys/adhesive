package cli

import (
	"github.com/mana-sys/adhesive/internal/cli/command"
	"github.com/mana-sys/adhesive/internal/cli/command/deploy"
	"github.com/mana-sys/adhesive/internal/cli/command/historyserver"
	"github.com/mana-sys/adhesive/internal/cli/command/local"
	package1 "github.com/mana-sys/adhesive/internal/cli/command/package"
	"github.com/mana-sys/adhesive/internal/cli/command/remove"
	"github.com/mana-sys/adhesive/internal/cli/command/startjobrun"
	"github.com/mana-sys/adhesive/internal/cli/config"
	"github.com/mana-sys/adhesive/internal/cli/version"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

func NewRootCommand(adhesiveCli *command.AdhesiveCli) *cobra.Command {
	var (
		conf       config.Config
		configFile string
		debug      bool
	)

	cmd := &cobra.Command{
		Use: "adhesive",
		PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
			// Load configuration from the specified configuration file, or
			// adhesive.toml if no file was specified.
			if err := adhesiveCli.ReloadConfigFile(configFile); err != nil {
				return err
			}

			// Set debug mode.
			if debug {
				logrus.SetLevel(logrus.DebugLevel)
			}

			// Merge configuration from flags.
			adhesiveCli.Config.MergeConfig(&conf)

			return nil
		},
		SilenceErrors:    true,
		SilenceUsage:     true,
		TraverseChildren: true,
		Version:          version.Version,
	}

	flags := cmd.Flags()
	flags.BoolVarP(&debug, "debug", "d", false, "Enable debug mode")
	flags.StringVar(&conf.Profile, "profile", "", "The profile to use")
	flags.StringVar(&conf.Region, "region", "", "Region to execute in")
	flags.StringVarP(&configFile, "config", "c", "",
		"Path to Adhesive configuration file")

	cmd.AddCommand(
		deploy.NewDeployCommand(adhesiveCli, &conf.Deploy),
		local.NewLocalCommand(adhesiveCli),
		package1.NewPackageCommand(adhesiveCli, &conf.Package),
		remove.NewRemoveCommand(adhesiveCli, &conf.Remove),
		historyserver.NewHistoryServerCommand(adhesiveCli, &conf.HistoryServer),
		startjobrun.NewStartJobRunCommand(adhesiveCli, &conf.StartJobRun),
	)

	cmd.SetVersionTemplate(version.Template)

	return cmd
}
