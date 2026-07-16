package verify

import (
	"context"

	"github.com/pkg/errors"
	"google.golang.org/api/googleapi"
)

// workloadIdentityPoolProviderVerifier probes a pool provider by its full
// resource name (projects/<number>/locations/global/workloadIdentityPools/
// <pool>/providers/<provider>) via the IAM API.
type workloadIdentityPoolProviderVerifier struct{}

func (v *workloadIdentityPoolProviderVerifier) IDOutputKey() string { return "name" }

func (v *workloadIdentityPoolProviderVerifier) VerifyExists(ctx context.Context, svc *Services, outputs map[string]string) error {
	provider, err := svc.Iam.Projects.Locations.WorkloadIdentityPools.Providers.Get(outputs["name"]).Context(ctx).Do()
	if err != nil {
		return errors.Wrapf(err, "workload identity pool provider %s not found after deploy", outputs["name"])
	}
	if provider.State != "ACTIVE" {
		return errors.Errorf("workload identity pool provider %s exists but is %s (want ACTIVE) after deploy", outputs["name"], provider.State)
	}
	return nil
}

func (v *workloadIdentityPoolProviderVerifier) VerifyAbsent(ctx context.Context, svc *Services, outputs map[string]string) error {
	provider, err := svc.Iam.Projects.Locations.WorkloadIdentityPools.Providers.Get(outputs["name"]).Context(ctx).Do()
	if err != nil {
		var apiErr *googleapi.Error
		if errors.As(err, &apiErr) && apiErr.Code == 404 {
			return nil
		}
		return errors.Wrapf(err, "unexpected error probing workload identity pool provider %s after destroy", outputs["name"])
	}
	// Soft-deleted providers remain readable in state DELETED for ~30 days —
	// that IS the destroyed state; only ACTIVE means the destroy failed.
	// (The API also 404s provider gets once the PARENT pool is soft-deleted,
	// which the error path above already treats as absent.)
	if provider.State == "ACTIVE" {
		return errors.Errorf("workload identity pool provider %s still ACTIVE (not even soft-deleted) after destroy", outputs["name"])
	}
	return nil
}
