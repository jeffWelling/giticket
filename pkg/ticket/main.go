package ticket

type Ticket struct {
	Title       string
	Id          int
	Description string
	Tags        []string
	Priority    int
	Severity    int
	Comments    []Comment
	Status      string
	Created     int64
}

type Comment struct {
	Id      string
	Created int64
	Body    string
	Author  string
}
