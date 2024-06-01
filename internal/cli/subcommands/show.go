package subcommands

import (
	"flag"
	"fmt"

	"github.com/jeffwelling/giticket/pkg/common"
	"github.com/jeffwelling/giticket/pkg/ticket"
)

func init() {
	subcommand := new(SubcommandShow)
	registerSubcommand("show", subcommand)
}

type SubcommandShow struct {
	debugFlag  bool
	helpFlag   bool
	output     string
	ticket_id  int
	flagset    *flag.FlagSet
	parameters map[string]interface{}
}

func (subcommand *SubcommandShow) Execute() {
	ticket.HandleShow(subcommand.ticket_id, subcommand.output, subcommand.debugFlag, subcommand.helpFlag)
}

func (subcommand *SubcommandShow) Help() {
	fmt.Println("  show - Show ticket")
	fmt.Println("    eg: giticket show [parameters]")
	fmt.Println("    parameters:")
	fmt.Println("      --ticketid | --id N")
	fmt.Println("      --output   | --o text|yaml|json")
	fmt.Println("      --debug")
	fmt.Println("      --help")
	fmt.Println("    examples:")
	fmt.Println("      - name: Show ticket with ID #1")
	fmt.Println("        example: giticket show --ticketid 1")

}

func (subcommand *SubcommandShow) InitFlags(args []string) error {
	subcommand.parameters = make(map[string]interface{})
	var (
		helpFlag  bool
		debugFlag bool
		output    string
		ticket_id int
	)
	subcommand.flagset = flag.NewFlagSet("show", flag.ExitOnError)

	subcommand.flagset.BoolVar(&subcommand.debugFlag, "debug", false, "Print debug info")
	subcommand.flagset.BoolVar(&subcommand.helpFlag, "help", false, "Print help for the show subcommand")
	subcommand.flagset.StringVar(&subcommand.output, "output", "text", "Output format")
	subcommand.flagset.StringVar(&subcommand.output, "o", "text", "Output format")
	subcommand.flagset.IntVar(&subcommand.ticket_id, "ticketid", 0, "Ticket ID")
	subcommand.flagset.IntVar(&subcommand.ticket_id, "id", 0, "Ticket ID")
	subcommand.flagset.Parse(args)

	subcommand.parameters["debugFlag"] = debugFlag
	subcommand.parameters["helpFlag"] = helpFlag
	subcommand.parameters["output"] = output
	subcommand.parameters["ticket_id"] = ticket_id

	if subcommand.helpFlag {
		common.PrintVersion()
		fmt.Println("giticket")
		subcommand.Help()
	}

	if subcommand.ticket_id == 0 {
		fmt.Println("Error: Ticket ID must be specified")
		// Print usage
		common.PrintGeneralUsage()
		subcommand.Help()
	}

	return nil
}

// Parameters()
func (subcommand *SubcommandShow) Parameters() map[string]interface{} {
	return subcommand.parameters
}

// DebugFlag()
func (subcommand *SubcommandShow) DebugFlag() bool {
	return subcommand.debugFlag
}
