package module

import (
	"github.com/pkg/errors"
	"github.com/pulumi/pulumi-digitalocean/sdk/v4/go/digitalocean"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

// certificate provisions the DigitalOcean SSL certificate and exports its outputs.
//
// The spec's certificate_source oneof picks the branch, and the branch derives
// DigitalOcean's `type` argument. Every argument is create-only, so any change
// replaces the certificate; DeleteBeforeReplace stays false (the default) so the
// replacement is created first and consumers referencing the certificate by its
// stable name never observe a gap.
func certificate(
	ctx *pulumi.Context,
	locals *Locals,
	digitalOceanProvider *digitalocean.Provider,
) (*digitalocean.Certificate, error) {
	spec := locals.DigitalOceanCertificate.Spec

	certArgs := &digitalocean.CertificateArgs{
		Name: pulumi.String(spec.CertificateName),
	}

	if letsEncrypt := spec.GetLetsEncrypt(); letsEncrypt != nil {
		certArgs.Type = pulumi.String("lets_encrypt")
		var domains pulumi.StringArray
		for _, d := range letsEncrypt.Domains {
			domains = append(domains, pulumi.String(d))
		}
		certArgs.Domains = domains
	}

	if custom := spec.GetCustom(); custom != nil {
		certArgs.Type = pulumi.String("custom")
		// The provider stores only hashes of the PEM material (the
		// DigitalOcean API never returns it).
		certArgs.LeafCertificate = pulumi.String(custom.LeafCertificate)
		certArgs.PrivateKey = pulumi.String(custom.PrivateKey)
		if custom.CertificateChain != "" {
			certArgs.CertificateChain = pulumi.StringPtr(custom.CertificateChain)
		}
	}

	createdCertificate, err := digitalocean.NewCertificate(
		ctx,
		"certificate",
		certArgs,
		pulumi.Provider(digitalOceanProvider),
	)
	if err != nil {
		return nil, errors.Wrap(err, "failed to create digitalocean certificate")
	}

	// The resource id is the certificate NAME at the current provider pin (a
	// Let's Encrypt certificate's UUID rotates on every auto-renewal).
	ctx.Export(OpCertificateId, createdCertificate.ID())
	ctx.Export(OpExpiryRfc3339, createdCertificate.NotAfter)

	return createdCertificate, nil
}
