package actions

import (
	"fmt"
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
		panic(err)
	}

	blobOid, err := repo.CreateBlobFromBuffer([]byte("0"))
	if err != nil {
		panic(err)
	}

	treeBuilderRoot, err := repo.TreeBuilder()
	if err != nil {
		panic(err)
	}
	treeBuilderGiticket, err := repo.TreeBuilder()
	if err != nil {
		panic(err)
	}
	defer treeBuilderRoot.Free()
	defer treeBuilderGiticket.Free()

	err = treeBuilderGiticket.Insert("next_ticket_id", blobOid, git.FilemodeBlob)
	if err != nil {
		panic(err)
	}

	giticketTreeID, err := treeBuilderGiticket.Write()
	if err != nil {
		panic(err)
	}

	err = treeBuilderRoot.Insert(".giticket", giticketTreeID, git.FilemodeTree)
	if err != nil {
		panic(err)
	}

	// Write the tree to the repository
	treeOid, err := treeBuilderRoot.Write()
	if err != nil {
		panic(err)
	}

	tree, err := repo.LookupTree(treeOid)
	if err != nil {
		panic(err)
	}

	// Load the configuration which merges global, system, and local configs
	cfg, err := repo.Config()
	if err != nil {
		fmt.Println("Error accessing config:", err)
		panic(err)
	}
	defer cfg.Free()

	// Retrieve user's name and email from the configuration
	name, err := cfg.LookupString("user.name")
	if err != nil {
		fmt.Println("Error retrieving user name:", err)
		panic(err)
	}

	email, err := cfg.LookupString("user.email")
	if err != nil {
		fmt.Println("Error retrieving user email:", err)
		panic(err)
	}

	// Create a new commit on the branch
	author := &git.Signature{
		Name:  name,
		Email: email,
		When:  time.Now(),
	}

	oid, err := repo.CreateCommit("refs/heads/giticket", author, author, "Initial commit", tree)
	if err != nil {
		panic(err)
	}

	// Lookup the commit from its OID
	commit, err := repo.LookupCommit(oid)
	if err != nil {
		panic(err)
	}
	defer commit.Free()

	// Create an orphan branch called giticket
	_, err = repo.CreateBranch("giticket", commit, true)
	if err != nil {
		panic(err)
	}
}

// Help() prints help for this action
func (action *ActionInit) Help() {
	fmt.Println("  init - Initialize giticket")
	fmt.Println("    eg: giticket init")
}

// InitFlags()
func (action *ActionInit) InitFlags() {}
