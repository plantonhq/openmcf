package module

import (
	"github.com/pkg/errors"
	gcplogmetricv1alpha1 "github.com/plantonhq/planton/catalog/gcp/gcplogmetric/v1alpha1"
	"github.com/pulumi/pulumi-gcp/sdk/v9/go/gcp"
	"github.com/pulumi/pulumi-gcp/sdk/v9/go/gcp/logging"
	"github.com/pulumi/pulumi-gcp/sdk/v9/go/gcp/projects"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

// logMetric provisions the log-based metric — the bridge from log entries
// matching `filter` to a chartable Cloud Monitoring metric.
//
// `disabled` is sent EXPLICITLY on every apply: it is Optional in the
// provider, and a spec transition true -> false must reach the API rather
// than being omitted (the send-true-or-omit class silently no-ops such
// transitions — a metric that silently stays paused hides the incident it
// was built to reveal).
func logMetric(ctx *pulumi.Context, locals *Locals, gcpProvider *gcp.Provider) error {
	spec := locals.GcpLogMetric.Spec

	// Enable the Cloud Logging API so a fresh project can host the metric.
	// disable_on_destroy stays false (the provider default): tearing down
	// one metric must never disable logging for everything else in the
	// project. Matches the Terraform module.
	serviceArgs := &projects.ServiceArgs{
		Service:                  pulumi.String("logging.googleapis.com"),
		DisableDependentServices: pulumi.BoolPtr(true),
	}
	if spec.ProjectId.GetValue() != "" {
		serviceArgs.Project = pulumi.String(spec.ProjectId.GetValue())
	}
	createdProjectService, err := projects.NewService(ctx,
		"logmetric-logging.googleapis.com", serviceArgs, pulumi.Provider(gcpProvider))
	if err != nil {
		return errors.Wrap(err, "failed to enable logging.googleapis.com api")
	}

	args := &logging.MetricArgs{
		Name:   pulumi.String(locals.MetricName),
		Filter: pulumi.String(spec.Filter),
		// Explicit send — see the function comment.
		Disabled: pulumi.Bool(spec.Disabled),
	}

	if spec.BucketName.GetValue() != "" {
		args.BucketName = pulumi.String(spec.BucketName.GetValue())
	}

	if spec.Description != "" {
		args.Description = pulumi.StringPtr(spec.Description)
	}

	if spec.ValueExtractor != "" {
		args.ValueExtractor = pulumi.StringPtr(spec.ValueExtractor)
	}

	if len(spec.LabelExtractors) > 0 {
		args.LabelExtractors = pulumi.ToStringMap(spec.LabelExtractors)
	}

	if spec.MetricDescriptor != nil {
		args.MetricDescriptor = expandMetricDescriptor(spec.MetricDescriptor)
	}

	if spec.BucketOptions != nil {
		args.BucketOptions = expandBucketOptions(spec.BucketOptions)
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

	createdMetric, err := logging.NewMetric(ctx, "log-metric", args,
		pulumi.Provider(gcpProvider), pulumi.DependsOn([]pulumi.Resource{createdProjectService}))
	if err != nil {
		return errors.Wrap(err, "failed to create log metric")
	}

	ctx.Export(OpMetricName, createdMetric.Name)

	return nil
}

// expandMetricDescriptor maps the spec's descriptor onto the provider's
// metric_descriptor block.
func expandMetricDescriptor(descriptor *gcplogmetricv1alpha1.GcpLogMetricDescriptor) *logging.MetricMetricDescriptorArgs {
	descriptorArgs := &logging.MetricMetricDescriptorArgs{
		MetricKind: pulumi.String(descriptor.MetricKind),
		ValueType:  pulumi.String(descriptor.ValueType),
	}
	if descriptor.Unit != "" {
		descriptorArgs.Unit = pulumi.StringPtr(descriptor.Unit)
	}
	if descriptor.DisplayName != "" {
		descriptorArgs.DisplayName = pulumi.StringPtr(descriptor.DisplayName)
	}
	if len(descriptor.Labels) > 0 {
		labels := logging.MetricMetricDescriptorLabelArray{}
		for _, label := range descriptor.Labels {
			labelArgs := &logging.MetricMetricDescriptorLabelArgs{
				Key: pulumi.String(label.Key),
			}
			if label.Description != "" {
				labelArgs.Description = pulumi.StringPtr(label.Description)
			}
			if label.ValueType != "" {
				labelArgs.ValueType = pulumi.StringPtr(label.ValueType)
			}
			labels = append(labels, labelArgs)
		}
		descriptorArgs.Labels = labels
	}
	return descriptorArgs
}

// expandBucketOptions maps the spec's histogram layouts onto the provider's
// bucket_options block. At least one layout is present (proto-CEL-enforced,
// mirroring the provider's AtLeastOneOf).
func expandBucketOptions(bucketOptions *gcplogmetricv1alpha1.GcpLogMetricBucketOptions) *logging.MetricBucketOptionsArgs {
	bucketOptionsArgs := &logging.MetricBucketOptionsArgs{}

	if explicit := bucketOptions.ExplicitBuckets; explicit != nil {
		bounds := pulumi.Float64Array{}
		for _, bound := range explicit.Bounds {
			bounds = append(bounds, pulumi.Float64(bound))
		}
		bucketOptionsArgs.ExplicitBuckets = &logging.MetricBucketOptionsExplicitBucketsArgs{
			Bounds: bounds,
		}
	}

	if exponential := bucketOptions.ExponentialBuckets; exponential != nil {
		bucketOptionsArgs.ExponentialBuckets = &logging.MetricBucketOptionsExponentialBucketsArgs{
			NumFiniteBuckets: pulumi.Int(int(exponential.NumFiniteBuckets)),
			GrowthFactor:     pulumi.Float64(exponential.GrowthFactor),
			Scale:            pulumi.Float64(exponential.Scale),
		}
	}

	if linear := bucketOptions.LinearBuckets; linear != nil {
		bucketOptionsArgs.LinearBuckets = &logging.MetricBucketOptionsLinearBucketsArgs{
			NumFiniteBuckets: pulumi.Int(int(linear.NumFiniteBuckets)),
			Offset:           pulumi.Float64(linear.Offset),
			Width:            pulumi.Float64(linear.Width),
		}
	}

	return bucketOptionsArgs
}
