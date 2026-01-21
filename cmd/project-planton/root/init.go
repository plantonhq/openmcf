package root

import (
	"fmt"
	"os"

	"github.com/plantonhq/project-planton/apis/org/project_planton/shared"
	"github.com/plantonhq/project-planton/apis/org/project_planton/shared/iac/terraform"
	"github.com/plantonhq/project-planton/internal/cli/cliprint"
	"github.com/plantonhq/project-planton/internal/cli/flag"
	climanifest "github.com/plantonhq/project-planton/internal/cli/manifest"
	"github.com/plantonhq/project-planton/internal/cli/prompt"
	"github.com/plantonhq/project-planton/internal/cli/workspace"
	"github.com/plantonhq/project-planton/internal/manifest"
	"github.com/plantonhq/project-planton/pkg/crkreflect"
	"github.com/plantonhq/project-planton/pkg/iac/localmodule"
	"github.com/plantonhq/project-planton/pkg/iac/provisioner"
	"github.com/plantonhq/project-planton/pkg/iac/pulumi/pulumistack"
	"github.com/plantonhq/project-planton/pkg/iac/stackinput"
	"github.com/plantonhq/project-planton/pkg/iac/stackinput/stackinputproviderconfig"
	"github.com/plantonhq/project-planton/pkg/iac/tofu/tfbackend"
	"github.com/plantonhq/project-planton/pkg/iac/tofu/tofumodule"
	"github.com/plantonhq/project-planton/pkg/kubernetes/kubecontext"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"google.golang.org/protobuf/proto"
)

var Init = &cobra.Command{
	Use:   "init",
	Short: "initialize backend/stack using the provisioner specified in manifest",
	Long: `Initialize infrastructure backend or stack by automatically routing to the appropriate provisioner
(Pulumi, Tofu, or Terraform) based on the manifest label 'project-planton.org/provisioner'.

If the provisioner label is not present, you will be prompted to select one interactively.`,
	Example: `
	# Initialize with manifest file
	project-planton init -f manifest.yaml
	project-planton init --manifest manifest.yaml

	# Initialize with kustomize
	project-planton init --kustomize-dir _kustomize --overlay prod

	# Initialize with tofu-specific backend config
	project-planton init -f manifest.yaml --backend-type s3 --backend-config bucket=my-bucket
	`,
	Run: initHandler,
}

func init() {
	pwd, err := os.Getwd()
	if err != nil {
		log.Fatal("failed to get current working directory")
	}

	// Use StringP to support both --manifest and -f
	Init.PersistentFlags().StringP(string(flag.Manifest), "f", "", "path of the deployment-component manifest file")

	Init.PersistentFlags().String(string(flag.InputDir), "", "directory containing target.yaml and credential yaml files")
	Init.PersistentFlags().String(string(flag.KustomizeDir), "", "directory containing kustomize configuration")
	Init.PersistentFlags().String(string(flag.Overlay), "", "kustomize overlay to use (e.g., prod, dev, staging)")
	Init.PersistentFlags().String(string(flag.ModuleDir), pwd, "directory containing the provisioner module")
	Init.PersistentFlags().StringToString(string(flag.Set), map[string]string{}, "override resource manifest values using key=value pairs")

	// Pulumi-specific flags
	Init.PersistentFlags().String(string(flag.Stack), "", "pulumi stack fqdn in the format of <org>/<project>/<stack>")

	// Tofu/Terraform-specific flags
	Init.PersistentFlags().String(string(flag.BackendType), terraform.TerraformBackendType_local.String(),
		"Specifies the backend type (Tofu/Terraform) - 'local', 's3', 'gcs', 'azurerm', etc.")
	Init.PersistentFlags().StringArray(string(flag.BackendConfig), []string{},
		"Backend configuration key=value pairs (Tofu/Terraform)")

	// Staging/cleanup flags
	Init.PersistentFlags().Bool(string(flag.NoCleanup), false, "Do not cleanup the workspace copy after execution (keeps cloned modules)")
	Init.PersistentFlags().String(string(flag.ModuleVersion), "",
		"Checkout a specific version (tag, branch, or commit SHA) of the IaC modules in the workspace copy.\n"+
			"This allows using a different module version than what's in the staging area without affecting it.")

	// Kubernetes context flag
	Init.PersistentFlags().String(string(flag.KubeContext), "", "kubectl context to use for Kubernetes deployments (overrides manifest label)")

	// Provider credential flags
	Init.PersistentFlags().String(string(flag.AtlasProviderConfig), "", "path of the mongodb-atlas-credential file")
	Init.PersistentFlags().String(string(flag.Auth0ProviderConfig), "", "path of the auth0-credential file")
	Init.PersistentFlags().String(string(flag.AwsProviderConfig), "", "path of the aws-credential file")
	Init.PersistentFlags().String(string(flag.AzureProviderConfig), "", "path of the azure-credential file")
	Init.PersistentFlags().String(string(flag.CloudflareProviderConfig), "", "path of the cloudflare-credential file")
	Init.PersistentFlags().String(string(flag.ConfluentProviderConfig), "", "path of the confluent-credential file")
	Init.PersistentFlags().String(string(flag.GcpProviderConfig), "", "path of the gcp-credential file")
	Init.PersistentFlags().String(string(flag.KubernetesProviderConfig), "", "path of the yaml file containing the kubernetes cluster configuration")
	Init.PersistentFlags().String(string(flag.SnowflakeProviderConfig), "", "path of the snowflake-credential file")
}

func initHandler(cmd *cobra.Command, args []string) {
	moduleDir, err := cmd.Flags().GetString(string(flag.ModuleDir))
	flag.HandleFlagErrAndValue(err, flag.ModuleDir, moduleDir)

	valueOverrides, err := cmd.Flags().GetStringToString(string(flag.Set))
	flag.HandleFlagErr(err, flag.Set)

	// Check which manifest source is being used for informative messages
	kustomizeDir, _ := cmd.Flags().GetString(string(flag.KustomizeDir))
	overlay, _ := cmd.Flags().GetString(string(flag.Overlay))

	if kustomizeDir != "" && overlay != "" {
		cliprint.PrintStep(fmt.Sprintf("Building manifest from kustomize overlay: %s", overlay))
	} else {
		cliprint.PrintStep("Loading manifest...")
	}

	// Resolve manifest path with priority: --manifest > --input-dir > --kustomize-dir + --overlay
	targetManifestPath, isTemp, err := climanifest.ResolveManifestPath(cmd)
	if err != nil {
		cliprint.PrintError(fmt.Sprintf("Failed to resolve manifest: %v", err))
		os.Exit(1)
	}
	if isTemp {
		defer os.Remove(targetManifestPath)
	}

	cliprint.PrintSuccess("Manifest loaded")

	// Apply value overrides if any (creates new temp file if overrides exist)
	if len(valueOverrides) > 0 {
		cliprint.PrintStep(fmt.Sprintf("Applying %d field override(s)...", len(valueOverrides)))
	}

	finalManifestPath, isTempOverrides, err := manifest.ApplyOverridesToFile(targetManifestPath, valueOverrides)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	if isTempOverrides {
		defer os.Remove(finalManifestPath)
		targetManifestPath = finalManifestPath
		cliprint.PrintSuccess("Overrides applied")
	}

	// Validate manifest before proceeding (after overrides are applied)
	cliprint.PrintStep("Validating manifest...")
	if err := manifest.Validate(targetManifestPath); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	cliprint.PrintSuccess("Manifest validated")

	// Load manifest to extract provisioner
	cliprint.PrintStep("Detecting provisioner...")
	manifestObject, err := manifest.LoadWithOverrides(targetManifestPath, valueOverrides)
	if err != nil {
		cliprint.PrintError(fmt.Sprintf("Failed to load manifest: %v", err))
		os.Exit(1)
	}

	// Extract provisioner from manifest
	provType, err := provisioner.ExtractFromManifest(manifestObject)
	if err != nil {
		cliprint.PrintError(fmt.Sprintf("Invalid provisioner in manifest: %v", err))
		os.Exit(1)
	}

	// If provisioner not specified in manifest, prompt user
	if provType == provisioner.ProvisionerTypeUnspecified {
		cliprint.PrintInfo("Provisioner not specified in manifest")
		provType, err = prompt.PromptForProvisioner()
		if err != nil {
			cliprint.PrintError(fmt.Sprintf("Failed to get provisioner: %v", err))
			os.Exit(1)
		}
	}

	cliprint.PrintSuccess(fmt.Sprintf("Using provisioner: %s", provType.String()))

	// Resolve kube context: flag takes priority over manifest label
	kubeCtx, _ := cmd.Flags().GetString(string(flag.KubeContext))
	if kubeCtx == "" {
		kubeCtx = kubecontext.ExtractFromManifest(manifestObject)
	}
	if kubeCtx != "" {
		cliprint.PrintInfo(fmt.Sprintf("Using kubectl context: %s", kubeCtx))
	}

	// Handle --local-module flag: derive module directory from local project-planton repo
	localModule, _ := cmd.Flags().GetBool(string(flag.LocalModule))
	if localModule {
		var iacProv shared.IacProvisioner
		switch provType {
		case provisioner.ProvisionerTypePulumi:
			iacProv = shared.IacProvisioner_pulumi
		case provisioner.ProvisionerTypeTofu, provisioner.ProvisionerTypeTerraform:
			iacProv = shared.IacProvisioner_terraform
		}
		moduleDir, err = localmodule.GetModuleDir(targetManifestPath, cmd, iacProv)
		if err != nil {
			if lmErr, ok := err.(*localmodule.Error); ok {
				lmErr.PrintError()
			} else {
				cliprint.PrintError(err.Error())
			}
			os.Exit(1)
		}
	}

	// Prepare provider configs
	cliprint.PrintStep("Preparing execution...")
	providerConfigOptions, err := stackinputproviderconfig.BuildWithFlags(cmd.Flags())
	if err != nil {
		cliprint.PrintError(fmt.Sprintf("Failed to build credential options: %v", err))
		os.Exit(1)
	}
	cliprint.PrintSuccess("Execution prepared")

	// Route to appropriate provisioner
	switch provType {
	case provisioner.ProvisionerTypePulumi:
		initWithPulumi(cmd, moduleDir, targetManifestPath, valueOverrides)
	case provisioner.ProvisionerTypeTofu:
		initWithTofu(cmd, moduleDir, targetManifestPath, valueOverrides, kubeCtx, manifestObject, providerConfigOptions)
	case provisioner.ProvisionerTypeTerraform:
		initWithTerraform(cmd, moduleDir, targetManifestPath, valueOverrides, kubeCtx, manifestObject, providerConfigOptions)
	default:
		cliprint.PrintError("Unknown provisioner type")
		os.Exit(1)
	}
}

func initWithPulumi(cmd *cobra.Command, moduleDir, targetManifestPath string, valueOverrides map[string]string) {
	// Stack can be provided via flag or extracted from manifest
	stackFqdn, err := cmd.Flags().GetString(string(flag.Stack))
	flag.HandleFlagErr(err, flag.Stack)

	noCleanup, _ := cmd.Flags().GetBool(string(flag.NoCleanup))
	moduleVersion, _ := cmd.Flags().GetString(string(flag.ModuleVersion))

	err = pulumistack.Init(moduleDir, stackFqdn, targetManifestPath, valueOverrides, moduleVersion, noCleanup)
	if err != nil {
		cliprint.PrintPulumiFailure()
		os.Exit(1)
	}
	cliprint.PrintPulumiSuccess()
}

func initWithTofu(cmd *cobra.Command, moduleDir, targetManifestPath string, valueOverrides map[string]string,
	kubeContext string, manifestObject proto.Message, providerConfigOptions []stackinputproviderconfig.StackInputProviderConfigOption) {

	backendTypeString, err := cmd.Flags().GetString(string(flag.BackendType))
	flag.HandleFlagErrAndValue(err, flag.BackendType, backendTypeString)

	backendConfigList, err := cmd.Flags().GetStringArray(string(flag.BackendConfig))
	flag.HandleFlagErr(err, flag.BackendConfig)

	backendType := tfbackend.BackendTypeFromString(backendTypeString)

	// Extract kind name for module path resolution
	kindName, err := crkreflect.ExtractKindFromProto(manifestObject)
	if err != nil {
		cliprint.PrintError(fmt.Sprintf("Failed to extract kind name from manifest proto: %v", err))
		os.Exit(1)
	}

	noCleanup, _ := cmd.Flags().GetBool(string(flag.NoCleanup))
	moduleVersion, _ := cmd.Flags().GetString(string(flag.ModuleVersion))

	pathResult, err := tofumodule.GetModulePath(moduleDir, kindName, moduleVersion, noCleanup)
	if err != nil {
		cliprint.PrintError(fmt.Sprintf("Failed to get tofu module directory: %v", err))
		os.Exit(1)
	}

	// Setup cleanup to run after execution
	if pathResult.ShouldCleanup {
		defer func() {
			if cleanupErr := pathResult.CleanupFunc(); cleanupErr != nil {
				fmt.Printf("Warning: failed to cleanup workspace copy: %v\n", cleanupErr)
			}
		}()
	}

	tofuModulePath := pathResult.ModulePath

	// Build stack input YAML
	opts := stackinputproviderconfig.StackInputProviderConfigOptions{}
	for _, opt := range providerConfigOptions {
		opt(&opts)
	}

	stackInputYaml, err := stackinput.BuildStackInputYaml(manifestObject, opts)
	if err != nil {
		cliprint.PrintError(fmt.Sprintf("Failed to build stack input yaml: %v", err))
		os.Exit(1)
	}

	workspaceDir, err := workspace.GetWorkspaceDir()
	if err != nil {
		cliprint.PrintError("Failed to get workspace directory")
		os.Exit(1)
	}

	providerConfigEnvVars, err := tofumodule.GetProviderConfigEnvVars(stackInputYaml, workspaceDir, kubeContext)
	if err != nil {
		cliprint.PrintError(fmt.Sprintf("Failed to get credential env vars: %v", err))
		os.Exit(1)
	}

	cliprint.PrintHandoff("OpenTofu")

	err = tofumodule.Init(
		"tofu",
		tofuModulePath,
		manifestObject,
		backendType,
		backendConfigList,
		providerConfigEnvVars,
		false,
		nil,
	)
	if err != nil {
		cliprint.PrintTofuFailure()
		os.Exit(1)
	}
	cliprint.PrintTofuSuccess()
}

func initWithTerraform(cmd *cobra.Command, moduleDir, targetManifestPath string, valueOverrides map[string]string,
	kubeContext string, manifestObject proto.Message, providerConfigOptions []stackinputproviderconfig.StackInputProviderConfigOption) {

	backendTypeString, err := cmd.Flags().GetString(string(flag.BackendType))
	flag.HandleFlagErrAndValue(err, flag.BackendType, backendTypeString)

	backendConfigList, err := cmd.Flags().GetStringArray(string(flag.BackendConfig))
	flag.HandleFlagErr(err, flag.BackendConfig)

	backendType := tfbackend.BackendTypeFromString(backendTypeString)

	kindName, err := crkreflect.ExtractKindFromProto(manifestObject)
	if err != nil {
		cliprint.PrintError(fmt.Sprintf("Failed to extract kind name from manifest proto: %v", err))
		os.Exit(1)
	}

	noCleanup, _ := cmd.Flags().GetBool(string(flag.NoCleanup))
	moduleVersion, _ := cmd.Flags().GetString(string(flag.ModuleVersion))

	pathResult, err := tofumodule.GetModulePath(moduleDir, kindName, moduleVersion, noCleanup)
	if err != nil {
		cliprint.PrintError(fmt.Sprintf("Failed to get terraform module directory: %v", err))
		os.Exit(1)
	}

	if pathResult.ShouldCleanup {
		defer func() {
			if cleanupErr := pathResult.CleanupFunc(); cleanupErr != nil {
				fmt.Printf("Warning: failed to cleanup workspace copy: %v\n", cleanupErr)
			}
		}()
	}

	modulePath := pathResult.ModulePath

	opts := stackinputproviderconfig.StackInputProviderConfigOptions{}
	for _, opt := range providerConfigOptions {
		opt(&opts)
	}

	stackInputYaml, err := stackinput.BuildStackInputYaml(manifestObject, opts)
	if err != nil {
		cliprint.PrintError(fmt.Sprintf("Failed to build stack input yaml: %v", err))
		os.Exit(1)
	}

	workspaceDir, err := workspace.GetWorkspaceDir()
	if err != nil {
		cliprint.PrintError("Failed to get workspace directory")
		os.Exit(1)
	}

	providerConfigEnvVars, err := tofumodule.GetProviderConfigEnvVars(stackInputYaml, workspaceDir, kubeContext)
	if err != nil {
		cliprint.PrintError(fmt.Sprintf("Failed to get credential env vars: %v", err))
		os.Exit(1)
	}

	cliprint.PrintHandoff("Terraform")

	err = tofumodule.Init(
		"terraform",
		modulePath,
		manifestObject,
		backendType,
		backendConfigList,
		providerConfigEnvVars,
		false,
		nil,
	)
	if err != nil {
		cliprint.PrintTerraformFailure()
		os.Exit(1)
	}
	cliprint.PrintTerraformSuccess()
}
