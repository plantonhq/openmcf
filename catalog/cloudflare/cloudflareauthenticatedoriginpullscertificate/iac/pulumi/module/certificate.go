package module

import (
	"github.com/pkg/errors"
	"github.com/pulumi/pulumi-cloudflare/sdk/v6/go/cloudflare"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

// certificate uploads the client certificate to the surface the scope
// selects: zone replaces the zone-wide client certificate; hostname uploads a
// certificate for per-hostname associations to reference. Rotation is
// replacement on both surfaces. Never rotate only the private key against
// the same certificate: the zone-scoped API silently ignores a key-only
// change (a measured provider defect at v5.23.0 -- its Update is empty and
// the key does not force replacement), so key and certificate must always
// change together.
func certificate(
	ctx *pulumi.Context,
	locals *Locals,
	cloudflareProvider *cloudflare.Provider,
) error {
	spec := locals.CloudflareAuthenticatedOriginPullsCertificate.Spec

	isZoneScope := spec.Scope == nil || spec.GetScope() == "zone"

	var certificateId pulumi.StringOutput
	var expiresOn pulumi.StringOutput

	if isZoneScope {
		createdCertificate, err := cloudflare.NewAuthenticatedOriginPullsCertificate(
			ctx,
			"certificate",
			&cloudflare.AuthenticatedOriginPullsCertificateArgs{
				ZoneId:      pulumi.String(spec.ZoneId.GetValue()),
				Certificate: pulumi.String(spec.Certificate),
				PrivateKey:  pulumi.String(spec.PrivateKey.GetValue()),
			},
			pulumi.Provider(cloudflareProvider),
		)
		if err != nil {
			return errors.Wrap(err, "failed to create zone-scoped certificate")
		}
		certificateId = createdCertificate.ID().ToStringOutput()
		expiresOn = createdCertificate.ExpiresOn
	} else {
		createdCertificate, err := cloudflare.NewAuthenticatedOriginPullsHostnameCertificate(
			ctx,
			"certificate",
			&cloudflare.AuthenticatedOriginPullsHostnameCertificateArgs{
				ZoneId:      pulumi.String(spec.ZoneId.GetValue()),
				Certificate: pulumi.String(spec.Certificate),
				PrivateKey:  pulumi.String(spec.PrivateKey.GetValue()),
			},
			pulumi.Provider(cloudflareProvider),
		)
		if err != nil {
			return errors.Wrap(err, "failed to create hostname-scoped certificate")
		}
		certificateId = createdCertificate.ID().ToStringOutput()
		expiresOn = createdCertificate.ExpiresOn
	}

	ctx.Export(OpCertificateId, certificateId)
	ctx.Export(OpZoneId, pulumi.String(spec.ZoneId.GetValue()))
	ctx.Export(OpExpiresOn, expiresOn)
	// status is deliberately NOT exported: deployment is asynchronous
	// (pending_deployment -> active seconds after create), so a
	// point-in-time phase flips on the first refresh after the transition
	// and re-plans forever (measured live 2026-08-28).

	return nil
}
