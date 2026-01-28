package webapp

import (
	"fmt"
	"os"
	"os/exec"
	"time"

	"github.com/spf13/cobra"
)

var StartCmd = &cobra.Command{
	Use:   "start",
	Short: "start the OpenMCF web app",
	Long:  `Start the OpenMCF web app container and wait for services to be ready`,
	Run:   startHandler,
}

func startHandler(cmd *cobra.Command, args []string) {
	fmt.Println("========================================")
	fmt.Println("🚀 Starting OpenMCF Web App")
	fmt.Println("========================================")
	fmt.Println()

	// Check if container exists
	if !containerExists() {
		fmt.Printf("❌ Container '%s' not found.\n", ContainerName)
		fmt.Println("   Please run: planton webapp init")
		os.Exit(1)
	}

	// Check if already running
	if isContainerRunning() {
		fmt.Println("✅ Web app is already running")
		printAccessInfo()
		return
	}

	// Start the container
	fmt.Println("🔄 Starting container...")
	startCmd := exec.Command("docker", "start", ContainerName)
	if err := startCmd.Run(); err != nil {
		fmt.Printf("❌ Failed to start container: %v\n", err)
		os.Exit(1)
	}

	// Wait for services to be healthy
	fmt.Println("⏳ Waiting for services to start (this may take 30-60 seconds)...")
	if err := waitForHealthy(60); err != nil {
		fmt.Printf("⚠️  Warning: %v\n", err)
		fmt.Println("   Check logs with: planton webapp logs")
	} else {
		fmt.Println("✅ All services are healthy")
	}

	fmt.Println()
	fmt.Println("========================================")
	fmt.Println("✨ Web App Started Successfully!")
	fmt.Println("========================================")
	fmt.Println()
	printAccessInfo()
}

func isContainerRunning() bool {
	cmd := exec.Command("docker", "ps", "--filter", fmt.Sprintf("name=^%s$", ContainerName), "--format", "{{.Names}}")
	output, err := cmd.Output()
	if err != nil {
		return false
	}
	return len(output) > 0
}

func waitForHealthy(timeoutSeconds int) error {
	timeout := time.After(time.Duration(timeoutSeconds) * time.Second)
	ticker := time.NewTicker(2 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-timeout:
			return fmt.Errorf("timeout waiting for services to be healthy")
		case <-ticker.C:
			cmd := exec.Command("docker", "inspect", "--format", "{{.State.Health.Status}}", ContainerName)
			output, err := cmd.Output()
			if err != nil {
				// Container might not have health check
				// Check if it's just running
				if isContainerRunning() {
					return nil
				}
				continue
			}

			status := string(output)
			if status == "healthy\n" {
				return nil
			}
		}
	}
}

func printAccessInfo() {
	fmt.Println("Access the web interface at:")
	fmt.Printf("  🌐 Frontend:  http://localhost:%s\n", FrontendPort)
	fmt.Printf("  🔌 Backend:   http://localhost:%s\n", BackendPort)
	fmt.Println()
	fmt.Println("Useful commands:")
	fmt.Println("  planton webapp status    # Check service status")
	fmt.Println("  planton webapp logs      # View service logs")
	fmt.Println("  planton webapp stop      # Stop the web app")
	fmt.Println()
}
