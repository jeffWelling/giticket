package actions

var registryActions map[string]ActionInterface

type ActionInterface interface {
	Execute()
	Help()
}

// registerAction() takes an action_name and an action_plugin
// and register the action_plugin under action_name. Subsequent calls
// with the same action_name will overwrite the previous registration
func registerAction(action_name string, action_plugin ActionInterface) {
	if len(registryActions) == 0 {
		registryActions = make(map[string]ActionInterface)
	}
	registryActions[action_name] = action_plugin
}

// isAction() takes an action_name and checks that it matches an action
// that's in the registry
func isAction(action_name string) bool {
	_, ok := registryActions[action_name]
	return ok
}

// Use the action with the given name by returning it
func Use(action_name string) ActionInterface {
	return registryActions[action_name]
}

// ListActions returns a list of strings which are the names
// of the available actions.
func ListActions() []string {
	keys := make([]string, 0, len(registryActions))
	for k := range registryActions {
		keys = append(keys, k)
	}
	return keys
}
