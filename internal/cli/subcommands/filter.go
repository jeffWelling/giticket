package subcommands

import (
	"flag"
	"fmt"
	"os"

	"github.com/jeffwelling/giticket/pkg/common"
	"github.com/jeffwelling/giticket/pkg/ticket"
)

func init() {
	subcommand := new(SubcommandFilter)
	registerSubcommand("filter", subcommand)
}

type SubcommandFilter struct {
	debugFlag    bool
	deleteFlag   bool
	filter       string
	filterName   string
	flagset      *flag.FlagSet
	helpFlag     bool
	listFlag     bool
	outputFormat string
	parameters   map[string]interface{}
}

func (subcommand *SubcommandFilter) InitFlags(args []string) error {
	subcommand.flagset = flag.NewFlagSet("filter", flag.ExitOnError)
	subcommand.flagset.BoolVar(&subcommand.debugFlag, "debug", false, "Print debug info")
	subcommand.flagset.BoolVar(&subcommand.helpFlag, "help", false, "Print help")
	subcommand.flagset.BoolVar(&subcommand.deleteFlag, "delete", false, "Delete filter")
	subcommand.flagset.BoolVar(&subcommand.deleteFlag, "d", false, "Delete filter")
	subcommand.flagset.StringVar(&subcommand.filter, "filter", "", "Filter")
	subcommand.flagset.StringVar(&subcommand.filter, "f", "", "Filter")
	subcommand.flagset.StringVar(&subcommand.filterName, "filter-name", "", "Filter name")
	subcommand.flagset.StringVar(&subcommand.filterName, "name", "", "Filter name")
	subcommand.flagset.BoolVar(&subcommand.listFlag, "list", false, "List filters")
	subcommand.flagset.BoolVar(&subcommand.listFlag, "l", false, "List filters")
	subcommand.flagset.StringVar(&subcommand.outputFormat, "output-format", "json", "Output format, default is json, and can be 'json' or 'yaml'")
	subcommand.flagset.StringVar(&subcommand.outputFormat, "o", "json", "Output format, default is json, and can be 'json' or 'yaml'")
	if err := subcommand.flagset.Parse(args); err != nil {
		return err
	}

	if subcommand.helpFlag {
		common.PrintVersion()
		fmt.Println("giticket")
		subcommand.Help()
	}

	// Setup parameters
	subcommand.parameters = make(map[string]interface{})
	subcommand.parameters["debugFlag"] = subcommand.debugFlag
	subcommand.parameters["helpFlag"] = subcommand.helpFlag
	subcommand.parameters["deleteFlag"] = subcommand.deleteFlag
	subcommand.parameters["filter"] = subcommand.filter
	subcommand.parameters["filterName"] = subcommand.filterName
	subcommand.parameters["listFlag"] = subcommand.listFlag
	subcommand.parameters["outputFormat"] = subcommand.outputFormat

	// Sanity checks
	// If delete flag is set, then filter name must also be set
	if subcommand.deleteFlag && subcommand.filterName == "" {
		return fmt.Errorf("filter name must be set if delete flag is set")
	}

	// If delete is false and list is false then both filter name and filter are
	// required
	if !subcommand.deleteFlag && !subcommand.listFlag && (subcommand.filterName == "" || subcommand.filter == "") {
		return fmt.Errorf("filter name and filter must be set if not deleting or listing filters")
	}

	// If list is true then debug must be false, and both filter and filter name
	// must be empty strings
	if subcommand.listFlag && subcommand.debugFlag {
		return fmt.Errorf("debug flag cannot be set if listing filters")
	}
	if subcommand.listFlag && subcommand.filterName != "" {
		return fmt.Errorf("filter name cannot be set if listing filters")
	}
	if subcommand.listFlag && subcommand.filter != "" {
		return fmt.Errorf("filter cannot be set if listing filters")
	}

	// Sanity check that outputFormat is set to json or yaml
	if subcommand.outputFormat != "json" && subcommand.outputFormat != "yaml" {
		return fmt.Errorf("output format must be 'json' or 'yaml'")
	}

	return nil
}

func (subcommand *SubcommandFilter) Execute() {
	if subcommand.deleteFlag {
		err := ticket.HandleFilterDelete(subcommand.filterName, subcommand.debugFlag)
		if err != nil {
			fmt.Println(err)
			return
		}
	} else if subcommand.listFlag {
		err := ticket.HandleFilterList(os.Stdout, subcommand.outputFormat, subcommand.debugFlag)
		if err != nil {
			fmt.Println(err)
			return
		}
	} else {
		err := ticket.HandleFilterCreate(subcommand.filter, subcommand.filterName, subcommand.debugFlag)
		if err != nil {
			fmt.Println(err)
			return
		}
	}
}

func (subcommand *SubcommandFilter) Help() {
	fmt.Println("  filter - Set or delete filters for listing tickets")
	fmt.Println("    eg: giticket filter [params]")
	fmt.Println("    parameters:")
	fmt.Println("      --debug")
	fmt.Println("      --delete   | -d")
	fmt.Println("      --filter   | --f \"my filter\"")
	fmt.Println("      --filter-name | --name \"my filter name\"")
	fmt.Println("      --help")
	fmt.Println("      --list")
	fmt.Println("      --output-format | --o 'json' or 'yaml'")
	fmt.Println("    examples:")
	fmt.Println("      - name: Add filter \"my filter\"")
	fmt.Println("        example: giticket filter --filter \"my filter\" --filter-name \"my filter name\"")
	fmt.Println("      - name: List filters")
	fmt.Println("        example: giticket filter --list")
	fmt.Println("      - name: List filters in yaml format")
	fmt.Println("        example: giticket filter --list --output-format 'yaml'")
	fmt.Println("      - name: Delete filter \"my filter\"")
	fmt.Println("        example: giticket filter --delete --filter-name \"my filter\"")
}

func (subcommand *SubcommandFilter) Parameters() map[string]interface{} {
	return subcommand.parameters
}

func (subcommand *SubcommandFilter) DebugFlag() bool {
	return subcommand.debugFlag
}
