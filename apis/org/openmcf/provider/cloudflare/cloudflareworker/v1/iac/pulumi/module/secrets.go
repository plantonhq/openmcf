package module

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"

	"github.com/pkg/errors"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

// uploadWorkerSecrets uploads secrets to Cloudflare Workers Secrets API.
// Secrets are encrypted at rest and separate from the worker version.
// This must be called AFTER the worker script exists.
func uploadWorkerSecrets(
	ctx *pulumi.Context,
	locals *Locals,
	scriptName string,
	secrets map[string]string,
) error {

	if len(secrets) == 0 {
		return nil // No secrets to upload
	}

	accountId := locals.CloudflareWorker.Spec.AccountId

	// Get API token from provider config or environment variable
	var apiToken string
	if locals.CloudflareProviderConfig != nil && locals.CloudflareProviderConfig.ApiToken != "" {
		apiToken = locals.CloudflareProviderConfig.ApiToken
	} else {
		// Fall back to CLOUDFLARE_API_TOKEN environment variable
		apiToken = os.Getenv("CLOUDFLARE_API_TOKEN")
		if apiToken == "" {
			return errors.New("Cloudflare API token not found in provider config or CLOUDFLARE_API_TOKEN environment variable")
		}
	}

	// Upload each secret individually (Cloudflare API requirement)
	for key, value := range secrets {
		if err := uploadSingleSecret(accountId, scriptName, key, value, apiToken); err != nil {
			return errors.Wrapf(err, "failed to upload secret '%s'", key)
		}
	}

	return nil
}

// uploadSingleSecret uploads a single secret via Cloudflare API.
// https://developers.cloudflare.com/api/operations/worker-script-upload-worker-secret
func uploadSingleSecret(accountId, scriptName, secretName, secretValue, apiToken string) error {
	url := fmt.Sprintf(
		"https://api.cloudflare.com/client/v4/accounts/%s/workers/scripts/%s/secrets",
		accountId, scriptName)

	payload := map[string]interface{}{
		"name": secretName,
		"text": secretValue,
		"type": "secret_text",
	}

	jsonPayload, err := json.Marshal(payload)
	if err != nil {
		return errors.Wrap(err, "failed to marshal secret payload")
	}

	req, err := http.NewRequest("PUT", url, bytes.NewBuffer(jsonPayload))
	if err != nil {
		return errors.Wrap(err, "failed to create request")
	}

	req.Header.Set("Authorization", "Bearer "+apiToken)
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{Timeout: 30 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return errors.Wrap(err, "failed to send request")
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)

	if resp.StatusCode != http.StatusOK {
		return errors.Errorf("Cloudflare API returned status %d: %s", resp.StatusCode, string(body))
	}

	return nil
}
