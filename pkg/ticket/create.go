package ticket

import (
	"strconv"
	"time"

	git "github.com/jeffwelling/git2go/v37"
	"github.com/jeffwelling/giticket/pkg/common"
	"github.com/jeffwelling/giticket/pkg/debug"
	"github.com/jeffwelling/giticket/pkg/repo"
	"github.com/jeffwelling/giticket/pkg/subcommand"
)

func Create(thisRepo *git.Repository, branchName string, subcommand subcommand.SubcommandInterface) (*git.Oid, string) {
	parentCommit, err := repo.GetParentCommit(thisRepo, branchName, subcommand.DebugFlag())
	if err != nil {
		panic(err)
	}

	// Get value for .giticket/next_ticket_id
	ticketID, err := ReadNextTicketID(thisRepo, parentCommit)
	if err != nil {
		panic(err)
	}
	debug.DebugMessage(subcommand.DebugFlag(), "Next ticket ID: "+strconv.Itoa(ticketID))

	// Increment ticketID and write it as a blob
	i := ticketID + 1
	debug.DebugMessage(subcommand.DebugFlag(), "incrementing next ticket ID in .giticket/next_ticket_id, is now: "+strconv.Itoa(i))
	NTIDBlobOID, err := thisRepo.CreateBlobFromBuffer([]byte(strconv.Itoa(i)))
	debug.DebugMessage(subcommand.DebugFlag(), "NTIDBlobOID: "+NTIDBlobOID.String())

	rootTreeBuilder, previousCommitTree, err := repo.TreeBuilderFromCommit(parentCommit, thisRepo, subcommand.DebugFlag())
	if err != nil {
		panic(err)
	}
	defer rootTreeBuilder.Free()

	giticketTree, err := repo.GetSubTreeByName(previousCommitTree, thisRepo, ".giticket", subcommand.DebugFlag())
	defer giticketTree.Free()

	// Create a TreeBuilder from the previous tree for giticket so we can add a
	// ticket under .gititcket/tickets and change the value of
	// .giticket/next_ticket_id
	debug.DebugMessage(subcommand.DebugFlag(), "creating tree builder for .giticket tree: "+giticketTree.Id().String())
	giticketTreeBuilder, err := thisRepo.TreeBuilderFromTree(giticketTree)
	if err != nil {
		panic(err)
	}
	defer giticketTreeBuilder.Free()

	// Insert the blob for next_ticket_id into the TreeBuilder for giticket
	// This essentially saves the file to .giticket/next_ticket_id, but we then need to
	// save the directory to the parent, all the way up to the commit.
	debug.DebugMessage(subcommand.DebugFlag(), "inserting next ticket ID into .giticket/next_ticket_id: "+NTIDBlobOID.String())
	err = giticketTreeBuilder.Insert("next_ticket_id", NTIDBlobOID, git.FilemodeBlob)
	if err != nil {
		panic(err)
	}

	var giticketTicketsTreeBuilder *git.TreeBuilder
	_giticketTicketsTreeID := giticketTree.EntryByName("tickets")
	if _giticketTicketsTreeID == nil {
		// Create the tickets directory

		// Get a TreeBuilder for .giticket/tickets so we can add the ticket to that
		// directory
		debug.DebugMessage(subcommand.DebugFlag(), "creating empty tree builder for .giticket/tickets")
		giticketTicketsTreeBuilder, err = thisRepo.TreeBuilder()
		if err != nil {
			panic(err)
		}
		defer giticketTicketsTreeBuilder.Free()
	} else {
		debug.DebugMessage(subcommand.DebugFlag(), "looking up tree for .giticket/tickets: "+_giticketTicketsTreeID.Id.String())
		giticketTicketsTree, err := thisRepo.LookupTree(_giticketTicketsTreeID.Id)
		if err != nil {
			panic(err)
		}
		defer giticketTicketsTree.Free()

		debug.DebugMessage(subcommand.DebugFlag(), "creating tree builder for .giticket/tickets tree: "+giticketTicketsTree.Id().String())
		giticketTicketsTreeBuilder, err = thisRepo.TreeBuilderFromTree(giticketTicketsTree)
		if err != nil {
			panic(err)
		}
	}

	debug.DebugMessage(subcommand.DebugFlag(), "creating and populating ticket")
	// Craft the ticket
	t := Ticket{}
	t.Created = time.Now().Unix()
	t.Title = subcommand.Parameters()["title"].(string)
	t.Description = subcommand.Parameters()["description"].(string)
	t.Labels = subcommand.Parameters()["labels"].([]string)
	t.Priority = subcommand.Parameters()["priority"].(int)
	t.Severity = subcommand.Parameters()["severity"].(int)
	t.Status = subcommand.Parameters()["status"].(string)
	t.ID = ticketID
	t.Comments = subcommand.Parameters()["comments"].([]Comment)
	t.NextCommentID = subcommand.Parameters()["next_comment_id"].(int)
	// FIXME Add a way to parse comments from initial ticket creation

	// Write ticket
	ticketBlobOID, err := thisRepo.CreateBlobFromBuffer(t.TicketToYaml())
	if err != nil {
		panic(err)
	}
	debug.DebugMessage(subcommand.DebugFlag(), "writing ticket to .giticket/tickets: "+ticketBlobOID.String())

	// Add ticket to .giticket/tickets
	debug.DebugMessage(subcommand.DebugFlag(), "adding ticket to .giticket/tickets: "+ticketBlobOID.String())
	err = giticketTicketsTreeBuilder.Insert(t.TicketFilename(), ticketBlobOID, git.FilemodeBlob)
	if err != nil {
		panic(err)
	}

	// Save the tree and get the tree ID for .giticket/tickets
	giticketTicketsTreeID, err := giticketTicketsTreeBuilder.Write()
	if err != nil {
		panic(err)
	}
	debug.DebugMessage(subcommand.DebugFlag(), "saving .giticket/tickets:"+giticketTicketsTreeID.String())

	// Add 'ticket' directory to '.giticket' TreeBuilder, to update the
	// .giticket directory with the new tree for the updated tickets directory
	debug.DebugMessage(subcommand.DebugFlag(), "adding ticket directory to .giticket: "+giticketTicketsTreeID.String())
	err = giticketTreeBuilder.Insert("tickets", giticketTicketsTreeID, git.FilemodeTree)
	if err != nil {
		panic(err)
	}

	giticketTreeID, err := giticketTreeBuilder.Write()
	if err != nil {
		panic(err)
	}
	debug.DebugMessage(subcommand.DebugFlag(), "saving .giticket tree to root tree: "+giticketTreeID.String())

	// Update the root tree builder with the new .giticket directory
	debug.DebugMessage(subcommand.DebugFlag(), "updating root tree with: "+giticketTreeID.String())
	err = rootTreeBuilder.Insert(".giticket", giticketTreeID, git.FilemodeTree)
	if err != nil {
		panic(err)
	}

	// Save the new root tree and get the ID so we can lookup the tree for the
	// commit
	debug.DebugMessage(subcommand.DebugFlag(), "saving root tree")
	rootTreeBuilderID, err := rootTreeBuilder.Write()
	if err != nil {
		panic(err)
	}

	// Lookup the tree so we can use it in the commit
	debug.DebugMessage(subcommand.DebugFlag(), "lookup root tree for commit: "+rootTreeBuilderID.String())
	rootTree, err := thisRepo.LookupTree(rootTreeBuilderID)
	if err != nil {
		panic(err)
	}
	defer rootTree.Free()

	// Get author data by reading .git configs
	debug.DebugMessage(subcommand.DebugFlag(), "getting author data")
	author := common.GetAuthor(thisRepo)

	// commit and update 'giticket' branch
	debug.DebugMessage(subcommand.DebugFlag(), "creating commit")
	commitID, err := thisRepo.CreateCommit("refs/heads/giticket", author, author, "Creating ticket "+t.TicketFilename(), rootTree, parentCommit)
	if err != nil {
		panic(err)
	}

	return commitID, t.TicketFilename()
}
