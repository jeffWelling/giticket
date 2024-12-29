package ticket

import (
	"fmt"
	"io"
	"reflect"
	"strings"

	git "github.com/jeffwelling/git2go/v37"
	"github.com/jeffwelling/giticket/pkg/common"
	"github.com/jeffwelling/giticket/pkg/debug"
	"gopkg.in/yaml.v2"
)

func HandleList(w io.Writer, windowWidth int, branchName string, filterName string, filterSet bool, debugFlag bool) error {
	debug.DebugMessage(debugFlag, "Opening git repository")
	thisRepo, err := git.OpenRepository(".")
	if err != nil {
		return err
	}

	output, err := ListTickets(
		thisRepo, branchName, windowWidth, filterName, filterSet, debugFlag)
	if err != nil {
		return err
	}
	fmt.Fprint(w, output)
	return nil
}

// GetTicketsList takes no parameters and returns a list of tickets
// It is intended to be used by giticket-webui, GetListOfTickets requires a repo
// parameter, this creates the repo and calls GetListOfTickets
func GetTicketsList() ([]Ticket, error) {
	thisRepo, err := git.OpenRepository(".")
	if err != nil {
		return nil, err
	}
	return GetListOfTickets(thisRepo, common.BranchName, false)
}

func ListTickets(thisRepo *git.Repository, branchName string, windowWidth int, filterName string, filterSet bool, debugFlag bool) (string, error) {
	output := ""

	// Get a list of tickets from the repo
	var ticketsList []Ticket
	ticketsList, err := GetListOfTickets(thisRepo, branchName, debugFlag)
	if err != nil {
		return "", fmt.Errorf("unable to list tickets: %s", err) // TODO: err
	}

	// Sanity check that a filter has been set before attempting to set
	// preferred filter
	currentFilter, err := GetCurrentFilter(debugFlag)
	if err != nil {
		return "", err
	}

	// If the user is trying to set the preferred filter, but the filter name is
	// empty, that's an error.
	if filterName == "" && filterSet {
		return "", fmt.Errorf("cannot set preferred filter when no filter has been configured yet, create one with the filter subcommand")
	}

	// Filter tickets
	filteredTicketsList := new([]Ticket)
	if filterName != "" {
		filteredTicketsList, err = FilterTickets(ticketsList, filterName, debugFlag)
		if err != nil {
			return "", err
		}
	} else if currentFilter != "" {
		filteredTicketsList, err = FilterTickets(ticketsList, currentFilter, debugFlag)
		if err != nil {
			return "", err
		}
	} else {
		filteredTicketsList = &ticketsList
	}

	widthOfID := widest(ticketsList, "ID")
	if widthOfID < 3 {
		widthOfID = 3
	}
	widthOfTitle := widest(ticketsList, "Title")
	if widthOfTitle < 20 {
		widthOfTitle = 20
	}
	widthOfSeverity := widest(ticketsList, "Severity")
	if widthOfSeverity < 9 {
		widthOfSeverity = 9
	}
	widthOfStatus := widest(ticketsList, "Status")
	if widthOfStatus < 10 {
		widthOfStatus = 10
	}

	// Print the header
	output += padRight("ID", widthOfID) + " | " + padRight("Title", widthOfTitle) + " | " + padRight("Severity", widthOfSeverity) + " | " + padRight("Status", widthOfStatus) + "\n"
	output += strings.Repeat("-", widthOfID+widthOfTitle+widthOfSeverity+widthOfStatus+4) + "\n"

	// Print the tickets
	for _, t := range *filteredTicketsList {
		IDAsString := fmt.Sprintf("%d", t.ID)
		SeverityAsString := fmt.Sprintf("%d", t.Severity)
		output += fmt.Sprintf("%s | %s | %s | %s\n", padRight(IDAsString, widthOfID), padRight(t.Title, widthOfTitle), padRight(SeverityAsString, widthOfSeverity), padRight(t.Status, widthOfStatus))
	}

	return output, nil
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

// widest() takes a list of tickets and a string representing the attribute name
// of each ticket to check to find the widest string, and return that value
func widest(tickets []Ticket, attr string) int {
	widest := 0
	for _, ticket := range tickets {
		v := reflect.ValueOf(ticket)
		if v.Kind() != reflect.Struct {
			panic("not a struct")
		}

		fieldVal := v.FieldByName(attr)
		if !fieldVal.IsValid() {
			panic("not a valid field")
		}
		if !fieldVal.CanInterface() {
			panic("cannot interface")
		}

		// if fieldVal is an int, conver it to string
		fieldValString := fmt.Sprintf("%v", fieldVal.Interface())

		if len(fieldValString) > widest {
			widest = len(fieldValString)
		}
	}
	return widest
}

func GetListOfTickets(thisRepo *git.Repository, branchName string, debugFlag bool) ([]Ticket, error) {
	debug.DebugMessage(debugFlag, "GetListOfTickets() start")

	// Find the branch and its target commit
	debug.DebugMessage(debugFlag, "Looking up branch: "+branchName)
	branch, err := thisRepo.LookupBranch(branchName, git.BranchLocal)
	if err != nil {
		return nil, fmt.Errorf("unable to list tickets because there was an error looking up the branch: %s", err)
	}

	// Lookup the commit the branch references
	debug.DebugMessage(debugFlag, "Looking up commit: "+branch.Target().String())
	parentCommit, err := thisRepo.LookupCommit(branch.Target())
	if err != nil {
		return nil, fmt.Errorf("unable to list tickets because there was an error looking up the commit: %s", err)
	}

	debug.DebugMessage(debugFlag, "looking up tree from parent commit, tree ID: "+parentCommit.TreeId().String())
	parentCommitTree, err := parentCommit.Tree()
	if err != nil {
		return nil, fmt.Errorf("unable to list tickets because there was an error looking up the tree of the parent commit: %s", err)
	}
	defer parentCommitTree.Free()

	debug.DebugMessage(debugFlag, "looking up .giticket tree entry from parent commit: "+parentCommitTree.Id().String())
	giticketTreeEntry, err := parentCommitTree.EntryByPath(".giticket")
	if err != nil {
		return nil, fmt.Errorf("unable to list tickets because there was an error looking up the .giticket tree entry: %s", err)
	}

	debug.DebugMessage(debugFlag, "looking up giticketTree from ID: "+giticketTreeEntry.Id.String())
	giticketTree, err := thisRepo.LookupTree(giticketTreeEntry.Id)
	if err != nil {
		return nil, fmt.Errorf("unable to list tickets because there was an error looking up the .giticket tree from the ID: %s", err)
	}
	defer giticketTree.Free()

	debug.DebugMessage(debugFlag, "looking up tickets tree from .giticket tree: "+giticketTreeEntry.Id.String())
	giticketTicketsTreeEntry, err := giticketTree.EntryByPath("tickets")
	if err != nil {
		return nil, fmt.Errorf("unable to list tickets because there was an error looking up the .giticket/tickets tree entry: %s", err)
	}

	debug.DebugMessage(debugFlag, "looking up giticketTicketsTree from ID: "+giticketTicketsTreeEntry.Id.String())
	giticketTicketsTree, err := thisRepo.LookupTree(giticketTicketsTreeEntry.Id)
	if err != nil {
		return nil, fmt.Errorf("unable to list tickets because there was an error looking up the .giticket/tickets tree from the ID: %s", err)
	}
	defer giticketTicketsTree.Free()

	var ticketList []Ticket
	var t Ticket
	err = giticketTicketsTree.Walk(func(name string, entry *git.TreeEntry) error {
		ticketFile, err := thisRepo.LookupBlob(entry.Id)
		if err != nil {
			return fmt.Errorf("error walking the tickets tree and looking up the entry ID: %s", err)
		}
		defer ticketFile.Free()

		t = Ticket{}
		// Unmarshal the ticket which is yaml
		err = yaml.Unmarshal(ticketFile.Contents(), &t)
		if err != nil {
			return fmt.Errorf("error unmarshalling yaml ticket from file in tickets directory: %s", err)
		}

		ticketList = append(ticketList, t)
		return nil
	})
	if err != nil {
		return nil, fmt.Errorf("error walking the tickets tree: %s", err)
	}

	debug.DebugMessage(debugFlag, "Number of tickets: "+fmt.Sprint(len(ticketList)))
	return ticketList, nil
}
