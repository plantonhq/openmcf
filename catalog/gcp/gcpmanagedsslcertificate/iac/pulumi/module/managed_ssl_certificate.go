package module

import (
	"strconv"

	"github.com/pkg/errors"
	"github.com/pulumi/pulumi-gcp/sdk/v9/go/gcp"
	"github.com/pulumi/pulumi-gcp/sdk/v9/go/gcp/compute"
	"github.com/pulumi/pulumi-gcp/sdk/v9/go/gcp/projects"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

// managedSslCertificate provisions a Google-managed SSL certificate — a global
// Compute Engine SSL certificate whose private key and issuance are handled
// entirely by Google. Attach its self_link to a target HTTPS proxy to
// terminate TLS at a global external Application Load Balancer.
//
// The whole resource is immutable in GCP (name and domains are ForceNew): any
// change destroys and recreates the certificate. Because a cert attached to a
// proxy cannot be deleted, rotate by creating the replacement first and
// swapping the proxy's ssl_certificates reference before destroying the old
// one (create-before-destroy) — otherwise the destroy fails and TLS drops
// during the gap.
//
// Provisioning is asynchronous and DNS-gated: creation returns immediately,
// but the certificate stays PROVISIONING until each domain's DNS points at the
// load balancer's IP. expire_time stays empty until provisioning completes.
func managedSslCertificate(ctx *pulumi.Context, locals *Locals, gcpProvider *gcp.Provider) error {
	spec := locals.GcpManagedSslCertificate.Spec

	// Enable the Compute Engine API so a fresh project can host the certificate.
	// disable_on_destroy stays false (the provider default): tearing down one
	// certificate must never disable the API for everything else in the project.
	// Matches the Terraform module.
	serviceArgs := &projects.ServiceArgs{
		Service:                  pulumi.String("compute.googleapis.com"),
		DisableDependentServices: pulumi.BoolPtr(true),
	}
	if spec.ProjectId.GetValue() != "" {
		serviceArgs.Project = pulumi.String(spec.ProjectId.GetValue())
	}
	createdProjectService, err := projects.NewService(ctx,
		"managedssl-compute.googleapis.com", serviceArgs, pulumi.Provider(gcpProvider))
	if err != nil {
		return errors.Wrap(err, "failed to enable compute.googleapis.com api")
	}

	args := &compute.ManagedSslCertificateArgs{
		Name: pulumi.String(locals.CertificateName),
		Managed: &compute.ManagedSslCertificateManagedArgs{
			Domains: pulumi.ToStringArray(spec.Domains),
		},
	}

	// Honor the spec contract: an empty project_id falls back to the provider's
	// default project. Leaving Project unset lets the gcp provider resolve its
	// own project (configuration or the GOOGLE_PROJECT / GOOGLE_CLOUD_PROJECT
	// environment chain); an empty string would be sent verbatim and rejected.
	if spec.ProjectId.GetValue() != "" {
		args.Project = pulumi.String(spec.ProjectId.GetValue())
	}

	if spec.Description != "" {
		args.Description = pulumi.String(spec.Description)
	}

	createdCert, err := compute.NewManagedSslCertificate(ctx, "managed-ssl-certificate", args,
		pulumi.Provider(gcpProvider), pulumi.DependsOn([]pulumi.Resource{createdProjectService}))
	if err != nil {
		return errors.Wrap(err, "failed to create managed ssl certificate")
	}

	ctx.Export(OpSelfLink, createdCert.SelfLink)
	ctx.Export(OpCertificateName, createdCert.Name)
	ctx.Export(OpCertificateId, createdCert.CertificateId.ApplyT(func(id int) string {
		return strconv.Itoa(id)
	}).(pulumi.StringOutput))
	ctx.Export(OpExpireTime, createdCert.ExpireTime)

	return nil
}
