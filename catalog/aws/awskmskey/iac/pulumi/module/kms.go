package module

import (
	"fmt"

	"github.com/pkg/errors"
	"github.com/pulumi/pulumi-aws/sdk/v7/go/aws"
	"github.com/pulumi/pulumi-aws/sdk/v7/go/aws/kms"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

type KmsKeyResult struct {
	KeyId      pulumi.StringOutput
	KeyArn     pulumi.StringOutput
	AliasNames pulumi.StringArrayOutput
}

func kmsKey(ctx *pulumi.Context, locals *Locals, provider *aws.Provider) (*KmsKeyResult, error) {
	spec := locals.AwsKmsKey.Spec

	args := &kms.KeyArgs{
		Tags: pulumi.ToStringMap(locals.AwsTags),
	}

	if spec.Description != "" {
		args.Description = pulumi.String(spec.Description)
	}

	// Empty keeps the AWS default (SYMMETRIC_DEFAULT / ENCRYPT_DECRYPT).
	if spec.KeySpec != "" {
		args.CustomerMasterKeySpec = pulumi.String(spec.KeySpec)
	}
	if spec.KeyUsage != "" {
		args.KeyUsage = pulumi.String(spec.KeyUsage)
	}

	if spec.Policy != "" {
		args.Policy = pulumi.String(spec.Policy)
	}

	if spec.BypassPolicyLockoutSafetyCheck {
		args.BypassPolicyLockoutSafetyCheck = pulumi.Bool(true)
	}

	args.IsEnabled = pulumi.Bool(!spec.Disabled)
	args.EnableKeyRotation = pulumi.Bool(spec.EnableKeyRotation)

	if spec.RotationPeriodInDays != 0 {
		args.RotationPeriodInDays = pulumi.Int(int(spec.RotationPeriodInDays))
	}

	if spec.MultiRegion {
		args.MultiRegion = pulumi.Bool(true)
	}

	if spec.DeletionWindowDays != 0 {
		args.DeletionWindowInDays = pulumi.Int(int(spec.DeletionWindowDays))
	}

	// Custom key store surface: setting the store id makes KMS create the key
	// material in the CloudHSM cluster (or, with xks_key_id, forward
	// operations to the named key in an external key manager). Both
	// create-time immutable.
	if spec.CustomKeyStoreId != "" {
		args.CustomKeyStoreId = pulumi.String(spec.CustomKeyStoreId)
	}
	if spec.XksKeyId != "" {
		args.XksKeyId = pulumi.String(spec.XksKeyId)
	}

	createdKey, err := kms.NewKey(ctx, locals.AwsKmsKey.Metadata.Name, args, pulumi.Provider(provider))
	if err != nil {
		return nil, errors.Wrap(err, "failed to create kms key")
	}

	for _, aliasName := range spec.Aliases {
		_, err := kms.NewAlias(ctx, locals.AwsKmsKey.Metadata.Name+"-"+sanitizeAliasResourceName(aliasName), &kms.AliasArgs{
			Name:        pulumi.String(aliasName),
			TargetKeyId: createdKey.KeyId,
		}, pulumi.Provider(provider))
		if err != nil {
			return nil, errors.Wrapf(err, "failed to create kms alias %s", aliasName)
		}
	}

	// One KMS grant per spec entry: scoped, revocable permissions for a
	// principal to use this key without editing the key policy. Every grant
	// argument is create-time immutable (a change replaces the grant -- safe,
	// grants carry no state). Entries are keyed by list position; grant
	// identity lives in AWS's generated grant id. valueFrom principal
	// references were resolved to ARNs before the module ran.
	for idx, grant := range spec.Grants {
		grantArgs := &kms.GrantArgs{
			KeyId:            createdKey.KeyId,
			GranteePrincipal: pulumi.String(grant.GranteePrincipal.GetValue()),
			Operations:       pulumi.ToStringArray(grant.Operations),
		}
		if grant.Name != "" {
			grantArgs.Name = pulumi.String(grant.Name)
		}
		if grant.RetiringPrincipal.GetValue() != "" {
			grantArgs.RetiringPrincipal = pulumi.String(grant.RetiringPrincipal.GetValue())
		}
		// false REVOKES the grant at teardown (immediate hard stop); true
		// RETIRES it (the graceful path AWS recommends once the grant's work
		// is done).
		if grant.RetireOnDelete {
			grantArgs.RetireOnDelete = pulumi.Bool(true)
		}
		// At most one encryption-context constraint per grant (spec CEL
		// enforces the exclusivity at validate time; the provider only fails
		// it at apply).
		if len(grant.EncryptionContextEquals) > 0 || len(grant.EncryptionContextSubset) > 0 {
			constraint := &kms.GrantConstraintArgs{}
			if len(grant.EncryptionContextEquals) > 0 {
				constraint.EncryptionContextEquals = pulumi.ToStringMap(grant.EncryptionContextEquals)
			}
			if len(grant.EncryptionContextSubset) > 0 {
				constraint.EncryptionContextSubset = pulumi.ToStringMap(grant.EncryptionContextSubset)
			}
			grantArgs.Constraints = kms.GrantConstraintArray{constraint}
		}
		grantName := fmt.Sprintf("%s-grant-%d", locals.AwsKmsKey.Metadata.Name, idx)
		_, err := kms.NewGrant(ctx, grantName, grantArgs,
			pulumi.Provider(provider), pulumi.Parent(createdKey))
		if err != nil {
			return nil, errors.Wrapf(err, "failed to create kms grant %d", idx)
		}
	}

	return &KmsKeyResult{
		KeyId:      createdKey.KeyId,
		KeyArn:     createdKey.Arn,
		AliasNames: pulumi.ToStringArray(spec.Aliases).ToStringArrayOutput(),
	}, nil
}

// sanitizeAliasResourceName strips the alias/ prefix for use in Pulumi logical names.
func sanitizeAliasResourceName(aliasName string) string {
	const prefix = "alias/"
	if len(aliasName) > len(prefix) && aliasName[:len(prefix)] == prefix {
		return aliasName[len(prefix):]
	}
	return aliasName
}
