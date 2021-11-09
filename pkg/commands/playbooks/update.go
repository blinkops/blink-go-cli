package playbooks

import (
	"fmt"
	"github.com/blinkops/blink-go-cli/pkg/consts"
	"github.com/spf13/cobra"
)

func UpdatePlaybooksCommand() *cobra.Command {

	command := &cobra.Command{
		Use:     "update",
		Short:   "Update playbook from file",
		Long:    `The following command will update a playbook from a given YAML file`,
		Example: "update -f /path/to/playbook.yaml",
		RunE:    updatePlaybooks,
	}
	command.PersistentFlags().String(consts.WorkspaceIDAutoGenFlagName, "", "Required. workspace ID")

	return command
}

func updatePlaybooks(command *cobra.Command, _ []string) error {
	//wsID, err := command.Flags().GetString("ws_id")
	filePath, err := command.Flags().GetString("file")
	if err != nil || filePath == "" {
		return fmt.Errorf("no file input is supplied for the playbook creation")
	}

	readPlaybookFile(filePath)
	return nil
}
