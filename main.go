//go:generate go run gen/generate.go

package main

import (
	"fmt"
	"os"
	"regexp"

	"github.com/blinkops/blink-go-cli/pkg/commands/initialize"

	"github.com/blinkops/blink-go-cli/gen/spec"
	"github.com/blinkops/blink-go-cli/pkg/commands/playbooks"
	"github.com/blinkops/blink-go-cli/pkg/formatter"

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

var matchFirstCap = regexp.MustCompile("(.)([A-Z][a-z]+)")
var matchAllCap = regexp.MustCompile("([a-z0-9])([A-Z])")

var rootCmd *cobra.Command

func main() {

	var err error

	rootCmd, err = cli.MakeRootCmd()
	if err != nil {
		fmt.Println("Cmd construction error: ", err)
		os.Exit(1)
	}

	rootCmd.Long = ICON

	spec, err := spec.GetSwaggerSpec()
	if err != nil {
		panic(err)
	}

	formatter.CMDFormat(rootCmd, spec)
	formatter.FlagFormat(rootCmd)

	cmds := rootCmd.Commands()
	for i := range cmds {
		if cmds[i].Use == "playbooks" {
			playbooks.RegisterCommands(cmds[i])
		}
	}

	rootCmd.AddCommand(initialize.CMD())

	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}

}
