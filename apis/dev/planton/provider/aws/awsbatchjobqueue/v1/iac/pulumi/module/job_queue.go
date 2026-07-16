package module

import (
	"github.com/pkg/errors"
	"github.com/pulumi/pulumi-aws/sdk/v7/go/aws"
	"github.com/pulumi/pulumi-aws/sdk/v7/go/aws/batch"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

// jobQueue creates the Batch job queue mapped onto its compute environments
// in preference order.
func jobQueue(
	ctx *pulumi.Context,
	locals *Locals,
	provider *aws.Provider,
) (*batch.JobQueue, error) {
	spec := locals.AwsBatchJobQueue.Spec

	// The compute-environment mapping is the queue's core: the scheduler
	// tries the lowest `order` first, which is how Spot-first-overflow-to-
	// On-Demand (and blue/green environment replacement) is expressed.
	var ceOrders batch.JobQueueComputeEnvironmentOrderArray
	for _, entry := range spec.ComputeEnvironmentOrder {
		ceOrders = append(ceOrders, &batch.JobQueueComputeEnvironmentOrderArgs{
			Order:              pulumi.Int(entry.Order),
			ComputeEnvironment: pulumi.String(entry.ComputeEnvironment.GetValue()),
		})
	}

	args := &batch.JobQueueArgs{
		// The cloud name comes from metadata.name (the catalog naming
		// basis) -- set explicitly so both engines create the same queue
		// name and Pulumi never auto-names.
		Name:                     pulumi.StringPtr(locals.AwsBatchJobQueue.Metadata.Name),
		Priority:                 pulumi.Int(spec.Priority),
		State:                    pulumi.String(spec.GetState()),
		ComputeEnvironmentOrders: ceOrders,
		Tags:                     pulumi.ToStringMap(locals.AwsTags),
	}

	// AWS quirk carried by the spec comment: once set, the scheduling
	// policy can be replaced but never removed from a live queue.
	if spec.SchedulingPolicy.GetValue() != "" {
		args.SchedulingPolicyArn = pulumi.StringPtr(spec.SchedulingPolicy.GetValue())
	}

	if len(spec.JobStateTimeLimitActions) > 0 {
		var actions batch.JobQueueJobStateTimeLimitActionArray
		for _, action := range spec.JobStateTimeLimitActions {
			actions = append(actions, &batch.JobQueueJobStateTimeLimitActionArgs{
				Action:         pulumi.String(action.Action),
				MaxTimeSeconds: pulumi.Int(action.MaxTimeSeconds),
				Reason:         pulumi.String(action.Reason),
				State:          pulumi.String(action.State),
			})
		}
		args.JobStateTimeLimitActions = actions
	}

	createdQueue, err := batch.NewJobQueue(ctx, "job-queue", args, pulumi.Provider(provider))
	if err != nil {
		return nil, errors.Wrap(err, "create batch job queue")
	}

	return createdQueue, nil
}
