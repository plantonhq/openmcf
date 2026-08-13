package module

import (
	"github.com/pkg/errors"
	"github.com/pulumi/pulumi-aws/sdk/v7/go/aws"
	"github.com/pulumi/pulumi-aws/sdk/v7/go/aws/redshift"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

// snapshotScheduleAssociation binds the cluster to an EXISTING snapshot
// schedule (AWS keeps exactly one schedule per cluster -- the schedule
// replaces the default automated-snapshot cadence). The schedule itself
// is an account-scoped resource shared by many clusters and is not
// created here.
func snapshotScheduleAssociation(
	ctx *pulumi.Context,
	locals *Locals,
	provider *aws.Provider,
	createdCluster *redshift.Cluster,
) error {
	spec := locals.AwsRedshiftCluster.Spec
	if spec.SnapshotScheduleIdentifier == "" {
		return nil
	}

	if _, err := redshift.NewSnapshotScheduleAssociation(ctx, "snapshot-schedule-association",
		&redshift.SnapshotScheduleAssociationArgs{
			ClusterIdentifier:  createdCluster.ClusterIdentifier,
			ScheduleIdentifier: pulumi.String(spec.SnapshotScheduleIdentifier),
		},
		pulumi.Provider(provider), pulumi.DependsOn([]pulumi.Resource{createdCluster})); err != nil {
		return errors.Wrap(err, "failed to associate snapshot schedule")
	}
	return nil
}
