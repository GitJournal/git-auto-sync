package main

import (
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"testing"

	cp "github.com/otiai10/copy"
)

func prepareFixture(name string) string {
	newRepoPath, err := ioutil.TempDir(os.TempDir(), name)
	if err != nil {
		log.Fatal(err)
	}

	fixturePath := filepath.Join("testdata", name)
	err = cp.Copy(fixturePath, newRepoPath)
	if err != nil {
		log.Fatal(err)
	}

	err = os.Rename(filepath.Join(newRepoPath, ".gitted"), filepath.Join(newRepoPath, ".git"))
	if err != nil {
		log.Fatal(err)
	}

	return newRepoPath
}

func Test_NoChanges(t *testing.T) {
	repoPath := prepareFixture("no_changes")

	err := commit(repoPath)
	if err != nil {
		t.Fatal(err)
	}
}

// Commit Tests
// 2. One file changes
// 3. Multiple files changed
// 4. A file in a subdirectory changed
