package module

import (
	"github.com/pkg/errors"
	"github.com/pulumi/pulumi-cloudflare/sdk/v6/go/cloudflare"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

// customSslCertificate uploads the bring-your-own certificate to the zone.
// Rotation is replacement: certificate and private key changes force a
// destroy-and-create at the provider (Cloudflare serves the previous
// certificate until the replacement deploys). priority is deliberately
// absent: at provider v5.23.0 it is read-only (the v4 reprioritization
// surface was dropped).
func customSslCertificate(
	ctx *pulumi.Context,
	locals *Locals,
	cloudflareProvider *cloudflare.Provider,
) error {
	spec := locals.CloudflareCustomSslCertificate.Spec

	args := &cloudflare.CustomSslArgs{
		ZoneId:      pulumi.String(spec.ZoneId.GetValue()),
		Certificate: pulumi.String(spec.Certificate),
		PrivateKey:  pulumi.String(spec.PrivateKey.GetValue()),
	}

	if spec.Type != nil {
		args.Type = pulumi.StringPtr(spec.GetType())
	}
	if spec.BundleMethod != nil {
		args.BundleMethod = pulumi.StringPtr(spec.GetBundleMethod())
	}
	if spec.Deploy != nil {
		args.Deploy = pulumi.StringPtr(spec.GetDeploy())
	}
	if spec.Policy != "" {
		args.Policy = pulumi.StringPtr(spec.Policy)
	}
	if spec.CustomCsrId != "" {
		args.CustomCsrId = pulumi.StringPtr(spec.CustomCsrId)
	}
	if spec.GeoRestrictions != nil && spec.GeoRestrictions.Label != nil {
		args.GeoRestrictions = &cloudflare.CustomSslGeoRestrictionsArgs{
			Label: pulumi.StringPtr(spec.GeoRestrictions.GetLabel()),
		}
	}

	createdCertificate, err := cloudflare.NewCustomSsl(
		ctx,
		"custom_ssl_certificate",
		args,
		pulumi.Provider(cloudflareProvider),
	)
	if err != nil {
		return errors.Wrap(err, "failed to create custom ssl certificate")
	}

	ctx.Export(OpCertificateId, createdCertificate.ID())
	ctx.Export(OpZoneId, pulumi.String(spec.ZoneId.GetValue()))
	ctx.Export(OpExpiresOn, createdCertificate.ExpiresOn)
	// status is deliberately NOT exported: deployment is asynchronous
	// (pending before active), so a point-in-time phase flips on the first
	// refresh after the transition and re-plans forever (the class was
	// measured live 2026-08-28 on the sibling AOP certificate).

	return nil
}
