/*
Giticket is a simple tool for tracking your bug tickets in a branch of your git repository.
It uses a git branch called 'giticket' to track your issues in the same repository as your code, letting you work on both entirely offline and then push your changes up later.

Usage:

	giticket {action} [flags]

	One action is accepted
	Zero or more parameters are accepted, parameters include: -help, -version
	giticket -help            will print this message
	giticket {action} -help   will print the help for that command
	giticket -version         will print the version of giticket

	Available Actions:
	-  comment
	-  create
	-  delete
	-  init
	-  label
	-  list
	-  priority
	-  severity
	-  show
	-  status
*/
package main

import (
	"github.com/jeffwelling/giticket/internal/cli"
)

func main() {
	cli.Exec()
}
