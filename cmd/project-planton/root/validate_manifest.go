package root

import (
	"fmt"
	"os"

	"github.com/plantonhq/project-planton/internal/cli/cliprint"
	"github.com/plantonhq/project-planton/internal/cli/iacflags"
	climanifest "github.com/plantonhq/project-planton/internal/cli/manifest"
	"github.com/plantonhq/project-planton/internal/manifest"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var ValidateManifest = &cobra.Command{
	Use:   "validate-manifest [manifest-path]",
	Short: "validate a project-planton manifest",
	Aliases: []string{
		"validate",
	},
	Example: `
	# Validate from clipboard
	project-planton validate --clipboard
	project-planton validate -c
	project-planton validate --clip
	project-planton validate --cb

	# Validate from file (positional argument)
	project-planton validate manifest.yaml

	# Validate from file (flag)
	project-planton validate -f manifest.yaml

	# Validate from kustomize
	project-planton validate --kustomize-dir _kustomize --overlay prod
	`,
	Args: cobra.MaximumNArgs(1), // Optional manifest path
	Run:  validateHandler,
}

func init() {
	iacflags.AddManifestSourceFlags(ValidateManifest)
}

func validateHandler(cmd *cobra.Command, args []string) {
	var manifestPath string
	var isTemp bool
	var err error

	// If a positional arg is provided, use it as manifest path (backward compatibility)
	if len(args) > 0 {
		manifestPath = args[0]
	} else {
		// Use unified resolver for --clipboard, --manifest, --kustomize-dir, etc.
		manifestPath, isTemp, err = climanifest.ResolveManifestPath(cmd)
		if err != nil {
			log.Fatalf("failed to resolve manifest: %v", err)
		}
		if isTemp {
			defer os.Remove(manifestPath)
		}
	}

	err = manifest.Validate(manifestPath)
	if err != nil {
		fmt.Printf("%v\n", err)
		return
	}
	cliprint.PrintSuccessMessage("manifest is valid")
}
