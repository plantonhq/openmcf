package module

import (
	"github.com/pkg/errors"
	"github.com/pulumi/pulumi-gcp/sdk/v9/go/gcp"
	"github.com/pulumi/pulumi-gcp/sdk/v9/go/gcp/certificatemanager"
	"github.com/pulumi/pulumi-gcp/sdk/v9/go/gcp/projects"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

// certificateMap provisions the certificate map and its entry fan-out.
// The map (and its entries) is a GLOBAL Certificate Manager resource —
// there is no location argument by API design.
//
// Entries are almost entirely IMMUTABLE (hostname, matcher, entry name,
// project all replace on change); the certificate LIST is the mutable
// part — certificate rotation edits the list in place while the entry
// keeps serving.
func certificateMap(ctx *pulumi.Context, locals *Locals, gcpProvider *gcp.Provider) error {
	spec := locals.GcpCertificateMap.Spec

	// Enable the Certificate Manager API so a fresh project can host the
	// map. disable_on_destroy stays false (the provider default): tearing
	// down one map must never disable Certificate Manager for everything
	// else in the project. Matches the Terraform module.
	serviceArgs := &projects.ServiceArgs{
		Service:                  pulumi.String("certificatemanager.googleapis.com"),
		DisableDependentServices: pulumi.BoolPtr(true),
	}
	if spec.ProjectId.GetValue() != "" {
		serviceArgs.Project = pulumi.String(spec.ProjectId.GetValue())
	}
	createdProjectService, err := projects.NewService(ctx,
		"certificatemap-certificatemanager.googleapis.com", serviceArgs, pulumi.Provider(gcpProvider))
	if err != nil {
		return errors.Wrap(err, "failed to enable certificatemanager.googleapis.com api")
	}

	// The bridge names this resource CertificateMapResource (the plain
	// CertificateMap name is taken by its data source).
	mapArgs := &certificatemanager.CertificateMapResourceArgs{
		Name:   pulumi.String(locals.MapName),
		Labels: pulumi.ToStringMap(locals.GcpLabels),
	}
	// Honor the spec contract: an empty project_id falls back to the
	// provider's default project.
	if spec.ProjectId.GetValue() != "" {
		mapArgs.Project = pulumi.String(spec.ProjectId.GetValue())
	}
	if spec.Description != "" {
		mapArgs.Description = pulumi.StringPtr(spec.Description)
	}
	// Unset defers to the provider default (DELETE). The same value is
	// wired to every entry below — one spec lever, every resource.
	if spec.DeletionPolicy != "" {
		mapArgs.DeletionPolicy = pulumi.String(spec.DeletionPolicy)
	}

	createdMap, err := certificatemanager.NewCertificateMapResource(ctx, "certificate-map", mapArgs,
		pulumi.Provider(gcpProvider), pulumi.DependsOn([]pulumi.Resource{createdProjectService}))
	if err != nil {
		return errors.Wrap(err, "failed to create certificate map")
	}

	for _, entry := range spec.Entries {
		certificates := pulumi.StringArray{}
		for _, certificate := range entry.Certificates {
			certificates = append(certificates, pulumi.String(certificate.GetValue()))
		}

		entryArgs := &certificatemanager.CertificateMapEntryArgs{
			Name:         pulumi.String(entry.EntryName),
			Map:          createdMap.Name,
			Certificates: certificates,
			Labels:       pulumi.ToStringMap(entryLabels(locals.GcpLabels, entry.Labels)),
		}
		if spec.ProjectId.GetValue() != "" {
			entryArgs.Project = pulumi.String(spec.ProjectId.GetValue())
		}
		// Exactly one of hostname / matcher (proto-CEL-enforced, mirroring
		// the provider's ExactlyOneOf).
		if entry.Hostname != "" {
			entryArgs.Hostname = pulumi.StringPtr(entry.Hostname)
		}
		if entry.Matcher != "" {
			entryArgs.Matcher = pulumi.StringPtr(entry.Matcher)
		}
		if entry.Description != "" {
			entryArgs.Description = pulumi.StringPtr(entry.Description)
		}
		if spec.DeletionPolicy != "" {
			entryArgs.DeletionPolicy = pulumi.String(spec.DeletionPolicy)
		}
		if _, err := certificatemanager.NewCertificateMapEntry(ctx, "entry-"+entry.EntryName, entryArgs,
			pulumi.Provider(gcpProvider)); err != nil {
			return errors.Wrapf(err, "failed to create certificate map entry %s", entry.EntryName)
		}
	}

	// The full resource name and the //certificatemanager.googleapis.com/
	// form a GcpTargetHttpsProxy's certificate_map argument consumes —
	// assembled from the map's own computed project attribute so the
	// ambient-project case renders correctly.
	mapId := pulumi.Sprintf("projects/%s/locations/global/certificateMaps/%s",
		createdMap.Project, createdMap.Name)
	ctx.Export(OpMapId, mapId)
	ctx.Export(OpMapUri, pulumi.Sprintf("//certificatemanager.googleapis.com/%s", mapId))
	ctx.Export(OpMapName, createdMap.Name)

	return nil
}
