package bash

import (
	"fmt"
	"os/exec"
)

func executeBash(bashCommand string) error {

	cmd := exec.Command("sh", "-c", bashCommand)
	combinedOutput, err := cmd.CombinedOutput()

	if err != nil {
		return fmt.Errorf("command execution failed: %v\noutput: %s", err, string(combinedOutput))
	}

	return nil
}
