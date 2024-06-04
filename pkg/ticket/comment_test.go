package ticket

import (
	"testing"

	"github.com/jeffwelling/giticket/pkg/common"
	"github.com/jeffwelling/giticket/pkg/repo"
)

func TestHandleComment(t *testing.T) {

	common.UseTempDir(t)

	// Initialize git and giticket
	err := repo.InitGitAndInitGiticket(t)
	if err != nil {
		t.Fatal(err)
	}

	testCases := []struct {
		deleteFlag            bool
		comment               string
		commentID             int
		ticketID              int
		debugFlag             bool
		expectedFullCommentID string
	}{
		{
			deleteFlag:            false,
			comment:               "test comment 1",
			commentID:             0,
			ticketID:              1,
			debugFlag:             true,
			expectedFullCommentID: "1-3",
		},
		{
			deleteFlag:            true,
			comment:               "",
			commentID:             3,
			ticketID:              1,
			debugFlag:             true,
			expectedFullCommentID: "1-3",
		},
		{
			deleteFlag:            false,
			comment:               "test comment 2",
			commentID:             0,
			ticketID:              1,
			debugFlag:             true,
			expectedFullCommentID: "1-4",
		},
	}

	for _, tc := range testCases {
		commentID, err := HandleComment(common.BranchName, tc.comment, tc.commentID, tc.ticketID, tc.deleteFlag, tc.debugFlag)
		if err != nil {
			t.Fatal(err)
		}
		if commentID != tc.expectedFullCommentID {
			t.Fatalf("Expected commentID '%s', got '%s'", tc.expectedFullCommentID, commentID)
		}

		// Check that the comment is actually created/deleted by looking at the
		// git branch and not just the response from the function. There have
		// been times when the function reported it wrote successfully but it
		// didn't actually write the tree as expected.
		exists, err := repo.CommentExists(common.BranchName, commentID, tc.debugFlag)
		if err != nil {
			t.Fatal(err)
		}
		if exists != !tc.deleteFlag {
			t.Errorf("Expected comment to exist %t, got %t", !tc.deleteFlag, exists)
		}
	}
}
