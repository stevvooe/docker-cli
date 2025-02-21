package cliplugins // import "docker.com/cli/v28/e2e/cli-plugins"

import (
	"fmt"
	"os"
	"testing"

	"github.com/docker/cli/v28/cli/config"
	"github.com/docker/cli/v28/cli/config/configfile"
	"gotest.tools/v3/assert"
	"gotest.tools/v3/fs"
	"gotest.tools/v3/icmd"
)

func prepare(t *testing.T) (func(args ...string) icmd.Cmd, *configfile.ConfigFile, func()) {
	t.Helper()
	cfg := fs.NewDir(t, "plugin-test",
		fs.WithFile("config.json", fmt.Sprintf(`{"cliPluginsExtraDirs": [%q]}`, os.Getenv("DOCKER_CLI_E2E_PLUGINS_EXTRA_DIRS"))),
	)
	run := func(args ...string) icmd.Cmd {
		return icmd.Command("docker", append([]string{"--config", cfg.Path()}, args...)...)
	}
	cleanup := func() {
		cfg.Remove()
	}
	cfgfile, err := config.Load(cfg.Path())
	assert.NilError(t, err)

	return run, cfgfile, cleanup
}
