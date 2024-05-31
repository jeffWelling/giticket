package subcommands

import (
	"flag"
	"fmt"

	git "github.com/jeffwelling/git2go/v37"
	"github.com/jeffwelling/giticket/pkg/common"
	"github.com/jeffwelling/giticket/pkg/debug"
	"github.com/jeffwelling/giticket/pkg/repo"
	"github.com/jeffwelling/giticket/pkg/ticket"
)

func init() {
	subcommand := new(SubcommandComment)
	registerSubcommand("comment", subcommand)
}

type SubcommandComment struct {
	flagset   *flag.FlagSet
	helpFlag  bool
	ticketID  int
	commentID int
	comment   string
	delete    bool
	debug     bool
	params    map[string]interface{}
}

func (subcommand *SubcommandComment) InitFlags(args []string) error {
	subcommand.params = make(map[string]interface{})
	var (
		helpFlag  bool
		ticketID  int
		commentID int
		comment   string
		delete    bool
		debugFlag bool
	)

	subcommand.flagset = flag.NewFlagSet("comment", flag.ExitOnError)

	subcommand.flagset.IntVar(&ticketID, "ticketid", 0, "Ticket ID")
	subcommand.flagset.IntVar(&ticketID, "id", 0, "Ticket ID")
	subcommand.flagset.IntVar(&commentID, "commentid", 0, "Comment ID")
	subcommand.flagset.IntVar(&commentID, "cid", 0, "Comment ID")
	subcommand.flagset.StringVar(&comment, "comment", "", "Comment")
	subcommand.flagset.StringVar(&comment, "c", "", "Comment")
	subcommand.flagset.BoolVar(&delete, "d", false, "Delete comment")
	subcommand.flagset.BoolVar(&delete, "delete", false, "Delete comment")
	subcommand.flagset.BoolVar(&debugFlag, "debug", false, "Print debug info")
	subcommand.flagset.BoolVar(&helpFlag, "help", false, "Print help for the comment subcommand")
	subcommand.flagset.Parse(args)

	subcommand.params["helpFlag"] = helpFlag
	subcommand.params["ticketID"] = ticketID
	subcommand.params["commentID"] = commentID
	subcommand.params["comment"] = comment
	subcommand.params["delete"] = delete
	subcommand.params["debug"] = debugFlag

	if subcommand.helpFlag {
		common.PrintVersion()
		fmt.Println("giticket")
		subcommand.Help()
	}

	// Sanity check of args
	if subcommand.delete && subcommand.commentID == 0 {
		fmt.Println("Error: When deleting a comment, the commment ID must be specified")
		// Print usage
		common.PrintGeneralUsage()
		subcommand.Help()
	}
	return nil
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

	tickets, err := ticket.GetListOfTickets(thisRepo, branchName, subcommand.debug)
	if err != nil {
		panic(err)
	}
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
	fmt.Println("  comment - Add or remove a comment from a ticket")
	fmt.Println("    eg: giticket comment [parameters]")
	fmt.Println("    parameters:")
	fmt.Println("      --ticketid  | --id N")
	fmt.Println("      --comment   | --c \"My comment\"")
	fmt.Println("      --commentid | --cid N")
	fmt.Println("      --delete    | -d")
	fmt.Println("      --debug")
	fmt.Println("    examples:")
	fmt.Println("      - name: Add a comment to ticket with ID #1")
	fmt.Println("        example: giticket comment --ticketid 1 --comment \"My new comment on ticket with ID #1\"")
	fmt.Println("      - name: Delete a comment with comment ID #1 on ticket with ID #1")
	fmt.Println("        example: giticket comment --ticketid 1 --commentid 1 --delete")
	fmt.Println("      - name: Add a multi-line comment to ticket with ID #1")
	fmt.Println("        example: giticket comment --ticketid 1 --comment \"This is a multi-line comment. \n        This is a new line in the same comment\"")
}

func (subcommand *SubcommandComment) Parameters() map[string]interface{} {
	return subcommand.params
}

func (subcommand *SubcommandComment) DebugFlag() bool {
	return subcommand.debug
}
