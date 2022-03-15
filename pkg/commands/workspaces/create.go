package workspaces

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/blinkops/blink-go-cli/gen/models"
	"github.com/blinkops/blink-go-cli/pkg/api_responses"
	"github.com/blinkops/blink-go-cli/pkg/consts"
	"github.com/blinkops/blink-go-cli/pkg/utils"
	"github.com/spf13/cobra"
	"io/ioutil"
	"net/http"
)

func CreateWorkspaceCommand() *cobra.Command {
	command := &cobra.Command{
		Use:     "create",
		Aliases: []string{"c", "cr"},
		Short:   "Create a new workspace",
		Long:    `The following command will create a workspace from a given name`,
		Example: "create -n workspaceName -d 'Describe workspace here'",
		RunE:    createWorkspace,
	}

	//command.PersistentFlags().String(consts.WorkspaceIDAutoGenFlagName, "", "Required. workspace ID")
	command.Flags().StringP(consts.NameFlagName, "n", "", "Workspace name")
	command.Flags().StringP(consts.DescriptionFlagName, "d", "", "Workspace description")

	return command
}

func createWorkspace(command *cobra.Command, _ []string) error {

	wsDescription, err := command.Flags().GetString(consts.DescriptionFlagName)
	wsName, err := command.Flags().GetString(consts.NameFlagName)
	if err != nil || wsName == "" {
		return fmt.Errorf("worspace name is required")
	}

	if err := performCreateWorkspace(wsName, wsDescription); err != nil {
		return err
	}

	return nil
}

func performCreateWorkspace(wsName string, wsDescription string) error {
	// { "name": "jon", "description": "desc" }
	workspace := models.ModelsWorkspaceInfo{
		Name:        wsName,
		Description: wsDescription,
	}

	workspaceData, err := json.Marshal(workspace)
	if err != nil {
		return fmt.Errorf("Failed to create workspace data: %s ", err)
	}

	url := utils.GetBaseURL() + "/api/v1/workspaces"

	request, err := utils.NewRequest(http.MethodPost, url, bytes.NewBuffer(workspaceData), map[string]string{
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
		//return ws_id
		var workspaceResponse api_responses.CreateWorkspaceResponse
		if err := json.Unmarshal(responseBody, &workspaceResponse); err != nil {
			return err
		}
		if workspaceResponse.Id != "" {
			fmt.Printf(workspaceResponse.Id)
		}

		return nil
	}

	return fmt.Errorf(string(responseBody))
}
