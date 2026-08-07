package module

import (
	"github.com/pkg/errors"
	"github.com/pulumi/pulumi-aws/sdk/v7/go/aws"
	"github.com/pulumi/pulumi-aws/sdk/v7/go/aws/efs"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

type FileSystemResult struct {
	FileSystem *efs.FileSystem
}

func fileSystem(ctx *pulumi.Context, locals *Locals, provider *aws.Provider) (*FileSystemResult, error) {
	spec := locals.AwsElasticFileSystem.Spec

	args := &efs.FileSystemArgs{
		// The creation token is the file system's idempotency key AND its only
		// human-controlled physical identity (EFS has no "name" argument — the
		// console name is the Name tag). Pinning it to metadata.name keeps the
		// physical identity identical to the Terraform module's; left unset,
		// Pulumi would auto-name with a random suffix and the two engines
		// would silently diverge.
		CreationToken: pulumi.String(locals.AwsElasticFileSystem.Metadata.Name),
		Tags:          pulumi.ToStringMap(locals.AwsTags),
	}

	// Encryption at rest. ForceNew — it cannot be enabled after creation.
	if spec.Encrypted {
		args.Encrypted = pulumi.Bool(true)
	}

	// Customer-managed KMS key (requires encrypted; CEL-enforced). ForceNew.
	if spec.KmsKeyId.GetValue() != "" {
		args.KmsKeyId = pulumi.StringPtr(spec.KmsKeyId.GetValue())
	}

	// Performance mode (generalPurpose or maxIO). ForceNew. Empty keeps the
	// AWS default (generalPurpose).
	if spec.PerformanceMode != "" {
		args.PerformanceMode = pulumi.StringPtr(spec.PerformanceMode)
	}

	// Throughput mode (bursting, provisioned, or elastic). Mutable, but AWS
	// enforces a 24-hour cooldown between throughput-mode changes.
	if spec.ThroughputMode != "" {
		args.ThroughputMode = pulumi.StringPtr(spec.ThroughputMode)
	}

	// Provisioned throughput (only valid with throughput_mode = "provisioned";
	// the coupling is CEL-enforced at validation time).
	if spec.ProvisionedThroughputInMibps > 0 {
		args.ProvisionedThroughputInMibps = pulumi.Float64Ptr(spec.ProvisionedThroughputInMibps)
	}

	// One Zone storage (single AZ). ForceNew.
	if spec.AvailabilityZoneName != "" {
		args.AvailabilityZoneName = pulumi.StringPtr(spec.AvailabilityZoneName)
	}

	// Replication-overwrite protection. AWS defaults to ENABLED; only send an
	// explicit value so unset stays indistinguishable from the AWS default.
	// DISABLED is required before this file system can be targeted as a
	// replication destination (or modified/deleted after having been one).
	if spec.ReplicationOverwriteProtection != "" {
		args.Protection = &efs.FileSystemProtectionArgs{
			ReplicationOverwrite: pulumi.StringPtr(spec.ReplicationOverwriteProtection),
		}
	}

	// Lifecycle policies — the provider models each transition rule as its own
	// lifecycle_policy element (up to 3: IA, Archive, back-to-primary).
	var lifecyclePolicies efs.FileSystemLifecyclePolicyArray
	if spec.TransitionToIa != "" {
		lifecyclePolicies = append(lifecyclePolicies, &efs.FileSystemLifecyclePolicyArgs{
			TransitionToIa: pulumi.StringPtr(spec.TransitionToIa),
		})
	}
	if spec.TransitionToArchive != "" {
		lifecyclePolicies = append(lifecyclePolicies, &efs.FileSystemLifecyclePolicyArgs{
			TransitionToArchive: pulumi.StringPtr(spec.TransitionToArchive),
		})
	}
	if spec.TransitionToPrimaryStorageClass != "" {
		lifecyclePolicies = append(lifecyclePolicies, &efs.FileSystemLifecyclePolicyArgs{
			TransitionToPrimaryStorageClass: pulumi.StringPtr(spec.TransitionToPrimaryStorageClass),
		})
	}
	if len(lifecyclePolicies) > 0 {
		args.LifecyclePolicies = lifecyclePolicies
	}

	// Stable logical name; the cloud identity travels through CreationToken
	// and the Name tag, never through the Pulumi resource name.
	createdFs, err := efs.NewFileSystem(ctx, "file-system", args, pulumi.Provider(provider))
	if err != nil {
		return nil, errors.Wrap(err, "failed to create efs file system")
	}

	return &FileSystemResult{FileSystem: createdFs}, nil
}
