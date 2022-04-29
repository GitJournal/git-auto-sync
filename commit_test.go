package main

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"

	cp "github.com/otiai10/copy"
	"gopkg.in/src-d/go-git.v4"
	"gopkg.in/src-d/go-git.v4/plumbing"

	"gotest.tools/v3/assert"
)

func prepareFixture(t *testing.T, name string) string {
	newRepoPath, err := ioutil.TempDir(os.TempDir(), name)
	assert.NilError(t, err)

	fixturePath := filepath.Join("testdata", name)
	err = cp.Copy(fixturePath, newRepoPath)
	assert.NilError(t, err)

	err = os.Rename(filepath.Join(newRepoPath, ".gitted"), filepath.Join(newRepoPath, ".git"))
	assert.NilError(t, err)

	return newRepoPath
}

func Test_NoChanges(t *testing.T) {
	repoPath := prepareFixture(t, "no_changes")

	err := commit(repoPath)
	assert.NilError(t, err)

	// Get the head commit
	r, err := git.PlainOpen(repoPath)
	assert.NilError(t, err)

	head, err := r.Head()
	assert.NilError(t, err)

	assert.Equal(t, head.Hash(), plumbing.NewHash("28cc969d97ddb7640f5e1428bbc8f2947d1ffd57"))
}

func Test_NewFile(t *testing.T) {
	repoPath := prepareFixture(t, "new_file")

	err := commit(repoPath)
	assert.NilError(t, err)

	// Get the head commit
	r, err := git.PlainOpen(repoPath)
	assert.NilError(t, err)

	head, err := r.Head()
	assert.NilError(t, err)

	assert.Assert(t, head.Hash() != plumbing.NewHash("28cc969d97ddb7640f5e1428bbc8f2947d1ffd57"))
}

// TODO
// * One file changes
// * Multiple files changed
// * A file in a subdirectory changed
// * The commit message is how we want it
