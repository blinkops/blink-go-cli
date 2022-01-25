package playbooks

import (
	"errors"
	"fmt"
	"github.com/blinkops/blink-go-cli/pkg/consts"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/blinkops/blink-go-cli/pkg/utils"

	"github.com/spf13/cobra"
)

func DeletePlaybookCommand() *cobra.Command {
	command := &cobra.Command{
		Use:     "delete",
		Aliases: []string{"d"},
		Short:   "Delete playbook by name or id",
		Long:    `The following command will delete a playbook by name or id`,
		Example: "delete --name my_playbook",
		RunE:    deletePlaybook,
	}

	command.PersistentFlags().String(consts.WorkspaceIDAutoGenFlagName, "", "Required. workspace ID")
	command.Flags().StringP(consts.NameFlagName, "n", "", "The name of the playbook")
	command.Flags().StringP(consts.IDFlagName, "i", "", "The id of the playbook")

	return command
}

func performDeletePlaybookById(playbookID, wsID string) error {

	url := utils.GetBaseURL() + fmt.Sprintf("/api/v1/workspace/%s/playbooks/%s", wsID, playbookID)
	request, err := utils.NewRequest(http.MethodDelete, url, nil, nil)
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

	if response.StatusCode != http.StatusOK {
		return errors.New(string(responseBody))
	}

	fmt.Printf(strings.Trim(string(responseBody), "\""))
	return nil

}

func deletePlaybook(command *cobra.Command, _ []string) error {

	wsID, err := command.Flags().GetString(consts.WorkspaceIDAutoGenFlagName)
	name, err := command.Flags().GetString(consts.NameFlagName)
	id, err := command.Flags().GetString(consts.IDFlagName)
	if err != nil || (name == "" && id == "") {
		return fmt.Errorf("no name or id is supplied")
	}

	// if both name and id are supplied, name takes priority
	if name != "" {
		if id, err = getPlaybookId(name, wsID); err != nil {
			return err
		}
	}

	if err := performDeletePlaybookById(id, wsID); err != nil {
		return err
	}

	return nil
}
