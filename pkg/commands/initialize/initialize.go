package initialize

import (
	"github.com/blinkops/blink-go-cli/pkg/commands/initialize/menu"
	"github.com/spf13/cobra"
)

func Command() *cobra.Command {
	command := &cobra.Command{
		Use:     "init",
		Aliases: []string{"i"},
		Short:   "Initialize a blink config file",
		Long:    `This utility will walk you through setting up the blink configuration`,
		RunE:    menu.Setup,
	}
	return command
}

