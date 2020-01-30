package local

import (
	"errors"
	"os"
	"os/exec"
	"path/filepath"
)

const (
	DockerImageName = "sysmana/aws-glue-dev-base"

	DistPackagesVolume    = "aws_glue_dist_packages"
	DistPackagesDirectory = "/usr/local/lib/python2.7/dist-packages"
)

// buildDockerCommand builds an exec.Cmd to run the Docker container with the provided options.
func buildDockerCommand(entrypoint string, options *dockerOptions, args []string) (*exec.Cmd, error) {
	var (
		err  error
		envs []string
		vols []string
	)

	for _, env := range options.env {
		envs = append(envs, "-e "+env)
	}

	for _, vol := range options.volumes {
		vols = append(vols, "-v "+vol)
	}

	credsDir := options.credentials
	if credsDir == "" {
		if credsDir, err = os.UserHomeDir(); err != nil {
			return nil, errors.New("unable to determine home directory")
		}
		credsDir = filepath.Join(credsDir, ".aws")
	}

	wd, err := os.Getwd()

	dockerArgs := []string{"run", "--rm", "-it"}
	dockerArgs = append(dockerArgs, envs...)
	dockerArgs = append(dockerArgs, vols...)
	dockerArgs = append(dockerArgs, "-v", credsDir+":/root/.aws",
		"-v", DistPackagesVolume+":"+DistPackagesDirectory,
		"-v", wd+":/project",
		"--entrypoint", entrypoint, DockerImageName+":0.9")
	dockerArgs = append(dockerArgs, args...)

	return exec.Command("docker", dockerArgs...), nil
}

func buildAndRunDockerCommand(entrypoint string, options *dockerOptions, args []string) error {
	dockerCmd, err := buildDockerCommand(entrypoint, options, args)
	if err != nil {
		return err
	}

	dockerCmd.Stdin = os.Stdin
	dockerCmd.Stdout = os.Stdout
	dockerCmd.Stderr = os.Stderr

	if err = dockerCmd.Start(); err != nil {
		return err
	}

	return dockerCmd.Wait()
}
