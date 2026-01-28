package root

import (
	"fmt"
	"os"

	"github.com/plantonhq/openmcf/internal/cli/cliprint"
	"github.com/plantonhq/openmcf/internal/cli/iacflags"
	climanifest "github.com/plantonhq/openmcf/internal/cli/manifest"
	"github.com/plantonhq/openmcf/internal/manifest"
	"github.com/spf13/cobra"
)

var ValidateManifest = &cobra.Command{
	Use:   "validate-manifest [manifest-path]",
	Short: "validate a openmcf manifest",
	Aliases: []string{
		"validate",
	},
	Example: `
	# Validate from clipboard
	openmcf validate --clipboard
	openmcf validate -c
	openmcf validate --clip
	openmcf validate --cb

	# Validate from file (positional argument)
	openmcf validate manifest.yaml

	# Validate from file (flag)
	openmcf validate -f manifest.yaml

	# Validate from kustomize
	openmcf validate --kustomize-dir _kustomize --overlay prod
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
			// Check for clipboard-specific errors and display beautifully
			if climanifest.HandleClipboardError(err) {
				os.Exit(1)
			}
			cliprint.PrintError(fmt.Sprintf("failed to resolve manifest: %v", err))
			os.Exit(1)
		}
		if isTemp {
			defer os.Remove(manifestPath)
		}
	}

	err = manifest.Validate(manifestPath)
	if err != nil {
		// Check for manifest load errors (proto unmarshaling) and display beautifully
		if manifest.HandleManifestLoadError(err) {
			os.Exit(1)
		}
		cliprint.PrintError(err.Error())
		os.Exit(1)
	}
	cliprint.PrintSuccessMessage("manifest is valid")
}
