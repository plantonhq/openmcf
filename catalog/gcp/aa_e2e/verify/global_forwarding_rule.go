package verify

import (
	"context"

	"github.com/pkg/errors"
	"google.golang.org/api/googleapi"
)

// globalForwardingRuleVerifier probes a global forwarding rule by name and
// confirms the frontend actually formed: a VIP was bound and the target
// wiring exists — the rule is the node DNS points at, so both must hold.
type globalForwardingRuleVerifier struct{}

func (v *globalForwardingRuleVerifier) IDOutputKey() string { return "self_link" }

func (v *globalForwardingRuleVerifier) VerifyExists(ctx context.Context, svc *Services, outputs map[string]string) error {
	name := outputs["forwarding_rule_name"]
	rule, err := svc.Compute.GlobalForwardingRules.Get(svc.Project, name).Context(ctx).Do()
	if err != nil {
		return errors.Wrapf(err, "global forwarding rule %s not found after deploy", name)
	}
	if rule.IPAddress == "" {
		return errors.Errorf("global forwarding rule %s has no VIP bound", name)
	}
	if rule.Target == "" {
		return errors.Errorf("global forwarding rule %s has no target wired", name)
	}
	return nil
}

func (v *globalForwardingRuleVerifier) VerifyAbsent(ctx context.Context, svc *Services, outputs map[string]string) error {
	name := outputs["forwarding_rule_name"]
	_, err := svc.Compute.GlobalForwardingRules.Get(svc.Project, name).Context(ctx).Do()
	if err != nil {
		var apiErr *googleapi.Error
		if errors.As(err, &apiErr) && apiErr.Code == 404 {
			return nil
		}
		return errors.Wrapf(err, "unexpected error probing global forwarding rule %s after destroy", name)
	}
	return errors.Errorf("global forwarding rule %s still exists after destroy", name)
}
