package tofumodule

import (
	"fmt"

	"github.com/pkg/errors"
	"github.com/plantonhq/project-planton/apis/org/project_planton/shared/iac/terraform"
	"github.com/plantonhq/project-planton/internal/cli/workspace"
	"github.com/plantonhq/project-planton/internal/manifest"
	"github.com/plantonhq/project-planton/pkg/crkreflect"
	"github.com/plantonhq/project-planton/pkg/iac/stackinput"
	"github.com/plantonhq/project-planton/pkg/iac/stackinput/stackinputproviderconfig"
	"github.com/plantonhq/project-planton/pkg/iac/tofu/backendconfig"
	"github.com/plantonhq/project-planton/pkg/iac/tofu/tfbackend"
	log "github.com/sirupsen/logrus"
)

// RunCommand executes an HCL-based IaC operation (init + operation) using the specified binary.
// The binaryName parameter specifies which CLI binary to use ("tofu" or "terraform").
func RunCommand(
	binaryName string,
	inputModuleDir string,
	targetManifestPath string,
	terraformOperation terraform.TerraformOperationType,
	valueOverrides map[string]string,
	isAutoApprove bool,
	isDestroyPlan bool,
	isReconfigure bool,
	moduleVersion string,
	noCleanup bool,
	kubeContext string,
	providerConfigOptions ...stackinputproviderconfig.StackInputProviderConfigOption,
) error {

	manifestObject, err := manifest.LoadWithOverrides(targetManifestPath, valueOverrides)
	if err != nil {
		return errors.Wrapf(err, "failed to override values in target manifest file")
	}

	// Extract backend configuration from manifest labels (optional)
	// Uses provisioner-specific labels (e.g., tofu.project-planton.org/backend.type)
	// with fallback to legacy terraform.* labels for backward compatibility
	var backendType terraform.TerraformBackendType = terraform.TerraformBackendType_local
	var backendConfigArgs []string

	tofuBackendConfig, err := backendconfig.ExtractFromManifest(manifestObject, binaryName)
	if err != nil {
		// Log but don't fail - backend config is optional
		log.Debugf("Could not extract %s backend config from manifest labels: %v", binaryName, err)
	}

	if tofuBackendConfig != nil {
		// Convert backend type string to enum
		backendType = tfbackend.BackendTypeFromString(tofuBackendConfig.BackendType)
		if backendType == terraform.TerraformBackendType_terraform_backend_type_unspecified {
			return errors.Errorf("unsupported backend type from manifest labels: %s", tofuBackendConfig.BackendType)
		}

		// Build backend config arguments based on backend type
		backendConfigArgs = buildBackendConfigArgs(tofuBackendConfig)
	} else {
		log.Debugf("No %s backend config in manifest labels, using default local backend", binaryName)
	}

	kindName, err := crkreflect.ExtractKindFromProto(manifestObject)
	if err != nil {
		return errors.Wrapf(err, "failed to extract kind name from manifest proto")
	}

	// Get module path using staging-based approach
	pathResult, err := GetModulePath(inputModuleDir, kindName, moduleVersion, noCleanup)
	if err != nil {
		return errors.Wrapf(err, "failed to get %s module directory", binaryName)
	}

	// Setup cleanup to run after execution
	if pathResult.ShouldCleanup {
		defer func() {
			if cleanupErr := pathResult.CleanupFunc(); cleanupErr != nil {
				fmt.Printf("Warning: failed to cleanup workspace copy: %v\n", cleanupErr)
			}
		}()
	}

	modulePath := pathResult.ModulePath

	// Gather credential options
	opts := stackinputproviderconfig.StackInputProviderConfigOptions{}
	for _, opt := range providerConfigOptions {
		opt(&opts)
	}

	stackInputYaml, err := stackinput.BuildStackInputYaml(manifestObject, opts)
	if err != nil {
		return errors.Wrap(err, "failed to build stack input yaml")
	}

	workspaceDir, err := workspace.GetWorkspaceDir()
	if err != nil {
		return errors.Wrap(err, "failed to get workspace directory")
	}

	providerConfigEnvVars, err := GetProviderConfigEnvVars(stackInputYaml, workspaceDir, kubeContext)
	if err != nil {
		return errors.Wrap(err, "failed to get provider config env vars")
	}

	// Initialize with backend configuration before any operation
	err = Init(binaryName, modulePath, manifestObject, backendType, backendConfigArgs,
		providerConfigEnvVars, isReconfigure, false, nil)
	if err != nil {
		return errors.Wrapf(err, "failed to initialize %s module", binaryName)
	}

	err = RunOperation(binaryName, modulePath, terraformOperation,
		isAutoApprove, isDestroyPlan, manifestObject,
		providerConfigEnvVars, false, nil)
	if err != nil {
		return errors.Wrapf(err, "failed to run %s operation", binaryName)
	}

	return nil
}

// buildBackendConfigArgs builds backend configuration arguments based on backend type
func buildBackendConfigArgs(config *backendconfig.TofuBackendConfig) []string {
	var args []string

	switch config.BackendType {
	case "s3":
		// S3 backend: bucket, key, and region
		if config.BackendBucket != "" {
			args = append(args, fmt.Sprintf("bucket=%s", config.BackendBucket))
		}
		if config.BackendKey != "" {
			args = append(args, fmt.Sprintf("key=%s", config.BackendKey))
		}
		if config.BackendRegion != "" {
			args = append(args, fmt.Sprintf("region=%s", config.BackendRegion))
		}

	case "gcs":
		// GCS backend: bucket and prefix (key is called prefix in GCS)
		if config.BackendBucket != "" {
			args = append(args, fmt.Sprintf("bucket=%s", config.BackendBucket))
		}
		if config.BackendKey != "" {
			args = append(args, fmt.Sprintf("prefix=%s", config.BackendKey))
		}

	case "azurerm":
		// Azure backend: container_name and key
		if config.BackendBucket != "" {
			args = append(args, fmt.Sprintf("container_name=%s", config.BackendBucket))
		}
		if config.BackendKey != "" {
			args = append(args, fmt.Sprintf("key=%s", config.BackendKey))
		}

	case "local":
		// Local backend doesn't need config args
		// The path is handled by terraform itself
	}

	return args
}
