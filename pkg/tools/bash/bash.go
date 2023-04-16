package bash

import (
	"os/exec"
	"strings"
)

func executeBash(bashCommand string) error {

	parts := strings.Fields(bashCommand)

	// Execute the command
	cmd := exec.Command(parts[0], parts[1:]...)
	if _, err := cmd.CombinedOutput(); err != nil {
		return err
	}

	return nil
}
