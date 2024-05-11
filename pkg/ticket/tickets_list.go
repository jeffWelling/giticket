package ticket

import (
	"fmt"

	git "github.com/jeffwelling/git2go/v37"
	"gopkg.in/yaml.v2"
)

func GetListOfTickets(debug bool) []Ticket {
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

	var ticketList []Ticket
	var t Ticket
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

		t = Ticket{}
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
