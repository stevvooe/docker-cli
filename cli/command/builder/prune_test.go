package builder // import "docker.com/cli/v28/cli/command/builder"

import (
	"context"
	"errors"
	"io"
	"testing"

	"github.com/docker/cli/v28/internal/test"
	"github.com/docker/docker/api/types"
)

func TestBuilderPromptTermination(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	t.Cleanup(cancel)

	cli := test.NewFakeCli(&fakeClient{
		builderPruneFunc: func(ctx context.Context, opts types.BuildCachePruneOptions) (*types.BuildCachePruneReport, error) {
			return nil, errors.New("fakeClient builderPruneFunc should not be called")
		},
	})
	cmd := NewPruneCommand(cli)
	cmd.SetOut(io.Discard)
	cmd.SetErr(io.Discard)
	test.TerminatePrompt(ctx, t, cmd, cli)
}
