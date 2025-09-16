package cmd

import (
	"errors"
	"io"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/tkr53/ghar/internal/config"
)

type MockCommandExecutor struct {
	RunCommandFunc func(name string, args ...string) (string, error)
}

func (m *MockCommandExecutor) RunCommand(name string, args ...string) (string, error) {
	return m.RunCommandFunc(name, args...)
}

type MockConfigLoader struct {
	LoadConfigFunc func() (*config.Config, error)
}

func (m *MockConfigLoader) LoadConfig() (*config.Config, error) {
	return m.LoadConfigFunc()
}

func TestIsProtected(t *testing.T) {
	// Setup a mock config for testing
	cfg = &config.Config{
		ProtectedBranches: []string{"main", "master", "develop", "release/v1.0"},
	}

	assert.True(t, isProtected("main"))
	assert.True(t, isProtected("master"))
	assert.True(t, isProtected("develop"))
	assert.True(t, isProtected("release/v1.0"))
	assert.False(t, isProtected("feature/new-feature"))
	assert.False(t, isProtected("bugfix/fix-bug"))
}

func TestRun(t *testing.T) {
	mockExecutor := &MockCommandExecutor{
		RunCommandFunc: func(name string, args ...string) (string, error) {
			if name == "git" {
				if args[0] == "rev-parse" {
					return "/mock/git/root", nil
				} else if args[0] == "fetch" {
					return "", nil
				} else if args[0] == "branch" && args[1] == "-vv" {
					return `  feature/gone  7890abc [origin/feature/gone: gone] Another commit
  bugfix/gone   56789de [origin/bugfix/gone: gone] Bug fix commit
  main        0123456 [origin/main] Commit message`, nil
				} else if args[0] == "branch" && (args[1] == "-d" || args[1] == "-D") {
					return "", nil
				}
			}
			return "", errors.New("unexpected git command")
		},
	}

	mockConfigLoader := &MockConfigLoader{
		LoadConfigFunc: func() (*config.Config, error) {
			return &config.Config{
				ProtectedBranches: []string{"main", "master"},
			}, nil
		},
	}

	// Capture stdout
	oldStdout := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	// Simulate user input 'y'
	oldStdin := os.Stdin
	pr, pw, _ := os.Pipe()
	pw.Write([]byte("y\n"))
	pw.Close()
	os.Stdin = pr

	run(nil, nil, mockExecutor, mockConfigLoader)

	// Restore stdout and stdin
	w.Close()
	os.Stdout = oldStdout
	os.Stdin = oldStdin
	pr.Close()

	out, _ := io.ReadAll(r)
	output := string(out)

	assert.Contains(t, output, "The following branches are gone from the remote and can be removed:")
	assert.Contains(t, output, "- feature/gone")
	assert.Contains(t, output, "- bugfix/gone")
	assert.NotContains(t, output, "- main") // main should be protected
	assert.Contains(t, output, "Deleted branch feature/gone")
	assert.Contains(t, output, "Deleted branch bugfix/gone")
}

func TestRunNoBranchesToRemove(t *testing.T) {
	mockExecutor := &MockCommandExecutor{
		RunCommandFunc: func(name string, args ...string) (string, error) {
			if name == "git" {
				if args[0] == "rev-parse" {
					return "/mock/git/root", nil
				} else if args[0] == "fetch" {
					return "", nil
				} else if args[0] == "branch" && args[1] == "-vv" {
					return "", nil // No gone branches
				}
			}
			return "", errors.New("unexpected git command")
		},
	}

	mockConfigLoader := &MockConfigLoader{
		LoadConfigFunc: func() (*config.Config, error) {
			return &config.Config{}, nil
		},
	}

	oldStdout := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	run(nil, nil, mockExecutor, mockConfigLoader)

	w.Close()
	os.Stdout = oldStdout

	out, _ := io.ReadAll(r)
	output := string(out)

	assert.Contains(t, output, "No branches to remove.")
}

func TestRunAborted(t *testing.T) {
	mockExecutor := &MockCommandExecutor{
		RunCommandFunc: func(name string, args ...string) (string, error) {
			if name == "git" {
				if args[0] == "rev-parse" {
					return "/mock/git/root", nil
				} else if args[0] == "fetch" {
					return "", nil
				} else if args[0] == "branch" && args[1] == "-vv" {
					return "  feature/gone  7890abc [origin/feature/gone: gone] Another commit", nil
				}
			}
			return "", errors.New("unexpected git command")
		},
	}

	mockConfigLoader := &MockConfigLoader{
		LoadConfigFunc: func() (*config.Config, error) {
			return &config.Config{}, nil
		},
	}

	oldStdout := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	oldStdin := os.Stdin
	pr, pw, _ := os.Pipe()
	pw.Write([]byte("n\n"))
	pw.Close()
	os.Stdin = pr // Simulate user input 'n'

	run(nil, nil, mockExecutor, mockConfigLoader)

	w.Close()
	os.Stdout = oldStdout
	os.Stdin = oldStdin
	pr.Close()

	out, _ := io.ReadAll(r)
	output := string(out)

	assert.Contains(t, output, "Aborted.")
}

func TestRunForce(t *testing.T) {
	mockExecutor := &MockCommandExecutor{
		RunCommandFunc: func(name string, args ...string) (string, error) {
			if name == "git" {
				if args[0] == "rev-parse" {
					return "/mock/git/root", nil
				} else if args[0] == "fetch" {
					return "", nil
				} else if args[0] == "branch" && args[1] == "-vv" {
					return "  feature/gone  7890abc [origin/feature/gone: gone] Another commit", nil
				} else if args[0] == "branch" && (args[1] == "-d" || args[1] == "-D") {
					return "", nil
				}
			}
			return "", errors.New("unexpected git command")
		},
	}

	mockConfigLoader := &MockConfigLoader{
		LoadConfigFunc: func() (*config.Config, error) {
			return &config.Config{}, nil
		},
	}

	// Set force flag
	force = true

	oldStdout := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	run(nil, nil, mockExecutor, mockConfigLoader)

	w.Close()
	os.Stdout = oldStdout

	out, _ := io.ReadAll(r)
	output := string(out)

	assert.Contains(t, output, "Deleted branch feature/gone")

	force = false // Reset force flag
}

func TestRunMerged(t *testing.T) {
	mockExecutor := &MockCommandExecutor{
		RunCommandFunc: func(name string, args ...string) (string, error) {
			if name == "git" {
				if args[0] == "rev-parse" {
					return "/mock/git/root", nil
				} else if args[0] == "fetch" {
					return "", nil
				} else if args[0] == "branch" && args[1] == "-vv" {
					return "  feature/merged  7890abc [origin/feature/merged: gone] Another commit", nil
				} else if args[0] == "branch" && (args[1] == "-d" || args[1] == "-D") {
					return "", nil
				}
			}
			return "", errors.New("unexpected git command")
		},
	}

	mockConfigLoader := &MockConfigLoader{
		LoadConfigFunc: func() (*config.Config, error) {
			return &config.Config{}, nil
		},
	}

	// Set merged flag
	merged = true

	oldStdout := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	oldStdin := os.Stdin
	pr, pw, _ := os.Pipe()
	pw.Write([]byte("y\n"))
	pw.Close()
	os.Stdin = pr

	run(nil, nil, mockExecutor, mockConfigLoader)

	w.Close()
	os.Stdout = oldStdout
	os.Stdin = oldStdin
	pr.Close()

	out, _ := io.ReadAll(r)
	output := string(out)

	assert.Contains(t, output, "Deleted branch feature/merged")

	merged = false // Reset merged flag
}

func TestRunGitRootError(t *testing.T) {
	// Skip this test because log.Fatalf exits the process
	// TODO: Refactor run() to return error instead of calling log.Fatalf
	t.Skip("Skipping TestRunGitRootError - log.Fatalf exits the process")
}
