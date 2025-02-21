package plugin // import "docker.com/cli/v28/e2e/plugin"

import (
	"context"
	"testing"

	"github.com/docker/cli/v28/e2e/internal/fixtures"
	"github.com/docker/cli/v28/e2e/testutils"
	"github.com/docker/cli/v28/internal/test/environment"
	"github.com/docker/docker/api/types/versions"
	"gotest.tools/v3/icmd"
	"gotest.tools/v3/skip"
)

const registryPrefix = "registry:5000"

func TestInstallWithContentTrust(t *testing.T) {
	// TODO(krissetto): remove this skip once the fix (see https://github.com/moby/moby/pull/47299) is deployed to moby versions < 25
	skip.If(t, versions.LessThan(environment.DaemonAPIVersion(t), "1.44"))
	skip.If(t, environment.SkipPluginTests())

	const pluginName = registryPrefix + "/plugin-content-trust"

	dir := fixtures.SetupConfigFile(t)
	defer dir.Remove()

	ctx, cancel := context.WithCancel(context.Background())
	t.Cleanup(cancel)

	pluginDir := testutils.SetupPlugin(t, ctx)
	t.Cleanup(pluginDir.Remove)

	icmd.RunCommand("docker", "plugin", "create", pluginName, pluginDir.Path()).Assert(t, icmd.Success)
	result := icmd.RunCmd(icmd.Command("docker", "plugin", "push", pluginName),
		fixtures.WithConfig(dir.Path()),
		fixtures.WithTrust,
		fixtures.WithNotary,
		fixtures.WithPassphrase("foo", "bar"),
	)
	result.Assert(t, icmd.Expected{
		Out: "Signing and pushing trust metadata",
	})

	icmd.RunCommand("docker", "plugin", "rm", "-f", pluginName).Assert(t, icmd.Success)

	result = icmd.RunCmd(icmd.Command("docker", "plugin", "install", "--grant-all-permissions", pluginName),
		fixtures.WithConfig(dir.Path()),
		fixtures.WithTrust,
		fixtures.WithNotary,
	)
	result.Assert(t, icmd.Expected{
		Out: "Installed plugin " + pluginName,
	})
}

func TestInstallWithContentTrustUntrusted(t *testing.T) {
	skip.If(t, environment.SkipPluginTests())

	dir := fixtures.SetupConfigFile(t)
	defer dir.Remove()

	result := icmd.RunCmd(icmd.Command("docker", "plugin", "install", "--grant-all-permissions", "tiborvass/sample-volume-plugin:latest"),
		fixtures.WithConfig(dir.Path()),
		fixtures.WithTrust,
		fixtures.WithNotary,
	)
	result.Assert(t, icmd.Expected{
		ExitCode: 1,
		Err:      "Error: remote trust data does not exist",
	})
}
