package verify

import (
	"context"
	"strings"

	"github.com/pkg/errors"
	"google.golang.org/api/googleapi"
)

// routerNatVerifier probes a Cloud Router and its NAT configuration through
// the compute routers API. The NAT has no standalone GET — it lives on the
// router object — so the verifier parses the router's region and name from
// the router_self_link output and asserts the named NAT's allocation posture.
type routerNatVerifier struct{}

func (v *routerNatVerifier) IDOutputKey() string { return "name" }

// routerLocator extracts (region, routerName) from a router self link of the
// form .../projects/{p}/regions/{region}/routers/{name}.
func routerLocator(selfLink string) (string, string, error) {
	parts := strings.Split(selfLink, "/")
	for i := 0; i < len(parts)-3; i++ {
		if parts[i] == "regions" && parts[i+2] == "routers" {
			return parts[i+1], parts[i+3], nil
		}
	}
	return "", "", errors.Errorf("cannot parse region/router from self link %q", selfLink)
}

func (v *routerNatVerifier) VerifyExists(ctx context.Context, svc *Services, outputs map[string]string) error {
	natName := outputs["name"]
	region, routerName, err := routerLocator(outputs["router_self_link"])
	if err != nil {
		return errors.Wrap(err, "after deploy")
	}

	router, err := svc.Compute.Routers.Get(svc.Project, region, routerName).Context(ctx).Do()
	if err != nil {
		return errors.Wrapf(err, "router %s in %s not found after deploy", routerName, region)
	}

	for _, nat := range router.Nats {
		if nat.Name != natName {
			continue
		}
		// Allocation posture: a non-empty nat_ips output means manual
		// allocation with exactly those reservations attached.
		if outputs["nat_ips"] != "" && outputs["nat_ips"] != "[]" {
			if nat.NatIpAllocateOption != "MANUAL_ONLY" {
				return errors.Errorf("nat %s: nat_ips output set but allocation is %q, want MANUAL_ONLY",
					natName, nat.NatIpAllocateOption)
			}
			if len(nat.NatIps) == 0 {
				return errors.Errorf("nat %s: manual allocation but no NAT IPs attached live", natName)
			}
		}
		return nil
	}
	return errors.Errorf("router %s in %s has no NAT named %s after deploy", routerName, region, natName)
}

func (v *routerNatVerifier) VerifyAbsent(ctx context.Context, svc *Services, outputs map[string]string) error {
	natName := outputs["name"]
	region, routerName, err := routerLocator(outputs["router_self_link"])
	if err != nil {
		return nil
	}

	router, err := svc.Compute.Routers.Get(svc.Project, region, routerName).Context(ctx).Do()
	if err != nil {
		var apiErr *googleapi.Error
		if errors.As(err, &apiErr) && apiErr.Code == 404 {
			// The router itself is gone — the NAT went with it.
			return nil
		}
		return errors.Wrapf(err, "unexpected error probing router %s in %s after destroy", routerName, region)
	}
	for _, nat := range router.Nats {
		if nat.Name == natName {
			return errors.Errorf("nat %s still exists on router %s after destroy", natName, routerName)
		}
	}
	return nil
}
