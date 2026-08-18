package module

import (
	"github.com/pkg/errors"
	"github.com/pulumi/pulumi-aws/sdk/v7/go/aws/iam"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

// samlProvider provisions the SAML identity provider and exports
// outputs.
//
// Lifecycle facts the render below depends on:
//   - the provider's name comes from metadata.name and is WRITE-ONCE
//     at AWS - a rename replaces the provider and invalidates every
//     role trust policy naming its ARN;
//   - the metadata document updates IN PLACE - certificate rotations
//     are ordinary updates, and valid_until (exported below) is the
//     date to rotate by;
//   - IAM is global: the provider exists account-wide regardless of
//     the endpoint region the stack ran against.
func samlProvider(ctx *pulumi.Context, locals *Locals, provider pulumi.ProviderResource) error {
	providerName := locals.Target.Metadata.Name
	spec := locals.Spec

	createdSamlProvider, err := iam.NewSamlProvider(ctx, providerName, &iam.SamlProviderArgs{
		Name:                 pulumi.String(providerName),
		SamlMetadataDocument: pulumi.String(spec.SamlMetadataDocument),
		Tags:                 pulumi.ToStringMap(locals.AwsTags),
	}, pulumi.Provider(provider))
	if err != nil {
		return errors.Wrap(err, "failed to create SAML provider")
	}

	ctx.Export(OpProviderArn, createdSamlProvider.Arn)
	ctx.Export(OpSamlProviderUuid, createdSamlProvider.SamlProviderUuid)
	ctx.Export(OpValidUntil, createdSamlProvider.ValidUntil)

	return nil
}
