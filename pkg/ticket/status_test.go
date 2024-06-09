package ticket

import (
	"testing"

	git "github.com/jeffwelling/git2go/v37"
	"github.com/jeffwelling/giticket/pkg/common"
	"github.com/jeffwelling/giticket/pkg/debug"
	"github.com/jeffwelling/giticket/pkg/repo"
)

func TestHandleStatus(t *testing.T) {
	common.UseTempDir(t)

	// Initialize git and giticket
	err := repo.InitGitAndInitGiticket(t)
	if err != nil {
		t.Fatal(err)
	}

	testCases := []struct {
		branchName   string
		statusBefore string
		statusAfter  string
		ticketID     int
		debugFlag    bool
	}{
		{
			branchName:   "giticket",
			statusBefore: "open",
			statusAfter:  "spoilers",
			ticketID:     1,
			debugFlag:    true,
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
		if preTicket.Status != tc.statusBefore {
			t.Fatal("Ticket priority before change should be " + tc.statusBefore + " but is " + preTicket.Status)
		}

		// Change ticket priority
		err = HandleStatus(tc.statusAfter, tc.ticketID, false, tc.debugFlag)
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
		if postTicket.Status != tc.statusAfter {
			t.Fatalf("After priority change, was expecting priority '%s' but got priority '%s'", tc.statusAfter, postTicket.Status)
		}
	}
}
