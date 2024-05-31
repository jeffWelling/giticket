package subcommand

type SubcommandInterface interface {
	Execute()
	Help()
	InitFlags([]string) error
	Parameters() map[string]interface{}
	DebugFlag() bool
}
