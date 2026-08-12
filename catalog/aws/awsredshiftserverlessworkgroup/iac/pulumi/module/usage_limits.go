package module

import (
	"fmt"

	"github.com/pkg/errors"
	"github.com/pulumi/pulumi-aws/sdk/v7/go/aws"
	"github.com/pulumi/pulumi-aws/sdk/v7/go/aws/redshiftserverless"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

// usageLimits cap the workgroup's consumption -- RPU-hours of serverless
// compute or terabytes of cross-region datasharing transfer. AWS
// generates the limit IDs at creation; the returned map exports them
// keyed by the same usage_type/period pair the Terraform module uses
// (unset period rendered as monthly -- the spec's uniqueness CEL applies
// the same normalization), so imports and out-of-band CLI operations can
// address each limit identically on both engines. The resource is
// untagged in AWS (unlike the provisioned cluster's usage limits).
func usageLimits(
	ctx *pulumi.Context,
	locals *Locals,
	provider *aws.Provider,
	createdWorkgroup *redshiftserverless.Workgroup,
) (pulumi.StringMap, error) {
	spec := locals.AwsRedshiftServerlessWorkgroup.Spec

	limitIds := pulumi.StringMap{}
	for _, limit := range spec.UsageLimits {
		period := limit.Period
		if period == "" {
			period = "monthly"
		}
		key := fmt.Sprintf("%s/%s", limit.UsageType, period)

		args := &redshiftserverless.UsageLimitArgs{
			ResourceArn: createdWorkgroup.Arn,
			UsageType:   pulumi.String(limit.UsageType),
			Amount:      pulumi.Int(int(limit.Amount)),
		}

		// Empty keeps the AWS defaults (monthly / log); the provider
		// defaults match, so sending only set values stays faithful.
		if limit.Period != "" {
			args.Period = pulumi.String(limit.Period)
		}
		if limit.BreachAction != "" {
			args.BreachAction = pulumi.String(limit.BreachAction)
		}

		createdLimit, err := redshiftserverless.NewUsageLimit(ctx, "usage-limit-"+key, args,
			pulumi.Provider(provider), pulumi.DependsOn([]pulumi.Resource{createdWorkgroup}))
		if err != nil {
			return nil, errors.Wrapf(err, "failed to create usage limit %s", key)
		}
		limitIds[key] = createdLimit.ID().ToStringOutput()
	}
	return limitIds, nil
}
