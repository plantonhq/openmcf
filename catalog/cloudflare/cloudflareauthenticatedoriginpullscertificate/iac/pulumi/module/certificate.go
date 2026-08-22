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
	var status pulumi.StringOutput

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
		status = createdCertificate.Status
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
		status = createdCertificate.Status
	}

	ctx.Export(OpCertificateId, certificateId)
	ctx.Export(OpZoneId, pulumi.String(spec.ZoneId.GetValue()))
	ctx.Export(OpExpiresOn, expiresOn)
	ctx.Export(OpStatus, status)

	return nil
}
