package historyserver

import (
	"github.com/mana-sys/adhesive/internal/cli/command"
	"github.com/spf13/cobra"
)

func NewHistoryServerCommand(adhesiveCli *command.AdhesiveCli) *cobra.Command {
	opts := &adhesiveCli.Config.HistoryServer

	cmd := &cobra.Command{
		Use:   "history-server",
		Short: "Launch the Spark history server",
		RunE: func(cmd *cobra.Command, args []string) error {
			return historyServer(adhesiveCli)
		},
	}

	flags := cmd.Flags()
	flags.IntVarP(&opts.Port, "port", "p", opts.Port, "The port to listen on")

	return cmd
}

//// buildDockerCommand builds an exec.Cmd to run the Docker container with the provided options.
//func buildDockerCommand(entrypoint string, args []string) (*exec.Cmd, error) {
//	var (
//		err  error
//		envs []string
//		vols []string
//	)
//
//	for _, env := range options.env {
//		envs = append(envs, "-e "+env)
//	}
//
//	for _, vol := range options.volumes {
//		vols = append(vols, "-v "+vol)
//	}
//
//	credsDir := options.credentials
//	if credsDir == "" {
//		if credsDir, err = os.UserHomeDir(); err != nil {
//			return nil, errors.New("unable to determine home directory")
//		}
//		credsDir = filepath.Join(credsDir, ".aws")
//	}
//
//	wd, err := os.Getwd()
//
//	dockerArgs := []string{"run", "--rm", "-t"}
//	dockerArgs = append(dockerArgs, envs...)
//	dockerArgs = append(dockerArgs, vols...)
//	dockerArgs = append(dockerArgs, "-v", credsDir+":/root/.aws",
//		"-v", DistPackagesVolume+":"+DistPackagesDirectory,
//		"-v", wd+":/project",
//		"--entrypoint", entrypoint, DockerImageName+":0.9")
//	dockerArgs = append(dockerArgs, args...)
//
//	return exec.Command("docker", dockerArgs...), nil
//}

func historyServer(adhesiveCli *command.AdhesiveCli) error {
	return nil
}
