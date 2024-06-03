package repo

import (
	"errors"
	"fmt"
	"slices"
	"strconv"
	"testing"

	git "github.com/jeffwelling/git2go/v37"
	"github.com/jeffwelling/giticket/pkg/common"
	"github.com/jeffwelling/giticket/pkg/debug"
	"gopkg.in/yaml.v2"
)

// Create a git repository, and then initialize giticket in that git repository
func InitGitAndInitGiticket(t *testing.T) {
	// Initialize git
	_ = common.InitGit(t)

	// Initialize the giticket branch
	HandleInitGiticket(false)
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

	// Open the git repository
	thisRepo, err := git.OpenRepository(".")
	if err != nil {
		return false, err
	}

	branch, err := thisRepo.LookupBranch(branchName, git.BranchLocal)
	if err != nil {
		return false, err
	}

	// Lookup the commit the branch references
	debug.DebugMessage(debugFlag, "Looking up commit: "+branch.Target().String())
	parentCommit, err := thisRepo.LookupCommit(branch.Target())
	if err != nil {
		return false, err
	}

	// Get the parent commit of the branch
	debug.DebugMessage(debugFlag, "Getting parent commit from branch '"+branchName+"'")
	parentCommit, err = GetParentCommit(thisRepo, branchName, debugFlag)
	if err != nil {
		return false, err
	}

	// Get tree from commit
	tree, err := parentCommit.Tree()
	if err != nil {
		return false, err
	}

	// Get '.gititickets' subtree
	giticketSubTreeEntry := tree.EntryByName(".giticket")
	if giticketSubTreeEntry == nil {
		return false, errors.New("Subtree 'tickets' not found")
	}
	debug.DebugMessage(debugFlag, "Looked up tree entry: "+giticketSubTreeEntry.Id.String())

	giticketSubTree, err := thisRepo.LookupTree(giticketSubTreeEntry.Id)
	if err != nil {
		return false, err
	}
	debug.DebugMessage(debugFlag, "Found tree: "+giticketSubTree.Id().String())

	giticketTicketsSubTreeEntry := tree.EntryByName("tickets")
	if giticketTicketsSubTreeEntry == nil {
		return false, errors.New("Subtree 'tickets' not found")
	}
	debug.DebugMessage(debugFlag, "Looked up tree entry: "+giticketTicketsSubTreeEntry.Id.String())

	giticketTicketsSubTree, err := thisRepo.LookupTree(giticketTicketsSubTreeEntry.Id)
	if err != nil {
		return false, err
	}
	debug.DebugMessage(debugFlag, "Found tree: "+giticketSubTree.Id().String())

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

	// Check if the comment exists
	return slices.Contains(commentIDs, commentID), nil
}
