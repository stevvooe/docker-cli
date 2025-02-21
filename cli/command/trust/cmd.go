package trust // import "docker.com/cli/v28/cli/command/trust"

import (
	"github.com/docker/cli/v28/cli"
	"github.com/docker/cli/v28/cli/command"
	"github.com/spf13/cobra"
)

// NewTrustCommand returns a cobra command for `trust` subcommands
func NewTrustCommand(dockerCli command.Cli) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "trust",
		Short: "Manage trust on Docker images",
		Args:  cli.NoArgs,
		RunE:  command.ShowHelp(dockerCli.Err()),
	}
	cmd.AddCommand(
		newRevokeCommand(dockerCli),
		newSignCommand(dockerCli),
		newTrustKeyCommand(dockerCli),
		newTrustSignerCommand(dockerCli),
		newInspectCommand(dockerCli),
	)
	return cmd
}
