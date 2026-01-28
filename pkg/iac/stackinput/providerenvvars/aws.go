package providerenvvars

import (
	"github.com/pkg/errors"
	awsprovider "github.com/plantonhq/openmcf/apis/org/openmcf/provider/aws"
)

// loadAwsEnvVars loads AWS provider config and returns environment variables.
func loadAwsEnvVars(providerConfigYaml []byte) (map[string]string, error) {
	config := new(awsprovider.AwsProviderConfig)
	if err := loadProviderConfigProto(providerConfigYaml, config); err != nil {
		return nil, errors.Wrap(err, "failed to load AWS provider config")
	}

	envVars := map[string]string{
		"AWS_REGION":            config.GetRegion(),
		"AWS_ACCESS_KEY_ID":     config.AccessKeyId,
		"AWS_SECRET_ACCESS_KEY": config.SecretAccessKey,
	}

	return envVars, nil
}
