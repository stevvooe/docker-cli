package plugin // import "docker.com/cli/v28/cli/command/plugin"

import (
	"context"
	"fmt"
	"strings"

	"github.com/distribution/reference"
	"github.com/docker/cli/v28/cli"
	"github.com/docker/cli/v28/cli/command"
	"github.com/docker/cli/v28/cli/command/image"
	"github.com/docker/cli/v28/cli/internal/jsonstream"
	"github.com/docker/docker/api/types"
	registrytypes "github.com/docker/docker/api/types/registry"
	"github.com/docker/docker/registry"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

type pluginOptions struct {
	remote          string
	localName       string
	grantPerms      bool
	disable         bool
	args            []string
	skipRemoteCheck bool
	untrusted       bool
}

func loadPullFlags(dockerCli command.Cli, opts *pluginOptions, flags *pflag.FlagSet) {
	flags.BoolVar(&opts.grantPerms, "grant-all-permissions", false, "Grant all permissions necessary to run the plugin")
	command.AddTrustVerificationFlags(flags, &opts.untrusted, dockerCli.ContentTrustEnabled())
}

func newInstallCommand(dockerCli command.Cli) *cobra.Command {
	var options pluginOptions
	cmd := &cobra.Command{
		Use:   "install [OPTIONS] PLUGIN [KEY=VALUE...]",
		Short: "Install a plugin",
		Args:  cli.RequiresMinArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			options.remote = args[0]
			if len(args) > 1 {
				options.args = args[1:]
			}
			return runInstall(cmd.Context(), dockerCli, options)
		},
	}

	flags := cmd.Flags()
	loadPullFlags(dockerCli, &options, flags)
	flags.BoolVar(&options.disable, "disable", false, "Do not enable the plugin on install")
	flags.StringVar(&options.localName, "alias", "", "Local name for plugin")
	return cmd
}

func buildPullConfig(ctx context.Context, dockerCli command.Cli, opts pluginOptions, cmdName string) (types.PluginInstallOptions, error) {
	// Names with both tag and digest will be treated by the daemon
	// as a pull by digest with a local name for the tag
	// (if no local name is provided).
	ref, err := reference.ParseNormalizedNamed(opts.remote)
	if err != nil {
		return types.PluginInstallOptions{}, err
	}

	repoInfo, err := registry.ParseRepositoryInfo(ref)
	if err != nil {
		return types.PluginInstallOptions{}, err
	}

	remote := ref.String()

	_, isCanonical := ref.(reference.Canonical)
	if !opts.untrusted && !isCanonical {
		ref = reference.TagNameOnly(ref)
		nt, ok := ref.(reference.NamedTagged)
		if !ok {
			return types.PluginInstallOptions{}, errors.Errorf("invalid name: %s", ref.String())
		}

		trusted, err := image.TrustedReference(ctx, dockerCli, nt)
		if err != nil {
			return types.PluginInstallOptions{}, err
		}
		remote = reference.FamiliarString(trusted)
	}

	authConfig := command.ResolveAuthConfig(dockerCli.ConfigFile(), repoInfo.Index)
	encodedAuth, err := registrytypes.EncodeAuthConfig(authConfig)
	if err != nil {
		return types.PluginInstallOptions{}, err
	}

	options := types.PluginInstallOptions{
		RegistryAuth:          encodedAuth,
		RemoteRef:             remote,
		Disabled:              opts.disable,
		AcceptAllPermissions:  opts.grantPerms,
		AcceptPermissionsFunc: acceptPrivileges(dockerCli, opts.remote),
		PrivilegeFunc:         command.RegistryAuthenticationPrivilegedFunc(dockerCli, repoInfo.Index, cmdName),
		Args:                  opts.args,
	}
	return options, nil
}

func runInstall(ctx context.Context, dockerCLI command.Cli, opts pluginOptions) error {
	var localName string
	if opts.localName != "" {
		aref, err := reference.ParseNormalizedNamed(opts.localName)
		if err != nil {
			return err
		}
		if _, ok := aref.(reference.Canonical); ok {
			return errors.Errorf("invalid name: %s", opts.localName)
		}
		localName = reference.FamiliarString(reference.TagNameOnly(aref))
	}

	options, err := buildPullConfig(ctx, dockerCLI, opts, "plugin install")
	if err != nil {
		return err
	}
	responseBody, err := dockerCLI.Client().PluginInstall(ctx, localName, options)
	if err != nil {
		if strings.Contains(err.Error(), "(image) when fetching") {
			return errors.New(err.Error() + " - Use \"docker image pull\"")
		}
		return err
	}
	defer responseBody.Close()
	if err := jsonstream.Display(ctx, responseBody, dockerCLI.Out()); err != nil {
		return err
	}
	_, _ = fmt.Fprintln(dockerCLI.Out(), "Installed plugin", opts.remote) // todo: return proper values from the API for this result
	return nil
}

func acceptPrivileges(dockerCLI command.Cli, name string) func(ctx context.Context, privileges types.PluginPrivileges) (bool, error) {
	return func(ctx context.Context, privileges types.PluginPrivileges) (bool, error) {
		_, _ = fmt.Fprintf(dockerCLI.Out(), "Plugin %q is requesting the following privileges:\n", name)
		for _, privilege := range privileges {
			_, _ = fmt.Fprintf(dockerCLI.Out(), " - %s: %v\n", privilege.Name, privilege.Value)
		}
		return command.PromptForConfirmation(ctx, dockerCLI.In(), dockerCLI.Out(), "Do you grant the above permissions?")
	}
}
