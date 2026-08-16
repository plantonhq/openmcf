package module

import (
	"fmt"

	"github.com/pkg/errors"
	"github.com/pulumi/pulumi-aws/sdk/v7/go/aws"
	"github.com/pulumi/pulumi-aws/sdk/v7/go/aws/ebs"
	"github.com/pulumi/pulumi-aws/sdk/v7/go/aws/ec2"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

// volume creates the EBS volume (create XOR copy arm), its in-line
// attachments, and exports outputs.
//
// Lifecycle facts the render below depends on:
//   - the two arms are two provider resources (ebs.Volume /
//     ebs.VolumeCopy); exactly one exists per the spec's union CEL,
//     and both expose the same downstream surface (id/arn);
//   - a copy inherits the source's availability zone, encryption
//     posture, and snapshot lineage - the provider offers no
//     override, which is why the spec forbids those fields on the
//     copy arm;
//   - Size/Type/Iops/Throughput update in place on both arms;
//     everything else replaces;
//   - attachments are ForceNew per (volume, instance, device);
//   - FinalSnapshot is config-only at AWS (never read back), so
//     imports do not round-trip it.
func volume(ctx *pulumi.Context, locals *Locals, provider *aws.Provider) error {
	spec := locals.Spec

	var volumeId pulumi.StringOutput
	var volumeArn pulumi.StringOutput
	var volumeZone pulumi.StringOutput
	var volumeSize pulumi.IntOutput
	var createTime pulumi.StringOutput

	if spec.CopyFrom != nil {
		// The COPY arm: clone an existing volume (same zone as the
		// source).
		args := &ebs.VolumeCopyArgs{
			SourceVolumeId: pulumi.String(spec.CopyFrom.SourceVolumeId.GetValue()),
			Tags:           pulumi.ToStringMap(locals.AwsTags),
		}
		if spec.SizeGb > 0 {
			args.Size = pulumi.Int(int(spec.SizeGb))
		}
		if spec.Type != "" {
			args.VolumeType = pulumi.String(spec.Type)
		}
		if spec.Iops > 0 {
			args.Iops = pulumi.Int(int(spec.Iops))
		}
		if spec.ThroughputMibps > 0 {
			args.Throughput = pulumi.Int(int(spec.ThroughputMibps))
		}

		createdCopy, err := ebs.NewVolumeCopy(ctx, "volume", args, pulumi.Provider(provider))
		if err != nil {
			return errors.Wrap(err, "create volume copy")
		}
		volumeId = createdCopy.ID().ToStringOutput()
		volumeArn = createdCopy.Arn
		volumeZone = createdCopy.AvailabilityZone
		volumeSize = createdCopy.Size
		// The provider does not expose create time on the copy arm.
		createTime = pulumi.String("").ToStringOutput()
	} else {
		// The CREATE arm: a fresh volume in a chosen availability
		// zone.
		args := &ebs.VolumeArgs{
			AvailabilityZone: pulumi.String(spec.AvailabilityZone),
			Tags:             pulumi.ToStringMap(locals.AwsTags),
		}
		if spec.SizeGb > 0 {
			args.Size = pulumi.Int(int(spec.SizeGb))
		}
		if spec.Type != "" {
			args.Type = pulumi.String(spec.Type)
		}
		if spec.Iops > 0 {
			args.Iops = pulumi.Int(int(spec.Iops))
		}
		if spec.ThroughputMibps > 0 {
			args.Throughput = pulumi.Int(int(spec.ThroughputMibps))
		}
		if spec.SnapshotId.GetValue() != "" {
			args.SnapshotId = pulumi.String(spec.SnapshotId.GetValue())
		}
		if spec.Encrypted {
			args.Encrypted = pulumi.Bool(true)
		}
		if spec.KmsKeyId.GetValue() != "" {
			args.KmsKeyId = pulumi.String(spec.KmsKeyId.GetValue())
		}
		if spec.MultiAttachEnabled {
			args.MultiAttachEnabled = pulumi.Bool(true)
		}
		if spec.FinalSnapshot {
			args.FinalSnapshot = pulumi.Bool(true)
		}
		if spec.VolumeInitializationRate > 0 {
			args.VolumeInitializationRate = pulumi.Int(int(spec.VolumeInitializationRate))
		}

		createdVolume, err := ebs.NewVolume(ctx, "volume", args, pulumi.Provider(provider))
		if err != nil {
			return errors.Wrap(err, "create volume")
		}
		volumeId = createdVolume.ID().ToStringOutput()
		volumeArn = createdVolume.Arn
		volumeZone = createdVolume.AvailabilityZone
		volumeSize = createdVolume.Size
		createTime = createdVolume.CreateTime
	}

	// In-line attachments, one resource per (device, instance) pair -
	// the same keying as the Terraform module.
	for _, attachment := range spec.Attachments {
		args := &ec2.VolumeAttachmentArgs{
			DeviceName: pulumi.String(attachment.DeviceName),
			VolumeId:   volumeId,
			InstanceId: pulumi.String(attachment.InstanceId.GetValue()),
		}
		if attachment.ForceDetach {
			args.ForceDetach = pulumi.Bool(true)
		}
		if attachment.SkipDestroy {
			args.SkipDestroy = pulumi.Bool(true)
		}
		if attachment.StopInstanceBeforeDetaching {
			args.StopInstanceBeforeDetaching = pulumi.Bool(true)
		}

		resourceName := fmt.Sprintf("attachment-%s-%s", attachment.DeviceName, attachment.InstanceId.GetValue())
		if _, err := ec2.NewVolumeAttachment(ctx, resourceName, args, pulumi.Provider(provider)); err != nil {
			return errors.Wrapf(err, "attach %s", attachment.DeviceName)
		}
	}

	ctx.Export(OpVolumeId, volumeId)
	ctx.Export(OpVolumeArn, volumeArn)
	ctx.Export(OpAvailabilityZone, volumeZone)
	// The size is an int at the provider; exported as a string to
	// match the outputs contract (string-typed observable state).
	ctx.Export(OpSizeGb, pulumi.Sprintf("%d", volumeSize))
	ctx.Export(OpCreateTime, createTime)
	return nil
}
