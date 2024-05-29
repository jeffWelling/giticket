package subcommands

import (
	"flag"
	"fmt"
	"strconv"

	git "github.com/jeffwelling/git2go/v37"
	"github.com/jeffwelling/giticket/pkg/common"
	"github.com/jeffwelling/giticket/pkg/debug"
	"github.com/jeffwelling/giticket/pkg/repo"
	"github.com/jeffwelling/giticket/pkg/ticket"
)

func init() {
	subcommand := new(SubcommandSeverity)
	registerSubcommand("severity", subcommand)
}

type SubcommandSeverity struct {
	debugFlag bool
	helpFlag  bool
	ticketID  int
	severity  int
	flagset   *flag.FlagSet
}

func (subcommand *SubcommandSeverity) InitFlags(args []string) error {
	subcommand.flagset = flag.NewFlagSet("severity", flag.ExitOnError)

	subcommand.flagset.BoolVar(&subcommand.debugFlag, "debug", false, "Print debug info")
	subcommand.flagset.BoolVar(&subcommand.helpFlag, "help", false, "Print help for the severity subcommand")
	subcommand.flagset.IntVar(&subcommand.severity, "severity", 1, "Severity of the ticket")
	subcommand.flagset.IntVar(&subcommand.severity, "s", 1, "Severity of the ticket")
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

func (subcommand *SubcommandSeverity) Execute() {
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
	t.Severity = subcommand.severity

	repo.Commit(&t, thisRepo, branchName, author, "Setting severity of ticket "+strconv.Itoa(t.ID)+" to "+strconv.Itoa(subcommand.severity), subcommand.debugFlag)
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
