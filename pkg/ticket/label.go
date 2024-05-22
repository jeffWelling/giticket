package ticket

import git "github.com/jeffwelling/git2go/v37"

// DeleteLabel() takes a pointer to a ticket, a label, a git repo, a branch
// name, and a debug flag. It deletes the label from the ticket and returns an
// error.
func DeleteLabel(t *Ticket, label string, repo *git.Repository, branchName string, debug bool) string {
	updatedLabels := []string{}
	for _, l := range t.Labels {
		if l != label {
			updatedLabels = append(updatedLabels, l)
		}
	}
	t.Labels = updatedLabels
	return label
}

// AddLabel() takes a pointer to a ticket, a label, a git repo, a branch name,
// and a debug flag. It adds the label to the ticket and returns an error.
func AddLabel(t *Ticket, label string, repo *git.Repository, branchName string, debug bool) string {
	t.Labels = append(t.Labels, label)
	return label
}
