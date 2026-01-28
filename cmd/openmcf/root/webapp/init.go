package webapp

import (
	"fmt"
	"os"
	"os/exec"
	"runtime"
	"strings"

	"github.com/plantonhq/openmcf/cmd/openmcf/root"
	"github.com/spf13/cobra"
)

const (
	DockerImageName   = "ghcr.io/plantonhq/openmcf"
	DockerImageTag    = "latest"
	ContainerName     = "openmcf-webapp"
	MongoDBVolume     = "openmcf-mongodb-data"
	PulumiVolume      = "openmcf-pulumi-state"
	GoCacheVolume     = "openmcf-go-cache"
	BackendPort       = "50051"
	FrontendPort      = "3000"
	DefaultBackendURL = "http://localhost:50051"
)

var InitCmd = &cobra.Command{
	Use:   "init",
	Short: "initialize and configure the OpenMCF web app",
	Long: `Initialize the OpenMCF web app by:
  - Checking Docker availability
  - Pulling the unified Docker image
  - Creating Docker volumes for data persistence
  - Creating and starting the container
  - Configuring the CLI to use the local backend`,
	Run: initHandler,
}

func initHandler(cmd *cobra.Command, args []string) {
	fmt.Println("========================================")
	fmt.Println("🚀 OpenMCF Web App Initialization")
	fmt.Println("========================================")
	fmt.Println()

	// Step 1: Check Docker availability
	fmt.Println("📋 Step 1/5: Checking Docker availability...")
	if err := checkDockerAvailable(); err != nil {
		printDockerInstallInstructions()
		os.Exit(1)
	}
	fmt.Println("✅ Docker is available and running")
	fmt.Println()

	// Step 2: Check if container already exists
	fmt.Println("📋 Step 2/5: Checking for existing installation...")
	if containerExists() {
		fmt.Printf("⚠️  Container '%s' already exists.\n", ContainerName)
		fmt.Println("   To reinitialize, first run: planton webapp uninstall")
		os.Exit(1)
	}
	fmt.Println("✅ No existing installation found")
	fmt.Println()

	// Step 3: Pull Docker image
	fmt.Println("📋 Step 3/5: Pulling Docker image...")
	fullImageName := fmt.Sprintf("%s:%s", DockerImageName, DockerImageTag)
	fmt.Printf("   Pulling %s (this may take a few minutes)...\n", fullImageName)
	if err := pullDockerImage(fullImageName); err != nil {
		fmt.Printf("❌ Failed to pull Docker image: %v\n", err)
		os.Exit(1)
	}
	fmt.Println("✅ Docker image pulled successfully")
	fmt.Println()

	// Step 4: Create volumes and container
	fmt.Println("📋 Step 4/5: Creating Docker volumes and container...")
	if err := createVolumes(); err != nil {
		fmt.Printf("❌ Failed to create volumes: %v\n", err)
		os.Exit(1)
	}
	fmt.Println("   ✓ Created MongoDB data volume")
	fmt.Println("   ✓ Created Pulumi state volume")
	fmt.Println("   ✓ Created Go cache volume")

	if err := createContainer(fullImageName); err != nil {
		fmt.Printf("❌ Failed to create container: %v\n", err)
		os.Exit(1)
	}
	fmt.Println("   ✓ Created container")
	fmt.Println("✅ Container created successfully")
	fmt.Println()

	// Step 5: Configure CLI
	fmt.Println("📋 Step 5/5: Configuring CLI...")
	if err := configureBackendURL(); err != nil {
		fmt.Printf("⚠️  Warning: Failed to configure backend URL: %v\n", err)
		fmt.Printf("   You can manually configure it later with: planton config set backend-url %s\n", DefaultBackendURL)
	} else {
		fmt.Println("✅ CLI configured to use local backend")
	}
	fmt.Println()

	// Success message
	fmt.Println("========================================")
	fmt.Println("✨ Initialization Complete!")
	fmt.Println("========================================")
	fmt.Println()
	fmt.Println("Next steps:")
	fmt.Printf("  1. Start the web app:     planton webapp start\n")
	fmt.Printf("  2. Check status:          planton webapp status\n")
	fmt.Printf("  3. View logs:             planton webapp logs\n")
	fmt.Println()
	fmt.Println("Once started, access the web interface at:")
	fmt.Printf("  Frontend:  http://localhost:%s\n", FrontendPort)
	fmt.Printf("  Backend:   http://localhost:%s\n", BackendPort)
	fmt.Println()
}

func checkDockerAvailable() error {
	// Check if docker command exists
	_, err := exec.LookPath("docker")
	if err != nil {
		return fmt.Errorf("docker command not found")
	}

	// Check if Docker daemon is running
	cmd := exec.Command("docker", "info")
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("docker daemon is not running")
	}

	return nil
}

func printDockerInstallInstructions() {
	fmt.Println()
	fmt.Println("❌ Error: Docker Engine is not installed or not running")
	fmt.Println()
	fmt.Println("OpenMCF web app requires Docker Engine to run.")
	fmt.Println()
	fmt.Println("Installation instructions:")

	switch runtime.GOOS {
	case "darwin":
		fmt.Println("  macOS:    brew install docker docker-compose")
		fmt.Println("            Or install Docker Desktop: https://docker.com/products/docker-desktop")
	case "linux":
		fmt.Println("  Linux:    https://docs.docker.com/engine/install/")
	case "windows":
		fmt.Println("  Windows:  https://docs.docker.com/desktop/install/windows-install/")
	default:
		fmt.Println("  Visit:    https://docs.docker.com/engine/install/")
	}

	fmt.Println()
	fmt.Println("After installation, ensure Docker daemon is running:")
	fmt.Println("  docker info")
	fmt.Println()
}

func containerExists() bool {
	cmd := exec.Command("docker", "ps", "-a", "--filter", fmt.Sprintf("name=^%s$", ContainerName), "--format", "{{.Names}}")
	output, err := cmd.Output()
	if err != nil {
		return false
	}
	return strings.TrimSpace(string(output)) == ContainerName
}

func pullDockerImage(imageName string) error {
	cmd := exec.Command("docker", "pull", imageName)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

func createVolumes() error {
	volumes := []string{MongoDBVolume, PulumiVolume, GoCacheVolume}
	for _, volume := range volumes {
		cmd := exec.Command("docker", "volume", "create", volume)
		if err := cmd.Run(); err != nil {
			return fmt.Errorf("failed to create volume %s: %w", volume, err)
		}
	}
	return nil
}

func createContainer(imageName string) error {
	args := []string{
		"create",
		"--name", ContainerName,
		"-p", fmt.Sprintf("%s:%s", FrontendPort, FrontendPort),
		"-p", fmt.Sprintf("%s:%s", BackendPort, BackendPort),
		"-v", fmt.Sprintf("%s:/data/db", MongoDBVolume),
		"-v", fmt.Sprintf("%s:/home/appuser/.pulumi", PulumiVolume),
		"-v", fmt.Sprintf("%s:/home/appuser/go", GoCacheVolume),
		"--restart", "unless-stopped",
		imageName,
	}

	cmd := exec.Command("docker", args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

func configureBackendURL() error {
	config, err := root.LoadConfigPublic()
	if err != nil {
		return err
	}

	config.BackendURL = DefaultBackendURL
	config.WebAppContainerID = ContainerName
	config.WebAppVersion = DockerImageTag

	return root.SaveConfigPublic(config)
}
