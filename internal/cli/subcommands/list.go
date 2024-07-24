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
	debugFlag   bool
	flagset     *flag.FlagSet
	filter      string
	filterSet   bool
	helpFlag    bool
	parameters  map[string]interface{}
	windowWidth int
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
	subcommand.flagset.StringVar(&subcommand.filter, "filter", "", "The filter name to use for listing tickets with")
	subcommand.flagset.StringVar(&subcommand.filter, "f", "", "The filter name to use for listing tickets with")
	subcommand.flagset.BoolVar(&subcommand.filterSet, "set-filter", false, "Requires the filter name parameter. If true, save the name of the filter as the default filter to use for future list operations.")
	if err := subcommand.flagset.Parse(args); err != nil {
		return err
	}

	// If filterSet is true, then filter is required
	if subcommand.filterSet && subcommand.filter == "" {
		return fmt.Errorf("Filter name is required when using the --set-filter flag")
	}

	subcommand.parameters["debugFlag"] = debugFlag
	subcommand.parameters["helpFlag"] = helpFlag
	subcommand.parameters["windowWidth"] = window
	return nil
}

// Execute is used to list tickets when the user uses the list subcommand from the CLI
func (subcommand *SubcommandList) Execute() {
	err := ticket.HandleList(os.Stdout, subcommand.windowWidth, common.BranchName, subcommand.filter, subcommand.filterSet, subcommand.debugFlag)
	if err != nil {
		fmt.Println(err)
		return
	}
}

// Help prints help information for the list subcommand
func (subcommand *SubcommandList) Help() {
	fmt.Println("  list - List tickets")
	fmt.Println("    eg: giticket list [params]")
}

// Parameters
func (subcommand *SubcommandList) Parameters() map[string]interface{} {
	return subcommand.parameters
}

// DebugFlag
func (subcommand *SubcommandList) DebugFlag() bool {
	return subcommand.debugFlag
}
