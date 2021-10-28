package playbooks

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"

	"github.com/blinkops/blink-go-cli/gen/models"
	"github.com/blinkops/blink-go-cli/pkg/commands/requests"
	"github.com/blinkops/blink-go-cli/pkg/consts"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v2"
)

func getCreatePlaybookCommand() *cobra.Command {
	command := &cobra.Command{
		Use:     "create",
		Aliases: []string{"c", "cr"},
		Short:   "Create playbook by file",
		Long:    `The following command will request Blink's system to create a playbook by a given YAML file`,
		Example: "create -f /path/to/playbook.yaml",
		RunE:    CreatePlaybook,
	}

	command.Flags().StringP(consts.FileFlag, "f", "", "The path to the playbook file")

	return command
}

func performCreatePlaybook(filePath, wsID string) error {
	if _, err := os.Stat(filePath); err != nil {
		return fmt.Errorf("%s does not exist", filePath)
	}

	data, err := ioutil.ReadFile(filePath)
	if err != nil {
		return err
	}

	playbookObject := &models.ModelsPlaybook{}
	if err := yaml.NewDecoder(bytes.NewBuffer(data)).Decode(playbookObject); err != nil {
		return fmt.Errorf("invalid playbook file: %s", err)
	}

	playbookAsYaml, err := yaml.Marshal(playbookObject)
	if err != nil {
		return fmt.Errorf("failed to create playbook object data: %s", err)
	}

	playbook := &models.ModelsPlaybook{
		Version:         playbookObject.Version,
		Playbook:        string(playbookAsYaml),
		NumOfExecutions: 0,
		Tags:            playbookObject.Tags,
	}
	playbook.Name = playbookObject.Name

	playbookData, err := json.Marshal(playbook)
	if err != nil {
		return fmt.Errorf("Failed to create playbook data: %s ", err)
	}

	request, err := requests.NewRequest(http.MethodPost, GetPlaybookURL(wsID),
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

func CreatePlaybook(command *cobra.Command, _ []string) error {

	wsID := getWorkspaceParamFlags(command)
	filePath, err := command.Flags().GetString(consts.FileFlag)
	if err != nil || filePath == "" {
		return fmt.Errorf("no file input is supplied for the playbook creation")
	}

	if err := performCreatePlaybook(filePath, wsID); err != nil {
		return err
	}

	return nil
}
