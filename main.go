//go:generate go run gen/generate.go

package main

import (
	"os"

	"github.com/blinkops/blink-go-cli/gen/spec"

	"github.com/blinkops/blink-go-cli/pkg/commands"
	"github.com/blinkops/blink-go-cli/pkg/normalizer"

	"github.com/blinkops/blink-go-cli/gen/cli"
	"github.com/spf13/cobra"
)

const (
	ICON = `
██████╗ ██╗     ██╗███╗   ██╗██╗  ██╗     ██████╗██╗     ██╗
██╔══██╗██║     ██║████╗  ██║██║ ██╔╝    ██╔════╝██║     ██║
██████╔╝██║     ██║██╔██╗ ██║█████╔╝     ██║     ██║     ██║
██╔══██╗██║     ██║██║╚██╗██║██╔═██╗     ██║     ██║     ██║
██████╔╝███████╗██║██║ ╚████║██║  ██╗    ╚██████╗███████╗██║
╚═════╝ ╚══════╝╚═╝╚═╝  ╚═══╝╚═╝  ╚═╝     ╚═════╝╚══════╝╚═╝`
)

var rootCmd *cobra.Command

func main() {

	var err error

	rootCmd, err = cli.MakeRootCmd()
	cobra.CheckErr(err)

	rootCmd.Long = ICON

	swaggerSpec, err := spec.GetSwaggerSpec()
	cobra.CheckErr(err)

	// Add the children commands
	parentCommands := rootCmd.Commands()
	childCommands := commands.GetRegisteredChildCommands()

	for i := range parentCommands {
		parent := parentCommands[i]
		if children, found := childCommands[parent.Use]; found {
			parent.AddCommand(children...)
		}
	}

	// Add the standalone commands
	rootCmd.AddCommand(
		commands.GetRegisteredStandaloneCommands()...,
	)

	normalizer.NormalizeCommands(rootCmd, swaggerSpec)
	normalizer.NormalizeFlags(rootCmd)

	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}

}
