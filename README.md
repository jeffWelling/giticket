# Giticket

Store your tickets alongside your code under a git branch using this golang tool.

## Intro

Giticket is a ticketing tool written in golang that allows you to store tickets beside your code under a branch in git, it is a spiritual successor to ticgit and ticgit-ng which were written in ruby.

## Install

```bash
go install github.com/jeffwelling/giticket/cmd/giticket@latest
```

This will install giticket to $GOPATH/bin, make sure it is in your $PATH.

## Usage

```bash
# Initialize giticket to create the ticket branch
$ giticket init

# Create a ticket
$ giticket create --title 'My first ticket' --comments '[{"Body": "First comment", "Author": "John Smith <jsmith@example.com>"}]' --description "This is an awesome description." --labels "bugfix,ux"

# List tickets
$ giticket list
ID  | Title                | Severity  | Status
----------------------------------------------
1   | My first ticket      | 1         | new

# View ticket
$ giticket show --id 1
ID: 1
Title: My first ticket
Description: This is an awesome description.
Status: new
Severity: 1
Labels: bugfix, ux
Created: 2024-05-24 01:11:03 -0700 PDT
NextTicketID: 2
Comments:
    Comment ID: 1-1
    Created: 2024-05-24 01:11:03 -0700 PDT
    Author: John Smith <jsmith@example.com>
    Body: First comment

# Set status to in progress
$ giticket status --id 1 --status "in progress"

# Comment in the ticket
$ giticket comment --id 1 --comment "Inverted tardis polarity"

# Close the ticket
$ giticket status --id 1 --status "closed"

## TBD
# Delete the ticket
```
