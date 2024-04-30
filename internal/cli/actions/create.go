package actions

import (
	"encoding/json"
	"flag"
	"fmt"
	"strings"
	"time"

	"github.com/jeffWelling/giticket/pkg/ticket"
	git "github.com/jeffwelling/git2go/v37"
)

// init() is used to register this action
func init() {
	action := new(ActionCreate)
	registerAction("create", action)
}

type ActionCreate struct {
	title       string
	description string
	labels      []string
	priority    int
	severity    int
	status      string
	id          int
	comments    []ticket.Comment
}

// Execute() is called when this action is invoked
func (action *ActionCreate) Execute() {
	fmt.Println("Creating ticket")

	branchName := "giticket"
	repo, err := git.OpenRepository(".")
	if err != nil {
		panic(err)
	}

	newTicket := ticket.Ticket{}
	newTicket.Created = time.Now().Unix() // Use the current time as the created time for the ticke
	newTicket.Title = action.title
	newTicket.Description = action.description
	newTicket.Labels = action.labels

	filePath := ".giticket/tickets/" + newTicket.TicketFilename()

	fmt.Println("finding branch")
	// Find the branch and its target commit
	branch, err := repo.LookupBranch(branchName, git.BranchLocal)
	if err != nil {
		panic(err)
	}

	fmt.Println("finding commit")
	commit, err := repo.LookupCommit(branch.Target())
	if err != nil {
		panic(err)
	}

	fmt.Println("creating tree builder")
	// Create a tree builder from the commit's tree
	old_tree, err := commit.Tree()
	if err != nil {
		panic(err)
	}

	fmt.Println("tree builder from tree")
	treeBuilder, err := repo.TreeBuilderFromTree(old_tree)
	if err != nil {
		panic(err)
	}
	defer treeBuilder.Free()

	fmt.Println("reading and incrementing ticket id")
	nextTicketID, err := ticket.ReadAndIncrementTicketID(repo, branchName, ".giticket/next_ticket_id", treeBuilder)
	if err != nil {
		panic(err)
	}
	newTicket.ID = nextTicketID

	fmt.Println("creating blob")
	// Create a blob from the content
	blobOid, err := repo.CreateBlobFromBuffer(newTicket.TicketToYaml())
	if err != nil {
		panic(err)
	}

	// Insert the new file into the tree
	err = treeBuilder.Insert(filePath, blobOid, git.FilemodeBlob)
	if err != nil {
		panic(err)
	}

	// Write the tree to the repository
	treeOid, err := treeBuilder.Write()
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
	parents := []*git.Commit{commit}

	objectID, err := repo.CreateCommit("refs/heads/"+branchName, author, author, "Add fubar file", tree, parents...)
	if err != nil {
		panic(err)
	}

	// Print the objectID
	fmt.Printf("Created ticket: %s\n", objectID.String())
}

// Help() prints help for this action
func (action *ActionCreate) Help() {
	fmt.Println("  create - Create a new ticket")
	fmt.Println("    eg: giticket create [parameters]")
	fmt.Println("    parameters:")
	fmt.Println("      -title \"Ticket Title\"")
	fmt.Println("      -description \"Ticket Description\"")
	fmt.Println("      -labels \"my first tag, my second tag, tag3\"")
	fmt.Println("      -priority 1")
	fmt.Println("      -severity 1")
	fmt.Println("      -status \"new\"")
}

// InitFlags()
func (action *ActionCreate) InitFlags() {
	titleFlag := flag.String("title", "", "Title of the ticket to create")
	descriptionFlag := flag.String("description", "", "Description of the ticket to create")
	labelsFlag := flag.String("labels", "", "Comma separated list of labels to apply to the ticket")
	priorityFlag := flag.Int("priority", 1, "Priority of the ticket")
	severityFlag := flag.Int("severity", 1, "Severity of the ticket")
	statusFlag := flag.String("status", "new", "Status of the ticket")
	commentsFlag := flag.String("comments", "", "Comma separated list of comments to add to the ticket")

	// If comments are provided, parse them.
	if *commentsFlag != "" {
		var comments []ticket.Comment
		err := json.Unmarshal([]byte(*commentsFlag), &comments)
		if err != nil {
			panic(err)
		}
		action.comments = comments
	}

	flag.Parse()

	// Check for required parameters
	if *titleFlag == "" {
		ticket.PrintParameterMissing("title")
		return
	}

	action.title = *titleFlag
	action.description = *descriptionFlag
	if *labelsFlag != "" {
		labels := strings.Split(*labelsFlag, ",")
		for i, label := range labels {
			labels[i] = strings.TrimSpace(label)
		}
		action.labels = labels
	}
	action.priority = *priorityFlag
	action.severity = *severityFlag
	action.status = *statusFlag
}

// ticketFilename() returns the ticket title encoded as a filename prefixed with
// the ticket ID
func (action *ActionCreate) ticketFilename() string {
	// turn spaces into underscores
	title := strings.ReplaceAll(action.title, " ", "_")
	return fmt.Sprintf("%d_%s", action.id, title)
}
