package common

import (
	"os"
	"testing"
	"time"

	git "github.com/jeffwelling/git2go/v37"
)

// UseTempDir takes a pointer to a testing T, it gets the temp directory for
// this test from T.TempDir(), and cd's into it. If there's an error, it
// fails the test.
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

// InitGit takes a pointer to a testing T, it initializes a git repository in
// the current directory and returns a pointer to the first commit ID. If there
// is an error, it fails the test.
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
