package iacrunner

import (
	"fmt"
	"os"

	"github.com/plantonhq/project-planton/apis/org/project_planton/shared/iac/terraform"
	"github.com/plantonhq/project-planton/internal/cli/cliprint"
	"github.com/plantonhq/project-planton/internal/cli/flag"
	"github.com/plantonhq/project-planton/internal/cli/prompt"
	"github.com/plantonhq/project-planton/internal/cli/ui"
	"github.com/plantonhq/project-planton/pkg/iac/provisioner"
	"github.com/plantonhq/project-planton/pkg/iac/tofu/backendconfig"
	"github.com/plantonhq/project-planton/pkg/iac/tofu/tofumodule"
	"github.com/spf13/cobra"
)

// RunTofu executes an OpenTofu operation using the resolved context.
func RunTofu(ctx *Context, cmd *cobra.Command, operation terraform.TerraformOperationType) error {
	return runHcl(ctx, cmd, operation, provisioner.HclBinaryTofu)
}

// runHcl executes an HCL-based IaC operation (tofu or terraform) using the resolved context.
func runHcl(ctx *Context, cmd *cobra.Command, operation terraform.TerraformOperationType, binary provisioner.HclBinary) error {
	// Get auto-approve flag if defined (ignore error for commands that don't register it)
	isAutoApprove, _ := cmd.Flags().GetBool(string(flag.AutoApprove))

	// Get reconfigure flag if defined (ignore error for commands that don't register it)
	isReconfigure, _ := cmd.Flags().GetBool(string(flag.Reconfigure))

	// For plan operation, check if it's a destroy plan
	isDestroyPlan := false
	if operation == terraform.TerraformOperationType_plan {
		isDestroyPlan, _ = cmd.Flags().GetBool(string(flag.Destroy))
		// Plan is always auto-approve (non-interactive)
		isAutoApprove = true
	}

	// For refresh operation, no approval needed (read-only state sync)
	if operation == terraform.TerraformOperationType_refresh {
		isAutoApprove = true
	}

	// Check if binary is available before proceeding
	if err := binary.CheckAvailable(); err != nil {
		cliprint.PrintError(err.Error())
		os.Exit(1)
	}

	// Build and validate backend configuration before handoff
	backendCfg, err := buildAndValidateBackendConfig(ctx, cmd, string(binary))
	if err != nil {
		cliprint.PrintError(err.Error())
		os.Exit(1)
	}

	// Display backend configuration if available (before handoff)
	if backendCfg != nil && backendCfg.BackendType != "" && backendCfg.BackendType != "local" {
		ui.BackendConfigSummary(backendCfg)
	}

	// Display module path
	cliprint.PrintModulePath(ctx.ModuleDir)

	cliprint.PrintHandoff(binary.DisplayName())

	err = tofumodule.RunCommand(
		string(binary),
		ctx.ModuleDir,
		ctx.ManifestPath,
		operation,
		ctx.ValueOverrides,
		isAutoApprove,
		isDestroyPlan,
		isReconfigure,
		ctx.ModuleVersion,
		ctx.NoCleanup,
		ctx.KubeContext,
		ctx.ProviderConfig,
	)
	if err != nil {
		printHclFailure(binary)
		os.Exit(1)
	}

	printHclSuccess(binary)
	return nil
}

// buildAndValidateBackendConfig builds backend config from CLI flags and manifest labels,
// validates it, and prompts for missing values if in interactive mode.
func buildAndValidateBackendConfig(ctx *Context, cmd *cobra.Command, provisionerType string) (*backendconfig.TofuBackendConfig, error) {
	// Extract CLI flags for backend configuration
	cliFlags := extractCLIBackendFlags(cmd)

	// Build merged configuration (CLI flags override manifest labels)
	config, err := backendconfig.BuildBackendConfig(ctx.ManifestObject, provisionerType, cliFlags)
	if err != nil {
		return nil, fmt.Errorf("failed to build backend configuration: %w", err)
	}

	// Skip validation for local backend or when no backend is configured
	if config.BackendType == "" || config.BackendType == "local" {
		return config, nil
	}

	// Detect and announce S3-compatible backend
	if config.S3Compatible {
		if config.BackendRegion == "auto" {
			ui.S3CompatibleDetected("Region is set to 'auto', indicating an S3-compatible backend")
		} else if config.BackendEndpoint != "" {
			ui.S3CompatibleDetected("Custom endpoint detected")
		}
	}

	// Validate configuration completeness
	validation := backendconfig.Validate(config)

	// If configuration is incomplete, handle based on interactivity
	if !validation.Valid {
		if !prompt.IsInteractive() {
			// Non-interactive mode: show error and fail
			ui.MissingBackendConfigError(validation.MissingFields, config.BackendType)
			return nil, fmt.Errorf("incomplete backend configuration - provide missing values via CLI flags or manifest labels")
		}

		// Interactive mode: prompt for missing values
		ui.MissingBackendConfigError(validation.MissingFields, config.BackendType)
		config, err = prompt.PromptForMissingBackendConfig(config, validation.MissingFields)
		if err != nil {
			return nil, fmt.Errorf("failed to get backend configuration: %w", err)
		}
	}

	return config, nil
}

// extractCLIBackendFlags extracts backend configuration from CLI flags.
func extractCLIBackendFlags(cmd *cobra.Command) backendconfig.CLIBackendFlags {
	// Get values from flags, ignoring errors for flags that aren't registered
	backendType, _ := cmd.Flags().GetString(string(flag.BackendType))
	backendBucket, _ := cmd.Flags().GetString(string(flag.BackendBucket))
	backendKey, _ := cmd.Flags().GetString(string(flag.BackendKey))
	backendRegion, _ := cmd.Flags().GetString(string(flag.BackendRegion))
	backendEndpoint, _ := cmd.Flags().GetString(string(flag.BackendEndpoint))

	return backendconfig.CLIBackendFlags{
		BackendType:     backendType,
		BackendBucket:   backendBucket,
		BackendKey:      backendKey,
		BackendRegion:   backendRegion,
		BackendEndpoint: backendEndpoint,
	}
}

// printHclSuccess prints a success message for the appropriate binary.
func printHclSuccess(binary provisioner.HclBinary) {
	switch binary {
	case provisioner.HclBinaryTofu:
		cliprint.PrintTofuSuccess()
	case provisioner.HclBinaryTerraform:
		cliprint.PrintTerraformSuccess()
	}
}

// printHclFailure prints a failure message for the appropriate binary.
func printHclFailure(binary provisioner.HclBinary) {
	switch binary {
	case provisioner.HclBinaryTofu:
		cliprint.PrintTofuFailure()
	case provisioner.HclBinaryTerraform:
		cliprint.PrintTerraformFailure()
	}
}
