package subcommands

import (
	"flag"
	"fmt"

	"github.com/jeffwelling/giticket/pkg/common"
	"github.com/jeffwelling/giticket/pkg/ticket"
)

// init registers the status subcommand
func init() {
	subcommand := new(SubcommandStatus)
	registerSubcommand("status", subcommand)
}

// SubcommandStatus implements SubcommandInterface and extends it with
// attributes common to the status subcommand
type SubcommandStatus struct {
	flagset    *flag.FlagSet
	debugFlag  bool
	helpFlag   bool
	ticketID   int
	status     string
	parameters map[string]interface{}
}

// InitFlags sets up the flags for the status subcommand, parses flags, and returns any errors
func (subcommand *SubcommandStatus) InitFlags(args []string) error {
	subcommand.parameters = make(map[string]interface{})
	var (
		helpFlag  bool
		debugFlag bool
		ticketID  int
		status    string
	)

	subcommand.flagset = flag.NewFlagSet("status", flag.ExitOnError)

	subcommand.flagset.BoolVar(&subcommand.debugFlag, "debug", false, "Print debug info")
	subcommand.flagset.BoolVar(&subcommand.helpFlag, "help", false, "Print help for the status subcommand")
	subcommand.flagset.StringVar(&subcommand.status, "status", "", "Status to set the ticket to")
	subcommand.flagset.StringVar(&subcommand.status, "s", "", "Status to set the ticket to")
	subcommand.flagset.IntVar(&subcommand.ticketID, "ticketid", 0, "Ticket ID")
	subcommand.flagset.IntVar(&subcommand.ticketID, "id", 0, "Ticket ID")
	subcommand.flagset.Parse(args)

	subcommand.parameters["debugFlag"] = debugFlag
	subcommand.parameters["helpFlag"] = helpFlag
	subcommand.parameters["status"] = status
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
	if subcommand.status == "" {
		fmt.Println("Error: Status must be specified")
		// Print usage
		common.PrintGeneralUsage()
		subcommand.Help()
	}
	return nil
}

// Execute is used to set the status of a ticket when the status subcommand is used from the CLI
func (subcommand *SubcommandStatus) Execute() {
	ticket.HandleStatus(subcommand.status, subcommand.ticketID, subcommand.helpFlag, subcommand.debugFlag)
}

// Help prints help information for the status subcommand
func (subcommand *SubcommandStatus) Help() {
	fmt.Println("  status - Set ticket status")
	fmt.Println("    eg. giticket status [parameters]")
	fmt.Println("    parameters:")
	fmt.Println("      --ticketid | --id 1")
	fmt.Println("      --status   | --s new")
	fmt.Println("      --debug")
	fmt.Println("      --help")
	fmt.Println("    examples:")
	fmt.Println("      - name: Set status of ticket with ID #1 to new")
	fmt.Println("        example: giticket status --ticketid 1 --status new")
}

// Parameters
func (subcommand *SubcommandStatus) Parameters() map[string]interface{} {
	return subcommand.parameters
}

// DebugFlag
func (subcommand *SubcommandStatus) DebugFlag() bool {
	return subcommand.debugFlag
}
