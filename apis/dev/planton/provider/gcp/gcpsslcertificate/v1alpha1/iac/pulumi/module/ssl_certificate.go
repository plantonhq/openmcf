package module

import (
	"strconv"

	"github.com/pkg/errors"
	"github.com/pulumi/pulumi-gcp/sdk/v9/go/gcp"
	"github.com/pulumi/pulumi-gcp/sdk/v9/go/gcp/compute"
	"github.com/pulumi/pulumi-gcp/sdk/v9/go/gcp/projects"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

// sslCertificate provisions a self-managed Compute Engine SSL certificate —
// the user brings the PEM chain and private key; the load balancer presents
// them to clients. Target HTTPS (and SSL) proxies reference the
// certificate's self_link exactly like a Google-managed certificate (the two
// share one API collection and name namespace), but nothing here renews
// itself: track expire_time and rotate.
//
// One kind, two provider resources: GCP models global and regional SSL
// certificates as separate API collections with an identical surface, so an
// empty spec.region creates compute.SSLCertificate and a set region creates
// compute.RegionSslCertificate — mirroring the Terraform module's count
// guards.
//
// EVERY argument is immutable (ForceNew in the provider): any change
// destroys and recreates the certificate. Because a certificate attached to
// a proxy cannot be deleted (GCP returns resourceInUseByAnotherResource),
// rotation is create-before-destroy: create the replacement under a new
// name, repoint the proxy's certificate list, then destroy this one.
//
// The private key is write-only in GCP: the API never returns it, and this
// module marks it a Pulumi secret so it is encrypted in state and never
// surfaced in outputs.
func sslCertificate(ctx *pulumi.Context, locals *Locals, gcpProvider *gcp.Provider) error {
	spec := locals.GcpSslCertificate.Spec

	// Enable the Compute Engine API so a fresh project can host the
	// certificate. disable_on_destroy stays false (the provider default):
	// tearing down one certificate must never disable the API for everything
	// else in the project. Matches the Terraform module.
	serviceArgs := &projects.ServiceArgs{
		Service:                  pulumi.String("compute.googleapis.com"),
		DisableDependentServices: pulumi.BoolPtr(true),
	}
	if spec.ProjectId.GetValue() != "" {
		serviceArgs.Project = pulumi.String(spec.ProjectId.GetValue())
	}
	createdProjectService, err := projects.NewService(ctx,
		"sslcert-compute.googleapis.com", serviceArgs, pulumi.Provider(gcpProvider))
	if err != nil {
		return errors.Wrap(err, "failed to enable compute.googleapis.com api")
	}

	if spec.Region == "" {
		return globalSslCertificate(ctx, locals, gcpProvider, createdProjectService)
	}
	return regionalSslCertificate(ctx, locals, gcpProvider, createdProjectService)
}

func globalSslCertificate(ctx *pulumi.Context, locals *Locals, gcpProvider *gcp.Provider, computeApiService pulumi.Resource) error {
	spec := locals.GcpSslCertificate.Spec

	args := &compute.SSLCertificateArgs{
		Name:        pulumi.String(locals.CertificateName),
		Certificate: pulumi.String(spec.Certificate),
		// Secret material; write-only in GCP and never surfaced in outputs.
		// ToSecret marks it encrypted in Pulumi state.
		PrivateKey: pulumi.ToSecret(pulumi.String(spec.PrivateKey)).(pulumi.StringOutput),
	}

	// Omitted description stays unset (matching the Terraform module's null)
	// rather than being sent as an empty string.
	if spec.Description != "" {
		args.Description = pulumi.String(spec.Description)
	}

	// Honor the spec contract: an empty project_id falls back to the provider's
	// default project. Leaving Project unset lets the gcp provider resolve its
	// own project (configuration or the GOOGLE_PROJECT / GOOGLE_CLOUD_PROJECT
	// environment chain); an empty string would be sent verbatim and rejected.
	if spec.ProjectId.GetValue() != "" {
		args.Project = pulumi.String(spec.ProjectId.GetValue())
	}

	createdCertificate, err := compute.NewSSLCertificate(ctx, "ssl-certificate", args,
		pulumi.Provider(gcpProvider), pulumi.DependsOn([]pulumi.Resource{computeApiService}))
	if err != nil {
		return errors.Wrap(err, "failed to create global ssl certificate")
	}

	ctx.Export(OpSelfLink, createdCertificate.SelfLink)
	ctx.Export(OpCertificateName, createdCertificate.Name)
	ctx.Export(OpCertificateId, createdCertificate.CertificateId.ApplyT(func(id int) string {
		return strconv.Itoa(id)
	}).(pulumi.StringOutput))
	ctx.Export(OpExpireTime, createdCertificate.ExpireTime)
	// Empty region marks the global scope for downstream composition checks.
	ctx.Export(OpRegion, pulumi.String(""))

	return nil
}

func regionalSslCertificate(ctx *pulumi.Context, locals *Locals, gcpProvider *gcp.Provider, computeApiService pulumi.Resource) error {
	spec := locals.GcpSslCertificate.Spec

	args := &compute.RegionSslCertificateArgs{
		Name:        pulumi.String(locals.CertificateName),
		Region:      pulumi.String(spec.Region),
		Certificate: pulumi.String(spec.Certificate),
		// Secret material; write-only in GCP and never surfaced in outputs.
		PrivateKey: pulumi.ToSecret(pulumi.String(spec.PrivateKey)).(pulumi.StringOutput),
	}

	if spec.Description != "" {
		args.Description = pulumi.String(spec.Description)
	}

	if spec.ProjectId.GetValue() != "" {
		args.Project = pulumi.String(spec.ProjectId.GetValue())
	}

	createdCertificate, err := compute.NewRegionSslCertificate(ctx, "ssl-certificate", args,
		pulumi.Provider(gcpProvider), pulumi.DependsOn([]pulumi.Resource{computeApiService}))
	if err != nil {
		return errors.Wrap(err, "failed to create regional ssl certificate")
	}

	ctx.Export(OpSelfLink, createdCertificate.SelfLink)
	ctx.Export(OpCertificateName, createdCertificate.Name)
	ctx.Export(OpCertificateId, createdCertificate.CertificateId.ApplyT(func(id int) string {
		return strconv.Itoa(id)
	}).(pulumi.StringOutput))
	ctx.Export(OpExpireTime, createdCertificate.ExpireTime)
	// Export the plain region NAME from the spec (matching the Terraform
	// module) — the provider's region attribute can carry a self-link, which
	// API callers and verification cannot use directly.
	ctx.Export(OpRegion, pulumi.String(spec.Region))

	return nil
}
