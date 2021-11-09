package playbooks

import (
	"fmt"

	"github.com/blinkops/blink-go-cli/gen/client/playbooks"
	httptransport "github.com/go-openapi/runtime/client"
	"github.com/go-openapi/strfmt"
	"github.com/spf13/viper"

	"github.com/blinkops/blink-go-cli/gen/client"
	"github.com/blinkops/blink-go-cli/pkg/consts"
	"github.com/spf13/cobra"
)

func UpdatePlaybooksCommand() *cobra.Command {

	command := &cobra.Command{
		Use:     "update",
		Short:   "Update playbook from file",
		Long:    `The following command will update a playbook from a given YAML file`,
		Example: "update -f /path/to/playbook.yaml",
		RunE:    updatePlaybooks,
	}
	command.PersistentFlags().String(consts.WorkspaceIDAutoGenFlagName, "", "Required. workspace ID")
	command.Flags().StringP(consts.FileFlagName, "f", "", "The path to the playbook file")

	return command
}

func updatePlaybooks(command *cobra.Command, _ []string) error {

	wsID, err := command.Flags().GetString("ws_id")
	filePath, err := command.Flags().GetString("file")
	if err != nil || filePath == "" {
		return fmt.Errorf("no file input is supplied for the playbook creation")
	}

	playbookObj, err := readPlaybookFile(filePath)
	if err != nil {
		return err
	}

	hostname := viper.GetString(consts.HostnameEntry)
	scheme := viper.GetString(consts.SchemeEntry)
	apiKey := viper.GetString(consts.ApiKeyEntry)

	if scheme == "" {
		scheme = "https"
	}

	r := httptransport.New(hostname, client.DefaultBasePath, []string{scheme})
	r.DefaultAuthentication = httptransport.Compose(
		httptransport.APIKeyAuth(consts.ApiKeyHeader, "header", apiKey),
	)

	searchParam := playbooks.NewPlaybookFindByFilterParams()
	searchParam.Q = fmt.Sprintf(`{"search":{"text":"%s"}}`, playbookObj.Name)
	searchParam.WsID = wsID

	playbookResponse, err := client.New(r, strfmt.Default).
		Playbooks.PlaybookFindByFilter(searchParam, nil)

	if err != nil {
		return err
	}

	for _, val := range playbookResponse.Payload.Results {
		if val.Name == playbookObj.Name {
			playbookObj.ID = val.ID
			break
		}
	}

	if playbookObj.ID == "" {
		return fmt.Errorf("could not find playbook [%s]", playbookObj.Name)
	}

	updateParam := playbooks.NewUpdatePlaybookParams()
	updateParam.ID = playbookObj.ID
	updateParam.Playbook = &playbookObj

	_, err = client.New(r, strfmt.Default).
		Playbooks.UpdatePlaybook(updateParam, nil)

	if err != nil {
		return err
	}

	return nil
}
