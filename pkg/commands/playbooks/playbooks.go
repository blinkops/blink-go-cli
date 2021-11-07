package playbooks

import (
	"fmt"
	"github.com/blinkops/blink-go-cli/pkg/utils"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func GetPlaybookURL(workspaceID string) string {
	return utils.GetBaseURL() +
		fmt.Sprintf("/api/v1/workspace/%s/table/playbooks", workspaceID)
}

func getWorkspaceParamFlags(cmd *cobra.Command) string {
	return viper.GetString("workspace-id")
}
