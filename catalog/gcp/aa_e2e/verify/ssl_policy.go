package verify

import (
	"context"

	"github.com/pkg/errors"
	"google.golang.org/api/googleapi"
)

// sslPolicyVerifier probes a Compute Engine SSL policy by name and confirms
// the hardening posture landed: the policy GCP stores must carry the
// non-empty enabled-features set the deploy reported. One kind maps to two
// API collections (global and regional SSL policies), so the verifier
// routes on the region output rather than requiring the caller to know the
// scope.
type sslPolicyVerifier struct{}

func (v *sslPolicyVerifier) IDOutputKey() string { return "self_link" }

func (v *sslPolicyVerifier) VerifyExists(ctx context.Context, svc *Services, outputs map[string]string) error {
	name := outputs["ssl_policy_name"]
	region := outputs["region"]

	if region != "" {
		policy, err := svc.Compute.RegionSslPolicies.Get(svc.Project, region, name).Context(ctx).Do()
		if err != nil {
			return errors.Wrapf(err, "regional ssl policy %s/%s not found after deploy", region, name)
		}
		if len(policy.EnabledFeatures) == 0 {
			return errors.Errorf("regional ssl policy %s/%s reports no enabled cipher features", region, name)
		}
		return nil
	}
	policy, err := svc.Compute.SslPolicies.Get(svc.Project, name).Context(ctx).Do()
	if err != nil {
		return errors.Wrapf(err, "global ssl policy %s not found after deploy", name)
	}
	if len(policy.EnabledFeatures) == 0 {
		return errors.Errorf("global ssl policy %s reports no enabled cipher features", name)
	}
	return nil
}

func (v *sslPolicyVerifier) VerifyAbsent(ctx context.Context, svc *Services, outputs map[string]string) error {
	name := outputs["ssl_policy_name"]
	region := outputs["region"]

	var err error
	if region != "" {
		_, err = svc.Compute.RegionSslPolicies.Get(svc.Project, region, name).Context(ctx).Do()
	} else {
		_, err = svc.Compute.SslPolicies.Get(svc.Project, name).Context(ctx).Do()
	}
	if err != nil {
		var apiErr *googleapi.Error
		if errors.As(err, &apiErr) && apiErr.Code == 404 {
			return nil
		}
		return errors.Wrapf(err, "unexpected error probing ssl policy %s after destroy", name)
	}
	return errors.Errorf("ssl policy %s still exists after destroy", name)
}
