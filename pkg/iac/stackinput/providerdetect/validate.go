package providerdetect

import (
	"os"

	"github.com/pkg/errors"
	"github.com/plantonhq/project-planton/apis/org/project_planton/shared/cloudresourcekind"
	"github.com/plantonhq/project-planton/pkg/protobufyaml"
	"google.golang.org/protobuf/proto"

	atlasprovider "github.com/plantonhq/project-planton/apis/org/project_planton/provider/atlas"
	auth0provider "github.com/plantonhq/project-planton/apis/org/project_planton/provider/auth0"
	awsprovider "github.com/plantonhq/project-planton/apis/org/project_planton/provider/aws"
	azureprovider "github.com/plantonhq/project-planton/apis/org/project_planton/provider/azure"
	cloudflareprovider "github.com/plantonhq/project-planton/apis/org/project_planton/provider/cloudflare"
	confluentprovider "github.com/plantonhq/project-planton/apis/org/project_planton/provider/confluent"
	gcpprovider "github.com/plantonhq/project-planton/apis/org/project_planton/provider/gcp"
	kubernetesprovider "github.com/plantonhq/project-planton/apis/org/project_planton/provider/kubernetes"
	openfgaprovider "github.com/plantonhq/project-planton/apis/org/project_planton/provider/openfga"
	snowflakeprovider "github.com/plantonhq/project-planton/apis/org/project_planton/provider/snowflake"
)

// ValidateProviderConfig validates that the provider config file can be loaded
// as the expected provider type.
func ValidateProviderConfig(providerConfigPath string, provider cloudresourcekind.CloudResourceProvider) error {
	// Read the provider config file
	configBytes, err := os.ReadFile(providerConfigPath)
	if err != nil {
		return errors.Wrapf(err, "failed to read provider config file %s", providerConfigPath)
	}

	// Get the proto message for this provider
	protoMsg, err := getProviderConfigProto(provider)
	if err != nil {
		return err
	}

	// Try to load the config into the proto message
	if err := protobufyaml.LoadYamlBytes(configBytes, protoMsg); err != nil {
		return errors.Wrapf(err, "failed to parse provider config as %s config", ProviderDisplayName(provider))
	}

	return nil
}

// getProviderConfigProto returns a new proto message for the given provider.
func getProviderConfigProto(provider cloudresourcekind.CloudResourceProvider) (proto.Message, error) {
	switch provider {
	case cloudresourcekind.CloudResourceProvider_atlas:
		return new(atlasprovider.AtlasProviderConfig), nil
	case cloudresourcekind.CloudResourceProvider_auth0:
		return new(auth0provider.Auth0ProviderConfig), nil
	case cloudresourcekind.CloudResourceProvider_aws:
		return new(awsprovider.AwsProviderConfig), nil
	case cloudresourcekind.CloudResourceProvider_azure:
		return new(azureprovider.AzureProviderConfig), nil
	case cloudresourcekind.CloudResourceProvider_cloudflare:
		return new(cloudflareprovider.CloudflareProviderConfig), nil
	case cloudresourcekind.CloudResourceProvider_confluent:
		return new(confluentprovider.ConfluentProviderConfig), nil
	case cloudresourcekind.CloudResourceProvider_gcp:
		return new(gcpprovider.GcpProviderConfig), nil
	case cloudresourcekind.CloudResourceProvider_kubernetes:
		return new(kubernetesprovider.KubernetesProviderConfig), nil
	case cloudresourcekind.CloudResourceProvider_open_fga:
		return new(openfgaprovider.OpenFgaProviderConfig), nil
	case cloudresourcekind.CloudResourceProvider_snowflake:
		return new(snowflakeprovider.SnowflakeProviderConfig), nil
	default:
		return nil, errors.Errorf("unsupported provider: %s", provider.String())
	}
}

// LoadProviderConfigBytes reads a provider config file and returns its contents.
func LoadProviderConfigBytes(providerConfigPath string) ([]byte, error) {
	configBytes, err := os.ReadFile(providerConfigPath)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to read provider config file %s", providerConfigPath)
	}
	return configBytes, nil
}
