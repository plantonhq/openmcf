package upgrade

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/plantonhq/openmcf/internal/cli/cliprint"
)

// UpgradeViaHomebrew upgrades the CLI using Homebrew cask
func UpgradeViaHomebrew() error {
	// Step 1: Update Homebrew
	fmt.Println()
	cliprint.PrintStep("Updating Homebrew...")

	updateCmd := exec.Command("brew", "update")
	updateCmd.Stdout = os.Stdout
	updateCmd.Stderr = os.Stderr
	if err := updateCmd.Run(); err != nil {
		return fmt.Errorf("failed to update Homebrew: %w", err)
	}

	// Step 2: Upgrade openmcf cask
	fmt.Println()
	cliprint.PrintStep("Upgrading openmcf...")

	upgradeCmd := exec.Command("brew", "upgrade", "--cask", "openmcf")
	upgradeCmd.Stdout = os.Stdout
	upgradeCmd.Stderr = os.Stderr
	if err := upgradeCmd.Run(); err != nil {
		// brew upgrade returns non-zero if already up to date, which is fine
		// Check if openmcf is actually installed and at latest
		return fmt.Errorf("Homebrew upgrade failed: %w", err)
	}

	return nil
}

// UninstallHomebrew uninstalls the CLI via Homebrew cask
// This is used when transitioning from Homebrew to direct-download management
// for version-specific installs
func UninstallHomebrew() error {
	fmt.Println()
	cliprint.PrintStep("Uninstalling via Homebrew...")

	cmd := exec.Command("brew", "uninstall", "--cask", "openmcf")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to uninstall via Homebrew: %w", err)
	}

	cliprint.PrintSuccess("Uninstalled via Homebrew")
	return nil
}
