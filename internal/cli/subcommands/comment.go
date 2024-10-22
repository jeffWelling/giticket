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
	debugFlag bool
	params    map[string]interface{}
}

// InitFlags sets up flags the command subcommand, parses flags, and returns any
// errors.
func (subcommand *SubcommandComment) InitFlags(args []string) error {

	subcommand.flagset = flag.NewFlagSet("comment", flag.ExitOnError)

	subcommand.flagset.BoolVar(&subcommand.debugFlag, "debug", false, "Print debug info")

	subcommand.flagset.IntVar(&subcommand.ticketID, "ticketid", 0, "Ticket ID")
	subcommand.flagset.IntVar(&subcommand.ticketID, "id", 0, "Ticket ID")
	subcommand.flagset.IntVar(&subcommand.commentID, "commentid", 0, "Comment ID")
	subcommand.flagset.IntVar(&subcommand.commentID, "cid", 0, "Comment ID")
	subcommand.flagset.StringVar(&subcommand.comment, "comment", "", "Comment")
	subcommand.flagset.StringVar(&subcommand.comment, "c", "", "Comment")
	subcommand.flagset.BoolVar(&subcommand.delete, "d", false, "Delete comment")
	subcommand.flagset.BoolVar(&subcommand.delete, "delete", false, "Delete comment")
	subcommand.flagset.BoolVar(&subcommand.helpFlag, "help", false, "Print help for the comment subcommand")
	if err := subcommand.flagset.Parse(args); err != nil {
		return err
	}

	subcommand.params = make(map[string]interface{})
	subcommand.params["helpFlag"] = subcommand.helpFlag
	subcommand.params["ticketID"] = subcommand.ticketID
	subcommand.params["commentID"] = subcommand.commentID
	subcommand.params["comment"] = subcommand.comment
	subcommand.params["delete"] = subcommand.delete
	subcommand.params["debugFlag"] = subcommand.debugFlag

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
	_, err := ticket.HandleComment(
		common.BranchName,
		subcommand.comment,
		subcommand.commentID,
		subcommand.ticketID,
		subcommand.delete,
		subcommand.debugFlag,
	)

	if err != nil {
		fmt.Println(err)
		return
	}
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
	return subcommand.debugFlag
}
