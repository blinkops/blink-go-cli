package normalizer

import (
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// better name normalizer

// NormalizeFlags formats the cobra flags according to the
func NormalizeFlags(root *cobra.Command) {

	// disable this stuff
	for _, flag := range []string{"BLINK-API-KEY", "hostname", "scheme"} {
		flag := root.PersistentFlags().Lookup(flag)
		if flag != nil {
			flag.Hidden = true
		}
	}

	for _, val := range root.Commands() {
		for _, subCmds := range val.Commands() {
			ws := subCmds.PersistentFlags().Lookup("ws_id")
			if ws != nil {
				// bindWorkspace needs to run after the global flags are parsed
				// but before the command is run
				subCmds.PersistentPreRun = bindWorkspace
				ws.Changed = true
				ws.Hidden = true
			}
		}
	}
}

func bindWorkspace(cmd *cobra.Command, args []string) {

	workspaceID := viper.GetString("workspace-id")
	ws := cmd.PersistentFlags().Lookup("ws_id")
	// this should never happen but let's play it safe
	if ws != nil {
		ws.Value.Set(workspaceID)
	}

}
