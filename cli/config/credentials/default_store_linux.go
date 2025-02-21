package credentials // import "docker.com/cli/v28/cli/config/credentials"

import (
	"os/exec"
)

func defaultCredentialsStore() string {
	if _, err := exec.LookPath("pass"); err == nil {
		return "pass"
	}

	return "secretservice"
}
