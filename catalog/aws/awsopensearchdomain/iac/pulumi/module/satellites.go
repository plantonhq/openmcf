package module

import (
	"fmt"

	"github.com/pkg/errors"
	"github.com/pulumi/pulumi-aws/sdk/v7/go/aws"
	"github.com/pulumi/pulumi-aws/sdk/v7/go/aws/opensearch"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

// satellites provisions the domain-scoped companions, each its own
// provider resource keyed off the domain: SAML authentication for
// Dashboards and the grantor side of cross-account VPC endpoint access.
// Both attach to exactly this domain and are meaningless without it,
// which is why they fold here instead of being their own kinds.
func satellites(ctx *pulumi.Context, locals *Locals, provider *aws.Provider,
	osDomain *opensearch.Domain) error {
	spec := locals.Spec

	// SAML sign-in for OpenSearch Dashboards. The block's presence is the
	// enable switch: removing spec.saml_options destroys this resource,
	// which disables SAML on the domain (the provider's delete calls the
	// disable API). Rides on fine-grained access control (CEL-enforced).
	if saml := spec.SamlOptions; saml != nil {
		samlArgs := &opensearch.DomainSamlOptionsSamlOptionsArgs{
			Enabled: pulumi.BoolPtr(true),
			Idp: &opensearch.DomainSamlOptionsSamlOptionsIdpArgs{
				EntityId:        pulumi.String(saml.IdpEntityId),
				MetadataContent: pulumi.String(saml.IdpMetadataContent),
			},
		}
		if saml.MasterBackendRole != "" {
			samlArgs.MasterBackendRole = pulumi.StringPtr(saml.MasterBackendRole)
		}
		if saml.MasterUserName != "" {
			samlArgs.MasterUserName = pulumi.StringPtr(saml.MasterUserName)
		}
		if saml.RolesKey != "" {
			samlArgs.RolesKey = pulumi.StringPtr(saml.RolesKey)
		}
		if saml.SubjectKey != "" {
			samlArgs.SubjectKey = pulumi.StringPtr(saml.SubjectKey)
		}
		// 0 keeps the AWS default (60 minutes).
		if saml.SessionTimeoutMinutes > 0 {
			samlArgs.SessionTimeoutMinutes = pulumi.IntPtr(int(saml.SessionTimeoutMinutes))
		}

		if _, err := opensearch.NewDomainSamlOptions(ctx, "saml-options",
			&opensearch.DomainSamlOptionsArgs{
				DomainName:  osDomain.DomainName,
				SamlOptions: samlArgs,
			}, pulumi.Provider(provider), pulumi.Parent(osDomain)); err != nil {
			return errors.Wrap(err, "failed to create SAML options")
		}
	}

	// Cross-account private access, grantor side: each listed account is
	// authorized to create OpenSearch-managed VPC endpoints against this
	// domain. One resource per account, keyed by the account ID, so
	// grants come and go independently.
	for _, account := range spec.AuthorizedVpcEndpointAccessAccounts {
		if _, err := opensearch.NewAuthorizeVpcEndpointAccess(ctx,
			fmt.Sprintf("authorize-vpc-endpoint-access-%s", account),
			&opensearch.AuthorizeVpcEndpointAccessArgs{
				DomainName: osDomain.DomainName,
				Account:    pulumi.String(account),
			}, pulumi.Provider(provider), pulumi.Parent(osDomain)); err != nil {
			return errors.Wrapf(err, "failed to authorize VPC endpoint access for account %s", account)
		}
	}

	return nil
}
