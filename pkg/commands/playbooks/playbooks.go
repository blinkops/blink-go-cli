package playbooks

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
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
	playbookType, _ := playbookPayload["type"].(string)

	var tags []string
	for _, tagI := range tagsI {
		if tag, ok := tagI.(string); ok {
			tags = append(tags, tag)
		}
	}

	playbook = models.ModelsPlaybook{
		Type:     playbookType,
		Version:  version,
		Playbook: string(data),
		Tags:     tags,
		Name:     name,
	}

	return playbook, nil

}

func extractPlaybookIdFromResponse(responseBody []byte, playbookName string) (string, error) {
	type Result struct {
		Results []struct {
			Id              string        `json:"id"`
		} `json:"results"`
	}

	var result Result
	err := json.Unmarshal(responseBody, &result)
	if err != nil {
		return "", err
	}

	if result.Results != nil && len(result.Results) > 0 && result.Results[0].Id != "" {
		return result.Results[0].Id, nil
	}

	return "", fmt.Errorf("cannot find playbook id for playbook [%s]", playbookName)
}

func getPlaybookId(playbookName string, workspaceID string) (string, error) {
	filter := fmt.Sprintf(`{"limit": 1, "offset": 0, "filter": {"name": {"$eq": "%s"}}}`, playbookName)
	url := utils.GetBaseURL() + fmt.Sprintf("/api/v1/workspace/%s/playbooks?q=%s", workspaceID, url.QueryEscape(filter))
	request, err := utils.NewRequest(http.MethodGet, url, nil, nil)
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

	if response.StatusCode != http.StatusOK {
		return "", nil
	}

	return extractPlaybookIdFromResponse(responseBody, playbookName)
}