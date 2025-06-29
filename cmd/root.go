package cmd

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/fatih/color"
	"github.com/spf13/cobra"
	"github.com/tkr53/git-branch-auto-remove/internal/config"
	"github.com/tkr53/git-branch-auto-remove/internal/git"
)

var (
	force  bool
	merged bool
	cfg    *config.Config
)

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
	for _, p := range cfg.ProtectedBranches {
		if branch == p {
			return true
		}
	}
	return false
}

func run(cmd *cobra.Command, args []string) {
	var err error
	cfg, err = config.LoadConfig()
	if err != nil {
		log.Fatalf(color.RedString("Error loading config: %v"), err)
	}

	// Check if it's a git repository
	if _, err := git.GetGitRoot(); err != nil {
		log.Fatalf(color.RedString("Error: %v"), err)
	}

	if err := git.Prune(); err != nil {
		log.Fatalf(color.RedString("Error pruning remote branches: %v"), err)
	}

	goneBranches, err := git.GetGoneBranches()
	if err != nil {
		log.Fatalf(color.RedString("Error getting gone branches: %v"), err)
	}

	// Filter out protected branches
	var branchesToRemove []string
	for _, branch := range goneBranches {
		if !isProtected(branch) {
			branchesToRemove = append(branchesToRemove, branch)
		}
	}

	if len(branchesToRemove) == 0 {
		fmt.Println(color.YellowString("No branches to remove."))
		return
	}

	fmt.Println(color.YellowString("The following branches are gone from the remote and can be removed:"))
	for _, branch := range branchesToRemove {
		fmt.Printf("- %s\n", color.GreenString(branch))
	}

	if !force {
		fmt.Print("\nDo you want to remove these branches? (y/n): ")
		reader := bufio.NewReader(os.Stdin)
		input, _ := reader.ReadString('\n')
		if strings.TrimSpace(input) != "y" {
			fmt.Println(color.YellowString("Aborted."))
			return
		}
	}

	deleteCmd := "-d"
	if merged {
		deleteCmd = "-D"
	}

	for _, branch := range branchesToRemove {
		if _, err := git.Run("branch", deleteCmd, branch); err != nil {
			log.Printf(color.RedString("Failed to delete branch %s: %v"), branch, err)
		} else {
			fmt.Printf(color.GreenString("Deleted branch %s\n"), branch)
		}
	}
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(color.RedString(err.Error()))
		os.Exit(1)
	}
}
