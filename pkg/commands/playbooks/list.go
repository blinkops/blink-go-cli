package playbooks

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/blinkops/blink-go-cli/gen/models"
	"github.com/blinkops/blink-go-cli/pkg/commands/requests"
	"github.com/blinkops/blink-go-cli/pkg/consts"
	"github.com/spf13/cobra"
)

func getListPlaybooksCommand() *cobra.Command {

	command := &cobra.Command{
		Use:     "list",
		Aliases: []string{"ls"},
		Short:   "Get all the available playbooks",
		Long:    `The following command will request Blink's system to list all available playbooks`,
		Example: "list",
		RunE:    ListPlaybooks,
	}

	command.Flags().StringP(consts.FileFlag, "f", "", "The path to the playbook file")

	return command
}

func performListPlaybooks(wsID string) (string, error) {

	request, err := requests.NewRequest(http.MethodGet, GetPlaybookURL(wsID), nil, nil)
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

	indentedOutput, err := json.MarshalIndent(responseObject.Results, "", consts.MarshalIndentation)
	if err != nil {
		return "", err
	}

	return string(indentedOutput), nil
}

func ListPlaybooks(command *cobra.Command, _ []string) error {

	wsID := getWorkspaceParamFlags(command)
	pagingInfo, err := performListPlaybooks(wsID)
	if err != nil {
		return err
	}

	fmt.Printf("%v", pagingInfo)
	return nil
}
