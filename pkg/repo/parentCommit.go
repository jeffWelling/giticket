package repo

import (
	"errors"

	git "github.com/jeffwelling/git2go/v37"
	"github.com/jeffwelling/giticket/pkg/debug"
)

// GetParentCommit takes a pointer to a git repository and a branch name and
// debugFlag, and returns a pointer to the commit at the tip of branchName. It
// returns a pointer to that commit and an error if there was one.
func GetParentCommit(repo *git.Repository, branchName string, debugFlag bool) (*git.Commit, error) {
	// Find the branch and its target commit
	debug.DebugMessage(debugFlag, "Looking up branch: "+branchName)
	branch, err := repo.LookupBranch(branchName, git.BranchLocal)
	if err != nil {
		return nil, err
	}

	// Lookup the commit the branch references
	debug.DebugMessage(debugFlag, "Looking up commit: "+branch.Target().String())
	parentCommit, err := repo.LookupCommit(branch.Target())
	if err != nil {
		return nil, err
	}

	return parentCommit, nil
}

// TreeBuilderFromCommit takes a pointer to a git commit and a pointer to a git
// repository and a debugFlag. It looks up the tree from the commit and returns
// a pointer to a treeBuilder from that commit with an error if there was one.
func TreeBuilderFromCommit(commit *git.Commit, thisRepo *git.Repository, debugFlag bool) (*git.TreeBuilder, *git.Tree, error) {
	debug.DebugMessage(debugFlag, "Looking up tree from parent commit, tree ID: "+commit.TreeId().String())
	previousCommitTree, err := commit.Tree()
	if err != nil {
		return nil, nil, err
	}

	debug.DebugMessage(debugFlag, "Creating root tree builder from previous commit")
	rootTreeBuilder, err := thisRepo.TreeBuilderFromTree(previousCommitTree)
	if err != nil {
		return nil, nil, err
	}

	return rootTreeBuilder, previousCommitTree, nil
}

// GetSubTreeByName takes a pointer to a git tree, a pointer to a git repository
// and a sub-tree name and a debugFlag. It returns a pointer to the sub-tree
// with the given name and an error if there was one.
func GetSubTreeByName(parentTree *git.Tree, repo *git.Repository, subTreeName string, debugFlag bool) (*git.Tree, error) {
	subTreeEntry := parentTree.EntryByName(subTreeName)
	if subTreeEntry == nil {
		return nil, errors.New("Subtree " + subTreeName + " not found")
	}
	debug.DebugMessage(debugFlag, "Looked up tree entry: "+subTreeEntry.Id.String())

	subTree, err := repo.LookupTree(subTreeEntry.Id)
	if err != nil {
		return nil, err
	}
	debug.DebugMessage(debugFlag, "Found tree: "+subTree.Id().String())

	return subTree, nil
}
