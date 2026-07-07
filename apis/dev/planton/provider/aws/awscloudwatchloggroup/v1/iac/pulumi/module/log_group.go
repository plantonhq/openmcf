package module

import (
	"encoding/json"
	"fmt"

	"github.com/pkg/errors"
	"github.com/pulumi/pulumi-aws/sdk/v7/go/aws"
	"github.com/pulumi/pulumi-aws/sdk/v7/go/aws/cloudwatch"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

type LogGroupResult struct {
	LogGroupArn  pulumi.StringOutput
	LogGroupName pulumi.StringOutput
}

func logGroup(ctx *pulumi.Context, locals *Locals, provider *aws.Provider) (*LogGroupResult, error) {
	spec := locals.AwsCloudwatchLogGroup.Spec

	args := &cloudwatch.LogGroupArgs{
		// The log group's cloud name is the resource's metadata.name — the same
		// basis the Terraform module uses. Setting it explicitly (instead of
		// relying on Pulumi auto-naming, which appends a random suffix) keeps
		// the physical name identical across both IaC engines and predictable
		// for services that address the group by name.
		Name: pulumi.String(locals.AwsCloudwatchLogGroup.Metadata.Name),
		Tags: pulumi.ToStringMap(locals.AwsTags),
	}

	// Retention: 0 means never expire (the default). Only set the field when a
	// non-zero retention is configured, since Pulumi treats nil as "do not manage".
	if spec.RetentionInDays > 0 {
		args.RetentionInDays = pulumi.IntPtr(int(spec.RetentionInDays))
	}

	// KMS encryption: customer-managed key for log data at rest. Associating the
	// key is an in-place update on the log group (never a replacement).
	if spec.KmsKeyId != nil {
		args.KmsKeyId = pulumi.StringPtr(spec.KmsKeyId.GetValue())
	}

	// Log group class: STANDARD (default), INFREQUENT_ACCESS, or DELIVERY.
	// Only set when explicitly specified; omitting lets AWS default to STANDARD.
	if spec.LogGroupClass != "" {
		args.LogGroupClass = pulumi.StringPtr(spec.LogGroupClass)
	}

	// Deletion protection: applied by the provider through a separate
	// PutLogGroupDeletionProtection call after create. Only set when enabled so
	// unset stays indistinguishable from AWS's default (unprotected).
	if spec.DeletionProtectionEnabled {
		args.DeletionProtectionEnabled = pulumi.BoolPtr(true)
	}

	createdLogGroup, err := cloudwatch.NewLogGroup(
		ctx,
		locals.AwsCloudwatchLogGroup.Metadata.Name,
		args,
		pulumi.Provider(provider),
	)
	if err != nil {
		return nil, errors.Wrap(err, "failed to create cloudwatch log group")
	}

	// ---------------------------------------------------------------------------
	// Folded satellites. All four resource types below are log-group-scoped:
	// they are keyed by the group's name, share its lifecycle, and are not
	// independently referenceable — which is why they live on this spec instead
	// of being separate kinds. Each depends on the created group so destroys
	// unwind in the right order.
	// ---------------------------------------------------------------------------

	// Metric filters: one provider resource per named filter (many-per-group).
	for _, filter := range spec.MetricFilters {
		transformationArgs := cloudwatch.LogMetricFilterMetricTransformationArgs{
			Name:      pulumi.String(filter.Transformation.MetricName),
			Namespace: pulumi.String(filter.Transformation.MetricNamespace),
			Value:     pulumi.String(filter.Transformation.MetricValue),
		}
		// default_value is a genuine tri-state: nil means "publish nothing for
		// non-matching periods". AWS forbids combining it with dimensions (the
		// spec CEL enforces that before we ever get here).
		if filter.Transformation.DefaultValue != nil {
			transformationArgs.DefaultValue = pulumi.StringPtr(
				fmt.Sprintf("%g", filter.Transformation.GetDefaultValue()))
		}
		if len(filter.Transformation.Dimensions) > 0 {
			transformationArgs.Dimensions = pulumi.ToStringMap(filter.Transformation.Dimensions)
		}
		if filter.Transformation.Unit != "" {
			transformationArgs.Unit = pulumi.StringPtr(filter.Transformation.Unit)
		}

		filterArgs := &cloudwatch.LogMetricFilterArgs{
			Name:         pulumi.String(filter.Name),
			LogGroupName: createdLogGroup.Name,
			// The provider requires pattern (empty string = match every event).
			Pattern:              pulumi.String(filter.Pattern),
			MetricTransformation: transformationArgs,
		}
		if filter.ApplyOnTransformedLogs {
			filterArgs.ApplyOnTransformedLogs = pulumi.BoolPtr(true)
		}

		if _, err := cloudwatch.NewLogMetricFilter(
			ctx,
			fmt.Sprintf("metric-filter-%s", filter.Name),
			filterArgs,
			pulumi.Provider(provider),
			pulumi.DependsOn([]pulumi.Resource{createdLogGroup}),
		); err != nil {
			return nil, errors.Wrapf(err, "failed to create metric filter %s", filter.Name)
		}
	}

	// Subscription filters: at most two per group (AWS limit, CEL-enforced).
	for _, filter := range spec.SubscriptionFilters {
		subscriptionArgs := &cloudwatch.LogSubscriptionFilterArgs{
			Name:           pulumi.String(filter.Name),
			LogGroup:       createdLogGroup.Name,
			DestinationArn: pulumi.String(filter.DestinationArn.GetValue()),
			// The provider requires filter_pattern (empty = deliver everything).
			FilterPattern: pulumi.String(filter.FilterPattern),
		}
		// role_arn is required for Kinesis/Firehose destinations (CloudWatch Logs
		// assumes it to put records); Lambda destinations authorize through a
		// Lambda permission instead, so the field stays optional.
		if filter.RoleArn.GetValue() != "" {
			subscriptionArgs.RoleArn = pulumi.StringPtr(filter.RoleArn.GetValue())
		}
		if filter.Distribution != "" {
			subscriptionArgs.Distribution = pulumi.StringPtr(filter.Distribution)
		}
		if len(filter.EmitSystemFields) > 0 {
			subscriptionArgs.EmitSystemFields = pulumi.ToStringArray(filter.EmitSystemFields)
		}
		if filter.ApplyOnTransformedLogs {
			subscriptionArgs.ApplyOnTransformedLogs = pulumi.BoolPtr(true)
		}

		if _, err := cloudwatch.NewLogSubscriptionFilter(
			ctx,
			fmt.Sprintf("subscription-filter-%s", filter.Name),
			subscriptionArgs,
			pulumi.Provider(provider),
			pulumi.DependsOn([]pulumi.Resource{createdLogGroup}),
		); err != nil {
			return nil, errors.Wrapf(err, "failed to create subscription filter %s", filter.Name)
		}
	}

	// Data protection policy: a single group-scoped policy document (PII
	// audit + masking). The spec models it as a Struct; the provider argument
	// wants the JSON document text.
	if spec.DataProtectionPolicy != nil {
		policyJson, err := json.Marshal(spec.DataProtectionPolicy.AsMap())
		if err != nil {
			return nil, errors.Wrap(err, "failed to marshal data protection policy")
		}
		if _, err := cloudwatch.NewLogDataProtectionPolicy(
			ctx,
			"data-protection-policy",
			&cloudwatch.LogDataProtectionPolicyArgs{
				LogGroupName:   createdLogGroup.Name,
				PolicyDocument: pulumi.String(policyJson),
			},
			pulumi.Provider(provider),
			pulumi.DependsOn([]pulumi.Resource{createdLogGroup}),
		); err != nil {
			return nil, errors.Wrap(err, "failed to create data protection policy")
		}
	}

	// Field index policy: a single group-scoped policy listing the log fields
	// to index for Logs Insights acceleration.
	if spec.FieldIndexPolicy != nil {
		policyJson, err := json.Marshal(spec.FieldIndexPolicy.AsMap())
		if err != nil {
			return nil, errors.Wrap(err, "failed to marshal field index policy")
		}
		if _, err := cloudwatch.NewLogIndexPolicy(
			ctx,
			"field-index-policy",
			&cloudwatch.LogIndexPolicyArgs{
				LogGroupName:   createdLogGroup.Name,
				PolicyDocument: pulumi.String(policyJson),
			},
			pulumi.Provider(provider),
			pulumi.DependsOn([]pulumi.Resource{createdLogGroup}),
		); err != nil {
			return nil, errors.Wrap(err, "failed to create field index policy")
		}
	}

	return &LogGroupResult{
		LogGroupArn:  createdLogGroup.Arn,
		LogGroupName: createdLogGroup.Name,
	}, nil
}
