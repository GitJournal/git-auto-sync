package main

import (
	"testing"

	git "gopkg.in/src-d/go-git.v4"
	"gopkg.in/src-d/go-git.v4/plumbing"
	"gotest.tools/v3/assert"
)

// No new commits on remote or local
func Test_RebaseNothing(t *testing.T) {
	repoPath := PrepareMultiFixtures(t, "rebase_nothing", []string{"rebase_parent"})

	err := rebase(repoPath)
	assert.NilError(t, err)

	r, err := git.PlainOpen(repoPath)
	assert.NilError(t, err)

	head, err := r.Head()
	assert.NilError(t, err)

	assert.Equal(t, head.Hash(), plumbing.NewHash("28cc969d97ddb7640f5e1428bbc8f2947d1ffd57"))
}

// New commits on local
func Test_RebaseLocalCommits(t *testing.T) {
	repoPath := PrepareMultiFixtures(t, "rebase_local_commits", []string{"rebase_parent"})

	err := rebase(repoPath)
	assert.NilError(t, err)

	r, err := git.PlainOpen(repoPath)
	assert.NilError(t, err)

	head, err := r.Head()
	assert.NilError(t, err)

	assert.Equal(t, head.Hash(), plumbing.NewHash("7fc438e0c9cc4f58178a1efe8521e52f0f8ee688"))
}

// New commits on remote
func Test_RebaseRemoteCommits(t *testing.T) {
	repoPath := PrepareMultiFixtures(t, "rebase_remote_commits", []string{"rebase_parent"})

	err := rebase(repoPath)
	assert.NilError(t, err)

	r, err := git.PlainOpen(repoPath)
	assert.NilError(t, err)

	head, err := r.Head()
	assert.NilError(t, err)

	assert.Equal(t, head.Hash(), plumbing.NewHash("ccda8f2e691aa416791a10afc74ccdbd1cb419fe"))
}

// New commits on both, no conflict
func Test_RebaseBothCommitsNoConflict(t *testing.T) {
	repoPath := PrepareMultiFixtures(t, "rebase_both_commits", []string{"rebase_parent"})

	err := rebase(repoPath)
	assert.NilError(t, err)

	r, err := git.PlainOpen(repoPath)
	assert.NilError(t, err)

	head, err := r.Head()
	assert.NilError(t, err)

	assert.Check(t, head.Hash() != plumbing.NewHash("ccda8f2e691aa416791a10afc74ccdbd1cb419fe"))
	assert.Check(t, head.Hash() != plumbing.NewHash("5779561afa9d074ae8d20974861c54757429aca9"))
	assert.Check(t, head.Hash() != plumbing.NewHash("7fc438e0c9cc4f58178a1efe8521e52f0f8ee688"))
}

// New commits on both, some kind of conflict
func Test_RebaseBothCommitsConflict(t *testing.T) {
	repoPath := PrepareMultiFixtures(t, "rebase_both_commits_conflict", []string{"rebase_parent"})

	r, err := git.PlainOpen(repoPath)
	assert.NilError(t, err)

	origHead, err := r.Head()
	assert.NilError(t, err)

	err = rebase(repoPath)
	assert.NilError(t, err)

	newHead, err := r.Head()
	assert.NilError(t, err)

	assert.Check(t, newHead.Hash() == origHead.Hash())
}
