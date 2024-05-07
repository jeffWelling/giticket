package ticket

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	git "github.com/jeffwelling/git2go/v37"
	"gopkg.in/yaml.v2"
)

type Ticket struct {
	Title         string
	Description   string
	Labels        []string
	Priority      int
	Severity      int
	Status        string
	Comments      []Comment
	nextCommentID int

	// Set automatically
	ID      int
	Created int64
}

type Comment struct {
	ID      string
	Created int64
	Body    string
	Author  string
}

// Return the value of ".giticket/next_ticket_id" from the given commit as an
// int, or returns 0 and an error. Make sure to write the incremented value back
// things to ".giticket/next_ticket_id" in the same commit. Repo is required to
// lookup treeIDs
func ReadNextTicketID(repo *git.Repository, commit *git.Commit) (int, error) {

	tree, err := commit.Tree()
	if err != nil {
		return 0, err
	}
	defer tree.Free()

	giticketTreeEntry, err := tree.EntryByPath(".giticket")
	if err != nil {
		return 0, err
	}

	// Lookup giticketTreeEntry
	giticketTree, err := repo.LookupTree(giticketTreeEntry.Id)
	if err != nil {
		return 0, err
	}
	defer giticketTree.Free()

	NTIDEntry, err := giticketTree.EntryByPath("next_ticket_id")
	if err != nil {
		return 0, err
	}

	NTIDBlob, err := repo.LookupBlob(NTIDEntry.Id)
	if err != nil {
		return 0, err
	}
	defer NTIDBlob.Free()

	// read value of blob as int
	s := strings.TrimSpace(string(NTIDBlob.Contents()))

	// Convert string s into int i
	i, err := strconv.Atoi(s)
	if err != nil {
		return 0, err
	}

	return i, nil
}

// TicketFilename() returns the filename of the ticket by cating the ticket
// ID and the ticket title, with spaces replaced by underscores
func (t *Ticket) TicketFilename() string {
	// turn spaces into underscores
	title := strings.ReplaceAll(t.Title, " ", "_")
	return fmt.Sprintf("%d__%s", t.ID, title)
}

// TicketToYaml() returns the ticket as a YAML string
// which is used to save to disk
func (t *Ticket) TicketToYaml() []byte {
	// Turn the ticket into a yaml string
	yamlTicket, err := yaml.Marshal(t)
	if err != nil {
		panic(err)
	}
	return yamlTicket
}

func PrintParameterMissing(param string) {
	fmt.Printf("A required parameter was not provided, check the '--help' output for the action for more details. Missing parameter: %s\n", param)
}

func GetAuthor(repo *git.Repository) *git.Signature {
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

	return author
}
