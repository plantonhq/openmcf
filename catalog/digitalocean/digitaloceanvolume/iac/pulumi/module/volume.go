package module

import (
	"github.com/pkg/errors"
	"github.com/pulumi/pulumi-digitalocean/sdk/v4/go/digitalocean"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

// volume provisions the block storage volume and exports its outputs.
//
// Attachment to Droplets is a property of the Droplet (its volume_ids list),
// never of the volume. Size can only be EXPANDED after creation -- the
// provider rejects a shrink at plan time. Every other argument is create-only
// and replaces the volume when changed.
func volume(
	ctx *pulumi.Context,
	locals *Locals,
	digitalOceanProvider *digitalocean.Provider,
) (*digitalocean.Volume, error) {
	spec := locals.DigitalOceanVolume.Spec

	volumeArgs := &digitalocean.VolumeArgs{
		Name:   pulumi.String(spec.VolumeName),
		Region: pulumi.String(spec.Region.String()),
		Size:   pulumi.Int(int(spec.SizeGib)),
	}

	// Create-only at the current provider pin: a description change REPLACES
	// the volume.
	if spec.Description != "" {
		volumeArgs.Description = pulumi.StringPtr(spec.Description)
	}

	// When set, the volume is created from the snapshot, inheriting its
	// region and minimum size. Create-only, never reported back by the API.
	if spec.SnapshotId != "" {
		volumeArgs.SnapshotId = pulumi.StringPtr(spec.SnapshotId)
	}

	// Formatting happens once at creation via the NON-deprecated
	// initial_filesystem_type argument (the SDK's FilesystemType maps to the
	// deprecated attribute and conflicts with it). The enum value names ARE
	// the strings the provider expects; unformatted stays unset.
	if fsType := spec.FilesystemType.String(); fsType != "unformatted" {
		volumeArgs.InitialFilesystemType = pulumi.StringPtr(fsType)
		// The label only means anything when a filesystem is being formatted.
		if spec.InitialFilesystemLabel != "" {
			volumeArgs.InitialFilesystemLabel = pulumi.StringPtr(spec.InitialFilesystemLabel)
		}
	}

	// User tags plus the standard Planton labels rendered as "key:value"
	// tags — the exact set the Terraform module applies.
	tagSet := map[string]bool{}
	var tagInputs pulumi.StringArray
	for _, t := range spec.Tags {
		if !tagSet[t] {
			tagSet[t] = true
			tagInputs = append(tagInputs, pulumi.String(t))
		}
	}
	for k, v := range locals.DigitalOceanLabels {
		t := k + ":" + v
		if !tagSet[t] {
			tagSet[t] = true
			tagInputs = append(tagInputs, pulumi.String(t))
		}
	}
	volumeArgs.Tags = tagInputs

	createdVolume, err := digitalocean.NewVolume(
		ctx,
		"volume",
		volumeArgs,
		pulumi.Provider(digitalOceanProvider),
	)
	if err != nil {
		return nil, errors.Wrap(err, "failed to create digitalocean volume")
	}

	// Stack outputs -- exactly the DigitalOceanVolumeStackOutputs contract,
	// from the SDK's real field names (the urn output is VolumeUrn).
	ctx.Export(OpVolumeId, createdVolume.ID())
	ctx.Export(OpUrn, createdVolume.VolumeUrn)

	return createdVolume, nil
}
