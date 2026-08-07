package module

import (
	"fmt"

	"github.com/pkg/errors"
	"github.com/pulumi/pulumi-aws/sdk/v7/go/aws"
	"github.com/pulumi/pulumi-aws/sdk/v7/go/aws/neptune"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

// clusterInstances provisions the cluster's compute: one
// aws_neptune_cluster_instance per folded instances entry, keyed by the
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
	createdCluster *neptune.Cluster) ([]*neptune.ClusterInstance, error) {
	spec := locals.AwsNeptuneCluster.Spec

	createdInstances := make([]*neptune.ClusterInstance, 0, len(spec.Instances))
	for _, instance := range spec.Instances {
		args := &neptune.ClusterInstanceArgs{
			Identifier:        pulumi.Sprintf("%s-%s", locals.ClusterIdentifier, instance.Name),
			ClusterIdentifier: createdCluster.ID(),
			Engine:            pulumi.String("neptune"),
			EngineVersion:     createdCluster.EngineVersion,

			// "db.serverless" makes this a Neptune Serverless instance
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
		// instance); the cluster-level group lives on the cluster
		// resource.
		if instance.NeptuneParameterGroupName != "" {
			args.NeptuneParameterGroupName = pulumi.String(instance.NeptuneParameterGroupName)
		}

		// Tri-state: unset keeps the AWS default (true).
		if instance.AutoMinorVersionUpgrade != nil {
			args.AutoMinorVersionUpgrade = pulumi.Bool(instance.GetAutoMinorVersionUpgrade())
		}

		// A per-instance maintenance window; empty inherits AWS
		// scheduling.
		if instance.PreferredMaintenanceWindow != "" {
			args.PreferredMaintenanceWindow = pulumi.String(instance.PreferredMaintenanceWindow)
		}

		createdInstance, err := neptune.NewClusterInstance(ctx,
			fmt.Sprintf("instance-%s", instance.Name), args, pulumi.Provider(provider))
		if err != nil {
			return nil, errors.Wrapf(err, "failed to create cluster instance %s", instance.Name)
		}
		createdInstances = append(createdInstances, createdInstance)
	}

	return createdInstances, nil
}
