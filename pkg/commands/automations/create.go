package automations

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/blinkops/blink-go-cli/pkg/api_responses"
	"github.com/blinkops/blink-go-cli/pkg/consts"
	"github.com/blinkops/blink-go-cli/pkg/utils"

	"github.com/spf13/cobra"
)

func CreateAutomationCommand() *cobra.Command {
	command := &cobra.Command{
		Use:     "create",
		Aliases: []string{"c", "cr"},
		Short:   "Create automation by file",
		Long:    `The following command will create a automation from a given YAML file`,
		Example: "create -f /path/to/automation.yaml",
		RunE:    createAutomation,
	}

	command.PersistentFlags().String(consts.WorkspaceIDAutoGenFlagName, "", "Required. workspace ID")
	command.Flags().StringP(consts.FileFlagName, "f", "", "The path to the automation file")
	command.Flags().StringP(consts.AutomationPackFlag, consts.AutomationPackShortFlag, "", "Name of an automation pack to create the automation in")
	command.Flags().BoolP(consts.PublishFlag, "a", true, "Publish and Activate the automation")

	return command
}

func performCreateAutomation(filePath, wsID, packName string, publish bool) error {
	automation, err := readAutomationFile(filePath)
	if err != nil {
		return err
	}

	automation.Active = publish
	packId, err := resolveAutomationPackId(packName, wsID)
	if err != nil {
		return err
	}
	automation.PackID = packId

	automationData, err := json.Marshal(automation)
	if err != nil {
		return fmt.Errorf("Failed to create automation data: %s ", err)
	}

	request, err := utils.NewRequest(http.MethodPost, GetAutomationURL(wsID),
		bytes.NewBuffer(automationData), map[string]string{
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
		if err != nil {
			return err
		}
		var automationResponse api_responses.CreateResponseWithId
		if err := json.Unmarshal(responseBody, &automationResponse); err != nil {
			return err
		}
		if automationResponse.Id != "" {
			fmt.Printf(automationResponse.Id)
		}
		return nil
	}

	return fmt.Errorf(string(responseBody))
}

func createAutomation(command *cobra.Command, _ []string) error {
	wsID, err := command.Flags().GetString(consts.WorkspaceNameFlagName)
	if err != nil {
		return err
	} else if wsID != "" {
		err = checkExistingWorkspace(wsID)
		if err != nil {
			return err
		}
	}

	filePath, err := command.Flags().GetString(consts.FileFlagName)
	if err != nil {
		return err
	}
	if filePath == "" {
		return fmt.Errorf("no file input is supplied for the automation creation")
	}
	packName, err := command.Flags().GetString(consts.AutomationPackFlag)
	if err != nil {
		return err
	}

	published, err := command.Flags().GetBool(consts.PublishFlag)
	if err != nil {
		return err
	}

	if err := performCreateAutomation(filePath, wsID, packName, published); err != nil {
		return err
	}

	return nil
}

func checkExistingWorkspace(wsID string) error {
	request, err := utils.NewRequest(http.MethodGet, GetFindWorkspaceURL(wsID), nil, nil)
	if err != nil {
		return err
	}

	resp, err := http.DefaultClient.Do(request)
	if err != nil {
		return err
	}
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("failed to find workspace with ID %s, status code %d returned", wsID, resp.StatusCode)
	}
	return nil
}
