package ticket

import (
	"fmt"

	git "github.com/jeffwelling/git2go/v37"
	"github.com/jeffwelling/giticket/pkg/debug"
	"gopkg.in/yaml.v2"
)

func GetListOfTickets(thisRepo *git.Repository, branchName string, debugFlag bool) ([]Ticket, error) {
	// Find the branch and its target commit
	debug.DebugMessage(debugFlag, "Looking up branch: "+branchName)
	branch, err := thisRepo.LookupBranch(branchName, git.BranchLocal)
	if err != nil {
		return nil, fmt.Errorf("Unable to list tickets because there was an error looking up the branch: %s", err)
	}

	// Lookup the commit the branch references
	debug.DebugMessage(debugFlag, "Looking up commit: "+branch.Target().String())
	parentCommit, err := thisRepo.LookupCommit(branch.Target())
	if err != nil {
		return nil, fmt.Errorf("Unable to list tickets because there was an error looking up the commit: %s", err)
	}

	debug.DebugMessage(debugFlag, "looking up tree from parent commit, tree ID: "+parentCommit.TreeId().String())
	parentCommitTree, err := parentCommit.Tree()
	if err != nil {
		return nil, fmt.Errorf("Unable to list tickets because there was an error looking up the tree of the parent commit: %s", err)
	}
	defer parentCommitTree.Free()

	debug.DebugMessage(debugFlag, "looking up .giticket tree entry from parent commit: "+parentCommitTree.Id().String())
	giticketTreeEntry, err := parentCommitTree.EntryByPath(".giticket")
	if err != nil {
		return nil, fmt.Errorf("Unable to list tickets because there was an error looking up the .giticket tree entry: %s", err)
	}

	debug.DebugMessage(debugFlag, "looking up giticketTree from ID: "+giticketTreeEntry.Id.String())
	giticketTree, err := thisRepo.LookupTree(giticketTreeEntry.Id)
	if err != nil {
		return nil, fmt.Errorf("Unable to list tickets because there was an error looking up the .giticket tree from the ID: %s", err)
	}
	defer giticketTree.Free()

	debug.DebugMessage(debugFlag, "looking up tickets tree from .giticket tree: "+giticketTreeEntry.Id.String())
	giticketTicketsTreeEntry, err := giticketTree.EntryByPath("tickets")
	if err != nil {
		return nil, fmt.Errorf("Unable to list tickets because there was an error looking up the .giticket/tickets tree entry: %s", err)
	}

	debug.DebugMessage(debugFlag, "looking up giticketTicketsTree from ID: "+giticketTicketsTreeEntry.Id.String())
	giticketTicketsTree, err := thisRepo.LookupTree(giticketTicketsTreeEntry.Id)
	if err != nil {
		return nil, fmt.Errorf("Unable to list tickets because there was an error looking up the .giticket/tickets tree from the ID: %s", err)
	}
	defer giticketTicketsTree.Free()

	var ticketList []Ticket
	var t Ticket
	giticketTicketsTree.Walk(func(name string, entry *git.TreeEntry) error {
		ticketFile, err := thisRepo.LookupBlob(entry.Id)
		if err != nil {
			return fmt.Errorf("Error walking the tickets tree and looking up the entry ID: %s", err)
		}
		defer ticketFile.Free()

		t = Ticket{}
		// Unmarshal the ticket which is yaml
		err = yaml.Unmarshal(ticketFile.Contents(), &t)
		if err != nil {
			return fmt.Errorf("Error unmarshalling yaml ticket from file in tickets directory: %s", err)
		}

		ticketList = append(ticketList, t)
		return nil
	})

	debug.DebugMessage(debugFlag, "Number of tickets: "+fmt.Sprint(len(ticketList)))
	return ticketList, nil
}
