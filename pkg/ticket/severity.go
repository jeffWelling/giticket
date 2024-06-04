package ticket

import (
	"strconv"

	git "github.com/jeffwelling/git2go/v37"
	"github.com/jeffwelling/giticket/pkg/common"
	"github.com/jeffwelling/giticket/pkg/debug"
	"github.com/jeffwelling/giticket/pkg/repo"
)

func HandleSeverity(ticketID int, severity int, debugFlag bool) error {
	debug.DebugMessage(debugFlag, "Opening git repository")
	thisRepo, err := git.OpenRepository(".")
	if err != nil {
		panic(err)
	}

	// Get author
	author, err := common.GetAuthor(thisRepo)
	if err != nil {
		panic(err)
	}

	tickets, err := GetListOfTickets(thisRepo, common.BranchName, debugFlag)
	if err != nil {
		panic(err)
	}
	t := FilterTicketsByID(tickets, ticketID)
	t.Severity = severity

	repo.Commit(&t, thisRepo, common.BranchName, author, "Setting severity of ticket "+strconv.Itoa(t.ID)+" to "+strconv.Itoa(severity), debugFlag)

	return nil
}
