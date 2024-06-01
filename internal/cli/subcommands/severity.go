package subcommands

import (
	"flag"
	"fmt"

	"github.com/jeffwelling/giticket/pkg/common"
	"github.com/jeffwelling/giticket/pkg/ticket"
)

func init() {
	subcommand := new(SubcommandSeverity)
	registerSubcommand("severity", subcommand)
}

type SubcommandSeverity struct {
	flagset    *flag.FlagSet
	debugFlag  bool
	helpFlag   bool
	ticketID   int
	severity   int
	parameters map[string]interface{}
}

func (subcommand *SubcommandSeverity) InitFlags(args []string) error {
	subcommand.parameters = make(map[string]interface{})

	var (
		helpFlag  bool
		ticketID  int
		severity  int
		debugFlag bool
	)

	subcommand.flagset = flag.NewFlagSet("severity", flag.ExitOnError)

	subcommand.flagset.BoolVar(&subcommand.debugFlag, "debug", false, "Print debug info")
	subcommand.flagset.BoolVar(&subcommand.helpFlag, "help", false, "Print help for the severity subcommand")
	subcommand.flagset.IntVar(&subcommand.severity, "severity", 1, "Severity of the ticket")
	subcommand.flagset.IntVar(&subcommand.severity, "s", 1, "Severity of the ticket")
	subcommand.flagset.IntVar(&subcommand.ticketID, "ticketid", 0, "Ticket ID")
	subcommand.flagset.IntVar(&subcommand.ticketID, "id", 0, "Ticket ID")
	subcommand.flagset.Parse(args)

	subcommand.parameters["debugFlag"] = debugFlag
	subcommand.parameters["helpFlag"] = helpFlag
	subcommand.parameters["severity"] = severity
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

func (subcommand *SubcommandSeverity) Execute() {
	ticket.HandleSeverity(subcommand.ticketID, subcommand.severity, subcommand.debugFlag)
}

func (subcommand *SubcommandSeverity) Help() {
	fmt.Println("  severity - Set severity")
	fmt.Println("    eg: giticket severity [parameters]")
	fmt.Println("    parameters:")
	fmt.Println("      --ticketid | --id 1")
	fmt.Println("      --severity | --s 1")
	fmt.Println("      --debug")
	fmt.Println("      --help")
	fmt.Println("    examples:")
	fmt.Println("      - name: Set severity of ticket with ID #1 to 1")
	fmt.Println("        example: giticket severity --ticketid 1 --severity 1")
}

// Parameters
func (subcommand *SubcommandSeverity) Parameters() map[string]interface{} {
	return subcommand.parameters
}

// DebugFlag
func (subcommand *SubcommandSeverity) DebugFlag() bool {
	return subcommand.debugFlag
}
