package cli

import (
	"fmt"
	"os"
	"strings"

	"github.com/jeffWelling/giticket/internal/cli/subcommands"
	"github.com/jeffWelling/giticket/pkg/common"
)

func Exec() {
	// Sanity check, are we being called with no subcommands
	if len(os.Args) <= 1 {
		common.PrintGeneralUsage()
		fmt.Print("\n")
		fmt.Println("Available Actions:\n-  " + strings.Join(subcommands.ListSubcommand(), "\n-  "))
		fmt.Print("\n")
		return
	}

	// Parse the subcommand
	subcommand_name := os.Args[1]

	// If no first argument is provided, or if the first argument is a flag
	if subcommand_name == "" || strings.HasPrefix(subcommand_name, "-") {

		if subcommand_name == "--version" {
			common.PrintVersion()
			return
		}

		if subcommand_name == "--help" {
			common.PrintGeneralUsage()
			fmt.Print("\n")
			fmt.Println("Available Actions:\n-  " + strings.Join(subcommands.ListSubcommand(), "\n-  "))
			fmt.Print("\n")
			return
		}

		printActionMissing()
		return
	}

	subcommand := subcommands.Use(subcommand_name)
	if len(os.Args) <= 2 && subcommand_name != "init" && subcommand_name != "list" {
		// Every subcommand except init requires one or more parameters
		subcommand.Help()
		return
	}
	err := subcommand.InitFlags(os.Args[2:])
	if err != nil {
		fmt.Println("There was a problem initializing the flags: ")
		fmt.Println(err)
		return
	}
	subcommand.Execute()
}

func printBanner() {
	fmt.Println("======================================")
}
func printActionMissing() {
	printBanner()
	fmt.Println("Warning: No action given, and no parameters given. Nothing to do.")
	printBanner()
	common.PrintGeneralUsage()
}
