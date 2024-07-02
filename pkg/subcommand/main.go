// package subcommand implements the interface for subcommands in giticket
package subcommand

// SubcommandInterface defines methods all subcommands must implement
type SubcommandInterface interface {
	Execute()
	Help()
	InitFlags([]string) error
	Parameters() map[string]interface{} // TODO remove me?
	DebugFlag() bool                    // TODO remove me?
}
