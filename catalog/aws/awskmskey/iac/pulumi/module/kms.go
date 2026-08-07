package module

import (
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
