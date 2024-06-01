package ticket

import (
	"strconv"
	"time"

	git "github.com/jeffwelling/git2go/v37"
	"github.com/jeffwelling/giticket/pkg/common"
	"github.com/jeffwelling/giticket/pkg/debug"
	"github.com/jeffwelling/giticket/pkg/repo"
)

type Comment struct {
	ID      int
	Created int64
	Body    string
	Author  string
}

func HandleComment(
	branchName string,
	comment string,
	commentID int,
	ticketID int,
	deleteFlag bool,
	debugFlag bool,
) {
	debug.DebugMessage(debugFlag, "Opening git repository")
	thisRepo, err := git.OpenRepository(".")
	if err != nil {
		panic(err)
	}

	// Get author
	author := common.GetAuthor(thisRepo)

	tickets, err := GetListOfTickets(thisRepo, common.BranchName, debugFlag)
	if err != nil {
		panic(err)
	}
	t := FilterTicketsByID(tickets, ticketID)

	if deleteFlag {
		commentID := DeleteComment(&t, commentID, thisRepo, branchName, debugFlag)
		err := repo.Commit(&t, thisRepo, branchName, author, "Deleting comment "+commentID, debugFlag)
		if err != nil {
			panic(err)
		}
	} else {
		commentID := AddComment(&t, comment, thisRepo, branchName, debugFlag)
		err := repo.Commit(&t, thisRepo, branchName, author, "Adding comment "+commentID, debugFlag)
		if err != nil {
			panic(err)
		}
	}
}

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

func AddComment(t *Ticket, comment string, thisRepo *git.Repository, branchName string, debug bool) string {
	author := common.GetAuthor(thisRepo)
	newComment := Comment{
		ID:      t.NextCommentID,
		Created: time.Now().Unix(),
		Body:    comment,
		Author:  author.Name + " <" + author.Email + ">",
	}
	t.NextCommentID++
	t.Comments = append(t.Comments, newComment)
	// Return a string of the ticketID-CommentID
	return strconv.Itoa(t.ID) + "-" + strconv.Itoa(newComment.ID)
}
