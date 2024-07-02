package common

import (
	"fmt"
)

// PrintGeneralUsage prints the general usage information for giticket
func PrintGeneralUsage() {
	fmt.Println("Giticket is a git based bug tracker written in golang.")
	PrintVersion()
	fmt.Println("Usgae: giticket {action} [parameters]")
	fmt.Println("One action is accepted")
	// fmt.Println("Actions:")
	// availableActions := subcommands.ListSubcommand()
	// for _, action := range availableActions {
	// 	subcommands.Use(action).Help()
	// }
	fmt.Println("Zero or more parameters are accepted, parameters include: -help, -version")
	fmt.Println("giticket -help            will print this message")
	fmt.Println("giticket {action} -help   will print the help for that command")
	fmt.Println("giticket -version         will print the version of giticket")
}

// PrintVersion prints giticket's version
func PrintVersion() {
	fmt.Println("Giticket Version: " + Version)
}
