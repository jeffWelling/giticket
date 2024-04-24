package actions

import (
	"fmt"
	"log"
	"time"

	git "github.com/jeffwelling/git2go/v37"
)

// init() is used to register this action
func init() {
	action := new(ActionInit)
	registerAction("init", action)
}

type ActionInit struct {
}

// Execute() is called when this action is invoked
func (action *ActionInit) Execute() {
	fmt.Println("Initializing giticket")

	// Open an existing repository in the current directory
	repo, err := git.OpenRepository(".")
	if err != nil {
		log.Fatalf("Failed to open repo: %v", err)
	}

	// Create an index (staging area)
	index, err := repo.Index()
	if err != nil {
		log.Fatalf("Failed to get repository index: %v", err)
	}
	defer index.Free()

	// Write an empty tree from the index
	treeID, err := index.WriteTree()
	if err != nil {
		log.Fatalf("Failed to write tree: %v", err)
	}

	// Get the tree object from its ID
	tree, err := repo.LookupTree(treeID)
	if err != nil {
		log.Fatalf("Failed to get tree: %v", err)
	}
	defer tree.Free()

	// Now create a commit with no parents (root commit)
	author := &git.Signature{
		Name:  "Your Name",
		Email: "your.email@example.com",
		When:  time.Now(),
	}
	oid, err := repo.CreateCommit("refs/heads/giticket", author, author, "Initial commit", tree)
	if err != nil {
		log.Fatalf("Failed to create commit: %v", err)
	}

	// Lookup the commit from its OID
	commit, err := repo.LookupCommit(oid)
	if err != nil {
		log.Fatalf("Failed to find commit: %v", err)
	}
	defer commit.Free()

	// Create an orphan branch called giticket
	_, err = repo.CreateBranch("giticket", commit, true)
	if err != nil {
		log.Fatalf("Failed to create branch: %v", err)
	}
}

// Help() prints help for this action
func (action *ActionInit) Help() {
	fmt.Println("  init - Initialize giticket")
	fmt.Println("    eg: giticket init")
}
