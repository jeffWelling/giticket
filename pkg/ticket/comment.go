package ticket

import (
	"strconv"
	"time"

	git "github.com/jeffwelling/git2go/v37"
	"github.com/jeffwelling/giticket/pkg/common"
	"github.com/jeffwelling/giticket/pkg/debug"
	"github.com/jeffwelling/giticket/pkg/repo"
)

// A comment for a giticket ticket
// The fields must be filled in to populate details of the ticket
type Comment struct {
	ID      int
	Created int64
	Body    string
	Author  string
}

// HandleComment handles the addition or deletion of a comment for a ticket.
// It takes in the branch name, comment content, comment ID, ticket ID,
// delete flag, and debug flag as parameters. If the delete flag is true,
// it deletes the specified comment from the ticket. Otherwise, it adds
// the comment to the ticket. The function returns the comment ID and an error
// if there is one.
func HandleComment(
	branchName string,
	comment string,
	commentID int,
	ticketID int,
	deleteFlag bool,
	debugFlag bool,
) (string, error) {
	debug.DebugMessage(debugFlag, "Handling comment")
	debug.DebugMessage(debugFlag, "Opening git repository")
	thisRepo, err := git.OpenRepository(".")
	if err != nil {
		return "", err
	}

	// Get author
	debug.DebugMessage(debugFlag, "Getting author")
	author, err := common.GetAuthor(thisRepo)
	if err != nil {
		return "", err
	}

	debug.DebugMessage(debugFlag, "Getting tickets")
	tickets, err := GetListOfTickets(thisRepo, common.BranchName, debugFlag)
	if err != nil {
		return "", err
	}
	debug.DebugMessage(debugFlag, "Filtering tickets")
	t := FilterTicketsByID(tickets, ticketID)

	var fullCommentID string
	if deleteFlag {
		debug.DebugMessage(debugFlag, "Deleting comment")
		fullCommentID = DeleteComment(&t, commentID, thisRepo, branchName, debugFlag)
		err := repo.Commit(&t, thisRepo, branchName, author, "Deleting comment "+fullCommentID, debugFlag)
		if err != nil {
			return "", err
		}
	} else {
		debug.DebugMessage(debugFlag, "Adding comment")
		fullCommentID, err = AddComment(&t, comment, thisRepo, branchName, debugFlag)
		if err != nil {
			return "", err
		}
		err = repo.Commit(&t, thisRepo, branchName, author, "Adding comment "+fullCommentID, debugFlag)
		if err != nil {
			return "", err
		}
	}

	debug.DebugMessage(debugFlag, "Returning comment ID "+fullCommentID)
	return fullCommentID, nil
}

// DeleteComment removes the comment identified by commentID from ticket t under
// the repo and branchName provided with a debug flag. It returns a string
// representing the unique ID of the deleted comment.
func DeleteComment(t *Ticket, commentID int, repo *git.Repository, branchName string, debug bool) string {
	updatedComments := []Comment{}
	for _, comment := range t.Comments {
		if comment.ID != commentID {
			updatedComments = append(updatedComments, comment)
		}
	}
	t.Comments = updatedComments
	return strconv.Itoa(t.ID) + "-" + strconv.Itoa(commentID)
}

// AddComment adds a comment to ticket t under the repo and branchName provided
// with a debug flag. It returns a string representing the unique ID of the added
// comment and an error if there is was one.
func AddComment(t *Ticket, comment string, thisRepo *git.Repository, branchName string, debug bool) (string, error) {
	author, err := common.GetAuthor(thisRepo)
	if err != nil {
		return "", err
	}
	newComment := Comment{
		ID:      t.NextCommentID,
		Created: time.Now().Unix(),
		Body:    comment,
		Author:  author.Name + " <" + author.Email + ">",
	}
	t.NextCommentID++
	t.Comments = append(t.Comments, newComment)
	// Return a string of the ticketID-CommentID
	return strconv.Itoa(t.ID) + "-" + strconv.Itoa(newComment.ID), nil
}
