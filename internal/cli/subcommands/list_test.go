package subcommands

import (
	"strings"
	"testing"

	"github.com/jeffWelling/giticket/pkg/ticket"
)

func TestWidest_FindWidestTitle(t *testing.T) {
	tickets := []ticket.Ticket{
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
	tickets := []ticket.Ticket{
		{Description: "a"},                      // 1
		{Description: "bb"},                     // 2
		{Description: strings.Repeat("X", 100)}, // 100
	}

	widest := widest(tickets, "Description")
	if widest != 100 {
		t.Errorf("Widest() = %d, want %d", widest, 100)
	}
}
