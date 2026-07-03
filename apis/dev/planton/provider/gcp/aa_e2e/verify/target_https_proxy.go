package verify

import (
	"context"

	"github.com/pkg/errors"
	"google.golang.org/api/googleapi"
)

// targetHttpsProxyVerifier probes a global target HTTPS proxy by name and
// confirms both wiring dimensions: the URL map (routing) and a TLS input
// (certificates, a certificate map, or a server TLS policy) — a TLS
// frontend without either is not serving.
type targetHttpsProxyVerifier struct{}

func (v *targetHttpsProxyVerifier) IDOutputKey() string { return "self_link" }

func (v *targetHttpsProxyVerifier) VerifyExists(ctx context.Context, svc *Services, outputs map[string]string) error {
	name := outputs["proxy_name"]
	proxy, err := svc.Compute.TargetHttpsProxies.Get(svc.Project, name).Context(ctx).Do()
	if err != nil {
		return errors.Wrapf(err, "target https proxy %s not found after deploy", name)
	}
	if proxy.UrlMap == "" {
		return errors.Errorf("target https proxy %s has no url map wired", name)
	}
	if len(proxy.SslCertificates) == 0 && proxy.CertificateMap == "" && proxy.ServerTlsPolicy == "" {
		return errors.Errorf("target https proxy %s has no TLS input (certificates, certificate map, or server TLS policy)", name)
	}
	return nil
}

func (v *targetHttpsProxyVerifier) VerifyAbsent(ctx context.Context, svc *Services, outputs map[string]string) error {
	name := outputs["proxy_name"]
	_, err := svc.Compute.TargetHttpsProxies.Get(svc.Project, name).Context(ctx).Do()
	if err != nil {
		var apiErr *googleapi.Error
		if errors.As(err, &apiErr) && apiErr.Code == 404 {
			return nil
		}
		return errors.Wrapf(err, "unexpected error probing target https proxy %s after destroy", name)
	}
	return errors.Errorf("target https proxy %s still exists after destroy", name)
}
