package subcommands

import (
	"flag"
	"fmt"
	"strings"
	"time"

	"github.com/jeffWelling/giticket/pkg/debug"
	git "github.com/jeffwelling/git2go/v37"
)

// init() is used to register this action
func init() {
	subcommand := new(SubcommandInit)
	registerSubcommand("init", subcommand)
}

type SubcommandInit struct {
	debugFlag bool
	helpFlag  bool
	flagset   *flag.FlagSet
}

// Execute() is called when this action is invoked
func (subcommand *SubcommandInit) Execute() {

	// Open an existing repository in the current directory
	debug.DebugMessage(subcommand.debugFlag, "Opening git repository '.'")
	repo, err := git.OpenRepository(".")
	if err != nil {
		panic(err)
	}

	// Set the first ticket ID to 1
	// Save "1" as a blob for next_ticket_id
	debug.DebugMessage(subcommand.debugFlag, "Setting first ticket ID to 1")
	blobOid, err := repo.CreateBlobFromBuffer([]byte("1"))
	if err != nil {
		panic(err)
	}

	// This is a root commit, and we're adding a directory with a file in it,
	// so we need two treeBuilders. One for the root tree of the commit, and
	// one for the directory
	debug.DebugMessage(subcommand.debugFlag, "Creating root tree builder")
	treeBuilderRoot, err := repo.TreeBuilder()
	if err != nil {
		panic(err)
	}
	debug.DebugMessage(subcommand.debugFlag, "Creating giticket tree builder")
	treeBuilderGiticket, err := repo.TreeBuilder()
	if err != nil {
		panic(err)
	}
	defer treeBuilderRoot.Free()
	defer treeBuilderGiticket.Free()

	// Create a file named next_ticket_id under the directory we will create
	debug.DebugMessage(subcommand.debugFlag, "Creating file named next_ticket_id")
	err = treeBuilderGiticket.Insert("next_ticket_id", blobOid, git.FilemodeBlob)
	if err != nil {
		panic(err)
	}

	// Write the tree for the directory and get the tree id
	debug.DebugMessage(subcommand.debugFlag, "Writing giticket tree")
	giticketTreeID, err := treeBuilderGiticket.Write()
	if err != nil {
		panic(err)
	}

	// Add the tree ID for the directory named ".giticket" to the root tree
	// builder
	debug.DebugMessage(subcommand.debugFlag, "Adding giticket tree ID to root tree")
	err = treeBuilderRoot.Insert(".giticket", giticketTreeID, git.FilemodeTree)
	if err != nil {
		panic(err)
	}

	// Write the root tree to the repository
	debug.DebugMessage(subcommand.debugFlag, "Writing root tree")
	treeOid, err := treeBuilderRoot.Write()
	if err != nil {
		panic(err)
	}

	// Lookup the tree ID to get the tree we just created
	debug.DebugMessage(subcommand.debugFlag, "Lookup tree ID")
	tree, err := repo.LookupTree(treeOid)
	if err != nil {
		panic(err)
	}

	// Load the configuration which merges global, system, and local configs
	debug.DebugMessage(subcommand.debugFlag, "Loading config")
	cfg, err := repo.Config()
	if err != nil {
		fmt.Println("Error accessing config:", err)
		panic(err)
	}
	defer cfg.Free()

	// Retrieve user's name and email from the configuration
	debug.DebugMessage(subcommand.debugFlag, "Retrieving user name and email")
	name, err := cfg.LookupString("user.name")
	if err != nil {
		fmt.Println("Error retrieving user name:", err)
		panic(err)
	}
	debug.DebugMessage(subcommand.debugFlag, "Retrieving user email")
	email, err := cfg.LookupString("user.email")
	if err != nil {
		fmt.Println("Error retrieving user email:", err)
		panic(err)
	}

	// Create a new commit on the branch
	debug.DebugMessage(subcommand.debugFlag, "Creating commit")
	author := &git.Signature{
		Name:  name,
		Email: email,
		When:  time.Now(),
	}

	// Raise shields, weapons to maximum!
	debug.DebugMessage(subcommand.debugFlag, "Committing")
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
	debug.DebugMessage(subcommand.debugFlag, "Lookup commit")
	commit, err := repo.LookupCommit(oid)
	if err != nil {
		panic(err)
	}
	defer commit.Free()

	// Create branch called giticket pointing to this commit
	debug.DebugMessage(subcommand.debugFlag, "Creating branch")
	_, err = repo.CreateBranch("giticket", commit, true)
	if err != nil {
		panic(err)
	}
}

// Help() prints help for this action
func (subcommand *SubcommandInit) Help() {
	fmt.Println("  init - Initialize giticket")
	fmt.Println("    eg: giticket init")
	fmt.Println("    parameters:")
	fmt.Println("      --debug")
	fmt.Println("      --help")
	fmt.Println("    examples:")
	fmt.Println("      - name: Initialize giticket")
	fmt.Println("        example: giticket init")
}

// InitFlags()
func (subcommand *SubcommandInit) InitFlags(args []string) error {
	subcommand.flagset = flag.NewFlagSet("init", flag.ExitOnError)
	subcommand.flagset.BoolVar(&subcommand.debugFlag, "debug", false, "Print debug info")
	subcommand.flagset.BoolVar(&subcommand.helpFlag, "help", false, "Print help")
	subcommand.flagset.Parse(args)

	// If help
	if subcommand.helpFlag {
		subcommand.Help()
		return nil
	}

	return nil
}
