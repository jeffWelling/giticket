package ticket

import (
	"strconv"

	git "github.com/jeffwelling/git2go/v37"
	"github.com/jeffwelling/giticket/pkg/common"
	"github.com/jeffwelling/giticket/pkg/debug"
	"github.com/jeffwelling/giticket/pkg/repo"
)

func HandlePriority(ticketID int, priority int, debugFlag bool) error {
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

	tickets, err := GetListOfTickets(thisRepo, common.BranchName, debugFlag)
	if err != nil {
		return err
	}
	t := FilterTicketsByID(tickets, ticketID)
	t.Priority = priority
	repo.Commit(&t, thisRepo, common.BranchName, author, "Setting priority of ticket "+strconv.Itoa(t.ID)+" to "+strconv.Itoa(priority)+"", debugFlag)

	return nil
}
