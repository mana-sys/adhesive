package main

//
//import (
//	"fmt"
//
//	"github.com/spf13/cobra"
//)
//
//type docker struct {
//	arg         []string
//	credentials string
//	env         []string
//	volumes     []string
//}
//
//type localOptions struct {
//	docker
//}
//
//const (
//	DockerImageName = "sysmana/aws-glue-dev-base"
//
//	DistPackagesVolume    = "aws_glue_dist_packages"
//	DistPackagesDirectory = "/usr/local/lib/python2.7/dist-packages"
//)
//
////var sparkSubmitCommand = cli.Command{
////	Name:  "spark-submit",
////	Usage: "start a Glue job locally",
////	Action: func(ctx *cli.Context) error {
////		args := []string{"run"}
////		// Build the Docker command, starting with the miscellaneous arguments.
////		args = append(ctx.GlobalStringSlice("docker-arg"))
////
////		// Always run in interactive mode and remove after container shutdown.
////		args = append(args, "-it", "--rm")
////
////		// Add credentials volume.
////		creds := ctx.GlobalString("docker-credentials")
////		if creds != "" {
////			args = append(args, "-v", creds+":/root/.aws")
////		}
////
////		// Add environment variables.
////		for _, env := range ctx.GlobalStringSlice("docker-env-var") {
////			args = append(args, "-e", env)
////		}
////
////		// Add additional volumes.
////		for _, vol := range ctx.GlobalStringSlice("docker-volumes") {
////			args = append(args, "-v", vol)
////		}
////
////		_ = exec.Command("docker", args...)
////
////		return nil
////		//cmd := exec.Command
////	},
////}
////
//
//var (
//	localConfig ocalConfig
//
//	pipCommand = &cobra.Command{
//		Use:                "pip",
//		Short:              "-- install Python dependencies for local job runs",
//		DisableFlagParsing: true,
//		Run: func(cmd *cobra.Command, args []string) {
//			fmt.Println("pip")
//		},
//	}
//
//	sparkSubmitCommand = &cobra.Command{
//		Use:                "spark-submit",
//		Short:              "-- start a Glue job locally",
//		DisableFlagParsing: true,
//		RunE: func(cmd *cobra.Command, args []string) error {
//			return buildAndRunDockerCommand("spark-submit", &localConfig.dockerOptions, args)
//		},
//	}
//
//	localCommand = &cobra.Command{
//		Use:              "local",
//		TraverseChildren: true,
//	}
//)
//
//func NewLocalCommand() *cobra.Command {
//	var opts localOptions
//	cmd := &cobra.Command{
//		Use:   "local",
//		Short: "Run AWS Glue jobs and test suites locally",
//	}
//
//	flags := cmd.Flags()
//	flags.StringSliceVarP(&opts.arg, "docker-env", "e", nil, "Set Docker environment variables.")
//	flags.StringSliceVarP(&opts.volumes, "docker-volumes", "v", nil, "Mount Docker volumes.")
//	flags.StringSliceVarP(&opts.arg, "docker-arg", "a", nil,
//		"Pass additional arguments to the \"docker run\" command.")
//
//	return cmd
//}
//
//func init() {
//	localCommand.Flags()
//	localCommand.Flags().StringVarP(&localConfig.dockerOptions.credentials, "docker-creds", "C",
//		"", "Set the AWS credentials directory. If not specified, $HOME/.aws is used.")
//	localCommand.Flags().StringSliceVarP(&localConfig.dockerOptions.env, "docker-env", "e",
//		nil, "Set Docker environment variables.")
//	localCommand.Flags().StringSliceVarP(&localConfig.dockerOptions.volumes, "docker-volumes", "v",
//		nil, "Mount Docker volumes.")
//	//localCommand.Flags().StringVarP(&localConfig.workdir, "docker-workdir", "w")
//	localCommand.AddCommand(pipCommand, pySparkCommand, pyTestCommand, sparkSubmitCommand)
//}
