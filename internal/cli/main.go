package cli

import (
	"fmt"
	"os"
	"strings"

	"github.com/jeffWelling/giticket/internal/cli/subcommands"
)

// constants
const version = "0.0.3"

func Exec() {
	// Sanity check, are we being called with no subcommands
	if len(os.Args) <= 1 {
		printGeneralUsage()
		return
	}

	// Parse the subcommand
	subcommand_name := os.Args[1]

	// If no first argument is provided, or if the first argument is a flag
	if subcommand_name == "" || strings.HasPrefix(subcommand_name, "-") {

		if subcommand_name == "--version" {
			fmt.Println("giticket version: " + version)
			return
		}

		printActionMissing()
		return
	}

	subcommand := subcommands.Use(subcommand_name)
	subcommand.InitFlags(os.Args[2:])
	subcommand.Execute()
}

func printGeneralUsage() {
	fmt.Println("Giticket is a git based bug tracker written in golang.")
	printVersion()
	fmt.Println("Usgae: giticket {action} [parameters]")
	fmt.Println("One action is accepted")
	fmt.Println("Actions:")
	availableActions := subcommands.ListSubcommand()
	for _, action := range availableActions {
		subcommands.Use(action).Help()
	}
	fmt.Println("Zero or more parameters are accepted, parameters include: -help, -version")
	fmt.Println("giticket -help            will print this message")
	fmt.Println("giticket {action} -help   will print the help for that command")
	fmt.Println("giticket -version         will print the version of giticket")
}
func printBanner() {
	fmt.Println("======================================")
}
func printActionMissing() {
	printBanner()
	fmt.Println("Warning: No action given, and no parameters given. Nothing to do.")
	printBanner()
	printGeneralUsage()
}

func printVersion() {
	fmt.Println("Giticket Version: " + version)
}
