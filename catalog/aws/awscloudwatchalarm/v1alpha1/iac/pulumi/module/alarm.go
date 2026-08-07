package module

import (
	"github.com/pkg/errors"
	fkv1 "github.com/plantonhq/planton/shared/foreignkey/v1"
	"github.com/pulumi/pulumi-aws/sdk/v7/go/aws"
	"github.com/pulumi/pulumi-aws/sdk/v7/go/aws/cloudwatch"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

type AlarmResult struct {
	AlarmArn  pulumi.StringOutput
	AlarmName pulumi.StringOutput
}

func alarm(ctx *pulumi.Context, locals *Locals, provider *aws.Provider) (*AlarmResult, error) {
	spec := locals.AwsCloudwatchAlarm.Spec

	args := &cloudwatch.MetricAlarmArgs{
		// The alarm's cloud name is the resource's metadata.name — the same
		// basis the Terraform module uses. Setting it explicitly (instead of
		// relying on Pulumi auto-naming, which appends a random suffix) keeps
		// the physical name identical across both IaC engines. The name is also
		// the alarm's identity for composite alarm rules, so it must be stable.
		Name: pulumi.String(locals.AwsCloudwatchAlarm.Metadata.Name),
		Tags: pulumi.ToStringMap(locals.AwsTags),
	}

	// PromQL alarms carry their alerting contract in the query itself, so the
	// threshold-evaluation arguments below only apply to the other two modes.
	// The spec CELs guarantee comparison_operator/evaluation_periods are set
	// exactly when the alarm is NOT a PromQL alarm.
	if spec.EvaluationCriteria == nil {
		args.ComparisonOperator = pulumi.String(spec.ComparisonOperator)
		args.EvaluationPeriods = pulumi.IntPtr(int(spec.EvaluationPeriods))

		// Threshold vs. anomaly detection: mutually exclusive.
		// When threshold_metric_id is set, the anomaly detection band acts as the
		// dynamic threshold. Otherwise, use the static threshold value.
		if spec.ThresholdMetricId != "" {
			args.ThresholdMetricId = pulumi.StringPtr(spec.ThresholdMetricId)
		} else {
			args.Threshold = pulumi.Float64Ptr(spec.Threshold)
		}
	}

	// Datapoints-to-alarm: M-of-N evaluation. Only set when explicitly provided
	// (non-zero), otherwise AWS defaults to evaluation_periods.
	if spec.DatapointsToAlarm > 0 {
		args.DatapointsToAlarm = pulumi.IntPtr(int(spec.DatapointsToAlarm))
	}

	// Treat missing data: controls alarm behavior when data points are absent.
	if spec.TreatMissingData != "" {
		args.TreatMissingData = pulumi.StringPtr(spec.TreatMissingData)
	}

	// Actions enabled is a presence-aware optional: unset means "let AWS
	// default to true", while an explicit false is a real user choice
	// (suppressing actions during maintenance) that must reach the provider.
	if spec.ActionsEnabled != nil {
		args.ActionsEnabled = pulumi.BoolPtr(spec.GetActionsEnabled())
	}

	// Alarm description.
	if spec.AlarmDescription != "" {
		args.AlarmDescription = pulumi.StringPtr(spec.AlarmDescription)
	}

	// Percentile low-sample-count behavior.
	if spec.EvaluateLowSampleCountPercentiles != "" {
		args.EvaluateLowSampleCountPercentiles = pulumi.StringPtr(spec.EvaluateLowSampleCountPercentiles)
	}

	// ---------------------------------------------------------------------------
	// Simple metric mode (metric_name, namespace, period, statistic/extended)
	// ---------------------------------------------------------------------------
	if spec.MetricName != "" {
		args.MetricName = pulumi.StringPtr(spec.MetricName)
		args.Namespace = pulumi.StringPtr(spec.Namespace)
		args.Period = pulumi.IntPtr(int(spec.Period))

		if spec.Statistic != "" {
			args.Statistic = pulumi.StringPtr(spec.Statistic)
		}
		if spec.ExtendedStatistic != "" {
			args.ExtendedStatistic = pulumi.StringPtr(spec.ExtendedStatistic)
		}
		if len(spec.Dimensions) > 0 {
			args.Dimensions = pulumi.ToStringMap(spec.Dimensions)
		}
		if spec.Unit != "" {
			args.Unit = pulumi.StringPtr(spec.Unit)
		}
	}

	// ---------------------------------------------------------------------------
	// Metric query mode (metric math, anomaly detection, multi-metric)
	// ---------------------------------------------------------------------------
	if len(spec.MetricQueries) > 0 {
		queries := cloudwatch.MetricAlarmMetricQueryArray{}
		for _, mq := range spec.MetricQueries {
			query := cloudwatch.MetricAlarmMetricQueryArgs{
				Id: pulumi.String(mq.Id),
			}

			if mq.Expression != "" {
				query.Expression = pulumi.StringPtr(mq.Expression)
			}
			if mq.Label != "" {
				query.Label = pulumi.StringPtr(mq.Label)
			}
			if mq.Period > 0 {
				query.Period = pulumi.IntPtr(int(mq.Period))
			}
			if mq.ReturnData {
				query.ReturnData = pulumi.BoolPtr(true)
			}
			if mq.AccountId != "" {
				query.AccountId = pulumi.StringPtr(mq.AccountId)
			}

			// Raw metric definition within this query.
			if mq.Metric != nil {
				m := mq.Metric
				metricArgs := cloudwatch.MetricAlarmMetricQueryMetricArgs{
					MetricName: pulumi.String(m.MetricName),
					Namespace:  pulumi.StringPtr(m.Namespace),
					Period:     pulumi.Int(int(m.Period)),
					Stat:       pulumi.String(m.Stat),
				}
				if len(m.Dimensions) > 0 {
					metricArgs.Dimensions = pulumi.ToStringMap(m.Dimensions)
				}
				if m.Unit != "" {
					metricArgs.Unit = pulumi.StringPtr(m.Unit)
				}
				query.Metric = metricArgs
			}

			queries = append(queries, query)
		}
		args.MetricQueries = queries
	}

	// ---------------------------------------------------------------------------
	// PromQL mode (evaluation_criteria)
	// ---------------------------------------------------------------------------
	if spec.EvaluationCriteria != nil {
		promql := spec.EvaluationCriteria.PromqlCriteria
		promqlArgs := &cloudwatch.MetricAlarmEvaluationCriteriaPromqlCriteriaArgs{
			Query: pulumi.String(promql.Query),
		}
		// pending/recovery periods are presence-aware: unset lets AWS apply its
		// default firing behavior, while an explicit 0 means "transition
		// immediately" — the two must not collapse into each other.
		if promql.PendingPeriod != nil {
			promqlArgs.PendingPeriod = pulumi.IntPtr(int(promql.GetPendingPeriod()))
		}
		if promql.RecoveryPeriod != nil {
			promqlArgs.RecoveryPeriod = pulumi.IntPtr(int(promql.GetRecoveryPeriod()))
		}
		args.EvaluationCriteria = &cloudwatch.MetricAlarmEvaluationCriteriaArgs{
			PromqlCriteria: promqlArgs,
		}
		if spec.EvaluationInterval > 0 {
			args.EvaluationInterval = pulumi.IntPtr(int(spec.EvaluationInterval))
		}
	}

	// ---------------------------------------------------------------------------
	// Actions: convert repeated StringValueOrRef to pulumi.Array
	// ---------------------------------------------------------------------------
	if len(spec.AlarmActions) > 0 {
		args.AlarmActions = buildActionArns(spec.AlarmActions)
	}
	if len(spec.OkActions) > 0 {
		args.OkActions = buildActionArns(spec.OkActions)
	}
	if len(spec.InsufficientDataActions) > 0 {
		args.InsufficientDataActions = buildActionArns(spec.InsufficientDataActions)
	}

	// ---------------------------------------------------------------------------
	// Create the metric alarm
	// ---------------------------------------------------------------------------
	createdAlarm, err := cloudwatch.NewMetricAlarm(
		ctx,
		locals.AwsCloudwatchAlarm.Metadata.Name,
		args,
		pulumi.Provider(provider),
	)
	if err != nil {
		return nil, errors.Wrap(err, "failed to create cloudwatch metric alarm")
	}

	return &AlarmResult{
		AlarmArn:  createdAlarm.Arn,
		AlarmName: createdAlarm.Name,
	}, nil
}

// buildActionArns converts a slice of StringValueOrRef into a pulumi.Array
// suitable for alarm/ok/insufficient-data action ARN fields.
func buildActionArns(actions []*fkv1.StringValueOrRef) pulumi.Array {
	result := pulumi.Array{}
	for _, action := range actions {
		if action.GetValue() != "" {
			result = append(result, pulumi.String(action.GetValue()))
		}
	}
	return result
}
