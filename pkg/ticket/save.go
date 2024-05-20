package ticket

import (
	"github.com/jeffWelling/giticket/pkg/common"
	"github.com/jeffWelling/giticket/pkg/debug"
	git "github.com/jeffwelling/git2go/v37"
	"gopkg.in/yaml.v2"
)

func SaveTicket(t *Ticket, repo *git.Repository, branchName string, debugFlag bool) {
	// turn the ticket into a yaml string
	debug.DebugMessage(debugFlag, "Yaml marshal of ticket")
	yamlTicket, err := yaml.Marshal(t)
	if err != nil {
		panic(err)
	}

	// Find the branch and its target commit
	debug.DebugMessage(debugFlag, "Looking up branch: "+branchName)
	branch, err := repo.LookupBranch(branchName, git.BranchLocal)
	if err != nil {
		panic(err)
	}

	// Lookup the commit the branch references
	debug.DebugMessage(debugFlag, "Looking up commit: "+branch.Target().String())
	parentCommit, err := repo.LookupCommit(branch.Target())
	if err != nil {
		panic(err)
	}

	ticketBlobID, err := repo.CreateBlobFromBuffer(yamlTicket)
	if err != nil {
		panic(err)
	}

	// Lookup the tree from the previous commit, we need this to build a tree
	// that includes the files from the previous commit.
	debug.DebugMessage(debugFlag, "looking up tree from parent commit, tree ID: "+parentCommit.TreeId().String())
	previousCommitTree, err := parentCommit.Tree()
	if err != nil {
		panic(err)
	}
	defer previousCommitTree.Free()

	// Create a TreeBuilder from the previous commit's tree, so that we can
	// update it with our changes to the .giticket directory
	debug.DebugMessage(debugFlag, "creating root tree builder from previous commit")
	rootTreeBuilder, err := repo.TreeBuilderFromTree(previousCommitTree)
	if err != nil {
		panic(err)
	}
	defer rootTreeBuilder.Free()

	// Get the TreeEntry for ".giticket" from the previous commit so we can get
	// the tree for .giticket
	debug.DebugMessage(debugFlag, "looking up tree entry for .giticket: "+previousCommitTree.EntryByName(".giticket").Id.String())
	giticketTreeEntry := previousCommitTree.EntryByName(".giticket")

	// Lookup tree for giticket
	debug.DebugMessage(debugFlag, "looking up tree for .giticket: "+giticketTreeEntry.Id.String())
	giticketTree, err := repo.LookupTree(giticketTreeEntry.Id)
	if err != nil {
		panic(err)
	}
	defer giticketTree.Free()

	// Create a TreeBuilder from the previous tree for giticket so we can add a
	// ticket under .gititcket/tickets and change the value of
	// .giticket/next_ticket_id
	debug.DebugMessage(debugFlag, "creating tree builder for .giticket tree: "+giticketTree.Id().String())
	giticketTreeBuilder, err := repo.TreeBuilderFromTree(giticketTree)
	if err != nil {
		panic(err)
	}
	defer giticketTreeBuilder.Free()

	giticketTicketsTreeID := giticketTree.EntryByName("tickets")
	debug.DebugMessage(debugFlag, "looking up tree for .giticket/tickets: "+giticketTicketsTreeID.Id.String())
	giticketTicketsTree, err := repo.LookupTree(giticketTicketsTreeID.Id)
	if err != nil {
		panic(err)
	}
	defer giticketTicketsTree.Free()

	// Tree builder for .giticket/tickets
	debug.DebugMessage(debugFlag, "creating tree builder for .giticket/tickets tree: "+giticketTicketsTree.Id().String())
	giticketTicketsTreeBuilder, err := repo.TreeBuilderFromTree(giticketTicketsTree)
	if err != nil {
		panic(err)
	}

	// Add ticket to .giticket/tickets
	debug.DebugMessage(debugFlag, "adding ticket to .giticket/tickets: "+ticketBlobID.String())
	err = giticketTicketsTreeBuilder.Insert(t.TicketFilename(), ticketBlobID, git.FilemodeBlob)
	if err != nil {
		panic(err)
	}

	// Save the tree and get the tree ID for .giticket/tickets
	newGiticketTicketsTreeID, err := giticketTicketsTreeBuilder.Write()
	if err != nil {
		panic(err)
	}
	debug.DebugMessage(debugFlag, "saving .giticket/tickets:"+newGiticketTicketsTreeID.String())

	// Add 'ticket' directory to '.giticket' TreeBuilder, to update the
	// .giticket directory with the new tree for the updated tickets directory
	debug.DebugMessage(debugFlag, "adding ticket directory to .giticket: "+newGiticketTicketsTreeID.String())
	err = giticketTreeBuilder.Insert("tickets", newGiticketTicketsTreeID, git.FilemodeTree)
	if err != nil {
		panic(err)
	}

	giticketTreeID, err := giticketTreeBuilder.Write()
	if err != nil {
		panic(err)
	}
	debug.DebugMessage(debugFlag, "saving .giticket tree to root tree: "+giticketTreeID.String())

	// Update the root tree builder with the new .giticket directory
	debug.DebugMessage(debugFlag, "updating root tree with: "+giticketTreeID.String())
	err = rootTreeBuilder.Insert(".giticket", giticketTreeID, git.FilemodeTree)
	if err != nil {
		panic(err)
	}

	// Save the new root tree and get the ID so we can lookup the tree for the
	// commit
	debug.DebugMessage(debugFlag, "saving root tree")
	rootTreeBuilderID, err := rootTreeBuilder.Write()
	if err != nil {
		panic(err)
	}

	// Lookup the tree so we can use it in the commit
	debug.DebugMessage(debugFlag, "lookup root tree for commit: "+rootTreeBuilderID.String())
	rootTree, err := repo.LookupTree(rootTreeBuilderID)
	if err != nil {
		panic(err)
	}
	defer rootTree.Free()

	// Get author data by reading .git configs
	debug.DebugMessage(debugFlag, "getting author data")
	author := common.GetAuthor(repo)

	// commit and update 'giticket' branch
	debug.DebugMessage(debugFlag, "creating commit")
	commitID, err := repo.CreateCommit("refs/heads/giticket", author, author, "Creating ticket "+t.TicketFilename(), rootTree, parentCommit)
	if err != nil {
		panic(err)
	}

	debug.DebugMessage(debugFlag, "commit created: "+commitID.String())
	return
}
