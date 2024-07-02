package ticket

import (
	"strconv"

	git "github.com/jeffwelling/git2go/v37"
	"github.com/jeffwelling/giticket/pkg/common"
	"github.com/jeffwelling/giticket/pkg/debug"
	"github.com/jeffwelling/giticket/pkg/repo"
)

func HandleStatus(
	status string,
	ticketID int,
	helpFlag bool,
	debugFlag bool,
) error {
	if helpFlag {
		return nil
	}

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

	t.Status = status
	err = repo.Commit(&t, thisRepo, common.BranchName, author, "Setting status of ticket "+strconv.Itoa(t.ID)+" to "+status, debugFlag)
	if err != nil {
		return err
	}
	return nil
}
