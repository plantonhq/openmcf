package module

import (
	"github.com/pkg/errors"
	"github.com/pulumi/pulumi-aws/sdk/v7/go/aws/elasticache"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

// user provisions the aws_elasticache_user. The modeled surface maps
// one-to-one onto AWS's CreateUser call: identity (user id + user name +
// engine), the ACL access string, and exactly one authentication mode.
func user(ctx *pulumi.Context, locals *Locals, provider pulumi.ProviderResource) error {
	spec := locals.AwsElasticacheUser.Spec

	// The nested authentication_mode block is the single authentication
	// surface (the provider's legacy flat passwords/no_password_required
	// arms model the same capability and are deliberately not used --
	// one honest shape). Passwords are only present for the "password"
	// type; CEL guarantees the other types carry none.
	authMode := &elasticache.UserAuthenticationModeArgs{
		Type: pulumi.String(spec.AuthenticationMode.Type),
	}
	if len(spec.AuthenticationMode.Passwords) > 0 {
		authMode.Passwords = pulumi.ToStringArray(spec.AuthenticationMode.Passwords)
	}

	created, err := elasticache.NewUser(ctx, locals.UserId, &elasticache.UserArgs{
		// The AWS user id is create-time immutable and doubles as the
		// Pulumi resource name -- metadata.name on both engines.
		UserId: pulumi.String(locals.UserId),
		// user_name is what clients present in AUTH; it is NOT unique
		// per user (AWS unions credentials of same-named users), which
		// is why it is a spec field instead of reusing metadata.name.
		UserName:           pulumi.String(spec.UserName),
		Engine:             pulumi.String(spec.Engine),
		AccessString:       pulumi.String(spec.AccessString),
		AuthenticationMode: authMode,
		Tags:               pulumi.ToStringMap(locals.AwsTags),
	}, pulumi.Provider(provider))
	if err != nil {
		return errors.Wrap(err, "failed to create user")
	}

	ctx.Export(OpUserId, created.UserId)
	ctx.Export(OpArn, created.Arn)
	ctx.Export(OpUserName, created.UserName)

	return nil
}
