package ticket

func FilterTicketsByID(tickets []Ticket, id int) Ticket {
	var t Ticket
	for _, t_ := range tickets {
		if t_.ID == id {
			t = t_
		}
	}
	return t
}
