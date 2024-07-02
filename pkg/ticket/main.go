// Ticket implemtns routines for working with Giticket Tickets
//
// The ticket package uses the git2go library to work with git repositories.
package ticket

import (
	"fmt"
	"strconv"
	"strings"

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
	NextCommentID int `yaml:"next_comment_id" json:"next_comment_id"`

	// Set automatically
	ID      int
	Created int64
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
