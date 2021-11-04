package initialize

import (
	"github.com/spf13/cobra"
)

func CMD() *cobra.Command {
	command := &cobra.Command{
		Use:     "init",
		Aliases: []string{"i"},
		Short:   "Initialize a blink config file",
		Long:    `This utility will walk you through setting up the blink configuration`,
		RunE:    setupMenuOptions,
	}
	return command
}

