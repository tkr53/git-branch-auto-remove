package git

import "strings"

// Prune fetches and prunes remote branches.
func Prune(executor CommandExecutor) error {
	_, err := Run(executor, "fetch", "--prune")
	return err
}

// GetLocalBranches returns a list of local branches.
func GetLocalBranches(executor CommandExecutor) ([]string, error) {
	output, err := Run(executor, "branch", "--format=%(refname:short)")
	if err != nil {
		return nil, err
	}
	return strings.Split(output, "\n"), nil
}

// GetGoneBranches returns a list of branches that are gone from the remote.
func GetGoneBranches(executor CommandExecutor) ([]string, error) {
	output, err := Run(executor, "branch", "-vv")
	if err != nil {
		return nil, err
	}

	var goneBranches []string
	lines := strings.Split(output, "\n")
	for _, line := range lines {
		if strings.Contains(line, ": gone]") {
			fields := strings.Fields(line)
			if len(fields) > 0 {
				goneBranches = append(goneBranches, strings.TrimSpace(fields[0]))
			}
		}
	}
	return goneBranches, nil
}