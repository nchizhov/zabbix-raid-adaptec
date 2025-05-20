package support

import (
	"bytes"
	"errors"
	"fmt"
	"os/exec"
)

func RunCommand(prg string, args ...string) (string, error) {
	var out bytes.Buffer
	cmd := exec.Command(prg, args...)
	cmd.Stdout = &out
	if err := cmd.Start(); err != nil {
		return "", err
	}
	if err := cmd.Wait(); err != nil {
		var exitError *exec.ExitError
		if errors.As(err, &exitError) {
			return "", fmt.Errorf("exit status: %d", exitError.ExitCode())
		}
		return "", err
	}
	return out.String(), nil
}
