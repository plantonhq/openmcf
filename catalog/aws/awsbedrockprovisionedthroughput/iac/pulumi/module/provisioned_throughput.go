package module

import (
	"github.com/pkg/errors"
	"github.com/pulumi/pulumi-aws/sdk/v7/go/aws"
	"github.com/pulumi/pulumi-aws/sdk/v7/go/aws/bedrock"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

// provisionedThroughput purchases the dedicated model serving capacity
// and exports outputs.
//
// COST: this resource bills from the moment it is created. Without a
// commitment it bills hourly and deletes any time; with a OneMonth/
// SixMonths commitment it bills for the FULL term and cannot be deleted
// until the term lapses (the provider refuses the destroy server-side).
// Every argument except tags is create-time-immutable.
func provisionedThroughput(ctx *pulumi.Context, locals *Locals, provider *aws.Provider) error {
	spec := locals.Spec

	args := &bedrock.ProvisionedModelThroughputArgs{
		ProvisionedModelName: pulumi.String(locals.ProvisionedModelName),
		// The model capacity is bought for -- a custom model's output ARN
		// (the default reference wiring) or a foundation-model ARN where
		// AWS allows direct provisioning.
		ModelArn: pulumi.String(spec.ModelArn.GetValue()),
		// A model unit's throughput (tokens/minute) is model-specific;
		// account quotas for NO-COMMITMENT units default low (often 0-2).
		ModelUnits: pulumi.Int(int(spec.ModelUnits)),
		Tags:       pulumi.ToStringMap(locals.AwsTags),
	}

	// Omitted entirely for no commitment (hourly billing) -- the provider
	// treats the absent argument as the no-commitment purchase.
	if spec.CommitmentDuration != "" {
		args.CommitmentDuration = pulumi.String(spec.CommitmentDuration)
	}

	createdThroughput, err := bedrock.NewProvisionedModelThroughput(ctx, locals.ProvisionedModelName, args, pulumi.Provider(provider))
	if err != nil {
		return errors.Wrap(err, "create provisioned model throughput")
	}

	ctx.Export(OpProvisionedModelArn, createdThroughput.ProvisionedModelArn)
	ctx.Export(OpProvisionedModelName, createdThroughput.ProvisionedModelName)

	return nil
}
