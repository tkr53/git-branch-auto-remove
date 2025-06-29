package git

import (
	"os/exec"
	"testing"

	"bou.ke/monkey"
	"github.com/stretchr/testify/assert"
)

func TestRun(t *testing.T) {
	// Test case 1: Successful command execution
	monkey.Patch(exec.Command, func(name string, arg ...string) *exec.Cmd {
		return exec.Command("echo", "test output")
	})
	out, err := Run("test")
	assert.NoError(t, err)
	assert.Equal(t, "test output", out)

	// Test case 2: Command execution with error
	monkey.Patch(exec.Command, func(name string, arg ...string) *exec.Cmd {
		return exec.Command("false") // Command that always fails
	})
	out, err = Run("test")
	assert.Error(t, err)
	assert.Empty(t, out)

	monkey.UnpatchAll()
}
