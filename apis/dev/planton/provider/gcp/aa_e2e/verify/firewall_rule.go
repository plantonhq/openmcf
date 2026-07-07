package verify

import (
	"context"
	"strings"

	"github.com/pkg/errors"
	"google.golang.org/api/googleapi"
)

// firewallRuleVerifier probes a VPC firewall rule via the Compute API.
// The firewall_self_link output carries the compound identifier
// (projects/{p}/global/firewalls/{name}); the probe parses it rather
// than assuming the harness project so cross-project rules verify
// honestly.
type firewallRuleVerifier struct{}

func (v *firewallRuleVerifier) IDOutputKey() string { return "firewall_self_link" }

// parseFirewallSelfLink extracts (project, name) from a firewall
// self-link URL or relative path.
func parseFirewallSelfLink(selfLink string) (project, name string, err error) {
	parts := strings.Split(selfLink, "/")
	for i := 0; i < len(parts)-1; i++ {
		if parts[i] == "projects" && i+1 < len(parts) {
			project = parts[i+1]
		}
		if parts[i] == "firewalls" && i+1 < len(parts) {
			name = parts[i+1]
		}
	}
	if project == "" || name == "" {
		return "", "", errors.Errorf("firewall_self_link %q does not carry projects/{p}/.../firewalls/{name}", selfLink)
	}
	return project, name, nil
}

func (v *firewallRuleVerifier) VerifyExists(ctx context.Context, svc *Services, outputs map[string]string) error {
	selfLink := outputs["firewall_self_link"]
	if selfLink == "" {
		return errors.New("firewall_self_link output missing after deploy")
	}
	project, name, err := parseFirewallSelfLink(selfLink)
	if err != nil {
		return err
	}

	firewall, err := svc.Compute.Firewalls.Get(project, name).Context(ctx).Do()
	if err != nil {
		return errors.Wrapf(err, "firewall rule %s not found after deploy", selfLink)
	}
	// Posture: a rule with neither allow nor deny entries matched nothing —
	// proof the protocol/port payload landed, not just the shell.
	if len(firewall.Allowed) == 0 && len(firewall.Denied) == 0 {
		return errors.Errorf("firewall rule %s has no allow or deny entries after deploy", selfLink)
	}
	return nil
}

func (v *firewallRuleVerifier) VerifyAbsent(ctx context.Context, svc *Services, outputs map[string]string) error {
	selfLink := outputs["firewall_self_link"]
	if selfLink == "" {
		return nil
	}
	project, name, err := parseFirewallSelfLink(selfLink)
	if err != nil {
		return err
	}

	_, err = svc.Compute.Firewalls.Get(project, name).Context(ctx).Do()
	if err == nil {
		return errors.Errorf("firewall rule %s still exists after destroy", selfLink)
	}
	var apiErr *googleapi.Error
	if errors.As(err, &apiErr) && apiErr.Code == 404 {
		return nil
	}
	return errors.Wrapf(err, "unexpected error probing firewall rule %s after destroy", selfLink)
}
