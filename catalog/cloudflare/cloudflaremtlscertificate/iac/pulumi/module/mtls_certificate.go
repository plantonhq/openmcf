package module

import (
	"github.com/pkg/errors"
	"github.com/pulumi/pulumi-cloudflare/sdk/v6/go/cloudflare"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

// mtlsCertificate uploads the certificate to the account-level mTLS store.
// Every argument is create-only at the API -- any change replaces the upload
// and the certificate id changes, so consumers must re-point at the new id
// (rotate by replace, never in place).
func mtlsCertificate(
	ctx *pulumi.Context,
	locals *Locals,
	cloudflareProvider *cloudflare.Provider,
) error {
	spec := locals.CloudflareMtlsCertificate.Spec

	args := &cloudflare.MtlsCertificateArgs{
		AccountId:    pulumi.String(spec.AccountId),
		Ca:           pulumi.Bool(spec.GetCa()),
		Certificates: pulumi.String(spec.Certificates),
	}

	if spec.Name != "" {
		args.Name = pulumi.StringPtr(spec.Name)
	}
	// The private key is optional: CA uploads used only to validate clients
	// carry no key.
	if spec.PrivateKey.GetValue() != "" {
		args.PrivateKey = pulumi.StringPtr(spec.PrivateKey.GetValue())
	}

	createdCertificate, err := cloudflare.NewMtlsCertificate(
		ctx,
		"mtls_certificate",
		args,
		pulumi.Provider(cloudflareProvider),
	)
	if err != nil {
		return errors.Wrap(err, "failed to create mtls certificate")
	}

	ctx.Export(OpCertificateId, createdCertificate.ID())
	ctx.Export(OpExpiresOn, createdCertificate.ExpiresOn)
	ctx.Export(OpSerialNumber, createdCertificate.SerialNumber)

	return nil
}
