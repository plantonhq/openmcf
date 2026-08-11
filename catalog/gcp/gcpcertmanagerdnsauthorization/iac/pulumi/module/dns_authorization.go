package module

import (
	"github.com/pkg/errors"
	"github.com/pulumi/pulumi-gcp/sdk/v9/go/gcp"
	"github.com/pulumi/pulumi-gcp/sdk/v9/go/gcp/certificatemanager"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

// dnsAuthorization provisions one Certificate Manager DNS authorization:
// the proof-of-domain-control a Google-managed certificate consumes. It
// covers the domain AND its wildcard, and exports the CNAME validation
// record a GcpDnsRecord composes into the zone — validation can complete
// BEFORE any certificate exists, which is what makes zero-downtime
// certificate migration possible.
func dnsAuthorization(ctx *pulumi.Context, locals *Locals, gcpProvider *gcp.Provider, projectService pulumi.Resource) error {
	spec := locals.GcpCertManagerDnsAuthorization.Spec

	args := &certificatemanager.DnsAuthorizationArgs{
		Name:   pulumi.String(locals.AuthorizationName),
		Domain: pulumi.String(spec.Domain),
		Labels: pulumi.ToStringMap(locals.GcpLabels),
	}

	// Empty project falls back to the provider's default project — the same
	// ambient contract the Terraform module honors.
	if locals.ProjectId != "" {
		args.Project = pulumi.StringPtr(locals.ProjectId)
	}

	if spec.Description != "" {
		args.Description = pulumi.StringPtr(spec.Description)
	}

	if spec.Location != "" {
		args.Location = pulumi.StringPtr(spec.Location)
	}

	if spec.Type != "" {
		args.Type = pulumi.StringPtr(spec.Type)
	}

	// Unset defers to the provider default (DELETE).
	if spec.DeletionPolicy != "" {
		args.DeletionPolicy = pulumi.StringPtr(spec.DeletionPolicy)
	}

	createdAuthorization, err := certificatemanager.NewDnsAuthorization(ctx,
		locals.GcpCertManagerDnsAuthorization.Metadata.Name,
		args,
		pulumi.Provider(gcpProvider),
		pulumi.DependsOn([]pulumi.Resource{projectService}))
	if err != nil {
		return errors.Wrap(err, "failed to create dns authorization")
	}

	ctx.Export(OpAuthorizationId, createdAuthorization.ID())
	ctx.Export(OpAuthorizationName, createdAuthorization.Name)
	ctx.Export(OpDomain, createdAuthorization.Domain)

	// The validation record is a computed one-element list. Export via one
	// struct-slice ApplyT with len guards degrading to "" — the
	// optional-output export contract shared with the Terraform module's
	// try() guards.
	records := createdAuthorization.DnsResourceRecords
	ctx.Export(OpDnsRecordName, records.ApplyT(func(items []certificatemanager.DnsAuthorizationDnsResourceRecord) string {
		if len(items) == 0 || items[0].Name == nil {
			return ""
		}
		return *items[0].Name
	}).(pulumi.StringOutput))
	ctx.Export(OpDnsRecordType, records.ApplyT(func(items []certificatemanager.DnsAuthorizationDnsResourceRecord) string {
		if len(items) == 0 || items[0].Type == nil {
			return ""
		}
		return *items[0].Type
	}).(pulumi.StringOutput))
	ctx.Export(OpDnsRecordData, records.ApplyT(func(items []certificatemanager.DnsAuthorizationDnsResourceRecord) string {
		if len(items) == 0 || items[0].Data == nil {
			return ""
		}
		return *items[0].Data
	}).(pulumi.StringOutput))

	return nil
}
