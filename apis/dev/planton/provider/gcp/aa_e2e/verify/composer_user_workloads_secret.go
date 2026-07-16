package verify

import (
	"context"

	"github.com/pkg/errors"
	"google.golang.org/api/googleapi"
)

// composerUserWorkloadsSecretVerifier probes a Composer user workloads
// Secret via the Composer API. The name output is the fully qualified
// resource name
// (projects/{p}/locations/{r}/environments/{e}/userWorkloadsSecrets/{n}).
// The get call returns metadata only (never the secret data), so the
// probe confirms existence without ever touching the material.
type composerUserWorkloadsSecretVerifier struct{}

func (v *composerUserWorkloadsSecretVerifier) IDOutputKey() string { return "name" }

func (v *composerUserWorkloadsSecretVerifier) VerifyExists(ctx context.Context, svc *Services, outputs map[string]string) error {
	name := outputs["name"]
	if name == "" {
		return errors.New("name output missing after deploy")
	}

	_, err := svc.Composer.Projects.Locations.Environments.UserWorkloadsSecrets.Get(name).Context(ctx).Do()
	if err != nil {
		return errors.Wrapf(err, "composer user workloads secret %s not found after deploy", name)
	}
	return nil
}

func (v *composerUserWorkloadsSecretVerifier) VerifyAbsent(ctx context.Context, svc *Services, outputs map[string]string) error {
	name := outputs["name"]
	if name == "" {
		return nil
	}

	_, err := svc.Composer.Projects.Locations.Environments.UserWorkloadsSecrets.Get(name).Context(ctx).Do()
	if err == nil {
		return errors.Errorf("composer user workloads secret %s still exists after destroy", name)
	}
	var apiErr *googleapi.Error
	if errors.As(err, &apiErr) && apiErr.Code == 404 {
		return nil
	}
	// When the whole chain (environment included) was destroyed, the API
	// answers for the missing parent instead of the secret.
	if apiErr != nil && apiErr.Code == 400 {
		return nil
	}
	return errors.Wrapf(err, "unexpected error probing composer user workloads secret %s after destroy", name)
}
