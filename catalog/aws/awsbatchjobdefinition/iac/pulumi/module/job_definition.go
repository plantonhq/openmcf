package module

import (
	"github.com/pkg/errors"
	"github.com/pulumi/pulumi-aws/sdk/v7/go/aws"
	"github.com/pulumi/pulumi-aws/sdk/v7/go/aws/batch"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

// jobDefinition registers the Batch job definition revision.
//
// Every meaningful change registers a NEW revision (revisions are immutable
// in AWS); with deregister_on_new_revision at its default the previous
// revision is deregistered so exactly one ACTIVE revision tracks this
// resource.
func jobDefinition(
	ctx *pulumi.Context,
	locals *Locals,
	provider *aws.Provider,
) (*batch.JobDefinition, error) {
	spec := locals.AwsBatchJobDefinition.Spec

	args := &batch.JobDefinitionArgs{
		// The cloud name comes from metadata.name (the catalog naming
		// basis) -- revisions register under this name in both engines.
		Name: pulumi.StringPtr(locals.AwsBatchJobDefinition.Metadata.Name),
		// Two workload arms are modeled, both of AWS type "container":
		// ECS-based container jobs (containerProperties) and Batch-on-EKS
		// pod jobs (eksProperties). The spec guarantees exactly one arm is
		// set. Multinode (nodeProperties, type "multinode") and
		// multi-container ECS (ecsProperties) remain unmodeled long-tail
		// shapes.
		Type: pulumi.String("container"),
		Tags: pulumi.ToStringMap(locals.AwsTags),
	}

	if spec.Container != nil {
		containerProperties, err := buildContainerProperties(spec)
		if err != nil {
			return nil, errors.Wrap(err, "build container properties")
		}
		args.ContainerProperties = pulumi.StringPtr(containerProperties)
	}

	if spec.Eks != nil {
		args.EksProperties = buildEksProperties(spec.Eks)
	}

	if len(spec.PlatformCapabilities) > 0 {
		args.PlatformCapabilities = pulumi.ToStringArray(spec.PlatformCapabilities)
	}
	if len(spec.Parameters) > 0 {
		args.Parameters = pulumi.ToStringMap(spec.Parameters)
	}
	if spec.SchedulingPriority > 0 {
		args.SchedulingPriority = pulumi.IntPtr(int(spec.SchedulingPriority))
	}
	args.PropagateTags = pulumi.BoolPtr(spec.PropagateTags)
	// Platform-defaulted to true; an explicit false keeps every historical
	// revision ACTIVE for out-of-band consumers.
	args.DeregisterOnNewRevision = pulumi.BoolPtr(spec.GetDeregisterOnNewRevision())

	if spec.RetryStrategy != nil {
		retry := &batch.JobDefinitionRetryStrategyArgs{}
		if spec.RetryStrategy.Attempts > 0 {
			retry.Attempts = pulumi.IntPtr(int(spec.RetryStrategy.Attempts))
		}
		if len(spec.RetryStrategy.EvaluateOnExit) > 0 {
			var conditions batch.JobDefinitionRetryStrategyEvaluateOnExitArray
			for _, condition := range spec.RetryStrategy.EvaluateOnExit {
				entry := &batch.JobDefinitionRetryStrategyEvaluateOnExitArgs{
					Action: pulumi.String(condition.Action),
				}
				if condition.OnExitCode != "" {
					entry.OnExitCode = pulumi.StringPtr(condition.OnExitCode)
				}
				if condition.OnReason != "" {
					entry.OnReason = pulumi.StringPtr(condition.OnReason)
				}
				if condition.OnStatusReason != "" {
					entry.OnStatusReason = pulumi.StringPtr(condition.OnStatusReason)
				}
				conditions = append(conditions, entry)
			}
			retry.EvaluateOnExits = conditions
		}
		args.RetryStrategy = retry
	}

	if spec.Timeout != nil {
		args.Timeout = &batch.JobDefinitionTimeoutArgs{
			AttemptDurationSeconds: pulumi.IntPtr(int(spec.Timeout.AttemptDurationSeconds)),
		}
	}

	createdJobDefinition, err := batch.NewJobDefinition(ctx, "job-definition", args, pulumi.Provider(provider))
	if err != nil {
		return nil, errors.Wrap(err, "create batch job definition")
	}

	return createdJobDefinition, nil
}
