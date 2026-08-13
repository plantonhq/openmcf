package module

import (
	"strconv"

	"github.com/pkg/errors"
	"github.com/pulumi/pulumi-aws/sdk/v7/go/aws"
	"github.com/pulumi/pulumi-aws/sdk/v7/go/aws/fsx"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

func volume(ctx *pulumi.Context, locals *Locals, provider *aws.Provider) (*fsx.OntapVolume, error) {
	spec := locals.AwsFsxOntapVolume.Spec
	name := locals.AwsFsxOntapVolume.Metadata.Name

	args := &fsx.OntapVolumeArgs{
		StorageVirtualMachineId: pulumi.String(spec.StorageVirtualMachineId.GetValue()),
		Name:                    pulumi.StringPtr(spec.Name),
		Tags:                    pulumi.ToStringMap(locals.AwsTags),
	}

	// Size flows through exactly one arm (spec-enforced XOR). SizeInBytes is
	// the byte-precise arm and the only way past 2 PiB — the provider types
	// it as a string because the value exceeds a 32-bit integer.
	if spec.SizeInMegabytes != nil {
		args.SizeInMegabytes = pulumi.IntPtr(int(spec.GetSizeInMegabytes()))
	}
	if spec.SizeInBytes != nil {
		args.SizeInBytes = pulumi.StringPtr(strconv.FormatInt(spec.GetSizeInBytes(), 10))
	}

	// Junction path (optional — unmounted volume if omitted).
	if spec.JunctionPath != "" {
		args.JunctionPath = pulumi.StringPtr(spec.JunctionPath)
	}

	// ONTAP volume type (ForceNew; spec default RW, materialized before the
	// module runs).
	if spec.GetOntapVolumeType() != "" {
		args.OntapVolumeType = pulumi.StringPtr(spec.GetOntapVolumeType())
	}

	// Volume style (ForceNew; spec default FLEXVOL, materialized before the
	// module runs).
	if spec.GetVolumeStyle() != "" {
		args.VolumeStyle = pulumi.StringPtr(spec.GetVolumeStyle())
	}

	// Security style (optional — inherits from SVM if omitted).
	if spec.SecurityStyle != "" {
		args.SecurityStyle = pulumi.StringPtr(spec.SecurityStyle)
	}

	// Snapshot policy (optional — ONTAP's "default" policy if omitted).
	if spec.SnapshotPolicy != "" {
		args.SnapshotPolicy = pulumi.StringPtr(spec.SnapshotPolicy)
	}

	// Tri-state efficiency switch: the resolved value is sent whenever
	// present so an explicit false genuinely disables dedup/compression —
	// an only-send-true guard would silently drop the disable.
	if spec.StorageEfficiencyEnabled != nil {
		args.StorageEfficiencyEnabled = pulumi.BoolPtr(spec.GetStorageEfficiencyEnabled())
	}

	// Backup and delete-time controls are presence-honest: explicit false
	// values must reach AWS rather than being coerced away. The delete-time
	// trio (skip final backup, its tags, the SnapLock bypass) is read from
	// state at destroy time, so it must be applied BEFORE the deletion.
	if spec.CopyTagsToBackups != nil {
		args.CopyTagsToBackups = pulumi.BoolPtr(spec.GetCopyTagsToBackups())
	}
	if spec.SkipFinalBackup != nil {
		args.SkipFinalBackup = pulumi.BoolPtr(spec.GetSkipFinalBackup())
	}
	if len(spec.FinalBackupTags) > 0 {
		args.FinalBackupTags = pulumi.ToStringMap(spec.FinalBackupTags)
	}
	if spec.BypassSnaplockEnterpriseRetention != nil {
		args.BypassSnaplockEnterpriseRetention = pulumi.BoolPtr(spec.GetBypassSnaplockEnterpriseRetention())
	}

	// Tiering policy (optional).
	if spec.TieringPolicy != nil && spec.TieringPolicy.Name != "" {
		tp := &fsx.OntapVolumeTieringPolicyArgs{
			Name: pulumi.StringPtr(spec.TieringPolicy.Name),
		}
		if spec.TieringPolicy.CoolingPeriod > 0 {
			tp.CoolingPeriod = pulumi.IntPtr(int(spec.TieringPolicy.CoolingPeriod))
		}
		args.TieringPolicy = tp
	}

	// SnapLock configuration (optional).
	if spec.SnaplockConfiguration != nil {
		sl := spec.SnaplockConfiguration

		slArgs := &fsx.OntapVolumeSnaplockConfigurationArgs{
			SnaplockType: pulumi.String(sl.SnaplockType),
		}

		// Presence-honest booleans: explicit false must reach AWS.
		if sl.AuditLogVolume != nil {
			slArgs.AuditLogVolume = pulumi.BoolPtr(sl.GetAuditLogVolume())
		}

		if sl.GetPrivilegedDelete() != "" {
			slArgs.PrivilegedDelete = pulumi.StringPtr(sl.GetPrivilegedDelete())
		}

		if sl.VolumeAppendModeEnabled != nil {
			slArgs.VolumeAppendModeEnabled = pulumi.BoolPtr(sl.GetVolumeAppendModeEnabled())
		}

		// Autocommit period.
		if sl.AutocommitPeriod != nil && sl.AutocommitPeriod.Type != "" {
			acArgs := &fsx.OntapVolumeSnaplockConfigurationAutocommitPeriodArgs{
				Type: pulumi.StringPtr(sl.AutocommitPeriod.Type),
			}
			if sl.AutocommitPeriod.Value > 0 {
				acArgs.Value = pulumi.IntPtr(int(sl.AutocommitPeriod.Value))
			}
			slArgs.AutocommitPeriod = acArgs
		}

		// Retention bounds (default, minimum, maximum). A value of 0 IS
		// meaningful for unit types (e.g. a 0-day minimum retention), so the
		// value is sent whenever the type is a unit type and only elided for
		// INFINITE/UNSPECIFIED, where AWS takes no value. The container is
		// attached only when at least one duration survives its type gate —
		// an empty retention_period would send an empty struct to the AWS API.
		if sl.RetentionPeriod != nil {
			retentionValue := func(retentionType string, value int32) pulumi.IntPtrInput {
				if retentionType == "INFINITE" || retentionType == "UNSPECIFIED" || retentionType == "" {
					return nil
				}
				return pulumi.IntPtr(int(value))
			}

			rp := &fsx.OntapVolumeSnaplockConfigurationRetentionPeriodArgs{}
			hasRetention := false

			if dr := sl.RetentionPeriod.DefaultRetention; dr != nil && dr.Type != "" {
				rp.DefaultRetention = &fsx.OntapVolumeSnaplockConfigurationRetentionPeriodDefaultRetentionArgs{
					Type:  pulumi.StringPtr(dr.Type),
					Value: retentionValue(dr.Type, dr.Value),
				}
				hasRetention = true
			}

			if mn := sl.RetentionPeriod.MinimumRetention; mn != nil && mn.Type != "" {
				rp.MinimumRetention = &fsx.OntapVolumeSnaplockConfigurationRetentionPeriodMinimumRetentionArgs{
					Type:  pulumi.StringPtr(mn.Type),
					Value: retentionValue(mn.Type, mn.Value),
				}
				hasRetention = true
			}

			if mx := sl.RetentionPeriod.MaximumRetention; mx != nil && mx.Type != "" {
				rp.MaximumRetention = &fsx.OntapVolumeSnaplockConfigurationRetentionPeriodMaximumRetentionArgs{
					Type:  pulumi.StringPtr(mx.Type),
					Value: retentionValue(mx.Type, mx.Value),
				}
				hasRetention = true
			}

			if hasRetention {
				slArgs.RetentionPeriod = rp
			}
		}

		args.SnaplockConfiguration = slArgs
	}

	// Aggregate configuration (for FLEXGROUP volumes). The block renders when
	// ANY leaf is set — constituents_per_aggregate WITHOUT aggregates is a
	// legitimate shape (AWS spreads constituents across all of the file
	// system's aggregates and pins only the constituent count). Gating on a
	// non-empty aggregates list alone would silently drop that shape.
	if ac := spec.AggregateConfiguration; ac != nil && (len(ac.Aggregates) > 0 || ac.ConstituentsPerAggregate > 0) {
		aggrArgs := &fsx.OntapVolumeAggregateConfigurationArgs{}

		if len(ac.Aggregates) > 0 {
			aggrs := make(pulumi.StringArray, 0, len(ac.Aggregates))
			for _, a := range ac.Aggregates {
				aggrs = append(aggrs, pulumi.String(a))
			}
			aggrArgs.Aggregates = aggrs
		}

		if ac.ConstituentsPerAggregate > 0 {
			aggrArgs.ConstituentsPerAggregate = pulumi.IntPtr(int(ac.ConstituentsPerAggregate))
		}

		args.AggregateConfiguration = aggrArgs
	}

	createdVolume, err := fsx.NewOntapVolume(ctx, name, args, pulumi.Provider(provider))
	if err != nil {
		return nil, errors.Wrap(err, "failed to create fsx ontap volume")
	}

	return createdVolume, nil
}
