package ticket

import (
	"testing"

	git "github.com/jeffwelling/git2go/v37"
	"github.com/jeffwelling/giticket/pkg/common"
	"github.com/jeffwelling/giticket/pkg/debug"
	"github.com/jeffwelling/giticket/pkg/repo"
)

func TestHandleLabel(t *testing.T) {
	common.UseTempDir(t)

	// Initialize git and giticket
	err := repo.InitGitAndInitGiticket(t)
	if err != nil {
		t.Fatal(err)
	}

	testCases := []struct {
		branchName string
		label      string
		deleteFlag bool
		ticketID   int
		debugFlag  bool
	}{
		{
			branchName: "giticket",
			label:      "label1",
			deleteFlag: false,
			ticketID:   1,
			debugFlag:  true,
		},
		{
			branchName: "giticket",
			label:      "label1",
			deleteFlag: true,
			ticketID:   1,
			debugFlag:  true,
		},
	}

	// Open repo
	debug.DebugMessage(false, "Opening git repository")
	thisRepo, err := git.OpenRepository(".")
	if err != nil {
		t.Fatal(err)
	}

	for _, testCase := range testCases {
		err := HandleLabel(
			testCase.branchName,
			testCase.label,
			testCase.deleteFlag,
			testCase.ticketID,
			testCase.debugFlag,
		)
		if err != nil {
			t.Fatal(err)
		}

		// Get tickets
		tickets, err := GetListOfTickets(thisRepo, testCase.branchName, testCase.debugFlag)
		if err != nil {
			t.Fatal(err)
		}

		ticket := FilterTicketsByID(tickets, testCase.ticketID)

		if repo.CheckLabel(ticket.Labels, testCase.branchName, testCase.label, testCase.deleteFlag) {
			return
		}
		t.Fatal("Label not found on ticket")
	}
}
