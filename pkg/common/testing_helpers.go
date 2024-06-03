package common

import (
	"os"
	"testing"
	"time"

	git "github.com/jeffwelling/git2go/v37"
)

// Create a temporary directory by calling testing.T.TempDir()
// then cd into the directory for the purposes of running tests
// from within that directory.
func UseTempDir(t *testing.T) {
	// Create a temporary directory
	tempDir := t.TempDir()

	// Cd into tempdir from t.tempDir
	err := os.Chdir(tempDir)
	if err != nil {
		// Fail the test
		t.Fatal(err)
	}
}

// Create a git repository in the current directory
// for use with testing, returns the commit ID
func InitGit(t *testing.T) *git.Oid {

	// https://libgit2.org/libgit2/#HEAD/group/repository/git_repository_init
	// `false` here means the .git directory will be created
	repo, err := git.InitRepository(".", false)
	if err != nil {
		t.Fatal(err)
	}

	// Just write a string to the repo so the repository isn't totally bare
	blobOid, err := repo.CreateBlobFromBuffer([]byte("words. big words. yuge words. the best words."))
	if err != nil {
		t.Fatal(err)
	}

	// Create a treeBuilder, for use as the tree in the commit
	treeBuilderRoot, err := repo.TreeBuilder()
	if err != nil {
		t.Fatal(err)
	}

	// Add the blobOid to the tree with the filename "testFile"
	err = treeBuilderRoot.Insert("testFile", blobOid, git.FilemodeBlob)
	if err != nil {
		t.Fatal(err)
	}

	rootTreeID, err := treeBuilderRoot.Write()
	if err != nil {
		t.Fatal(err)
	}
	rootTree, err := repo.LookupTree(rootTreeID)
	if err != nil {
		t.Fatal(err)
	}

	// Create the git author
	author := git.Signature{
		Name:  "test user",
		Email: "test@example.com",
		When:  time.Now(),
	}

	// Commit
	commitID, err := repo.CreateCommit("refs/heads/main", &author, &author, "test commit", rootTree)
	if err != nil {
		t.Fatal(err)
	}

	return commitID
}
