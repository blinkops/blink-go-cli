package playbooks

import (
	"fmt"

	"github.com/blinkops/blink-go-cli/pkg/commands/requests"
	"github.com/spf13/cobra"
)

func getSubCommands() []*cobra.Command {
	return []*cobra.Command{
		getListPlaybooksCommand(),
		getCreatePlaybookCommand(),
	}
}

func RegisterCommands(root *cobra.Command) {
	commands := getSubCommands()
	// we need to add workspace id to the playbooks endpoints
	for i := range commands {
		registerWorkspaceParamFlags(commands[i])
	}
	root.AddCommand(commands...)
}

func GetPlaybookURL(workspaceID string) string {
	return requests.GetBaseURL() +
		fmt.Sprintf("/api/v1/workspace/%s/table/playbooks", workspaceID)
}

func registerWorkspaceParamFlags(cmd *cobra.Command) error {
	var wsIdFlagDefault string
	_ = cmd.PersistentFlags().String("ws_id", wsIdFlagDefault, "Required. workspace ID")
	return nil
}

func getWorkspaceParamFlags(cmd *cobra.Command) string {
	wsID, _ := cmd.Flags().GetString("ws_id")
	return wsID
}
