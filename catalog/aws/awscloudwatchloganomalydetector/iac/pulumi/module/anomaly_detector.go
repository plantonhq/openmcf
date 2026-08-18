package module

import (
	"github.com/pkg/errors"
	"github.com/pulumi/pulumi-aws/sdk/v7/go/aws"
	"github.com/pulumi/pulumi-aws/sdk/v7/go/aws/cloudwatch"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

// anomalyDetector creates the CloudWatch Logs anomaly detector and
// exports outputs.
//
// Lifecycle facts the render below depends on:
//   - Enabled is REQUIRED by the provider and always rendered - false
//     pauses analysis without losing the trained model;
//   - KmsKeyId replaces the detector on change (AWS cannot re-encrypt
//     a trained model in place); everything else updates in place;
//   - anomaly_visibility_time is presence-typed in the spec so the 7
//     and 90 boundary values are expressible; unset keeps AWS's
//     default (21);
//   - the provider treats an AccessDeniedException on read as
//     "detector gone" and drops it from state - a permissions
//     regression can look like a deleted detector.
func anomalyDetector(ctx *pulumi.Context, locals *Locals, provider *aws.Provider) error {
	spec := locals.Spec

	logGroupArns := pulumi.StringArray{}
	for _, arn := range spec.LogGroupArns {
		logGroupArns = append(logGroupArns, pulumi.String(arn.GetValue()))
	}

	args := &cloudwatch.LogAnomalyDetectorArgs{
		LogGroupArnLists: logGroupArns,
		Enabled:          pulumi.Bool(spec.Enabled),
		Tags:             pulumi.ToStringMap(locals.AwsTags),
	}

	if spec.DetectorName != "" {
		args.DetectorName = pulumi.String(spec.DetectorName)
	}
	if spec.EvaluationFrequency != "" {
		args.EvaluationFrequency = pulumi.String(spec.EvaluationFrequency)
	}
	if spec.FilterPattern != "" {
		args.FilterPattern = pulumi.String(spec.FilterPattern)
	}
	if spec.KmsKeyId.GetValue() != "" {
		args.KmsKeyId = pulumi.String(spec.KmsKeyId.GetValue())
	}
	if spec.AnomalyVisibilityTime != nil {
		args.AnomalyVisibilityTime = pulumi.Int(int(*spec.AnomalyVisibilityTime))
	}

	createdDetector, err := cloudwatch.NewLogAnomalyDetector(ctx, "anomaly_detector", args, pulumi.Provider(provider))
	if err != nil {
		return errors.Wrap(err, "create anomaly detector")
	}

	ctx.Export(OpAnomalyDetectorArn, createdDetector.Arn)
	return nil
}
