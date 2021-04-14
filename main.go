package main

import (
	"os"

	"github.com/omegion/go-cli/cmd"

	"github.com/spf13/cobra"
)

// RootCommand is the main entry point of this application.
func RootCommand() *cobra.Command {
	root := &cobra.Command{
		Use:          "go-cli",
		Short:        "Go CLI application template",
		Long:         "Go CLI application template for Go projects.",
		SilenceUsage: true,
	}

	root.AddCommand(cmd.Version())

	return root
}

func main() {
	if err := RootCommand().Execute(); err != nil {
		os.Exit(1)
	}
}
