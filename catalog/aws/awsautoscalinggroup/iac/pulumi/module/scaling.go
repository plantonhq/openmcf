package module

import (
	"fmt"
	"strconv"

	"github.com/pkg/errors"
	awsautoscalinggroupv1alpha1 "github.com/plantonhq/planton/catalog/aws/awsautoscalinggroup/v1alpha1"
	"github.com/pulumi/pulumi-aws/sdk/v7/go/aws/autoscaling"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

// scaling provisions the folded sub-resources of the group: scaling
// policies, scheduled actions, lifecycle hooks, and event notifications.
// Each is an AWS sub-resource of exactly ONE group with no referenceable
// identity of its own -- which is why they live in this spec instead of
// being standalone kinds. Keying them on the group's Name output gives
// every one an implicit dependency on the group.
func scaling(ctx *pulumi.Context, locals *Locals, provider pulumi.ProviderResource, createdGroup *autoscaling.Group) error {
	spec := locals.AwsAutoScalingGroup.Spec

	for _, policy := range spec.ScalingPolicies {
		if err := scalingPolicy(ctx, provider, createdGroup, policy); err != nil {
			return errors.Wrapf(err, "failed to create scaling policy %q", policy.Name)
		}
	}

	for _, action := range spec.ScheduledActions {
		args := &autoscaling.ScheduleArgs{
			AutoscalingGroupName: createdGroup.Name,
			ScheduledActionName:  pulumi.String(action.Name),
		}
		if action.Recurrence != "" {
			args.Recurrence = pulumi.StringPtr(action.Recurrence)
		}
		if action.TimeZone != "" {
			args.TimeZone = pulumi.StringPtr(action.TimeZone)
		}
		if action.StartTime != "" {
			args.StartTime = pulumi.StringPtr(action.StartTime)
		}
		if action.EndTime != "" {
			args.EndTime = pulumi.StringPtr(action.EndTime)
		}
		// Absent capacity values mean "leave unchanged": AWS expresses that
		// as -1, which is why the proto models them as optional ints.
		args.MinSize = scheduleCapacity(action.MinSize)
		args.MaxSize = scheduleCapacity(action.MaxSize)
		args.DesiredCapacity = scheduleCapacity(action.DesiredCapacity)

		if _, err := autoscaling.NewSchedule(ctx,
			fmt.Sprintf("%s-schedule-%s", locals.AwsAutoScalingGroup.Metadata.Name, action.Name),
			args, pulumi.Provider(provider)); err != nil {
			return errors.Wrapf(err, "failed to create scheduled action %q", action.Name)
		}
	}

	for _, hook := range spec.LifecycleHooks {
		// Hooks flagged apply_at_launch render inline on the group
		// (attached atomically at creation); only the rest are standalone
		// hook resources here.
		if hook.ApplyAtLaunch {
			continue
		}
		args := &autoscaling.LifecycleHookArgs{
			AutoscalingGroupName: createdGroup.Name,
			Name:                 pulumi.StringPtr(hook.Name),
			LifecycleTransition:  pulumi.String(hook.LifecycleTransition),
		}
		if hook.DefaultResult != "" {
			args.DefaultResult = pulumi.StringPtr(hook.DefaultResult)
		}
		if hook.HeartbeatTimeoutSeconds > 0 {
			args.HeartbeatTimeout = pulumi.IntPtr(int(hook.HeartbeatTimeoutSeconds))
		}
		if hook.NotificationTargetArn.GetValue() != "" {
			args.NotificationTargetArn = pulumi.StringPtr(hook.NotificationTargetArn.GetValue())
		}
		if hook.RoleArn.GetValue() != "" {
			args.RoleArn = pulumi.StringPtr(hook.RoleArn.GetValue())
		}
		if hook.NotificationMetadata != "" {
			args.NotificationMetadata = pulumi.StringPtr(hook.NotificationMetadata)
		}

		if _, err := autoscaling.NewLifecycleHook(ctx,
			fmt.Sprintf("%s-hook-%s", locals.AwsAutoScalingGroup.Metadata.Name, hook.Name),
			args, pulumi.Provider(provider)); err != nil {
			return errors.Wrapf(err, "failed to create lifecycle hook %q", hook.Name)
		}
	}

	if spec.Notifications != nil {
		notificationTypes := make(autoscaling.NotificationTypeArray, 0, len(spec.Notifications.EventTypes))
		for _, eventType := range spec.Notifications.EventTypes {
			notificationTypes = append(notificationTypes, autoscaling.NotificationType(eventType))
		}
		if _, err := autoscaling.NewNotification(ctx,
			fmt.Sprintf("%s-notifications", locals.AwsAutoScalingGroup.Metadata.Name),
			&autoscaling.NotificationArgs{
				GroupNames:    pulumi.StringArray{createdGroup.Name},
				Notifications: notificationTypes,
				TopicArn:      pulumi.String(spec.Notifications.Topic.GetValue()),
			}, pulumi.Provider(provider)); err != nil {
			return errors.Wrap(err, "failed to create notifications")
		}
	}

	return nil
}

// scheduleCapacity converts an optional capacity to the AWS convention:
// absent means "leave unchanged", expressed as -1.
func scheduleCapacity(value *int32) pulumi.IntPtrInput {
	if value == nil {
		return pulumi.IntPtr(-1)
	}
	return pulumi.IntPtr(int(*value))
}

// scalingPolicy provisions one policy. policy_type decides which
// configuration block applies (the spec's CEL enforces the discriminated
// union), mirroring how PutScalingPolicy interprets its input.
func scalingPolicy(ctx *pulumi.Context, provider pulumi.ProviderResource, createdGroup *autoscaling.Group, policy *awsautoscalinggroupv1alpha1.AwsAutoScalingGroupScalingPolicy) error {
	args := &autoscaling.PolicyArgs{
		AutoscalingGroupName: createdGroup.Name,
		Name:                 pulumi.StringPtr(policy.Name),
		PolicyType:           pulumi.StringPtr(policy.PolicyType),
	}
	// The pause button: a disabled policy stays configured (alarms,
	// history, forecast state) but stops acting on the group. AWS
	// defaults to enabled, so only an explicit disable is sent.
	if policy.Disabled {
		args.Enabled = pulumi.BoolPtr(false)
	}
	if policy.EstimatedInstanceWarmupSeconds > 0 {
		args.EstimatedInstanceWarmup = pulumi.IntPtr(int(policy.EstimatedInstanceWarmupSeconds))
	}

	switch {
	case policy.TargetTracking != nil:
		args.TargetTrackingConfiguration = targetTrackingArgs(policy.TargetTracking)

	case policy.StepScaling != nil:
		args.AdjustmentType = pulumi.StringPtr(policy.StepScaling.AdjustmentType)
		if policy.StepScaling.MetricAggregationType != "" {
			args.MetricAggregationType = pulumi.StringPtr(policy.StepScaling.MetricAggregationType)
		}
		if policy.StepScaling.MinAdjustmentMagnitude > 0 {
			args.MinAdjustmentMagnitude = pulumi.IntPtr(int(policy.StepScaling.MinAdjustmentMagnitude))
		}
		steps := make(autoscaling.PolicyStepAdjustmentArray, 0, len(policy.StepScaling.StepAdjustments))
		for _, step := range policy.StepScaling.StepAdjustments {
			stepArgs := &autoscaling.PolicyStepAdjustmentArgs{
				ScalingAdjustment: pulumi.Int(int(step.ScalingAdjustment)),
			}
			if step.MetricIntervalLowerBound != "" {
				stepArgs.MetricIntervalLowerBound = pulumi.StringPtr(step.MetricIntervalLowerBound)
			}
			if step.MetricIntervalUpperBound != "" {
				stepArgs.MetricIntervalUpperBound = pulumi.StringPtr(step.MetricIntervalUpperBound)
			}
			steps = append(steps, stepArgs)
		}
		args.StepAdjustments = steps

	case policy.SimpleScaling != nil:
		args.AdjustmentType = pulumi.StringPtr(policy.SimpleScaling.AdjustmentType)
		args.ScalingAdjustment = pulumi.IntPtr(int(policy.SimpleScaling.ScalingAdjustment))
		if policy.SimpleScaling.CooldownSeconds > 0 {
			args.Cooldown = pulumi.IntPtr(int(policy.SimpleScaling.CooldownSeconds))
		}
		if policy.SimpleScaling.MinAdjustmentMagnitude > 0 {
			args.MinAdjustmentMagnitude = pulumi.IntPtr(int(policy.SimpleScaling.MinAdjustmentMagnitude))
		}

	case policy.PredictiveScaling != nil:
		args.PredictiveScalingConfiguration = predictiveScalingArgs(policy.PredictiveScaling)
	}

	_, err := autoscaling.NewPolicy(ctx, policy.Name, args, pulumi.Provider(provider))
	return err
}

// targetTrackingArgs maps the thermostat model: hold a predefined or
// customized metric (single-metric or metric-math form) at a target value.
func targetTrackingArgs(tracking *awsautoscalinggroupv1alpha1.AwsAutoScalingGroupTargetTrackingConfig) *autoscaling.PolicyTargetTrackingConfigurationArgs {
	args := &autoscaling.PolicyTargetTrackingConfigurationArgs{
		TargetValue: pulumi.Float64(tracking.TargetValue),
	}
	if tracking.DisableScaleIn {
		args.DisableScaleIn = pulumi.BoolPtr(true)
	}
	if tracking.PredefinedMetricType != "" {
		predefined := &autoscaling.PolicyTargetTrackingConfigurationPredefinedMetricSpecificationArgs{
			PredefinedMetricType: pulumi.String(tracking.PredefinedMetricType),
		}
		if tracking.ResourceLabel != "" {
			predefined.ResourceLabel = pulumi.StringPtr(tracking.ResourceLabel)
		}
		args.PredefinedMetricSpecification = predefined
	}
	if tracking.CustomizedMetric != nil {
		args.CustomizedMetricSpecification = customizedMetricArgs(tracking.CustomizedMetric)
	}
	return args
}

// customizedMetricArgs maps a custom CloudWatch metric: either the
// single-metric form or the metric-math form (a set of query expressions).
func customizedMetricArgs(metric *awsautoscalinggroupv1alpha1.AwsAutoScalingGroupCustomizedMetric) *autoscaling.PolicyTargetTrackingConfigurationCustomizedMetricSpecificationArgs {
	args := &autoscaling.PolicyTargetTrackingConfigurationCustomizedMetricSpecificationArgs{}

	if metric.MetricName != "" {
		args.MetricName = pulumi.StringPtr(metric.MetricName)
	}
	if metric.Namespace != "" {
		args.Namespace = pulumi.StringPtr(metric.Namespace)
	}
	if metric.Statistic != "" {
		args.Statistic = pulumi.StringPtr(metric.Statistic)
	}
	if metric.Unit != "" {
		args.Unit = pulumi.StringPtr(metric.Unit)
	}
	if metric.PeriodSeconds > 0 {
		args.Period = pulumi.IntPtr(int(metric.PeriodSeconds))
	}
	if len(metric.Dimensions) > 0 {
		dimensions := make(autoscaling.PolicyTargetTrackingConfigurationCustomizedMetricSpecificationMetricDimensionArray, 0, len(metric.Dimensions))
		for _, dimension := range metric.Dimensions {
			dimensions = append(dimensions, &autoscaling.PolicyTargetTrackingConfigurationCustomizedMetricSpecificationMetricDimensionArgs{
				Name:  pulumi.String(dimension.Name),
				Value: pulumi.String(dimension.Value),
			})
		}
		args.MetricDimensions = dimensions
	}

	if len(metric.Metrics) > 0 {
		queries := make(autoscaling.PolicyTargetTrackingConfigurationCustomizedMetricSpecificationMetricArray, 0, len(metric.Metrics))
		for _, query := range metric.Metrics {
			queryArgs := &autoscaling.PolicyTargetTrackingConfigurationCustomizedMetricSpecificationMetricArgs{
				Id: pulumi.String(query.Id),
			}
			if query.Expression != "" {
				queryArgs.Expression = pulumi.StringPtr(query.Expression)
			}
			if query.Label != "" {
				queryArgs.Label = pulumi.StringPtr(query.Label)
			}
			// return_data defaults to true at AWS; only an explicit value
			// is sent, so intermediate entries carry an explicit false.
			if query.ReturnData != nil {
				queryArgs.ReturnData = pulumi.BoolPtr(query.GetReturnData())
			}
			if query.MetricStat != nil {
				statArgs := &autoscaling.PolicyTargetTrackingConfigurationCustomizedMetricSpecificationMetricMetricStatArgs{
					Stat: pulumi.String(query.MetricStat.Stat),
					Metric: &autoscaling.PolicyTargetTrackingConfigurationCustomizedMetricSpecificationMetricMetricStatMetricArgs{
						MetricName: pulumi.String(query.MetricStat.MetricName),
						Namespace:  pulumi.String(query.MetricStat.Namespace),
						Dimensions: metricStatDimensions(query.MetricStat.Dimensions),
					},
				}
				if query.MetricStat.Unit != "" {
					statArgs.Unit = pulumi.StringPtr(query.MetricStat.Unit)
				}
				if query.MetricStat.PeriodSeconds > 0 {
					statArgs.Period = pulumi.IntPtr(int(query.MetricStat.PeriodSeconds))
				}
				queryArgs.MetricStat = statArgs
			}
			queries = append(queries, queryArgs)
		}
		args.Metrics = queries
	}

	return args
}

func metricStatDimensions(dimensions []*awsautoscalinggroupv1alpha1.AwsAutoScalingGroupMetricDimension) autoscaling.PolicyTargetTrackingConfigurationCustomizedMetricSpecificationMetricMetricStatMetricDimensionArray {
	if len(dimensions) == 0 {
		return nil
	}
	result := make(autoscaling.PolicyTargetTrackingConfigurationCustomizedMetricSpecificationMetricMetricStatMetricDimensionArray, 0, len(dimensions))
	for _, dimension := range dimensions {
		result = append(result, &autoscaling.PolicyTargetTrackingConfigurationCustomizedMetricSpecificationMetricMetricStatMetricDimensionArgs{
			Name:  pulumi.String(dimension.Name),
			Value: pulumi.String(dimension.Value),
		})
	}
	return result
}

// predictiveScalingArgs maps forecast-driven pre-provisioning. The metrics
// come in three forms (the spec's CEL enforces the choice): one predefined
// PAIR, SPLIT predefined load + scaling metrics, or fully CUSTOMIZED
// metric-math query sets.
func predictiveScalingArgs(predictive *awsautoscalinggroupv1alpha1.AwsAutoScalingGroupPredictiveScalingConfig) *autoscaling.PolicyPredictiveScalingConfigurationArgs {
	metricSpecification := &autoscaling.PolicyPredictiveScalingConfigurationMetricSpecificationArgs{
		TargetValue: pulumi.Float64(predictive.TargetValue),
	}
	if predictive.PredefinedMetricPairType != "" {
		pairArgs := &autoscaling.PolicyPredictiveScalingConfigurationMetricSpecificationPredefinedMetricPairSpecificationArgs{
			PredefinedMetricType: pulumi.String(predictive.PredefinedMetricPairType),
		}
		if predictive.ResourceLabel != "" {
			pairArgs.ResourceLabel = pulumi.StringPtr(predictive.ResourceLabel)
		}
		metricSpecification.PredefinedMetricPairSpecification = pairArgs
	}
	if predictive.PredefinedLoadMetric != nil {
		loadArgs := &autoscaling.PolicyPredictiveScalingConfigurationMetricSpecificationPredefinedLoadMetricSpecificationArgs{
			PredefinedMetricType: pulumi.String(predictive.PredefinedLoadMetric.MetricType),
		}
		if predictive.PredefinedLoadMetric.ResourceLabel != "" {
			loadArgs.ResourceLabel = pulumi.StringPtr(predictive.PredefinedLoadMetric.ResourceLabel)
		}
		metricSpecification.PredefinedLoadMetricSpecification = loadArgs
	}
	if predictive.PredefinedScalingMetric != nil {
		scalingArgs := &autoscaling.PolicyPredictiveScalingConfigurationMetricSpecificationPredefinedScalingMetricSpecificationArgs{
			PredefinedMetricType: pulumi.String(predictive.PredefinedScalingMetric.MetricType),
		}
		if predictive.PredefinedScalingMetric.ResourceLabel != "" {
			scalingArgs.ResourceLabel = pulumi.StringPtr(predictive.PredefinedScalingMetric.ResourceLabel)
		}
		metricSpecification.PredefinedScalingMetricSpecification = scalingArgs
	}
	if len(predictive.CustomizedLoadMetricQueries) > 0 {
		queries := make(autoscaling.PolicyPredictiveScalingConfigurationMetricSpecificationCustomizedLoadMetricSpecificationMetricDataQueryArray, 0, len(predictive.CustomizedLoadMetricQueries))
		for _, query := range predictive.CustomizedLoadMetricQueries {
			queryArgs := &autoscaling.PolicyPredictiveScalingConfigurationMetricSpecificationCustomizedLoadMetricSpecificationMetricDataQueryArgs{
				Id: pulumi.String(query.Id),
			}
			if query.Expression != "" {
				queryArgs.Expression = pulumi.StringPtr(query.Expression)
			}
			if query.Label != "" {
				queryArgs.Label = pulumi.StringPtr(query.Label)
			}
			if query.ReturnData != nil {
				queryArgs.ReturnData = pulumi.BoolPtr(query.GetReturnData())
			}
			if query.MetricStat != nil {
				statArgs := &autoscaling.PolicyPredictiveScalingConfigurationMetricSpecificationCustomizedLoadMetricSpecificationMetricDataQueryMetricStatArgs{
					Stat: pulumi.String(query.MetricStat.Stat),
					Metric: &autoscaling.PolicyPredictiveScalingConfigurationMetricSpecificationCustomizedLoadMetricSpecificationMetricDataQueryMetricStatMetricArgs{
						MetricName: pulumi.String(query.MetricStat.MetricName),
						Namespace:  pulumi.String(query.MetricStat.Namespace),
						Dimensions: predictiveLoadDimensions(query.MetricStat.Dimensions),
					},
				}
				if query.MetricStat.Unit != "" {
					statArgs.Unit = pulumi.StringPtr(query.MetricStat.Unit)
				}
				queryArgs.MetricStat = statArgs
			}
			queries = append(queries, queryArgs)
		}
		metricSpecification.CustomizedLoadMetricSpecification = &autoscaling.PolicyPredictiveScalingConfigurationMetricSpecificationCustomizedLoadMetricSpecificationArgs{
			MetricDataQueries: queries,
		}
	}
	if len(predictive.CustomizedScalingMetricQueries) > 0 {
		queries := make(autoscaling.PolicyPredictiveScalingConfigurationMetricSpecificationCustomizedScalingMetricSpecificationMetricDataQueryArray, 0, len(predictive.CustomizedScalingMetricQueries))
		for _, query := range predictive.CustomizedScalingMetricQueries {
			queryArgs := &autoscaling.PolicyPredictiveScalingConfigurationMetricSpecificationCustomizedScalingMetricSpecificationMetricDataQueryArgs{
				Id: pulumi.String(query.Id),
			}
			if query.Expression != "" {
				queryArgs.Expression = pulumi.StringPtr(query.Expression)
			}
			if query.Label != "" {
				queryArgs.Label = pulumi.StringPtr(query.Label)
			}
			if query.ReturnData != nil {
				queryArgs.ReturnData = pulumi.BoolPtr(query.GetReturnData())
			}
			if query.MetricStat != nil {
				statArgs := &autoscaling.PolicyPredictiveScalingConfigurationMetricSpecificationCustomizedScalingMetricSpecificationMetricDataQueryMetricStatArgs{
					Stat: pulumi.String(query.MetricStat.Stat),
					Metric: &autoscaling.PolicyPredictiveScalingConfigurationMetricSpecificationCustomizedScalingMetricSpecificationMetricDataQueryMetricStatMetricArgs{
						MetricName: pulumi.String(query.MetricStat.MetricName),
						Namespace:  pulumi.String(query.MetricStat.Namespace),
						Dimensions: predictiveScalingDimensions(query.MetricStat.Dimensions),
					},
				}
				if query.MetricStat.Unit != "" {
					statArgs.Unit = pulumi.StringPtr(query.MetricStat.Unit)
				}
				queryArgs.MetricStat = statArgs
			}
			queries = append(queries, queryArgs)
		}
		metricSpecification.CustomizedScalingMetricSpecification = &autoscaling.PolicyPredictiveScalingConfigurationMetricSpecificationCustomizedScalingMetricSpecificationArgs{
			MetricDataQueries: queries,
		}
	}
	if len(predictive.CustomizedCapacityMetricQueries) > 0 {
		queries := make(autoscaling.PolicyPredictiveScalingConfigurationMetricSpecificationCustomizedCapacityMetricSpecificationMetricDataQueryArray, 0, len(predictive.CustomizedCapacityMetricQueries))
		for _, query := range predictive.CustomizedCapacityMetricQueries {
			queryArgs := &autoscaling.PolicyPredictiveScalingConfigurationMetricSpecificationCustomizedCapacityMetricSpecificationMetricDataQueryArgs{
				Id: pulumi.String(query.Id),
			}
			if query.Expression != "" {
				queryArgs.Expression = pulumi.StringPtr(query.Expression)
			}
			if query.Label != "" {
				queryArgs.Label = pulumi.StringPtr(query.Label)
			}
			if query.ReturnData != nil {
				queryArgs.ReturnData = pulumi.BoolPtr(query.GetReturnData())
			}
			if query.MetricStat != nil {
				statArgs := &autoscaling.PolicyPredictiveScalingConfigurationMetricSpecificationCustomizedCapacityMetricSpecificationMetricDataQueryMetricStatArgs{
					Stat: pulumi.String(query.MetricStat.Stat),
					Metric: &autoscaling.PolicyPredictiveScalingConfigurationMetricSpecificationCustomizedCapacityMetricSpecificationMetricDataQueryMetricStatMetricArgs{
						MetricName: pulumi.String(query.MetricStat.MetricName),
						Namespace:  pulumi.String(query.MetricStat.Namespace),
						Dimensions: predictiveCapacityDimensions(query.MetricStat.Dimensions),
					},
				}
				if query.MetricStat.Unit != "" {
					statArgs.Unit = pulumi.StringPtr(query.MetricStat.Unit)
				}
				queryArgs.MetricStat = statArgs
			}
			queries = append(queries, queryArgs)
		}
		metricSpecification.CustomizedCapacityMetricSpecification = &autoscaling.PolicyPredictiveScalingConfigurationMetricSpecificationCustomizedCapacityMetricSpecificationArgs{
			MetricDataQueries: queries,
		}
	}

	args := &autoscaling.PolicyPredictiveScalingConfigurationArgs{
		MetricSpecification: metricSpecification,
	}
	if predictive.Mode != "" {
		args.Mode = pulumi.StringPtr(predictive.Mode)
	}
	// The provider models buffer time and capacity buffer as strings
	// (nullable ints at AWS); the proto keeps honest ints and converts.
	if predictive.SchedulingBufferTimeSeconds > 0 {
		args.SchedulingBufferTime = pulumi.StringPtr(strconv.Itoa(int(predictive.SchedulingBufferTimeSeconds)))
	}
	if predictive.MaxCapacityBreachBehavior != "" {
		args.MaxCapacityBreachBehavior = pulumi.StringPtr(predictive.MaxCapacityBreachBehavior)
	}
	if predictive.MaxCapacityBuffer > 0 {
		args.MaxCapacityBuffer = pulumi.StringPtr(strconv.Itoa(int(predictive.MaxCapacityBuffer)))
	}

	return args
}

// The three predictive dimension helpers exist only because the provider
// SDK generates a distinct (structurally identical) type per metric form.
func predictiveLoadDimensions(dimensions []*awsautoscalinggroupv1alpha1.AwsAutoScalingGroupMetricDimension) autoscaling.PolicyPredictiveScalingConfigurationMetricSpecificationCustomizedLoadMetricSpecificationMetricDataQueryMetricStatMetricDimensionArray {
	if len(dimensions) == 0 {
		return nil
	}
	result := make(autoscaling.PolicyPredictiveScalingConfigurationMetricSpecificationCustomizedLoadMetricSpecificationMetricDataQueryMetricStatMetricDimensionArray, 0, len(dimensions))
	for _, dimension := range dimensions {
		result = append(result, &autoscaling.PolicyPredictiveScalingConfigurationMetricSpecificationCustomizedLoadMetricSpecificationMetricDataQueryMetricStatMetricDimensionArgs{
			Name:  pulumi.String(dimension.Name),
			Value: pulumi.String(dimension.Value),
		})
	}
	return result
}

func predictiveScalingDimensions(dimensions []*awsautoscalinggroupv1alpha1.AwsAutoScalingGroupMetricDimension) autoscaling.PolicyPredictiveScalingConfigurationMetricSpecificationCustomizedScalingMetricSpecificationMetricDataQueryMetricStatMetricDimensionArray {
	if len(dimensions) == 0 {
		return nil
	}
	result := make(autoscaling.PolicyPredictiveScalingConfigurationMetricSpecificationCustomizedScalingMetricSpecificationMetricDataQueryMetricStatMetricDimensionArray, 0, len(dimensions))
	for _, dimension := range dimensions {
		result = append(result, &autoscaling.PolicyPredictiveScalingConfigurationMetricSpecificationCustomizedScalingMetricSpecificationMetricDataQueryMetricStatMetricDimensionArgs{
			Name:  pulumi.String(dimension.Name),
			Value: pulumi.String(dimension.Value),
		})
	}
	return result
}

func predictiveCapacityDimensions(dimensions []*awsautoscalinggroupv1alpha1.AwsAutoScalingGroupMetricDimension) autoscaling.PolicyPredictiveScalingConfigurationMetricSpecificationCustomizedCapacityMetricSpecificationMetricDataQueryMetricStatMetricDimensionArray {
	if len(dimensions) == 0 {
		return nil
	}
	result := make(autoscaling.PolicyPredictiveScalingConfigurationMetricSpecificationCustomizedCapacityMetricSpecificationMetricDataQueryMetricStatMetricDimensionArray, 0, len(dimensions))
	for _, dimension := range dimensions {
		result = append(result, &autoscaling.PolicyPredictiveScalingConfigurationMetricSpecificationCustomizedCapacityMetricSpecificationMetricDataQueryMetricStatMetricDimensionArgs{
			Name:  pulumi.String(dimension.Name),
			Value: pulumi.String(dimension.Value),
		})
	}
	return result
}
