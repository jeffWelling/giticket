package subcommands

import (
	"flag"
	"fmt"

	"github.com/jeffWelling/giticket/pkg/common"
	"github.com/jeffWelling/giticket/pkg/debug"
	"github.com/jeffWelling/giticket/pkg/repo"
	"github.com/jeffWelling/giticket/pkg/ticket"
	git "github.com/jeffwelling/git2go/v37"
)

func init() {
	subcommand := new(SubcommandComment)
	registerSubcommand("comment", subcommand)
}

type SubcommandComment struct {
	flagset   *flag.FlagSet
	ticketID  int
	commentID int
	comment   string
	delete    bool
	debug     bool
}

func (subcommand *SubcommandComment) InitFlags(args []string) {
	subcommand.flagset = flag.NewFlagSet("comment", flag.ExitOnError)

	subcommand.flagset.IntVar(&subcommand.ticketID, "id", 0, "Ticket ID")
	subcommand.flagset.IntVar(&subcommand.commentID, "commentid", 0, "Comment ID")
	subcommand.flagset.StringVar(&subcommand.comment, "comment", "", "Comment")
	subcommand.flagset.BoolVar(&subcommand.delete, "delete", false, "Delete comment")
	subcommand.flagset.BoolVar(&subcommand.debug, "debug", false, "Print debug info")
	subcommand.flagset.Parse(args)

	// Sanity check of args
	if subcommand.delete && subcommand.commentID == 0 {
		fmt.Println("Error: When deleting a comment, the commment ID must be specified")
		// Print usage
		common.PrintGeneralUsage()
		subcommand.Help()
	}
}

func (subcommand *SubcommandComment) Execute() {
	branchName := "giticket"

	debug.DebugMessage(subcommand.debug, "Opening git repository")
	thisRepo, err := git.OpenRepository(".")
	if err != nil {
		panic(err)
	}

	// Get author
	author := common.GetAuthor(thisRepo)

	tickets := ticket.GetListOfTickets(subcommand.debug)
	t := ticket.FilterTicketsByID(tickets, subcommand.ticketID)

	if subcommand.delete {
		commentID := ticket.DeleteComment(&t, subcommand.commentID, thisRepo, branchName, subcommand.debug)
		err := repo.Commit(&t, thisRepo, branchName, author, "Deleting comment "+commentID, subcommand.debug)
		if err != nil {
			panic(err)
		}
	} else {
		commentID := ticket.AddComment(&t, subcommand.comment, thisRepo, branchName, subcommand.debug)
		err := repo.Commit(&t, thisRepo, branchName, author, "Adding comment "+commentID, subcommand.debug)
		if err != nil {
			panic(err)
		}
	}
}

func (subcommand *SubcommandComment) Help() {
	fmt.Println("  comment - Add comment to ticket")
	fmt.Println("    eg: giticket comment [parameters]")
	fmt.Println("    parameters:")
	fmt.Println("      -id N")
	fmt.Println("      -comment \"My comment\"")
	fmt.Println("      -delete")
}
