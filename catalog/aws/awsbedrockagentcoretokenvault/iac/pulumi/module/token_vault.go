package module

import (
	"github.com/pkg/errors"
	"github.com/pulumi/pulumi-aws/sdk/v7/go/aws"
	"github.com/pulumi/pulumi-aws/sdk/v7/go/aws/bedrock"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

// tokenVault sets the KMS key on the targeted token vault and exports
// outputs.
//
// Lifecycle facts the render below depends on:
//   - an unset token_vault_id targets AWS's one default vault (locals
//     resolve "default", mirroring the provider's own default);
//     changing the id replaces the configuration onto another vault;
//   - key_type/kms_key_arn pairing is CEL-enforced upstream of this
//     module (CustomerManagedKey requires the ARN, ServiceManagedKey
//     forbids it), so the conditional below never sends a stray ARN;
//   - destroy is a NO-OP at AWS: the last-applied key setting REMAINS
//     in effect. Reverting to AWS-managed encryption is an APPLY with
//     key_type ServiceManagedKey, never a destroy.
func tokenVault(ctx *pulumi.Context, locals *Locals, provider *aws.Provider) error {
	spec := locals.Spec

	kmsConfiguration := &bedrock.AgentcoreTokenVaultCmkKmsConfigurationArgs{
		KeyType: pulumi.String(spec.KeyType),
	}
	if spec.KmsKeyArn.GetValue() != "" {
		kmsConfiguration.KmsKeyArn = pulumi.String(spec.KmsKeyArn.GetValue())
	}

	created, err := bedrock.NewAgentcoreTokenVaultCmk(ctx, "token-vault-cmk",
		&bedrock.AgentcoreTokenVaultCmkArgs{
			TokenVaultId:     pulumi.String(locals.TokenVaultId),
			KmsConfiguration: kmsConfiguration,
		}, pulumi.Provider(provider))
	if err != nil {
		return errors.Wrap(err, "set token vault cmk")
	}

	ctx.Export(OpTokenVaultId, created.TokenVaultId)
	ctx.Export(OpKeyType, created.KmsConfiguration.KeyType())
	ctx.Export(OpKmsKeyArn, created.KmsConfiguration.KmsKeyArn().Elem())
	return nil
}
