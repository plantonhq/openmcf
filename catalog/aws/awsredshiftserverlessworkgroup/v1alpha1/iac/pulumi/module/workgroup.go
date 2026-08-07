package module

import (
	"github.com/pkg/errors"
	"github.com/pulumi/pulumi-aws/sdk/v7/go/aws"
	"github.com/pulumi/pulumi-aws/sdk/v7/go/aws/redshiftserverless"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

// workgroup provisions the workgroup itself -- RPU capacity, VPC
// placement, reachability, and query-level configuration. Create-only in
// AWS: the workgroup name and the namespace it serves. Capacity edits
// are in place, though AWS serializes them (the provider issues one
// capacity change per update call for exactly that reason).
func workgroup(
	ctx *pulumi.Context,
	locals *Locals,
	provider *aws.Provider,
) (*redshiftserverless.Workgroup, error) {
	spec := locals.AwsRedshiftServerlessWorkgroup.Spec

	args := &redshiftserverless.WorkgroupArgs{
		WorkgroupName: pulumi.String(locals.WorkgroupName),
		NamespaceName: pulumi.String(spec.NamespaceName.GetValue()),

		// Reachability is opt-in; a private workgroup is the norm and
		// stays reachable through VPC routing.
		PubliclyAccessible: pulumi.Bool(spec.PubliclyAccessible),

		// Enhanced VPC routing forces COPY/UNLOAD data movement through
		// the VPC where flow logs and endpoints can see and govern it.
		EnhancedVpcRouting: pulumi.Bool(spec.EnhancedVpcRouting),

		Tags: pulumi.ToStringMap(locals.AwsTags),
	}

	// The capacity contract (CEL enforces the shape): a fixed RPU
	// baseline OR an enabled price-performance target where AWS owns the
	// baseline. 0 keeps the AWS default (128 RPU); AWS validates the
	// exact RPU increments at deploy, since they have changed over time.
	if spec.BaseCapacity != 0 {
		args.BaseCapacity = pulumi.Int(int(spec.BaseCapacity))
	}

	// 0 leaves scaling uncapped (the AWS default): the provider treats
	// an unset max as "no ceiling", so it is forwarded only when set.
	if spec.MaxCapacity != 0 {
		args.MaxCapacity = pulumi.Int(int(spec.MaxCapacity))
	}

	// The price-performance dial is forwarded even when disabled IF the
	// message is present, so a later enable/disable edits in place; an
	// absent message keeps the argument off entirely.
	if spec.PricePerformanceTarget != nil {
		pricePerformanceTarget := &redshiftserverless.WorkgroupPricePerformanceTargetArgs{
			Enabled: pulumi.Bool(spec.PricePerformanceTarget.Enabled),
		}
		if spec.PricePerformanceTarget.Level != 0 {
			pricePerformanceTarget.Level = pulumi.Int(int(spec.PricePerformanceTarget.Level))
		}
		args.PricePerformanceTarget = pricePerformanceTarget
	}

	// VPC placement: at least three subnets in three AZs (CEL-enforced
	// when set); empty falls back to the account's default VPC. The VPC
	// default security group applies when none are given (AWS's own
	// default).
	if len(spec.SubnetIds) > 0 {
		subnetIds := pulumi.StringArray{}
		for _, subnetId := range spec.SubnetIds {
			subnetIds = append(subnetIds, pulumi.String(subnetId.GetValue()))
		}
		args.SubnetIds = subnetIds
	}
	if len(spec.SecurityGroupIds) > 0 {
		securityGroupIds := pulumi.StringArray{}
		for _, securityGroupId := range spec.SecurityGroupIds {
			securityGroupIds = append(securityGroupIds, pulumi.String(securityGroupId.GetValue()))
		}
		args.SecurityGroupIds = securityGroupIds
	}

	// 0 keeps the AWS default (5439); CEL constrains set values to the
	// only ranges Redshift Serverless accepts.
	if spec.Port != 0 {
		args.Port = pulumi.Int(int(spec.Port))
	}

	// Query-level parameters apply directly to the workgroup --
	// serverless has no parameter groups, so there is nothing to fold or
	// reference.
	if len(spec.ConfigParameters) > 0 {
		configParameters := redshiftserverless.WorkgroupConfigParameterArray{}
		for _, configParameter := range spec.ConfigParameters {
			configParameters = append(configParameters, &redshiftserverless.WorkgroupConfigParameterArgs{
				ParameterKey:   pulumi.String(configParameter.Name),
				ParameterValue: pulumi.String(configParameter.Value),
			})
		}
		args.ConfigParameters = configParameters
	}

	// Empty keeps the AWS default release track ("current").
	if spec.TrackName != "" {
		args.TrackName = pulumi.String(spec.TrackName)
	}

	createdWorkgroup, err := redshiftserverless.NewWorkgroup(ctx, "workgroup", args, pulumi.Provider(provider))
	if err != nil {
		return nil, errors.Wrap(err, "failed to create Redshift Serverless workgroup")
	}
	return createdWorkgroup, nil
}
