package ticket

import (
	"testing"

	"github.com/jeffwelling/giticket/pkg/common"
	"github.com/jeffwelling/giticket/pkg/repo"
)

func TestHandleCreate(t *testing.T) {
	common.UseTempDir(t)

	// Initialize git and giticket
	err := repo.InitGitAndInitGiticket(t)
	if err != nil {
		t.Fatal(err)
	}
	testCases := []struct {
		branchName    string
		created       int64
		title         string
		description   string
		labels        []string
		priority      int
		severity      int
		status        string
		comments      []Comment
		nextCommentId int
		debugFlag     bool
	}{
		{
			branchName:    "giticket",
			created:       1716538263,
			title:         "test HandleCreate",
			description:   "test description",
			labels:        []string{"label1", "label2"},
			priority:      1,
			severity:      1,
			status:        "new",
			comments:      []Comment{},
			nextCommentId: 1,
			debugFlag:     true,
		},
		{
			branchName:    "giticket",
			created:       1716538264,
			title:         "test HandleCreate 2",
			description:   "test description",
			labels:        []string{},
			priority:      1,
			severity:      1,
			status:        "old",
			comments:      []Comment{},
			nextCommentId: 1,
			debugFlag:     true,
		},
	}

	for _, tc := range testCases {
		_, ticketFilename, err := HandleCreate(
			tc.branchName,
			tc.created,
			tc.title,
			tc.description,
			tc.labels,
			tc.priority,
			tc.severity,
			tc.status,
			tc.comments,
			tc.nextCommentId,
			tc.debugFlag,
		)
		if err != nil {
			t.Fatal(err)
		}

		exists, err := repo.TicketExists(ticketFilename, tc.debugFlag)
		if err != nil {
			t.Fatal(err)
		}
		if exists != true {
			t.Fatal("Ticket '" + ticketFilename + "' does not exist")
		}
	}
}
