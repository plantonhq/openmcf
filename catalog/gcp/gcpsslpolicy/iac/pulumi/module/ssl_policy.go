package module

import (
	"github.com/pkg/errors"
	"github.com/pulumi/pulumi-gcp/sdk/v9/go/gcp"
	"github.com/pulumi/pulumi-gcp/sdk/v9/go/gcp/compute"
	"github.com/pulumi/pulumi-gcp/sdk/v9/go/gcp/projects"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

// sslPolicy provisions the Compute Engine SSL policy — the control for which
// TLS versions and cipher suites a load balancer accepts from clients.
// Target HTTPS (and SSL) proxies reference the policy's self_link; without
// one, GCP's permissive default applies (min TLS 1.0, COMPATIBLE ciphers).
//
// One kind, two provider resources: GCP models global and regional SSL
// policies as separate API collections with an identical surface, so an
// empty spec.region creates compute.SSLPolicy and a set region creates
// compute.RegionSslPolicy — mirroring the Terraform module's count guards.
//
// name, project, and description are immutable (ForceNew in the provider):
// changing any of them destroys and recreates the policy, briefly breaking
// every proxy that references the old self_link. profile, min_tls_version,
// and custom_features all update IN PLACE — a policy is shared
// configuration, so tightening the TLS floor for an entire proxy fleet is a
// single-resource change that applies on the next client handshake.
//
// profile and min_tls_version deliberately fall through to the API defaults
// (COMPATIBLE / TLS_1_0) when unset — hardcoding them here would silently
// pin behavior the provider may evolve.
func sslPolicy(ctx *pulumi.Context, locals *Locals, gcpProvider *gcp.Provider) error {
	spec := locals.GcpSslPolicy.Spec

	// Enable the Compute Engine API so a fresh project can host the SSL
	// policy. disable_on_destroy stays false (the provider default): tearing
	// down one policy must never disable the API for everything else in the
	// project. Matches the Terraform module.
	serviceArgs := &projects.ServiceArgs{
		Service:                  pulumi.String("compute.googleapis.com"),
		DisableDependentServices: pulumi.BoolPtr(true),
	}
	if spec.ProjectId.GetValue() != "" {
		serviceArgs.Project = pulumi.String(spec.ProjectId.GetValue())
	}
	createdProjectService, err := projects.NewService(ctx,
		"sslpolicy-compute.googleapis.com", serviceArgs, pulumi.Provider(gcpProvider))
	if err != nil {
		return errors.Wrap(err, "failed to enable compute.googleapis.com api")
	}

	if spec.Region == "" {
		return globalSslPolicy(ctx, locals, gcpProvider, createdProjectService)
	}
	return regionalSslPolicy(ctx, locals, gcpProvider, createdProjectService)
}

func globalSslPolicy(ctx *pulumi.Context, locals *Locals, gcpProvider *gcp.Provider, computeApiService pulumi.Resource) error {
	spec := locals.GcpSslPolicy.Spec

	args := &compute.SSLPolicyArgs{
		Name: pulumi.String(locals.SslPolicyName),
	}

	// Omitted enum-like strings stay unset (matching the Terraform module's
	// null) so GCP applies its own defaults (COMPATIBLE / TLS_1_0) rather
	// than receiving empty strings it would reject.
	if spec.Profile != "" {
		args.Profile = pulumi.String(spec.Profile)
	}
	if spec.MinTlsVersion != "" {
		args.MinTlsVersion = pulumi.String(spec.MinTlsVersion)
	}
	if spec.Description != "" {
		args.Description = pulumi.String(spec.Description)
	}

	// Only sent with the CUSTOM profile (the proto CEL enforces the
	// pairing); GCP rejects the field on predefined profiles.
	if len(spec.CustomFeatures) > 0 {
		args.CustomFeatures = pulumi.ToStringArray(spec.CustomFeatures)
	}

	// Post-quantum rollout stance; unset falls through to the API default
	// (DEFAULT — GCP's own timeline). Matches the Terraform module's null.
	if spec.PostQuantumKeyExchange != "" {
		args.PostQuantumKeyExchange = pulumi.String(spec.PostQuantumKeyExchange)
	}

	// Unset defers to the provider default (DELETE).
	if spec.DeletionPolicy != "" {
		args.DeletionPolicy = pulumi.String(spec.DeletionPolicy)
	}

	// Honor the spec contract: an empty project_id falls back to the provider's
	// default project. Leaving Project unset lets the gcp provider resolve its
	// own project (configuration or the GOOGLE_PROJECT / GOOGLE_CLOUD_PROJECT
	// environment chain); an empty string would be sent verbatim and rejected.
	if spec.ProjectId.GetValue() != "" {
		args.Project = pulumi.String(spec.ProjectId.GetValue())
	}

	createdSslPolicy, err := compute.NewSSLPolicy(ctx, "ssl-policy", args,
		pulumi.Provider(gcpProvider), pulumi.DependsOn([]pulumi.Resource{computeApiService}))
	if err != nil {
		return errors.Wrap(err, "failed to create global ssl policy")
	}

	ctx.Export(OpSelfLink, createdSslPolicy.SelfLink)
	ctx.Export(OpSslPolicyName, createdSslPolicy.Name)
	ctx.Export(OpEnabledFeatures, createdSslPolicy.EnabledFeatures)
	// Empty region marks the global scope for downstream composition checks.
	ctx.Export(OpRegion, pulumi.String(""))

	return nil
}

func regionalSslPolicy(ctx *pulumi.Context, locals *Locals, gcpProvider *gcp.Provider, computeApiService pulumi.Resource) error {
	spec := locals.GcpSslPolicy.Spec

	args := &compute.RegionSslPolicyArgs{
		Name:   pulumi.String(locals.SslPolicyName),
		Region: pulumi.String(spec.Region),
	}

	if spec.Profile != "" {
		args.Profile = pulumi.String(spec.Profile)
	}
	if spec.MinTlsVersion != "" {
		args.MinTlsVersion = pulumi.String(spec.MinTlsVersion)
	}
	if spec.Description != "" {
		args.Description = pulumi.String(spec.Description)
	}

	if len(spec.CustomFeatures) > 0 {
		args.CustomFeatures = pulumi.ToStringArray(spec.CustomFeatures)
	}

	if spec.PostQuantumKeyExchange != "" {
		args.PostQuantumKeyExchange = pulumi.String(spec.PostQuantumKeyExchange)
	}

	// Unset defers to the provider default (DELETE).
	if spec.DeletionPolicy != "" {
		args.DeletionPolicy = pulumi.String(spec.DeletionPolicy)
	}

	if spec.ProjectId.GetValue() != "" {
		args.Project = pulumi.String(spec.ProjectId.GetValue())
	}

	createdSslPolicy, err := compute.NewRegionSslPolicy(ctx, "ssl-policy", args,
		pulumi.Provider(gcpProvider), pulumi.DependsOn([]pulumi.Resource{computeApiService}))
	if err != nil {
		return errors.Wrap(err, "failed to create regional ssl policy")
	}

	ctx.Export(OpSelfLink, createdSslPolicy.SelfLink)
	ctx.Export(OpSslPolicyName, createdSslPolicy.Name)
	ctx.Export(OpEnabledFeatures, createdSslPolicy.EnabledFeatures)
	// Export the plain region NAME from the spec (matching the Terraform
	// module) — the provider's region attribute can carry a self-link, which
	// API callers and verification cannot use directly.
	ctx.Export(OpRegion, pulumi.String(spec.Region))

	return nil
}
