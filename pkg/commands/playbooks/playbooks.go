package playbooks

import (
	"fmt"

	"github.com/blinkops/blink-go-cli/pkg/commands/requests"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func getSubCommands() []*cobra.Command {
	return []*cobra.Command{
		getListPlaybooksCommand(),
		getCreatePlaybookCommand(),
	}
}

func RegisterCommands(root *cobra.Command) {
	commands := getSubCommands()
	root.AddCommand(commands...)
}

func GetPlaybookURL(workspaceID string) string {
	return requests.GetBaseURL() +
		fmt.Sprintf("/api/v1/workspace/%s/table/playbooks", workspaceID)
}

func getWorkspaceParamFlags(cmd *cobra.Command) string {
	return viper.GetString("workspace-id")
}
