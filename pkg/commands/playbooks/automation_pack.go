package playbooks

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/blinkops/blink-go-cli/gen/models"
	"github.com/blinkops/blink-go-cli/pkg/api_responses"
	"github.com/blinkops/blink-go-cli/pkg/consts"
	"github.com/blinkops/blink-go-cli/pkg/utils"
)

func checkExistingAutomationPack(packName, wsID string) (*api_responses.GetIdByNameResponse, error) {
	query := fmt.Sprintf(consts.GetAutomationPackByDisplayNameQueryFormat, packName)
	request, err := utils.NewRequest(http.MethodGet, GetFindAutomationPackURL(wsID, query), nil, nil)
	if err != nil {
		return nil, err
	}
	resp, err := http.DefaultClient.Do(request)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("invalid status code %d returned from blink when looking for automation pack named '%s'", resp.StatusCode, packName)
	}
	result := &api_responses.GetIdByNameResponse{}
	if err = json.NewDecoder(resp.Body).Decode(result); err != nil {
		return nil, err
	}
	return result, nil
}

func createAutomationPack(packName, wsID string) (string, error) {
	creationPayload := &models.APIRequestsAutomationPackPayload{
		Name: packName,
	}
	creationPayloadBytes, err := json.Marshal(creationPayload)
	if err != nil {
		return "", err
	}
	request, err := utils.NewRequest(http.MethodPost, GetCreateAutomationPackURL(wsID), bytes.NewBuffer(creationPayloadBytes), map[string]string{
		"Content-Type": "application/json",
	})
	resp, err := http.DefaultClient.Do(request)
	if err != nil {
		return "", err
	}
	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("invalid status code %d returned from blink when creating automation pack named '%s'", resp.StatusCode, packName)
	}
	result := &api_responses.CreateResponseWithId{}
	if err = json.NewDecoder(resp.Body).Decode(result); err != nil {
		return "", err
	}
	return result.Id, nil
}

func resolveAutomationPackId(packName, wsID string) (string, error) {
	if packName == "" {
		// no pack name so let the logic on the controller put the playbook in a 'default' pack in case of creation
		// or keep it in the same pack in case of update
		return "", nil
	}
	result, err := checkExistingAutomationPack(packName, wsID)
	if err != nil {
		return "", nil
	}
	if len(result.Results) > 0 {
		return result.Results[0].Id, nil
	}
	id, err := createAutomationPack(packName, wsID)
	if err != nil {
		return "", err
	}
	return id, nil
}
