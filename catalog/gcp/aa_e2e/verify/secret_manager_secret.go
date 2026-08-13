package verify

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	"github.com/pkg/errors"
	"google.golang.org/api/googleapi"
)

// secretManagerSecretVerifier probes a Secret Manager secret by its full
// resource name. GLOBAL secrets go through the typed client; REGIONAL
// secrets (name carries /locations/) are served ONLY by the location's
// regional endpoint (secretmanager.{location}.rep.googleapis.com), which
// the typed client does not route to — those probe through the
// ADC-authenticated REST client, the established pattern for endpoints
// outside the pinned typed surface. When the deploy seeded a version, its
// existence is asserted too — the "readable secret from one manifest"
// contract.
type secretManagerSecretVerifier struct{}

func (v *secretManagerSecretVerifier) IDOutputKey() string { return "secret_name" }

// regionalLocation extracts the location segment from a regional secret
// name (projects/{p}/locations/{l}/secrets/{id}); empty for global names.
func regionalLocation(name string) string {
	parts := strings.Split(name, "/")
	for i := 0; i < len(parts)-1; i++ {
		if parts[i] == "locations" {
			return parts[i+1]
		}
	}
	return ""
}

// probeRegional issues an authenticated GET against the regional Secret
// Manager endpoint and returns the HTTP status code.
func probeRegional(ctx context.Context, svc *Services, location, name string) (int, error) {
	url := fmt.Sprintf("https://secretmanager.%s.rep.googleapis.com/v1/%s", location, name)
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return 0, errors.Wrap(err, "failed to build regional secretmanager GET request")
	}
	resp, err := svc.RestClient.Do(req)
	if err != nil {
		return 0, errors.Wrap(err, "regional secretmanager GET request failed")
	}
	defer resp.Body.Close()
	return resp.StatusCode, nil
}

func (v *secretManagerSecretVerifier) VerifyExists(ctx context.Context, svc *Services, outputs map[string]string) error {
	name := outputs["secret_name"]
	versionName := outputs["latest_version_name"]

	if location := regionalLocation(name); location != "" {
		status, err := probeRegional(ctx, svc, location, name)
		if err != nil {
			return err
		}
		if status != http.StatusOK {
			return errors.Errorf("regional secret %s probe returned %d after deploy", name, status)
		}
		if versionName != "" {
			status, err := probeRegional(ctx, svc, location, versionName)
			if err != nil {
				return err
			}
			if status != http.StatusOK {
				return errors.Errorf("regional secret version %s probe returned %d after deploy", versionName, status)
			}
		}
		return nil
	}

	if _, err := svc.SecretManager.Projects.Secrets.Get(name).Context(ctx).Do(); err != nil {
		return errors.Wrapf(err, "secret %s not found after deploy", name)
	}
	if versionName != "" {
		if _, err := svc.SecretManager.Projects.Secrets.Versions.Get(versionName).Context(ctx).Do(); err != nil {
			return errors.Wrapf(err, "secret version %s not found after deploy", versionName)
		}
	}
	return nil
}

func (v *secretManagerSecretVerifier) VerifyAbsent(ctx context.Context, svc *Services, outputs map[string]string) error {
	name := outputs["secret_name"]

	if location := regionalLocation(name); location != "" {
		status, err := probeRegional(ctx, svc, location, name)
		if err != nil {
			return err
		}
		if status == http.StatusNotFound {
			return nil
		}
		return errors.Errorf("regional secret %s still answers %d after destroy", name, status)
	}

	_, err := svc.SecretManager.Projects.Secrets.Get(name).Context(ctx).Do()
	if err != nil {
		var apiErr *googleapi.Error
		if errors.As(err, &apiErr) && apiErr.Code == 404 {
			return nil
		}
		return errors.Wrapf(err, "unexpected error probing secret %s after destroy", name)
	}
	return errors.Errorf("secret %s still exists after destroy", name)
}
