package ticket

import (
	"testing"
	"time"

	"github.com/jeffwelling/giticket/pkg/common"
	"github.com/jeffwelling/giticket/pkg/repo"
)

func TestHandleDelete(t *testing.T) {
	common.UseTempDir(t)

	// Initialize git and giticket
	err := repo.InitGitAndInitGiticket(t)
	if err != nil {
		t.Fatal(err)
	}

	testCases := []struct {
		ticketID      int  // A ticketID of 0 indicates that the test should create a ticket and use that ticket's ID
		expectDeleted bool // If the test should expect a success response from HandleDelete
	}{
		{
			ticketID:      1,
			expectDeleted: true,
		},
	}

	for _, tc := range testCases {
		if tc.ticketID == 0 {
			// Create a ticket for testing
			ticketID, _ := HandleCreate(
				common.BranchName,
				time.Now().Unix(),
				"test ticket",
				"test description",
				[]string{"label1", "label2"},
				1,
				1,
				"new",
				[]Comment{},
				1,
				false,
			)

			tc.ticketID = ticketID
		}

		// Delete the ticket!
		deleted, err := HandleDelete(tc.ticketID, common.BranchName, false)
		if err != nil {
			t.Fatal(err)
		}
		if deleted != tc.expectDeleted {
			t.Errorf("Expected %t, got %t", tc.expectDeleted, deleted)
		}
	}
}
