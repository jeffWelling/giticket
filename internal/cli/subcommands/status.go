package subcommands

import (
	"flag"
	"fmt"
	"strconv"

	"github.com/jeffWelling/giticket/pkg/common"
	"github.com/jeffWelling/giticket/pkg/debug"
	"github.com/jeffWelling/giticket/pkg/repo"
	"github.com/jeffWelling/giticket/pkg/ticket"
	git "github.com/jeffwelling/git2go/v37"
)

func init() {
	subcommand := new(SubcommandStatus)
	registerSubcommand("status", subcommand)
}

type SubcommandStatus struct {
	flagset   *flag.FlagSet
	debugFlag bool
	helpFlag  bool
	ticketID  int
	status    string
}

func (subcommand *SubcommandStatus) InitFlags(args []string) error {
	subcommand.flagset = flag.NewFlagSet("status", flag.ExitOnError)

	subcommand.flagset.BoolVar(&subcommand.debugFlag, "debug", false, "Print debug info")
	subcommand.flagset.BoolVar(&subcommand.helpFlag, "help", false, "Print help for the status subcommand")
	subcommand.flagset.StringVar(&subcommand.status, "status", "", "Status to set the ticket to")
	subcommand.flagset.StringVar(&subcommand.status, "s", "", "Status to set the ticket to")
	subcommand.flagset.IntVar(&subcommand.ticketID, "ticketid", 0, "Ticket ID")
	subcommand.flagset.IntVar(&subcommand.ticketID, "id", 0, "Ticket ID")
	subcommand.flagset.Parse(args)

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
	if subcommand.helpFlag {
		return
	}
	branchName := "giticket"

	debug.DebugMessage(subcommand.debugFlag, "Opening git repository")
	thisRepo, err := git.OpenRepository(".")
	if err != nil {
		panic(err)
	}

	// Get author
	author := common.GetAuthor(thisRepo)

	tickets, err := ticket.GetListOfTickets(thisRepo, branchName, subcommand.debugFlag)
	if err != nil {
		panic(err)
	}
	t := ticket.FilterTicketsByID(tickets, subcommand.ticketID)

	t.Status = subcommand.status
	repo.Commit(&t, thisRepo, branchName, author, "Setting status of ticket "+strconv.Itoa(t.ID)+" to "+subcommand.status, subcommand.debugFlag)
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