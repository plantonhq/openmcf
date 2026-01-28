package webapp

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/plantonhq/openmcf/cmd/openmcf/root"
	"github.com/spf13/cobra"
)

var (
	uninstallPurgeData bool
	uninstallForce     bool
)

var UninstallCmd = &cobra.Command{
	Use:   "uninstall",
	Short: "uninstall the OpenMCF web app",
	Long: `Uninstall the OpenMCF web app by removing the container.
Data volumes are preserved by default unless --purge-data is specified.`,
	Run: uninstallHandler,
}

func init() {
	UninstallCmd.Flags().BoolVar(&uninstallPurgeData, "purge-data", false, "also remove data volumes (WARNING: this deletes all data)")
	UninstallCmd.Flags().BoolVarP(&uninstallForce, "force", "f", false, "skip confirmation prompts")
}

func uninstallHandler(cmd *cobra.Command, args []string) {
	fmt.Println("========================================")
	fmt.Println("🗑️  Uninstalling OpenMCF Web App")
	fmt.Println("========================================")
	fmt.Println()

	// Check if container exists
	if !containerExists() {
		fmt.Printf("ℹ️  Container '%s' not found.\n", ContainerName)
		fmt.Println("   Nothing to uninstall.")
		return
	}

	// Confirmation prompt
	if !uninstallForce {
		fmt.Println("This will:")
		fmt.Println("  - Stop the web app container")
		fmt.Println("  - Remove the container")
		if uninstallPurgeData {
			fmt.Println("  - ⚠️  DELETE ALL DATA (MongoDB, Pulumi state, Go cache)")
		} else {
			fmt.Println("  - Keep data volumes (MongoDB, Pulumi state, Go cache)")
		}
		fmt.Println()
		fmt.Print("Are you sure you want to continue? (yes/no): ")

		reader := bufio.NewReader(os.Stdin)
		response, err := reader.ReadString('\n')
		if err != nil {
			fmt.Printf("❌ Failed to read input: %v\n", err)
			os.Exit(1)
		}

		response = strings.TrimSpace(strings.ToLower(response))
		if response != "yes" && response != "y" {
			fmt.Println("❌ Uninstall cancelled")
			os.Exit(0)
		}
		fmt.Println()
	}

	// Stop container if running
	if isContainerRunning() {
		fmt.Println("🔄 Stopping container...")
		stopCmd := exec.Command("docker", "stop", ContainerName)
		if err := stopCmd.Run(); err != nil {
			fmt.Printf("⚠️  Warning: Failed to stop container: %v\n", err)
		} else {
			fmt.Println("✅ Container stopped")
		}
	}

	// Remove container
	fmt.Println("🔄 Removing container...")
	rmCmd := exec.Command("docker", "rm", ContainerName)
	if err := rmCmd.Run(); err != nil {
		fmt.Printf("❌ Failed to remove container: %v\n", err)
		os.Exit(1)
	}
	fmt.Println("✅ Container removed")

	// Remove volumes if requested
	if uninstallPurgeData {
		fmt.Println("🔄 Removing data volumes...")
		volumes := []string{MongoDBVolume, PulumiVolume, GoCacheVolume}
		for _, volume := range volumes {
			rmVolCmd := exec.Command("docker", "volume", "rm", volume)
			if err := rmVolCmd.Run(); err != nil {
				fmt.Printf("⚠️  Warning: Failed to remove volume %s: %v\n", volume, err)
			} else {
				fmt.Printf("   ✓ Removed %s\n", volume)
			}
		}
		fmt.Println("✅ Data volumes removed")
	} else {
		fmt.Println("ℹ️  Data volumes preserved:")
		fmt.Printf("   - %s\n", MongoDBVolume)
		fmt.Printf("   - %s\n", PulumiVolume)
		fmt.Printf("   - %s\n", GoCacheVolume)
		fmt.Println()
		fmt.Println("   To remove them manually, run:")
		fmt.Println("     docker volume rm " + MongoDBVolume)
		fmt.Println("     docker volume rm " + PulumiVolume)
		fmt.Println("     docker volume rm " + GoCacheVolume)
	}

	// Clean up CLI configuration
	fmt.Println("🔄 Cleaning up CLI configuration...")
	if err := cleanupConfig(); err != nil {
		fmt.Printf("⚠️  Warning: Failed to clean up configuration: %v\n", err)
	} else {
		fmt.Println("✅ CLI configuration cleaned up")
	}

	fmt.Println()
	fmt.Println("========================================")
	fmt.Println("✨ Uninstall Complete!")
	fmt.Println("========================================")
	fmt.Println()

	if !uninstallPurgeData {
		fmt.Println("To reinstall with existing data:")
		fmt.Println("  planton webapp init")
		fmt.Println()
	}
}

func cleanupConfig() error {
	config, err := root.LoadConfigPublic()
	if err != nil {
		return err
	}

	// Clear web app related config
	config.WebAppContainerID = ""
	config.WebAppVersion = ""
	// Note: We keep BackendURL in case user wants to use a different backend

	return root.SaveConfigPublic(config)
}
