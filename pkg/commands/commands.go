package commands

import (
	"github.com/blinkops/blink-go-cli/pkg/commands/initialize"
	"github.com/blinkops/blink-go-cli/pkg/commands/invites"
	"github.com/blinkops/blink-go-cli/pkg/commands/playbooks"
	"github.com/spf13/cobra"
)

// GetRegisteredChildCommands
// Specify the parent usage when this command should
// be grouped a part of the generated operation
func GetRegisteredChildCommands() map[string][]*cobra.Command {
	return map[string][]*cobra.Command{
		"playbooks": {
			playbooks.ListPlaybooksCommand(),
			playbooks.CreatePlaybookCommand(),
			playbooks.UpdatePlaybooksCommand(),
		},
		"invites": {
			invites.InviteCommand(),
		},
	}
}

func GetRegisteredStandaloneCommands() []*cobra.Command {
	return []*cobra.Command{
		initialize.Command(),
	}
}
