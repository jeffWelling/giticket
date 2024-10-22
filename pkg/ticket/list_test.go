package ticket

import (
	"strings"
	"testing"

	"github.com/jeffwelling/giticket/pkg/common"
	"github.com/jeffwelling/giticket/pkg/repo"
)

func TestWidest_FindWidestTitle(t *testing.T) {
	tickets := []Ticket{
		{Title: "a"},                      // 1
		{Title: "bb"},                     // 2
		{Title: strings.Repeat("X", 100)}, // 100
	}

	widest := widest(tickets, "Title")
	if widest != 100 {
		t.Errorf("Widest() = %d, want %d", widest, 100)
	}
}

func TestWidest_FindWidestDescription(t *testing.T) {
	tickets := []Ticket{
		{Description: "a"},                      // 1
		{Description: "bb"},                     // 2
		{Description: strings.Repeat("X", 100)}, // 100
	}

	widest := widest(tickets, "Description")
	if widest != 100 {
		t.Errorf("Widest() = %d, want %d", widest, 100)
	}
}

func TestHandleList(t *testing.T) {
	common.UseTempDir(t)

	// Initialize git and giticket
	err := repo.InitGitAndInitGiticket(t)
	if err != nil {
		t.Fatal(err)
	}

	testCases := []struct {
		branchName          string
		debugFlag           bool
		ticketTitle         string
		ticketDescription   string
		ticketID            int
		ticketSeverity      int
		ticketStatus        string
		expectedOutput      string
		expectOutputMissing bool
	}{
		{ // Test that you can see a ticket that was just added during init
			branchName:          "giticket",
			debugFlag:           true,
			ticketTitle:         "test HandleList",
			ticketDescription:   "test description",
			ticketID:            1,
			ticketSeverity:      1,
			ticketStatus:        "new",
			expectedOutput:      "1",
			expectOutputMissing: false,
		},
		{ // Test that searching for '99' in list output returns nothing
			branchName:          "giticket",
			debugFlag:           true,
			ticketTitle:         "test HandleList",
			ticketDescription:   "test description",
			ticketID:            1,
			ticketSeverity:      1,
			ticketStatus:        "newtwo",
			expectedOutput:      "99",
			expectOutputMissing: true,
		},
	}

	for _, testCase := range testCases {
		// Create a writer
		w := &strings.Builder{}

		// list tickets
		err := HandleList(w, 0, testCase.branchName, "", false, testCase.debugFlag)
		if err != nil {
			t.Fatal(err)
		}

		// Check output
		if testCase.expectOutputMissing {
			if strings.Contains(w.String(), testCase.expectedOutput) {
				t.Errorf("HandleList() included '%s' in the output but this value should be missing", testCase.expectedOutput)
			}
		} else {
			if !strings.Contains(w.String(), testCase.expectedOutput) {
				t.Errorf("HandleList() returned '%s' which does not include '%s'", w.String(), testCase.expectedOutput)
			}
		}
	}
}
