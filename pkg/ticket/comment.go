package ticket

import (
	"strconv"
	"time"

	"github.com/jeffWelling/giticket/pkg/common"
	git "github.com/jeffwelling/git2go/v37"
)

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
