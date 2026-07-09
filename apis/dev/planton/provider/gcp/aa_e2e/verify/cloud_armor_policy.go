package verify

import (
	"context"

	"github.com/pkg/errors"
	"google.golang.org/api/googleapi"
)

// cloudArmorPolicyVerifier probes a Cloud Armor security policy via the
// compute API. Posture assertions confirm the default rule invariant every
// policy carries (a rule at priority 2147483647) and that the exported
// self_link matches the live resource — the value every backend-service and
// backend-bucket FK consumes.
type cloudArmorPolicyVerifier struct{}

const cloudArmorDefaultRulePriority = int64(2147483647)

func (v *cloudArmorPolicyVerifier) IDOutputKey() string { return "policy_name" }

func (v *cloudArmorPolicyVerifier) VerifyExists(ctx context.Context, svc *Services, outputs map[string]string) error {
	policyName := outputs["policy_name"]
	if policyName == "" {
		return errors.New("policy_name output missing after deploy")
	}

	policy, err := svc.Compute.SecurityPolicies.Get(svc.Project, policyName).Context(ctx).Do()
	if err != nil {
		return errors.Wrapf(err, "cloud armor policy %s not found after deploy", policyName)
	}

	if wantSelfLink := outputs["policy_self_link"]; wantSelfLink != "" && policy.SelfLink != wantSelfLink {
		return errors.Errorf("cloud armor policy %s self_link mismatch: output %q, live %q", policyName, wantSelfLink, policy.SelfLink)
	}

	// Every Cloud Armor policy carries the default rule; its absence would
	// mean the rule set deployed incompletely.
	hasDefaultRule := false
	for _, rule := range policy.Rules {
		if rule.Priority == cloudArmorDefaultRulePriority {
			hasDefaultRule = true
			break
		}
	}
	if !hasDefaultRule {
		return errors.Errorf("cloud armor policy %s has no default rule at priority %d after deploy", policyName, cloudArmorDefaultRulePriority)
	}
	return nil
}

func (v *cloudArmorPolicyVerifier) VerifyAbsent(ctx context.Context, svc *Services, outputs map[string]string) error {
	policyName := outputs["policy_name"]
	if policyName == "" {
		return nil
	}

	_, err := svc.Compute.SecurityPolicies.Get(svc.Project, policyName).Context(ctx).Do()
	if err != nil {
		var apiErr *googleapi.Error
		if errors.As(err, &apiErr) && apiErr.Code == 404 {
			return nil
		}
		return errors.Wrapf(err, "unexpected error probing cloud armor policy %s after destroy", policyName)
	}
	return errors.Errorf("cloud armor policy %s still exists after destroy", policyName)
}
