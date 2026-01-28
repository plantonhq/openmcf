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

var CloudResourceUpdateCmd = &cobra.Command{
	Use:   "cloud-resource:update",
	Short: "update a cloud resource from YAML manifest",
	Long:  "Update an existing cloud resource by providing its ID and a YAML manifest file. The manifest's 'kind' and 'metadata.name' must match the existing resource.",
	Run:   cloudResourceUpdateHandler,
}

func init() {
	CloudResourceUpdateCmd.Flags().StringP("id", "i", "", "unique identifier of the cloud resource (required)")
	CloudResourceUpdateCmd.Flags().StringP("arg", "a", "", "path to the YAML manifest file (required)")
	CloudResourceUpdateCmd.MarkFlagRequired("id")
	CloudResourceUpdateCmd.MarkFlagRequired("arg")
}

func cloudResourceUpdateHandler(cmd *cobra.Command, args []string) {
	// Get resource ID
	id, _ := cmd.Flags().GetString("id")
	if id == "" {
		fmt.Println("Error: --id flag is required. Provide the cloud resource ID")
		fmt.Println("Usage: openmcf cloud-resource:update --id=<resource-id> --arg=<yaml-file>")
		os.Exit(1)
	}

	// Get YAML file path
	yamlFile, _ := cmd.Flags().GetString("arg")
	if yamlFile == "" {
		fmt.Println("Error: --arg flag is required. Provide path to YAML manifest file")
		fmt.Println("Usage: openmcf cloud-resource:update --id=<resource-id> --arg=<yaml-file>")
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
	req := &cloudresourcev1.UpdateCloudResourceRequest{
		Id:       id,
		Manifest: string(manifestContent),
	}

	// Create context with timeout
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// Make the API call
	resp, err := client.Update(ctx, connect.NewRequest(req))
	if err != nil {
		if connect.CodeOf(err) == connect.CodeUnavailable {
			fmt.Printf("Error: Cannot connect to backend service at %s. Please check:\n", backendURL)
			fmt.Printf("  1. The backend service is running\n")
			fmt.Printf("  2. The backend URL is correct\n")
			fmt.Printf("  3. Network connectivity\n")
			os.Exit(1)
		}
		if connect.CodeOf(err) == connect.CodeNotFound {
			fmt.Printf("Error: Cloud resource with ID '%s' not found\n", id)
			os.Exit(1)
		}
		if connect.CodeOf(err) == connect.CodeInvalidArgument {
			fmt.Printf("Error: Invalid manifest - %v\n", err)
			os.Exit(1)
		}
		fmt.Printf("Error updating cloud resource: %v\n", err)
		os.Exit(1)
	}

	resource := resp.Msg.Resource

	// Display success message
	fmt.Println("✅ Cloud resource updated successfully!")
	fmt.Printf("\nID: %s\n", resource.Id)
	fmt.Printf("Name: %s\n", resource.Name)
	fmt.Printf("Kind: %s\n", resource.Kind)
	if resource.UpdatedAt != nil {
		fmt.Printf("Updated At: %s\n", resource.UpdatedAt.AsTime().Format("2006-01-02 15:04:05"))
	}
}
