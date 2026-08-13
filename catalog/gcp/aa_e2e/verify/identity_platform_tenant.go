package verify

import (
	"context"

	"github.com/pkg/errors"
	"google.golang.org/api/googleapi"
)

// identityPlatformTenantVerifier probes an Identity Platform tenant through
// the Identity Toolkit admin API (projects/{project}/tenants/{tenant_id}).
// Unlike the project config, tenants are fully deletable, so the verifier
// asserts real absence after destroy.
type identityPlatformTenantVerifier struct{}

func (v *identityPlatformTenantVerifier) IDOutputKey() string { return "tenant_name" }

func (v *identityPlatformTenantVerifier) VerifyExists(ctx context.Context, svc *Services, outputs map[string]string) error {
	name := outputs["tenant_name"]
	if _, err := svc.IdentityToolkit.Projects.Tenants.Get(name).Context(ctx).Do(); err != nil {
		return errors.Wrapf(err, "identity platform tenant %s not found after deploy", name)
	}
	return nil
}

func (v *identityPlatformTenantVerifier) VerifyAbsent(ctx context.Context, svc *Services, outputs map[string]string) error {
	name := outputs["tenant_name"]
	_, err := svc.IdentityToolkit.Projects.Tenants.Get(name).Context(ctx).Do()
	if err != nil {
		var apiErr *googleapi.Error
		if errors.As(err, &apiErr) && apiErr.Code == 404 {
			return nil
		}
		return errors.Wrapf(err, "unexpected error probing tenant %s after destroy", name)
	}
	return errors.Errorf("identity platform tenant %s still exists after destroy", name)
}
