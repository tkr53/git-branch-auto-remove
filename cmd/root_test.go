package cmd

import (
	"errors"
	"io"
	"os"
	"testing"

	"bou.ke/monkey"
	"github.com/ryosuke/git-branch-auto-remove/internal/config"
	"github.com/ryosuke/git-branch-auto-remove/internal/git"
	"github.com/stretchr/testify/assert"
)

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
	// Mock git.Prune
	monkey.Patch(git.Prune, func() error {
		return nil
	})

	// Mock git.GetGitRoot
	monkey.Patch(git.GetGitRoot, func() (string, error) {
		return "/mock/git/root", nil
	})

	// Mock git.GetGoneBranches
	monkey.Patch(git.GetGoneBranches, func() ([]string, error) {
		return []string{"feature/gone", "bugfix/gone", "main"}, nil
	})

	// Mock git.Run for branch deletion
	monkey.Patch(git.Run, func(args ...string) (string, error) {
		assert.Contains(t, []string{"-d", "-D"}, args[0])
		return "", nil
	})

	// Mock config loading
	monkey.Patch(config.LoadConfig, func() (*config.Config, error) {
		return &config.Config{
			ProtectedBranches: []string{"main", "master"},
		}, nil
	})

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

	// Restore stdout and stdin
	w.Close()
	os.Stdout = oldStdout
	os.Stdin = oldStdin
	pr.Close()
	w.Close()
	os.Stdout = oldStdout
	os.Stdin = oldStdin

	out, _ := io.ReadAll(r)
	output := string(out)

	assert.Contains(t, output, "The following branches are gone from the remote and can be removed:")
	assert.Contains(t, output, "- feature/gone")
	assert.Contains(t, output, "- bugfix/gone")
	assert.NotContains(t, output, "- main") // main should be protected
	assert.Contains(t, output, "Deleted branch feature/gone")
	assert.Contains(t, output, "Deleted branch bugfix/gone")

	monkey.UnpatchAll()
}

func TestRunNoBranchesToRemove(t *testing.T) {
	monkey.Patch(git.Prune, func() error {
		return nil
	})
	monkey.Patch(git.GetGitRoot, func() (string, error) {
		return "/mock/git/root", nil
	})
	monkey.Patch(git.GetGoneBranches, func() ([]string, error) {
		return []string{}, nil
	})
	monkey.Patch(config.LoadConfig, func() (*config.Config, error) {
		return &config.Config{}, nil
	})

	oldStdout := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	run(nil, nil)

	w.Close()
	os.Stdout = oldStdout

	out, _ := io.ReadAll(r)
	output := string(out)

	assert.Contains(t, output, "No branches to remove.")

	monkey.UnpatchAll()
}

func TestRunAborted(t *testing.T) {
	monkey.Patch(git.Prune, func() error {
		return nil
	})
	monkey.Patch(git.GetGitRoot, func() (string, error) {
		return "/mock/git/root", nil
	})
	monkey.Patch(git.GetGoneBranches, func() ([]string, error) {
		return []string{"feature/gone"}, nil
	})
	monkey.Patch(config.LoadConfig, func() (*config.Config, error) {
		return &config.Config{}, nil
	})

	oldStdout := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	oldStdin := os.Stdin
	pr, pw, _ := os.Pipe()
	pw.Write([]byte("n\n"))
	pw.Close()
	os.Stdin = pr // Simulate user input 'n'

	run(nil, nil)

	w.Close()
	os.Stdout = oldStdout
	os.Stdin = oldStdin
	pr.Close()

	out, _ := io.ReadAll(r)
	output := string(out)

	assert.Contains(t, output, "Aborted.")

	monkey.UnpatchAll()
}

func TestRunForce(t *testing.T) {
	monkey.Patch(git.Prune, func() error {
		return nil
	})
	monkey.Patch(git.GetGitRoot, func() (string, error) {
		return "/mock/git/root", nil
	})
	monkey.Patch(git.GetGoneBranches, func() ([]string, error) {
		return []string{"feature/gone"}, nil
	})
	monkey.Patch(git.Run, func(args ...string) (string, error) {
		assert.Contains(t, []string{"-d", "-D"}, args[0])
		return "", nil
	})
	monkey.Patch(config.LoadConfig, func() (*config.Config, error) {
		return &config.Config{}, nil
	})

	// Set force flag
	force = true

	oldStdout := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	run(nil, nil)

	w.Close()
	os.Stdout = oldStdout

	out, _ := io.ReadAll(r)
	output := string(out)

	assert.Contains(t, output, "Deleted branch feature/gone")

	force = false // Reset force flag
	monkey.UnpatchAll()
}

func TestRunMerged(t *testing.T) {
	monkey.Patch(git.Prune, func() error {
		return nil
	})
	monkey.Patch(git.GetGitRoot, func() (string, error) {
		return "/mock/git/root", nil
	})
	monkey.Patch(git.GetGoneBranches, func() ([]string, error) {
		return []string{"feature/merged"}, nil
	})
	monkey.Patch(git.Run, func(args ...string) (string, error) {
		assert.Equal(t, "-D", args[0]) // Should use -D for merged branches
		return "", nil
	})
	monkey.Patch(config.LoadConfig, func() (*config.Config, error) {
		return &config.Config{}, nil
	})

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

	run(nil, nil)

	w.Close()
	os.Stdout = oldStdout
	os.Stdin = oldStdin
	pr.Close()

	out, _ := io.ReadAll(r)
	output := string(out)

	assert.Contains(t, output, "Deleted branch feature/merged")

	merged = false // Reset merged flag
	monkey.UnpatchAll()
}

func TestRunGitRootError(t *testing.T) {
	monkey.Patch(git.GetGitRoot, func() (string, error) {
		return "", errors.New("not a git repository")
	})

	// Capture os.Stderr output
	oldStderr := os.Stderr
	r, w, _ := os.Pipe()
	os.Stderr = w

	// Run in a goroutine to catch os.Exit
	done := make(chan struct{})
	go func() {
		defer func() {
			if recover() != nil {
				// ignore panic from os.Exit in test
			}
			close(done)
		}()
		run(nil, nil)
	}()

	w.Close()
	os.Stderr = oldStderr
	<-done

	out, _ := io.ReadAll(r)
	output := string(out)

	assert.Contains(t, output, "Error: not a git repository")

	monkey.UnpatchAll()
}
