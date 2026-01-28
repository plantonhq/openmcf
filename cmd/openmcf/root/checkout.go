package root

import (
	"fmt"
	"os"

	"github.com/plantonhq/openmcf/internal/cli/cliprint"
	"github.com/plantonhq/openmcf/internal/cli/staging"
	"github.com/spf13/cobra"
)

var Checkout = &cobra.Command{
	Use:   "checkout <version>",
	Short: "Checkout a specific version tag in the staging area",
	Long: `Checkout a specific version tag or commit in the local staging area.

The staging area (~/.openmcf/staging/openmcf) maintains a local copy
of the OpenMCF repository. This command allows you to switch between different
versions of the IaC modules without affecting your CLI binary version.

This is useful when:
- You need to use IaC modules from a specific release version
- You want to test with a newer version before upgrading the CLI
- You need to rollback to an older module version for compatibility
- You're debugging and need to match the exact module version used previously

The version argument can be:
- "latest" - automatically resolves to the most recent release tag
- A release tag (e.g., v0.2.273)
- A branch name (e.g., main)
- A commit SHA

If the staging area does not exist, it will be automatically cloned first.`,
	Example: `  # Checkout the latest release tag
  openmcf checkout latest

  # Checkout a specific release version
  openmcf checkout v0.2.273

  # Checkout the main branch for latest development
  openmcf checkout main

  # Checkout a specific commit
  openmcf checkout abc1234`,
	Args: cobra.ExactArgs(1),
	Run:  checkoutHandler,
}

func checkoutHandler(cmd *cobra.Command, args []string) {
	version := args[0]

	if err := staging.Checkout(version); err != nil {
		cliprint.PrintError(fmt.Sprintf("Failed to checkout: %v", err))
		os.Exit(1)
	}

	// Get current staging info
	_, _, repoPath, err := staging.GetStagingInfo()
	if err != nil {
		cliprint.PrintError(fmt.Sprintf("Failed to get staging info: %v", err))
		os.Exit(1)
	}

	fmt.Printf("\nStaging area: %s\n", repoPath)
}
