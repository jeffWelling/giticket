package subcommands

import (
	"flag"
	"fmt"
	"reflect"
	"strings"

	"github.com/jeffWelling/giticket/pkg/ticket"
)

func init() {
	subcommand := new(SubcommandList)
	registerSubcommand("list", subcommand)
}

type SubcommandList struct {
	flagset *flag.FlagSet
	debug   bool
}

func (subcommand *SubcommandList) InitFlags(args []string) {
	subcommand.flagset = flag.NewFlagSet("list", flag.ExitOnError)
}

func (subcommand *SubcommandList) Execute() {
	listParams := listParams{
		windowLength: 0,
		debug:        subcommand.debug,
	}

	output := ListTickets(listParams)
	fmt.Print(output)
}

func (subcommand *SubcommandList) Help() {
	fmt.Println("  list - List tickets")
	fmt.Println("    eg: giticket list [params]")
}

type listParams struct {
	windowLength int
	debug        bool
}

// ListTickets() takes a listParams parameter which contains the optional
// and mandatory parameters for ListTickets(). The only mandatory parameter is
// windowLength which is the length of the window to list tickets in.
func ListTickets(params listParams) string {
	output := ""

	// Get a list of tickets from the repo
	var ticketsList []ticket.Ticket
	ticketsList = ticket.GetListOfTickets(params.debug)

	widthOfID := widest(ticketsList, "ID")
	if widthOfID < 3 {
		widthOfID = 3
	}
	widthOfTitle := widest(ticketsList, "Title")
	if widthOfTitle < 20 {
		widthOfTitle = 20
	}
	widthOfSeverity := widest(ticketsList, "Severity")
	if widthOfSeverity < 9 {
		widthOfSeverity = 9
	}
	widthOfStatus := widest(ticketsList, "Status")
	if widthOfStatus < 10 {
		widthOfStatus = 10
	}

	// Print the header
	output += padRight("ID", widthOfID) + " | " + padRight("Title", widthOfTitle) + " | " + padRight("Severity", widthOfSeverity) + " | " + padRight("Status", widthOfStatus) + "\n"
	output += strings.Repeat("-", widthOfID+widthOfTitle+widthOfSeverity+widthOfStatus+4) + "\n"

	// Print the tickets
	for _, t := range ticketsList {
		IDAsString := fmt.Sprintf("%d", t.ID)
		SeverityAsString := fmt.Sprintf("%d", t.Severity)
		output += fmt.Sprintf("%s | %s | %s | %s\n", padRight(IDAsString, widthOfID), padRight(t.Title, widthOfTitle), padRight(SeverityAsString, widthOfSeverity), padRight(t.Status, widthOfStatus))
	}

	return output
}

// padRight() takes string s and width int, it finds the difference in length
// between len(s) and width and adds that many spaces to the string to ensure
// the returned string is exactly width len long
func padRight(s string, width int) string {
	diff := width - len(s)
	if diff <= 0 {
		return s[0:width]
	}
	return s + strings.Repeat(" ", diff)
}

// widest() takes a list of tickets and a string representing the attribute name
// of each ticket to check to find the widest string, and return that value
func widest(tickets []ticket.Ticket, attr string) int {
	widest := 0
	for _, ticket := range tickets {
		v := reflect.ValueOf(ticket)
		if v.Kind() != reflect.Struct {
			panic("not a struct")
		}

		fieldVal := v.FieldByName(attr)
		if !fieldVal.IsValid() {
			panic("not a valid field")
		}
		if !fieldVal.CanInterface() {
			panic("cannot interface")
		}

		// if fieldVal is an int, conver it to string
		fieldValString := fmt.Sprintf("%v", fieldVal.Interface())

		if len(fieldValString) > widest {
			widest = len(fieldValString)
		}
	}
	return widest
}