package providerenvvars

import (
	"github.com/pkg/errors"
	awsprovider "github.com/plantonhq/planton/catalog/aws"
	"github.com/plantonhq/planton/pkg/crkreflect"
	"github.com/plantonhq/planton/shared/cloudresourcekind"
	"gopkg.in/yaml.v3"
)

// ExtractAwsProviderConfig decodes the stack input's provider_config into the typed
// AwsProviderConfig when (and only when) the target resource belongs to the AWS
// provider. Returns (nil, nil) for non-AWS targets and for stack inputs carrying no
// provider config -- both are ordinary, not errors.
//
// This is the single stack-input -> typed-AWS-config decode point shared by the
// env-var loader above and the tofu provider-override writer (pkg/iac/tofu/tfoverride);
// consumers must never re-implement the target/kind/provider dispatch, so the two
// paths cannot drift on how a provider config is recognized.
func ExtractAwsProviderConfig(stackInputYaml string) (*awsprovider.AwsProviderConfig, error) {
	stackInputMap := map[string]interface{}{}
	if err := yaml.Unmarshal([]byte(stackInputYaml), &stackInputMap); err != nil {
		return nil, errors.Wrap(err, "failed to unmarshal stack input yaml")
	}

	targetYaml, err := extractTargetYaml(stackInputMap)
	if err != nil {
		return nil, errors.Wrap(err, "failed to extract target from stack input")
	}

	kind, err := crkreflect.ExtractKindFromYaml(targetYaml)
	if err != nil {
		return nil, errors.Wrap(err, "failed to extract cloud resource kind from target")
	}

	if crkreflect.GetProvider(kind) != cloudresourcekind.CloudResourceProvider_aws {
		return nil, nil
	}

	providerConfigYaml, hasProviderConfig := extractProviderConfigYaml(stackInputMap)
	if !hasProviderConfig {
		return nil, nil
	}

	config := new(awsprovider.AwsProviderConfig)
	if err := loadProviderConfigProto(providerConfigYaml, config); err != nil {
		return nil, errors.Wrap(err, "failed to load AWS provider config")
	}
	return config, nil
}
