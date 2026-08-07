package verify

import (
	"context"

	"github.com/pkg/errors"
	"google.golang.org/api/googleapi"
)

// workloadIdentityPoolVerifier probes a workload identity pool by its full
// resource name (projects/<number>/locations/global/workloadIdentityPools/<id>)
// via the IAM API.
type workloadIdentityPoolVerifier struct{}

func (v *workloadIdentityPoolVerifier) IDOutputKey() string { return "name" }

func (v *workloadIdentityPoolVerifier) VerifyExists(ctx context.Context, svc *Services, outputs map[string]string) error {
	pool, err := svc.Iam.Projects.Locations.WorkloadIdentityPools.Get(outputs["name"]).Context(ctx).Do()
	if err != nil {
		return errors.Wrapf(err, "workload identity pool %s not found after deploy", outputs["name"])
	}
	if pool.State != "ACTIVE" {
		return errors.Errorf("workload identity pool %s exists but is %s (want ACTIVE) after deploy", outputs["name"], pool.State)
	}
	return nil
}

func (v *workloadIdentityPoolVerifier) VerifyAbsent(ctx context.Context, svc *Services, outputs map[string]string) error {
	pool, err := svc.Iam.Projects.Locations.WorkloadIdentityPools.Get(outputs["name"]).Context(ctx).Do()
	if err != nil {
		var apiErr *googleapi.Error
		if errors.As(err, &apiErr) && apiErr.Code == 404 {
			return nil
		}
		return errors.Wrapf(err, "unexpected error probing workload identity pool %s after destroy", outputs["name"])
	}
	// GCP soft-deletes pools: for ~30 days the pool remains readable in state
	// DELETED while rejecting token exchanges. That IS the destroyed state —
	// only an ACTIVE pool means the destroy failed.
	if pool.State == "ACTIVE" {
		return errors.Errorf("workload identity pool %s still ACTIVE (not even soft-deleted) after destroy", outputs["name"])
	}
	return nil
}
