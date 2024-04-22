package actions

import "fmt"

// init() is used to register this action
func init() {
	action := new(ActionCreate)
	registerAction("create", action)
}

type ActionCreate struct {
}

// Execute() is called when this action is invoked
func (action *ActionCreate) Execute() {
	// TODO
	fmt.Println("Creating ticket")
}

// Help() prints help for this action
func (action *ActionCreate) Help() {
	fmt.Println("  create - Create a new ticket")
	fmt.Println("    eg: giticket create [parameters]")
	fmt.Println("    parameters:")
	fmt.Println("      -title \"Ticket Title\"")
	fmt.Println("      -description \"Ticket Description\"")
	fmt.Println("      -tags \"tag1, tag2, tag3\"")
}
