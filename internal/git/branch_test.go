package git

import (
	"os/exec"
	"testing"

	"bou.ke/monkey"
	"github.com/stretchr/testify/assert"
)

func TestPrune(t *testing.T) {
	monkey.Patch(exec.Command, func(name string, arg ...string) *exec.Cmd {
		assert.Equal(t, "git", name)
		assert.Equal(t, []string{"fetch", "--prune"}, arg)
		return exec.Command("echo", "")
	})

	err := Prune()
	assert.NoError(t, err)

	monkey.UnpatchAll()
}

func TestGetLocalBranches(t *testing.T) {
	monkey.Patch(exec.Command, func(name string, arg ...string) *exec.Cmd {
		assert.Equal(t, "git", name)
		assert.Equal(t, []string{"branch", "--format=%(refname:short)"}, arg)
		return exec.Command("echo", "main\ndevelop\nfeature/test")
	})

	branches, err := GetLocalBranches()
	assert.NoError(t, err)
	assert.Equal(t, []string{"main", "develop", "feature/test"}, branches)

	monkey.UnpatchAll()
}

func TestGetGoneBranches(t *testing.T) {
	// Test case 1: Branches with : gone]
	monkey.Patch(exec.Command, func(name string, arg ...string) *exec.Cmd {
		assert.Equal(t, "git", name)
		assert.Equal(t, []string{"branch", "-vv"}, arg)
		return exec.Command("echo", `  main        0123456 [origin/main] Commit message\n  feature/gone  7890abc [origin/feature/gone: gone] Another commit\n  develop     def1234 [origin/develop] Yet another commit\n  bugfix/gone   56789de [origin/bugfix/gone: gone] Bug fix commit`)
	})

	branches, err := GetGoneBranches()
	assert.NoError(t, err)
	assert.Equal(t, []string{"feature/gone", "bugfix/gone"}, branches)

	monkey.UnpatchAll()

	// Test case 2: No gone branches
	monkey.Patch(exec.Command, func(name string, arg ...string) *exec.Cmd {
		assert.Equal(t, "git", name)
		assert.Equal(t, []string{"branch", "-vv"}, arg)
		return exec.Command("echo", `  main        0123456 [origin/main] Commit message\n  develop     def1234 [origin/develop] Yet another commit`)
	})

	branches, err = GetGoneBranches()
	assert.NoError(t, err)
	assert.Empty(t, branches)

	monkey.UnpatchAll()
}
