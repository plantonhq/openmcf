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

var CloudResourceGetCmd = &cobra.Command{
	Use:   "cloud-resource:get",
	Short: "get a cloud resource by ID",
	Long:  "Retrieve detailed information about a cloud resource by providing its unique ID.",
	Run:   cloudResourceGetHandler,
}

func init() {
	CloudResourceGetCmd.Flags().StringP("id", "i", "", "unique identifier of the cloud resource (required)")
	CloudResourceGetCmd.MarkFlagRequired("id")
}

func cloudResourceGetHandler(cmd *cobra.Command, args []string) {
	// Get resource ID
	id, _ := cmd.Flags().GetString("id")
	if id == "" {
		fmt.Println("Error: --id flag is required. Provide the cloud resource ID")
		fmt.Println("Usage: openmcf cloud-resource:get --id=<resource-id>")
		os.Exit(1)
	}

	// Get backend URL from configuration
	backendURL, err := GetBackendURL()
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}

	// Create Connect-RPC client
	client := cloudresourcev1connect.NewCloudResourceQueryControllerClient(
		http.DefaultClient,
		backendURL,
	)

	// Prepare request
	req := &cloudresourcev1.GetCloudResourceRequest{
		Id: id,
	}

	// Create context with timeout
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// Make the API call
	resp, err := client.Get(ctx, connect.NewRequest(req))
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
		fmt.Printf("Error getting cloud resource: %v\n", err)
		os.Exit(1)
	}

	resource := resp.Msg.Resource

	// Display resource details
	fmt.Println("Cloud Resource Details:")
	fmt.Println("======================")
	fmt.Printf("ID:         %s\n", resource.Id)
	fmt.Printf("Name:       %s\n", resource.Name)
	fmt.Printf("Kind:       %s\n", resource.Kind)
	if resource.CreatedAt != nil {
		fmt.Printf("Created At: %s\n", resource.CreatedAt.AsTime().Format("2006-01-02 15:04:05"))
	}
	if resource.UpdatedAt != nil {
		fmt.Printf("Updated At: %s\n", resource.UpdatedAt.AsTime().Format("2006-01-02 15:04:05"))
	}
	fmt.Println("\nManifest:")
	fmt.Println("----------")
	fmt.Println(resource.Manifest)
}
