package ticket

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
	"time"

	git "github.com/jeffwelling/git2go/v37"
	"github.com/jeffwelling/giticket/pkg/common"
	"github.com/jeffwelling/giticket/pkg/debug"
	"gopkg.in/yaml.v2"
)

// HandleShow is used to print a list of giticket tickets in a number of formats
func HandleShow(ticketID int, output string, debugFlag bool, helpFlag bool) error {
	debug.DebugMessage(debugFlag, "Opening git repository")
	thisRepo, err := git.OpenRepository(".")
	if err != nil {
		return err
	}

	if helpFlag {
		return nil
	}

	tickets, err := GetListOfTickets(thisRepo, common.BranchName, debugFlag)
	if err != nil {
		return err
	}
	t := FilterTicketsByID(tickets, ticketID)
	ShowTicket(t, output, debugFlag)

	return nil
}

// ShowTicket takes a ticket, an output type, and a debug flag and prints the
// ticket details in the given format.
func ShowTicket(ticket Ticket, output string, debug bool) {
	// switch on output type
	switch output {
	case "text":
		ShowTicketsText(ticket, debug)
	case "yaml":
		ShowTicketsYaml(ticket, debug)
	case "json":
		ShowTicketsJson(ticket, debug)
	default:
		ShowTicketsText(ticket, debug)
	}
}

// ShowTicketsText is used to print a list of giticket tickets in text format
func ShowTicketsText(t Ticket, debug bool) {
	fmt.Println("ID: " + strconv.Itoa(t.ID))
	fmt.Println("Title: " + t.Title)
	fmt.Println("Description: " + t.Description)
	fmt.Println("Status: " + t.Status)
	fmt.Println("Severity: " + strconv.Itoa(t.Severity))
	fmt.Println("Labels: " + strings.Join(t.Labels, ", "))
	fmt.Println("Created: " + time.Unix(t.Created, 0).String())
	fmt.Println("NextTicketID: " + strconv.Itoa(t.NextCommentID))
	fmt.Println("Comments: ")

	for _, comment := range t.Comments {
		fmt.Println("    Comment ID: " + strconv.Itoa(t.ID) + "-" + strconv.Itoa(comment.ID))
		fmt.Println("    Created: " + time.Unix(comment.Created, 0).String())
		fmt.Println("    Author: " + comment.Author)
		fmt.Println("    Body: " + comment.Body)
	}
	fmt.Println("")
}

// ShowTicketsYaml is used to print a list of giticket tickets in yaml format
func ShowTicketsYaml(t Ticket, debug bool) {
	// turn the ticket into a yaml string
	yamlTicket, err := yaml.Marshal(t)
	if err != nil {
		panic(err)
	}
	fmt.Println(string(yamlTicket))
}

// ShowTicketsJson is used to print a list of giticket tickets in json format
func ShowTicketsJson(t Ticket, debug bool) {
	// turn the ticket into a json string
	jsonTicket, err := json.Marshal(t)
	if err != nil {
		panic(err)
	}
	fmt.Println(string(jsonTicket))
}
