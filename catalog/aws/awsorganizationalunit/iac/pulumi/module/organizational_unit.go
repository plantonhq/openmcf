package module

import (
	"github.com/pkg/errors"
	"github.com/pulumi/pulumi-aws/sdk/v7/go/aws"
	"github.com/pulumi/pulumi-aws/sdk/v7/go/aws/organizations"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

// organizationalUnit creates the OU and exports outputs.
//
// Lifecycle facts the render below depends on:
//   - the parent reference is immutable (AWS moves accounts between
//     OUs, never OUs themselves) - a parent change forces replacement;
//   - the display name renames in place;
//   - creation retries through the organization's finalization window
//     (the provider handles FinalizingOrganizationException for up to
//     four minutes after CreateOrganization);
//   - AWS identifies the OU as "ou-..." (the import ID).
func organizationalUnit(ctx *pulumi.Context, locals *Locals, provider *aws.Provider) error {
	spec := locals.Spec

	createdOrganizationalUnit, err := organizations.NewOrganizationalUnit(ctx, "organizational-unit",
		&organizations.OrganizationalUnitArgs{
			Name:     pulumi.StringPtr(spec.OuName),
			ParentId: pulumi.String(spec.ParentId.GetValue()),
			Tags:     pulumi.ToStringMap(locals.AwsTags),
		}, pulumi.Provider(provider))
	if err != nil {
		return errors.Wrap(err, "create organizational unit")
	}

	ctx.Export(OpOuId, createdOrganizationalUnit.ID())
	ctx.Export(OpArn, createdOrganizationalUnit.Arn)
	return nil
}
