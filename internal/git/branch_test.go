package git

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPrune(t *testing.T) {
	mockExecutor := &MockCommandExecutor{
		RunCommandFunc: func(name string, args ...string) (string, error) {
			assert.Equal(t, "git", name)
			assert.Equal(t, []string{"fetch", "--prune"}, args)
			return "", nil
		},
	}
	err := Prune(mockExecutor)
	assert.NoError(t, err)
}

func TestGetLocalBranches(t *testing.T) {
	mockExecutor := &MockCommandExecutor{
		RunCommandFunc: func(name string, args ...string) (string, error) {
			assert.Equal(t, "git", name)
			assert.Equal(t, []string{"branch", "--format=%(refname:short)"}, args)
			return "main\ndevelop\nfeature/test", nil
		},
	}
	branches, err := GetLocalBranches(mockExecutor)
	assert.NoError(t, err)
	assert.Equal(t, []string{"main", "develop", "feature/test"}, branches)
}

func TestGetGoneBranches(t *testing.T) {
	// Test case 1: Branches with : gone]
	mockExecutor := &MockCommandExecutor{
		RunCommandFunc: func(name string, args ...string) (string, error) {
			assert.Equal(t, "git", name)
			assert.Equal(t, []string{"branch", "-vv"}, args)
			return `  feature/gone  7890abc [origin/feature/gone: gone] Another commit
  bugfix/gone   56789de [origin/bugfix/gone: gone] Bug fix commit`, nil
		},
	}
	branches, err := GetGoneBranches(mockExecutor)
	assert.NoError(t, err)
	assert.Equal(t, []string{"feature/gone", "bugfix/gone"}, branches)

	// Test case 2: No gone branches
	mockExecutor = &MockCommandExecutor{
		RunCommandFunc: func(name string, args ...string) (string, error) {
			assert.Equal(t, "git", name)
			assert.Equal(t, []string{"branch", "-vv"}, args)
			return `  main        0123456 [origin/main] Commit message\n  develop     def1234 [origin/develop] Yet another commit`, nil
		},
	}
	branches, err = GetGoneBranches(mockExecutor)
	assert.NoError(t, err)
	assert.Empty(t, branches)

	// Test case 3: Error from git command
	mockExecutor = &MockCommandExecutor{
		RunCommandFunc: func(name string, args ...string) (string, error) {
			return "", errors.New("git command failed")
		},
	}
	branches, err = GetGoneBranches(mockExecutor)
	assert.Error(t, err)
	assert.Nil(t, branches)
}