package module

import (
	"github.com/pkg/errors"
	"github.com/pulumi/pulumi-gcp/sdk/v9/go/gcp"
	"github.com/pulumi/pulumi-gcp/sdk/v9/go/gcp/cloudrun"
	"github.com/pulumi/pulumi-gcp/sdk/v9/go/gcp/organizations"
	"github.com/pulumi/pulumi-gcp/sdk/v9/go/gcp/projects"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

// domainMapping provisions the Cloud Run domain mapping.
//
// The underlying resource is fully IMMUTABLE (every argument is
// create-only), so this module has no update semantics to manage — any
// spec change replaces the mapping, which GCP re-creates in seconds.
//
// The Cloud Run v1 API requires a metadata block whose namespace equals
// the project ID or project number, so the module always renders one:
// spec.namespace when set, else the spec's project, else the provider's
// resolved default project (read from the client config only in that last
// case — the count-gated ambient-project pattern).
func domainMapping(ctx *pulumi.Context, locals *Locals, gcpProvider *gcp.Provider) error {
	spec := locals.GcpCloudRunDomainMapping.Spec

	// Enable the Cloud Run Admin API, which serves domain mappings.
	// disable_on_destroy stays false: tearing down one mapping must never
	// disable Cloud Run for everything else in the project. Matches the
	// Terraform module.
	serviceArgs := &projects.ServiceArgs{
		Service:                  pulumi.String("run.googleapis.com"),
		DisableDependentServices: pulumi.BoolPtr(true),
	}
	if spec.ProjectId.GetValue() != "" {
		serviceArgs.Project = pulumi.String(spec.ProjectId.GetValue())
	}
	createdProjectService, err := projects.NewService(ctx,
		"mapping-run.googleapis.com", serviceArgs, pulumi.Provider(gcpProvider))
	if err != nil {
		return errors.Wrap(err, "failed to enable run.googleapis.com api")
	}

	// The namespace fallback chain: spec.namespace → spec project →
	// provider default project (client-config read gated to the one case
	// that needs it, so ordinary project-named deployments stay
	// credential-free offline).
	namespace := spec.Namespace
	if namespace == "" {
		namespace = spec.ProjectId.GetValue()
	}
	if namespace == "" {
		clientConfig, err := organizations.GetClientConfig(ctx, pulumi.Provider(gcpProvider))
		if err != nil {
			return errors.Wrap(err, "failed to read provider client config for the default project")
		}
		if clientConfig.Project == "" {
			return errors.New("namespace and project_id are empty and the provider has no default project configured")
		}
		namespace = clientConfig.Project
	}

	metadataArgs := &cloudrun.DomainMappingMetadataArgs{
		Namespace: pulumi.String(namespace),
		Labels:    pulumi.ToStringMap(locals.GcpLabels),
	}
	if len(spec.Annotations) > 0 {
		metadataArgs.Annotations = pulumi.ToStringMap(spec.Annotations)
	}

	specArgs := &cloudrun.DomainMappingSpecArgs{
		RouteName: pulumi.String(spec.Route.GetValue()),
	}
	if spec.CertificateMode != "" {
		specArgs.CertificateMode = pulumi.StringPtr(spec.CertificateMode)
	}
	if spec.ForceOverride {
		specArgs.ForceOverride = pulumi.BoolPtr(true)
	}

	args := &cloudrun.DomainMappingArgs{
		// The mapping's name IS the domain being mapped.
		Name:     pulumi.StringPtr(spec.Domain),
		Location: pulumi.String(spec.Region),
		Spec:     specArgs,
		Metadata: metadataArgs,
	}
	if spec.ProjectId.GetValue() != "" {
		args.Project = pulumi.String(spec.ProjectId.GetValue())
	}
	if spec.DeletionPolicy != "" {
		args.DeletionPolicy = pulumi.StringPtr(spec.DeletionPolicy)
	}

	createdDomainMapping, err := cloudrun.NewDomainMapping(ctx, "domain-mapping", args,
		pulumi.Provider(gcpProvider), pulumi.DependsOn([]pulumi.Resource{createdProjectService}))
	if err != nil {
		return errors.Wrap(err, "failed to create domain mapping")
	}

	ctx.Export(OpDomain, createdDomainMapping.Name)
	ctx.Export(OpRegion, createdDomainMapping.Location)

	// The record count is SERVER-decided (a root domain receives A/AAAA
	// sets, a subdomain one CNAME), so the records export as ONE
	// structured list output — registered synchronously here; only the
	// VALUE resolves inside ApplyT (never ctx.Export itself, which races
	// the engine's output marshaling). The outputs transformer flattens
	// the list onto the repeated proto message by dot-indexed keys.
	ctx.Export(OpResourceRecords, createdDomainMapping.Statuses.ApplyT(
		func(statuses []cloudrun.DomainMappingStatus) []map[string]string {
			records := []map[string]string{}
			if len(statuses) == 0 {
				return records
			}
			for _, record := range statuses[0].ResourceRecords {
				entry := map[string]string{
					"record_type": "",
					"record_name": "",
					"rrdata":      "",
				}
				if record.Type != nil {
					entry["record_type"] = *record.Type
				}
				if record.Name != nil {
					entry["record_name"] = *record.Name
				}
				if record.Rrdata != nil {
					entry["rrdata"] = *record.Rrdata
				}
				records = append(records, entry)
			}
			return records
		}))
	ctx.Export(OpMappedRouteName, createdDomainMapping.Statuses.ApplyT(
		func(statuses []cloudrun.DomainMappingStatus) string {
			if len(statuses) > 0 && statuses[0].MappedRouteName != nil {
				return *statuses[0].MappedRouteName
			}
			return ""
		}).(pulumi.StringOutput))

	return nil
}
