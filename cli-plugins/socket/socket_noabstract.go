//go:build !windows && !linux

package socket // import "docker.com/cli/v28/cli-plugins/socket"

import (
	"os"
	"path/filepath"
)

func socketName(basename string) string {
	// Because abstract sockets are unavailable, use a socket path in the
	// system temporary directory.
	return filepath.Join(os.TempDir(), basename)
}
