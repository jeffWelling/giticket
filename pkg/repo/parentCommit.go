package repo

import (
	"errors"

	"github.com/jeffWelling/giticket/pkg/debug"
	git "github.com/jeffwelling/git2go/v37"
)

// Take a repo and a branch name, lookup the branch name and return the parent
// commit
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

// Take a git commit, and return a tree builder for the tree the commit points
// to, the tree builder must be freed when done
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

// Take a tree, a repo, and an entry name for a sub-tree, and return the
// sub-tree
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
