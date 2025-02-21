package main // import "docker.com/cli/v28/cmd/docker"

import (
	"github.com/docker/cli/v28/cli/context/store"
	"github.com/spf13/cobra"
)

type contextStoreProvider interface {
	ContextStore() store.Store
}

func completeContextNames(dockerCLI contextStoreProvider) func(*cobra.Command, []string, string) ([]string, cobra.ShellCompDirective) {
	return func(*cobra.Command, []string, string) ([]string, cobra.ShellCompDirective) {
		names, _ := store.Names(dockerCLI.ContextStore())
		return names, cobra.ShellCompDirectiveNoFileComp
	}
}

var logLevels = []string{"debug", "info", "warn", "error", "fatal", "panic"}

func completeLogLevels(*cobra.Command, []string, string) ([]string, cobra.ShellCompDirective) {
	return cobra.FixedCompletions(logLevels, cobra.ShellCompDirectiveNoFileComp)(nil, nil, "")
}
