package repo

import (
	"github.com/jeffWelling/giticket/pkg/debug"
	"github.com/jeffWelling/giticket/pkg/ticket"
	git "github.com/jeffwelling/git2go/v37"
)

// Commit creates a new commit on the branch with the given message, ticket, and
// author.
func Commit(
	t *ticket.Ticket,
	thisRepo *git.Repository,
	branchName string,
	author *git.Signature,
	commitMessage string,
	debugFlag bool,
) error {
	// Get the parent commit of the branch
	debug.DebugMessage(debugFlag, "Getting parent commit from branch '"+branchName+"'")
	parentCommit, err := GetParentCommit(thisRepo, branchName, debugFlag)
	if err != nil {
		return err
	}

	debug.DebugMessage(debugFlag, "Creating tree builder from parent commit")
	rootTreeBuilder, previousCommitTree, err := TreeBuilderFromCommit(parentCommit, thisRepo, debugFlag)
	if err != nil {
		return err
	}
	defer rootTreeBuilder.Free()

	debug.DebugMessage(debugFlag, "Getting .giticket subtree from previous commit")
	giticketTree, err := GetSubTreeByName(previousCommitTree, thisRepo, ".giticket", debugFlag)
	if err != nil {
		return err
	}
	defer giticketTree.Free()
	debug.DebugMessage(debugFlag, "Found .giticket tree: "+giticketTree.Id().String())

	debug.DebugMessage(debugFlag, "Getting tickets subtree from .giticket")
	giticketTicketsTree, err := GetSubTreeByName(giticketTree, thisRepo, "tickets", debugFlag)
	if err != nil {
		return err
	}
	defer giticketTicketsTree.Free()
	debug.DebugMessage(debugFlag, "Found tickets tree: "+giticketTicketsTree.Id().String())

	debug.DebugMessage(debugFlag, "Creating new ticket blob")
	NewTicketOID, err := thisRepo.CreateBlobFromBuffer([]byte(t.TicketToYaml()))
	if err != nil {
		return err
	}
	debug.DebugMessage(debugFlag, "Created new ticket blob: "+NewTicketOID.String())

	debug.DebugMessage(debugFlag, "Getting tree builder for tickets tree")
	giticketTicketsTreeBuilder, err := thisRepo.TreeBuilderFromTree(giticketTicketsTree)
	if err != nil {
		return err
	}

	debug.DebugMessage(debugFlag, "Inserting new ticket ("+NewTicketOID.String()+") into tickets tree")
	err = giticketTicketsTreeBuilder.Insert(t.TicketFilename(), NewTicketOID, git.FilemodeBlob)
	if err != nil {
		return err
	}

	debug.DebugMessage(debugFlag, "Writing tickets tree")
	giticketTicketsTreeID, err := giticketTicketsTreeBuilder.Write()
	if err != nil {
		return err
	}
	debug.DebugMessage(debugFlag, "Wrote tickets tree: "+giticketTicketsTreeID.String())

	debug.DebugMessage(debugFlag, "Creating new .giticket tree")
	giticketTreeBuilder, err := thisRepo.TreeBuilderFromTree(giticketTree)
	if err != nil {
		return err
	}

	debug.DebugMessage(debugFlag, "Inserting tickets tree ("+giticketTicketsTreeID.String()+") into .giticket tree")
	err = giticketTreeBuilder.Insert("tickets", giticketTicketsTreeID, git.FilemodeTree)
	if err != nil {
		return err
	}

	debug.DebugMessage(debugFlag, "Writing .giticket tree")
	giticketTreeID, err := giticketTreeBuilder.Write()
	if err != nil {
		return err
	}
	debug.DebugMessage(debugFlag, "Wrote .giticket tree: "+giticketTreeID.String())

	debug.DebugMessage(debugFlag, "Inserting .giticket tree ("+giticketTreeID.String()+") into root tree")
	err = rootTreeBuilder.Insert(".giticket", giticketTreeID, git.FilemodeTree)
	if err != nil {
		return err
	}

	debug.DebugMessage(debugFlag, "Writing root tree")
	newRootTreeID, err := rootTreeBuilder.Write()
	if err != nil {
		return err
	}
	debug.DebugMessage(debugFlag, "Wrote root tree: "+newRootTreeID.String())

	debug.DebugMessage(debugFlag, "Getting new root tree for commit")
	newRootTree, err := thisRepo.LookupTree(newRootTreeID)
	if err != nil {
		return err
	}
	defer newRootTree.Free()

	debug.DebugMessage(debugFlag, "Creating commit")
	commitID, err := thisRepo.CreateCommit("refs/heads/"+branchName, author, author, commitMessage, newRootTree, parentCommit)
	if err != nil {
		return err
	}

	debug.DebugMessage(debugFlag, "Created commit: "+commitID.String())
	return nil
}
