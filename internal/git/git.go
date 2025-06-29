package git

import (
	"os/exec"
	"strings"
)

// Run a git command and returns its output.
func Run(args ...string) (string, error) {
	cmd := exec.Command("git", args...)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(string(output)), nil
}
