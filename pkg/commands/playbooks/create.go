package playbooks

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/blinkops/blink-go-cli/pkg/consts"
	"io/ioutil"
	"net/http"

	"github.com/blinkops/blink-go-cli/pkg/utils"

	"github.com/spf13/cobra"
)

func CreatePlaybookCommand() *cobra.Command {
	command := &cobra.Command{
		Use:     "create",
		Aliases: []string{"c", "cr"},
		Short:   "Create playbook by file",
		Long:    `The following command will create a playbook from a given YAML file`,
		Example: "create -f /path/to/playbook.yaml",
		RunE:    createPlaybook,
	}

	command.PersistentFlags().String(consts.WorkspaceIDAutoGenFlagName, "", "Required. workspace ID")
	command.Flags().StringP(consts.FileFlagName, "f", "", "The path to the playbook file")

	return command
}

func performCreatePlaybook(filePath, wsID string) error {

	playbook, err := readPlaybookFile(filePath)
	if err != nil {
		return err
	}

	playbookData, err := json.Marshal(playbook)
	if err != nil {
		return fmt.Errorf("Failed to create playbook data: %s ", err)
	}

	request, err := utils.NewRequest(http.MethodPost, GetPlaybookURL(wsID),
		bytes.NewBuffer(playbookData), map[string]string{
			"Content-Type": "application/json",
		})
	if err != nil {
		return err
	}

	client := &http.Client{}
	response, err := client.Do(request)
	if err != nil {
		return err
	}

	defer func() { _ = response.Body.Close() }()

	responseBody, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return err
	}

	if response.StatusCode >= http.StatusOK && response.StatusCode < http.StatusBadRequest {
		return nil
	}

	return fmt.Errorf(string(responseBody))

}

func createPlaybook(command *cobra.Command, _ []string) error {

	wsID, err := command.Flags().GetString(consts.WorkspaceIDAutoGenFlagName)
	filePath, err := command.Flags().GetString(consts.FileFlagName)
	if err != nil || filePath == "" {
		return fmt.Errorf("no file input is supplied for the playbook creation")
	}

	if err := performCreatePlaybook(filePath, wsID); err != nil {
		return err
	}

	return nil
}
