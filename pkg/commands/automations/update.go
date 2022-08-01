package automations

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

func UpdateAutomationsCommand() *cobra.Command {
	command := &cobra.Command{
		Use:     "update",
		Short:   "Update automation from file",
		Long:    `The following command will update an automation from a given YAML file`,
		Example: "update -f /path/to/automation.yaml",
		RunE:    updateAutomations,
	}
	command.PersistentFlags().String(consts.WorkspaceIDAutoGenFlagName, "", "Required. workspace ID")
	command.Flags().StringP(consts.FileFlagName, "f", "", "The path to the automation file")
	command.Flags().StringP(consts.AutomationPackFlag, consts.AutomationPackShortFlag, "", "Name of an automation pack to put the updated automation in")
	command.Flags().BoolP(consts.PublishFlag, "a", true, "Publish and Activate the automation")

	return command
}

func updateAutomations(command *cobra.Command, _ []string) error {
	wsID, err := command.Flags().GetString(consts.WorkspaceIDAutoGenFlagName)
	if err != nil {
		return err
	}
	filePath, err := command.Flags().GetString(consts.FileFlagName)
	if err != nil {
		return err
	}
	if filePath == "" {
		return fmt.Errorf("no file input is supplied for the automation update")
	}
	packName, err := command.Flags().GetString(consts.AutomationPackFlag)
	if err != nil {
		return err
	}

	automationObj, err := readAutomationFile(filePath)
	if err != nil {
		return err
	}

	r := utils.NewTransport()

	searchParam := automations.NewAutomationFindByFilterParams()
	searchParam.Q = fmt.Sprintf(`{"search":{"text":"%s","fields":["name"]}}`, automationObj.Name)
	searchParam.WsID = wsID

	automationResponse, err := client.New(r, strfmt.Default).
		Automations.AutomationFindByFilter(searchParam, nil)
	if err != nil {
		return err
	}

	for _, val := range automationResponse.Payload.Results {
		if val.Name == automationObj.Name {
			automationObj.ID = val.ID
			automationObj.Active = val.Active
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
		automationObj.Active = published
	}

	if automationObj.ID == "" {
		return fmt.Errorf("could not find automation [%s]", automationObj.Name)
	}

	packId, err := resolveAutomationPackId(packName, wsID)
	if err != nil {
		return err
	}
	automationObj.PackID = packId

	updateParam := automations.NewUpdateAutomationParams()
	updateParam.ID = automationObj.ID
	updateParam.WsID = wsID
	updateParam.Automation = &automationObj

	_, err = client.New(r, strfmt.Default).
		Automations.UpdateAutomation(updateParam, nil)

	if err != nil {
		return err
	}

	// print the automation id for automation purposes
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
