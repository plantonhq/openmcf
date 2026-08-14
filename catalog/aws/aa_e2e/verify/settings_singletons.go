package verify

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/apigateway"
	"github.com/aws/aws-sdk-go-v2/service/bedrock"
	"github.com/aws/aws-sdk-go-v2/service/bedrockagentcorecontrol"
	"github.com/aws/aws-sdk-go-v2/service/sesv2"
	sesv2types "github.com/aws/aws-sdk-go-v2/service/sesv2/types"
	"github.com/pkg/errors"
)

// The settings-singleton verifiers. These kinds manage ONE
// account/region-scoped settings object whose identity is the region,
// not a name -- and each carries its OWN destroy/absent contract
// (source-verified against the provider at the pin):
//
//   - AwsBedrockInvocationLogging: destroy DELETES the configuration
//     (absent means gone);
//   - AwsApiGatewayAccountSettings: destroy RESETS the CloudWatch role
//     (absent means the role is empty; the account object always
//     exists);
//   - AwsBedrockAgentCoreTokenVault: destroy is a NO-OP (absent means
//     the vault STILL EXISTS with the last-applied setting; reverting
//     is an apply with ServiceManagedKey);
//   - AwsSesAccountSettings: ASYMMETRIC -- VDM resets to DISABLED,
//     suppression persists (the recorded settings-retention class).
//
// There is deliberately no shared "reset" helper: four kinds with
// three different contracts do not justify a new verifier interface.

// apiGatewayAccountSettingsVerifier verifies AwsApiGatewayAccountSettings
// via GetAccount, keyed on account_id (the singleton has no per-name
// identity; the ID output exists for the import recipe and progress
// output).
type apiGatewayAccountSettingsVerifier struct{}

func (*apiGatewayAccountSettingsVerifier) IDOutputKey() string { return "account_id" }

func (*apiGatewayAccountSettingsVerifier) VerifyExists(ctx context.Context, cfg aws.Config, id, region string) error {
	out, err := apigateway.NewFromConfig(cfg, func(o *apigateway.Options) {
		if region != "" {
			o.Region = region
		}
	}).GetAccount(ctx, &apigateway.GetAccountInput{})
	if err != nil {
		return errors.Wrap(err, "GetAccount")
	}
	// The scenario's whole point is setting the CloudWatch role --
	// verify the setting that was written, not mere object presence
	// (the account object always exists).
	if out.CloudwatchRoleArn == nil || *out.CloudwatchRoleArn == "" {
		return errors.Errorf("api gateway account (%s) has no CloudWatch role set after deploy", id)
	}
	return nil
}

func (*apiGatewayAccountSettingsVerifier) VerifyAbsent(ctx context.Context, cfg aws.Config, id, region string) error {
	out, err := apigateway.NewFromConfig(cfg, func(o *apigateway.Options) {
		if region != "" {
			o.Region = region
		}
	}).GetAccount(ctx, &apigateway.GetAccountInput{})
	if err != nil {
		return errors.Wrap(err, "GetAccount")
	}
	// Destroy RESETS the role; the account object itself always exists.
	if out.CloudwatchRoleArn != nil && *out.CloudwatchRoleArn != "" {
		return errors.Errorf("api gateway account (%s) still has CloudWatch role %q after destroy (expected reset)", id, *out.CloudwatchRoleArn)
	}
	return nil
}

// bedrockInvocationLoggingVerifier verifies AwsBedrockInvocationLogging
// via GetModelInvocationLoggingConfiguration, keyed on
// configured_region. Unlike its class siblings, destroy is a REAL
// delete for this singleton.
type bedrockInvocationLoggingVerifier struct{}

func (*bedrockInvocationLoggingVerifier) IDOutputKey() string { return "configured_region" }

func (*bedrockInvocationLoggingVerifier) VerifyExists(ctx context.Context, cfg aws.Config, id, region string) error {
	out, err := bedrock.NewFromConfig(cfg, func(o *bedrock.Options) {
		if region != "" {
			o.Region = region
		}
	}).GetModelInvocationLoggingConfiguration(ctx, &bedrock.GetModelInvocationLoggingConfigurationInput{})
	if err != nil {
		return errors.Wrap(err, "GetModelInvocationLoggingConfiguration")
	}
	if out.LoggingConfig == nil {
		return errors.Errorf("invocation logging configuration (%s) not found after deploy", id)
	}
	// At least one destination is the spec's own contract -- assert the
	// applied configuration carries one.
	if out.LoggingConfig.CloudWatchConfig == nil && out.LoggingConfig.S3Config == nil {
		return errors.Errorf("invocation logging configuration (%s) has no delivery destination after deploy", id)
	}
	return nil
}

func (*bedrockInvocationLoggingVerifier) VerifyAbsent(ctx context.Context, cfg aws.Config, id, region string) error {
	out, err := bedrock.NewFromConfig(cfg, func(o *bedrock.Options) {
		if region != "" {
			o.Region = region
		}
	}).GetModelInvocationLoggingConfiguration(ctx, &bedrock.GetModelInvocationLoggingConfigurationInput{})
	if err != nil {
		// The API returns an empty body rather than NotFound for an
		// unconfigured region; a hard error here is a real failure.
		return errors.Wrap(err, "GetModelInvocationLoggingConfiguration")
	}
	if out.LoggingConfig != nil {
		return errors.Errorf("invocation logging configuration (%s) still present after destroy", id)
	}
	return nil
}

// agentCoreTokenVaultVerifier verifies AwsBedrockAgentCoreTokenVault
// via GetTokenVault, keyed on token_vault_id. DESTROY IS A NO-OP for
// this singleton: absence means the vault still exists with the
// last-applied setting -- never "gone".
type agentCoreTokenVaultVerifier struct{}

func (*agentCoreTokenVaultVerifier) IDOutputKey() string { return "token_vault_id" }

func (*agentCoreTokenVaultVerifier) VerifyExists(ctx context.Context, cfg aws.Config, id, region string) error {
	out, err := agentCoreClient(cfg, region).GetTokenVault(ctx, &bedrockagentcorecontrol.GetTokenVaultInput{
		TokenVaultId: aws.String(id),
	})
	if err != nil {
		return errors.Wrapf(err, "GetTokenVault(%s)", id)
	}
	if out.KmsConfiguration == nil {
		return errors.Errorf("token vault %s has no KMS configuration after deploy", id)
	}
	return nil
}

func (*agentCoreTokenVaultVerifier) VerifyAbsent(ctx context.Context, cfg aws.Config, id, region string) error {
	// Provider delete is a no-op: the vault and its last-applied key
	// setting REMAIN after destroy. Asserting disappearance would fail
	// every honest run -- the contract is that the vault still answers.
	_, err := agentCoreClient(cfg, region).GetTokenVault(ctx, &bedrockagentcorecontrol.GetTokenVaultInput{
		TokenVaultId: aws.String(id),
	})
	if err != nil {
		return errors.Wrapf(err, "GetTokenVault(%s) after destroy (the vault must persist -- delete is a no-op)", id)
	}
	return nil
}

// sesAccountSettingsVerifier verifies AwsSesAccountSettings via the
// SESv2 GetAccount, keyed on account_id. The destroy contract is
// ASYMMETRIC: the VDM resource's delete resets VDM to DISABLED, while
// the suppression setting persists (upstream delete is a no-op).
type sesAccountSettingsVerifier struct{}

func (*sesAccountSettingsVerifier) IDOutputKey() string { return "account_id" }

func sesAccount(ctx context.Context, cfg aws.Config, region string) (*sesv2.GetAccountOutput, error) {
	return sesv2.NewFromConfig(cfg, func(o *sesv2.Options) {
		if region != "" {
			o.Region = region
		}
	}).GetAccount(ctx, &sesv2.GetAccountInput{})
}

func (*sesAccountSettingsVerifier) VerifyExists(ctx context.Context, cfg aws.Config, id, region string) error {
	out, err := sesAccount(ctx, cfg, region)
	if err != nil {
		return errors.Wrap(err, "sesv2 GetAccount")
	}
	suppressionManaged := out.SuppressionAttributes != nil && len(out.SuppressionAttributes.SuppressedReasons) > 0
	vdmEnabled := out.VdmAttributes != nil && out.VdmAttributes.VdmEnabled == sesv2types.FeatureStatusEnabled
	// Evidence that the applied arms took effect (the scenarios set at
	// least one; both is the canonical posture).
	if !suppressionManaged && !vdmEnabled {
		return errors.Errorf("ses account (%s) shows neither suppression reasons nor VDM enabled after deploy", id)
	}
	return nil
}

func (*sesAccountSettingsVerifier) VerifyAbsent(ctx context.Context, cfg aws.Config, id, region string) error {
	out, err := sesAccount(ctx, cfg, region)
	if err != nil {
		return errors.Wrap(err, "sesv2 GetAccount")
	}
	// VDM is the arm destroy genuinely reverts; suppression persisting
	// is the recorded settings-retention class, never a failure.
	if out.VdmAttributes != nil && out.VdmAttributes.VdmEnabled == sesv2types.FeatureStatusEnabled {
		return errors.Errorf("ses account (%s) still has VDM enabled after destroy (expected DISABLED reset)", id)
	}
	return nil
}
