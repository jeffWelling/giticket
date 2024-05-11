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
	debug   bool
	output  string
	flagset *flag.FlagSet
}

func (subcommand *SubcommandShow) Execute() {
	tickets := ticket.GetListOfTickets(subcommand.debug)
	fmt.Println(tickets)
}

func (subcommand *SubcommandShow) Help() {
	fmt.Println("  show - Show ticket")
	fmt.Println("    eg: giticket show [parameters]")
	fmt.Println("    parameters:")
	fmt.Println("      -debug")
	fmt.Println("      -output text|yaml|json")
}

func (subcommand *SubcommandShow) InitFlags(args []string) {
	subcommand.flagset = flag.NewFlagSet("show", flag.ExitOnError)

	subcommand.flagset.BoolVar(&subcommand.debug, "debug", false, "Print debug info")
	subcommand.flagset.StringVar(&subcommand.output, "output", "markdown", "Output format")
	subcommand.flagset.Parse(args)
}
