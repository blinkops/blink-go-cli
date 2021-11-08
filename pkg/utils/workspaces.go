package utils

import (
	"fmt"

	"github.com/blinkops/blink-go-cli/gen/client"
	httptransport "github.com/go-openapi/runtime/client"
	"github.com/go-openapi/strfmt"
	"github.com/spf13/viper"
)

func GetWorkspaceID(workspaceName, apiKey string) (string, error) {

	hostname := viper.GetString("hostname")
	scheme := viper.GetString("scheme")

	if scheme == "" {
		scheme = "https"
	}

	r := httptransport.New(hostname, client.DefaultBasePath, []string{scheme})
	r.DefaultAuthentication = httptransport.Compose(
		httptransport.APIKeyAuth("BLINK-API-KEY", "header", apiKey),
	)

	userDetails, err := client.New(r, strfmt.Default).UserInfo.GetUserDetails(nil, nil)
	if err != nil {
		return "", err
	}

	for _, val := range userDetails.Payload.Workspaces {
		if val.Name == workspaceName {
			return val.ID, nil
		}
	}

	return "", fmt.Errorf("cannot find workspace id for workspace [%s]", workspaceName)

}
