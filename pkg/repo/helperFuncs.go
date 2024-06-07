package repo

import (
	"errors"
	"fmt"
	"os"
	"slices"
	"strconv"
	"strings"
	"testing"

	git "github.com/jeffwelling/git2go/v37"
	"github.com/jeffwelling/giticket/pkg/common"
	"github.com/jeffwelling/giticket/pkg/debug"
	"gopkg.in/yaml.v2"
)

// Create a git repository, and then initialize giticket in that git repository
// and create a ticket
func InitGitAndInitGiticket(t *testing.T) error {
	debugFlag := true
	// Initialize git
	_ = common.InitGit(t)

	gitConfigContent := `# This is Git's per-user configuration file.
[user]
# Please adapt and uncomment the following lines:
	name = John Smith
	email = jsmith@example.com`

	// Write gitConfigContent to ./.gitconfig
	err := os.WriteFile(".gitconfig", []byte(gitConfigContent), 0644)
	if err != nil {
		return err
	}

	// Initialize the giticket branch
	HandleInitGiticket(debugFlag)

	debug.DebugMessage(debugFlag, "Opening git repository")
	thisRepo, err := git.OpenRepository(".")
	if err != nil {
		return err
	}

	debug.DebugMessage(debugFlag, "Getting parent commit from branch '"+common.BranchName+"'")
	parentCommit, err := GetParentCommit(thisRepo, common.BranchName, debugFlag)
	if err != nil {
		return err
	}

	// Get value for .giticket/next_ticket_id
	tree, err := parentCommit.Tree()
	if err != nil {
		return err
	}
	defer tree.Free()

	giticketTreeEntry, err := tree.EntryByPath(".giticket")
	if err != nil {
		return err
	}

	// Lookup giticketTreeEntry
	giticketTree, err := thisRepo.LookupTree(giticketTreeEntry.Id)
	if err != nil {
		return err
	}
	defer giticketTree.Free()

	NTIDEntry, err := giticketTree.EntryByPath("next_ticket_id")
	if err != nil {
		return err
	}

	NTIDBlob, err := thisRepo.LookupBlob(NTIDEntry.Id)
	if err != nil {
		return err
	}
	defer NTIDBlob.Free()

	// read value of blob as int
	s := strings.TrimSpace(string(NTIDBlob.Contents()))

	// Convert string to int
	ticketID, err := strconv.Atoi(s)
	if err != nil {
		return err
	}

	// Increment ticketID and write it as a blob
	i := ticketID + 1
	debug.DebugMessage(debugFlag, "incrementing next ticket ID in .giticket/next_ticket_id, is now: "+strconv.Itoa(i))
	NTIDBlobOID, err := thisRepo.CreateBlobFromBuffer([]byte(strconv.Itoa(i)))
	if err != nil {
		return err
	}
	debug.DebugMessage(debugFlag, "NTIDBlobOID: "+NTIDBlobOID.String())

	rootTreeBuilder, previousCommitTree, err := TreeBuilderFromCommit(parentCommit, thisRepo, debugFlag)
	if err != nil {
		return err
	}
	defer rootTreeBuilder.Free()

	giticketTree, err = GetSubTreeByName(previousCommitTree, thisRepo, ".giticket", debugFlag)
	if err != nil {
		return err
	}
	defer giticketTree.Free()

	// Create a TreeBuilder from the previous tree for giticket so we can add a
	// ticket under .gititcket/tickets and change the value of
	// .giticket/next_ticket_id
	debug.DebugMessage(debugFlag, "creating tree builder for .giticket tree: "+giticketTree.Id().String())
	giticketTreeBuilder, err := thisRepo.TreeBuilderFromTree(giticketTree)
	if err != nil {
		return err
	}
	defer giticketTreeBuilder.Free()

	// Insert the blob for next_ticket_id into the TreeBuilder for giticket
	// This essentially saves the file to .giticket/next_ticket_id, but we then need to
	// save the directory to the parent, all the way up to the commit.
	debug.DebugMessage(debugFlag, "inserting next ticket ID into .giticket/next_ticket_id: "+NTIDBlobOID.String())
	err = giticketTreeBuilder.Insert("next_ticket_id", NTIDBlobOID, git.FilemodeBlob)
	if err != nil {
		return err
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
			return err
		}
		defer giticketTicketsTreeBuilder.Free()
	} else {
		debug.DebugMessage(debugFlag, "looking up tree for .giticket/tickets: "+_giticketTicketsTreeID.Id.String())
		giticketTicketsTree, err := thisRepo.LookupTree(_giticketTicketsTreeID.Id)
		if err != nil {
			return err
		}
		defer giticketTicketsTree.Free()

		debug.DebugMessage(debugFlag, "creating tree builder for .giticket/tickets tree: "+giticketTicketsTree.Id().String())
		giticketTicketsTreeBuilder, err = thisRepo.TreeBuilderFromTree(giticketTicketsTree)
		if err != nil {
			return err
		}
	}

	debug.DebugMessage(debugFlag, "creating and populating ticket")
	// Craft the ticket, but avoid importing 'ticket' due to cyclical
	// dependencies
	ticketHereDoc := `title: My first ticket
description: This is an awesome description.
labels:
- bugfix
- ux
priority: 1
severity: 1
status: open
comments:
- id: 1
  created: 1716538263
  body: First comment
  author: John Smith <jsmith@example.com>
- id: 2
  created: 1716538445
  body: Inverted the tardis polarity
  author: Bob Franks <bfranks@example.com>
next_comment_id: 3
id: 1
created: 1716538263`

	// Write ticket
	ticketBlobOID, err := thisRepo.CreateBlobFromBuffer([]byte(ticketHereDoc))
	if err != nil {
		return err
	}
	debug.DebugMessage(debugFlag, "writing ticket to .giticket/tickets: "+ticketBlobOID.String())

	// Add ticket to .giticket/tickets
	debug.DebugMessage(debugFlag, "adding ticket to .giticket/tickets: "+ticketBlobOID.String())
	err = giticketTicketsTreeBuilder.Insert("1__My_first_ticket", ticketBlobOID, git.FilemodeBlob)
	if err != nil {
		return err
	}

	// Save the tree and get the tree ID for .giticket/tickets
	giticketTicketsTreeID, err := giticketTicketsTreeBuilder.Write()
	if err != nil {
		return err
	}
	debug.DebugMessage(debugFlag, "saving .giticket/tickets:"+giticketTicketsTreeID.String())

	// Add 'ticket' directory to '.giticket' TreeBuilder, to update the
	// .giticket directory with the new tree for the updated tickets directory
	debug.DebugMessage(debugFlag, "adding ticket directory to .giticket: "+giticketTicketsTreeID.String())
	err = giticketTreeBuilder.Insert("tickets", giticketTicketsTreeID, git.FilemodeTree)
	if err != nil {
		return err
	}

	giticketTreeID, err := giticketTreeBuilder.Write()
	if err != nil {
		return err
	}
	debug.DebugMessage(debugFlag, "saving .giticket tree to root tree: "+giticketTreeID.String())

	// Update the root tree builder with the new .giticket directory
	debug.DebugMessage(debugFlag, "updating root tree with: "+giticketTreeID.String())
	err = rootTreeBuilder.Insert(".giticket", giticketTreeID, git.FilemodeTree)
	if err != nil {
		return err
	}

	// Save the new root tree and get the ID so we can lookup the tree for the
	// commit
	debug.DebugMessage(debugFlag, "saving root tree")
	rootTreeBuilderID, err := rootTreeBuilder.Write()
	if err != nil {
		return err
	}

	// Lookup the tree so we can use it in the commit
	debug.DebugMessage(debugFlag, "lookup root tree for commit: "+rootTreeBuilderID.String())
	rootTree, err := thisRepo.LookupTree(rootTreeBuilderID)
	if err != nil {
		return err
	}
	defer rootTree.Free()

	// Get author data by reading .git configs
	debug.DebugMessage(debugFlag, "getting author data")
	author, err := common.GetAuthor(thisRepo)
	if err != nil {
		return err
	}

	// commit and update 'giticket' branch
	debug.DebugMessage(debugFlag, "creating commit with message 'Creating ticket 1__My_first_ticket'")
	commitID, err := thisRepo.CreateCommit("refs/heads/giticket", author, author, "Creating ticket 1__My_first_ticket", rootTree, parentCommit)
	if err != nil {
		return err
	}
	debug.DebugMessage(debugFlag, "successfully created commit with ID: "+commitID.String())
	return nil
}

// Check whether the comment identified by commentID exists in the branchName,
// emit debug messages if debugFlag is set.
// This function is intended for use in tests, because of that it doesn't re-use
// functions that duplicate some of the functionality seen here.
func CommentExists(
	branchName string,
	commentID string,
	debugFlag bool,
) (bool, error) {
	debug.DebugMessage(debugFlag, "Checking if comment with ID "+commentID+" exists in branch "+branchName)

	type localComment struct {
		ID      int
		Created int64
		Body    string
		Author  string
	}
	type localTicket struct {
		Title         string
		Description   string
		Labels        []string
		Priority      int
		Severity      int
		Status        string
		Comments      []localComment
		NextCommentID int `yaml:"next_comment_id" json:"next_comment_id"`

		// Set automatically
		ID      int
		Created int64
	}
	var (
		t          localTicket
		commentIDs []string
	)

	// Open the git repository
	thisRepo, _, giticketTicketsSubTree, _, err := openGitAndReturnGiticketThings(branchName, debugFlag)
	if err != nil {
		return false, err
	}
	debug.DebugMessage(debugFlag, "Walking giticket tickets tree")
	giticketTicketsSubTree.Walk(func(name string, entry *git.TreeEntry) error {
		ticketFile, err := thisRepo.LookupBlob(entry.Id)
		if err != nil {
			return fmt.Errorf("Error walking the tickets tree and looking up the entry ID: %s", err)
		}
		defer ticketFile.Free()

		t = localTicket{}
		// Unmarshal the ticket which is yaml
		err = yaml.Unmarshal(ticketFile.Contents(), &t)
		if err != nil {
			return fmt.Errorf("Error unmarshalling yaml ticket from file in tickets directory: %s", err)
		}

		// For each comment in the ticket
		for _, comment := range t.Comments {
			commentIDs = append(commentIDs, strconv.Itoa(t.ID)+"-"+strconv.Itoa(comment.ID))
		}

		return nil
	})
	debug.DebugMessage(debugFlag, "Finished walking giticket tickets tree")

	// Check if the comment exists
	return slices.Contains(commentIDs, commentID), nil
}

func TicketExists(
	ticketFilename string,
	debugFlag bool,
) (bool, error) {
	debug.DebugMessage(debugFlag, "Checking if ticket with filename "+ticketFilename+" exists")

	// Open the git repository
	_, _, giticketTicketsSubTree, _, err := openGitAndReturnGiticketThings("giticket", true)
	if err != nil {
		debug.DebugMessage(debugFlag, "Error opening git and returning giticket things: "+err.Error())
		return false, err
	}

	// Print every entry name in giticketTicketsSubTree
	gotcha := false
	giticketTicketsSubTree.Walk(func(name string, entry *git.TreeEntry) error {
		debug.DebugMessage(true, "Found entry: "+entry.Name)
		if entry.Name == ticketFilename {
			gotcha = true
		}
		return nil
	})

	if gotcha {
		return true, nil
	}

	debug.DebugMessage(debugFlag, "Checking if ticket with filename "+ticketFilename+" exists in .giticket/tickets")
	fileEntry, err := giticketTicketsSubTree.EntryByPath(ticketFilename)
	if err != nil {
		debug.DebugMessage(debugFlag, "Error getting file entry called '"+ticketFilename+"'")
		return false, err
	}
	if fileEntry != nil {
		debug.DebugMessage(debugFlag, "Error getting file entry called '"+ticketFilename+"', was nil")
		return false, nil
	}

	// Print every entry name in giticketTicketsSubTree
	giticketTicketsSubTree.Walk(func(name string, entry *git.TreeEntry) error {
		debug.DebugMessage(true, "Found entry: "+entry.Id.String())
		return nil
	})

	return true, nil
}

// Returns:
// - Pointer to the git repository
// - A pointer to the git tree for the .giticket sub-tree
// - A pointer to the git tree for the .giticket/tickets sub-tree
// - A pointer to the parent commit
// - An error if something goes wrong
func openGitAndReturnGiticketThings(branchName string, debugFlag bool) (*git.Repository, *git.Tree, *git.Tree, *git.Commit, error) {
	debug.DebugMessage(debugFlag, "Opening git and returning giticket things")
	// Open the git repository
	debug.DebugMessage(debugFlag, "Opening git repository")
	thisRepo, err := git.OpenRepository(".")
	if err != nil {
		return nil, nil, nil, nil, err
	}

	// Lookup the branch
	debug.DebugMessage(debugFlag, "Looking up branch: "+branchName)
	branch, err := thisRepo.LookupBranch(branchName, git.BranchLocal)
	if err != nil {
		return nil, nil, nil, nil, err
	}

	// Lookup the commit the branch references
	debug.DebugMessage(debugFlag, "Looking up commit: "+branch.Target().String())
	parentCommit, err := thisRepo.LookupCommit(branch.Target())
	if err != nil {
		return nil, nil, nil, nil, err
	}

	// Get the parent commit of the branch
	debug.DebugMessage(debugFlag, "Getting parent commit from branch '"+branchName+"'")
	parentCommit, err = GetParentCommit(thisRepo, branchName, debugFlag)
	if err != nil {
		return nil, nil, nil, nil, err
	}

	// Get root tree from commit
	debug.DebugMessage(debugFlag, "Getting root tree from commit: "+parentCommit.Id().String())
	tree, err := parentCommit.Tree()
	if err != nil {
		return nil, nil, nil, nil, err
	}

	// Get '.gititickets' subtree
	giticketSubTreeEntry := tree.EntryByName(".giticket")
	if giticketSubTreeEntry == nil {
		return nil, nil, nil, nil, errors.New("Subtree '.giticket' not found")
	}
	debug.DebugMessage(debugFlag, "Looked up tree entry for giticket: "+giticketSubTreeEntry.Id.String())

	giticketSubTree, err := thisRepo.LookupTree(giticketSubTreeEntry.Id)
	if err != nil {
		return nil, nil, nil, nil, err
	}
	debug.DebugMessage(debugFlag, "Found giticket tree: "+giticketSubTree.Id().String())

	giticketTicketsSubTreeEntry := giticketSubTree.EntryByName("tickets")
	if giticketTicketsSubTreeEntry == nil {
		return nil, nil, nil, nil, errors.New("Subtree 'tickets' not found")
	}
	debug.DebugMessage(debugFlag, "Looked up giticket tickets tree entry: "+giticketTicketsSubTreeEntry.Id.String())

	giticketTicketsSubTree, err := thisRepo.LookupTree(giticketTicketsSubTreeEntry.Id)
	if err != nil {
		return nil, nil, nil, nil, err
	}
	debug.DebugMessage(debugFlag, "Found giticket tickets tree: "+giticketTicketsSubTree.Id().String())

	return thisRepo, giticketSubTree, giticketTicketsSubTree, parentCommit, nil
}
