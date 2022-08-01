package automations

import (
	"errors"
	"fmt"
	"github.com/blinkops/blink-go-cli/pkg/consts"
	"io/ioutil"
	"net/http"

	"github.com/blinkops/blink-go-cli/pkg/utils"

	"github.com/spf13/cobra"
)

func GetAutomationCommand() *cobra.Command {
	command := &cobra.Command{
		Use:     "get",
		Aliases: []string{"g"},
		Short:   "Get automation by name or id",
		Long:    `The following command will get an automation by name or id`,
		Example: "get --name my_automation",
		RunE:    getAutomation,
	}

	command.PersistentFlags().String(consts.WorkspaceIDAutoGenFlagName, "", "Required. workspace ID")
	command.Flags().StringP(consts.NameFlagName, "n", "", "The name of the automation")
	command.Flags().StringP(consts.IDFlagName, "i", "", "The id of the automation")

	return command
}

func performGetAutomationById(automationID, wsID string) error {

	url := utils.GetBaseURL() + fmt.Sprintf("/api/v1/workspace/%s/playbooks/%s", wsID, automationID)
	request, err := utils.NewRequest(http.MethodGet, url, nil, nil)
	if err != nil {
		return err
	}

	client := &http.Client{}
	response, err := client.Do(request)
	if err != nil {
		return err
	}

	if response.Body == nil {
		return errors.New("invalid response body")
	}

	defer func() { _ = response.Body.Close() }()

	responseBody, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return err
	}

	if response.StatusCode != http.StatusOK {
		return errors.New(string(responseBody))
	}

	fmt.Printf(string(responseBody))
	return nil

}

func getAutomation(command *cobra.Command, _ []string) error {

	wsID, err := command.Flags().GetString(consts.WorkspaceIDAutoGenFlagName)
	name, err := command.Flags().GetString(consts.NameFlagName)
	id, err := command.Flags().GetString(consts.IDFlagName)
	if err != nil || (name == "" && id == "") {
		return fmt.Errorf("no name or id is supplied")
	}

	// if both name and id are supplied, name takes priority
	if name != "" {
		if id, err = getAutomationIdByName(name, wsID); err != nil {
			return err
		}
	}

	if err := performGetAutomationById(id, wsID); err != nil {
		return err
	}

	return nil
}
