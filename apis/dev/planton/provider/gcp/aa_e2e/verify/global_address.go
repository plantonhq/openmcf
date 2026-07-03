package verify

import (
	"context"
	"strings"

	"github.com/pkg/errors"
	"google.golang.org/api/googleapi"
)

// globalAddressVerifier probes a global address reservation. The kind's
// outputs carry no separate name field, so the name is derived from the
// self_link's final path segment (the stable GCP naming shape).
type globalAddressVerifier struct{}

func (v *globalAddressVerifier) IDOutputKey() string { return "self_link" }

func globalAddressName(outputs map[string]string) string {
	selfLink := outputs["self_link"]
	if idx := strings.LastIndex(selfLink, "/"); idx >= 0 {
		return selfLink[idx+1:]
	}
	return selfLink
}

func (v *globalAddressVerifier) VerifyExists(ctx context.Context, svc *Services, outputs map[string]string) error {
	name := globalAddressName(outputs)
	address, err := svc.Compute.GlobalAddresses.Get(svc.Project, name).Context(ctx).Do()
	if err != nil {
		return errors.Wrapf(err, "global address %s not found after deploy", name)
	}
	if address.Address == "" {
		return errors.Errorf("global address %s has no IP assigned", name)
	}
	return nil
}

func (v *globalAddressVerifier) VerifyAbsent(ctx context.Context, svc *Services, outputs map[string]string) error {
	name := globalAddressName(outputs)
	_, err := svc.Compute.GlobalAddresses.Get(svc.Project, name).Context(ctx).Do()
	if err != nil {
		var apiErr *googleapi.Error
		if errors.As(err, &apiErr) && apiErr.Code == 404 {
			return nil
		}
		return errors.Wrapf(err, "unexpected error probing global address %s after destroy", name)
	}
	return errors.Errorf("global address %s still exists after destroy", name)
}
