package module

import (
	"github.com/pkg/errors"
	"github.com/pulumi/pulumi-digitalocean/sdk/v4/go/digitalocean"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

// cdn provisions the CDN endpoint and exports its outputs. The certificate
// is wired through the SDK's CertificateName input ONLY: the numeric
// certificate_id argument is deprecated upstream (Let's Encrypt renewals
// rotate the UUID) and its update path can silently detach the certificate
// when the custom domain changes, so it is deliberately never set.
func cdn(
	ctx *pulumi.Context,
	locals *Locals,
	digitalOceanProvider *digitalocean.Provider,
) (*digitalocean.Cdn, error) {
	spec := locals.DigitalOceanCdn.Spec

	args := &digitalocean.CdnArgs{
		// The origin reference resolves to the Space's fully-qualified
		// domain name before the module runs. Create-only: changing it
		// replaces the endpoint.
		Origin: pulumi.String(spec.Origin.GetValue()),
	}

	// Unset defers to DigitalOcean's default of 3600 seconds, which the
	// provider reads back without a perpetual diff. Spec validation floors
	// ttl at 1 -- an explicit zero can never reach the API.
	if spec.Ttl != nil {
		args.Ttl = pulumi.IntPtr(int(*spec.Ttl))
	}

	// The certificate reference resolves to the certificate's stable NAME
	// (also its resource id at the current provider pin). The
	// "needs-cloudflare-cert" sentinel passes through verbatim.
	if spec.Certificate.GetValue() != "" {
		args.CertificateName = pulumi.String(spec.Certificate.GetValue())
	}

	if spec.CustomDomain != "" {
		args.CustomDomain = pulumi.String(spec.CustomDomain)
	}

	createdCdn, err := digitalocean.NewCdn(
		ctx,
		"cdn",
		args,
		pulumi.Provider(digitalOceanProvider),
	)
	if err != nil {
		return nil, errors.Wrap(err, "failed to create digitalocean cdn")
	}

	ctx.Export(OpCdnId, createdCdn.ID())
	ctx.Export(OpEndpoint, createdCdn.Endpoint)

	return createdCdn, nil
}
