package git

import (
	"fmt"
	"os/exec"
	"strings"
)

// CommandExecutor defines the interface for executing commands.
type CommandExecutor interface {
	RunCommand(name string, args ...string) (string, error)
}

// OSCommandExecutor implements CommandExecutor using os/exec.
type OSCommandExecutor struct{}

// RunCommand executes a command using os/exec.Command.
func (e *OSCommandExecutor) RunCommand(name string, args ...string) (string, error) {
	cmd := exec.Command(name, args...)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(string(output)), nil
}

// Run a git command and returns its output.
func Run(executor CommandExecutor, args ...string) (string, error) {
	return executor.RunCommand("git", args...)
}

// GetGitRoot attempts to find the Git repository root.
func GetGitRoot(executor CommandExecutor) (string, error) {
	output, err := executor.RunCommand("git", "rev-parse", "--show-toplevel")
	if err != nil {
		return "", fmt.Errorf("not a git repository: %w", err)
	}
	return output, nil
}
