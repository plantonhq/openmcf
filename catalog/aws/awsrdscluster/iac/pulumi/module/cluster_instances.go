package module

import (
	"fmt"

	"github.com/pkg/errors"
	"github.com/pulumi/pulumi-aws/sdk/v7/go/aws"
	"github.com/pulumi/pulumi-aws/sdk/v7/go/aws/rds"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

// clusterInstances provisions the cluster's compute: one
// aws_rds_cluster_instance per folded instances entry, keyed by the
// entry's name so adding or removing a reader is an in-place update that
// never touches the cluster or its siblings. The instance with the lowest
// promotion tier that is available becomes the writer; all others serve
// reads from the shared cluster storage.
//
// Engine and version are inherited from the cluster resource by output
// reference -- so an unpinned cluster (AWS-default version) still stamps
// its instances with the resolved version, and a version change rolls the
// instances in the same update.
func clusterInstances(ctx *pulumi.Context, locals *Locals, provider *aws.Provider,
	createdCluster *rds.Cluster) ([]*rds.ClusterInstance, error) {
	spec := locals.AwsRdsCluster.Spec

	createdInstances := make([]*rds.ClusterInstance, 0, len(spec.Instances))
	for _, instance := range spec.Instances {
		args := &rds.ClusterInstanceArgs{
			Identifier:        pulumi.Sprintf("%s-%s", locals.ClusterIdentifier, instance.Name),
			ClusterIdentifier: createdCluster.ID(),
			// The instance arg is the provider's EngineType enum; the
			// cluster exports a plain string -- bridge with a typed apply.
			Engine: createdCluster.Engine.ApplyT(func(engine string) rds.EngineType {
				return rds.EngineType(engine)
			}).(rds.EngineTypeOutput),
			EngineVersion: createdCluster.EngineVersion,

			// "db.serverless" makes this an Aurora Serverless v2 instance
			// scaling within the cluster's serverless_v2_scaling bounds;
			// any provisioned class fixes its size.
			InstanceClass: pulumi.String(instance.InstanceClass),

			// Failover priority (0 promoted first). AWS default is already
			// 0 -- forwarded unconditionally since 0 is the meaningful
			// base tier.
			PromotionTier: pulumi.Int(int(instance.PromotionTier)),

			PubliclyAccessible: pulumi.Bool(instance.PubliclyAccessible),

			Tags: pulumi.ToStringMap(locals.AwsTags),
		}

		// Empty lets AWS spread instances across the cluster's zones --
		// the right call almost always; a pin is create-only.
		if instance.AvailabilityZone != "" {
			args.AvailabilityZone = pulumi.String(instance.AvailabilityZone)
		}

		// Instance-LEVEL parameter group (engine tunables scoped to one
		// instance); the cluster-level group lives on the cluster resource.
		if instance.DbParameterGroupName != "" {
			args.DbParameterGroupName = pulumi.String(instance.DbParameterGroupName)
		}

		// Tri-state: unset keeps the AWS default (true).
		if instance.AutoMinorVersionUpgrade != nil {
			args.AutoMinorVersionUpgrade = pulumi.Bool(instance.GetAutoMinorVersionUpgrade())
		}

		// Tri-state: unset inherits the cluster-level Performance Insights
		// posture; an explicit value overrides it for this instance only.
		if instance.PerformanceInsightsEnabled != nil {
			args.PerformanceInsightsEnabled = pulumi.Bool(instance.GetPerformanceInsightsEnabled())
		}
		if instance.PerformanceInsightsKmsKeyId.GetValue() != "" {
			args.PerformanceInsightsKmsKeyId = pulumi.String(instance.PerformanceInsightsKmsKeyId.GetValue())
		}
		if instance.PerformanceInsightsRetentionPeriod != 0 {
			args.PerformanceInsightsRetentionPeriod = pulumi.Int(int(instance.PerformanceInsightsRetentionPeriod))
		}

		// Per-instance Enhanced Monitoring cadence, publishing through the
		// instance's own role when given, else the cluster spec's
		// monitoring role (AWS requires a role whenever the interval is
		// set -- the spec's CEL guarantees the cluster-level pairing).
		if instance.MonitoringInterval != 0 {
			args.MonitoringInterval = pulumi.Int(int(instance.MonitoringInterval))
			if instance.MonitoringRoleArn.GetValue() != "" {
				args.MonitoringRoleArn = pulumi.String(instance.MonitoringRoleArn.GetValue())
			} else if spec.MonitoringRoleArn.GetValue() != "" {
				args.MonitoringRoleArn = pulumi.String(spec.MonitoringRoleArn.GetValue())
			}
		}

		// Per-instance windows: empty inherits AWS scheduling -- stagger
		// these so readers never back up or patch simultaneously.
		if instance.PreferredBackupWindow != "" {
			args.PreferredBackupWindow = pulumi.String(instance.PreferredBackupWindow)
		}
		if instance.PreferredMaintenanceWindow != "" {
			args.PreferredMaintenanceWindow = pulumi.String(instance.PreferredMaintenanceWindow)
		}

		if instance.CaCertIdentifier != "" {
			args.CaCertIdentifier = pulumi.String(instance.CaCertIdentifier)
		}
		args.CopyTagsToSnapshot = pulumi.Bool(instance.CopyTagsToSnapshot)
		if instance.ApplyImmediately {
			args.ApplyImmediately = pulumi.Bool(true)
		}

		createdInstance, err := rds.NewClusterInstance(ctx,
			fmt.Sprintf("instance-%s", instance.Name), args, pulumi.Provider(provider))
		if err != nil {
			return nil, errors.Wrapf(err, "failed to create cluster instance %s", instance.Name)
		}
		createdInstances = append(createdInstances, createdInstance)
	}

	return createdInstances, nil
}
