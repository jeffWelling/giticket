package subcommands

import (
	"flag"
	"fmt"

	git "github.com/jeffwelling/git2go/v37"
	"github.com/jeffwelling/giticket/pkg/common"
	"github.com/jeffwelling/giticket/pkg/debug"
	"github.com/jeffwelling/giticket/pkg/ticket"
)

func init() {
	subcommand := new(SubcommandShow)
	registerSubcommand("show", subcommand)
}

type SubcommandShow struct {
	debugFlag bool
	helpFlag  bool
	output    string
	ticket_id int
	flagset   *flag.FlagSet
}

func (subcommand *SubcommandShow) Execute() {
	branchName := "giticket"

	debug.DebugMessage(subcommand.debugFlag, "Opening git repository")
	thisRepo, err := git.OpenRepository(".")
	if err != nil {
		panic(err)
	}

	if subcommand.helpFlag {
		return
	}

	tickets, err := ticket.GetListOfTickets(thisRepo, branchName, subcommand.debugFlag)
	if err != nil {
		panic(err)
	}
	t := ticket.FilterTicketsByID(tickets, subcommand.ticket_id)
	ticket.ShowTicket(t, subcommand.output, subcommand.debugFlag)
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
	subcommand.flagset = flag.NewFlagSet("show", flag.ExitOnError)

	subcommand.flagset.BoolVar(&subcommand.debugFlag, "debug", false, "Print debug info")
	subcommand.flagset.BoolVar(&subcommand.helpFlag, "help", false, "Print help for the show subcommand")
	subcommand.flagset.StringVar(&subcommand.output, "output", "text", "Output format")
	subcommand.flagset.StringVar(&subcommand.output, "o", "text", "Output format")
	subcommand.flagset.IntVar(&subcommand.ticket_id, "ticketid", 0, "Ticket ID")
	subcommand.flagset.IntVar(&subcommand.ticket_id, "id", 0, "Ticket ID")
	subcommand.flagset.Parse(args)

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
