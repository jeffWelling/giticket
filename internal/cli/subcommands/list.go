package subcommands

import (
	"flag"
	"fmt"
	"os"

	"github.com/jeffwelling/giticket/pkg/common"
	"github.com/jeffwelling/giticket/pkg/ticket"
)

// init registers the list subcommand
func init() {
	subcommand := new(SubcommandList)
	registerSubcommand("list", subcommand)
}

// SubcommandList implements SubcommandInterface and extends it with
// attributes common to the list subcommand
type SubcommandList struct {
	flagset     *flag.FlagSet
	debugFlag   bool
	helpFlag    bool
	windowWidth int
	parameters  map[string]interface{}
}

// InitFlags sets up the flags for the list subcommand, parses flags, and returns any errors
func (subcommand *SubcommandList) InitFlags(args []string) error {
	subcommand.parameters = make(map[string]interface{})
	var (
		helpFlag  bool
		window    int
		debugFlag bool
	)
	subcommand.flagset = flag.NewFlagSet("list", flag.ExitOnError)

	subcommand.flagset.BoolVar(&subcommand.debugFlag, "debug", false, "Print debug info")
	subcommand.flagset.BoolVar(&subcommand.helpFlag, "help", false, "Print help")
	subcommand.flagset.IntVar(&subcommand.windowWidth, "window", 0, "Window width")
	subcommand.flagset.IntVar(&subcommand.windowWidth, "w", 0, "Window width")
	subcommand.flagset.Parse(args)

	subcommand.parameters["debugFlag"] = debugFlag
	subcommand.parameters["helpFlag"] = helpFlag
	subcommand.parameters["windowWidth"] = window
	return nil
}

// Execute is used to list tickets when the user uses the list subcommand from the CLI
func (subcommand *SubcommandList) Execute() {
	ticket.HandleList(subcommand.debugFlag, common.BranchName, subcommand.windowWidth, os.Stdout)
}

// Help prints help information for the list subcommand
func (subcommand *SubcommandList) Help() {
	fmt.Println("  list - List tickets")
	fmt.Println("    eg: giticket list [params]")
}

type listParams struct {
	windowLength int
	debugFlag    bool
}

// Parameters
func (subcommand *SubcommandList) Parameters() map[string]interface{} {
	return subcommand.parameters
}

// DebugFlag
func (subcommand *SubcommandList) DebugFlag() bool {
	return subcommand.debugFlag
}
