package verify

import (
	"context"
	"strings"
	"time"

	"github.com/pkg/errors"
	"google.golang.org/api/googleapi"
)

// serviceNetworkingConnectionVerifier probes the VPC peering created by a
// private services access connection. The connection has no standalone GET
// API — existence is verified by finding the servicenetworking peering on the
// parent network.
type serviceNetworkingConnectionVerifier struct{}

func (v *serviceNetworkingConnectionVerifier) IDOutputKey() string { return "peering" }

func (v *serviceNetworkingConnectionVerifier) VerifyExists(ctx context.Context, svc *Services, outputs map[string]string) error {
	networkURL := outputs["network"]
	if networkURL == "" {
		return errors.New("network output missing after deploy")
	}
	networkName := networkNameFromSelfLink(networkURL)
	if networkName == "" {
		return errors.Errorf("could not parse network name from %q", networkURL)
	}

	peeringName := outputs["peering"]
	if peeringName == "" {
		return errors.New("peering output missing after deploy")
	}

	// The peering is created by the Service Networking API but read back
	// through the Compute API, and that cross-API view is eventually
	// consistent: a networks.get issued immediately after a successful
	// connection create can miss a peering that is already there. Poll
	// briefly before declaring it absent.
	var lastErr error
	deadline := time.Now().Add(2 * time.Minute)
	for {
		net, err := svc.Compute.Networks.Get(svc.Project, networkName).Context(ctx).Do()
		if err != nil {
			return errors.Wrapf(err, "network %s not found after deploy", networkName)
		}

		lastErr = errors.Errorf("peering %s not found on network %s after deploy", peeringName, networkName)
		for _, p := range net.Peerings {
			if p.Name == peeringName {
				if p.State != "ACTIVE" && p.State != "" {
					lastErr = errors.Errorf("peering %s exists but state is %q, want ACTIVE", peeringName, p.State)
					break
				}
				return nil
			}
		}

		if time.Now().After(deadline) {
			return lastErr
		}
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-time.After(10 * time.Second):
		}
	}
}

func (v *serviceNetworkingConnectionVerifier) VerifyAbsent(ctx context.Context, svc *Services, outputs map[string]string) error {
	networkURL := outputs["network"]
	networkName := networkNameFromSelfLink(networkURL)
	if networkName == "" {
		return nil
	}

	net, err := svc.Compute.Networks.Get(svc.Project, networkName).Context(ctx).Do()
	if err != nil {
		var apiErr *googleapi.Error
		if errors.As(err, &apiErr) && apiErr.Code == 404 {
			return nil
		}
		return errors.Wrapf(err, "unexpected error probing network %s after destroy", networkName)
	}

	peeringName := outputs["peering"]
	for _, p := range net.Peerings {
		if p.Name == peeringName {
			return errors.Errorf("peering %s still exists on network %s after destroy", peeringName, networkName)
		}
	}
	return nil
}

func networkNameFromSelfLink(selfLink string) string {
	if idx := strings.LastIndex(selfLink, "/"); idx >= 0 {
		return selfLink[idx+1:]
	}
	return selfLink
}
