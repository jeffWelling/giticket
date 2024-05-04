package actions

import (
	"encoding/json"
	"flag"
	"fmt"
	"strconv"
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
	debug       bool
}

// Execute() is called when this action is invoked
func (action *ActionCreate) Execute() {
	branchName := "giticket"

	if action.debug {
		fmt.Println("Opening repository '.'")
	}
	repo, err := git.OpenRepository(".")
	if err != nil {
		panic(err)
	}

	if action.debug {
		fmt.Println("looking up branch: ", branchName)
	}
	// Find the branch and its target commit
	branch, err := repo.LookupBranch(branchName, git.BranchLocal)
	if err != nil {
		panic(err)
	}

	if action.debug {
		fmt.Println("looking up commit: ", branch.Target())
	}
	commit, err := repo.LookupCommit(branch.Target())
	if err != nil {
		panic(err)
	}

	if action.debug {
		fmt.Println("looking up tree from commit")
	}
	// Create a tree builder from the commit's tree
	old_tree, err := commit.Tree()
	if err != nil {
		panic(err)
	}

	if action.debug {
		fmt.Println("looking up .giticket tree entry from parent commit tree")
	}
	giticketTreeEntry := old_tree.EntryByName(".giticket")

	if action.debug {
		fmt.Println("looking up tree for .giticket from parent commit tree")
	}
	giticketTree, err := repo.LookupTree(giticketTreeEntry.Id)
	if err != nil {
		panic(err)
	}

	if action.debug {
		fmt.Println("looking up tree entry for tickets from .giticket tree")
	}
	giticketTicketsTreeID := giticketTree.EntryByName("tickets")

	if giticketTicketsTreeID == nil {
		if action.debug {
			fmt.Println("giticket tickets tree entry not found, needs creating")
		}
		_, filename := createTicketAndDirectories(repo, branchName, action)
		fmt.Println("Ticket created: ", filename)
	} else {
		if action.debug {
			fmt.Println("giticket tickets tree entry found")
		}
		_, filename := createTicket(repo, branchName, action)
		fmt.Println("Ticket created: ", filename)
	}

	return
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
	fmt.Println("      -debug")
}

// InitFlags()
func (action *ActionCreate) InitFlags() {
	flag.BoolVar(&action.debug, "debug", false, "Print debug info")
	flag.StringVar(&action.title, "title", "", "Title of the ticket to create")
	flag.StringVar(&action.description, "description", "", "Description of the ticket to create")
	labelsFlag := flag.String("labels", "", "Comma separated list of labels to apply to the ticket")
	flag.IntVar(&action.priority, "priority", 1, "Priority of the ticket")
	flag.IntVar(&action.severity, "severity", 1, "Severity of the ticket")
	flag.StringVar(&action.status, "status", "new", "Status of the ticket")
	commentsFlag := flag.String("comments", "", "Comma separated list of comments to add to the ticket")

	flag.Parse()

	// Handle labels separately to split them into a slice
	if *labelsFlag != "" {
		labels := strings.Split(*labelsFlag, ",")
		for i, label := range labels {
			labels[i] = strings.TrimSpace(label)
		}
		action.labels = labels
	}

	// Handle comments separately to parse them into a slice of Comments
	if *commentsFlag != "" {
		var comments []ticket.Comment
		err := json.Unmarshal([]byte(*commentsFlag), &comments)
		if err != nil {
			panic(err)
		}
		action.comments = comments
	}

	// Check for required parameters
	if action.title == "" {
		ticket.PrintParameterMissing("title")
		return
	}
}

// ticketFilename() returns the ticket title encoded as a filename prefixed with
// the ticket ID
func (action *ActionCreate) ticketFilename() string {
	// turn spaces into underscores
	title := strings.ReplaceAll(action.title, " ", "_")
	return fmt.Sprintf("%d_%s", action.id, title)
}

func createTicketAndDirectories(repo *git.Repository, branchName string, action *ActionCreate) (*git.Oid, string) {
	// Find the branch and its target commit
	if action.debug {
		fmt.Println("looking up branch: ", branchName)
	}
	branch, err := repo.LookupBranch(branchName, git.BranchLocal)
	if err != nil {
		panic(err)
	}

	// Lookup the commit the branch references
	if action.debug {
		fmt.Println("looking up commit: ", branch.Target())
	}
	parentCommit, err := repo.LookupCommit(branch.Target())
	if err != nil {
		panic(err)
	}

	// Get value for .giticket/next_ticket_id
	if action.debug {
		fmt.Println("reading next ticket ID from .giticket/next_ticket_id")
	}
	ticketID, err := ticket.ReadNextTicketID(repo, parentCommit)
	if err != nil {
		panic(err)
	}

	// Increment ticketID and write it as a blob
	if action.debug {
		fmt.Println("incrementing next ticket ID in .giticket/next_ticket_id")
	}
	i := ticketID + 1
	NTIDBlobOID, err := repo.CreateBlobFromBuffer([]byte(strconv.Itoa(i)))

	// Lookup the tree from the previous commit, we need this to build a tree
	// that includes the files from the previous commit.
	if action.debug {
		fmt.Println("looking up tree from parent commit")
	}
	previousCommitTree, err := parentCommit.Tree()
	if err != nil {
		panic(err)
	}
	defer previousCommitTree.Free()

	// Create a TreeBuilder from the previous commit's tree, so that we can
	// update it with our changes to the .giticket directory
	rootTreeBuilder, err := repo.TreeBuilderFromTree(previousCommitTree)
	if err != nil {
		panic(err)
	}
	defer rootTreeBuilder.Free()

	// Get the TreeEntry for ".giticket" from the previous commit so we can get
	// the tree for .giticket
	if action.debug {
		fmt.Println("looking up tree entry for .giticket")
	}
	giticketTreeEntry := previousCommitTree.EntryByName(".giticket")

	// Lookup tree for giticket
	if action.debug {
		fmt.Println("looking up tree for .giticket")
	}
	giticketTree, err := repo.LookupTree(giticketTreeEntry.Id)
	if err != nil {
		panic(err)
	}
	defer giticketTree.Free()

	// Create a TreeBuilder from the previous tree for giticket so we can add a
	// ticket under .gititcket/tickets and change the value of
	// .giticket/next_ticket_id
	if action.debug {
		fmt.Println("creating tree builder for .giticket")
	}
	giticketTreeBuilder, err := repo.TreeBuilderFromTree(giticketTree)
	if err != nil {
		panic(err)
	}
	defer giticketTreeBuilder.Free()

	// Insert the blob for next_ticket_id into the TreeBuilder for giticket
	// This essentially saves the file to .giticket/next_ticket_id, but we then need to
	// save the directory to the parent, all the way up to the commit.
	if action.debug {
		fmt.Println("inserting next ticket ID into .giticket/next_ticket_id")
	}
	err = giticketTreeBuilder.Insert("next_ticket_id", NTIDBlobOID, git.FilemodeBlob)
	if err != nil {
		panic(err)
	}

	// Get a TreeBuilder for .giticket/tickets so we can add the ticket to that
	// directory
	if action.debug {
		fmt.Println("creating empty tree builder for .giticket/tickets")
	}
	giticketTicketsTreeBuilder, err := repo.TreeBuilder()
	if err != nil {
		panic(err)
	}
	defer giticketTicketsTreeBuilder.Free()

	if action.debug {
		fmt.Println("creating and populating ticket")
	}
	// Craft the ticket
	t := ticket.Ticket{}
	t.Created = time.Now().Unix()
	t.Title = action.title
	t.Description = action.description
	t.Labels = action.labels
	t.Priority = action.priority
	t.Severity = action.severity
	t.Status = action.status
	t.ID = ticketID
	// FIXME Add a way to parse comments from initial ticket creation

	// Write ticket
	if action.debug {
		fmt.Println("writing ticket")
	}
	ticketBlobOID, err := repo.CreateBlobFromBuffer(t.TicketToYaml())
	if err != nil {
		panic(err)
	}

	// Add ticket to .giticket/tickets
	if action.debug {
		fmt.Println("adding ticket to .giticket/tickets")
	}
	err = giticketTicketsTreeBuilder.Insert(t.TicketFilename(), ticketBlobOID, git.FilemodeBlob)
	if err != nil {
		panic(err)
	}

	// Save the tree and get the tree ID for .giticket/tickets
	if action.debug {
		fmt.Println("saving .giticket/tickets")
	}
	giticketTicketsTreeID, err := giticketTicketsTreeBuilder.Write()
	if err != nil {
		panic(err)
	}

	// Add 'ticket' directory to '.giticket' TreeBuilder, to update the
	// .giticket directory with the new tree for the updated tickets directory
	if action.debug {
		fmt.Println("adding ticket directory to .giticket")
	}
	err = giticketTreeBuilder.Insert("tickets", giticketTicketsTreeID, git.FilemodeTree)
	if err != nil {
		panic(err)
	}

	// Save it so we can update the root tree with the new .giticket directory
	if action.debug {
		fmt.Println("creating tree builder from giticket tree")
	}
	giticketTreeBuilder, err = repo.TreeBuilderFromTree(giticketTree)
	if err != nil {
		panic(err)
	}
	defer giticketTreeBuilder.Free()

	if action.debug {
		fmt.Println("saving .giticket tree to root tree")
	}
	giticketTreeID, err := giticketTreeBuilder.Write()
	if err != nil {
		panic(err)
	}

	// Update the root tree builder with the new .giticket directory
	if action.debug {
		fmt.Println("updating root tree")
	}
	err = rootTreeBuilder.Insert(".giticket", giticketTreeID, git.FilemodeTree)
	if err != nil {
		panic(err)
	}

	// Save the new root tree and get the ID so we can lookup the tree for the
	// commit
	if action.debug {
		fmt.Println("saving root tree")
	}
	rootTreeBuilderID, err := rootTreeBuilder.Write()
	if err != nil {
		panic(err)
	}

	// Lookup the tree so we can use it in the commit
	if action.debug {
		fmt.Println("lookup root tree for commit")
	}
	rootTree, err := repo.LookupTree(rootTreeBuilderID)
	if err != nil {
		panic(err)
	}
	defer rootTree.Free()

	// Get author data by reading .git configs
	if action.debug {
		fmt.Println("getting author data")
	}
	author := ticket.GetAuthor(repo)

	// commit and update 'giticket' branch
	if action.debug {
		fmt.Println("creating commit")
	}
	commitID, err := repo.CreateCommit("refs/heads/giticket", author, author, "Creating ticket "+t.TicketFilename(), rootTree, parentCommit)
	if err != nil {
		panic(err)
	}

	return commitID, t.TicketFilename()
}

func createTicket(repo *git.Repository, branchName string, action *ActionCreate) (*git.Commit, string) {
	return nil, ""
}
