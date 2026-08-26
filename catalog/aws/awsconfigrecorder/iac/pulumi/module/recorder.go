package module

import (
	"github.com/pkg/errors"
	"github.com/pulumi/pulumi-aws/sdk/v7/go/aws"
	"github.com/pulumi/pulumi-aws/sdk/v7/go/aws/cfg"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

// recorder creates the region's configuration recorder, delivery
// channel, recorder status, and retention configuration, and exports
// outputs.
//
// Lifecycle facts the renders below depend on:
//   - AWS allows ONE recorder and ONE delivery channel per region,
//     both named "default" by convention -- the names are hardcoded
//     here and metadata.name never reaches AWS (the settings-singleton
//     contract);
//   - the delivery channel cannot be created before a recorder exists,
//     and cannot be DELETED while the recorder is running -- the
//     dependency chain below fixes create order, and the provider
//     retries channel deletion for 30s while the recorder stop lands;
//   - the recorder-status resource is the folded start/stop toggle:
//     Create/Update/Delete all map to Start/StopConfigurationRecorder;
//   - starting a recorder without a delivery channel fails
//     (NoAvailableDeliveryChannelException) -- the spec CEL guarantees
//     a channel arrives whenever recording_enabled is not false;
//   - AWS validates the delivery bucket's POLICY ("insufficient
//     delivery policy") -- the bucket policy is the consumer's
//     contract (AwsS3Bucket spec.policy), never this module's.
func recorder(ctx *pulumi.Context, locals *Locals, provider *aws.Provider) error {
	spec := locals.Spec

	args := &cfg.RecorderArgs{
		// The regional singleton's conventional name -- an AWS-side
		// constant, deliberately not configurable.
		Name:    pulumi.String("default"),
		RoleArn: pulumi.String(spec.RoleArn.GetValue()),
	}

	if spec.RecordingGroup != nil {
		g := spec.RecordingGroup
		groupArgs := &cfg.RecorderRecordingGroupArgs{}
		if g.AllSupported != nil {
			groupArgs.AllSupported = pulumi.Bool(*g.AllSupported)
		}
		if g.IncludeGlobalResourceTypes != nil {
			groupArgs.IncludeGlobalResourceTypes = pulumi.Bool(*g.IncludeGlobalResourceTypes)
		}
		if len(g.ResourceTypes) > 0 {
			groupArgs.ResourceTypes = pulumi.ToStringArray(g.ResourceTypes)
		}
		if len(g.ExclusionByResourceTypes) > 0 {
			groupArgs.ExclusionByResourceTypes = cfg.RecorderRecordingGroupExclusionByResourceTypeArray{
				&cfg.RecorderRecordingGroupExclusionByResourceTypeArgs{
					ResourceTypes: pulumi.ToStringArray(g.ExclusionByResourceTypes),
				},
			}
		}
		if g.RecordingStrategy != "" {
			groupArgs.RecordingStrategies = cfg.RecorderRecordingGroupRecordingStrategyArray{
				&cfg.RecorderRecordingGroupRecordingStrategyArgs{
					UseOnly: pulumi.String(g.RecordingStrategy),
				},
			}
		}
		args.RecordingGroup = groupArgs
	}

	if spec.RecordingMode != nil {
		m := spec.RecordingMode
		modeArgs := &cfg.RecorderRecordingModeArgs{}
		if m.RecordingFrequency != "" {
			modeArgs.RecordingFrequency = pulumi.String(m.RecordingFrequency)
		} else {
			modeArgs.RecordingFrequency = pulumi.String("CONTINUOUS")
		}
		if m.Override != nil {
			overrideArgs := &cfg.RecorderRecordingModeRecordingModeOverrideArgs{
				RecordingFrequency: pulumi.String(m.Override.RecordingFrequency),
				ResourceTypes:      pulumi.ToStringArray(m.Override.ResourceTypes),
			}
			if m.Override.Description != "" {
				overrideArgs.Description = pulumi.String(m.Override.Description)
			}
			modeArgs.RecordingModeOverride = overrideArgs
		}
		args.RecordingMode = modeArgs
	}

	createdRecorder, err := cfg.NewRecorder(ctx, "recorder", args, pulumi.Provider(provider))
	if err != nil {
		return errors.Wrap(err, "create configuration recorder")
	}

	// The delivery channel (AWS refuses a channel without a recorder).
	var createdChannel *cfg.DeliveryChannel
	if spec.DeliveryChannel != nil {
		c := spec.DeliveryChannel
		channelArgs := &cfg.DeliveryChannelArgs{
			// The regional singleton's conventional name (see the
			// recorder).
			Name:         pulumi.String("default"),
			S3BucketName: pulumi.String(c.S3BucketName.GetValue()),
		}
		if c.S3KeyPrefix != "" {
			channelArgs.S3KeyPrefix = pulumi.String(c.S3KeyPrefix)
		}
		if c.S3KmsKeyArn.GetValue() != "" {
			channelArgs.S3KmsKeyArn = pulumi.String(c.S3KmsKeyArn.GetValue())
		}
		if c.SnsTopicArn.GetValue() != "" {
			channelArgs.SnsTopicArn = pulumi.String(c.SnsTopicArn.GetValue())
		}
		if c.SnapshotDeliveryFrequency != "" {
			channelArgs.SnapshotDeliveryProperties = &cfg.DeliveryChannelSnapshotDeliveryPropertiesArgs{
				DeliveryFrequency: pulumi.String(c.SnapshotDeliveryFrequency),
			}
		}
		createdChannel, err = cfg.NewDeliveryChannel(ctx, "delivery-channel", channelArgs,
			pulumi.Provider(provider), pulumi.DependsOn([]pulumi.Resource{createdRecorder}))
		if err != nil {
			return errors.Wrap(err, "create delivery channel")
		}
	}

	// The folded recorder toggle: unset recording_enabled means
	// RUNNING (the reason this component exists). Starting requires
	// the delivery channel to exist.
	recordingEnabled := true
	if spec.RecordingEnabled != nil {
		recordingEnabled = *spec.RecordingEnabled
	}
	statusDeps := []pulumi.Resource{createdRecorder}
	if createdChannel != nil {
		statusDeps = append(statusDeps, createdChannel)
	}
	createdStatus, err := cfg.NewRecorderStatus(ctx, "recorder-status", &cfg.RecorderStatusArgs{
		Name:      createdRecorder.Name,
		IsEnabled: pulumi.Bool(recordingEnabled),
	}, pulumi.Provider(provider), pulumi.DependsOn(statusDeps))
	if err != nil {
		return errors.Wrap(err, "set recorder status")
	}

	// The retention singleton (AWS names it "default"; the name is
	// API-computed and cannot be chosen). Managed only when the spec
	// sets a window.
	if spec.RetentionPeriodInDays > 0 {
		_, err := cfg.NewRetentionConfiguration(ctx, "retention", &cfg.RetentionConfigurationArgs{
			RetentionPeriodInDays: pulumi.Int(int(spec.RetentionPeriodInDays)),
		}, pulumi.Provider(provider))
		if err != nil {
			return errors.Wrap(err, "set retention configuration")
		}
	}

	ctx.Export(OpRecorderName, createdRecorder.Name)
	if createdChannel != nil {
		ctx.Export(OpDeliveryChannelName, createdChannel.Name)
	} else {
		ctx.Export(OpDeliveryChannelName, pulumi.String(""))
	}
	ctx.Export(OpRecordingEnabled, createdStatus.IsEnabled)
	// Config's singletons are addressed by REGION + the literal name
	// "default"; consumers (and the harness verifier) reaching the
	// recorder off the ambient region need the resolved region
	// alongside recorder_name.
	ctx.Export(OpRegion, createdRecorder.Region)
	return nil
}
