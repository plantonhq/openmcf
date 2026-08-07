package module

import (
	"github.com/pkg/errors"
	"github.com/pulumi/pulumi-aws/sdk/v7/go/aws"
	"github.com/pulumi/pulumi-aws/sdk/v7/go/aws/efs"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

// ReplicationResult carries the replication configuration when one was
// requested; nil when the spec omits replication.
type ReplicationResult struct {
	ReplicationConfiguration *efs.ReplicationConfiguration
}

func replication(ctx *pulumi.Context, locals *Locals, provider *aws.Provider, fs *efs.FileSystem) (*ReplicationResult, error) {
	spec := locals.AwsElasticFileSystem.Spec

	if spec.Replication == nil {
		return &ReplicationResult{}, nil
	}

	// EFS replication is one-per-file-system and create-time immutable: the
	// provider has no Update — any destination change replaces the
	// configuration (the destination file system itself survives that
	// replacement; deleting it additionally requires its own
	// replication_overwrite_protection to be DISABLED).
	destination := &efs.ReplicationConfigurationDestinationArgs{}

	if spec.Replication.DestinationRegion != "" {
		destination.Region = pulumi.StringPtr(spec.Replication.DestinationRegion)
	}

	// An AZ destination creates the replica as a One Zone file system — the
	// cheaper DR shape. At least one of region/AZ is CEL-enforced.
	if spec.Replication.DestinationAvailabilityZoneName != "" {
		destination.AvailabilityZoneName = pulumi.StringPtr(spec.Replication.DestinationAvailabilityZoneName)
	}

	if spec.Replication.DestinationKmsKeyId.GetValue() != "" {
		destination.KmsKeyId = pulumi.StringPtr(spec.Replication.DestinationKmsKeyId.GetValue())
	}

	// Replicate into an existing file system instead of letting AWS mint one.
	if spec.Replication.DestinationFileSystemId.GetValue() != "" {
		destination.FileSystemId = pulumi.StringPtr(spec.Replication.DestinationFileSystemId.GetValue())
	}

	createdReplication, err := efs.NewReplicationConfiguration(ctx, "replication", &efs.ReplicationConfigurationArgs{
		SourceFileSystemId: fs.ID(),
		Destination:        destination,
	}, pulumi.Provider(provider))
	if err != nil {
		return nil, errors.Wrap(err, "failed to create efs replication configuration")
	}

	return &ReplicationResult{ReplicationConfiguration: createdReplication}, nil
}
