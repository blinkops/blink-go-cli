package automations

import (
	"encoding/json"
	"fmt"
	"github.com/blinkops/blink-go-cli/pkg/consts"
	"io/ioutil"
	"net/http"

	"github.com/blinkops/blink-go-cli/pkg/utils"

	"github.com/blinkops/blink-go-cli/gen/models"
	"github.com/spf13/cobra"
)

func ListAutomationsCommand() *cobra.Command {

	command := &cobra.Command{
		Use:     "list",
		Aliases: []string{"ls"},
		Short:   "Get all the available automations",
		Long:    `The following command will request Blink's system to list all available automations`,
		Example: "list",
		RunE:    listAutomations,
	}

	command.PersistentFlags().String(consts.WorkspaceIDAutoGenFlagName, "", "Required. workspace ID")

	return command
}

func performListAutomations(wsID string) (string, error) {

	request, err := utils.NewRequest(http.MethodGet, GetAutomationURL(wsID), nil, nil)
	if err != nil {
		return "", err
	}

	client := &http.Client{}
	response, err := client.Do(request)
	if err != nil {
		return "", err
	}

	defer func() { _ = response.Body.Close() }()

	responseBody, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return "", err
	}

	var responseObject = struct {
		models.APIResponsesPagingInfo
		Results []interface{} `json:"results"`
	}{}

	if err := json.Unmarshal(responseBody, &responseObject); err != nil {
		return "", err
	}

	if responseObject.Results == nil {
		return "No automations are available", nil
	}

	marshaled, err := json.Marshal(responseObject.Results)
	if err != nil {
		return "", err
	}

	return string(marshaled), nil

}

func listAutomations(command *cobra.Command, _ []string) error {

	wsID, err := command.Flags().GetString("ws_id")
	pagingInfo, err := performListAutomations(wsID)
	if err != nil {
		return err
	}

	fmt.Printf("%v", pagingInfo)
	return nil
}
