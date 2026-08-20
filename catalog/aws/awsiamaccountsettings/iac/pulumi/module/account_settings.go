package module

import (
	"github.com/pkg/errors"
	"github.com/pulumi/pulumi-aws/sdk/v7/go/aws"
	"github.com/pulumi/pulumi-aws/sdk/v7/go/aws/iam"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

// accountSettings manages IAM's account-level settings - the sign-in
// alias, the password policy, and the STS global-endpoint token
// version - and exports outputs.
//
// Lifecycle facts the render below depends on:
//   - each arm renders ONLY when its spec field/message is present (an
//     omitted arm leaves the account's current setting untouched -
//     that omission is meaningful and deliberate);
//   - destroy semantics DIFFER per arm: the alias truly DELETES
//     (sign-in URLs revert to the bare account ID), the password
//     policy RESETS to AWS's defaults, and the STS preference's delete
//     is a NO-OP (the last-applied token version persists; reverting
//     is an apply with the other version);
//   - the password policy is replaced WHOLE on every apply (AWS's
//     UpdateAccountPasswordPolicy semantics): an unset field is AWS's
//     default, never "keep the current setting". The provider also
//     never sends false/0 values on the wire (they are
//     indistinguishable from unset at the SDK layer) - AWS treats a
//     missing toggle as false, so the rendered result is identical;
//   - none of these resources is taggable at AWS.
func accountSettings(ctx *pulumi.Context, locals *Locals, provider *aws.Provider) error {
	spec := locals.Spec

	if spec.AccountAlias != "" {
		createdAlias, err := iam.NewAccountAlias(ctx, "account-alias", &iam.AccountAliasArgs{
			AccountAlias: pulumi.String(spec.AccountAlias),
		}, pulumi.Provider(provider))
		if err != nil {
			return errors.Wrap(err, "put account alias")
		}
		ctx.Export(OpAccountAlias, createdAlias.AccountAlias)
	} else {
		ctx.Export(OpAccountAlias, pulumi.String(""))
	}

	if policy := spec.PasswordPolicy; policy != nil {
		args := &iam.AccountPasswordPolicyArgs{
			// Plain-bool toggles render unconditionally: false and
			// unset are the same posture at AWS (the policy is
			// replaced whole), and rendering both engines identically
			// keeps state parity.
			RequireLowercaseCharacters: pulumi.BoolPtr(policy.RequireLowercaseCharacters),
			RequireNumbers:             pulumi.BoolPtr(policy.RequireNumbers),
			RequireSymbols:             pulumi.BoolPtr(policy.RequireSymbols),
			RequireUppercaseCharacters: pulumi.BoolPtr(policy.RequireUppercaseCharacters),
			HardExpiry:                 pulumi.BoolPtr(policy.HardExpiry),
		}
		// Presence-typed knobs render only when set - their AWS
		// defaults (6-character minimum, self-service changes allowed,
		// no expiry, no reuse prevention) apply otherwise.
		if policy.MinimumPasswordLength != nil {
			args.MinimumPasswordLength = pulumi.IntPtr(int(policy.GetMinimumPasswordLength()))
		}
		if policy.AllowUsersToChangePassword != nil {
			args.AllowUsersToChangePassword = pulumi.BoolPtr(policy.GetAllowUsersToChangePassword())
		}
		if policy.MaxPasswordAge != nil {
			args.MaxPasswordAge = pulumi.IntPtr(int(policy.GetMaxPasswordAge()))
		}
		if policy.PasswordReusePrevention != nil {
			args.PasswordReusePrevention = pulumi.IntPtr(int(policy.GetPasswordReusePrevention()))
		}

		createdPolicy, err := iam.NewAccountPasswordPolicy(ctx, "password-policy", args, pulumi.Provider(provider))
		if err != nil {
			return errors.Wrap(err, "put account password policy")
		}
		// AWS derives expire_passwords from max_password_age; exported
		// as a string to match the Terraform module key-for-key.
		ctx.Export(OpExpirePasswords, pulumi.Sprintf("%t", createdPolicy.ExpirePasswords))
	} else {
		ctx.Export(OpExpirePasswords, pulumi.String(""))
	}

	if spec.Sts != nil {
		if _, err := iam.NewSecurityTokenServicePreferences(ctx, "sts-preferences",
			&iam.SecurityTokenServicePreferencesArgs{
				GlobalEndpointTokenVersion: pulumi.String(spec.Sts.GlobalEndpointTokenVersion),
			}, pulumi.Provider(provider)); err != nil {
			return errors.Wrap(err, "put sts preferences")
		}
	}

	// The account id feeds the output regardless of which arms render
	// (every resource here is account-scoped).
	callerIdentity, err := aws.GetCallerIdentity(ctx, &aws.GetCallerIdentityArgs{}, pulumi.Provider(provider))
	if err != nil {
		return errors.Wrap(err, "get caller identity")
	}
	ctx.Export(OpAccountId, pulumi.String(callerIdentity.AccountId))

	return nil
}
