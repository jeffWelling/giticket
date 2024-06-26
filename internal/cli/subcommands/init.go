package subcommands

import (
	"flag"
	"fmt"

	"github.com/jeffwelling/giticket/pkg/repo"
)

// init() is used to register this action
func init() {
	subcommand := new(SubcommandInit)
	registerSubcommand("init", subcommand)
}

type SubcommandInit struct {
	debugFlag  bool
	helpFlag   bool
	flagset    *flag.FlagSet
	parameters map[string]interface{}
}

// Execute() is called when this action is invoked
func (subcommand *SubcommandInit) Execute() {
	repo.HandleInitGiticket(subcommand.debugFlag)
}

// Help() prints help for this action
func (subcommand *SubcommandInit) Help() {
	fmt.Println("  init - Initialize giticket")
	fmt.Println("    eg: giticket init")
	fmt.Println("    parameters:")
	fmt.Println("      --debug")
	fmt.Println("      --help")
	fmt.Println("    examples:")
	fmt.Println("      - name: Initialize giticket")
	fmt.Println("        example: giticket init")
}

// InitFlags()
func (subcommand *SubcommandInit) InitFlags(args []string) error {
	var (
		debugFlag bool
		helpFlag  bool
	)
	subcommand.flagset = flag.NewFlagSet("init", flag.ExitOnError)
	subcommand.flagset.BoolVar(&subcommand.debugFlag, "debug", false, "Print debug info")
	subcommand.flagset.BoolVar(&subcommand.helpFlag, "help", false, "Print help")
	subcommand.flagset.Parse(args)

	subcommand.parameters = make(map[string]interface{})
	subcommand.parameters["debugFlag"] = debugFlag
	subcommand.parameters["helpFlag"] = helpFlag

	// If help
	if subcommand.helpFlag {
		subcommand.Help()
		return nil
	}

	return nil
}

func (subcommand *SubcommandInit) DebugFlag() bool {
	return subcommand.debugFlag
}

func (subcommand *SubcommandInit) Parameters() map[string]interface{} {
	return make(map[string]interface{})
}
