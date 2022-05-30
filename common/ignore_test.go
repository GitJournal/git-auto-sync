package common

import (
	"testing"

	"gotest.tools/v3/assert"
)

func Test_SimpleIgnore(t *testing.T) {
	repoPath := PrepareFixture(t, "ignore").RepoPath

	ignore, err := isFileIgnoredByGit(repoPath, "1.txt")
	assert.NilError(t, err)
	assert.Equal(t, ignore, true)

	ignore, err = isFileIgnoredByGit(repoPath, "2.md")
	assert.NilError(t, err)
	assert.Equal(t, ignore, false)
}

func Test_HiddenFilesIgnore(t *testing.T) {
	repoPath := PrepareFixture(t, "ignore").RepoPath

	ignore, err := ShouldIgnoreFile(repoPath, ".hidden")
	assert.NilError(t, err)
	assert.Equal(t, ignore, false)
}
