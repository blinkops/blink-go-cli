package runner_groups

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/blinkops/blink-go-cli/gen/models"
	"github.com/blinkops/blink-go-cli/pkg/consts"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/blinkops/blink-go-cli/pkg/utils"

	"github.com/spf13/cobra"
)

func CreateRunnerGroupCommand() *cobra.Command {
	command := &cobra.Command{
		Use:     "create",
		Aliases: []string{"c"},
		Short:   "Create a new runner group",
		Long:    `The following command will create a new runner group`,
		Example: "create --name my_runner_group --is-default true",
		RunE:    createRunnerGroup,
	}

	command.PersistentFlags().String(consts.WorkspaceIDAutoGenFlagName, "", "Required. workspace ID")
	command.Flags().StringP(consts.NameFlagName, "n", "", "The name for the runner group")
	command.Flags().BoolP(consts.IsDefaultFlagName, "", false, "Set the runner group as default")
	command.Flags().StringP(consts.TagsFlagName, "t", "", "Tags for the runner group, comma separated")

	err := cobra.MarkFlagRequired(command.Flags(), consts.NameFlagName)
	if err != nil {
		return nil
	}

	return command
}

func performCreateRunnerGroup(runnerGroupData models.APIRequestsCreateRunnerGroupRequest, wsID string) error {

	payload, err := json.Marshal(runnerGroupData)
	if err != nil {
		return fmt.Errorf("Failed to create runner group data: %s ", err)
	}

	url := utils.GetBaseURL() + fmt.Sprintf("/api/v1/workspace/%s/runner_groups", wsID)
	request, err := utils.NewRequest(http.MethodPost, url,
		bytes.NewBuffer(payload), map[string]string{
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

	if response.StatusCode != http.StatusOK {
		return errors.New(string(responseBody))
	}

	fmt.Printf(string(responseBody))
	return nil

}

func createRunnerGroup(command *cobra.Command, _ []string) error {
	runnerGroupData := models.APIRequestsCreateRunnerGroupRequest{
		IsDefault: false,
		Name:      "",
		Tags:      nil,
	}
	wsID, err := command.Flags().GetString(consts.WorkspaceIDAutoGenFlagName)
	runnerGroupData.Name, err = command.Flags().GetString(consts.NameFlagName)
	runnerGroupData.IsDefault, err = command.Flags().GetBool(consts.IsDefaultFlagName)
	tags, err := command.Flags().GetString(consts.TagsFlagName)
	if err != nil {
		return err
	}

	// handle tags if passed
	if tags != "" {
		runnerGroupData.Tags = strings.Split(tags, ",")
	}

	if err := performCreateRunnerGroup(runnerGroupData, wsID); err != nil {
		return err
	}

	return nil
}
