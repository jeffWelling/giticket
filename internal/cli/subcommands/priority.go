package subcommands

import (
	"flag"
	"fmt"

	"github.com/jeffwelling/giticket/pkg/common"
	"github.com/jeffwelling/giticket/pkg/ticket"
)

// init registers the priority subcommand
func init() {
	subcommand := new(SubcommandPriority)
	registerSubcommand("priority", subcommand)
}

// SubcommandPriority implements SubcommandInterface and extends it with
// attributes common to the priority subcommand
type SubcommandPriority struct {
	flagset    *flag.FlagSet
	debugFlag  bool
	helpFlag   bool
	priority   int
	ticketID   int
	parameters map[string]interface{}
}

// InitFlags sets up the flags for the priority subcommand, parses flags, and returns any errors
func (subcommand *SubcommandPriority) InitFlags(args []string) error {
	subcommand.parameters = make(map[string]interface{})
	var (
		helpFlag  bool
		ticketID  int
		priority  int
		debugFlag bool
	)
	subcommand.flagset = flag.NewFlagSet("priority", flag.ExitOnError)

	subcommand.flagset.BoolVar(&subcommand.debugFlag, "debug", false, "Print debug info")
	subcommand.flagset.BoolVar(&subcommand.helpFlag, "help", false, "Print help for the priority subcommand")
	subcommand.flagset.IntVar(&subcommand.priority, "priority", 1, "Priority of the ticket")
	subcommand.flagset.IntVar(&subcommand.priority, "p", 1, "Priority of the ticket")
	subcommand.flagset.IntVar(&subcommand.ticketID, "ticketid", 0, "Ticket ID")
	subcommand.flagset.IntVar(&subcommand.ticketID, "id", 0, "Ticket ID")
	if err := subcommand.flagset.Parse(args); err != nil {
		return err
	}

	subcommand.parameters["debugFlag"] = debugFlag
	subcommand.parameters["helpFlag"] = helpFlag
	subcommand.parameters["priority"] = priority
	subcommand.parameters["ticketID"] = ticketID

	if subcommand.helpFlag {
		common.PrintVersion()
		fmt.Println("giticket")
		subcommand.Help()
	}

	if subcommand.ticketID == 0 {
		fmt.Println("Error: Ticket ID must be specified")
		// Print usage
		common.PrintGeneralUsage()
		subcommand.Help()
	}
	return nil
}

// Execute executes the priority subcommand when the priority subcommand is used
// from the CLI
func (subcommand *SubcommandPriority) Execute() {
	err := ticket.HandlePriority(subcommand.ticketID, subcommand.priority, subcommand.debugFlag)
	if err != nil {
		fmt.Println(err)
		return
	}
}

// Help prints help information for the priority subcommand
func (subcommand *SubcommandPriority) Help() {
	fmt.Println("  priority - Set priority")
	fmt.Println("    eg: giticket priority [parameters]")
	fmt.Println("    parameters:")
	fmt.Println("      --ticketid | --id 1")
	fmt.Println("      --priority | --p 1")
	fmt.Println("      --debug")
	fmt.Println("      --help")
	fmt.Println("    examples:")
	fmt.Println("      - name: Set priority of ticket with ID #1 to 1")
	fmt.Println("        example: giticket priority --ticketid 1 --priority 1")
}

// Parameters
func (subcommand *SubcommandPriority) Parameters() map[string]interface{} {
	return subcommand.parameters
}

// DebugFlag()
func (subcommand *SubcommandPriority) DebugFlag() bool {
	return subcommand.debugFlag
}
