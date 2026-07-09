package providerenvvars

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestLoadKubernetesEnvVars_ExportsBothKubeconfigEnvVars guards the Tofu/Pulumi parity
// fix: the generated kubeconfig must be advertised under BOTH KUBECONFIG (Pulumi) and
// KUBE_CONFIG_PATH (Terraform/OpenTofu hashicorp/kubernetes + helm). Dropping
// KUBE_CONFIG_PATH makes the tofu provider silently fall back to in-cluster auth.
func TestLoadKubernetesEnvVars_ExportsBothKubeconfigEnvVars(t *testing.T) {
	// Minimal valid GCP GKE provider config. protojson accepts the proto field names
	// (snake_case) and the enum value name as a string.
	providerConfigYaml := []byte(`
provider: gcp_gke
gcp_gke:
  cluster_endpoint: "34.100.155.147"
  cluster_ca_data: "dGVzdC1jYS1kYXRh"
  service_account_key: "{\"type\":\"service_account\"}"
`)

	cacheDir := t.TempDir()

	envVars, err := loadKubernetesEnvVars(providerConfigYaml, cacheDir)
	require.NoError(t, err)

	kubeconfig, ok := envVars["KUBECONFIG"]
	assert.True(t, ok, "KUBECONFIG must be set for the Pulumi kubernetes provider")

	kubeConfigPath, ok := envVars["KUBE_CONFIG_PATH"]
	assert.True(t, ok, "KUBE_CONFIG_PATH must be set for the Terraform/OpenTofu kubernetes provider")

	assert.Equal(t, kubeconfig, kubeConfigPath,
		"both env vars must point at the same generated kubeconfig file")

	// The kubeconfig file the env vars reference must actually exist on disk.
	_, statErr := os.Stat(kubeConfigPath)
	assert.NoError(t, statErr, "the kubeconfig file referenced by the env vars should exist")
}

// TestLoadKubernetesEnvVars_AwsEks proves the EKS arm end to end on this path: the
// rendered kubeconfig lands on disk (0600 -- it can carry credentials in exec env
// entries) with an exec entry pointing back at this very binary.
func TestLoadKubernetesEnvVars_AwsEks(t *testing.T) {
	providerConfigYaml := []byte(`
provider: aws_eks
aws_eks:
  cluster_name: demo-cluster
  cluster_endpoint: "https://ABC123.gr7.us-east-1.eks.amazonaws.com"
  cluster_ca_data: "dGVzdC1jYS1kYXRh"
  region: us-east-1
`)

	envVars, err := loadKubernetesEnvVars(providerConfigYaml, t.TempDir())
	require.NoError(t, err)

	kubeConfigPath := envVars["KUBECONFIG"]
	require.NotEmpty(t, kubeConfigPath)
	assert.Equal(t, kubeConfigPath, envVars["KUBE_CONFIG_PATH"])

	info, err := os.Stat(kubeConfigPath)
	require.NoError(t, err)
	assert.Equal(t, os.FileMode(0600), info.Mode().Perm())

	content, err := os.ReadFile(kubeConfigPath)
	require.NoError(t, err)

	executable, err := os.Executable()
	require.NoError(t, err)
	assert.Contains(t, string(content), executable,
		"the exec credential command must be the engine-spawning binary itself")
	assert.Contains(t, string(content), "client.authentication.k8s.io/v1")
}

// TestLoadKubernetesEnvVars_DigitalOceanDoks: the DOKS kubeconfig passes through to
// disk unchanged and is advertised under both engine env vars.
func TestLoadKubernetesEnvVars_DigitalOceanDoks(t *testing.T) {
	providerConfigYaml := []byte(`
provider: digital_ocean_doks
digital_ocean_doks:
  kube_config: |
    apiVersion: v1
    kind: Config
`)

	envVars, err := loadKubernetesEnvVars(providerConfigYaml, t.TempDir())
	require.NoError(t, err)

	content, err := os.ReadFile(envVars["KUBE_CONFIG_PATH"])
	require.NoError(t, err)
	assert.Contains(t, string(content), "kind: Config")
}
