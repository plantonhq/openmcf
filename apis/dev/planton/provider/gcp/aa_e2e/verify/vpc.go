package verify

import (
	"context"

	"github.com/pkg/errors"
	"google.golang.org/api/googleapi"
)

// vpcVerifier probes a VPC network by name via the compute API. Registered so
// the VPC can serve as a verified E2E prerequisite for subnet-scoped kinds.
type vpcVerifier struct{}

func (v *vpcVerifier) IDOutputKey() string { return "network_self_link" }

func (v *vpcVerifier) VerifyExists(ctx context.Context, svc *Services, outputs map[string]string) error {
	name := outputs["network_name"]
	if _, err := svc.Compute.Networks.Get(svc.Project, name).Context(ctx).Do(); err != nil {
		return errors.Wrapf(err, "vpc network %s not found after deploy", name)
	}
	return nil
}

func (v *vpcVerifier) VerifyAbsent(ctx context.Context, svc *Services, outputs map[string]string) error {
	name := outputs["network_name"]
	_, err := svc.Compute.Networks.Get(svc.Project, name).Context(ctx).Do()
	if err != nil {
		var apiErr *googleapi.Error
		if errors.As(err, &apiErr) && apiErr.Code == 404 {
			return nil
		}
		return errors.Wrapf(err, "unexpected error probing vpc network %s after destroy", name)
	}
	return errors.Errorf("vpc network %s still exists after destroy", name)
}
