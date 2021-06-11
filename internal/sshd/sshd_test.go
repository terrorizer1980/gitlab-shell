package sshd

import (
	"context"
	"testing"
	"path"
	"gitlab.com/gitlab-org/gitlab-shell/internal/testhelper"

	"gitlab.com/gitlab-org/gitlab-shell/client/testserver"
	"github.com/stretchr/testify/require"
	"gitlab.com/gitlab-org/gitlab-shell/internal/config"
)

func TestRunIsCancelable(t *testing.T) {
	testDirCleanup, err := testhelper.PrepareTestRootDir()
	require.NoError(t, err)
	defer testDirCleanup()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	url := testserver.StartSocketHttpServer(t, nil)

	errCh := make(chan error, 1)

	go func() {
		cfg := &config.Config{
			GitlabUrl: url,
			Server: config.ServerConfig{
				HostKeyFiles: []string{path.Join(testhelper.TestRoot, "certs/valid/server.key")},
			},
		}

		errCh <- Run(ctx, cfg)
	}()

	cancel()
	require.NoError(t, <-errCh)
}
