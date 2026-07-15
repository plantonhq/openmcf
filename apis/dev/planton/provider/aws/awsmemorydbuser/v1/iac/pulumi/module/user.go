package module

import (
	"github.com/pkg/errors"
	"github.com/pulumi/pulumi-aws/sdk/v7/go/aws/memorydb"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

// user provisions the aws_memorydb_user. The modeled surface maps one-to-one
// onto AWS's CreateUser call: the name (the AUTH identity), the ACL access
// string, and exactly one authentication mode.
func user(ctx *pulumi.Context, locals *Locals, provider pulumi.ProviderResource) error {
	spec := locals.AwsMemorydbUser.Spec

	// Exactly two input types exist ("password" carries 1-2 secrets; "iam"
	// carries none) -- CEL guarantees the coupling, so presence of passwords
	// is the only branch the module needs.
	authMode := &memorydb.UserAuthenticationModeArgs{
		Type: pulumi.String(spec.AuthenticationMode.Type),
	}
	if len(spec.AuthenticationMode.Passwords) > 0 {
		authMode.Passwords = pulumi.ToStringArray(spec.AuthenticationMode.Passwords)
	}

	created, err := memorydb.NewUser(ctx, locals.UserName, &memorydb.UserArgs{
		// The user name is the user's single AWS identity, create-time
		// immutable, and doubles as the Pulumi resource name --
		// metadata.name on both engines.
		UserName:           pulumi.String(locals.UserName),
		AccessString:       pulumi.String(spec.AccessString),
		AuthenticationMode: authMode,
		Tags:               pulumi.ToStringMap(locals.AwsTags),
	}, pulumi.Provider(provider))
	if err != nil {
		return errors.Wrap(err, "failed to create user")
	}

	ctx.Export(OpUserName, created.UserName)
	ctx.Export(OpUserArn, created.Arn)
	ctx.Export(OpMinimumEngineVersion, created.MinimumEngineVersion)

	return nil
}
