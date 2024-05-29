package subcommands

import (
	"encoding/json"
	"flag"
	"fmt"
	"strconv"
	"strings"
	"time"

	git "github.com/jeffwelling/git2go/v37"
	"github.com/jeffwelling/giticket/pkg/common"
	"github.com/jeffwelling/giticket/pkg/debug"
	"github.com/jeffwelling/giticket/pkg/repo"
	"github.com/jeffwelling/giticket/pkg/ticket"
)

// init() is used to register this subcommand
func init() {
	subcommand := new(SubcommandCreate)
	registerSubcommand("create", subcommand)
}

type SubcommandCreate struct {
	title           string
	description     string
	labels          []string
	priority        int
	severity        int
	status          string
	id              int
	comments        []ticket.Comment
	debug           bool
	flagset         *flag.FlagSet
	next_comment_id int
	helpFlag        bool
}

// Execute() is called when this subcommand is invoked
func (subcommand *SubcommandCreate) Execute() {
	branchName := "giticket"

	debug.DebugMessage(subcommand.debug, "Opening git repository")
	repo, err := git.OpenRepository(".")
	if err != nil {
		panic(err)
	}

	_, filename := createTicketAndDirectories(repo, branchName, subcommand)
	fmt.Println("Ticket created: ", filename)
	return
}

// Help() prints help for this subcommand
func (subcommand *SubcommandCreate) Help() {
	fmt.Println("  create - Create a new ticket")
	fmt.Println("    eg: giticket create [parameters]")
	fmt.Println("    parameters:")
	fmt.Println("      --title       | -t \"Ticket Title\"")
	fmt.Println("      --description | -d \"Ticket Description\"")
	fmt.Println("      --priority    | -p 1")
	fmt.Println("      --severity    | -sev 1")
	fmt.Println("      --status      | -s \"new\"")
	fmt.Println("      --comments '[{\"Body\":\"My comment\", \"Author\": \"John Smith <smith@example.com>\", \"Created\": 1816534799}]'")
	fmt.Println("      --labels \"my first tag,tag2,tag3\"")
	fmt.Println("      --debug")
	fmt.Println("      --help")
	fmt.Println("    examples:")
	fmt.Println("      - name: Create a new ticket with title \"Ticket Title\" and description \"Ticket Description\"")
	fmt.Println("        example: giticket create --title \"Ticket Title\" --description \"Ticket Description\"")
	fmt.Println("      - name: Create a new ticket with title \"Ticket Title\" and description \"Ticket Description\" and priority 1")
	fmt.Println("        example: giticket create --title \"Ticket Title\" --description \"Ticket Description\" --priority 1")
	fmt.Println("      - name: Create a new ticket with title \"Ticket Title\" and the label \"first tag\"")
	fmt.Println("        example: giticket create --title \"Ticket Title\" --labels \"first tag\"")
	fmt.Println("      - name: Create a new ticket with title \"Ticket Title\" and the label \"first label\" and \"second label\"")
	fmt.Println("        example: giticket create --title \"Ticket Title\" --labels \"first label,second label\"")
	fmt.Println("      - name: Create a new ticket with title \"Ticket Title\" and a single comment")
	fmt.Println("        example: giticket create --title \"Ticket Title\" --comments '[{\"Body\":\"My comment\", \"Author\": \"John Smith <smith@example.com>\"}]'")
	fmt.Println("      - name: Create a new ticket with title \"Ticket Title\" and two comments with one Created date set manually")
	fmt.Println("        example: giticket create --title \"Ticket Title\" --comments '[{\"Body\":\"My comment\", \"Author\": \"John Smith <smith@example.com>\"}, {\"Body\":\"My second comment\", \"Author\": \"John Smith <smith@example.com>\", \"Created\": 1816534799}]'")
}

// InitFlags()
func (subcommand *SubcommandCreate) InitFlags(args []string) error {
	subcommand.flagset = flag.NewFlagSet("create", flag.ExitOnError)

	subcommand.flagset.BoolVar(&subcommand.debug, "debug", false, "Print debug info")

	subcommand.flagset.StringVar(&subcommand.title, "title", "", "Title for the new ticket")
	subcommand.flagset.StringVar(&subcommand.title, "t", "", "Title for the new ticket")
	subcommand.flagset.StringVar(&subcommand.description, "description", "", "Description of the ticket to create")
	subcommand.flagset.StringVar(&subcommand.description, "d", "", "Description of the ticket to create")
	labelsFlag := subcommand.flagset.String("labels", "", "Comma separated list of labels to apply to the ticket")
	subcommand.flagset.IntVar(&subcommand.priority, "priority", 1, "Priority of the ticket")
	subcommand.flagset.IntVar(&subcommand.priority, "p", 1, "Priority of the ticket")
	subcommand.flagset.IntVar(&subcommand.severity, "severity", 1, "Severity of the ticket")
	subcommand.flagset.IntVar(&subcommand.severity, "sev", 1, "Severity of the ticket")
	subcommand.flagset.StringVar(&subcommand.status, "status", "new", "Status of the ticket")
	subcommand.flagset.StringVar(&subcommand.status, "s", "new", "Status of the ticket")
	subcommand.flagset.BoolVar(&subcommand.helpFlag, "help", false, "Print help for the create subcommand")
	commentsFlag := subcommand.flagset.String("comments", "", "Comma separated list of comments to add to the ticket")

	subcommand.flagset.Parse(args)

	if subcommand.helpFlag {
		common.PrintVersion()
		fmt.Println("giticket")
		subcommand.Help()
	}

	// get the author
	debug.DebugMessage(subcommand.debug, "Opening git repository to get author")
	repo, err := git.OpenRepository(".")
	if err != nil {
		return err
	}
	author := common.GetAuthor(repo)

	// Handle labels separately to split them into a slice
	if *labelsFlag != "" {
		labels := strings.Split(*labelsFlag, ",")
		for i, label := range labels {
			labels[i] = strings.TrimSpace(label)
		}
		subcommand.labels = labels
	}

	// Handle comments separately to parse them into a slice of Comments
	if *commentsFlag != "" {
		// First comment ID starts at 1
		subcommand.next_comment_id = 1
		var comments []ticket.Comment
		err := json.Unmarshal([]byte(*commentsFlag), &comments)
		if err != nil {
			return err
		}
		for i := range comments {
			comments[i].ID = subcommand.next_comment_id
			subcommand.next_comment_id++

			if comments[i].Created == 0 {
				comments[i].Created = time.Now().Unix()
			}
			if comments[i].Author == "" {
				comments[i].Author = author.Name + " <" + author.Email + ">"
			}
		}
		subcommand.comments = comments
	}

	// Check for required parameters
	if subcommand.title == "" {
		ticket.PrintParameterMissing("title")
		return fmt.Errorf("Cannot create a ticket without a title")
	}
	return nil
}

// ticketFilename() returns the ticket title encoded as a filename prefixed with
// the ticket ID
func (subcommand *SubcommandCreate) ticketFilename() string {
	// turn spaces into underscores
	title := strings.ReplaceAll(subcommand.title, " ", "_")
	return fmt.Sprintf("%d_%s", subcommand.id, title)
}

func createTicketAndDirectories(thisRepo *git.Repository, branchName string, subcommand *SubcommandCreate) (*git.Oid, string) {
	parentCommit, err := repo.GetParentCommit(thisRepo, branchName, subcommand.debug)
	if err != nil {
		panic(err)
	}

	// Get value for .giticket/next_ticket_id
	ticketID, err := ticket.ReadNextTicketID(thisRepo, parentCommit)
	if err != nil {
		panic(err)
	}
	debug.DebugMessage(subcommand.debug, "Next ticket ID: "+strconv.Itoa(ticketID))

	// Increment ticketID and write it as a blob
	i := ticketID + 1
	debug.DebugMessage(subcommand.debug, "incrementing next ticket ID in .giticket/next_ticket_id, is now: "+strconv.Itoa(i))
	NTIDBlobOID, err := thisRepo.CreateBlobFromBuffer([]byte(strconv.Itoa(i)))
	debug.DebugMessage(subcommand.debug, "NTIDBlobOID: "+NTIDBlobOID.String())

	rootTreeBuilder, previousCommitTree, err := repo.TreeBuilderFromCommit(parentCommit, thisRepo, subcommand.debug)
	if err != nil {
		panic(err)
	}
	defer rootTreeBuilder.Free()

	giticketTree, err := repo.GetSubTreeByName(previousCommitTree, thisRepo, ".giticket", subcommand.debug)
	defer giticketTree.Free()

	// Create a TreeBuilder from the previous tree for giticket so we can add a
	// ticket under .gititcket/tickets and change the value of
	// .giticket/next_ticket_id
	debug.DebugMessage(subcommand.debug, "creating tree builder for .giticket tree: "+giticketTree.Id().String())
	giticketTreeBuilder, err := thisRepo.TreeBuilderFromTree(giticketTree)
	if err != nil {
		panic(err)
	}
	defer giticketTreeBuilder.Free()

	// Insert the blob for next_ticket_id into the TreeBuilder for giticket
	// This essentially saves the file to .giticket/next_ticket_id, but we then need to
	// save the directory to the parent, all the way up to the commit.
	debug.DebugMessage(subcommand.debug, "inserting next ticket ID into .giticket/next_ticket_id: "+NTIDBlobOID.String())
	err = giticketTreeBuilder.Insert("next_ticket_id", NTIDBlobOID, git.FilemodeBlob)
	if err != nil {
		panic(err)
	}

	var giticketTicketsTreeBuilder *git.TreeBuilder
	_giticketTicketsTreeID := giticketTree.EntryByName("tickets")
	if _giticketTicketsTreeID == nil {
		// Create the tickets directory

		// Get a TreeBuilder for .giticket/tickets so we can add the ticket to that
		// directory
		debug.DebugMessage(subcommand.debug, "creating empty tree builder for .giticket/tickets")
		giticketTicketsTreeBuilder, err = thisRepo.TreeBuilder()
		if err != nil {
			panic(err)
		}
		defer giticketTicketsTreeBuilder.Free()
	} else {
		debug.DebugMessage(subcommand.debug, "looking up tree for .giticket/tickets: "+_giticketTicketsTreeID.Id.String())
		giticketTicketsTree, err := thisRepo.LookupTree(_giticketTicketsTreeID.Id)
		if err != nil {
			panic(err)
		}
		defer giticketTicketsTree.Free()

		debug.DebugMessage(subcommand.debug, "creating tree builder for .giticket/tickets tree: "+giticketTicketsTree.Id().String())
		giticketTicketsTreeBuilder, err = thisRepo.TreeBuilderFromTree(giticketTicketsTree)
		if err != nil {
			panic(err)
		}
	}

	debug.DebugMessage(subcommand.debug, "creating and populating ticket")
	// Craft the ticket
	t := ticket.Ticket{}
	t.Created = time.Now().Unix()
	t.Title = subcommand.title
	t.Description = subcommand.description
	t.Labels = subcommand.labels
	t.Priority = subcommand.priority
	t.Severity = subcommand.severity
	t.Status = subcommand.status
	t.ID = ticketID
	t.Comments = subcommand.comments
	t.NextCommentID = subcommand.next_comment_id
	// FIXME Add a way to parse comments from initial ticket creation

	// Write ticket
	ticketBlobOID, err := thisRepo.CreateBlobFromBuffer(t.TicketToYaml())
	if err != nil {
		panic(err)
	}
	debug.DebugMessage(subcommand.debug, "writing ticket to .giticket/tickets: "+ticketBlobOID.String())

	// Add ticket to .giticket/tickets
	debug.DebugMessage(subcommand.debug, "adding ticket to .giticket/tickets: "+ticketBlobOID.String())
	err = giticketTicketsTreeBuilder.Insert(t.TicketFilename(), ticketBlobOID, git.FilemodeBlob)
	if err != nil {
		panic(err)
	}

	// Save the tree and get the tree ID for .giticket/tickets
	giticketTicketsTreeID, err := giticketTicketsTreeBuilder.Write()
	if err != nil {
		panic(err)
	}
	debug.DebugMessage(subcommand.debug, "saving .giticket/tickets:"+giticketTicketsTreeID.String())

	// Add 'ticket' directory to '.giticket' TreeBuilder, to update the
	// .giticket directory with the new tree for the updated tickets directory
	debug.DebugMessage(subcommand.debug, "adding ticket directory to .giticket: "+giticketTicketsTreeID.String())
	err = giticketTreeBuilder.Insert("tickets", giticketTicketsTreeID, git.FilemodeTree)
	if err != nil {
		panic(err)
	}

	giticketTreeID, err := giticketTreeBuilder.Write()
	if err != nil {
		panic(err)
	}
	debug.DebugMessage(subcommand.debug, "saving .giticket tree to root tree: "+giticketTreeID.String())

	// Update the root tree builder with the new .giticket directory
	debug.DebugMessage(subcommand.debug, "updating root tree with: "+giticketTreeID.String())
	err = rootTreeBuilder.Insert(".giticket", giticketTreeID, git.FilemodeTree)
	if err != nil {
		panic(err)
	}

	// Save the new root tree and get the ID so we can lookup the tree for the
	// commit
	debug.DebugMessage(subcommand.debug, "saving root tree")
	rootTreeBuilderID, err := rootTreeBuilder.Write()
	if err != nil {
		panic(err)
	}

	// Lookup the tree so we can use it in the commit
	debug.DebugMessage(subcommand.debug, "lookup root tree for commit: "+rootTreeBuilderID.String())
	rootTree, err := thisRepo.LookupTree(rootTreeBuilderID)
	if err != nil {
		panic(err)
	}
	defer rootTree.Free()

	// Get author data by reading .git configs
	debug.DebugMessage(subcommand.debug, "getting author data")
	author := common.GetAuthor(thisRepo)

	// commit and update 'giticket' branch
	debug.DebugMessage(subcommand.debug, "creating commit")
	commitID, err := thisRepo.CreateCommit("refs/heads/giticket", author, author, "Creating ticket "+t.TicketFilename(), rootTree, parentCommit)
	if err != nil {
		panic(err)
	}

	return commitID, t.TicketFilename()
}
