package common

import (
	"bytes"
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"

	cp "github.com/otiai10/copy"
	git "gopkg.in/src-d/go-git.v4"
	"gopkg.in/src-d/go-git.v4/plumbing"
	"gotest.tools/v3/assert"
)

func PrepareMultiFixtures(t *testing.T, name string, deps []string) RepoConfig {
	newTestDataPath, err := ioutil.TempDir(os.TempDir(), "mutli_fixture")
	assert.NilError(t, err)

	for _, name := range deps {
		fixturePath := filepath.Join("testdata", name)
		newPath := filepath.Join(newTestDataPath, name)
		err = cp.Copy(fixturePath, newPath)
		assert.NilError(t, err)

		err = os.Rename(filepath.Join(newPath, ".gitted"), filepath.Join(newPath, ".git"))
		assert.NilError(t, err)
	}

	newRepoConfig := PrepareFixture(t, name)
	FixFixtureGitConfig(t, newRepoConfig.RepoPath, newTestDataPath)

	return newRepoConfig
}

func FixFixtureGitConfig(t *testing.T, newRepoPath string, testDataPath string) {
	// Fix remote paths
	dotGitPath := filepath.Join(newRepoPath, ".git")
	gitConfigFilePath := filepath.Join(dotGitPath, "config")
	input, err := ioutil.ReadFile(gitConfigFilePath)
	assert.NilError(t, err)

	output := bytes.Replace(input, []byte("$TESTDATA$"), []byte(testDataPath), -1)

	err = ioutil.WriteFile(gitConfigFilePath, output, 0666)
	assert.NilError(t, err)
}

func Test_SimpleFetch(t *testing.T) {
	repoConfig := PrepareMultiFixtures(t, "simple_fetch", []string{"multiple_file_change"})

	err := fetch(repoConfig)
	assert.NilError(t, err)

	r, err := git.PlainOpen(repoConfig.RepoPath)
	assert.NilError(t, err)

	head, err := r.Head()
	assert.NilError(t, err)

	assert.Equal(t, head.Hash(), plumbing.NewHash("28cc969d97ddb7640f5e1428bbc8f2947d1ffd57"))

	ref, err := r.Reference(plumbing.NewRemoteReferenceName("origin1", "master"), true)
	assert.NilError(t, err)

	assert.Equal(t, ref.Hash(), plumbing.NewHash("7058b6b292ee3d1382670334b5f29570a1117ef1"))
}
