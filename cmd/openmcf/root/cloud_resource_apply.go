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
	"gopkg.in/yaml.v3"
)

var CloudResourceApplyCmd = &cobra.Command{
	Use:   "cloud-resource:apply",
	Short: "create or update a cloud resource from YAML manifest (upsert)",
	Long:  "Apply a cloud resource by providing a YAML manifest file. If a resource with the same name and kind exists, it will be updated. Otherwise, a new resource will be created.",
	Run:   cloudResourceApplyHandler,
}

func init() {
	CloudResourceApplyCmd.Flags().StringP("arg", "a", "", "path to the YAML manifest file (required)")
	CloudResourceApplyCmd.MarkFlagRequired("arg")
}

func cloudResourceApplyHandler(cmd *cobra.Command, args []string) {
	// Get YAML file path
	yamlFile, _ := cmd.Flags().GetString("arg")
	if yamlFile == "" {
		fmt.Println("Error: --arg flag is required. Provide path to YAML manifest file")
		fmt.Println("Usage: openmcf cloud-resource:apply --arg=<yaml-file>")
		os.Exit(1)
	}

	// Read YAML file
	manifestContent, err := os.ReadFile(yamlFile)
	if err != nil {
		fmt.Printf("Error: Failed to read YAML file '%s': %v\n", yamlFile, err)
		os.Exit(1)
	}

	// Parse YAML to validate format and extract name and kind for display
	var yamlData map[string]interface{}
	if err := yaml.Unmarshal(manifestContent, &yamlData); err != nil {
		fmt.Printf("Error: Invalid YAML format: %v\n", err)
		os.Exit(1)
	}

	// Extract kind for validation
	kind, ok := yamlData["kind"].(string)
	if !ok || kind == "" {
		fmt.Println("Error: Manifest must contain 'kind' field")
		os.Exit(1)
	}

	// Extract metadata.name for validation
	metadata, ok := yamlData["metadata"].(map[string]interface{})
	if !ok {
		fmt.Println("Error: Manifest must contain 'metadata' field")
		os.Exit(1)
	}

	name, ok := metadata["name"].(string)
	if !ok || name == "" {
		fmt.Println("Error: Manifest must contain 'metadata.name' field")
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

	// Create context with timeout
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// Use ApplyCloudResource API which handles both create and update automatically
	// It checks if resource exists by (kind + metadata.name) and performs:
	// - CREATE if resource doesn't exist + triggers Pulumi deployment
	// - UPDATE if resource exists + triggers Pulumi deployment
	fmt.Printf("Applying cloud resource: kind=%s, name=%s\n", kind, name)
	fmt.Printf("Checking if resource exists...\n")

	applyReq := &cloudresourcev1.ApplyCloudResourceRequest{
		Manifest: string(manifestContent),
	}

	applyResp, err := client.Apply(ctx, connect.NewRequest(applyReq))
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
		fmt.Printf("Error applying cloud resource: %v\n", err)
		os.Exit(1)
	}

	resource := applyResp.Msg.Resource
	created := applyResp.Msg.Created

	// Display success message
	var action string
	if created {
		action = "Created"
		fmt.Printf("✅ Cloud resource created successfully!\n")
	} else {
		action = "Updated"
		fmt.Printf("✅ Cloud resource updated successfully!\n")
	}

	fmt.Printf("\nAction: %s\n", action)
	fmt.Printf("ID: %s\n", resource.Id)
	fmt.Printf("Name: %s\n", resource.Name)
	fmt.Printf("Kind: %s\n", resource.Kind)
	if resource.CreatedAt != nil {
		fmt.Printf("Created At: %s\n", resource.CreatedAt.AsTime().Format("2006-01-02 15:04:05"))
	}
	if resource.UpdatedAt != nil {
		fmt.Printf("Updated At: %s\n", resource.UpdatedAt.AsTime().Format("2006-01-02 15:04:05"))
	}

	fmt.Printf("\n🚀 Pulumi deployment has been triggered automatically.\n")
	fmt.Printf("   Deployment is running in the background.\n")
	fmt.Printf("   Use 'openmcf stack-update:list' to check deployment status.\n")
}
