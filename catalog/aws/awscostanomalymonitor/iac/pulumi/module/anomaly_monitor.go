package module

import (
	"encoding/json"
	"fmt"

	"github.com/pkg/errors"
	awscostanomalymonitorv1alpha1 "github.com/plantonhq/planton/catalog/aws/awscostanomalymonitor/v1alpha1"
	"github.com/pulumi/pulumi-aws/sdk/v7/go/aws/costexplorer"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
	structpb "google.golang.org/protobuf/types/known/structpb"
)

// anomalyMonitor provisions the Cost Explorer anomaly monitor and its
// folded alert subscriptions, and exports outputs.
//
// Lifecycle facts the render below depends on:
//   - the monitor's SHAPE arms are create-only: monitor_type,
//     monitor_dimension, and monitor_specification all force
//     replacement (only the display name updates in place);
//   - the spec's CUSTOM-arm Struct renders as the provider's raw
//     Expression JSON string (the provider takes the AWS document
//     verbatim, not typed blocks) - the AUTHOR owns the canonical
//     form (server-verified 2026-08-25): all root members present
//     (unused ones null, matching the provider's re-marshaled READ
//     form) and tag keys "user:"-prefixed (CE echoes them back
//     prefixed). A sparser or unprefixed document is a perpetual
//     replace, not a module defect;
//   - each subscription is one aws_ce_anomaly_subscription bound to
//     THIS monitor's ARN, keyed by its spec name (the for_each key
//     and the outputs-map key);
//   - the threshold expression is leveled (root -> leaf) to the exact
//     nesting AWS accepts on subscriptions; the conversion below is
//     1:1 with no depth checks because the spec shape cannot express
//     an illegal tree;
//   - a DIMENSIONAL monitor is an account SINGLETON that AWS
//     auto-creates ("Default-Services-Monitor") on post-2023 Cost
//     Explorer accounts - CreateAnomalyMonitor then fails with
//     "Limit exceeded on dimensional spend monitor creation"
//     (server-verified 2026-08-25). Not a module defect: import the
//     existing monitor or use the CUSTOM arm.
func anomalyMonitor(ctx *pulumi.Context, locals *Locals, provider pulumi.ProviderResource) error {
	spec := locals.Spec

	args := &costexplorer.AnomalyMonitorArgs{
		Name:        pulumi.StringPtr(spec.MonitorName),
		MonitorType: pulumi.String(spec.MonitorType),
		Tags:        pulumi.ToStringMap(locals.AwsTags),
	}
	if spec.MonitorDimension != "" {
		args.MonitorDimension = pulumi.StringPtr(spec.MonitorDimension)
	}
	if spec.MonitorSpecification != nil {
		specificationJSON, err := structToJSONString(spec.MonitorSpecification)
		if err != nil {
			return errors.Wrap(err, "marshal monitor specification")
		}
		args.MonitorSpecification = pulumi.StringPtr(specificationJSON)
	}

	createdMonitor, err := costexplorer.NewAnomalyMonitor(ctx, "anomaly-monitor", args, pulumi.Provider(provider))
	if err != nil {
		return errors.Wrap(err, "create anomaly monitor")
	}

	subscriptionArns := pulumi.StringMap{}
	for _, sub := range spec.Subscriptions {
		subscribers := costexplorer.AnomalySubscriptionSubscriberArray{}
		for _, s := range sub.Subscribers {
			subscribers = append(subscribers, costexplorer.AnomalySubscriptionSubscriberArgs{
				Address: pulumi.String(s.Address.GetValue()),
				Type:    pulumi.String(s.Type),
			})
		}

		subArgs := &costexplorer.AnomalySubscriptionArgs{
			Name:            pulumi.StringPtr(sub.Name),
			Frequency:       pulumi.String(sub.Frequency),
			MonitorArnLists: pulumi.StringArray{createdMonitor.Arn},
			Subscribers:     subscribers,
			Tags:            pulumi.ToStringMap(locals.AwsTags),
		}
		if sub.ThresholdExpression != nil {
			subArgs.ThresholdExpression = buildThresholdExpression(sub.ThresholdExpression)
		}

		createdSubscription, err := costexplorer.NewAnomalySubscription(ctx,
			fmt.Sprintf("subscription-%s", sub.Name), subArgs,
			pulumi.Provider(provider), pulumi.Parent(createdMonitor))
		if err != nil {
			return errors.Wrapf(err, "create subscription %s", sub.Name)
		}
		subscriptionArns[sub.Name] = createdSubscription.Arn
	}

	ctx.Export(OpMonitorArn, createdMonitor.Arn)
	ctx.Export(OpSubscriptionArns, subscriptionArns)

	return nil
}

// buildThresholdExpression converts the spec's leveled threshold tree
// (root -> leaf) to the provider's per-path typed args.
func buildThresholdExpression(root *awscostanomalymonitorv1alpha1.AwsCostAnomalyMonitorExpression) *costexplorer.AnomalySubscriptionThresholdExpressionArgs {
	out := &costexplorer.AnomalySubscriptionThresholdExpressionArgs{}
	if root.Dimension != nil {
		out.Dimension = &costexplorer.AnomalySubscriptionThresholdExpressionDimensionArgs{
			Key:          pulumi.StringPtr(root.Dimension.Key),
			MatchOptions: pulumi.ToStringArray(root.Dimension.MatchOptions),
			Values:       pulumi.ToStringArray(root.Dimension.Values),
		}
	}
	if root.Tag != nil {
		out.Tags = &costexplorer.AnomalySubscriptionThresholdExpressionTagsArgs{
			Key:          pulumi.StringPtr(root.Tag.Key),
			MatchOptions: pulumi.ToStringArray(root.Tag.MatchOptions),
			Values:       pulumi.ToStringArray(root.Tag.Values),
		}
	}
	if root.CostCategory != nil {
		out.CostCategory = &costexplorer.AnomalySubscriptionThresholdExpressionCostCategoryArgs{
			Key:          pulumi.StringPtr(root.CostCategory.Key),
			MatchOptions: pulumi.ToStringArray(root.CostCategory.MatchOptions),
			Values:       pulumi.ToStringArray(root.CostCategory.Values),
		}
	}
	if len(root.And) > 0 {
		ands := costexplorer.AnomalySubscriptionThresholdExpressionAndArray{}
		for _, leaf := range root.And {
			child := costexplorer.AnomalySubscriptionThresholdExpressionAndArgs{}
			if leaf.Dimension != nil {
				child.Dimension = &costexplorer.AnomalySubscriptionThresholdExpressionAndDimensionArgs{
					Key:          pulumi.StringPtr(leaf.Dimension.Key),
					MatchOptions: pulumi.ToStringArray(leaf.Dimension.MatchOptions),
					Values:       pulumi.ToStringArray(leaf.Dimension.Values),
				}
			}
			if leaf.Tag != nil {
				child.Tags = &costexplorer.AnomalySubscriptionThresholdExpressionAndTagsArgs{
					Key:          pulumi.StringPtr(leaf.Tag.Key),
					MatchOptions: pulumi.ToStringArray(leaf.Tag.MatchOptions),
					Values:       pulumi.ToStringArray(leaf.Tag.Values),
				}
			}
			if leaf.CostCategory != nil {
				child.CostCategory = &costexplorer.AnomalySubscriptionThresholdExpressionAndCostCategoryArgs{
					Key:          pulumi.StringPtr(leaf.CostCategory.Key),
					MatchOptions: pulumi.ToStringArray(leaf.CostCategory.MatchOptions),
					Values:       pulumi.ToStringArray(leaf.CostCategory.Values),
				}
			}
			ands = append(ands, child)
		}
		out.Ands = ands
	}
	if len(root.Or) > 0 {
		ors := costexplorer.AnomalySubscriptionThresholdExpressionOrArray{}
		for _, leaf := range root.Or {
			child := costexplorer.AnomalySubscriptionThresholdExpressionOrArgs{}
			if leaf.Dimension != nil {
				child.Dimension = &costexplorer.AnomalySubscriptionThresholdExpressionOrDimensionArgs{
					Key:          pulumi.StringPtr(leaf.Dimension.Key),
					MatchOptions: pulumi.ToStringArray(leaf.Dimension.MatchOptions),
					Values:       pulumi.ToStringArray(leaf.Dimension.Values),
				}
			}
			if leaf.Tag != nil {
				child.Tags = &costexplorer.AnomalySubscriptionThresholdExpressionOrTagsArgs{
					Key:          pulumi.StringPtr(leaf.Tag.Key),
					MatchOptions: pulumi.ToStringArray(leaf.Tag.MatchOptions),
					Values:       pulumi.ToStringArray(leaf.Tag.Values),
				}
			}
			if leaf.CostCategory != nil {
				child.CostCategory = &costexplorer.AnomalySubscriptionThresholdExpressionOrCostCategoryArgs{
					Key:          pulumi.StringPtr(leaf.CostCategory.Key),
					MatchOptions: pulumi.ToStringArray(leaf.CostCategory.MatchOptions),
					Values:       pulumi.ToStringArray(leaf.CostCategory.Values),
				}
			}
			ors = append(ors, child)
		}
		out.Ors = ors
	}
	if root.Not != nil {
		child := &costexplorer.AnomalySubscriptionThresholdExpressionNotArgs{}
		if root.Not.Dimension != nil {
			child.Dimension = &costexplorer.AnomalySubscriptionThresholdExpressionNotDimensionArgs{
				Key:          pulumi.StringPtr(root.Not.Dimension.Key),
				MatchOptions: pulumi.ToStringArray(root.Not.Dimension.MatchOptions),
				Values:       pulumi.ToStringArray(root.Not.Dimension.Values),
			}
		}
		if root.Not.Tag != nil {
			child.Tags = &costexplorer.AnomalySubscriptionThresholdExpressionNotTagsArgs{
				Key:          pulumi.StringPtr(root.Not.Tag.Key),
				MatchOptions: pulumi.ToStringArray(root.Not.Tag.MatchOptions),
				Values:       pulumi.ToStringArray(root.Not.Tag.Values),
			}
		}
		if root.Not.CostCategory != nil {
			child.CostCategory = &costexplorer.AnomalySubscriptionThresholdExpressionNotCostCategoryArgs{
				Key:          pulumi.StringPtr(root.Not.CostCategory.Key),
				MatchOptions: pulumi.ToStringArray(root.Not.CostCategory.MatchOptions),
				Values:       pulumi.ToStringArray(root.Not.CostCategory.Values),
			}
		}
		out.Not = child
	}
	return out
}

// structToJSONString converts a google.protobuf.Struct to a raw JSON string.
func structToJSONString(s *structpb.Struct) (string, error) {
	if s == nil {
		return "{}", nil
	}
	bytes, err := json.Marshal(s.AsMap())
	if err != nil {
		return "", err
	}
	return string(bytes), nil
}
