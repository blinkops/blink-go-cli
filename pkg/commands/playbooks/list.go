package playbooks

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/blinkops/blink-go-cli/pkg/utils"

	"github.com/blinkops/blink-go-cli/gen/models"
	"github.com/spf13/cobra"
)

func ListPlaybooksCommand() *cobra.Command {

	command := &cobra.Command{
		Use:     "list",
		Aliases: []string{"ls"},
		Short:   "Get all the available playbooks",
		Long:    `The following command will request Blink's system to list all available playbooks`,
		Example: "list",
		RunE:    listPlaybooks,
	}

	command.Flags().StringP("file", "f", "", "The path to the playbook file")

	return command
}

func performListPlaybooks(wsID string) (string, error) {

	request, err := utils.NewRequest(http.MethodGet, GetPlaybookURL(wsID), nil, nil)
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
		return "No playbooks are available", nil
	}

	marshaled, err := json.Marshal(responseObject.Results)
	if err != nil {
		return "", err
	}

	return string(marshaled), nil

}

func listPlaybooks(command *cobra.Command, _ []string) error {

	wsID := getWorkspaceParamFlags(command)
	pagingInfo, err := performListPlaybooks(wsID)
	if err != nil {
		return err
	}

	fmt.Printf("%v", pagingInfo)
	return nil
}
