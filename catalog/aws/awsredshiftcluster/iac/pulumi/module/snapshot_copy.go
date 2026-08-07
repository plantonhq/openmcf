package module

import (
	"github.com/pkg/errors"
	"github.com/pulumi/pulumi-aws/sdk/v7/go/aws"
	"github.com/pulumi/pulumi-aws/sdk/v7/go/aws/redshift"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

// clusterSnapshotCopy configures cross-region snapshot copy -- a cluster
// setting keyed by the cluster identifier (AWS
// EnableSnapshotCopy/DisableSnapshotCopy), not a resource with its own
// identity, which is why it is folded for the same reason as logging.
// Changing the destination region tears the configuration down and
// re-enables it against the new region (provider replace semantics).
func clusterSnapshotCopy(ctx *pulumi.Context, locals *Locals, provider *aws.Provider, createdCluster *redshift.Cluster) error {
	spec := locals.AwsRedshiftCluster.Spec
	if spec.SnapshotCopy == nil {
		return nil
	}

	args := &redshift.SnapshotCopyArgs{
		ClusterIdentifier: createdCluster.ClusterIdentifier,
		DestinationRegion: pulumi.String(spec.SnapshotCopy.DestinationRegion),
	}

	// 0 keeps the AWS defaults: 7 days for copied automated snapshots,
	// indefinite (-1) for copied manual snapshots.
	if spec.SnapshotCopy.RetentionPeriod != 0 {
		args.RetentionPeriod = pulumi.Int(int(spec.SnapshotCopy.RetentionPeriod))
	}
	if spec.SnapshotCopy.ManualSnapshotRetentionPeriod != 0 {
		args.ManualSnapshotRetentionPeriod = pulumi.Int(int(spec.SnapshotCopy.ManualSnapshotRetentionPeriod))
	}

	// Required by AWS when the cluster is KMS-encrypted: the grant lets
	// Redshift encrypt copied snapshots with a key in the destination
	// region.
	if spec.SnapshotCopy.SnapshotCopyGrantName != "" {
		args.SnapshotCopyGrantName = pulumi.String(spec.SnapshotCopy.SnapshotCopyGrantName)
	}

	if _, err := redshift.NewSnapshotCopy(ctx, "snapshot-copy", args,
		pulumi.Provider(provider), pulumi.DependsOn([]pulumi.Resource{createdCluster})); err != nil {
		return errors.Wrap(err, "failed to enable snapshot copy")
	}
	return nil
}
