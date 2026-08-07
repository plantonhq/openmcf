package verify

import (
	"context"

	"github.com/pkg/errors"
	"google.golang.org/api/googleapi"
)

// serverlessVpcConnectorVerifier probes a Serverless VPC Access connector via
// the vpcaccess API. The posture assertion confirms the connector reconciled
// to READY — a connector stuck in CREATING or ERROR would "exist" but carry
// no traffic.
type serverlessVpcConnectorVerifier struct{}

func (v *serverlessVpcConnectorVerifier) IDOutputKey() string { return "self_link" }

// connectorPath returns the fully qualified resource name. The self_link
// output IS the projects/*/locations/*/connectors/* path the API expects.
func (v *serverlessVpcConnectorVerifier) connectorPath(outputs map[string]string) (string, error) {
	path := outputs["self_link"]
	if path == "" {
		return "", errors.New("self_link output missing")
	}
	return path, nil
}

func (v *serverlessVpcConnectorVerifier) VerifyExists(ctx context.Context, svc *Services, outputs map[string]string) error {
	path, err := v.connectorPath(outputs)
	if err != nil {
		return errors.Wrap(err, "after deploy")
	}

	connector, err := svc.VpcAccess.Projects.Locations.Connectors.Get(path).Context(ctx).Do()
	if err != nil {
		return errors.Wrapf(err, "serverless vpc connector %s not found after deploy", path)
	}

	if connector.State != "READY" {
		return errors.Errorf("serverless vpc connector %s state is %q, want READY", path, connector.State)
	}
	// Placement posture mirrors the API's two arms: network placement
	// reports the carved range in ip_cidr_range, while subnet placement
	// leaves ip_cidr_range EMPTY and reports the occupied subnet instead.
	// A connector with neither never materialized a placement.
	if connector.IpCidrRange == "" && (connector.Subnet == nil || connector.Subnet.Name == "") {
		return errors.Errorf("serverless vpc connector %s has neither an ip cidr range nor a subnet after deploy", path)
	}
	return nil
}

func (v *serverlessVpcConnectorVerifier) VerifyAbsent(ctx context.Context, svc *Services, outputs map[string]string) error {
	path, err := v.connectorPath(outputs)
	if err != nil {
		return nil
	}

	_, err = svc.VpcAccess.Projects.Locations.Connectors.Get(path).Context(ctx).Do()
	if err == nil {
		return errors.Errorf("serverless vpc connector %s still exists after destroy", path)
	}
	var apiErr *googleapi.Error
	if errors.As(err, &apiErr) && apiErr.Code == 404 {
		return nil
	}
	return errors.Wrapf(err, "unexpected error probing serverless vpc connector %s after destroy", path)
}
