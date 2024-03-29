package playbooks

import (
	"fmt"
	"github.com/blinkops/blink-go-cli/pkg/utils"
	"github.com/spf13/pflag"

	"github.com/blinkops/blink-go-cli/gen/client"
	"github.com/blinkops/blink-go-cli/gen/client/automations"
	"github.com/blinkops/blink-go-cli/pkg/consts"
	"github.com/go-openapi/strfmt"
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
	command.Flags().StringP(consts.AutomationPackFlag, consts.AutomationPackShortFlag, "", "Name of an automation pack to put the updated playbook in")
	command.Flags().BoolP(consts.PublishFlag, "a", true, "Publish and Activate the playbook")

	return command
}

func updatePlaybooks(command *cobra.Command, _ []string) error {
	wsID, err := command.Flags().GetString(consts.WorkspaceIDAutoGenFlagName)
	if err != nil {
		return err
	}
	filePath, err := command.Flags().GetString(consts.FileFlagName)
	if err != nil {
		return err
	}
	if filePath == "" {
		return fmt.Errorf("no file input is supplied for the playbook update")
	}
	packName, err := command.Flags().GetString(consts.AutomationPackFlag)
	if err != nil {
		return err
	}

	playbookObj, err := readPlaybookFile(filePath)
	if err != nil {
		return err
	}

	r := utils.NewTransport()

	searchParam := automations.NewPlaybookFindByFilterParams()
	searchParam.Q = fmt.Sprintf(`{"search":{"text":"%s","fields":["name"]}}`, playbookObj.Name)
	searchParam.WsID = wsID

	playbookResponse, err := client.New(r, strfmt.Default).
		Automations.PlaybookFindByFilter(searchParam, nil)
	if err != nil {
		return err
	}

	for _, val := range playbookResponse.Payload.Results {
		if val.Name == playbookObj.Name {
			playbookObj.ID = val.ID
			playbookObj.Active = val.Active
			break
		}
	}

	// handle publish case
	isSet := isFlagPassed(command, consts.PublishFlag)
	if isSet {
		published, err := command.Flags().GetBool(consts.PublishFlag)
		if err != nil {
			return err
		}
		playbookObj.Active = published
	}

	if playbookObj.ID == "" {
		return fmt.Errorf("could not find playbook [%s]", playbookObj.Name)
	}

	packId, err := resolveAutomationPackId(packName, wsID)
	if err != nil {
		return err
	}
	playbookObj.PackID = packId

	updateParam := automations.NewUpdatePlaybookParams()
	updateParam.ID = playbookObj.ID
	updateParam.WsID = wsID
	updateParam.Playbook = &playbookObj

	_, err = client.New(r, strfmt.Default).
		Automations.UpdatePlaybook(updateParam, nil)

	if err != nil {
		return err
	}

	// print the playbook id for automation purposes
	fmt.Println(updateParam.ID)

	return nil
}

func isFlagPassed(command *cobra.Command, name string) bool {
	found := false
	command.Flags().Visit(func(f *pflag.Flag) {
		if f.Name == name {
			found = true
		}
	})
	return found
}
