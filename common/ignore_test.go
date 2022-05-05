package common

import (
	"testing"

	"gotest.tools/v3/assert"
)

func Test_SimpleIgnore(t *testing.T) {
	repoPath := PrepareFixture(t, "ignore")

	ignore, err := isFileIgnoredByGit(repoPath, "1.txt")
	assert.NilError(t, err)
	assert.Equal(t, ignore, true)

	ignore, err = isFileIgnoredByGit(repoPath, "2.md")
	assert.NilError(t, err)
	assert.Equal(t, ignore, false)
}
