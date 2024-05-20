package common

import (
	"fmt"
	"time"

	git "github.com/jeffwelling/git2go/v37"
)

func GetAuthor(repo *git.Repository) *git.Signature {
	// Load the configuration which merges global, system, and local configs
	cfg, err := repo.Config()
	if err != nil {
		fmt.Println("Error accessing config:", err)
		panic(err)
	}
	defer cfg.Free()

	// Retrieve user's name and email from the configuration
	name, err := cfg.LookupString("user.name")
	if err != nil {
		fmt.Println("Error retrieving user name:", err)
		panic(err)
	}
	email, err := cfg.LookupString("user.email")
	if err != nil {
		fmt.Println("Error retrieving user email:", err)
		panic(err)
	}

	// Create a new commit on the branch
	author := &git.Signature{
		Name:  name,
		Email: email,
		When:  time.Now(),
	}

	return author
}
