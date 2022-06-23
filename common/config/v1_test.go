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

func Test_SimpleWriteReadV1(t *testing.T) {
	setup(t, "SimpleWriteRead")

	c := &ConfigV1{
		Repos: []string{"/home/xyz/hello"},
		Envs:  []string{"SSH_AUTH_SOCK=/private/tmp/com.apple.launchd.74ZznY1v1F/Listeners"},
	}
	err := WriteV1(c)
	assert.NilError(t, err)

	c2, err := ReadV1()
	assert.NilError(t, err)

	assert.DeepEqual(t, c, c2)
}

func Test_ReadEmptyV1(t *testing.T) {
	setup(t, "ReadEmpty")

	c, err := ReadV1()
	assert.NilError(t, err)
	assert.Assert(t, len(c.Repos) == 0)
}
