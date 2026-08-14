package module

import (
	"github.com/pkg/errors"
	"github.com/pulumi/pulumi-aws/sdk/v7/go/aws"
	"github.com/pulumi/pulumi-aws/sdk/v7/go/aws/apigateway"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

// accountSettings manages the region's ONE API Gateway account object
// and exports outputs.
//
// Lifecycle facts the render below depends on:
//   - the one configurable lever is the CloudWatch Logs role; an
//     empty string resets it (the provider patches /cloudwatchRoleArn
//     to nil for "" and unset alike), so passing the spec value
//     through unconditionally is faithful for both the set and clear
//     postures;
//   - destroy RESETS the role to none -- the account object itself
//     always exists and cannot be deleted;
//   - AWS validates the role at apply (trust: apigateway.amazonaws.com
//     plus CloudWatch Logs write permissions) and the provider retries
//     through IAM propagation lag.
func accountSettings(ctx *pulumi.Context, locals *Locals, provider *aws.Provider) error {
	spec := locals.Spec

	created, err := apigateway.NewAccount(ctx, "account-settings", &apigateway.AccountArgs{
		CloudwatchRoleArn: pulumi.String(spec.CloudwatchRoleArn.GetValue()),
	}, pulumi.Provider(provider))
	if err != nil {
		return errors.Wrap(err, "update account settings")
	}

	ctx.Export(OpAccountId, created.ID())
	ctx.Export(OpApiKeyVersion, created.ApiKeyVersion)
	ctx.Export(OpFeatures, created.Features)
	// throttle_settings is a computed one-element list upstream.
	ctx.Export(OpThrottleBurstLimit, created.ThrottleSettings.Index(pulumi.Int(0)).BurstLimit())
	ctx.Export(OpThrottleRateLimit, created.ThrottleSettings.Index(pulumi.Int(0)).RateLimit())
	return nil
}
