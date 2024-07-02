package ticket

import (
	"strconv"

	git "github.com/jeffwelling/git2go/v37"
	"github.com/jeffwelling/giticket/pkg/common"
	"github.com/jeffwelling/giticket/pkg/debug"
	"github.com/jeffwelling/giticket/pkg/repo"
)

func HandleCreate(
	branchName string,
	created int64,
	title string,
	description string,
	labels []string,
	priority int,
	severity int,
	status string,
	comments []Comment,
	nextCommentId int,
	debugFlag bool,
) (int, string, error) {
	debug.DebugMessage(debugFlag, "Opening git repository")
	thisRepo, err := git.OpenRepository(".")
	if err != nil {
		return 0, "", err
	}

	debug.DebugMessage(debugFlag, "Getting parent commit from branch '"+branchName+"'")
	parentCommit, err := repo.GetParentCommit(thisRepo, branchName, debugFlag)
	if err != nil {
		return 0, "", err
	}

	// Get value for .giticket/next_ticket_id
	ticketID, err := ReadNextTicketID(thisRepo, parentCommit)
	if err != nil {
		return 0, "", err
	}
	debug.DebugMessage(debugFlag, "Next ticket ID: "+strconv.Itoa(ticketID))

	// Increment ticketID and write it as a blob
	i := ticketID + 1
	debug.DebugMessage(debugFlag, "incrementing next ticket ID in .giticket/next_ticket_id, is now: "+strconv.Itoa(i))
	NTIDBlobOID, err := thisRepo.CreateBlobFromBuffer([]byte(strconv.Itoa(i)))
	if err != nil {
		return 0, "", err
	}
	debug.DebugMessage(debugFlag, "NTIDBlobOID: "+NTIDBlobOID.String())

	rootTreeBuilder, previousCommitTree, err := repo.TreeBuilderFromCommit(parentCommit, thisRepo, debugFlag)
	if err != nil {
		return 0, "", err
	}
	defer rootTreeBuilder.Free()

	giticketTree, err := repo.GetSubTreeByName(previousCommitTree, thisRepo, ".giticket", debugFlag)
	if err != nil {
		return 0, "", err
	}
	defer giticketTree.Free()

	// Create a TreeBuilder from the previous tree for giticket so we can add a
	// ticket under .gititcket/tickets and change the value of
	// .giticket/next_ticket_id
	debug.DebugMessage(debugFlag, "creating tree builder for .giticket tree: "+giticketTree.Id().String())
	giticketTreeBuilder, err := thisRepo.TreeBuilderFromTree(giticketTree)
	if err != nil {
		return 0, "", err
	}
	defer giticketTreeBuilder.Free()

	// Insert the blob for next_ticket_id into the TreeBuilder for giticket
	// This essentially saves the file to .giticket/next_ticket_id, but we then need to
	// save the directory to the parent, all the way up to the commit.
	debug.DebugMessage(debugFlag, "inserting next ticket ID into .giticket/next_ticket_id: "+NTIDBlobOID.String())
	err = giticketTreeBuilder.Insert("next_ticket_id", NTIDBlobOID, git.FilemodeBlob)
	if err != nil {
		return 0, "", err
	}

	var giticketTicketsTreeBuilder *git.TreeBuilder
	_giticketTicketsTreeID := giticketTree.EntryByName("tickets")
	if _giticketTicketsTreeID == nil {
		// Create the tickets directory

		// Get a TreeBuilder for .giticket/tickets so we can add the ticket to that
		// directory
		debug.DebugMessage(debugFlag, "creating empty tree builder for .giticket/tickets")
		giticketTicketsTreeBuilder, err = thisRepo.TreeBuilder()
		if err != nil {
			return 0, "", err
		}
		defer giticketTicketsTreeBuilder.Free()
	} else {
		debug.DebugMessage(debugFlag, "looking up tree for .giticket/tickets: "+_giticketTicketsTreeID.Id.String())
		giticketTicketsTree, err := thisRepo.LookupTree(_giticketTicketsTreeID.Id)
		if err != nil {
			return 0, "", err
		}
		defer giticketTicketsTree.Free()

		debug.DebugMessage(debugFlag, "creating tree builder for .giticket/tickets tree: "+giticketTicketsTree.Id().String())
		giticketTicketsTreeBuilder, err = thisRepo.TreeBuilderFromTree(giticketTicketsTree)
		if err != nil {
			return 0, "", err
		}
	}

	debug.DebugMessage(debugFlag, "creating and populating ticket")
	// Craft the ticket
	t := Ticket{}
	t.Created = created
	t.Title = title
	t.Description = description
	t.Labels = labels
	t.Priority = priority
	t.Severity = severity
	t.Status = status
	t.ID = ticketID
	t.Comments = comments
	t.NextCommentID = nextCommentId

	// Write ticket
	ticketBlobOID, err := thisRepo.CreateBlobFromBuffer(t.TicketToYaml())
	if err != nil {
		return 0, "", err
	}
	debug.DebugMessage(debugFlag, "writing ticket to .giticket/tickets: "+ticketBlobOID.String())

	// Add ticket to .giticket/tickets
	debug.DebugMessage(debugFlag, "adding ticket to .giticket/tickets: "+ticketBlobOID.String())
	err = giticketTicketsTreeBuilder.Insert(t.TicketFilename(), ticketBlobOID, git.FilemodeBlob)
	if err != nil {
		return 0, "", err
	}

	// Save the tree and get the tree ID for .giticket/tickets
	giticketTicketsTreeID, err := giticketTicketsTreeBuilder.Write()
	if err != nil {
		return 0, "", err
	}
	debug.DebugMessage(debugFlag, "saving .giticket/tickets:"+giticketTicketsTreeID.String())

	// Add 'ticket' directory to '.giticket' TreeBuilder, to update the
	// .giticket directory with the new tree for the updated tickets directory
	debug.DebugMessage(debugFlag, "adding ticket directory to .giticket: "+giticketTicketsTreeID.String())
	err = giticketTreeBuilder.Insert("tickets", giticketTicketsTreeID, git.FilemodeTree)
	if err != nil {
		return 0, "", err
	}

	giticketTreeID, err := giticketTreeBuilder.Write()
	if err != nil {
		return 0, "", err
	}
	debug.DebugMessage(debugFlag, "saving .giticket tree to root tree: "+giticketTreeID.String())

	// Update the root tree builder with the new .giticket directory
	debug.DebugMessage(debugFlag, "updating root tree with: "+giticketTreeID.String())
	err = rootTreeBuilder.Insert(".giticket", giticketTreeID, git.FilemodeTree)
	if err != nil {
		return 0, "", err
	}

	// Save the new root tree and get the ID so we can lookup the tree for the
	// commit
	debug.DebugMessage(debugFlag, "saving root tree")
	rootTreeBuilderID, err := rootTreeBuilder.Write()
	if err != nil {
		return 0, "", err
	}

	// Lookup the tree so we can use it in the commit
	debug.DebugMessage(debugFlag, "lookup root tree for commit: "+rootTreeBuilderID.String())
	rootTree, err := thisRepo.LookupTree(rootTreeBuilderID)
	if err != nil {
		return 0, "", err
	}
	defer rootTree.Free()

	// Get author data by reading .git configs
	debug.DebugMessage(debugFlag, "getting author data")
	author, err := common.GetAuthor(thisRepo)
	if err != nil {
		return 0, "", err
	}

	// commit and update 'giticket' branch
	debug.DebugMessage(debugFlag, "creating commit")
	_, err = thisRepo.CreateCommit("refs/heads/giticket", author, author, "Creating ticket "+t.TicketFilename(), rootTree, parentCommit)
	if err != nil {
		return 0, "", err
	}

	return ticketID, t.TicketFilename(), nil
}
