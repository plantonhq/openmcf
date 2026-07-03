package module

import (
	"github.com/pkg/errors"
	"github.com/pulumi/pulumi-gcp/sdk/v9/go/gcp"
	"github.com/pulumi/pulumi-gcp/sdk/v9/go/gcp/iam"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

// workloadIdentityPoolProvider provisions one external issuer inside a
// Workload Identity Pool — the piece that makes keyless federation work.
//
// workload_identity_pool_id, workload_identity_pool_provider_id, and project
// are immutable (ForceNew): changing any of them destroys and recreates the
// provider, invalidating tokens minted for the old audience. The issuer arm's
// contents, attribute_mapping, attribute_condition, display_name,
// description, and disabled all update in place. The issuer TYPE (aws vs oidc
// vs saml vs x509) cannot change on a live provider — the API rejects
// cross-type updates, so switching issuers means a new provider.
//
// GCP soft-deletes providers: after destroy, the provider remains DELETED for
// ~30 days and its ID cannot be reused until permanent deletion (no
// undelete-on-create). Prefer `disabled` for temporary shutoffs.
func workloadIdentityPoolProvider(ctx *pulumi.Context, locals *Locals, gcpProvider *gcp.Provider) error {
	spec := locals.GcpWorkloadIdentityPoolProvider.Spec

	args := &iam.WorkloadIdentityPoolProviderArgs{
		// The pool reference arrives resolved to the bare pool ID (the FK
		// resolver dereferences GcpWorkloadIdentityPool outputs before the
		// module runs).
		WorkloadIdentityPoolId:         pulumi.String(spec.WorkloadIdentityPoolId.GetValue()),
		WorkloadIdentityPoolProviderId: pulumi.String(spec.WorkloadIdentityPoolProviderId),
		// Disabled is always sent: it is the documented kill switch, and an
		// explicit false keeps re-enable flows (true -> false) working.
		Disabled: pulumi.Bool(spec.Disabled),
	}

	// Omitted optionals stay unset (matching the Terraform module's null)
	// rather than being sent as empty strings.
	if spec.DisplayName != "" {
		args.DisplayName = pulumi.String(spec.DisplayName)
	}
	if spec.Description != "" {
		args.Description = pulumi.String(spec.Description)
	}
	if spec.AttributeCondition != "" {
		args.AttributeCondition = pulumi.String(spec.AttributeCondition)
	}

	// AWS, SAML, and X.509 issuers have server-side default mappings, so an
	// empty map stays unset; OIDC requires an explicit mapping, which the spec
	// validation already guarantees carries google.subject.
	if len(spec.AttributeMapping) > 0 {
		mapping := pulumi.StringMap{}
		for attribute, expression := range spec.AttributeMapping {
			mapping[attribute] = pulumi.String(expression)
		}
		args.AttributeMapping = mapping
	}

	// Honor the spec contract: an empty project_id falls back to the provider's
	// default project. Leaving Project unset lets the gcp provider resolve its
	// own project (configuration or the GOOGLE_PROJECT / GOOGLE_CLOUD_PROJECT
	// environment chain); an empty string would be sent verbatim and rejected.
	if spec.ProjectId.GetValue() != "" {
		args.Project = pulumi.String(spec.ProjectId.GetValue())
	}

	// Exactly one issuer arm is set (enforced by the spec's oneof); the API
	// enforces the same exclusivity server-side.
	switch {
	case spec.GetAws() != nil:
		args.Aws = &iam.WorkloadIdentityPoolProviderAwsArgs{
			AccountId: pulumi.String(spec.GetAws().AccountId),
		}

	case spec.GetOidc() != nil:
		oidc := spec.GetOidc()
		oidcArgs := &iam.WorkloadIdentityPoolProviderOidcArgs{
			IssuerUri: pulumi.String(oidc.IssuerUri),
		}
		// An empty audience list means "audience must equal the provider's own
		// canonical resource name" — the safest default; only send overrides.
		if len(oidc.AllowedAudiences) > 0 {
			audiences := pulumi.StringArray{}
			for _, audience := range oidc.AllowedAudiences {
				audiences = append(audiences, pulumi.String(audience))
			}
			oidcArgs.AllowedAudiences = audiences
		}
		// Unset JWKS means keys are fetched from the issuer's .well-known
		// discovery document — the normal path for public issuers.
		if oidc.JwksJson != "" {
			oidcArgs.JwksJson = pulumi.String(oidc.JwksJson)
		}
		args.Oidc = oidcArgs

	case spec.GetSaml() != nil:
		args.Saml = &iam.WorkloadIdentityPoolProviderSamlArgs{
			IdpMetadataXml: pulumi.String(spec.GetSaml().IdpMetadataXml),
		}

	case spec.GetX509() != nil:
		trustStore := spec.GetX509().TrustStore
		anchors := iam.WorkloadIdentityPoolProviderX509TrustStoreTrustAnchorArray{}
		for _, anchor := range trustStore.TrustAnchors {
			anchors = append(anchors, &iam.WorkloadIdentityPoolProviderX509TrustStoreTrustAnchorArgs{
				PemCertificate: pulumi.String(anchor.PemCertificate),
			})
		}
		trustStoreArgs := &iam.WorkloadIdentityPoolProviderX509TrustStoreArgs{
			TrustAnchors: anchors,
		}
		if len(trustStore.IntermediateCas) > 0 {
			intermediates := iam.WorkloadIdentityPoolProviderX509TrustStoreIntermediateCaArray{}
			for _, intermediate := range trustStore.IntermediateCas {
				intermediates = append(intermediates, &iam.WorkloadIdentityPoolProviderX509TrustStoreIntermediateCaArgs{
					PemCertificate: pulumi.String(intermediate.PemCertificate),
				})
			}
			trustStoreArgs.IntermediateCas = intermediates
		}
		args.X509 = &iam.WorkloadIdentityPoolProviderX509Args{
			TrustStore: trustStoreArgs,
		}

	default:
		return errors.New("exactly one issuer (aws, oidc, saml, or x509) must be configured")
	}

	createdProvider, err := iam.NewWorkloadIdentityPoolProvider(ctx, "workload-identity-pool-provider", args, pulumi.Provider(gcpProvider))
	if err != nil {
		return errors.Wrap(err, "failed to create workload identity pool provider")
	}

	ctx.Export(OpName, createdProvider.Name)
	ctx.Export(OpWorkloadIdentityPoolProviderId, createdProvider.WorkloadIdentityPoolProviderId)
	ctx.Export(OpState, createdProvider.State)

	return nil
}
