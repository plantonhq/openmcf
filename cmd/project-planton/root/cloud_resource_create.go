package root

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"time"

	"connectrpc.com/connect"
	cloudresourcev1 "github.com/plantonhq/openmcf/apis/org/openmcf/app/cloudresource/v1"
	cloudresourcev1connect "github.com/plantonhq/openmcf/apis/org/openmcf/app/cloudresource/v1/cloudresourcev1connect"
	"github.com/spf13/cobra"
)

var CloudResourceCreateCmd = &cobra.Command{
	Use:   "cloud-resource:create",
	Short: "create a new cloud resource from YAML manifest",
	Long:  "Create a new cloud resource by providing a YAML manifest file. The manifest must contain 'kind' and 'metadata.name' fields.",
	Run:   cloudResourceCreateHandler,
}

func init() {
	CloudResourceCreateCmd.Flags().StringP("arg", "a", "", "path to the YAML manifest file (required)")
	CloudResourceCreateCmd.MarkFlagRequired("arg")
}

func cloudResourceCreateHandler(cmd *cobra.Command, args []string) {
	// Get YAML file path
	yamlFile, _ := cmd.Flags().GetString("arg")
	if yamlFile == "" {
		fmt.Println("Error: --arg flag is required. Provide path to YAML manifest file")
		fmt.Println("Usage: openmcf cloud-resource:create --arg=<yaml-file>")
		os.Exit(1)
	}

	// Read YAML file
	manifestContent, err := os.ReadFile(yamlFile)
	if err != nil {
		fmt.Printf("Error: Failed to read YAML file '%s': %v\n", yamlFile, err)
		os.Exit(1)
	}

	// Get backend URL from configuration
	backendURL, err := GetBackendURL()
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}

	// Create Connect-RPC client
	client := cloudresourcev1connect.NewCloudResourceCommandControllerClient(
		http.DefaultClient,
		backendURL,
	)

	// Prepare request
	req := &cloudresourcev1.CreateCloudResourceRequest{
		Manifest: string(manifestContent),
	}

	// Create context with timeout
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// Make the API call
	resp, err := client.Create(ctx, connect.NewRequest(req))
	if err != nil {
		if connect.CodeOf(err) == connect.CodeUnavailable {
			fmt.Printf("Error: Cannot connect to backend service at %s. Please check:\n", backendURL)
			fmt.Printf("  1. The backend service is running\n")
			fmt.Printf("  2. The backend URL is correct\n")
			fmt.Printf("  3. Network connectivity\n")
			os.Exit(1)
		}
		if connect.CodeOf(err) == connect.CodeInvalidArgument {
			fmt.Printf("Error: Invalid manifest - %v\n", err)
			os.Exit(1)
		}
		fmt.Printf("Error creating cloud resource: %v\n", err)
		os.Exit(1)
	}

	resource := resp.Msg.Resource

	// Display success message
	fmt.Println("✅ Cloud resource created successfully!")
	fmt.Printf("\nID: %s\n", resource.Id)
	fmt.Printf("Name: %s\n", resource.Name)
	fmt.Printf("Kind: %s\n", resource.Kind)
	if resource.CreatedAt != nil {
		fmt.Printf("Created At: %s\n", resource.CreatedAt.AsTime().Format("2006-01-02 15:04:05"))
	}
}
