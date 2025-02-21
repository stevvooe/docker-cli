package swarm // import "docker.com/cli/v28/cli/command/stack/swarm"

import (
	"context"
	"fmt"

	"github.com/docker/cli/v28/cli/command"
	"github.com/docker/cli/v28/cli/command/idresolver"
	"github.com/docker/cli/v28/cli/command/stack/options"
	"github.com/docker/cli/v28/cli/command/task"
	"github.com/docker/docker/api/types"
)

// RunPS is the swarm implementation of docker stack ps
func RunPS(ctx context.Context, dockerCli command.Cli, opts options.PS) error {
	filter := getStackFilterFromOpt(opts.Namespace, opts.Filter)

	client := dockerCli.Client()
	tasks, err := client.TaskList(ctx, types.TaskListOptions{Filters: filter})
	if err != nil {
		return err
	}

	if len(tasks) == 0 {
		return fmt.Errorf("nothing found in stack: %s", opts.Namespace)
	}

	format := opts.Format
	if len(format) == 0 {
		format = task.DefaultFormat(dockerCli.ConfigFile(), opts.Quiet)
	}

	return task.Print(ctx, dockerCli, tasks, idresolver.New(client, opts.NoResolve), !opts.NoTrunc, opts.Quiet, format)
}
