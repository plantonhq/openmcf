package module

import (
	"github.com/pkg/errors"
	gcpmonitoringalertpolicyv1alpha1 "github.com/plantonhq/planton/catalog/gcp/gcpmonitoringalertpolicy/v1alpha1"
	"github.com/pulumi/pulumi-gcp/sdk/v9/go/gcp"
	"github.com/pulumi/pulumi-gcp/sdk/v9/go/gcp/monitoring"
	"github.com/pulumi/pulumi-gcp/sdk/v9/go/gcp/projects"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

// alertPolicy provisions the Cloud Monitoring alert policy — the rule that
// watches metrics or logs and pages the referenced notification channels.
//
// Each spec condition carries exactly one condition-type arm (the proto CEL
// enforces the API's oneof, which the provider leaves unchecked
// client-side), so the expand code maps arms unconditionally when present.
//
// `enabled` is sent EXPLICITLY on every apply: it is Optional in the
// provider with a server default of true, and a spec transition
// true -> false must reach the API rather than being omitted (the
// send-true-or-omit class silently no-ops such transitions — a silently
// still-enabled alert policy pages people at 3am).
func alertPolicy(ctx *pulumi.Context, locals *Locals, gcpProvider *gcp.Provider) error {
	spec := locals.GcpMonitoringAlertPolicy.Spec

	// Enable the Cloud Monitoring API so a fresh project can host the
	// policy. disable_on_destroy stays false (the provider default):
	// tearing down one policy must never disable monitoring for everything
	// else in the project. Matches the Terraform module.
	serviceArgs := &projects.ServiceArgs{
		Service:                  pulumi.String("monitoring.googleapis.com"),
		DisableDependentServices: pulumi.BoolPtr(true),
	}
	if spec.ProjectId.GetValue() != "" {
		serviceArgs.Project = pulumi.String(spec.ProjectId.GetValue())
	}
	createdProjectService, err := projects.NewService(ctx,
		"alertpolicy-monitoring.googleapis.com", serviceArgs, pulumi.Provider(gcpProvider))
	if err != nil {
		return errors.Wrap(err, "failed to enable monitoring.googleapis.com api")
	}

	conditions := monitoring.AlertPolicyConditionArray{}
	for _, condition := range spec.Conditions {
		conditions = append(conditions, expandCondition(condition))
	}

	args := &monitoring.AlertPolicyArgs{
		DisplayName: pulumi.String(locals.DisplayName),
		Combiner:    pulumi.String(spec.Combiner),
		Conditions:  conditions,
		// Explicit send — see the function comment.
		Enabled:    pulumi.Bool(spec.Enabled == nil || spec.GetEnabled()),
		UserLabels: pulumi.ToStringMap(locals.GcpLabels),
	}

	if spec.Severity != "" {
		args.Severity = pulumi.String(spec.Severity)
	}

	if len(spec.NotificationChannels) > 0 {
		channels := pulumi.StringArray{}
		for _, channel := range spec.NotificationChannels {
			channels = append(channels, pulumi.String(channel.GetValue()))
		}
		args.NotificationChannels = channels
	}

	if spec.AlertStrategy != nil {
		args.AlertStrategy = expandAlertStrategy(spec.AlertStrategy)
	}

	if spec.Documentation != nil {
		documentationArgs := &monitoring.AlertPolicyDocumentationArgs{}
		if spec.Documentation.Content != "" {
			documentationArgs.Content = pulumi.StringPtr(spec.Documentation.Content)
		}
		if spec.Documentation.MimeType != "" {
			documentationArgs.MimeType = pulumi.StringPtr(spec.Documentation.MimeType)
		}
		if spec.Documentation.Subject != "" {
			documentationArgs.Subject = pulumi.StringPtr(spec.Documentation.Subject)
		}
		if len(spec.Documentation.Links) > 0 {
			links := monitoring.AlertPolicyDocumentationLinkArray{}
			for _, link := range spec.Documentation.Links {
				linkArgs := &monitoring.AlertPolicyDocumentationLinkArgs{}
				if link.DisplayName != "" {
					linkArgs.DisplayName = pulumi.StringPtr(link.DisplayName)
				}
				if link.Url != "" {
					linkArgs.Url = pulumi.StringPtr(link.Url)
				}
				links = append(links, linkArgs)
			}
			documentationArgs.Links = links
		}
		args.Documentation = documentationArgs
	}

	// Unset defers to the provider default (DELETE).
	if spec.DeletionPolicy != "" {
		args.DeletionPolicy = pulumi.String(spec.DeletionPolicy)
	}

	// Honor the spec contract: an empty project_id falls back to the
	// provider's default project (empty string would be sent verbatim and
	// rejected).
	if spec.ProjectId.GetValue() != "" {
		args.Project = pulumi.String(spec.ProjectId.GetValue())
	}

	createdAlertPolicy, err := monitoring.NewAlertPolicy(ctx, "alert-policy", args,
		pulumi.Provider(gcpProvider), pulumi.DependsOn([]pulumi.Resource{createdProjectService}))
	if err != nil {
		return errors.Wrap(err, "failed to create alert policy")
	}

	ctx.Export(OpPolicyName, createdAlertPolicy.Name)

	return nil
}

// expandCondition maps one spec condition onto the provider's condition
// block. Exactly one arm is present (proto-CEL-enforced).
func expandCondition(condition *gcpmonitoringalertpolicyv1alpha1.GcpMonitoringAlertPolicyCondition) *monitoring.AlertPolicyConditionArgs {
	conditionArgs := &monitoring.AlertPolicyConditionArgs{
		DisplayName: pulumi.String(condition.DisplayName),
	}

	if threshold := condition.ConditionThreshold; threshold != nil {
		thresholdArgs := &monitoring.AlertPolicyConditionConditionThresholdArgs{
			Comparison: pulumi.String(threshold.Comparison),
			Duration:   pulumi.String(threshold.Duration),
		}
		if threshold.Filter != "" {
			thresholdArgs.Filter = pulumi.StringPtr(threshold.Filter)
		}
		// threshold_value 0 is a legal threshold (e.g. COMPARISON_GT 0 for
		// "any errors at all"), so it is always sent — zero and unset are
		// deliberately NOT distinguished away for this field.
		thresholdArgs.ThresholdValue = pulumi.Float64Ptr(threshold.ThresholdValue)
		if len(threshold.Aggregations) > 0 {
			aggregations := monitoring.AlertPolicyConditionConditionThresholdAggregationArray{}
			for _, aggregation := range threshold.Aggregations {
				aggregationArgs := &monitoring.AlertPolicyConditionConditionThresholdAggregationArgs{}
				if aggregation.AlignmentPeriod != "" {
					aggregationArgs.AlignmentPeriod = pulumi.StringPtr(aggregation.AlignmentPeriod)
				}
				if aggregation.PerSeriesAligner != "" {
					aggregationArgs.PerSeriesAligner = pulumi.StringPtr(aggregation.PerSeriesAligner)
				}
				if aggregation.CrossSeriesReducer != "" {
					aggregationArgs.CrossSeriesReducer = pulumi.StringPtr(aggregation.CrossSeriesReducer)
				}
				if len(aggregation.GroupByFields) > 0 {
					aggregationArgs.GroupByFields = pulumi.ToStringArray(aggregation.GroupByFields)
				}
				aggregations = append(aggregations, aggregationArgs)
			}
			thresholdArgs.Aggregations = aggregations
		}
		if threshold.DenominatorFilter != "" {
			thresholdArgs.DenominatorFilter = pulumi.StringPtr(threshold.DenominatorFilter)
		}
		if len(threshold.DenominatorAggregations) > 0 {
			denominatorAggregations := monitoring.AlertPolicyConditionConditionThresholdDenominatorAggregationArray{}
			for _, aggregation := range threshold.DenominatorAggregations {
				aggregationArgs := &monitoring.AlertPolicyConditionConditionThresholdDenominatorAggregationArgs{}
				if aggregation.AlignmentPeriod != "" {
					aggregationArgs.AlignmentPeriod = pulumi.StringPtr(aggregation.AlignmentPeriod)
				}
				if aggregation.PerSeriesAligner != "" {
					aggregationArgs.PerSeriesAligner = pulumi.StringPtr(aggregation.PerSeriesAligner)
				}
				if aggregation.CrossSeriesReducer != "" {
					aggregationArgs.CrossSeriesReducer = pulumi.StringPtr(aggregation.CrossSeriesReducer)
				}
				if len(aggregation.GroupByFields) > 0 {
					aggregationArgs.GroupByFields = pulumi.ToStringArray(aggregation.GroupByFields)
				}
				denominatorAggregations = append(denominatorAggregations, aggregationArgs)
			}
			thresholdArgs.DenominatorAggregations = denominatorAggregations
		}
		if threshold.ForecastOptions != nil {
			thresholdArgs.ForecastOptions = &monitoring.AlertPolicyConditionConditionThresholdForecastOptionsArgs{
				ForecastHorizon: pulumi.String(threshold.ForecastOptions.ForecastHorizon),
			}
		}
		if threshold.Trigger != nil {
			triggerArgs := &monitoring.AlertPolicyConditionConditionThresholdTriggerArgs{}
			if threshold.Trigger.Count != 0 {
				triggerArgs.Count = pulumi.IntPtr(int(threshold.Trigger.Count))
			}
			if threshold.Trigger.Percent != 0 {
				triggerArgs.Percent = pulumi.Float64Ptr(threshold.Trigger.Percent)
			}
			thresholdArgs.Trigger = triggerArgs
		}
		if threshold.EvaluationMissingData != "" {
			thresholdArgs.EvaluationMissingData = pulumi.StringPtr(threshold.EvaluationMissingData)
		}
		conditionArgs.ConditionThreshold = thresholdArgs
	}

	if absent := condition.ConditionAbsent; absent != nil {
		absentArgs := &monitoring.AlertPolicyConditionConditionAbsentArgs{
			Duration: pulumi.String(absent.Duration),
		}
		if absent.Filter != "" {
			absentArgs.Filter = pulumi.StringPtr(absent.Filter)
		}
		if len(absent.Aggregations) > 0 {
			aggregations := monitoring.AlertPolicyConditionConditionAbsentAggregationArray{}
			for _, aggregation := range absent.Aggregations {
				aggregationArgs := &monitoring.AlertPolicyConditionConditionAbsentAggregationArgs{}
				if aggregation.AlignmentPeriod != "" {
					aggregationArgs.AlignmentPeriod = pulumi.StringPtr(aggregation.AlignmentPeriod)
				}
				if aggregation.PerSeriesAligner != "" {
					aggregationArgs.PerSeriesAligner = pulumi.StringPtr(aggregation.PerSeriesAligner)
				}
				if aggregation.CrossSeriesReducer != "" {
					aggregationArgs.CrossSeriesReducer = pulumi.StringPtr(aggregation.CrossSeriesReducer)
				}
				if len(aggregation.GroupByFields) > 0 {
					aggregationArgs.GroupByFields = pulumi.ToStringArray(aggregation.GroupByFields)
				}
				aggregations = append(aggregations, aggregationArgs)
			}
			absentArgs.Aggregations = aggregations
		}
		if absent.Trigger != nil {
			triggerArgs := &monitoring.AlertPolicyConditionConditionAbsentTriggerArgs{}
			if absent.Trigger.Count != 0 {
				triggerArgs.Count = pulumi.IntPtr(int(absent.Trigger.Count))
			}
			if absent.Trigger.Percent != 0 {
				triggerArgs.Percent = pulumi.Float64Ptr(absent.Trigger.Percent)
			}
			absentArgs.Trigger = triggerArgs
		}
		conditionArgs.ConditionAbsent = absentArgs
	}

	if matchedLog := condition.ConditionMatchedLog; matchedLog != nil {
		matchedLogArgs := &monitoring.AlertPolicyConditionConditionMatchedLogArgs{
			Filter: pulumi.String(matchedLog.Filter),
		}
		if len(matchedLog.LabelExtractors) > 0 {
			matchedLogArgs.LabelExtractors = pulumi.ToStringMap(matchedLog.LabelExtractors)
		}
		conditionArgs.ConditionMatchedLog = matchedLogArgs
	}

	if mql := condition.ConditionMonitoringQueryLanguage; mql != nil {
		mqlArgs := &monitoring.AlertPolicyConditionConditionMonitoringQueryLanguageArgs{
			Query:    pulumi.String(mql.Query),
			Duration: pulumi.String(mql.Duration),
		}
		if mql.Trigger != nil {
			triggerArgs := &monitoring.AlertPolicyConditionConditionMonitoringQueryLanguageTriggerArgs{}
			if mql.Trigger.Count != 0 {
				triggerArgs.Count = pulumi.IntPtr(int(mql.Trigger.Count))
			}
			if mql.Trigger.Percent != 0 {
				triggerArgs.Percent = pulumi.Float64Ptr(mql.Trigger.Percent)
			}
			mqlArgs.Trigger = triggerArgs
		}
		if mql.EvaluationMissingData != "" {
			mqlArgs.EvaluationMissingData = pulumi.StringPtr(mql.EvaluationMissingData)
		}
		conditionArgs.ConditionMonitoringQueryLanguage = mqlArgs
	}

	if promql := condition.ConditionPrometheusQueryLanguage; promql != nil {
		promqlArgs := &monitoring.AlertPolicyConditionConditionPrometheusQueryLanguageArgs{
			Query: pulumi.String(promql.Query),
		}
		if promql.Duration != "" {
			promqlArgs.Duration = pulumi.StringPtr(promql.Duration)
		}
		if promql.EvaluationInterval != "" {
			promqlArgs.EvaluationInterval = pulumi.StringPtr(promql.EvaluationInterval)
		}
		if len(promql.Labels) > 0 {
			promqlArgs.Labels = pulumi.ToStringMap(promql.Labels)
		}
		if promql.RuleGroup != "" {
			promqlArgs.RuleGroup = pulumi.StringPtr(promql.RuleGroup)
		}
		if promql.AlertRule != "" {
			promqlArgs.AlertRule = pulumi.StringPtr(promql.AlertRule)
		}
		if promql.DisableMetricValidation {
			promqlArgs.DisableMetricValidation = pulumi.BoolPtr(true)
		}
		conditionArgs.ConditionPrometheusQueryLanguage = promqlArgs
	}

	if sql := condition.ConditionSql; sql != nil {
		sqlArgs := &monitoring.AlertPolicyConditionConditionSqlArgs{
			Query: pulumi.String(sql.Query),
		}
		if sql.Minutes != nil {
			sqlArgs.Minutes = &monitoring.AlertPolicyConditionConditionSqlMinutesArgs{
				Periodicity: pulumi.Int(int(sql.Minutes.Periodicity)),
			}
		}
		if sql.Hourly != nil {
			hourlyArgs := &monitoring.AlertPolicyConditionConditionSqlHourlyArgs{
				Periodicity: pulumi.Int(int(sql.Hourly.Periodicity)),
			}
			if sql.Hourly.MinuteOffset != nil {
				hourlyArgs.MinuteOffset = pulumi.IntPtr(int(sql.Hourly.GetMinuteOffset()))
			}
			sqlArgs.Hourly = hourlyArgs
		}
		if sql.Daily != nil {
			dailyArgs := &monitoring.AlertPolicyConditionConditionSqlDailyArgs{
				Periodicity: pulumi.Int(int(sql.Daily.Periodicity)),
			}
			if sql.Daily.ExecutionTime != nil {
				dailyArgs.ExecutionTime = &monitoring.AlertPolicyConditionConditionSqlDailyExecutionTimeArgs{
					Hours:   pulumi.IntPtr(int(sql.Daily.ExecutionTime.Hours)),
					Minutes: pulumi.IntPtr(int(sql.Daily.ExecutionTime.Minutes)),
					Seconds: pulumi.IntPtr(int(sql.Daily.ExecutionTime.Seconds)),
					Nanos:   pulumi.IntPtr(int(sql.Daily.ExecutionTime.Nanos)),
				}
			}
			sqlArgs.Daily = dailyArgs
		}
		if sql.RowCountTest != nil {
			sqlArgs.RowCountTest = &monitoring.AlertPolicyConditionConditionSqlRowCountTestArgs{
				Comparison: pulumi.String(sql.RowCountTest.Comparison),
				Threshold:  pulumi.Int(int(sql.RowCountTest.Threshold)),
			}
		}
		if sql.BooleanTest != nil {
			sqlArgs.BooleanTest = &monitoring.AlertPolicyConditionConditionSqlBooleanTestArgs{
				Column: pulumi.String(sql.BooleanTest.Column),
			}
		}
		conditionArgs.ConditionSql = sqlArgs
	}

	return conditionArgs
}

// expandAlertStrategy maps the notification-behavior block.
func expandAlertStrategy(strategy *gcpmonitoringalertpolicyv1alpha1.GcpMonitoringAlertPolicyAlertStrategy) *monitoring.AlertPolicyAlertStrategyArgs {
	strategyArgs := &monitoring.AlertPolicyAlertStrategyArgs{}

	if strategy.AutoClose != "" {
		strategyArgs.AutoClose = pulumi.StringPtr(strategy.AutoClose)
	}
	if strategy.NotificationRateLimit != nil {
		rateLimitArgs := &monitoring.AlertPolicyAlertStrategyNotificationRateLimitArgs{}
		if strategy.NotificationRateLimit.Period != "" {
			rateLimitArgs.Period = pulumi.StringPtr(strategy.NotificationRateLimit.Period)
		}
		strategyArgs.NotificationRateLimit = rateLimitArgs
	}
	if len(strategy.NotificationChannelStrategy) > 0 {
		channelStrategies := monitoring.AlertPolicyAlertStrategyNotificationChannelStrategyArray{}
		for _, channelStrategy := range strategy.NotificationChannelStrategy {
			channelStrategyArgs := &monitoring.AlertPolicyAlertStrategyNotificationChannelStrategyArgs{}
			if len(channelStrategy.NotificationChannelNames) > 0 {
				names := pulumi.StringArray{}
				for _, name := range channelStrategy.NotificationChannelNames {
					names = append(names, pulumi.String(name.GetValue()))
				}
				channelStrategyArgs.NotificationChannelNames = names
			}
			if channelStrategy.RenotifyInterval != "" {
				channelStrategyArgs.RenotifyInterval = pulumi.StringPtr(channelStrategy.RenotifyInterval)
			}
			channelStrategies = append(channelStrategies, channelStrategyArgs)
		}
		strategyArgs.NotificationChannelStrategies = channelStrategies
	}
	if len(strategy.NotificationPrompts) > 0 {
		strategyArgs.NotificationPrompts = pulumi.ToStringArray(strategy.NotificationPrompts)
	}

	return strategyArgs
}
