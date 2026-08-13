package module

import (
	"fmt"

	"github.com/pkg/errors"
	awsdynamodbv1alpha1 "github.com/plantonhq/planton/catalog/aws/awsdynamodb/v1alpha1"
	"github.com/pulumi/pulumi-aws/sdk/v7/go/aws"
	"github.com/pulumi/pulumi-aws/sdk/v7/go/aws/appautoscaling"
	"github.com/pulumi/pulumi-aws/sdk/v7/go/aws/dynamodb"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

// autoscaling folds capacity management: Application Auto Scaling owns
// the table's live read/write capacity on EVERY provisioned table, in
// one of two modes.
//
//   - spec.autoscaling set: the user's min/max bounds with
//     target-tracking policies (and optional scheduled adjustments).
//   - spec.autoscaling unset: pinned targets (min = max = the declared
//     provisioned_throughput values). A declared capacity change updates
//     the target, and AWS moves out-of-range capacity into the new
//     bounds by contract (live-verified 2026-08-13: re-registering a
//     pinned write target 5/5 -> 8/8 moved DescribeTable's
//     WriteCapacityUnits to 8 within ~15 seconds, no table operation
//     involved) -- so capacity stays fully declarative even though the
//     table resource ignores it (table.go).
//
// One target per dimension either way, which is what makes adding or
// removing the autoscaling block an in-place update -- the table
// resource never changes shape. The scaler's identity IS this table
// (one scalable target per table dimension), so it lives here rather
// than as its own kind. PAY_PER_REQUEST tables scale natively and get
// none of this.
func autoscaling(ctx *pulumi.Context, locals *Locals, provider *aws.Provider,
	createdTable *dynamodb.Table) error {
	spec := locals.AwsDynamodb.Spec

	if spec.BillingMode == "PAY_PER_REQUEST" || spec.ProvisionedThroughput == nil {
		return nil
	}

	resourceID := pulumi.Sprintf("table/%s", createdTable.Name)

	// Per-dimension bounds: the user's autoscaling bounds, or the pinned
	// declared capacity.
	readMin, readMax := spec.ProvisionedThroughput.ReadCapacityUnits, spec.ProvisionedThroughput.ReadCapacityUnits
	writeMin, writeMax := spec.ProvisionedThroughput.WriteCapacityUnits, spec.ProvisionedThroughput.WriteCapacityUnits
	if spec.Autoscaling != nil && spec.Autoscaling.Read != nil {
		readMin, readMax = spec.Autoscaling.Read.MinCapacity, spec.Autoscaling.Read.MaxCapacity
	}
	if spec.Autoscaling != nil && spec.Autoscaling.Write != nil {
		writeMin, writeMax = spec.Autoscaling.Write.MinCapacity, spec.Autoscaling.Write.MaxCapacity
	}

	readTarget, err := appautoscaling.NewTarget(ctx, "read-target", &appautoscaling.TargetArgs{
		ServiceNamespace:  pulumi.String("dynamodb"),
		ResourceId:        resourceID,
		ScalableDimension: pulumi.String("dynamodb:table:ReadCapacityUnits"),
		MinCapacity:       pulumi.Int(int(readMin)),
		MaxCapacity:       pulumi.Int(int(readMax)),
	}, pulumi.Provider(provider), pulumi.Parent(createdTable))
	if err != nil {
		return errors.Wrap(err, "failed to register read scalable target")
	}

	writeTarget, err := appautoscaling.NewTarget(ctx, "write-target", &appautoscaling.TargetArgs{
		ServiceNamespace:  pulumi.String("dynamodb"),
		ResourceId:        resourceID,
		ScalableDimension: pulumi.String("dynamodb:table:WriteCapacityUnits"),
		MinCapacity:       pulumi.Int(int(writeMin)),
		MaxCapacity:       pulumi.Int(int(writeMax)),
	}, pulumi.Provider(provider), pulumi.Parent(createdTable))
	if err != nil {
		return errors.Wrap(err, "failed to register write scalable target")
	}

	// Target tracking holds consumed-to-provisioned utilization near the
	// target percentage -- only rendered when the user configured real
	// autoscaling for the dimension (pinned targets need no policy:
	// min = max leaves the scaler nothing to decide).
	if spec.Autoscaling != nil && spec.Autoscaling.Read != nil {
		if err := trackingPolicy(ctx, provider, locals.TableName, "read", readTarget,
			"DynamoDBReadCapacityUtilization", spec.Autoscaling.Read); err != nil {
			return err
		}
	}
	if spec.Autoscaling != nil && spec.Autoscaling.Write != nil {
		if err := trackingPolicy(ctx, provider, locals.TableName, "write", writeTarget,
			"DynamoDBWriteCapacityUtilization", spec.Autoscaling.Write); err != nil {
			return err
		}
	}

	// Scheduled capacity adjustments, keyed by name so entries come and
	// go independently. Each targets its dimension's registered scalable
	// target (CEL guarantees the dimension's autoscaling config exists).
	if spec.Autoscaling != nil {
		for _, adjustment := range spec.Autoscaling.ScheduledAdjustments {
			dimension := "dynamodb:table:WriteCapacityUnits"
			parentTarget := writeTarget
			if adjustment.Dimension == "READ" {
				dimension = "dynamodb:table:ReadCapacityUnits"
				parentTarget = readTarget
			}

			// 0 leaves the bound unchanged when the schedule fires (CEL
			// requires at least one).
			targetAction := &appautoscaling.ScheduledActionScalableTargetActionArgs{}
			if adjustment.MinCapacity > 0 {
				targetAction.MinCapacity = pulumi.IntPtr(int(adjustment.MinCapacity))
			}
			if adjustment.MaxCapacity > 0 {
				targetAction.MaxCapacity = pulumi.IntPtr(int(adjustment.MaxCapacity))
			}

			actionArgs := &appautoscaling.ScheduledActionArgs{
				Name:                 pulumi.String(adjustment.Name),
				ServiceNamespace:     pulumi.String("dynamodb"),
				ResourceId:           resourceID,
				ScalableDimension:    pulumi.String(dimension),
				Schedule:             pulumi.String(adjustment.Schedule),
				ScalableTargetAction: targetAction,
			}
			if adjustment.Timezone != "" {
				actionArgs.Timezone = pulumi.StringPtr(adjustment.Timezone)
			}
			if adjustment.StartTime != "" {
				actionArgs.StartTime = pulumi.StringPtr(adjustment.StartTime)
			}
			if adjustment.EndTime != "" {
				actionArgs.EndTime = pulumi.StringPtr(adjustment.EndTime)
			}

			if _, err := appautoscaling.NewScheduledAction(ctx,
				fmt.Sprintf("scheduled-adjustment-%s", adjustment.Name),
				actionArgs, pulumi.Provider(provider), pulumi.Parent(parentTarget)); err != nil {
				return errors.Wrapf(err, "failed to create scheduled adjustment %s", adjustment.Name)
			}
		}
	}

	return nil
}

// trackingPolicy renders one dimension's target-tracking policy, named
// "{table}-{dimension}-utilization" to match the Terraform module.
func trackingPolicy(ctx *pulumi.Context, provider *aws.Provider, tableName string,
	dimension string, target *appautoscaling.Target, metricType string,
	cfg *awsdynamodbv1alpha1.AwsDynamodbCapacityAutoscaling) error {
	policyCfg := &appautoscaling.PolicyTargetTrackingScalingPolicyConfigurationArgs{
		TargetValue: pulumi.Float64(float64(cfg.TargetUtilizationPercent)),
		PredefinedMetricSpecification: &appautoscaling.PolicyTargetTrackingScalingPolicyConfigurationPredefinedMetricSpecificationArgs{
			PredefinedMetricType: pulumi.String(metricType),
		},
	}
	// 0 keeps the AWS default cooldown.
	if cfg.ScaleInCooldownSeconds > 0 {
		policyCfg.ScaleInCooldown = pulumi.IntPtr(int(cfg.ScaleInCooldownSeconds))
	}
	if cfg.ScaleOutCooldownSeconds > 0 {
		policyCfg.ScaleOutCooldown = pulumi.IntPtr(int(cfg.ScaleOutCooldownSeconds))
	}

	if _, err := appautoscaling.NewPolicy(ctx, dimension+"-utilization-policy", &appautoscaling.PolicyArgs{
		Name:                                     pulumi.Sprintf("%s-%s-utilization", tableName, dimension),
		PolicyType:                               pulumi.String("TargetTrackingScaling"),
		ServiceNamespace:                         target.ServiceNamespace,
		ResourceId:                               target.ResourceId,
		ScalableDimension:                        target.ScalableDimension,
		TargetTrackingScalingPolicyConfiguration: policyCfg,
	}, pulumi.Provider(provider), pulumi.Parent(target)); err != nil {
		return errors.Wrapf(err, "failed to create %s target-tracking policy", dimension)
	}
	return nil
}
