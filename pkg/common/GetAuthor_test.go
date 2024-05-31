package common

import (
	"os"
	"testing"
)

func TestGetAuthor(t *testing.T) {
	// Cd into tempdir from t.tempDir
	err := os.Chdir(t.TempDir())
	if err != nil {
		// Fail the test
		t.Fatal(err)
	}

}
