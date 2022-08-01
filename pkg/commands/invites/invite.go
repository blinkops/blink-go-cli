package invites

import (
	"encoding/json"

	"github.com/blinkops/blink-go-cli/gen/cli"
	"github.com/blinkops/blink-go-cli/gen/client/invites"
	"github.com/spf13/cobra"
)

// InviteCommand returns a cmd to handle operation invite
func InviteCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "invite",
		Short:   `Inviting a user`,
		RunE:    runOperationInvitesInvite,
		Example: `blink invites invite --invite '[{"email":"user1@company.com"}, {"email":"user2@company.com"}]'`,
	}

	_ = cmd.PersistentFlags().String("invite", "", "json string for [invites]")

	return cmd
}

// runOperationInvitesInvite uses cmd flags to call endpoint api
func runOperationInvitesInvite(cmd *cobra.Command, args []string) error {
	appCli, err := cli.MakeClient(cmd)
	if err != nil {
		return err
	}
	// retrieve flag values from cmd and fill params
	params := invites.NewInviteParams()
	if err := retrieveOperationInvitesInviteInvitationFlag(params, cmd); err != nil {
		return err
	}

	_, err = appCli.Invites.Invite(params, nil)
	if err != nil {
		return err
	}
	return nil
}

func retrieveOperationInvitesInviteInvitationFlag(m *invites.InviteParams, cmd *cobra.Command) error {

	automationFlagValue, err := cmd.Flags().GetString("invite")
	if err != nil {
		return err
	}
	err = json.Unmarshal([]byte(automationFlagValue), &m.Invitation)
	if err != nil {
		return err
	}

	return nil
}
