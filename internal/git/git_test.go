package git

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

type MockCommandExecutor struct {
	RunCommandFunc func(name string, args ...string) (string, error)
}

func (m *MockCommandExecutor) RunCommand(name string, args ...string) (string, error) {
	return m.RunCommandFunc(name, args...)
}

func TestRun(t *testing.T) {
	// Test case 1: Successful command execution
	mockExecutor := &MockCommandExecutor{
		RunCommandFunc: func(name string, args ...string) (string, error) {
			assert.Equal(t, "git", name)
			assert.Equal(t, []string{"test"}, args)
			return "test output", nil
		},
	}
	out, err := Run(mockExecutor, "test")
	assert.NoError(t, err)
	assert.Equal(t, "test output", out)

	// Test case 2: Command execution with error
	mockExecutor = &MockCommandExecutor{
		RunCommandFunc: func(name string, args ...string) (string, error) {
			assert.Equal(t, "git", name)
			assert.Equal(t, []string{"test"}, args)
			return "", errors.New("command failed")
		},
	}
	out, err = Run(mockExecutor, "test")
	assert.Error(t, err)
	assert.Empty(t, out)
}

func TestGetGitRoot(t *testing.T) {
	// Test case 1: Successful git root retrieval
	mockExecutor := &MockCommandExecutor{
		RunCommandFunc: func(name string, args ...string) (string, error) {
			assert.Equal(t, "git", name)
			assert.Equal(t, []string{"rev-parse", "--show-toplevel"}, args)
			return "/mock/git/root", nil
		},
	}
	out, err := GetGitRoot(mockExecutor)
	assert.NoError(t, err)
	assert.Equal(t, "/mock/git/root", out)

	// Test case 2: Git root retrieval with error
	mockExecutor = &MockCommandExecutor{
		RunCommandFunc: func(name string, args ...string) (string, error) {
			assert.Equal(t, "git", name)
			assert.Equal(t, []string{"rev-parse", "--show-toplevel"}, args)
			return "", errors.New("not a git repository")
		},
	}
	out, err = GetGitRoot(mockExecutor)
	assert.Error(t, err)
	assert.Empty(t, out)
}