package subcommands

import (
	"sort"

	"github.com/jeffwelling/giticket/pkg/subcommand"
)

var registrySubcommands map[string]subcommand.SubcommandInterface

// registerSubcommand() takes a subcommand_name and a subcommand_plugin
// and register the action_plugin under subcommand_name. Subsequent calls
// with the same subcommand_name will overwrite the previous registration
func registerSubcommand(subcommand_name string, action_plugin subcommand.SubcommandInterface) {
	if len(registrySubcommands) == 0 {
		registrySubcommands = make(map[string]subcommand.SubcommandInterface)
	}
	registrySubcommands[subcommand_name] = action_plugin
}

// isAction() takes an subcommand_name and checks that it matches an action
// that's in the registry
func isAction(subcommand_name string) bool {
	_, ok := registrySubcommands[subcommand_name]
	return ok
}

// Use the action with the given name by returning it
func Use(subcommand_name string) subcommand.SubcommandInterface {
	return registrySubcommands[subcommand_name]
}

// ListSubcommand returns a list of strings which are the names
// of the available actions.
func ListSubcommand() []string {
	keys := make([]string, 0, len(registrySubcommands))
	for k := range registrySubcommands {
		keys = append(keys, k)
	}
	// Sort keys alphabetically and return the sorted value
	sort.Strings(keys)
	return keys
}
