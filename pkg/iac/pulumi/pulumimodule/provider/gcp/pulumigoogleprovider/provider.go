package pulumigoogleprovider

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"reflect"

	gcpprovider "github.com/plantonhq/openmcf/apis/org/openmcf/provider/gcp"
	"github.com/plantonhq/openmcf/pkg/iac/pulumi/pulumimodule/pulumi/pulumioutput"

	"github.com/pkg/errors"
	"github.com/pulumi/pulumi-gcp/sdk/v9/go/gcp"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

func Get(ctx *pulumi.Context, gcpProviderConfig *gcpprovider.GcpProviderConfig,
	nameSuffixes ...string) (*gcp.Provider, error) {
	gcpProviderArgs := &gcp.ProviderArgs{}

	//if stack-input contains gcp-credentials, provider will be created with those credentials
	if gcpProviderConfig != nil {
		serviceAccountKeyBytes, err := base64.StdEncoding.DecodeString(gcpProviderConfig.ServiceAccountKeyBase64)
		if err != nil {
			return nil, errors.Wrap(err, "failed to decode base64 encoded"+
				" google service account credential")
		}
		// Convert bytes to string - the GCP provider expects the JSON content as a string
		serviceAccountKeyJSON := string(serviceAccountKeyBytes)

		// Validate that the decoded content is valid JSON
		var serviceAccountKeyMap map[string]interface{}
		if err := json.Unmarshal(serviceAccountKeyBytes, &serviceAccountKeyMap); err != nil {
			return nil, errors.Wrap(err, "failed to parse service account key JSON. "+
				"Ensure the base64-encoded value contains valid JSON with fields: type, project_id, private_key_id, private_key, client_email, client_id, auth_uri, token_uri")
		}

		// Validate required fields are present
		requiredFields := []string{"type", "project_id", "private_key", "client_email"}
		for _, field := range requiredFields {
			if _, ok := serviceAccountKeyMap[field]; !ok {
				return nil, errors.Errorf("service account key JSON is missing required field: %s", field)
			}
		}

		// Validate private_key format - it should be a PEM-encoded key
		privateKey, ok := serviceAccountKeyMap["private_key"].(string)
		if !ok {
			return nil, errors.New("service account key 'private_key' field must be a string")
		}
		// Check if private key looks like PEM format (starts with -----BEGIN)
		if len(privateKey) > 0 && !(privateKey[:11] == "-----BEGIN " || privateKey[:15] == "-----BEGIN RSA ") {
			return nil, errors.New("service account key 'private_key' field must be a PEM-encoded key " +
				"(starting with '-----BEGIN PRIVATE KEY-----' or '-----BEGIN RSA PRIVATE KEY-----'). " +
				"Ensure you're using a JSON key file from GCP, not a P12/PKCS12 key")
		}

		gcpProviderArgs = &gcp.ProviderArgs{Credentials: pulumi.String(serviceAccountKeyJSON)}
	}

	googleProvider, err := gcp.NewProvider(ctx, ProviderResourceName(nameSuffixes), gcpProviderArgs)
	if err != nil {
		return nil, errors.Wrap(err, "failed to create google provider")
	}

	return googleProvider, nil
}

func ProviderResourceName(suffixes []string) string {
	name := "google"
	for _, s := range suffixes {
		name = fmt.Sprintf("%s-%s", name, s)
	}
	return name
}

func PulumiOutputName(r interface{}, name string, suffixes ...string) string {
	outputName := fmt.Sprintf("gcp_%s", pulumioutput.Name(reflect.TypeOf(r), name))
	for _, s := range suffixes {
		outputName = fmt.Sprintf("%s_%s", outputName, s)
	}
	return outputName
}
