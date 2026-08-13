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
	createdCluster *neptune.Cluster,
	createdInstanceParameterGroup *neptune.ParameterGroup) ([]*neptune.ClusterInstance, map[string]*neptune.ClusterInstance, error) {
	spec := locals.AwsNeptuneCluster.Spec

	createdInstances := make([]*neptune.ClusterInstance, 0, len(spec.Instances))
	createdInstancesByName := make(map[string]*neptune.ClusterInstance, len(spec.Instances))
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

			// The cluster's resolved port, pinned on every instance.
			// Instances have no port of their own (they listen on the
			// cluster's), but the instance schema carries its own default
			// (8182) -- left unset, a cluster on any other port would read
			// back that port and fight the 8182 default with a replacement
			// diff on every update. Pinning the cluster's own attribute
			// keeps instance state converged by construction.
			Port: createdCluster.Port,

			// The manifest's immediate-vs-maintenance-window intent
			// applies to instance-scope changes too (class resizes, window
			// moves, parameter group switches) -- without this the
			// provider defers them regardless of what the manifest asked.
			ApplyImmediately: pulumi.Bool(spec.ApplyImmediately),

			// Delete-time intent forwarded from the cluster. Cluster
			// members take no instance-level final snapshot (backups are
			// cluster-storage scoped), but the flag keeps teardown intent
			// consistent on every resource.
			SkipFinalSnapshot: pulumi.Bool(spec.SkipFinalSnapshot),

			Tags: pulumi.ToStringMap(locals.AwsTags),
		}

		// Empty lets AWS spread instances across the cluster's zones --
		// the right call almost always; a pin is create-only.
		if instance.AvailabilityZone != "" {
			args.AvailabilityZone = pulumi.String(instance.AvailabilityZone)
		}

		// Instance-LEVEL parameter group (engine tunables scoped to one
		// instance): an explicit per-instance group wins; otherwise
		// instances adopt the module-managed group from
		// spec.instance_parameters when one exists; unset keeps the engine
		// default group.
		if instance.NeptuneParameterGroupName != "" {
			args.NeptuneParameterGroupName = pulumi.String(instance.NeptuneParameterGroupName)
		} else if createdInstanceParameterGroup != nil {
			args.NeptuneParameterGroupName = createdInstanceParameterGroup.Name
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
			return nil, nil, errors.Wrapf(err, "failed to create cluster instance %s", instance.Name)
		}
		createdInstances = append(createdInstances, createdInstance)
		createdInstancesByName[instance.Name] = createdInstance
	}

	return createdInstances, createdInstancesByName, nil
}
