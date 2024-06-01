package subcommands

import (
	"flag"
	"fmt"

	"github.com/jeffwelling/giticket/pkg/common"
	"github.com/jeffwelling/giticket/pkg/ticket"
)

func init() {
	subcommand := new(SubcommandStatus)
	registerSubcommand("status", subcommand)
}

type SubcommandStatus struct {
	flagset    *flag.FlagSet
	debugFlag  bool
	helpFlag   bool
	ticketID   int
	status     string
	parameters map[string]interface{}
}

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

func (subcommand *SubcommandStatus) Execute() {
	ticket.HandleStatus(subcommand.status, subcommand.ticketID, subcommand.helpFlag, subcommand.debugFlag)
}

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
