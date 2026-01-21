// Package project_planton provides the root command for the Project Planton CLI.
// Auto-release test: CLI change triggers v{semver}.{YYYYMMDD}.{N} tag format.
package project_planton

import (
	"fmt"
	"os"

	"github.com/plantonhq/project-planton/cmd/project-planton/root"
	"github.com/plantonhq/project-planton/cmd/project-planton/root/webapp"
	"github.com/plantonhq/project-planton/internal/cli/flag"
	"github.com/spf13/cobra"
)

// DefaultProjectPlantonGitRepo is the default path for the local project-planton git repository
const DefaultProjectPlantonGitRepo = "~/scm/github.com/plantonhq/project-planton"

var rootCmd = &cobra.Command{
	Use:   "project-planton",
	Short: "Unified Interface for Multi-Cloud Infrastructure",
	Run: func(cmd *cobra.Command, args []string) {
		// Check if version flag was passed
		showVersion, _ := cmd.Flags().GetBool("version")
		if showVersion {
			root.PrintVersion()
			return
		}
		// Otherwise show help
		cmd.Help()
	},
}

func init() {
	rootCmd.CompletionOptions.DisableDefaultCmd = true
	rootCmd.DisableSuggestions = true

	// Enable -v as shorthand for --version (handled in Run function for colorful output)
	rootCmd.Flags().BoolP("version", "v", false, "show version information")

	// Local module flags - inherited by all subcommands
	rootCmd.PersistentFlags().Bool(string(flag.LocalModule), false,
		"Use local project-planton git repository for IaC modules instead of downloading")
	rootCmd.PersistentFlags().String(string(flag.ProjectPlantonGitRepo), DefaultProjectPlantonGitRepo,
		"Path to local project-planton git repository (used with --local-module)")

	rootCmd.AddCommand(
		root.Apply,
		root.Checkout,
		root.CloudResourceApplyCmd,
		root.CloudResourceCreateCmd,
		root.CloudResourceDeleteCmd,
		root.CloudResourceGetCmd,
		root.CloudResourceListCmd,
		root.CloudResourceUpdateCmd,
		root.ConfigCmd,
		root.CredentialCreateCmd,
		root.CredentialDeleteCmd,
		root.CredentialGetCmd,
		root.CredentialListCmd,
		root.CredentialUpdateCmd,
		root.Destroy,
		root.Downgrade,
		root.Init,
		root.LoadManifest,
		root.ModulesVersion,
		root.Plan,
		root.Pull,
		root.Pulumi,
		root.Refresh,
		root.StackUpdateStreamOutputCmd,
		root.Terraform,
		root.Tofu,
		root.Upgrade,
		root.ValidateManifest,
		root.Version,
		webapp.WebAppCmd,
	)
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
