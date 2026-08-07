package verify

import (
	"context"

	"github.com/pkg/errors"
	"google.golang.org/api/googleapi"
)

// targetHttpProxyVerifier probes a global target HTTP proxy by name and
// confirms the URL-map wiring — a proxy without a routing table is not a
// working frontend.
type targetHttpProxyVerifier struct{}

func (v *targetHttpProxyVerifier) IDOutputKey() string { return "self_link" }

func (v *targetHttpProxyVerifier) VerifyExists(ctx context.Context, svc *Services, outputs map[string]string) error {
	name := outputs["proxy_name"]
	proxy, err := svc.Compute.TargetHttpProxies.Get(svc.Project, name).Context(ctx).Do()
	if err != nil {
		return errors.Wrapf(err, "target http proxy %s not found after deploy", name)
	}
	if proxy.UrlMap == "" {
		return errors.Errorf("target http proxy %s has no url map wired", name)
	}
	return nil
}

func (v *targetHttpProxyVerifier) VerifyAbsent(ctx context.Context, svc *Services, outputs map[string]string) error {
	name := outputs["proxy_name"]
	_, err := svc.Compute.TargetHttpProxies.Get(svc.Project, name).Context(ctx).Do()
	if err != nil {
		var apiErr *googleapi.Error
		if errors.As(err, &apiErr) && apiErr.Code == 404 {
			return nil
		}
		return errors.Wrapf(err, "unexpected error probing target http proxy %s after destroy", name)
	}
	return errors.Errorf("target http proxy %s still exists after destroy", name)
}
