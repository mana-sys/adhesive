package main

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/mana-sys/adhesive/internal/cli"
	"github.com/mana-sys/adhesive/internal/cli/command"
)

func main() {
	adhesiveCli, err := command.NewAdhesiveCli("")
	if err != nil {
		fmt.Fprintf(os.Stderr, "%s: %s\n", os.Args[0], err)
		os.Exit(1)
	}

	cmd := cli.NewRootCommand(adhesiveCli)

	if err := cmd.Execute(); err != nil {
		if exitErr, ok := err.(*exec.ExitError); ok {
			os.Exit(exitErr.ExitCode())
		}

		fmt.Fprintf(os.Stderr, "%s: %s\n", os.Args[0], err)
		os.Exit(1)
	}
}
