package verify

import (
	"context"

	"github.com/pkg/errors"
	"google.golang.org/api/googleapi"
)

// addressVerifier probes a regional Compute Engine address reservation.
type addressVerifier struct{}

func (v *addressVerifier) IDOutputKey() string { return "self_link" }

func (v *addressVerifier) VerifyExists(ctx context.Context, svc *Services, outputs map[string]string) error {
	name := outputs["name"]
	region := outputs["region"]
	if name == "" || region == "" {
		return errors.New("name or region output missing after deploy")
	}

	addr, err := svc.Compute.Addresses.Get(svc.Project, region, name).Context(ctx).Do()
	if err != nil {
		return errors.Wrapf(err, "regional address %s in %s not found after deploy", name, region)
	}
	if addr.Address == "" {
		return errors.Errorf("regional address %s in %s has no IP assigned", name, region)
	}
	return nil
}

func (v *addressVerifier) VerifyAbsent(ctx context.Context, svc *Services, outputs map[string]string) error {
	name := outputs["name"]
	region := outputs["region"]
	if name == "" || region == "" {
		return nil
	}

	_, err := svc.Compute.Addresses.Get(svc.Project, region, name).Context(ctx).Do()
	if err != nil {
		var apiErr *googleapi.Error
		if errors.As(err, &apiErr) && apiErr.Code == 404 {
			return nil
		}
		return errors.Wrapf(err, "unexpected error probing regional address %s in %s after destroy", name, region)
	}
	return errors.Errorf("regional address %s in %s still exists after destroy", name, region)
}
