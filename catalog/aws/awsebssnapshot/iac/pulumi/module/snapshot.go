package module

import (
	"github.com/pkg/errors"
	awsebssnapshotv1alpha1 "github.com/plantonhq/planton/catalog/aws/awsebssnapshot/v1alpha1"
	"github.com/pulumi/pulumi-aws/sdk/v7/go/aws"
	"github.com/pulumi/pulumi-aws/sdk/v7/go/aws/ebs"
	"github.com/pulumi/pulumi-aws/sdk/v7/go/aws/ec2"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

// snapshot creates the EBS snapshot (whichever source arm), its fast
// snapshot restore and share satellites, and exports outputs.
//
// Lifecycle facts the render below depends on:
//   - the three arms are three provider resources (ebs.Snapshot /
//     ebs.SnapshotCopy / ebs.SnapshotImport); exactly one exists per
//     the spec's union CEL, and all three expose the same downstream
//     surface (id/arn/owner/size);
//   - only the volume arm is importable at the provider - the copy
//     and import resources ship no importer (declared honestly in the
//     import catalog);
//   - StorageTier / PermanentRestore / TemporaryRestoreDays update in
//     place on all arms; every source field replaces;
//   - fast snapshot restore bills per zone-hour while enabled - one
//     resource per zone;
//   - createVolumePermission grants are per-account resources;
//     encrypted snapshots additionally need the KMS key shared
//     out-of-band.
func snapshot(ctx *pulumi.Context, locals *Locals, provider *aws.Provider) error {
	spec := locals.Spec

	var snapshotId pulumi.StringOutput
	var snapshotArn pulumi.StringOutput
	var ownerId pulumi.StringOutput
	var volumeSize pulumi.IntOutput

	switch {
	case spec.CopyFrom != nil:
		// The COPY arm: copy an existing snapshot (same- or
		// cross-region), optionally re-encrypting.
		args := &ebs.SnapshotCopyArgs{
			SourceSnapshotId: pulumi.String(spec.CopyFrom.SourceSnapshotId.GetValue()),
			SourceRegion:     pulumi.String(spec.CopyFrom.SourceRegion),
			Tags:             pulumi.ToStringMap(locals.AwsTags),
		}
		if spec.Description != "" {
			args.Description = pulumi.String(spec.Description)
		}
		if spec.CopyFrom.Encrypted {
			args.Encrypted = pulumi.Bool(true)
		}
		if spec.CopyFrom.KmsKeyId.GetValue() != "" {
			args.KmsKeyId = pulumi.String(spec.CopyFrom.KmsKeyId.GetValue())
		}
		if spec.CopyFrom.CompletionDurationMinutes > 0 {
			args.CompletionDurationMinutes = pulumi.Int(int(spec.CopyFrom.CompletionDurationMinutes))
		}
		applyTieringToCopy(args, spec)

		createdCopy, err := ebs.NewSnapshotCopy(ctx, "snapshot", args, pulumi.Provider(provider))
		if err != nil {
			return errors.Wrap(err, "create snapshot copy")
		}
		snapshotId = createdCopy.ID().ToStringOutput()
		snapshotArn = createdCopy.Arn
		ownerId = createdCopy.OwnerId
		volumeSize = createdCopy.VolumeSize

	case spec.ImportFrom != nil:
		// The IMPORT arm: build the snapshot from a disk image via VM
		// Import/Export.
		diskContainer := ebs.SnapshotImportDiskContainerArgs{
			Format: pulumi.String(spec.ImportFrom.DiskContainer.Format),
		}
		if spec.ImportFrom.DiskContainer.Description != "" {
			diskContainer.Description = pulumi.String(spec.ImportFrom.DiskContainer.Description)
		}
		if spec.ImportFrom.DiskContainer.Url != "" {
			diskContainer.Url = pulumi.String(spec.ImportFrom.DiskContainer.Url)
		}
		if spec.ImportFrom.DiskContainer.S3Bucket != "" {
			diskContainer.UserBucket = &ebs.SnapshotImportDiskContainerUserBucketArgs{
				S3Bucket: pulumi.String(spec.ImportFrom.DiskContainer.S3Bucket),
				S3Key:    pulumi.String(spec.ImportFrom.DiskContainer.S3Key),
			}
		}

		args := &ebs.SnapshotImportArgs{
			DiskContainer: diskContainer,
			Tags:          pulumi.ToStringMap(locals.AwsTags),
		}
		if spec.Description != "" {
			args.Description = pulumi.String(spec.Description)
		}
		if spec.ImportFrom.RoleName != "" {
			args.RoleName = pulumi.String(spec.ImportFrom.RoleName)
		}
		if spec.ImportFrom.Encrypted {
			args.Encrypted = pulumi.Bool(true)
		}
		if spec.ImportFrom.KmsKeyId.GetValue() != "" {
			args.KmsKeyId = pulumi.String(spec.ImportFrom.KmsKeyId.GetValue())
		}
		applyTieringToImport(args, spec)

		createdImport, err := ebs.NewSnapshotImport(ctx, "snapshot", args, pulumi.Provider(provider))
		if err != nil {
			return errors.Wrap(err, "create snapshot import")
		}
		snapshotId = createdImport.ID().ToStringOutput()
		snapshotArn = createdImport.Arn
		ownerId = createdImport.OwnerId
		volumeSize = createdImport.VolumeSize

	default:
		// The VOLUME arm: snapshot a live volume.
		args := &ebs.SnapshotArgs{
			VolumeId: pulumi.String(spec.VolumeId.GetValue()),
			Tags:     pulumi.ToStringMap(locals.AwsTags),
		}
		if spec.Description != "" {
			args.Description = pulumi.String(spec.Description)
		}
		applyTieringToSnapshot(args, spec)

		createdSnapshot, err := ebs.NewSnapshot(ctx, "snapshot", args, pulumi.Provider(provider))
		if err != nil {
			return errors.Wrap(err, "create snapshot")
		}
		snapshotId = createdSnapshot.ID().ToStringOutput()
		snapshotArn = createdSnapshot.Arn
		ownerId = createdSnapshot.OwnerId
		volumeSize = createdSnapshot.VolumeSize
	}

	// Fast snapshot restore, one resource per availability zone.
	// Billed per zone-hour while enabled.
	for _, zone := range spec.FastRestoreAvailabilityZones {
		if _, err := ebs.NewFastSnapshotRestore(ctx, "fast-restore-"+zone, &ebs.FastSnapshotRestoreArgs{
			AvailabilityZone: pulumi.String(zone),
			SnapshotId:       snapshotId,
		}, pulumi.Provider(provider)); err != nil {
			return errors.Wrapf(err, "fast snapshot restore %s", zone)
		}
	}

	// createVolumePermission grants, one resource per account.
	for _, accountId := range spec.ShareWithAccountIds {
		if _, err := ec2.NewSnapshotCreateVolumePermission(ctx, "share-"+accountId, &ec2.SnapshotCreateVolumePermissionArgs{
			SnapshotId: snapshotId,
			AccountId:  pulumi.String(accountId),
		}, pulumi.Provider(provider)); err != nil {
			return errors.Wrapf(err, "share with %s", accountId)
		}
	}

	ctx.Export(OpSnapshotId, snapshotId)
	ctx.Export(OpSnapshotArn, snapshotArn)
	ctx.Export(OpOwnerId, ownerId)
	// The size is an int at the provider; exported as a string to
	// match the outputs contract (string-typed observable state).
	ctx.Export(OpVolumeSizeGb, pulumi.Sprintf("%d", volumeSize))
	return nil
}

// The tiering dials are identical across the three arm types; the
// bridge just generates three distinct args structs.

func applyTieringToSnapshot(args *ebs.SnapshotArgs, spec *awsebssnapshotv1alpha1.AwsEbsSnapshotSpec) {
	if spec.StorageTier != "" {
		args.StorageTier = pulumi.String(spec.StorageTier)
	}
	if spec.PermanentRestore {
		args.PermanentRestore = pulumi.Bool(true)
	}
	if spec.TemporaryRestoreDays > 0 {
		args.TemporaryRestoreDays = pulumi.Int(int(spec.TemporaryRestoreDays))
	}
}

func applyTieringToCopy(args *ebs.SnapshotCopyArgs, spec *awsebssnapshotv1alpha1.AwsEbsSnapshotSpec) {
	if spec.StorageTier != "" {
		args.StorageTier = pulumi.String(spec.StorageTier)
	}
	if spec.PermanentRestore {
		args.PermanentRestore = pulumi.Bool(true)
	}
	if spec.TemporaryRestoreDays > 0 {
		args.TemporaryRestoreDays = pulumi.Int(int(spec.TemporaryRestoreDays))
	}
}

func applyTieringToImport(args *ebs.SnapshotImportArgs, spec *awsebssnapshotv1alpha1.AwsEbsSnapshotSpec) {
	if spec.StorageTier != "" {
		args.StorageTier = pulumi.String(spec.StorageTier)
	}
	if spec.PermanentRestore {
		args.PermanentRestore = pulumi.Bool(true)
	}
	if spec.TemporaryRestoreDays > 0 {
		args.TemporaryRestoreDays = pulumi.Int(int(spec.TemporaryRestoreDays))
	}
}
