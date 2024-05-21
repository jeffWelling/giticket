package ticket

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
	"time"

	"gopkg.in/yaml.v2"
)

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

func ShowTicketsYaml(t Ticket, debug bool) {
	// turn the ticket into a yaml string
	yamlTicket, err := yaml.Marshal(t)
	if err != nil {
		panic(err)
	}
	fmt.Println(string(yamlTicket))
}

func ShowTicketsJson(t Ticket, debug bool) {
	// turn the ticket into a json string
	jsonTicket, err := json.Marshal(t)
	if err != nil {
		panic(err)
	}
	fmt.Println(string(jsonTicket))
}
