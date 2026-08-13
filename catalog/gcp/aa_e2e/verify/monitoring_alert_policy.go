package verify

import (
	"context"

	"github.com/pkg/errors"
	"google.golang.org/api/googleapi"
)

// monitoringAlertPolicyVerifier probes a Cloud Monitoring alert policy by
// its server-assigned resource name and confirms the policy carries at
// least one condition — a policy without conditions cannot exist through
// this kind's spec, so an empty condition list would mean the expand
// dropped the family.
type monitoringAlertPolicyVerifier struct{}

func (v *monitoringAlertPolicyVerifier) IDOutputKey() string { return "policy_name" }

func (v *monitoringAlertPolicyVerifier) VerifyExists(ctx context.Context, svc *Services, outputs map[string]string) error {
	name := outputs["policy_name"]
	policy, err := svc.Monitoring.Projects.AlertPolicies.Get(name).Context(ctx).Do()
	if err != nil {
		return errors.Wrapf(err, "alert policy %s not found after deploy", name)
	}
	if len(policy.Conditions) == 0 {
		return errors.Errorf("alert policy %s reports no conditions", name)
	}
	return nil
}

func (v *monitoringAlertPolicyVerifier) VerifyAbsent(ctx context.Context, svc *Services, outputs map[string]string) error {
	name := outputs["policy_name"]
	_, err := svc.Monitoring.Projects.AlertPolicies.Get(name).Context(ctx).Do()
	if err != nil {
		var apiErr *googleapi.Error
		if errors.As(err, &apiErr) && apiErr.Code == 404 {
			return nil
		}
		return errors.Wrapf(err, "unexpected error probing alert policy %s after destroy", name)
	}
	return errors.Errorf("alert policy %s still exists after destroy", name)
}
