package providerdetect

import (
	"fmt"
	"strings"

	"github.com/plantonhq/project-planton/apis/org/project_planton/shared/cloudresourcekind"
)

// ProviderConfigExample returns an example YAML configuration for the given provider.
func ProviderConfigExample(provider cloudresourcekind.CloudResourceProvider) string {
	switch provider {
	case cloudresourcekind.CloudResourceProvider_atlas:
		return `public_key: "<your-atlas-public-key>"
private_key: "<your-atlas-private-key>"`

	case cloudresourcekind.CloudResourceProvider_auth0:
		return `domain: "<your-auth0-domain>"
client_id: "<your-auth0-client-id>"
client_secret: "<your-auth0-client-secret>"`

	case cloudresourcekind.CloudResourceProvider_aws:
		return `access_key_id: "<your-aws-access-key-id>"
secret_access_key: "<your-aws-secret-access-key>"
region: "us-east-1"`

	case cloudresourcekind.CloudResourceProvider_azure:
		return `client_id: "<your-azure-client-id>"
client_secret: "<your-azure-client-secret>"
tenant_id: "<your-azure-tenant-id>"
subscription_id: "<your-azure-subscription-id>"`

	case cloudresourcekind.CloudResourceProvider_cloudflare:
		return `api_token: "<your-cloudflare-api-token>"`

	case cloudresourcekind.CloudResourceProvider_confluent:
		return `api_key: "<your-confluent-api-key>"
api_secret: "<your-confluent-api-secret>"`

	case cloudresourcekind.CloudResourceProvider_gcp:
		return `service_account_key_base64: "<base64-encoded-service-account-json>"`

	case cloudresourcekind.CloudResourceProvider_kubernetes:
		return `provider: 1  # 1=GCP_GKE, 2=AWS_EKS, 3=AZURE_AKS
gcp_gke:
  cluster_endpoint: "<cluster-endpoint>"
  cluster_ca_data: "<base64-encoded-ca-cert>"
  service_account_key_base64: "<base64-encoded-service-account-json>"`

	case cloudresourcekind.CloudResourceProvider_open_fga:
		return `api_url: "<your-openfga-api-url>"
api_token: "<your-openfga-api-token>"`

	case cloudresourcekind.CloudResourceProvider_snowflake:
		return `account: "<your-snowflake-account>"
region: "<your-snowflake-region>"
username: "<your-snowflake-username>"
password: "<your-snowflake-password>"`

	default:
		return "# Provider config format not available"
	}
}

// ProviderConfigFilename returns the suggested filename for the provider config.
func ProviderConfigFilename(provider cloudresourcekind.CloudResourceProvider) string {
	switch provider {
	case cloudresourcekind.CloudResourceProvider_atlas:
		return "atlas-provider-config.yaml"
	case cloudresourcekind.CloudResourceProvider_auth0:
		return "auth0-provider-config.yaml"
	case cloudresourcekind.CloudResourceProvider_aws:
		return "aws-provider-config.yaml"
	case cloudresourcekind.CloudResourceProvider_azure:
		return "azure-provider-config.yaml"
	case cloudresourcekind.CloudResourceProvider_cloudflare:
		return "cloudflare-provider-config.yaml"
	case cloudresourcekind.CloudResourceProvider_confluent:
		return "confluent-provider-config.yaml"
	case cloudresourcekind.CloudResourceProvider_gcp:
		return "gcp-provider-config.yaml"
	case cloudresourcekind.CloudResourceProvider_kubernetes:
		return "kubernetes-provider-config.yaml"
	case cloudresourcekind.CloudResourceProvider_open_fga:
		return "openfga-provider-config.yaml"
	case cloudresourcekind.CloudResourceProvider_snowflake:
		return "snowflake-provider-config.yaml"
	default:
		return "provider-config.yaml"
	}
}

// MissingProviderConfigGuidance returns a helpful message when provider config is missing.
func MissingProviderConfigGuidance(result *DetectionResult) string {
	var sb strings.Builder

	sb.WriteString(fmt.Sprintf("The %s resource requires %s credentials.\n\n",
		result.KindName, ProviderDisplayName(result.Provider)))

	sb.WriteString(fmt.Sprintf("Create a file '%s' with:\n\n",
		ProviderConfigFilename(result.Provider)))

	// Indent the example
	example := ProviderConfigExample(result.Provider)
	for _, line := range strings.Split(example, "\n") {
		sb.WriteString("  " + line + "\n")
	}

	sb.WriteString("\nThen run:\n\n")
	sb.WriteString(fmt.Sprintf("  project-planton plan -f manifest.yaml -p %s\n",
		ProviderConfigFilename(result.Provider)))

	return sb.String()
}

// KindDetectionErrorGuidance returns a helpful message when kind detection fails.
func KindDetectionErrorGuidance() string {
	return `The manifest must contain valid 'apiVersion' and 'kind' fields:

  apiVersion: gcp.project-planton.org/v1
  kind: GkeCluster
  metadata:
    name: my-cluster
  spec:
    # ... resource configuration

Check your manifest file for:
  - Missing or misspelled 'apiVersion'
  - Missing or misspelled 'kind'
  - Invalid YAML syntax

For supported resource kinds, see: https://project-planton.org/docs/resources`
}

// InvalidProviderConfigGuidance returns a helpful message when provider config is invalid.
func InvalidProviderConfigGuidance(result *DetectionResult, parseErr error) string {
	var sb strings.Builder

	sb.WriteString(fmt.Sprintf("The provider config file could not be parsed as %s credentials.\n\n",
		ProviderDisplayName(result.Provider)))

	sb.WriteString("Parse error: " + parseErr.Error() + "\n\n")

	sb.WriteString(fmt.Sprintf("Expected format for %s provider config:\n\n",
		ProviderDisplayName(result.Provider)))

	// Indent the example
	example := ProviderConfigExample(result.Provider)
	for _, line := range strings.Split(example, "\n") {
		sb.WriteString("  " + line + "\n")
	}

	return sb.String()
}
