package webapp

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/spf13/cobra"
)

var StopCmd = &cobra.Command{
	Use:   "stop",
	Short: "stop the OpenMCF web app",
	Long:  `Stop the OpenMCF web app container (data is preserved)`,
	Run:   stopHandler,
}

func stopHandler(cmd *cobra.Command, args []string) {
	fmt.Println("========================================")
	fmt.Println("🛑 Stopping OpenMCF Web App")
	fmt.Println("========================================")
	fmt.Println()

	// Check if container exists
	if !containerExists() {
		fmt.Printf("❌ Container '%s' not found.\n", ContainerName)
		fmt.Println("   Nothing to stop.")
		os.Exit(1)
	}

	// Check if running
	if !isContainerRunning() {
		fmt.Println("ℹ️  Web app is already stopped")
		return
	}

	// Stop the container
	fmt.Println("🔄 Stopping container...")
	stopCmd := exec.Command("docker", "stop", ContainerName)
	if err := stopCmd.Run(); err != nil {
		fmt.Printf("❌ Failed to stop container: %v\n", err)
		os.Exit(1)
	}

	fmt.Println()
	fmt.Println("✅ Web app stopped successfully")
	fmt.Println()
	fmt.Println("Data is preserved. To start again, run:")
	fmt.Println("  planton webapp start")
	fmt.Println()
}
