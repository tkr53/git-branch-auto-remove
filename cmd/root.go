package cmd

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/ryosuke/git-branch-auto-remove/internal/git"
	"github.com/spf13/cobra"
)

var (
	force  bool
	merged bool
)

// Add a list of protected branches that should not be deleted.
var protectedBranches = []string{"main", "master", "develop"}

var rootCmd = &cobra.Command{
	Use:   "git-branch-auto-remove",
	Short: "A CLI tool to remove local branches that are gone from the remote",
	Run:   run,
}

func init() {
	rootCmd.Flags().BoolVar(&force, "force", false, "Force execute delete branches.")
	rootCmd.Flags().BoolVarP(&merged, "merged", "D", false, "Delete merged branches.")
}

// isProtected checks if a branch is in the protected list.
func isProtected(branch string) bool {
	for _, p := range protectedBranches {
		if branch == p {
			return true
		}
	}
	return false
}

func run(cmd *cobra.Command, args []string) {
	if err := git.Prune(); err != nil {
		log.Fatalf("Error pruning remote branches: %v", err)
	}

	goneBranches, err := git.GetGoneBranches()
	if err != nil {
		log.Fatalf("Error getting gone branches: %v", err)
	}

	// Filter out protected branches
	var branchesToRemove []string
	for _, branch := range goneBranches {
		if !isProtected(branch) {
			branchesToRemove = append(branchesToRemove, branch)
		}
	}

	if len(branchesToRemove) == 0 {
		fmt.Println("No branches to remove.")
		return
	}

	fmt.Println("The following branches are gone from the remote and can be removed:")
	for _, branch := range branchesToRemove {
		fmt.Printf("- %s\n", branch)
	}

	if !force {
		fmt.Print("\nDo you want to remove these branches? (y/n): ")
		reader := bufio.NewReader(os.Stdin)
		input, _ := reader.ReadString('\n')
		if strings.TrimSpace(input) != "y" {
			fmt.Println("Aborted.")
			return
		}
	}

	deleteCmd := "-d"
	if merged {
		deleteCmd = "-D"
	}

	for _, branch := range branchesToRemove {
		if _, err := git.Run("branch", deleteCmd, branch); err != nil {
			log.Printf("Failed to delete branch %s: %v", branch, err)
		} else {
			fmt.Printf("Deleted branch %s\n", branch)
		}
	}
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}