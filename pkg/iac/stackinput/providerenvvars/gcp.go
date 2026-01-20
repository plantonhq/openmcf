package providerenvvars

import (
	"encoding/base64"

	"github.com/pkg/errors"
	gcpprovider "github.com/plantonhq/project-planton/apis/org/project_planton/provider/gcp"
)

// loadGcpEnvVars loads GCP provider config and returns environment variables.
func loadGcpEnvVars(providerConfigYaml []byte) (map[string]string, error) {
	config := new(gcpprovider.GcpProviderConfig)
	if err := loadProviderConfigProto(providerConfigYaml, config); err != nil {
		return nil, errors.Wrap(err, "failed to load GCP provider config")
	}

	serviceAccountKey, err := base64.StdEncoding.DecodeString(config.ServiceAccountKeyBase64)
	if err != nil {
		return nil, errors.Wrap(err, "failed to decode service account key from base64")
	}

	envVars := map[string]string{
		"GOOGLE_CREDENTIALS": string(serviceAccountKey),
	}

	return envVars, nil
}
