package root

import (
	"fmt"
	"os"

	"github.com/plantonhq/project-planton/internal/cli/cliprint"
	"github.com/plantonhq/project-planton/internal/cli/flag"
	"github.com/plantonhq/project-planton/internal/cli/iacflags"
	climanifest "github.com/plantonhq/project-planton/internal/cli/manifest"
	"github.com/plantonhq/project-planton/internal/manifest"
	"github.com/spf13/cobra"
)

var LoadManifest = &cobra.Command{
	Use:     "load-manifest [manifest-path]",
	Aliases: []string{"load"},
	Short:   "load a project-planton manifest from provided path or kustomize",
	Example: `
	# Load from clipboard
	project-planton load --clipboard
	project-planton load -c
	project-planton load --clip
	project-planton load --cb

	# Load from file (positional argument)
	project-planton load manifest.yaml

	# Load from file (flag)
	project-planton load -f manifest.yaml

	# Load from kustomize
	project-planton load --kustomize-dir _kustomize --overlay prod

	# Load with overrides
	project-planton load -f manifest.yaml --set spec.version=v1.2.3
	project-planton load --clipboard --set spec.replicas=3
	`,
	Args: cobra.MaximumNArgs(1), // Optional manifest path
	Run:  loadManifestHandler,
}

func init() {
	iacflags.AddManifestSourceFlags(LoadManifest)
	LoadManifest.PersistentFlags().StringToString(string(flag.Set), map[string]string{}, "override resource manifest values using key=value pairs")
}

func loadManifestHandler(cmd *cobra.Command, args []string) {
	valueOverrides, err := cmd.Flags().GetStringToString(string(flag.Set))
	flag.HandleFlagErr(err, flag.Set)

	var manifestPath string
	var isTemp bool

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

	updatedManifest, err := manifest.LoadWithOverrides(manifestPath, valueOverrides)
	if err != nil {
		// Check for manifest load errors (proto unmarshaling) and display beautifully
		if manifest.HandleManifestLoadError(err) {
			os.Exit(1)
		}
		cliprint.PrintError(fmt.Sprintf("failed to load manifest: %v", err))
		os.Exit(1)
	}
	if err := manifest.Print(updatedManifest); err != nil {
		cliprint.PrintError(fmt.Sprintf("failed to print manifest: %v", err))
		os.Exit(1)
	}
}
