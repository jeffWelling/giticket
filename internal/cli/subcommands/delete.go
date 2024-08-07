package subcommands

import (
	"flag"
	"fmt"

	"github.com/jeffwelling/giticket/pkg/common"
	"github.com/jeffwelling/giticket/pkg/ticket"
)

// init is used to register the delete subcommand
func init() {
	subcommand := new(SubcommandDelete)
	registerSubcommand("delete", subcommand)
}

// SubcommandDelete implements SubcommandInterface and extends it with
// attributes common to the delete subcommand
type SubcommandDelete struct {
	flagset   *flag.FlagSet
	params    map[string]interface{}
	ticketID  int
	debugFlag bool
}

// InitFlags sets up the flags for the delete subcommand, parses flags, and returns any errors
func (subcommand *SubcommandDelete) InitFlags(args []string) error {
	subcommand.params = make(map[string]interface{})
	var (
		helpFlag  bool
		ticketID  int
		debugFlag bool
	)

	subcommand.flagset = flag.NewFlagSet("delete", flag.ExitOnError)

	subcommand.flagset.BoolVar(&debugFlag, "debug", false, "Print debug info")
	subcommand.flagset.BoolVar(&helpFlag, "help", false, "Print help for the delete subcommand")
	subcommand.flagset.IntVar(&ticketID, "ticketid", 0, "Ticket ID")
	subcommand.flagset.IntVar(&ticketID, "id", 0, "Ticket ID")
	if err := subcommand.flagset.Parse(args); err != nil {
		return err
	}

	subcommand.ticketID = ticketID
	subcommand.debugFlag = debugFlag
	subcommand.params["helpFlag"] = helpFlag
	subcommand.params["ticketID"] = ticketID
	subcommand.params["debugFlag"] = debugFlag

	// Sanity check
	if ticketID == 0 {
		fmt.Println("Error: ticketID is missing but is required to delete a ticket")
		common.PrintGeneralUsage()
		subcommand.Help()
	}

	if helpFlag {
		common.PrintVersion()
		fmt.Println("giticket")
		subcommand.Help()
	}

	return nil
}

// Execute is used to delete a ticket when the user uses the delete subcommand
// from the CLI
func (subcommand *SubcommandDelete) Execute() {
	_, err := ticket.HandleDelete(subcommand.ticketID, common.BranchName, subcommand.debugFlag)
	if err != nil {
		fmt.Println(err)
		return
	}
}

// Help prints help information for the delete subcommand
func (subcommand *SubcommandDelete) Help() {
	fmt.Println("  delete - Delete ticket from the tree")
	fmt.Println("    eg: giticket delete [parameters]")
	fmt.Println("    parameters:")
	fmt.Println("      --ticketid | --id N")
	fmt.Println("      --debug")
	fmt.Println("      --help")
	fmt.Println("    examples:")
	fmt.Println("      - name: Delete ticket with ID #1")
	fmt.Println("        example: giticket delete --ticketid 1")
}

func (subcommand *SubcommandDelete) Parameters() map[string]interface{} {
	return subcommand.params
}

func (subcommand *SubcommandDelete) DebugFlag() bool {
	return subcommand.debugFlag
}
