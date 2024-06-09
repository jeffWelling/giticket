package ticket

import (
	"strconv"
	"testing"

	git "github.com/jeffwelling/git2go/v37"
	"github.com/jeffwelling/giticket/pkg/common"
	"github.com/jeffwelling/giticket/pkg/debug"
	"github.com/jeffwelling/giticket/pkg/repo"
)

func TestHandlePriority(t *testing.T) {
	common.UseTempDir(t)

	// Initialize git and giticket
	err := repo.InitGitAndInitGiticket(t)
	if err != nil {
		t.Fatal(err)
	}

	testCases := []struct {
		branchName     string
		priorityBefore int
		priorityAfter  int
		ticketID       int
		debugFlag      bool
	}{
		{
			branchName:     "giticket",
			priorityBefore: 1,
			priorityAfter:  2,
			ticketID:       1,
			debugFlag:      true,
		},
	}

	debug.DebugMessage(true, "Opening git repository")
	thisRepo, err := git.OpenRepository(".")
	if err != nil {
		t.Fatal(err)
	}

	for _, tc := range testCases {
		// Get ticket
		tickets, err := GetListOfTickets(thisRepo, common.BranchName, true)
		if err != nil {
			t.Fatal(err)
		}
		preTicket := FilterTicketsByID(tickets, tc.ticketID)

		// Get ticket priority
		if preTicket.Priority != tc.priorityBefore {
			t.Fatal("Ticket priority before change should be " + strconv.Itoa(tc.priorityBefore) + " but is " + strconv.Itoa(preTicket.Priority))
		}

		// Change ticket priority
		err = HandlePriority(preTicket.ID, tc.priorityAfter, tc.debugFlag)
		if err != nil {
			t.Fatal(err)
		}

		// Get new ticket priority to compare
		newTicketsList, err := GetListOfTickets(thisRepo, common.BranchName, true)
		if err != nil {
			t.Fatal(err)
		}
		postTicket := FilterTicketsByID(newTicketsList, tc.ticketID)

		// Compare priority
		if postTicket.Priority != tc.priorityAfter {
			t.Fatalf("After priority change, was expecting priority '%d' but got priority '%d'", postTicket.Priority, tc.priorityAfter)
		}
	}
}
