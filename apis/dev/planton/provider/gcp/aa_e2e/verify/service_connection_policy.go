package verify

import (
	"context"

	"github.com/pkg/errors"
	"google.golang.org/api/googleapi"
)

// serviceConnectionPolicyVerifier probes a service connection policy via the
// Network Connectivity API. Posture assertions confirm the policy's
// infrastructure output matches the live object (PSC) and that the live
// policy carries the service class the automation matches producers on.
type serviceConnectionPolicyVerifier struct{}

func (v *serviceConnectionPolicyVerifier) IDOutputKey() string { return "policy_id" }

func (v *serviceConnectionPolicyVerifier) VerifyExists(ctx context.Context, svc *Services, outputs map[string]string) error {
	policyId := outputs["policy_id"]
	if policyId == "" {
		return errors.New("policy_id output missing after deploy")
	}

	policy, err := svc.NetworkConnectivity.Projects.Locations.ServiceConnectionPolicies.Get(policyId).Context(ctx).Do()
	if err != nil {
		return errors.Wrapf(err, "service connection policy %s not found after deploy", policyId)
	}
	if policy.ServiceClass == "" {
		return errors.Errorf("service connection policy %s has no service class — the automation cannot match producers", policyId)
	}
	if infra := outputs["infrastructure"]; infra != "" && policy.Infrastructure != infra {
		return errors.Errorf("service connection policy %s infrastructure mismatch: output %q, live %q",
			policyId, infra, policy.Infrastructure)
	}
	return nil
}

func (v *serviceConnectionPolicyVerifier) VerifyAbsent(ctx context.Context, svc *Services, outputs map[string]string) error {
	policyId := outputs["policy_id"]
	if policyId == "" {
		return nil
	}

	_, err := svc.NetworkConnectivity.Projects.Locations.ServiceConnectionPolicies.Get(policyId).Context(ctx).Do()
	if err != nil {
		var apiErr *googleapi.Error
		if errors.As(err, &apiErr) && apiErr.Code == 404 {
			return nil
		}
		return errors.Wrapf(err, "unexpected error probing service connection policy %s after destroy", policyId)
	}
	return errors.Errorf("service connection policy %s still exists after destroy", policyId)
}
