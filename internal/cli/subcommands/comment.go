package subcommands

import (
	"flag"
	"fmt"

	"github.com/jeffwelling/giticket/pkg/common"
	"github.com/jeffwelling/giticket/pkg/ticket"
)

// Register the comment subcommand
func init() {
	subcommand := new(SubcommandComment)
	registerSubcommand("comment", subcommand)
}

// A SubcommandComment implements SubcommandInterface, extending it to
// include attributes specific to the comment subcommand
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

// InitFlags sets up flags the command subcommand, parses flags, and returns any
// errors.
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

// Execute is used to add a comment when the comment subcommand is used from the
// CLI
func (subcommand *SubcommandComment) Execute() {
	ticket.HandleComment(
		common.BranchName,
		subcommand.comment,
		subcommand.commentID,
		subcommand.ticketID,
		subcommand.delete,
		subcommand.debug,
	)
}

// Help prints information for the comment subcommand, it is called from CLI
// Exec() if the user uses the --help flag with the comment subcommand
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
