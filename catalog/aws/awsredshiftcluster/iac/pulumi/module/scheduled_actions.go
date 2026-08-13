package module

import (
	"github.com/pkg/errors"
	"github.com/pulumi/pulumi-aws/sdk/v7/go/aws"
	"github.com/pulumi/pulumi-aws/sdk/v7/go/aws/redshift"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

// scheduledActions pause, resume, or resize this cluster on a cron/at
// schedule -- the nights-and-weekends cost lever. Each entry renders one
// scheduled action keyed by its (account-unique) name; the spec's
// exactly-one-arm CEL guarantees a single target arm per entry.
//
// The IAM role's trust policy must allow
// scheduler.redshift.amazonaws.com to sts:AssumeRole -- AWS validates
// the trust at create, and the provider retries briefly while a freshly
// created trust propagates.
func scheduledActions(
	ctx *pulumi.Context,
	locals *Locals,
	provider *aws.Provider,
	createdCluster *redshift.Cluster,
) error {
	spec := locals.AwsRedshiftCluster.Spec

	for _, action := range spec.ScheduledActions {
		targetAction := &redshift.ScheduledActionTargetActionArgs{}
		if action.PauseCluster {
			targetAction.PauseCluster = &redshift.ScheduledActionTargetActionPauseClusterArgs{
				ClusterIdentifier: createdCluster.ClusterIdentifier,
			}
		}
		if action.ResumeCluster {
			targetAction.ResumeCluster = &redshift.ScheduledActionTargetActionResumeClusterArgs{
				ClusterIdentifier: createdCluster.ClusterIdentifier,
			}
		}
		if action.ResizeCluster != nil {
			resizeArgs := &redshift.ScheduledActionTargetActionResizeClusterArgs{
				ClusterIdentifier: createdCluster.ClusterIdentifier,
				Classic:           pulumi.Bool(action.ResizeCluster.Classic),
			}
			// Unset members keep the cluster's current topology.
			if action.ResizeCluster.ClusterType != "" {
				resizeArgs.ClusterType = pulumi.String(action.ResizeCluster.ClusterType)
			}
			if action.ResizeCluster.NodeType != "" {
				resizeArgs.NodeType = pulumi.String(action.ResizeCluster.NodeType)
			}
			if action.ResizeCluster.NumberOfNodes != 0 {
				resizeArgs.NumberOfNodes = pulumi.Int(int(action.ResizeCluster.NumberOfNodes))
			}
			targetAction.ResizeCluster = resizeArgs
		}

		args := &redshift.ScheduledActionArgs{
			Name: pulumi.String(action.Name),

			// The spec's zero value (disabled=false) is AWS's default:
			// enabled.
			Enable: pulumi.Bool(!action.Disabled),

			Schedule:     pulumi.String(action.Schedule),
			IamRole:      pulumi.String(action.IamRoleArn.GetValue()),
			TargetAction: targetAction,
		}
		if action.Description != "" {
			args.Description = pulumi.String(action.Description)
		}
		if action.StartTime != "" {
			args.StartTime = pulumi.String(action.StartTime)
		}
		if action.EndTime != "" {
			args.EndTime = pulumi.String(action.EndTime)
		}

		if _, err := redshift.NewScheduledAction(ctx, "scheduled-action-"+action.Name, args,
			pulumi.Provider(provider), pulumi.DependsOn([]pulumi.Resource{createdCluster})); err != nil {
			return errors.Wrapf(err, "failed to create scheduled action %s", action.Name)
		}
	}
	return nil
}
