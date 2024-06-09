package ticket

import "testing"

func TestFilterTicketsByID(t *testing.T) {
	testCases := []struct {
		tickets  []Ticket
		ticketID int
		expected Ticket
	}{
		{
			tickets: []Ticket{
				{
					ID: 1,
				},
				{
					ID: 2,
				},
				{
					ID: 3,
				},
			},
			ticketID: 2,
			expected: Ticket{
				ID: 2,
			},
		},
	}

	for _, testCase := range testCases {
		actual := FilterTicketsByID(testCase.tickets, testCase.ticketID)
		if actual.ID != testCase.expected.ID {
			t.Errorf("Expected %v, got %v", testCase.expected.ID, actual.ID)
		}
	}
}
