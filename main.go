package main

import (
	"fmt"
	"log"

	"github.com/ryosuke/git-branch-auto-remove/internal/git"
)

func main() {
	if err := git.Prune(); err != nil {
		log.Fatalf("Error pruning remote branches: %v", err)
	}

	goneBranches, err := git.GetGoneBranches()
	if err != nil {
		log.Fatalf("Error getting gone branches: %v", err)
	}

	if len(goneBranches) == 0 {
		fmt.Println("No branches to remove.")
		return
	}

	fmt.Println("The following branches are gone from the remote and can be removed:")
	for _, branch := range goneBranches {
		fmt.Printf("- %s\n", branch)
	}
}
