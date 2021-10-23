package main

import (
	"fmt"
	"os"
	"regexp"
	"strings"

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

	SetupOperations(rootCmd.Commands())

	// maybe just check the config file?
	if true {
		setGlobalWorkspace()
	}

	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}

}

func setGlobalWorkspace() {
	for _, val := range rootCmd.Commands() {
		for _, subCmds := range val.Commands() {
			ws := subCmds.PersistentFlags().Lookup("ws_id")
			if ws != nil {
				ws.Changed = true
				ws.Value.Set("global_ws_id_elie")
				ws.Hidden = true
			}
		}
	}
}

func unsetGlobalWorkspace() {
	for _, val := range rootCmd.Commands() {
		for _, subCmds := range val.Commands() {
			ws := subCmds.PersistentFlags().Lookup("ws_id")
			if ws != nil {
				ws.Changed = false
				ws.Value.Set("")
				ws.Hidden = false
			}
		}
	}
}

func SetupOperations(operations []*cobra.Command) {
	for o := range operations {
		operation := operations[o]
		// always keep this
		if operation.Use == "completion [bash|zsh|fish|powershell]" {
			continue
		}
		if !isOperationEnabled(operation.Use) {
			operation.Hidden = true
			continue
		}
		commands := operation.Commands()
		for c := range commands {
			command := commands[c]
			if !isCommandEnabled(operation.Use, command.Use) {
				// hide the command
				command.Hidden = true
				continue
			}
			stripped := removeGroupNameFromOperation(operation.Use, command.Use)
			command.Use = toSnakeCase(stripped)
		}
	}
}

// provide example with configured
func isOperationEnabled(op string) bool {
	spaced := strings.ReplaceAll(op, "_", " ")
	titled := strings.Title(spaced)
	_, exists := OpenAPIEnabledOperations[titled]
	return exists
}

func isCommandEnabled(op, cmd string) bool {
	spaced := strings.ReplaceAll(op, "_", " ")
	titled := strings.Title(spaced)
	_, exists := OpenAPIEnabledOperations[titled][cmd]
	return exists
}

func removeGroupNameFromOperation(groupName, operation string) string {
	// operations are in lowercase, to match we need to convert
	return strings.TrimPrefix(operation, strings.ToUpper(operation))
}

func toSnakeCase(str string) string {
	snake := matchFirstCap.ReplaceAllString(str, "${1}_${2}")
	snake = matchAllCap.ReplaceAllString(snake, "${1}_${2}")
	return strings.ToLower(snake)
}
