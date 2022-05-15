package common

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"

	cp "github.com/otiai10/copy"
	git "gopkg.in/src-d/go-git.v4"
	"gopkg.in/src-d/go-git.v4/plumbing"

	"gotest.tools/v3/assert"
)

func PrepareFixture(t *testing.T, name string) string {
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
	repoPath := PrepareFixture(t, "no_changes")

	err := commit(repoPath)
	assert.NilError(t, err)

	r, err := git.PlainOpen(repoPath)
	assert.NilError(t, err)

	head, err := r.Head()
	assert.NilError(t, err)

	assert.Equal(t, head.Hash(), plumbing.NewHash("28cc969d97ddb7640f5e1428bbc8f2947d1ffd57"))
}

func HasHeadCommit(t *testing.T, repoPath string, hash string, msg string) {
	r, err := git.PlainOpen(repoPath)
	assert.NilError(t, err)

	head, err := r.Head()
	assert.NilError(t, err)

	assert.Assert(t, head.Hash() != plumbing.NewHash(hash))

	commit, err := r.CommitObject(head.Hash())
	assert.NilError(t, err)

	parent, err := commit.Parent(0)
	assert.NilError(t, err)
	assert.Equal(t, parent.ID(), plumbing.NewHash(hash))
	assert.Equal(t, commit.Message, msg)
}

func Test_NewFile(t *testing.T) {
	repoPath := PrepareFixture(t, "new_file")

	err := commit(repoPath)
	assert.NilError(t, err)

	HasHeadCommit(t, repoPath, "28cc969d97ddb7640f5e1428bbc8f2947d1ffd57", "?? 2.md\n")
}

func Test_OneFileChange(t *testing.T) {
	repoPath := PrepareFixture(t, "one_file_change")

	err := commit(repoPath)
	assert.NilError(t, err)

	HasHeadCommit(t, repoPath, "28cc969d97ddb7640f5e1428bbc8f2947d1ffd57", " M 1.md\n")
}

func Test_VimSwapFile(t *testing.T) {
	repoPath := PrepareFixture(t, "vim_swap_file")

	err := commit(repoPath)
	assert.NilError(t, err)

	r, err := git.PlainOpen(repoPath)
	assert.NilError(t, err)

	head, err := r.Head()
	assert.NilError(t, err)

	assert.Equal(t, head.Hash(), plumbing.NewHash("28cc969d97ddb7640f5e1428bbc8f2947d1ffd57"))
}

func Test_MultipleFileChange(t *testing.T) {
	repoPath := PrepareFixture(t, "multiple_file_change")

	err := commit(repoPath)
	assert.NilError(t, err)

	HasHeadCommit(t, repoPath, "7058b6b292ee3d1382670334b5f29570a1117ef1", " M 1.md\n D dirA/2.md\n?? dirB/3.md\n")
}
