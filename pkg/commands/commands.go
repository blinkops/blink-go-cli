package commands

import (
	"github.com/blinkops/blink-go-cli/pkg/commands/initialize"
	"github.com/blinkops/blink-go-cli/pkg/commands/invites"
	"github.com/blinkops/blink-go-cli/pkg/commands/playbooks"
	"github.com/blinkops/blink-go-cli/pkg/commands/runner-groups"
	"github.com/blinkops/blink-go-cli/pkg/commands/table"
	"github.com/blinkops/blink-go-cli/pkg/consts"
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
			playbooks.GetPlaybookCommand(),
			playbooks.DeletePlaybookCommand(),
		},
		"invites": {
			invites.InviteCommand(),
		},
		"runners": {
			runner_groups.CreateRunnerGroupCommand(),
		},
	}
}

func GetRegisteredStandaloneCommands() []*cobra.Command {
	commands := []*cobra.Command{
		initialize.Command(),
		table.CRUD("connection_types", "connection-types", "Manage your connection types"),
	}

	for _, command := range commands {
		consts.AddAllowedOperation(command.Name())
	}

	return commands
}
