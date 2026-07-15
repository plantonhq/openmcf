package module

import (
	"github.com/pkg/errors"
	"github.com/pulumi/pulumi-aws/sdk/v7/go/aws"
	"github.com/pulumi/pulumi-aws/sdk/v7/go/aws/batch"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

// schedulingPolicy creates the Batch fair-share scheduling policy.
//
// The fair_share_policy block is emitted whenever ANY dial is set; an empty
// block is also valid AWS-side (all defaults), and emitting it
// unconditionally keeps the resource shape stable as dials are added or
// removed -- the whole surface updates in place.
func schedulingPolicy(
	ctx *pulumi.Context,
	locals *Locals,
	provider *aws.Provider,
) (*batch.SchedulingPolicy, error) {
	spec := locals.AwsBatchSchedulingPolicy.Spec

	fairShare := &batch.SchedulingPolicyFairSharePolicyArgs{}
	if spec.ComputeReservation != nil {
		fairShare.ComputeReservation = pulumi.IntPtr(int(spec.GetComputeReservation()))
	}
	if spec.ShareDecaySeconds != nil {
		fairShare.ShareDecaySeconds = pulumi.IntPtr(int(spec.GetShareDecaySeconds()))
	}
	if len(spec.ShareDistributions) > 0 {
		var distributions batch.SchedulingPolicyFairSharePolicyShareDistributionArray
		for _, dist := range spec.ShareDistributions {
			entry := &batch.SchedulingPolicyFairSharePolicyShareDistributionArgs{
				ShareIdentifier: pulumi.String(dist.ShareIdentifier),
			}
			// weight_factor 0 means "unset" in the spec (AWS then defaults
			// the share's weight to 1.0) -- never send a literal zero,
			// which is below AWS's 0.0001 minimum.
			if dist.WeightFactor > 0 {
				entry.WeightFactor = pulumi.Float64Ptr(dist.WeightFactor)
			}
			distributions = append(distributions, entry)
		}
		fairShare.ShareDistributions = distributions
	}

	args := &batch.SchedulingPolicyArgs{
		// The cloud name comes from metadata.name (the catalog naming
		// basis) -- set explicitly so both engines create the same policy
		// name and Pulumi never auto-names.
		Name:            pulumi.StringPtr(locals.AwsBatchSchedulingPolicy.Metadata.Name),
		FairSharePolicy: fairShare,
		Tags:            pulumi.ToStringMap(locals.AwsTags),
	}

	createdPolicy, err := batch.NewSchedulingPolicy(ctx, "scheduling-policy", args, pulumi.Provider(provider))
	if err != nil {
		return nil, errors.Wrap(err, "create batch scheduling policy")
	}

	return createdPolicy, nil
}
