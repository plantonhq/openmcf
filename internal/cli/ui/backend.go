package ui

import (
	"fmt"
	"strings"

	"github.com/plantonhq/project-planton/pkg/iac/tofu/backendconfig"
)

// S3CompatibleDetected displays a helpful message when S3-compatible backend is detected.
// This informs the user that additional configuration flags will be applied automatically.
func S3CompatibleDetected(reason string) {
	fmt.Println()
	fmt.Printf("%s  %s\n",
		infoIcon.Render("ℹ"),
		infoTitle.Render("S3-Compatible Backend Detected"))
	fmt.Printf("   %s\n", dimStyle.Render(reason))
	fmt.Printf("   %s\n", dimStyle.Render("Additional configuration will be applied automatically:"))
	fmt.Printf("   %s\n", dimStyle.Render("  • skip_credentials_validation = true"))
	fmt.Printf("   %s\n", dimStyle.Render("  • skip_region_validation = true"))
	fmt.Printf("   %s\n", dimStyle.Render("  • skip_requesting_account_id = true"))
	fmt.Printf("   %s\n", dimStyle.Render("  • skip_metadata_api_check = true"))
	fmt.Printf("   %s\n", dimStyle.Render("  • use_path_style = true"))
	fmt.Println()
}

// BackendConfigSummary displays the resolved backend configuration in a clean format.
func BackendConfigSummary(config *backendconfig.TofuBackendConfig) {
	fmt.Println()
	fmt.Printf("%s  %s\n",
		successIcon.Render("📦"),
		successTitle.Render("Backend Configuration"))
	fmt.Printf("   %-12s %s\n", dimStyle.Render("Type:"), infoMessage.Render(config.BackendType))
	fmt.Printf("   %-12s %s\n", dimStyle.Render("Bucket:"), infoMessage.Render(config.BackendBucket))
	fmt.Printf("   %-12s %s\n", dimStyle.Render("Key:"), infoMessage.Render(config.BackendKey))
	if config.BackendRegion != "" {
		fmt.Printf("   %-12s %s\n", dimStyle.Render("Region:"), infoMessage.Render(config.BackendRegion))
	}
	if config.BackendEndpoint != "" {
		fmt.Printf("   %-12s %s\n", dimStyle.Render("Endpoint:"), infoMessage.Render(config.BackendEndpoint))
	}
	if config.S3Compatible {
		fmt.Printf("   %-12s %s\n", dimStyle.Render("Mode:"), warningTitle.Render("S3-Compatible (R2/MinIO)"))
	}
	fmt.Println()
}

// MissingBackendConfigError displays a helpful error for missing backend configuration.
// This provides actionable guidance with CLI flags and examples for each missing field.
func MissingBackendConfigError(missing []backendconfig.MissingField, backendType string) {
	fmt.Println()
	fmt.Printf("%s  %s\n",
		errorIcon.Render(iconError),
		errorTitle.Render("Incomplete Backend Configuration"))
	fmt.Println()
	fmt.Printf("   %s backend requires the following configuration:\n\n",
		infoMessage.Render(strings.ToUpper(backendType)))

	for _, field := range missing {
		if !field.Required {
			continue
		}
		fmt.Printf("   %s %s\n",
			errorIcon.Render("•"),
			errorMessage.Render(field.Description))
		fmt.Printf("     %s %s\n",
			dimStyle.Render("Flag:"),
			cmdStyle.Render(field.FlagName))
		fmt.Printf("     %s %s\n",
			dimStyle.Render("Example:"),
			cmdStyle.Render(field.Example))
		fmt.Println()
	}
}

// MissingFieldPrompt displays the prompt for a single missing field.
func MissingFieldPrompt(field backendconfig.MissingField) {
	fmt.Println()
	fmt.Printf("%s %s\n",
		warningIcon.Render("!"),
		warningTitle.Render(fmt.Sprintf("Missing required field: %s", field.Name)))
	fmt.Printf("   %s\n", dimStyle.Render(field.Description))
	fmt.Printf("   %s %s\n", dimStyle.Render("Example:"), cmdStyle.Render(field.Example))
	fmt.Printf("   %s %s\n", dimStyle.Render("CLI flag:"), cmdStyle.Render(field.FlagName))
	fmt.Println()
}

// WarningIcon returns the styled warning icon for external use.
func WarningIcon() string {
	return warningIcon.Render(iconWarning)
}

// Bold returns text styled as bold for external use.
func Bold(text string) string {
	return warningTitle.Render(text)
}
