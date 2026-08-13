package verify

import (
	"context"

	"github.com/pkg/errors"
)

// identityPlatformConfigVerifier probes a project's Identity Platform
// configuration through the Identity Toolkit admin API
// (projects/{project}/config).
//
// The config resource is UNDELETABLE by API design: destroy abandons the
// configuration in place (GCP has no de-initialize). VERIFY-CLN therefore
// asserts the ABANDON contract — the configuration must still be readable
// after destroy — rather than absence.
type identityPlatformConfigVerifier struct{}

func (v *identityPlatformConfigVerifier) IDOutputKey() string { return "config_name" }

func (v *identityPlatformConfigVerifier) VerifyExists(ctx context.Context, svc *Services, outputs map[string]string) error {
	name := outputs["config_name"]
	if _, err := svc.IdentityToolkit.Projects.GetConfig(name).Context(ctx).Do(); err != nil {
		return errors.Wrapf(err, "identity platform config %s not found after deploy", name)
	}
	return nil
}

func (v *identityPlatformConfigVerifier) VerifyAbsent(ctx context.Context, svc *Services, outputs map[string]string) error {
	name := outputs["config_name"]
	// The abandon contract: the project keeps its (permanently initialized)
	// Identity Platform configuration after destroy — a readable config IS
	// the expected clean state for this kind.
	if _, err := svc.IdentityToolkit.Projects.GetConfig(name).Context(ctx).Do(); err != nil {
		return errors.Wrapf(err, "identity platform config %s unreadable after destroy — the abandon contract expects it to remain", name)
	}
	return nil
}
