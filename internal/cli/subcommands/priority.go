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
	subcommand := new(SubcommandPriority)
	registerSubcommand("priority", subcommand)
}

type SubcommandPriority struct {
	flagset   *flag.FlagSet
	debugFlag bool
	helpFlag  bool
	priority  int
	ticketID  int
}

func (subcommand *SubcommandPriority) InitFlags(args []string) error {
	subcommand.flagset = flag.NewFlagSet("priority", flag.ExitOnError)

	subcommand.flagset.BoolVar(&subcommand.debugFlag, "debug", false, "Print debug info")
	subcommand.flagset.BoolVar(&subcommand.helpFlag, "help", false, "Print help for the priority subcommand")
	subcommand.flagset.IntVar(&subcommand.priority, "priority", 1, "Priority of the ticket")
	subcommand.flagset.IntVar(&subcommand.priority, "p", 1, "Priority of the ticket")
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
	return nil
}

func (subcommand *SubcommandPriority) Execute() {
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
	t.Priority = subcommand.priority
	repo.Commit(&t, thisRepo, branchName, author, "Setting priority of ticket "+strconv.Itoa(t.ID)+" to "+strconv.Itoa(subcommand.priority)+"", subcommand.debugFlag)
}

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
