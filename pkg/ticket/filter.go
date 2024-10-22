package ticket

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"strconv"
	"time"

	git "github.com/jeffwelling/git2go/v37"
	"github.com/jeffwelling/giticket/pkg/common"
	"github.com/jeffwelling/giticket/pkg/debug"
	"github.com/jeffwelling/giticket/pkg/repo"
	"gopkg.in/yaml.v2"

	"github.com/itchyny/gojq"
)

// Filter is used in ticket list operations to return a subset of tickets
type Filter struct {
	Name   string
	Filter string

	CreatedAt string
}

// FilterList is a list of Filters and a value that represents the 'current'
// filter to use
type FilterList struct {
	CurrentFilter string
	Filters       map[string]Filter
}

// FilterTicketsByID takes a list of tickets and an integer representing the
// ticket ID to filter for, and return that ticket.
func FilterTicketsByID(tickets []Ticket, id int) Ticket {
	var t Ticket
	for _, t_ := range tickets {
		if t_.ID == id {
			t = t_
		}
	}
	return t
}

// HandleFilterDelete takes the name of a filter and a debug flag, and deletes
// the filter. It returns an error if there is one.
func HandleFilterDelete(filterName string, debugFlag bool) error {
	debug.DebugMessage(debugFlag, "Deleting filter: "+filterName)

	// Get list of filters
	filters, err := GetFilters(common.BranchName, debugFlag)
	if err != nil {
		return err
	}

	// Delete filter identified by filterName
	debug.DebugMessage(debugFlag, "Deleting filter: "+filterName+" from list of filters")
	for loadedFilterName, filter := range filters.Filters {
		if filter.Name == filterName {
			delete(filters.Filters, loadedFilterName)
			break
		}
	}

	// Write filters
	err = WriteFilters(filters, "Deleted filter: "+filterName, common.BranchName, debugFlag)
	if err != nil {
		return err
	}
	return nil
}

// HandleFilterList takes a debug flag and lists all filters. It returns an
// error if there is one.
func HandleFilterList(writer io.Writer, outputFormat string, debugFlag bool) error {
	// Get list of filters
	filters, err := GetFilters(common.BranchName, debugFlag)
	if err != nil {
		// Return a helpful error if the error message is:
		// the path 'filters.json' does not exist in the given tree
		if err.Error() == "the path 'filters.json' does not exist in the given tree" {
			return errors.New("There are no filters to list yet.")
		}

		return err
	}

	// Print the filters
	err = printFilters(filters, writer, outputFormat, debugFlag)
	if err != nil {
		return err
	}

	return nil
}

// HandleFilterCreate takes a filter string, a filter name, and a debug flag
// and creates a filter. It returns an error if there is one.
func HandleFilterCreate(filter string, filterName string, debugFlag bool) error {
	debug.DebugMessage(debugFlag, "Creating filter: "+filterName)

	// Check the filter is valid
	err := checkFilterIsValid(filter, filterName, debugFlag)
	if err != nil {
		return err
	}

	// Load the list of filters to add the new one too
	listOfFilters, err := GetFilters(common.BranchName, debugFlag)
	if err != nil {
		// If the error is that filters.json doesn't exist, eat the error and
		// continue reasonably
		if err.Error() != "the path 'filters.json' does not exist in the given tree" {
			return err
		} else {
			debug.DebugMessage(debugFlag, "Creating empty list of filters because filters.json doesn't exist yet")
			listOfFilters = new(FilterList)
			listOfFilters.Filters = make(map[string]Filter)
		}
	}

	// Add filter to list
	debug.DebugMessage(debugFlag, "Adding filter: "+filterName+" to list of filters")
	listOfFilters.Filters[filterName] = filterFromString(filter, filterName)

	// Write list
	err = WriteFilters(listOfFilters, "Created new filter", common.BranchName, debugFlag)
	if err != nil {
		return err
	}

	return nil
}

func GetFilters(branchName string, debugFlag bool) (*FilterList, error) {
	debug.DebugMessage(debugFlag, "GetFilters() start")
	debug.DebugMessage(debugFlag, "Getting filters from branch '"+branchName+"'")
	var filters FilterList

	debug.DebugMessage(debugFlag, "Opening git repository")
	thisRepo, err := git.OpenRepository(".")
	if err != nil {
		return &filters, err
	}

	debug.DebugMessage(debugFlag, "Getting parent commit from branch '"+branchName+"'")
	parentCommit, err := repo.GetParentCommit(thisRepo, branchName, debugFlag)
	if err != nil {
		return &filters, err
	}

	debug.DebugMessage(debugFlag, "Getting rootTreeBuilder and previousCommitTree from parent commit")
	rootTreeBuilder, previousCommitTree, err := repo.TreeBuilderFromCommit(parentCommit, thisRepo, debugFlag)
	if err != nil {
		return &filters, err
	}
	defer rootTreeBuilder.Free()

	debug.DebugMessage(debugFlag, "Getting .giticket subtree from previous commit")
	giticketTree, err := repo.GetSubTreeByName(previousCommitTree, thisRepo, ".giticket", debugFlag)
	if err != nil {
		return &filters, err
	}
	defer giticketTree.Free()

	debug.DebugMessage(debugFlag, "Getting filters.json from .giticket subtree")
	filtersFileEntry, err := giticketTree.EntryByPath("filters.json")
	if err != nil {
		return &filters, err
	}

	debug.DebugMessage(debugFlag, "Reading filters.json from .giticket subtree")
	filtersFileBlob, err := thisRepo.LookupBlob(filtersFileEntry.Id)
	if err != nil {
		return &filters, err
	}

	debug.DebugMessage(debugFlag, "Getting filters.json contents")
	filtersFileContents := filtersFileBlob.Contents()

	// JSON decode filtersFileContents into filters
	debug.DebugMessage(debugFlag, "Decoding filters.json contents")
	err = json.Unmarshal(filtersFileContents, &filters)
	if err != nil {
		return &filters, err
	}

	debug.DebugMessage(debugFlag, "Done getting filters, returning them")
	return &filters, nil
}

func WriteFilters(filters *FilterList, commitMessage string, branchName string, debugFlag bool) error {
	debug.DebugMessage(debugFlag, "Writing "+strconv.Itoa(len(filters.Filters))+" filters to branch '"+branchName+"'")

	// Open git repository
	debug.DebugMessage(debugFlag, "Opening git repository")
	thisRepo, err := git.OpenRepository(".")
	if err != nil {
		return err
	}

	debug.DebugMessage(debugFlag, "Getting parent commit from branch '"+branchName+"'")
	parentCommit, err := repo.GetParentCommit(thisRepo, branchName, debugFlag)
	if err != nil {
		return err
	}

	rootTreeBuilder, previousCommitTree, err := repo.TreeBuilderFromCommit(parentCommit, thisRepo, debugFlag)
	if err != nil {
		return err
	}
	defer rootTreeBuilder.Free()

	giticketTree, err := repo.GetSubTreeByName(previousCommitTree, thisRepo, ".giticket", debugFlag)
	if err != nil {
		return err
	}
	defer giticketTree.Free()

	// Create giticketTreeBuilder
	debug.DebugMessage(debugFlag, "Creating giticketTreeBuilder")
	giticketTreeBuilder, err := thisRepo.TreeBuilderFromTree(giticketTree)
	if err != nil {
		return err
	}
	defer giticketTreeBuilder.Free()

	// Convert filters into json string
	debug.DebugMessage(debugFlag, "Converting filters into json string")
	filtersJSON, err := json.Marshal(filters)
	if err != nil {
		return err
	}

	// Write filtersJSON to repo
	debug.DebugMessage(debugFlag, "Writing filtersJSON to repo")
	filtersFileBlobID, err := thisRepo.CreateBlobFromBuffer(filtersJSON)
	if err != nil {
		return err
	}

	// Insert filtersFileBlobID to giticketTreeBuilder
	debug.DebugMessage(debugFlag, "Inserting filtersFileBlobID to giticketTreeBuilder")
	err = giticketTreeBuilder.Insert("filters.json", filtersFileBlobID, git.FilemodeBlob)
	if err != nil {
		return err
	}

	// Write giticketTreeBuilder
	debug.DebugMessage(debugFlag, "Writing giticketTreeBuilder")
	giticketTreeOID, err := giticketTreeBuilder.Write()
	if err != nil {
		return err
	}

	// Insert giticketTreeOID to rootTreeBuilder
	debug.DebugMessage(debugFlag, "Inserting giticketTreeOID to rootTreeBuilder")
	err = rootTreeBuilder.Insert(".giticket", giticketTreeOID, git.FilemodeTree)
	if err != nil {
		return err
	}

	// Write rootTreeBuilder
	debug.DebugMessage(debugFlag, "Writing rootTreeBuilder")
	updatedRootTreeOID, err := rootTreeBuilder.Write()
	if err != nil {
		return err
	}

	// Lookup updatedRootTree
	debug.DebugMessage(debugFlag, "Lookup updatedRootTreeOID")
	updatedRootTree, err := thisRepo.LookupTree(updatedRootTreeOID)
	if err != nil {
		return err
	}
	defer updatedRootTree.Free()

	// Get git commit author
	debug.DebugMessage(debugFlag, "Getting git commit author")
	author, err := common.GetAuthor(thisRepo)
	if err != nil {
		return err
	}

	// Commit
	debug.DebugMessage(debugFlag, "Creating commit, branch: "+branchName+" message: "+commitMessage)
	commitID, err := thisRepo.CreateCommit("refs/heads/"+branchName, author, author, "Updated filters: "+commitMessage, updatedRootTree, parentCommit)
	if err != nil {
		return err
	}
	debug.DebugMessage(debugFlag, "Created commit: "+commitID.String())

	debug.DebugMessage(debugFlag, "Done writing filters")
	return nil
}

// checkFilterIsValid() takes a filter, the filter name, and a debug flag. It
// checks that the filter submitted is valid by testing that it can be used
// against a test list of tickets to ensure it doesn't throw an error. Intended
// for use as part of the 'add filter' code flow. Returns an error if there is
// one.
// We can't check that the filter returns the expected value but we can check
// that it can be used without throwing an error.
func checkFilterIsValid(filter string, name string, debugFlag bool) error {
	debug.DebugMessage(debugFlag, "Checking filter validity for filter: "+name)
	if filter == "" {
		return errors.New("Error validating filter: Filter cannot be empty")
	}

	// Create a set of test tickets to work with and turn them into JSON
	jsonListOfTickets, err := json.Marshal([]Ticket{
		{
			ID: 1,
		},
		{
			ID: 2,
		},
		{
			ID: 3,
		},
	})
	if err != nil {
		return err
	}

	// Turn jsonListOfTickets into a map[string]interface{}
	var listOfTickets []Ticket
	err = json.Unmarshal(jsonListOfTickets, &listOfTickets)
	if err != nil {
		return fmt.Errorf("Error unmarshalling jsonListOfTickets to validate filter: " + err.Error())
	}

	queryObj, err := gojq.Parse(filter)
	if err != nil {
		return fmt.Errorf("Filter validation error, unable to parse: " + err.Error())
	}

	// Just check that the filter can be used, we don't care about the result of
	// the filter operation
	iter := queryObj.Run(listOfTickets)
	for {
		result, ok := iter.Next()
		if !ok {
			break
		}
		if err, ok := result.(error); ok {
			return fmt.Errorf("Filter validation error: " + err.Error())
		}
	}

	// The filter appears valid
	debug.DebugMessage(debugFlag, "Filter "+name+" is valid")
	return nil
}

// filterFromString takes a filter in string form and returns a Filter
// No validation is performed. The filter is returned.
func filterFromString(filter string, filterName string) Filter {
	return Filter{
		Name:      filterName,
		Filter:    filter,
		CreatedAt: time.Now().UTC().String(), // The current time and date in UTC timezone
	}
}

// printFilters takes a list of filters, a writer, an output format, and a
// debug flag. It prints the filters to the writer in the output format
// specified. Returns an error if there is one.
func printFilters(filters *FilterList, writer io.Writer, outputFormat string, debugFlag bool) error {
	// switch on output type
	switch outputFormat {
	case "yaml":
		return printFiltersYaml(filters, writer, debugFlag)
	case "json":
		return printFiltersJson(filters, writer, debugFlag)
	}
	return nil
}

// printFiltersJson takes a list of filters, a writer, and a debug flag. It
// prints the filters to the writer in JSON format. Returns an error if there
// is one.
func printFiltersJson(filters *FilterList, writer io.Writer, debugFlag bool) error {
	encoder := json.NewEncoder(writer)
	encoder.SetIndent("", "  ")
	err := encoder.Encode(filters)
	if err != nil {
		debug.DebugMessage(debugFlag, err.Error())
		return err
	}
	return nil
}

// printFiltersYaml takes a list of filters, a writer, and a debug flag. It
// prints the filters to the writer in YAML format. Returns an error if there
// is one.
func printFiltersYaml(filters *FilterList, writer io.Writer, debugFlag bool) error {
	err := yaml.NewEncoder(writer).Encode(filters)
	if err != nil {
		debug.DebugMessage(debugFlag, err.Error())
		return err
	}
	return nil
}

// GetFilter takes a filter name and a debug flag, and returns that filter and
// an error if one was encountered. Returns error if named filter was not found.
func GetFilter(filterName string, debugFlag bool) (Filter, error) {
	filters, err := GetFilters(common.BranchName, debugFlag)
	if err != nil {
		return Filter{}, err
	}
	return filters.Filters[filterName], nil
}

// FilterTickets takes a list of tickets, a filter name, and a debug flag. It
// returns a list of tickets that match the filter. Returns an error if there
// is one.
func FilterTickets(tickets []Ticket, filterName string, debugFlag bool) (*[]Ticket, error) {
	// Get the filter
	filter, err := GetFilter(filterName, debugFlag)
	if err != nil {
		return nil, err
	}

	// Parse the filter
	queryObj, err := gojq.Parse(filter.Filter)
	if err != nil {
		return nil, fmt.Errorf("Error parsing filter: " + err.Error())
	}

	// Convert []Ticket into []map[string]interface{} for gojq
	var listOfTickets []map[string]interface{}
	ticketsJSON, err := json.Marshal(tickets)
	if err != nil {
		return nil, err
	}
	fmt.Println("The list of tickets as JSON: " + string(ticketsJSON))
	err = json.Unmarshal(ticketsJSON, &listOfTickets)
	if err != nil {
		return nil, err
	}
	fmt.Println("The length of listOfTickets is " + strconv.Itoa(len(listOfTickets)))

	// Apply the filter
	iter := queryObj.Run(listOfTickets)
	var iterTicket Ticket
	var filteredTickets []Ticket
	for {
		result, ok := iter.Next()
		if !ok {
			break
		}
		if err, ok := result.(error); ok {
			return nil, fmt.Errorf("Error applying filter: " + err.Error())
		}
		// Turn result into JSON and then into Ticket
		resultJSON, err := json.Marshal(result)
		if err != nil {
			return nil, err
		}
		fmt.Println("Trying to unmarshal: " + string(resultJSON))
		err = json.Unmarshal(resultJSON, &iterTicket)
		if err != nil {
			return nil, err
		}
		filteredTickets = append(filteredTickets, iterTicket)
	}

	return &filteredTickets, nil
}

// GetCurrentFilter takes a debug flag and returns the name of the current
// filter. Returns an error if there is one. This function may return an empty
// string if the current filter has not yet been set.
func GetCurrentFilter(debugFlag bool) (string, error) {
	debug.DebugMessage(debugFlag, "GetCurrentFilter() start")
	filters, err := GetFilters(common.BranchName, debugFlag)
	if err.Error() == "the path 'filters.json' does not exist in the given tree" {
		return "", nil
	}
	if err != nil {
		return "", err
	}
	return filters.CurrentFilter, nil
}
