package subcommands

import (
	"flag"
	"fmt"
	"strings"

	"github.com/jeffWelling/giticket/pkg/ticket"
	git "github.com/jeffwelling/git2go/v37"
	"gopkg.in/yaml.v2"
)

func init() {
	subcommand := new(SubcommandList)
	registerSubcommand("list", subcommand)
}

type SubcommandList struct {
	flagset *flag.FlagSet
	debug   bool
}

func (subcommand *SubcommandList) InitFlags(args []string) {
	subcommand.flagset = flag.NewFlagSet("list", flag.ExitOnError)
}

func (subcommand *SubcommandList) Execute() {
	listParams := listParams{
		windowLength: 0,
		debug:        subcommand.debug,
	}

	output := ListTickets(listParams)
	fmt.Print(output)
}

func (subcommand *SubcommandList) Help() {
	fmt.Println("  list - List tickets")
	fmt.Println("    eg: giticket list [params]")
}

type listParams struct {
	windowLength int
	debug        bool
}

// ListTickets() takes a listParams parameter which contains the optional
// and mandatory parameters for ListTickets(). The only mandatory parameter is
// windowLength which is the length of the window to list tickets in.
func ListTickets(params listParams) string {
	output := ""

	// Get a list of tickets from the repo
	var ticketsList []ticket.Ticket
	ticketsList = getListOfTickets(params.debug)

	// Print the header
	output += padRight("ID", 5) + "| " + padRight("Title", 30) + " | " + padRight("Severity", 10) + " | " + padRight("Status", 10) + "\n"

	// Print the tickets
	for _, t := range ticketsList {
		IDAsString := fmt.Sprintf("%d", t.ID)
		SeverityAsString := fmt.Sprintf("%d", t.Severity)
		output += fmt.Sprintf("%s | %s | %s | %s\n", padRight(IDAsString, 4), padRight(t.Title, 30), padRight(SeverityAsString, 10), padRight(t.Status, 10))
	}

	return output
}

func getListOfTickets(debug bool) []ticket.Ticket {
	branchName := "giticket"

	if debug {
		fmt.Println("Opening repository '.'")
	}
	repo, err := git.OpenRepository(".")
	if err != nil {
		panic(err)
	}

	// Find the branch and its target commit
	if debug {
		fmt.Println("looking up branch: ", branchName)
	}
	branch, err := repo.LookupBranch(branchName, git.BranchLocal)
	if err != nil {
		panic(err)
	}

	// Lookup the commit the branch references
	if debug {
		fmt.Println("looking up commit: ", branch.Target())
	}
	parentCommit, err := repo.LookupCommit(branch.Target())
	if err != nil {
		panic(err)
	}

	if debug {
		fmt.Println("looking up tree from parent commit, tree ID:", parentCommit.TreeId())
	}
	parentCommitTree, err := parentCommit.Tree()
	if err != nil {
		panic(err)
	}
	defer parentCommitTree.Free()

	if debug {
		fmt.Println("looking up .giticket tree entry from parent commit: ", parentCommitTree.Id())
	}
	giticketTreeEntry, err := parentCommitTree.EntryByPath(".giticket")
	if err != nil {
		panic(err)
	}

	if debug {
		fmt.Println("looking up giticketTree from ID: ", giticketTreeEntry.Id)
	}
	giticketTree, err := repo.LookupTree(giticketTreeEntry.Id)
	if err != nil {
		panic(err)
	}
	defer giticketTree.Free()

	if debug {
		fmt.Println("looking up tickets tree from .giticket tree: ", giticketTreeEntry.Id)
	}
	giticketTicketsTreeEntry, err := giticketTree.EntryByPath("tickets")
	if err != nil {
		panic(err)
	}

	if debug {
		fmt.Println("looking up giticketTicketsTree from ID: ", giticketTicketsTreeEntry.Id)
	}
	giticketTicketsTree, err := repo.LookupTree(giticketTicketsTreeEntry.Id)
	if err != nil {
		panic(err)
	}
	defer giticketTicketsTree.Free()

	var ticketList []ticket.Ticket
	var t ticket.Ticket
	giticketTicketsTree.Walk(func(name string, entry *git.TreeEntry) error {

		blob, err := repo.LookupBlob(entry.Id)
		if err != nil {
			panic(err)
		}
		defer blob.Free()

		ticketFile, err := repo.LookupBlob(entry.Id)
		if err != nil {
			panic(err)
		}

		t = ticket.Ticket{}
		// Unmarshal the ticket which is yaml
		err = yaml.Unmarshal(ticketFile.Contents(), &t)
		if err != nil {
			fmt.Println("Error unmarshalling yaml: ", err)
			fmt.Println("Contents: ", string(ticketFile.Contents()))
			panic(err)
		}

		ticketList = append(ticketList, t)
		return nil
	})

	return ticketList
}

// padRight() takes string s and width int, it finds the difference in length
// between len(s) and width and adds that many spaces to the string to ensure
// the returned string is exactly width len long
func padRight(s string, width int) string {
	diff := width - len(s)
	if diff <= 0 {
		return s[0:width]
	}
	return s + strings.Repeat(" ", diff)
}
