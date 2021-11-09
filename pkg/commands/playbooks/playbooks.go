package playbooks

import (
	"bytes"
	"fmt"

	"io/ioutil"
	"os"

	"github.com/blinkops/blink-go-cli/pkg/consts"

	"github.com/blinkops/blink-go-cli/gen/models"
	"github.com/blinkops/blink-go-cli/pkg/utils"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"gopkg.in/yaml.v2"
)

func GetPlaybookURL(workspaceID string) string {
	return utils.GetBaseURL() +
		fmt.Sprintf("/api/v1/workspace/%s/table/playbooks", workspaceID)
}

func getWorkspaceParamFlags(cmd *cobra.Command) string {
	return viper.GetString(consts.WorkspaceIDCobraKey)
}

func readPlaybookFile(filePath string) (playbook models.ModelsPlaybook, err error) {

	playbookPayload := make(map[string]interface{})

	if _, err := os.Stat(filePath); err != nil {
		return playbook, fmt.Errorf("%s does not exist", filePath)
	}

	data, err := ioutil.ReadFile(filePath)
	if err != nil {
		return playbook, err
	}

	// this needs to be the playbook object from the server
	//playbookObject := &models.ModelsPlaybook{}
	if err := yaml.NewDecoder(bytes.NewBuffer(data)).Decode(playbookPayload); err != nil {
		return playbook, fmt.Errorf("invalid playbook file: %s", err)
	}

	version, _ := playbookPayload["version"].(string)
	tagsI, _ := playbookPayload["tags"].([]interface{})
	name, _ := playbookPayload["name"].(string)

	var tags []string
	for _, tagI := range tagsI {
		if tag, ok := tagI.(string); ok {
			tags = append(tags, tag)
		}
	}

	playbook = models.ModelsPlaybook{
		Version:         version,
		Playbook:        string(data),
		NumOfExecutions: 0,
		Tags:            tags,
		Name:            name,
	}

	return playbook, nil

}
