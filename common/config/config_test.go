package config

import (
	"io/ioutil"
	"os"
	"testing"

	"github.com/kirsle/configdir"
	"gotest.tools/v3/assert"
)

func setup(t *testing.T, name string) {
	newConfigDir, err := ioutil.TempDir(os.TempDir(), "config_"+name)
	assert.NilError(t, err)
	t.Setenv("XDG_CONFIG_HOME", newConfigDir)
	t.Setenv("HOME", newConfigDir)

	configdir.Refresh()
}

func Test_SimpleWriteRead(t *testing.T) {
	setup(t, "SimpleWriteRead")

	c := &Config{Repos: []string{"/home/xyz/hello"}}
	err := Write(c)
	assert.NilError(t, err)

	c2, err := Read()
	assert.NilError(t, err)

	assert.DeepEqual(t, c, c2)
}

func Test_ReadEmpty(t *testing.T) {
	setup(t, "ReadEmpty")

	c, err := Read()
	assert.NilError(t, err)
	assert.Assert(t, len(c.Repos) == 0)
}
