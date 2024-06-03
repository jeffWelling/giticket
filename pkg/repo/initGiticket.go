package repo

import (
	"fmt"
	"strings"
	"time"

	git "github.com/jeffwelling/git2go/v37"
	"github.com/jeffwelling/giticket/pkg/debug"
)

func HandleInitGiticket(debugFlag bool) {
	// Open an existing repository in the current directory
	debug.DebugMessage(debugFlag, "Opening git repository '.'")
	repo, err := git.OpenRepository(".")
	if err != nil {
		panic(err)
	}

	// Set the first ticket ID to 1
	// Save "1" as a blob for next_ticket_id
	debug.DebugMessage(debugFlag, "Setting first ticket ID to 1")
	blobOid, err := repo.CreateBlobFromBuffer([]byte("1"))
	if err != nil {
		panic(err)
	}

	// This is a root commit, and we're adding a directory with a file in it,
	// so we need two treeBuilders. One for the root tree of the commit, and
	// one for the directory
	debug.DebugMessage(debugFlag, "Creating root tree builder")
	treeBuilderRoot, err := repo.TreeBuilder()
	if err != nil {
		panic(err)
	}
	debug.DebugMessage(debugFlag, "Creating giticket tree builder")
	treeBuilderGiticket, err := repo.TreeBuilder()
	if err != nil {
		panic(err)
	}
	defer treeBuilderRoot.Free()
	defer treeBuilderGiticket.Free()

	// Create a file named next_ticket_id under the directory we will create
	debug.DebugMessage(debugFlag, "Creating file named next_ticket_id")
	err = treeBuilderGiticket.Insert("next_ticket_id", blobOid, git.FilemodeBlob)
	if err != nil {
		panic(err)
	}

	// Write the tree for the directory and get the tree id
	debug.DebugMessage(debugFlag, "Writing giticket tree")
	giticketTreeID, err := treeBuilderGiticket.Write()
	if err != nil {
		panic(err)
	}

	// Add the tree ID for the directory named ".giticket" to the root tree
	// builder
	debug.DebugMessage(debugFlag, "Adding giticket tree ID to root tree")
	err = treeBuilderRoot.Insert(".giticket", giticketTreeID, git.FilemodeTree)
	if err != nil {
		panic(err)
	}

	// Write the root tree to the repository
	debug.DebugMessage(debugFlag, "Writing root tree")
	treeOid, err := treeBuilderRoot.Write()
	if err != nil {
		panic(err)
	}

	// Lookup the tree ID to get the tree we just created
	debug.DebugMessage(debugFlag, "Lookup tree ID")
	tree, err := repo.LookupTree(treeOid)
	if err != nil {
		panic(err)
	}

	// Load the configuration which merges global, system, and local configs
	debug.DebugMessage(debugFlag, "Loading config")
	cfg, err := repo.Config()
	if err != nil {
		fmt.Println("Error accessing config:", err)
		panic(err)
	}
	defer cfg.Free()

	// Retrieve user's name and email from the configuration
	debug.DebugMessage(debugFlag, "Retrieving user name and email")
	name, err := cfg.LookupString("user.name")
	if err != nil {
		fmt.Println("Error retrieving user name:", err)
		panic(err)
	}
	debug.DebugMessage(debugFlag, "Retrieving user email")
	email, err := cfg.LookupString("user.email")
	if err != nil {
		fmt.Println("Error retrieving user email:", err)
		panic(err)
	}

	// Create a new commit on the branch
	debug.DebugMessage(debugFlag, "Creating commit")
	author := &git.Signature{
		Name:  name,
		Email: email,
		When:  time.Now(),
	}

	// Raise shields, weapons to maximum!
	debug.DebugMessage(debugFlag, "Committing")
	oid, err := repo.CreateCommit("refs/heads/giticket", author, author, "Initial commit", tree)
	if err != nil {
		// If text of error includes the string "current tip is not the first parent" then
		// return "fubar" error
		if strings.Contains(err.Error(), "current tip is not the first parent") {
			fmt.Println("giticket already initialized")
			return
		} else {
			panic(err)
		}
	}

	// Lookup the commit from its OID, to set the branch
	debug.DebugMessage(debugFlag, "Lookup commit")
	commit, err := repo.LookupCommit(oid)
	if err != nil {
		panic(err)
	}
	defer commit.Free()

	// Create branch called giticket pointing to this commit
	debug.DebugMessage(debugFlag, "Creating branch")
	_, err = repo.CreateBranch("giticket", commit, true)
	if err != nil {
		panic(err)
	}
}
