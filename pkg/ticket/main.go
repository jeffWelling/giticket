package ticket

import (
	"fmt"
	"strconv"
	"strings"

	git "github.com/jeffwelling/git2go/v37"
	"gopkg.in/yaml.v2"
)

type Ticket struct {
	Title       string
	Description string
	Labels      []string
	Priority    int
	Severity    int
	Status      string
	Comments    []Comment

	// Set automatically
	ID      int
	Created int64
}

type Comment struct {
	Id      string
	Created int64
	Body    string
	Author  string
}

func ReadAndIncrementTicketID(repo *git.Repository, branchName, filePath string, treeBuilder *git.TreeBuilder) (int, error) {
	branch, err := repo.LookupBranch(branchName, git.BranchLocal)
	if err != nil {
		return 0, err
	}

	commit, err := repo.LookupCommit(branch.Target())
	if err != nil {
		return 0, err
	}

	tree, err := commit.Tree()
	if err != nil {
		return 0, err
	}

	entry, err := tree.EntryByPath(filePath)
	if err != nil {
		return 0, err
	}

	blob, err := repo.LookupBlob(entry.Id)
	if err != nil {
		return 0, err
	}

	s := strings.TrimSpace(string(blob.Contents()))

	// Convert string s into int i
	i, err := strconv.Atoi(s)
	if err != nil {
		return 0, err
	}

	// Increment i
	i++

	blobOid, err := repo.CreateBlobFromBuffer([]byte(strconv.Itoa(i)))
	if err != nil {
		return 0, err
	}

	err = treeBuilder.Insert(filePath, blobOid, git.FilemodeBlob)
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
