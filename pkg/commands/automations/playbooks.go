package automations

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"

	"github.com/blinkops/blink-go-cli/gen/models"
	"github.com/blinkops/blink-go-cli/pkg/api_responses"
	"github.com/blinkops/blink-go-cli/pkg/consts"
	"github.com/blinkops/blink-go-cli/pkg/utils"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"gopkg.in/yaml.v2"
)

func GetAutomationURL(workspaceID string) string {
	return utils.GetBaseURL() +
		fmt.Sprintf("/api/v1/workspace/%s/automations", workspaceID)
}

func GetFindAutomationPackURL(workspaceID, query string) string {
	return fmt.Sprintf("%s/api/v1/workspace/%s/table/automation_packs?q=%s", utils.GetBaseURL(), workspaceID, query)
}

func GetFindWorkspaceURL(workspaceID string) string {
	return fmt.Sprintf("%s/api/v1/workspaces/%s", utils.GetBaseURL(), workspaceID)
}

func GetCreateAutomationPackURL(workspaceID string) string {
	return fmt.Sprintf("%s/api/v1/workspace/%s/automation_packs", utils.GetBaseURL(), workspaceID)
}

func getWorkspaceParamFlags(cmd *cobra.Command) string {
	return viper.GetString(consts.WorkspaceIDCobraKey)
}

func readAutomationFile(filePath string) (automation models.ModelsAutomation, err error) {
	automationPayload := make(map[string]interface{})

	if _, err := os.Stat(filePath); err != nil {
		return automation, fmt.Errorf("%s does not exist", filePath)
	}

	data, err := ioutil.ReadFile(filePath)
	if err != nil {
		return automation, err
	}

	// this needs to be the automation object from the server
	// automationObject := &models.ModelsAutomation{}
	if err := yaml.NewDecoder(bytes.NewBuffer(data)).Decode(automationPayload); err != nil {
		return automation, fmt.Errorf("invalid automation file: %s", err)
	}

	version, _ := automationPayload["version"].(string)
	tagsI, _ := automationPayload["tags"].([]interface{})
	name, _ := automationPayload["name"].(string)
	automationType, _ := automationPayload["type"].(string)

	var tags []string
	for _, tagI := range tagsI {
		if tag, ok := tagI.(string); ok {
			tags = append(tags, tag)
		}
	}

	automation = models.ModelsAutomation{
		Type:       automationType,
		Version:    version,
		Automation: string(data),
		Tags:       tags,
		Name:       name,
	}

	return automation, nil
}

func extractAutomationIdFromResponse(responseBody []byte, automationName string) (string, error) {
	var result api_responses.GetIdByNameResponse
	if err := json.Unmarshal(responseBody, &result); err != nil {
		return "", err
	}

	if result.Results != nil && len(result.Results) > 0 && result.Results[0].Id != "" {
		return result.Results[0].Id, nil
	}

	return "", fmt.Errorf("cannot find automation id for automation [%s]", automationName)
}

func getAutomationIdByName(automationName string, workspaceID string) (string, error) {
	filter := fmt.Sprintf(`{"limit": 1, "offset": 0, "filter": {"name": {"$eq": "%s"}}, "select": ["id"]}`, automationName)
	url := utils.GetBaseURL() + fmt.Sprintf("/api/v1/workspace/%s/automations?q=%s", workspaceID, url.QueryEscape(filter))
	request, err := utils.NewRequest(http.MethodGet, url, nil, nil)
	if err != nil {
		return "", err
	}

	client := &http.Client{}
	response, err := client.Do(request)
	if err != nil {
		return "", err
	}

	if response.Body == nil {
		return "", errors.New("invalid response body")
	}

	defer func() { _ = response.Body.Close() }()

	responseBody, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return "", err
	}

	if response.StatusCode != http.StatusOK {
		return "", nil
	}

	return extractAutomationIdFromResponse(responseBody, automationName)
}
