package ticket

import (
	"fmt"
	"io"
	"reflect"
	"strings"

	git "github.com/jeffwelling/git2go/v37"
	"github.com/jeffwelling/giticket/pkg/debug"
)

func HandleList(debugFlag bool, branchName string, windowWidth int, w io.Writer) error {
	debug.DebugMessage(debugFlag, "Opening git repository")
	thisRepo, err := git.OpenRepository(".")
	if err != nil {
		return err
	}

	output, err := ListTickets(
		thisRepo, branchName, windowWidth, debugFlag)
	if err != nil {
		return err
	}
	fmt.Fprint(w, output)
	return nil
}

func ListTickets(thisRepo *git.Repository, branchName string, windowWidth int, debugFlag bool) (string, error) {
	output := ""

	// Get a list of tickets from the repo
	var ticketsList []Ticket
	ticketsList, err := GetListOfTickets(thisRepo, branchName, debugFlag)
	if err != nil {
		return "", fmt.Errorf("Unable to list tickets: %s", err) // TODO: err
	}

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

	return output, nil
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
func widest(tickets []Ticket, attr string) int {
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
