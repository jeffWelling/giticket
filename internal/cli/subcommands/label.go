package subcommands

import (
	"flag"
	"fmt"

	"github.com/jeffwelling/giticket/pkg/common"
	"github.com/jeffwelling/giticket/pkg/ticket"
)

func init() {
	subcommand := new(SubcommandLabel)
	registerSubcommand("label", subcommand)
}

type SubcommandLabel struct {
	flagset    *flag.FlagSet
	debugFlag  bool
	helpFlag   bool
	label      string
	deleteFlag bool
	ticketID   int
	parameters map[string]interface{}
}

func (subcommand *SubcommandLabel) InitFlags(args []string) error {

	subcommand.parameters = make(map[string]interface{})
	var (
		helpFlag  bool
		ticketID  int
		label     string
		delete    bool
		debugFlag bool
	)

	subcommand.flagset = flag.NewFlagSet("label", flag.ExitOnError)

	subcommand.flagset.BoolVar(&subcommand.debugFlag, "debug", false, "Print debug info")
	subcommand.flagset.BoolVar(&subcommand.helpFlag, "help", false, "Print help")

	subcommand.flagset.StringVar(&subcommand.label, "label", "", "Label to add")
	subcommand.flagset.StringVar(&subcommand.label, "l", "", "Label to add")
	subcommand.flagset.BoolVar(&subcommand.deleteFlag, "delete", false, "Delete label")
	subcommand.flagset.BoolVar(&subcommand.deleteFlag, "d", false, "Delete label")
	subcommand.flagset.IntVar(&subcommand.ticketID, "ticketid", 0, "Ticket ID")
	subcommand.flagset.IntVar(&subcommand.ticketID, "id", 0, "Ticket ID")
	subcommand.flagset.Parse(args)

	subcommand.parameters["debugFlag"] = debugFlag
	subcommand.parameters["helpFlag"] = helpFlag
	subcommand.parameters["label"] = label
	subcommand.parameters["deleteFlag"] = delete
	subcommand.parameters["ticketID"] = ticketID

	if subcommand.helpFlag {
		common.PrintVersion()
		fmt.Println("giticket")
		subcommand.Help()
	}

	// Sanity check of args
	if subcommand.ticketID == 0 {
		fmt.Println("Error: Ticket ID must be specified")
		// Print usage
		common.PrintGeneralUsage()
		subcommand.Help()
	}
	if subcommand.deleteFlag && subcommand.label == "" {
		fmt.Println("Error: When deleting a label, the label must be specified")
		// Print usage
		common.PrintGeneralUsage()
		subcommand.Help()
	}

	if subcommand.label == "" {
		fmt.Println("Error: Label must be specified in order to add a label")
		// Print usage
		common.PrintGeneralUsage()
		subcommand.Help()
	}

	return nil
}

func (subcommand *SubcommandLabel) Execute() {
	ticket.HandleLabel(
		common.BranchName,
		subcommand.label,
		subcommand.deleteFlag,
		subcommand.ticketID,
		subcommand.debugFlag,
	)
}

func (subcommand *SubcommandLabel) Help() {
	fmt.Println("  label - Add or delete labels")
	fmt.Println("    eg: giticket label [params]")
	fmt.Println("    parameters:")
	fmt.Println("      --ticketid | --id 1")
	fmt.Println("      --label    | --l \"my first label\"")
	fmt.Println("      --delete   | -d")
	fmt.Println("      --debug")
	fmt.Println("      --help")
	fmt.Println("    examples:")
	fmt.Println("      - name: Add label \"my first label\" to ticket with ID #1")
	fmt.Println("        example: giticket label --ticketid 1 --label \"my first label\"")
	fmt.Println("      - name: Delete label \"my first label\" from ticket with ID #1")
	fmt.Println("        example: giticket label --ticketid 1 --label \"my first label\" --delete")
}

// Parameters
func (subcommand *SubcommandLabel) Parameters() map[string]interface{} {
	return subcommand.parameters
}

// DebugFlag
func (subcommand *SubcommandLabel) DebugFlag() bool {
	return subcommand.debugFlag
}
