package module

import (
	"fmt"

	"github.com/pkg/errors"
	"github.com/pulumi/pulumi-aws/sdk/v7/go/aws"
	"github.com/pulumi/pulumi-aws/sdk/v7/go/aws/redshift"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

// usageLimits caps what individual Redshift features may consume on this
// cluster (Spectrum scans, concurrency-scaling time, cross-region
// datasharing transfer). AWS generates the limit IDs at creation; the
// returned map exports them keyed by the same
// feature_type/limit_type/period triple the Terraform module uses (unset
// period rendered as monthly -- the spec's uniqueness CEL applies the
// same normalization), so imports and out-of-band CLI operations can
// address each limit identically on both engines.
func usageLimits(
	ctx *pulumi.Context,
	locals *Locals,
	provider *aws.Provider,
	createdCluster *redshift.Cluster,
) (pulumi.StringMap, error) {
	spec := locals.AwsRedshiftCluster.Spec

	limitIds := pulumi.StringMap{}
	for _, limit := range spec.UsageLimits {
		period := limit.Period
		if period == "" {
			period = "monthly"
		}
		key := fmt.Sprintf("%s/%s/%s", limit.FeatureType, limit.LimitType, period)

		args := &redshift.UsageLimitArgs{
			ClusterIdentifier: createdCluster.ClusterIdentifier,
			FeatureType:       pulumi.String(limit.FeatureType),
			LimitType:         pulumi.String(limit.LimitType),
			Amount:            pulumi.Int(int(limit.Amount)),
			Tags:              pulumi.ToStringMap(locals.AwsTags),
		}

		// Empty keeps the AWS defaults (monthly / log); the provider
		// defaults match, so sending only set values stays faithful.
		if limit.Period != "" {
			args.Period = pulumi.String(limit.Period)
		}
		if limit.BreachAction != "" {
			args.BreachAction = pulumi.String(limit.BreachAction)
		}

		createdLimit, err := redshift.NewUsageLimit(ctx, "usage-limit-"+key, args,
			pulumi.Provider(provider), pulumi.DependsOn([]pulumi.Resource{createdCluster}))
		if err != nil {
			return nil, errors.Wrapf(err, "failed to create usage limit %s", key)
		}
		limitIds[key] = createdLimit.ID().ToStringOutput()
	}
	return limitIds, nil
}
