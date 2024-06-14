package ticket

import (
	"strconv"

	git "github.com/jeffwelling/git2go/v37"
	"github.com/jeffwelling/giticket/pkg/common"
	"github.com/jeffwelling/giticket/pkg/debug"
	"github.com/jeffwelling/giticket/pkg/repo"
)

// Delete the comment from the tree for the ticket with the given ID
func HandleDelete(ticketID int, branchName string, debugFlag bool) {
	debug.DebugMessage(debugFlag, "Deleting ticket "+strconv.Itoa(ticketID))

}

func deleteTicket(ticketID int, branchName string, debugFlag bool) (bool, error) {
	debug.DebugMessage(debugFlag, "Opening git repository")
	thisRepo, err := git.OpenRepository(".")
	if err != nil {
		return false, err
	}

	debug.DebugMessage(debugFlag, "Getting parent commit from branch '"+branchName+"'")
	parentCommit, err := repo.GetParentCommit(thisRepo, branchName, debugFlag)
	if err != nil {
		return false, err
	}

	debug.DebugMessage(debugFlag, "Getting rootTreeBuilder and previousCommitTree from parent commit")
	rootTreeBuilder, previousCommitTree, err := repo.TreeBuilderFromCommit(parentCommit, thisRepo, debugFlag)
	if err != nil {
		return false, err
	}
	defer rootTreeBuilder.Free()

	debug.DebugMessage(debugFlag, "Getting .giticket subtree from previous commit")
	giticketTree, err := repo.GetSubTreeByName(previousCommitTree, thisRepo, ".giticket", debugFlag)
	if err != nil {
		return false, err
	}
	defer giticketTree.Free()

	debug.DebugMessage(debugFlag, "Getting tickets subtree from giticket tree")
	giticketTreeBuilder, err := thisRepo.TreeBuilderFromTree(giticketTree)
	if err != nil {
		return false, err
	}
	defer giticketTreeBuilder.Free()

	debug.DebugMessage(debugFlag, "Getting tickets subtree from giticket tree")
	giticketTicketsTree, err := repo.GetSubTreeByName(giticketTree, thisRepo, "tickets", debugFlag)
	if err != nil {
		return false, err
	}
	defer giticketTicketsTree.Free()

	debug.DebugMessage(debugFlag, "Getting tickets subtree builder from repo")
	giticketTicketsTreeBuilder, err := thisRepo.TreeBuilderFromTree(giticketTicketsTree)
	if err != nil {
		return false, err
	}
	defer giticketTicketsTreeBuilder.Free()

	// Get ticket filename
	ticketsList, err := GetListOfTickets(thisRepo, branchName, debugFlag)
	// Filter for ticketID
	theTicket := FilterTicketsByID(ticketsList, ticketID)

	// Remove ticket from tickets subtree
	debug.DebugMessage(debugFlag, "Removing ticket "+strconv.Itoa(theTicket.ID)+" from tickets subtree")
	err = giticketTicketsTreeBuilder.Remove(theTicket.TicketFilename())
	if err != nil {
		return false, err
	}

	debug.DebugMessage(debugFlag, "Writing giticket tickets subtree")
	newGiticketTicketsTreeID, err := giticketTicketsTreeBuilder.Write()
	if err != nil {
		return false, err
	}

	debug.DebugMessage(debugFlag, "Adding tickets to giticket tree builder")
	err = giticketTreeBuilder.Insert("tickets", newGiticketTicketsTreeID, git.FilemodeTree)
	if err != nil {
		return false, err
	}

	debug.DebugMessage(debugFlag, "Writing giticket tree")
	newGiticketTreeID, err := giticketTreeBuilder.Write()
	if err != nil {
		return false, err
	}

	debug.DebugMessage(debugFlag, "Adding .giticket to root tree builder")
	err = rootTreeBuilder.Insert(".giticket", newGiticketTreeID, git.FilemodeTree)

	debug.DebugMessage(debugFlag, "Writing root tree")
	newRootTreeID, err := rootTreeBuilder.Write()
	if err != nil {
		return false, err
	}

	debug.DebugMessage(debugFlag, "Looking up root tree")
	newRootTree, err := thisRepo.LookupTree(newRootTreeID)
	if err != nil {
		return false, err
	}

	// Get author data by reading .git configs
	debug.DebugMessage(debugFlag, "getting author data")
	author, err := common.GetAuthor(thisRepo)
	if err != nil {
		return false, err
	}

	// Commit
	commitID, err := thisRepo.CreateCommit("refs/heads/giticket", author, author, "Deleting ticket "+theTicket.TicketFilename(), newRootTree, parentCommit)
	debug.DebugMessage(debugFlag, "Commit ID "+commitID.String()+" created for deleting ticket "+theTicket.TicketFilename())
	return true, err
}
