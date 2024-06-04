package subcommands

import (
	"encoding/json"
	"flag"
	"fmt"
	"strings"
	"time"

	git "github.com/jeffwelling/git2go/v37"
	"github.com/jeffwelling/giticket/pkg/common"
	"github.com/jeffwelling/giticket/pkg/debug"
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
	debugFlag       bool
	flagset         *flag.FlagSet
	next_comment_id int
	helpFlag        bool
	params          map[string]interface{}
}

// Execute() is called when this subcommand is invoked
func (subcommand *SubcommandCreate) Execute() {
	_, filename := ticket.HandleCreate(
		common.BranchName, time.Now().Unix(),
		subcommand.title, subcommand.description,
		subcommand.labels, subcommand.priority,
		subcommand.severity, subcommand.status,
		subcommand.comments, subcommand.next_comment_id,
		subcommand.debugFlag,
	)
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
	var (
		helpFlag  bool
		ticketID  int
		commentID int
		comment   string
		delete    bool
		debugFlag bool
	)

	subcommand.flagset = flag.NewFlagSet("create", flag.ExitOnError)

	subcommand.flagset.BoolVar(&subcommand.debugFlag, "debug", false, "Print debug info")

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

	subcommand.params["helpFlag"] = helpFlag
	subcommand.params["ticketID"] = ticketID
	subcommand.params["commentID"] = commentID
	subcommand.params["comment"] = comment
	subcommand.params["delete"] = delete
	subcommand.params["debugFlag"] = debugFlag
	subcommand.params["labelsFlag"] = labelsFlag
	subcommand.params["commentsFlag"] = commentsFlag

	if subcommand.helpFlag {
		common.PrintVersion()
		fmt.Println("giticket")
		subcommand.Help()
	}

	// get the author
	debug.DebugMessage(subcommand.debugFlag, "Opening git repository to get author")
	repo, err := git.OpenRepository(".")
	if err != nil {
		return err
	}
	author, err := common.GetAuthor(repo)
	if err != nil {
		return err
	}

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

func (subcommand *SubcommandCreate) DebugFlag() bool {
	return subcommand.Parameters()["debug"].(bool)
}

func (subcommand *SubcommandCreate) Parameters() map[string]interface{} {
	return subcommand.params
}
