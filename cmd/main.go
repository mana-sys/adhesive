package main

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/mana-sys/adhesive/internal/cli/command"

	"github.com/mana-sys/adhesive/internal/cli"
)

func main() {
	adhesiveCli, err := command.NewAdhesiveCli("")
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	cmd := cli.NewRootCommand(adhesiveCli)

	if err := cmd.Execute(); err != nil {
		if exitErr, ok := err.(*exec.ExitError); ok {
			os.Exit(exitErr.ExitCode())
		}

		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
