package subcommands

import (
	"flag"
	"fmt"

	"github.com/jeffWelling/giticket/pkg/ticket"
)

func init() {
	subcommand := new(SubcommandShow)
	registerSubcommand("show", subcommand)
}

type SubcommandShow struct {
	debug     bool
	output    string
	ticket_id int
	flagset   *flag.FlagSet
}

func (subcommand *SubcommandShow) Execute() {
	tickets := ticket.GetListOfTickets(subcommand.debug)
	t := ticket.FilterTicketsByID(tickets, subcommand.ticket_id)
	ticket.ShowTicket(t, subcommand.output, subcommand.debug)
}

func (subcommand *SubcommandShow) Help() {
	fmt.Println("  show - Show ticket")
	fmt.Println("    eg: giticket show [parameters]")
	fmt.Println("    parameters:")
	fmt.Println("      -id N")
	fmt.Println("      -debug")
	fmt.Println("      -output text|yaml|json")
}

func (subcommand *SubcommandShow) InitFlags(args []string) {
	subcommand.flagset = flag.NewFlagSet("show", flag.ExitOnError)

	subcommand.flagset.BoolVar(&subcommand.debug, "debug", false, "Print debug info")
	subcommand.flagset.StringVar(&subcommand.output, "output", "markdown", "Output format")
	subcommand.flagset.IntVar(&subcommand.ticket_id, "id", 0, "Ticket ID")
	subcommand.flagset.Parse(args)
}
