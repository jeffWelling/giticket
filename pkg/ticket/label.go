package ticket

import (
	git "github.com/jeffwelling/git2go/v37"
	"github.com/jeffwelling/giticket/pkg/common"
	"github.com/jeffwelling/giticket/pkg/debug"
	"github.com/jeffwelling/giticket/pkg/repo"
)

// HandleLabel handles adding or deleting a label from a ticket based on the
// value of deleteFlag, returning an error if there was one.
func HandleLabel(
	branchName string,
	label string,
	deleteFlag bool,
	ticketID int,
	debugFlag bool,
) error {
	debug.DebugMessage(debugFlag, "Opening git repository")
	thisRepo, err := git.OpenRepository(".")
	if err != nil {
		return err
	}

	// Get author
	author, err := common.GetAuthor(thisRepo)
	if err != nil {
		return err
	}

	tickets, err := GetListOfTickets(thisRepo, branchName, debugFlag)
	if err != nil {
		return err
	}
	t := FilterTicketsByID(tickets, ticketID)

	if deleteFlag {
		labelID := DeleteLabel(&t, label, thisRepo, branchName, debugFlag)
		err := repo.Commit(&t, thisRepo, branchName, author, "Deleting label "+labelID, debugFlag)
		if err != nil {
			return err
		}
	} else {
		labelID := AddLabel(&t, label, thisRepo, branchName, debugFlag)
		err := repo.Commit(&t, thisRepo, branchName, author, "Adding label "+labelID, debugFlag)
		if err != nil {
			return err
		}
	}
	return nil
}

// DeleteLabel() takes a pointer to a ticket, a label, a git repo, a branch
// name, and a debug flag. It deletes the label from the ticket and returns an
// error.
func DeleteLabel(t *Ticket, label string, repo *git.Repository, branchName string, debug bool) string {
	updatedLabels := []string{}
	for _, l := range t.Labels {
		if l != label {
			updatedLabels = append(updatedLabels, l)
		}
	}
	t.Labels = updatedLabels
	return label
}

// AddLabel() takes a pointer to a ticket, a label, a git repo, a branch name,
// and a debug flag. It adds the label to the ticket and returns an error.
func AddLabel(t *Ticket, label string, repo *git.Repository, branchName string, debug bool) string {
	t.Labels = append(t.Labels, label)
	return label
}
