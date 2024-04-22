package actions

import "fmt"

// init() is used to register this action
func init() {
	action := new(ActionInit)
	registerAction("init", action)
}

type ActionInit struct {
}

// Execute() is called when this action is invoked
func (action *ActionInit) Execute() {
	// TODO
	fmt.Println("Initializing giticket")
}

// Help() prints help for this action
func (action *ActionInit) Help() {
	fmt.Println("  init - Initialize giticket")
	fmt.Println("    eg: giticket init")
}
