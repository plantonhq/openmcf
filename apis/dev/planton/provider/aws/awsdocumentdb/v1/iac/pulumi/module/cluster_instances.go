package module

import (
	"fmt"

	"github.com/pkg/errors"
	"github.com/pulumi/pulumi-aws/sdk/v7/go/aws"
	"github.com/pulumi/pulumi-aws/sdk/v7/go/aws/docdb"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

// clusterInstances provisions the cluster's compute: one
// aws_docdb_cluster_instance per folded instances entry, keyed by the
// entry's name so adding or removing a reader is an in-place update that
// never touches the cluster or its siblings. The instance with the lowest
// promotion tier that is available becomes the writer; all others serve
// reads from the shared cluster storage.
//
// The engine version is inherited from the cluster by AWS itself --
// DocumentDB stamps every instance with the cluster's resolved version,
// and a version change rolls the instances in the same update.
func clusterInstances(ctx *pulumi.Context, locals *Locals, provider *aws.Provider,
	createdCluster *docdb.Cluster) ([]*docdb.ClusterInstance, error) {
	spec := locals.AwsDocumentDb.Spec

	createdInstances := make([]*docdb.ClusterInstance, 0, len(spec.Instances))
	for _, instance := range spec.Instances {
		args := &docdb.ClusterInstanceArgs{
			Identifier:        pulumi.Sprintf("%s-%s", locals.ClusterIdentifier, instance.Name),
			ClusterIdentifier: createdCluster.ID(),

			// "db.serverless" makes this a DocumentDB Serverless instance
			// scaling within the cluster's serverless_v2_scaling bounds;
			// any provisioned class fixes its size.
			InstanceClass: pulumi.String(instance.InstanceClass),

			// Failover priority (0 promoted first). AWS default is already
			// 0 -- forwarded unconditionally since 0 is the meaningful
			// base tier.
			PromotionTier: pulumi.Int(int(instance.PromotionTier)),

			// Copy this instance's tags onto its snapshots.
			CopyTagsToSnapshot: pulumi.Bool(instance.CopyTagsToSnapshot),

			Tags: pulumi.ToStringMap(locals.AwsTags),
		}

		// Empty lets AWS spread instances across the cluster's zones --
		// the right call almost always; a pin is create-only.
		if instance.AvailabilityZone != "" {
			args.AvailabilityZone = pulumi.String(instance.AvailabilityZone)
		}

		// Tri-state: unset keeps the AWS default (true).
		if instance.AutoMinorVersionUpgrade != nil {
			args.AutoMinorVersionUpgrade = pulumi.Bool(instance.GetAutoMinorVersionUpgrade())
		}

		// Performance Insights is instance-scoped on DocumentDB (there is
		// no cluster-level setting).
		if instance.PerformanceInsightsEnabled {
			args.EnablePerformanceInsights = pulumi.Bool(true)
		}
		if instance.PerformanceInsightsKmsKeyId.GetValue() != "" {
			args.PerformanceInsightsKmsKeyId = pulumi.String(instance.PerformanceInsightsKmsKeyId.GetValue())
		}

		// A per-instance maintenance window; empty inherits AWS
		// scheduling.
		if instance.PreferredMaintenanceWindow != "" {
			args.PreferredMaintenanceWindow = pulumi.String(instance.PreferredMaintenanceWindow)
		}

		if instance.CaCertIdentifier != "" {
			args.CaCertIdentifier = pulumi.String(instance.CaCertIdentifier)
		}

		createdInstance, err := docdb.NewClusterInstance(ctx,
			fmt.Sprintf("instance-%s", instance.Name), args, pulumi.Provider(provider))
		if err != nil {
			return nil, errors.Wrapf(err, "failed to create cluster instance %s", instance.Name)
		}
		createdInstances = append(createdInstances, createdInstance)
	}

	return createdInstances, nil
}
