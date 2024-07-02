package ticket

// FilterTicketsByID takes a list of tickets and an integer representing the
// ticket ID to filter for, and return that ticket.
func FilterTicketsByID(tickets []Ticket, id int) Ticket {
	var t Ticket
	for _, t_ := range tickets {
		if t_.ID == id {
			t = t_
		}
	}
	return t
}
