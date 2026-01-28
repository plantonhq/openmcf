package root

import (
	"fmt"
	"os"

	"github.com/plantonhq/openmcf/internal/cli/cliprint"
	"github.com/plantonhq/openmcf/internal/cli/staging"
	"github.com/spf13/cobra"
)

var Pull = &cobra.Command{
	Use:   "pull",
	Short: "Pull latest changes from upstream into staging area",
	Long: `Pull the latest changes from the upstream OpenMCF repository into the local staging area.

The staging area (~/.openmcf/staging/openmcf) is a local cache of the 
OpenMCF repository that enables fast infrastructure operations without requiring 
a network clone on every command execution.

This command performs a 'git fetch --all' followed by 'git pull' in the staging directory,
updating your local cache with the latest IaC modules, API definitions, and tooling from
the upstream repository.

If the staging area does not exist, it will be automatically cloned first.

Use this command periodically to ensure you have access to the latest deployment components
and bug fixes, especially before running apply/preview/destroy operations.`,
	Example: `  # Pull latest changes from upstream
  openmcf pull

  # Typical workflow: pull latest, then apply
  openmcf pull
  openmcf apply -f manifest.yaml`,
	Run: pullHandler,
}

func pullHandler(cmd *cobra.Command, args []string) {
	if err := staging.Pull(); err != nil {
		cliprint.PrintError(fmt.Sprintf("Failed to pull: %v", err))
		os.Exit(1)
	}

	// Get current staging info
	exists, version, repoPath, err := staging.GetStagingInfo()
	if err != nil {
		cliprint.PrintError(fmt.Sprintf("Failed to get staging info: %v", err))
		os.Exit(1)
	}

	if exists {
		fmt.Printf("\nStaging area: %s\n", repoPath)
		if version != "" {
			fmt.Printf("Current version: %s\n", version)
		}
	}
}
