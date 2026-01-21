package providerdetect

import (
	"fmt"
	"strings"

	"github.com/plantonhq/project-planton/apis/org/project_planton/provider/atlas"
	"github.com/plantonhq/project-planton/apis/org/project_planton/provider/auth0"
	"github.com/plantonhq/project-planton/apis/org/project_planton/provider/aws"
	"github.com/plantonhq/project-planton/apis/org/project_planton/provider/azure"
	"github.com/plantonhq/project-planton/apis/org/project_planton/provider/civo"
	"github.com/plantonhq/project-planton/apis/org/project_planton/provider/cloudflare"
	"github.com/plantonhq/project-planton/apis/org/project_planton/provider/confluent"
	"github.com/plantonhq/project-planton/apis/org/project_planton/provider/digitalocean"
	"github.com/plantonhq/project-planton/apis/org/project_planton/provider/gcp"
	"github.com/plantonhq/project-planton/apis/org/project_planton/provider/kubernetes"
	"github.com/plantonhq/project-planton/apis/org/project_planton/provider/openfga"
	"github.com/plantonhq/project-planton/apis/org/project_planton/provider/snowflake"
	"github.com/plantonhq/project-planton/apis/org/project_planton/shared/cloudresourcekind"
)

// ProviderConfigExample returns an example YAML configuration for the given provider.
func ProviderConfigExample(provider cloudresourcekind.CloudResourceProvider) string {
	switch provider {
	case cloudresourcekind.CloudResourceProvider_atlas:
		return atlas.ConfigFileExample
	case cloudresourcekind.CloudResourceProvider_auth0:
		return auth0.ConfigFileExample
	case cloudresourcekind.CloudResourceProvider_aws:
		return aws.ConfigFileExample
	case cloudresourcekind.CloudResourceProvider_azure:
		return azure.ConfigFileExample
	case cloudresourcekind.CloudResourceProvider_civo:
		return civo.ConfigFileExample
	case cloudresourcekind.CloudResourceProvider_cloudflare:
		return cloudflare.ConfigFileExample
	case cloudresourcekind.CloudResourceProvider_confluent:
		return confluent.ConfigFileExample
	case cloudresourcekind.CloudResourceProvider_digital_ocean:
		return digitalocean.ConfigFileExample
	case cloudresourcekind.CloudResourceProvider_gcp:
		return gcp.ConfigFileExample
	case cloudresourcekind.CloudResourceProvider_kubernetes:
		return kubernetes.ConfigFileExample
	case cloudresourcekind.CloudResourceProvider_open_fga:
		return openfga.ConfigFileExample
	case cloudresourcekind.CloudResourceProvider_snowflake:
		return snowflake.ConfigFileExample
	default:
		return "# Provider config format not available"
	}
}

// ProviderConfigFilename returns the suggested filename for the provider config.
func ProviderConfigFilename(provider cloudresourcekind.CloudResourceProvider) string {
	switch provider {
	case cloudresourcekind.CloudResourceProvider_atlas:
		return atlas.ConfigFileName
	case cloudresourcekind.CloudResourceProvider_auth0:
		return auth0.ConfigFileName
	case cloudresourcekind.CloudResourceProvider_aws:
		return aws.ConfigFileName
	case cloudresourcekind.CloudResourceProvider_azure:
		return azure.ConfigFileName
	case cloudresourcekind.CloudResourceProvider_civo:
		return civo.ConfigFileName
	case cloudresourcekind.CloudResourceProvider_cloudflare:
		return cloudflare.ConfigFileName
	case cloudresourcekind.CloudResourceProvider_confluent:
		return confluent.ConfigFileName
	case cloudresourcekind.CloudResourceProvider_digital_ocean:
		return digitalocean.ConfigFileName
	case cloudresourcekind.CloudResourceProvider_gcp:
		return gcp.ConfigFileName
	case cloudresourcekind.CloudResourceProvider_kubernetes:
		return kubernetes.ConfigFileName
	case cloudresourcekind.CloudResourceProvider_open_fga:
		return openfga.ConfigFileName
	case cloudresourcekind.CloudResourceProvider_snowflake:
		return snowflake.ConfigFileName
	default:
		return "provider-config.yaml"
	}
}

// ProviderEnvironmentVariablesHelp returns the environment variable export commands for the provider.
func ProviderEnvironmentVariablesHelp(provider cloudresourcekind.CloudResourceProvider) string {
	switch provider {
	case cloudresourcekind.CloudResourceProvider_atlas:
		return atlas.EnvironmentVariablesHelp
	case cloudresourcekind.CloudResourceProvider_auth0:
		return auth0.EnvironmentVariablesHelp
	case cloudresourcekind.CloudResourceProvider_aws:
		return aws.EnvironmentVariablesHelp
	case cloudresourcekind.CloudResourceProvider_azure:
		return azure.EnvironmentVariablesHelp
	case cloudresourcekind.CloudResourceProvider_civo:
		return civo.EnvironmentVariablesHelp
	case cloudresourcekind.CloudResourceProvider_cloudflare:
		return cloudflare.EnvironmentVariablesHelp
	case cloudresourcekind.CloudResourceProvider_confluent:
		return confluent.EnvironmentVariablesHelp
	case cloudresourcekind.CloudResourceProvider_digital_ocean:
		return digitalocean.EnvironmentVariablesHelp
	case cloudresourcekind.CloudResourceProvider_gcp:
		return gcp.EnvironmentVariablesHelp
	case cloudresourcekind.CloudResourceProvider_kubernetes:
		return kubernetes.EnvironmentVariablesHelp
	case cloudresourcekind.CloudResourceProvider_open_fga:
		return openfga.EnvironmentVariablesHelp
	case cloudresourcekind.CloudResourceProvider_snowflake:
		return snowflake.EnvironmentVariablesHelp
	default:
		return "# Environment variables not available for this provider"
	}
}

// ProviderDocsURL returns the documentation URL for the provider.
func ProviderDocsURL(provider cloudresourcekind.CloudResourceProvider) string {
	switch provider {
	case cloudresourcekind.CloudResourceProvider_atlas:
		return atlas.ProviderDocsURL
	case cloudresourcekind.CloudResourceProvider_auth0:
		return auth0.ProviderDocsURL
	case cloudresourcekind.CloudResourceProvider_aws:
		return aws.ProviderDocsURL
	case cloudresourcekind.CloudResourceProvider_azure:
		return azure.ProviderDocsURL
	case cloudresourcekind.CloudResourceProvider_civo:
		return civo.ProviderDocsURL
	case cloudresourcekind.CloudResourceProvider_cloudflare:
		return cloudflare.ProviderDocsURL
	case cloudresourcekind.CloudResourceProvider_confluent:
		return confluent.ProviderDocsURL
	case cloudresourcekind.CloudResourceProvider_digital_ocean:
		return digitalocean.ProviderDocsURL
	case cloudresourcekind.CloudResourceProvider_gcp:
		return gcp.ProviderDocsURL
	case cloudresourcekind.CloudResourceProvider_kubernetes:
		return kubernetes.ProviderDocsURL
	case cloudresourcekind.CloudResourceProvider_open_fga:
		return openfga.ProviderDocsURL
	case cloudresourcekind.CloudResourceProvider_snowflake:
		return snowflake.ProviderDocsURL
	default:
		return ""
	}
}

// MissingProviderConfigGuidance returns a helpful message when provider config is missing.
// It shows both options: environment variables (default) and explicit config file.
func MissingProviderConfigGuidance(result *DetectionResult) string {
	var sb strings.Builder

	sb.WriteString(fmt.Sprintf("The %s resource requires %s credentials.\n\n",
		result.KindName, ProviderDisplayName(result.Provider)))

	// Option 1: Environment variables (recommended for local development)
	sb.WriteString("Option 1: Set environment variables\n\n")
	envHelp := ProviderEnvironmentVariablesHelp(result.Provider)
	for _, line := range strings.Split(envHelp, "\n") {
		sb.WriteString("  " + line + "\n")
	}

	// Option 2: Provider config file
	sb.WriteString("\nOption 2: Create a provider config file\n\n")
	sb.WriteString(fmt.Sprintf("  Create '%s' with:\n\n",
		ProviderConfigFilename(result.Provider)))

	example := ProviderConfigExample(result.Provider)
	for _, line := range strings.Split(example, "\n") {
		sb.WriteString("    " + line + "\n")
	}

	sb.WriteString("\n  Then run:\n\n")
	sb.WriteString(fmt.Sprintf("    project-planton plan -f manifest.yaml -p %s\n",
		ProviderConfigFilename(result.Provider)))

	// Add documentation link if available
	docsURL := ProviderDocsURL(result.Provider)
	if docsURL != "" {
		sb.WriteString(fmt.Sprintf("\nFor more information: %s\n", docsURL))
	}

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

	example := ProviderConfigExample(result.Provider)
	for _, line := range strings.Split(example, "\n") {
		sb.WriteString("  " + line + "\n")
	}

	docsURL := ProviderDocsURL(result.Provider)
	if docsURL != "" {
		sb.WriteString(fmt.Sprintf("\nFor more information: %s\n", docsURL))
	}

	return sb.String()
}
