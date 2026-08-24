package runner

import (
	"os"
	"path/filepath"
	"strings"

	"github.com/pkg/errors"
	"github.com/plantonhq/planton/internal/manifest"
	"github.com/plantonhq/planton/pkg/iac/stackinput"
	"github.com/plantonhq/planton/pkg/iac/stackinput/providerenvvars"
	"github.com/plantonhq/planton/pkg/iac/stackinput/stackinputproviderconfig"
	"github.com/plantonhq/planton/pkg/iac/tofu/generators"
	"github.com/plantonhq/planton/pkg/iac/tofu/tfbackend"
	"github.com/plantonhq/planton/pkg/iac/tofu/tfoverride"
	"github.com/plantonhq/planton/shared/iac/terraform"
)

// TerraformInput holds the prepared inputs for a Terraform E2E test run.
type TerraformInput struct {
	// TfvarsPath is the absolute path to the generated terraform.tfvars file.
	TfvarsPath string

	// EnvVars holds provider-specific environment variables (e.g., KUBECONFIG).
	EnvVars map[string]string
}

// BuildTerraformInput prepares a Terraform module working directory for E2E testing.
// It loads the manifest, generates a tfvars file, writes the backend configuration,
// writes the provider-override file when the component's provider-config fixture
// carries provider-block arguments, and extracts provider environment variables.
//
// The workDir must already contain the TF module files (copied by PrepareWorkDir).
// providerConfig is nil for the harness's default posture (ambient credentials,
// empty provider block); see LoadProviderConfigFixture.
func BuildTerraformInput(manifestPath, workDir string,
	providerConfig *stackinputproviderconfig.ProviderConfig) (*TerraformInput, error) {
	manifestObject, err := manifest.LoadManifest(manifestPath)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to load manifest from %s", manifestPath)
	}

	// Generate terraform.tfvars from the proto manifest.
	// The tfvars file is placed inside the working directory so tofu init
	// can find it alongside the module's .tf files.
	tfvarsPath := filepath.Join(workDir, "terraform.tfvars")
	if err := generators.WriteVarFile(manifestObject, tfvarsPath); err != nil {
		return nil, errors.Wrap(err, "failed to generate terraform.tfvars from manifest")
	}

	// Write backend.tf with local backend for ephemeral E2E state.
	if err := tfbackend.WriteBackendFile(workDir, terraform.TerraformBackendType_local); err != nil {
		return nil, errors.Wrap(err, "failed to write backend.tf")
	}

	// Build stack-input YAML to extract provider environment variables.
	// For Kubernetes on kind, this produces KUBECONFIG.
	// For cloud providers, this produces AWS_ACCESS_KEY_ID, GOOGLE_CREDENTIALS, etc.
	stackInputYaml, err := stackinput.BuildStackInputYaml(manifestObject, providerConfig)
	if err != nil {
		return nil, errors.Wrap(err, "failed to build stack-input YAML for provider env var extraction")
	}

	// Provider-block arguments (assume-role chain, default tags, ...) cannot
	// ride env vars; they reach the module through the generated override
	// file -- the same seam the CLI and the platform runner use. A no-op when
	// the fixture carries none (or there is no fixture); workDir is a
	// disposable per-test copy, so no cleanup is needed.
	if _, err := tfoverride.WriteProviderOverrideFile(workDir, stackInputYaml); err != nil {
		return nil, errors.Wrap(err, "failed to write provider override file")
	}

	providerEnvVarMap, err := providerenvvars.GetEnvVarsWithOptions(stackInputYaml, providerenvvars.Options{})
	if err != nil {
		return nil, errors.Wrap(err, "failed to extract provider environment variables")
	}

	// The Terraform Kubernetes provider uses KUBE_CONFIG_PATH (not KUBECONFIG).
	// This bridges a DIFFERENT case than providerenvvars.loadKubernetesEnvVars (which now
	// sets both names for connection-derived kubeconfigs): here the kind harness exports
	// KUBECONFIG into the process for an in-cluster test kubeconfig, so we forward it to
	// KUBE_CONFIG_PATH for the TF provider. Kept distinct on purpose.
	if kubeconfig := os.Getenv("KUBECONFIG"); kubeconfig != "" {
		if _, exists := providerEnvVarMap["KUBE_CONFIG_PATH"]; !exists {
			providerEnvVarMap["KUBE_CONFIG_PATH"] = kubeconfig
		}
	}

	return &TerraformInput{
		TfvarsPath: tfvarsPath,
		EnvVars:    providerEnvVarMap,
	}, nil
}

// PrepareWorkDir creates an isolated temporary copy of a TF module
// directory. Terraform state files (.terraform/, terraform.tfstate) live in
// the working directory, so each test needs its own copy to avoid state
// conflicts.
//
// THE SIBLING-DIRECTORY CONTRACT: modules may read staged files through
// relative paths OUTSIDE the module directory — the module-owned-CRD class
// stages its CRD files at iac/crds/ and both engines read `../crds`. A copy
// of the module directory ALONE makes fileset() return empty and for_each
// silently plan ZERO resources (caught live: a lane "passed" riding a
// previous lane's retained CRDs). The copy therefore reproduces the
// module's PARENT directory tree (iac/), excluding engine state and the
// sibling Pulumi module, and returns the module subdirectory inside it.
//
// Returns the working directory path and a cleanup function.
func PrepareWorkDir(sourceModuleDir string) (string, func(), error) {
	tempRoot, err := os.MkdirTemp("", "planton-e2e-tf-*")
	if err != nil {
		return "", nil, errors.Wrap(err, "failed to create temp directory for TF module")
	}

	cleanup := func() {
		os.RemoveAll(tempRoot)
	}

	sourceModuleDir = filepath.Clean(sourceModuleDir)
	parentDir := filepath.Dir(sourceModuleDir)

	// Never copied: engine-local state (a shared checkout may carry a
	// developer's .terraform plugin tree and state files) and the sibling
	// Pulumi module (its own engine's world — irrelevant to a TF run).
	skip := map[string]bool{
		".terraform": true,
		"pulumi":     true,
	}

	var copyTree func(src, dst string) error
	copyTree = func(src, dst string) error {
		entries, err := os.ReadDir(src)
		if err != nil {
			return errors.Wrapf(err, "failed to read directory %s", src)
		}
		if err := os.MkdirAll(dst, 0755); err != nil {
			return errors.Wrapf(err, "failed to create directory %s", dst)
		}
		for _, entry := range entries {
			name := entry.Name()
			if skip[name] || strings.HasPrefix(name, "terraform.tfstate") || name == "terraform.tfplan" {
				continue
			}
			srcPath := filepath.Join(src, name)
			dstPath := filepath.Join(dst, name)
			if entry.IsDir() {
				if err := copyTree(srcPath, dstPath); err != nil {
					return err
				}
				continue
			}
			content, err := os.ReadFile(srcPath)
			if err != nil {
				return errors.Wrapf(err, "failed to read %s", srcPath)
			}
			if err := os.WriteFile(dstPath, content, 0644); err != nil {
				return errors.Wrapf(err, "failed to write %s", dstPath)
			}
		}
		return nil
	}

	if err := copyTree(parentDir, tempRoot); err != nil {
		cleanup()
		return "", nil, err
	}

	return filepath.Join(tempRoot, filepath.Base(sourceModuleDir)), cleanup, nil
}
