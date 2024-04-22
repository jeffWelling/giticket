package cli

import (
	"flag"
	"fmt"

	"github.com/jeffWelling/giticket/internal/cli/actions"
)

// constants
const version = "0.0.1"

func Exec() {

	// Set flags
	helpFlag := flag.Bool("help", false, "Print usage info and detailed help")
	versionFlag := flag.Bool("version", false, "Print version info")

	// Parse args
	flag.Parse()

	// Parse the action
	action := getAction(flag.Args())
	if action == "" {
		// No action, print general usage if help flag used
		if *helpFlag {
			printGeneralUsage()
			return
		}

		// If version flag set, print version
		if *versionFlag {
			printVersion()
			return
		}

		// No action and no flag set
		printActionMissing()
		return
	}

	// Execute
	actions.Use(action).Execute()
}

func printGeneralUsage() {
	fmt.Println("Giticket is a git based bug tracker written in golang.")
	printVersion()
	fmt.Println("Usgae: giticket {action} [parameters]")
	fmt.Println("One action is accepted")
	fmt.Println("Actions:")
	availableActions := actions.ListActions()
	for _, action := range availableActions {
		actions.Use(action).Help()
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

// getAction takes a list of strings, assumed to be the output of flags.Args()
// It returns args[0] unless args[0] starts with '-' in which case it returns ""
func getAction(args []string) string {
	if len(args) == 0 {
		return ""
	}
	if args[0][0] == '-' {
		return ""
	}
	return args[0]
}
