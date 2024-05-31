package subcommands

import (
	"flag"
	"fmt"
	"reflect"
	"strings"

	git "github.com/jeffwelling/git2go/v37"
	"github.com/jeffwelling/giticket/pkg/debug"
	"github.com/jeffwelling/giticket/pkg/ticket"
)

func init() {
	subcommand := new(SubcommandList)
	registerSubcommand("list", subcommand)
}

type SubcommandList struct {
	flagset     *flag.FlagSet
	debugFlag   bool
	helpFlag    bool
	windowWidth int
	parameters  map[string]interface{}
}

func (subcommand *SubcommandList) InitFlags(args []string) error {
	subcommand.parameters = make(map[string]interface{})
	var (
		helpFlag  bool
		window    int
		debugFlag bool
	)
	subcommand.flagset = flag.NewFlagSet("list", flag.ExitOnError)

	subcommand.flagset.BoolVar(&subcommand.debugFlag, "debug", false, "Print debug info")
	subcommand.flagset.BoolVar(&subcommand.helpFlag, "help", false, "Print help")
	subcommand.flagset.IntVar(&subcommand.windowWidth, "window", 0, "Window width")
	subcommand.flagset.IntVar(&subcommand.windowWidth, "w", 0, "Window width")
	subcommand.flagset.Parse(args)

	subcommand.parameters["debugFlag"] = debugFlag
	subcommand.parameters["helpFlag"] = helpFlag
	subcommand.parameters["windowWidth"] = window
	return nil
}

func (subcommand *SubcommandList) Execute() {
	branchName := "giticket"

	debug.DebugMessage(subcommand.debugFlag, "Opening git repository")
	thisRepo, err := git.OpenRepository(".")
	if err != nil {
		panic(err)
	}

	output, err := ListTickets(thisRepo, branchName, subcommand.windowWidth, subcommand.debugFlag)
	if err != nil {
		panic(err)
	}
	fmt.Print(output)
}

func (subcommand *SubcommandList) Help() {
	fmt.Println("  list - List tickets")
	fmt.Println("    eg: giticket list [params]")
}

type listParams struct {
	windowLength int
	debugFlag    bool
}

// ListTickets() takes a listParams parameter which contains the optional
// and mandatory parameters for ListTickets(). The only mandatory parameter is
// windowLength which is the length of the window to list tickets in.
func ListTickets(thisRepo *git.Repository, branchName string, windowWidth int, debugFlag bool) (string, error) {
	output := ""

	// Get a list of tickets from the repo
	var ticketsList []ticket.Ticket
	ticketsList, err := ticket.GetListOfTickets(thisRepo, branchName, debugFlag)
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

// Parameters
func (subcommand *SubcommandList) Parameters() map[string]interface{} {
	return subcommand.parameters
}

// DebugFlag
func (subcommand *SubcommandList) DebugFlag() bool {
	return subcommand.debugFlag
}
