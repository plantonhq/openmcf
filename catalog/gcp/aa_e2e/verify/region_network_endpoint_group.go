package verify

import (
	"context"

	"github.com/pkg/errors"
	"google.golang.org/api/googleapi"
)

// regionNetworkEndpointGroupVerifier probes a regional network endpoint group
// by name and region via the compute API.
type regionNetworkEndpointGroupVerifier struct{}

func (v *regionNetworkEndpointGroupVerifier) IDOutputKey() string { return "self_link" }

func (v *regionNetworkEndpointGroupVerifier) VerifyExists(ctx context.Context, svc *Services, outputs map[string]string) error {
	name := outputs["network_endpoint_group_name"]
	region := outputs["region"]
	neg, err := svc.Compute.RegionNetworkEndpointGroups.Get(svc.Project, region, name).Context(ctx).Do()
	if err != nil {
		return errors.Wrapf(err, "region network endpoint group %s/%s not found after deploy", region, name)
	}
	if neg.NetworkEndpointType == "" {
		return errors.Errorf("region network endpoint group %s/%s has empty network_endpoint_type", region, name)
	}
	return nil
}

func (v *regionNetworkEndpointGroupVerifier) VerifyAbsent(ctx context.Context, svc *Services, outputs map[string]string) error {
	name := outputs["network_endpoint_group_name"]
	region := outputs["region"]
	_, err := svc.Compute.RegionNetworkEndpointGroups.Get(svc.Project, region, name).Context(ctx).Do()
	if err != nil {
		var apiErr *googleapi.Error
		if errors.As(err, &apiErr) && apiErr.Code == 404 {
			return nil
		}
		return errors.Wrapf(err, "unexpected error probing region NEG %s/%s after destroy", region, name)
	}
	return errors.Errorf("region network endpoint group %s/%s still exists after destroy", region, name)
}
